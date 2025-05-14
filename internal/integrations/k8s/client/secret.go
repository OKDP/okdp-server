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
	"os"

	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/okdp/okdp-server/internal/model"
	"github.com/skeema/knownhosts"
	corev1 "k8s.io/api/core/v1"
	k8s "sigs.k8s.io/controller-runtime/pkg/client"

	log "github.com/okdp/okdp-server/internal/common/logging"
)

type K8SSecret struct {
	corev1.Secret
}

func (c KubeClient) GetSecret(ctx context.Context, name string, namespace string) (*K8SSecret, *model.ServerResponse) {
	secretKey := k8s.ObjectKey{
		Namespace: namespace,
		Name:      name,
	}
	var secret corev1.Secret
	err := c.Get(ctx, secretKey, &secret)
	if err != nil {
		return nil, model.
			NewServerResponse(model.K8sClusterResponse).
			UnprocessableEntity("Failed to get Kubernetes secret '%s' in namespace '%s', details: '%s'", name, namespace, err.Error())
	}

	return &K8SSecret{secret}, nil

}

func (c KubeClient) GetAuthMethod(secretName string, namespace string) (transport.AuthMethod, *model.ServerResponse) {
	secret, err := c.GetSecret(context.Background(), secretName, namespace)
	if err != nil {
		return nil, err
	}
	auth, err := secret.ToAuthMethod()
	if err != nil {
		return nil, err
	}
	return auth, nil
}

func (s K8SSecret) ToAuthMethod() (transport.AuthMethod, *model.ServerResponse) {
	switch {
	case s.Data != nil && s.Data["password"] != nil:
		log.Fatal("GIT repository access will use git token")
		return &http.BasicAuth{
			Username: string(s.Data["username"]),
			Password: string(s.Data["password"]),
		}, nil

	case s.Data != nil && s.Data["identity"] != nil:
		return BuildKey(s.Secret)

	default:
		return nil, model.
			NewServerResponse(model.K8sClusterResponse).
			UnprocessableEntity("Invalid secret %s=%s", s.Name, s.Namespace)
	}
}

func BuildKey(secret corev1.Secret) (*ssh.PublicKeys, *model.ServerResponse) {
	key, err := ssh.NewPublicKeys("git", secret.Data["identity"], "")
	if err != nil {
		return nil, model.
			NewServerResponse(model.K8sClusterResponse).
			UnprocessableEntity("Failed to create new public key from secret %s=%s, details: %s", secret.Name, secret.Namespace, err.Error())
	}
	file, err := os.CreateTemp(os.TempDir(), "git_known_hosts")
	if err != nil {
		return nil, model.
			NewServerResponse(model.K8sClusterResponse).
			UnprocessableEntity("Failed to create new tmp dir for git_known_hosts from secret %s=%s, details: %s", secret.Name, secret.Namespace, err.Error())
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Error("Error closing file: %v", err)
		}
		_ = os.Remove(file.Name())
	}()

	_, err = file.Write(secret.Data["known_hosts"])
	if err != nil {
		return nil, model.
			NewServerResponse(model.K8sClusterResponse).
			UnprocessableEntity("Failed to write known_hosts from secret %s=%s, details: %s", secret.Name, secret.Namespace, err.Error())
	}
	db, err := knownhosts.NewDB(file.Name())
	if err != nil {
		return nil, model.
			NewServerResponse(model.K8sClusterResponse).
			UnprocessableEntity("Failed to create a new known hosts database from secret %s=%s, details: %s", secret.Name, secret.Namespace, err.Error())
	}
	key.HostKeyCallback = db.HostKeyCallback()
	return key, nil
}
