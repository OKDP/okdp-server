/*
 *    Copyright 2025 okdp.io
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package client

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"sync"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content/memory"
	"oras.land/oras-go/v2/registry/remote"
	"oras.land/oras-go/v2/registry/remote/auth"
	"oras.land/oras-go/v2/registry/remote/errcode"
	"oras.land/oras-go/v2/registry/remote/retry"

	"github.com/okdp/okdp-server/internal/config"
	log "github.com/okdp/okdp-server/internal/logging"
	"github.com/okdp/okdp-server/internal/model"
	"github.com/okdp/okdp-server/internal/servererrors"
	"github.com/okdp/okdp-server/internal/utils"
)

var (
	instance *RepositoryClients
	once     sync.Once
)

type RepositoryClients struct {
	clients map[string]*RepositoryClient
}

type RepositoryClient struct {
	*remote.Repository
}

func GetClients() *RepositoryClients {
	once.Do(func() {
		clients := make(map[string]*RepositoryClient)
		catalogs := config.GetAppConfig().Catalogs
		for _, catalog := range catalogs {
			log.Info("Container Registry configuration: %+v", catalog)

			var creds auth.CredentialFunc = func(_ context.Context, _ string) (auth.Credential, error) {
				return auth.EmptyCredential, nil
			}

			if catalog.IsAuthenticated() {
				log.Info("The catalog ID '%s' is authenticated", catalog.ID)
				username := utils.GetEnv(*catalog.Credentials.RobotAccountName)
				password := utils.GetEnv(*catalog.Credentials.RobotAccountToken)
				if username == "" || password == "" {
					log.Fatal("Robot account name and robot account token: %s/%s are empty for catalog %s",
						*catalog.Credentials.RobotAccountName, *catalog.Credentials.RobotAccountToken, catalog.ID)
				}

				creds = auth.StaticCredential(catalog.RepoHost(), auth.Credential{
					Username: username,
					Password: password,
				})
			} else {
				log.Info("The catalog ID '%s' is not authenticated", catalog.ID)
			}

			authClient := &auth.Client{
				Client:     retry.DefaultClient,
				Cache:      auth.NewCache(),
				Credential: creds,
			}

			for _, p := range catalog.Packages {
				repo, err := remote.NewRepository(strings.TrimSuffix(catalog.RepoURL, "/") + "/" + p.Name)
				if err != nil {
					log.Fatal("Failed to create a client to the remote repository %s: %v", catalog.RepoURL, err)
				}
				repo.Client = authClient
				clients[repoClientKey(catalog.ID, p.Name)] = &RepositoryClient{repo}
			}
		}
		instance = &RepositoryClients{clients: clients}
	})
	return instance
}

func ListCatalogs() []*model.Catalog {
	catalogs := config.GetAppConfig().Catalogs
	if catalogs == nil {
		return []*model.Catalog{}
	}
	return catalogs
}

func GetCatalog(catalogID string) (*model.Catalog, *servererrors.ServerError) {
	catalogs := config.GetAppConfig().Catalogs
	for _, catalog := range catalogs {
		if strings.EqualFold(catalog.ID, catalogID) {
			return catalog, nil
		}
	}
	return nil, CatalogNotFoundError(catalogID)
}

func GetPackages(catalogID string) ([]*model.Package, *servererrors.ServerError) {
	catalog, err := GetCatalog(catalogID)
	if err != nil {
		return nil, err
	}
	packages := make([]*model.Package, 0, len(catalog.Packages))
	for _, p := range catalog.Packages {
		result, err := getPackage(catalogID, p.Name)
		if err != nil {
			return nil, err
		}
		packages = append(packages, result)
	}
	return packages, nil
}

func GetPackageByName(catalogID string, name string) (*model.Package, *servererrors.ServerError) {
	catalog, err := GetCatalog(catalogID)
	if err != nil {
		return nil, err
	}
	for _, p := range catalog.Packages {
		if strings.EqualFold(p.Name, name) {
			return getPackage(catalogID, p.Name)
		}
	}
	return nil, CatalogPackageNotFoundError(catalogID, name)
}

func GetPackageDefinition(catalogID string, name string, version string) (map[string]interface{}, *servererrors.ServerError) {
	repo, err := getRepoClient(catalogID, name)
	if err != nil {
		return nil, err
	}
	return repo.fetchDefinition(version)
}

func getPackage(catalogID string, name string) (*model.Package, *servererrors.ServerError) {
	repo, err := getRepoClient(catalogID, name)
	if err != nil {
		return nil, err
	}
	versions, err := repo.listTags()
	if err != nil {
		return nil, err
	}
	return &model.Package{
		Name:     name,
		Versions: utils.SortVersions(versions),
	}, nil
}

func (r RepositoryClient) listTags() ([]string, *servererrors.ServerError) {
	ctx := context.Background()
	var allTags []string
	var last string
	var err *servererrors.ServerError
	for {
		er := r.Tags(ctx, last, func(tags []string) error {
			if len(tags) == 0 {
				return io.EOF
			}
			last = tags[len(tags)-1]
			allTags = append(allTags, tags...)
			return nil
		})

		statusCode := http.StatusBadGateway
		if er != nil {
			if er != io.EOF {
				var httpErr *errcode.ErrorResponse
				if errors.As(er, &httpErr) {
					err = servererrors.OfType(servererrors.Registry).GenericError(httpErr.StatusCode, httpErr.Error())
				} else {
					err = servererrors.OfType(servererrors.Registry).GenericError(statusCode, er.Error())
				}
			}
			break
		}
	}

	return allTags, err
}

func (r RepositoryClient) fetchDefinition(version string) (map[string]interface{}, *servererrors.ServerError) {
	ctx := context.Background()
	// Pull to memory
	memStore := memory.New()
	desc, err := oras.Copy(ctx, r.Repository, version, memStore, version, oras.DefaultCopyOptions)
	if err != nil {
		log.Error("Failed to copy definition into the memstore: %v", err)
		return nil, servererrors.OfType(servererrors.OkdpServer).UnprocessableEntity(err.Error())
	}

	// Decode manifest
	manifestReader, err := memStore.Fetch(ctx, desc)
	if err != nil {
		log.Error("Failed to fetch definition from the memstore: %v", err)
		return nil, servererrors.OfType(servererrors.OkdpServer).UnprocessableEntity(err.Error())
	}
	defer manifestReader.Close()

	var manifest ocispec.Manifest
	if err := json.NewDecoder(manifestReader).Decode(&manifest); err != nil {
		log.Error("Failed to decode definition: %v", err)
		return nil, servererrors.OfType(servererrors.OkdpServer).UnprocessableEntity(err.Error())
	}

	configReader, err := memStore.Fetch(ctx, manifest.Config)
	if err != nil {
		log.Error("Failed to fetch the definition content: %v", err)
		return nil, servererrors.OfType(servererrors.OkdpServer).UnprocessableEntity(err.Error())
	}
	defer configReader.Close()

	var definition map[string]interface{}
	if err := json.NewDecoder(configReader).Decode(&definition); err != nil {
		log.Error("Failed to decode the definition content: %v", err)
		return nil, servererrors.OfType(servererrors.OkdpServer).UnprocessableEntity(err.Error())
	}
	return definition, nil
}

func getRepoClient(catalogID string, packageName string) (*RepositoryClient, *servererrors.ServerError) {
	instance, found := instance.clients[repoClientKey(catalogID, packageName)]
	if !found {
		return nil, CatalogPackageNotFoundError(catalogID, packageName)
	}
	return instance, nil
}

func repoClientKey(catalogID string, packageName string) string {
	return catalogID + packageName
}
