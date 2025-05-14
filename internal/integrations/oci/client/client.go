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

	log "github.com/okdp/okdp-server/internal/common/logging"
	"github.com/okdp/okdp-server/internal/config"
	"github.com/okdp/okdp-server/internal/model"
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
			for _, p := range catalog.Packages {
				repo, err := remote.NewRepository(strings.TrimSuffix(catalog.RepoURL, "/") + "/" + p.Name)
				if err != nil {
					log.Fatal("Failed to create a client to the remote repository %s: %v", catalog.RepoURL, err)
				}

				creds, err := getOCIRepoCredentials(catalog)

				if err != nil {
					log.Fatal("Unable to get login and password from dockerconfigjson for the repo: %s/%s", catalog.ID, catalog.RepoURL)
				}

				repo.Client = &auth.Client{
					Client:     retry.DefaultClient,
					Cache:      auth.NewCache(),
					Credential: creds,
				}
				clients[utils.MapKey(catalog.ID, p.Name)] = &RepositoryClient{repo}
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

func GetCatalog(catalogID string) (*model.Catalog, *model.ServerResponse) {
	catalogs := config.GetAppConfig().Catalogs
	for _, catalog := range catalogs {
		if strings.EqualFold(catalog.ID, catalogID) {
			return catalog, nil
		}
	}
	return nil, model.CatalogNotFoundError(catalogID)
}

func GetPackages(catalogID string) ([]*model.Package, *model.ServerResponse) {
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

func GetPackage(catalogID string, name string) (*model.Package, *model.ServerResponse) {
	catalog, err := GetCatalog(catalogID)
	if err != nil {
		return nil, err
	}
	for _, p := range catalog.Packages {
		if strings.EqualFold(p.Name, name) {
			return getPackage(catalogID, p.Name)
		}
	}
	return nil, model.CatalogPackageNotFoundError(catalogID, name)
}

func GetPackageDefinition(catalogID string, name string, version string) (map[string]interface{}, *model.ServerResponse) {
	repo, err := getRepoClient(catalogID, name)
	if err != nil {
		return nil, err
	}
	return repo.fetchDefinition(version)
}

func getPackage(catalogID string, name string) (*model.Package, *model.ServerResponse) {
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

func (r RepositoryClient) listTags() ([]string, *model.ServerResponse) {
	ctx := context.Background()
	var allTags []string
	var last string
	var err *model.ServerResponse
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
					err = model.
						NewServerResponse(model.RegistryResponse).GenericError(httpErr.StatusCode, httpErr.Error())
				} else {
					err = model.
						NewServerResponse(model.RegistryResponse).GenericError(statusCode, er.Error())
				}
			}
			break
		}
	}

	return allTags, err
}

func (r RepositoryClient) fetchDefinition(version string) (map[string]interface{}, *model.ServerResponse) {
	ctx := context.Background()
	// Pull to memory
	memStore := memory.New()
	desc, err := oras.Copy(ctx, r.Repository, version, memStore, version, oras.DefaultCopyOptions)
	if err != nil {
		log.Error("Failed to copy definition into the memstore: %v", err)
		return nil, model.
			NewServerResponse(model.OkdpServerResponse).UnprocessableEntity(err.Error())
	}

	// Decode manifest
	manifestReader, err := memStore.Fetch(ctx, desc)
	if err != nil {
		log.Error("Failed to fetch definition from the memstore: %v", err)
		return nil, model.
			NewServerResponse(model.OkdpServerResponse).UnprocessableEntity(err.Error())
	}
	defer manifestReader.Close()

	var manifest ocispec.Manifest
	if err := json.NewDecoder(manifestReader).Decode(&manifest); err != nil {
		log.Error("Failed to decode definition: %v", err)
		return nil, model.
			NewServerResponse(model.OkdpServerResponse).UnprocessableEntity(err.Error())
	}

	configReader, err := memStore.Fetch(ctx, manifest.Config)
	if err != nil {
		log.Error("Failed to fetch the definition content: %v", err)
		return nil, model.
			NewServerResponse(model.OkdpServerResponse).UnprocessableEntity(err.Error())
	}
	defer configReader.Close()

	var definition map[string]interface{}
	if err := json.NewDecoder(configReader).Decode(&definition); err != nil {
		log.Error("Failed to decode the definition content: %v", err)
		return nil, model.
			NewServerResponse(model.OkdpServerResponse).UnprocessableEntity(err.Error())
	}
	return definition, nil
}

func getRepoClient(catalogID string, packageName string) (*RepositoryClient, *model.ServerResponse) {
	instance, found := instance.clients[utils.MapKey(catalogID, packageName)]
	if !found {
		return nil, model.CatalogPackageNotFoundError(catalogID, packageName)
	}
	return instance, nil
}

func getOCIRepoCredentials(catalog *model.Catalog) (auth.CredentialFunc, error) {

	var empty auth.CredentialFunc = func(_ context.Context, _ string) (auth.Credential, error) {
		return auth.EmptyCredential, nil
	}

	if !catalog.IsAuthenticated() {
		return empty, nil
	}

	login := utils.ResolveEnv(*catalog.Credentials.RobotAccountName)
	passwd := utils.ResolveEnv(*catalog.Credentials.RobotAccountToken)
	dockerjson := utils.ResolveEnv(*catalog.Credentials.Dockerconfigjson)

	if login == "" && passwd == "" {
		if dockerjson == "" {
			return empty, nil
		}
		var err error
		login, passwd, err = utils.ToLoginPassword(dockerjson)
		if err != nil {
			return nil, err
		}
	}
	return auth.StaticCredential(catalog.RepoHost(), auth.Credential{
		Username: login,
		Password: passwd,
	}), nil
}
