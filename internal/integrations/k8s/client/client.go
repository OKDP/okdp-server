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
	"sync"

	kustomizev1 "github.com/fluxcd/kustomize-controller/api/v1"
	sourcev1 "github.com/fluxcd/source-controller/api/v1"
	corev1 "k8s.io/api/core/v1"
	apiruntime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/clientcmd"
	k8s "sigs.k8s.io/controller-runtime/pkg/client"

	kubocdv1alpha1 "kubocd/api/v1alpha1"

	log "github.com/okdp/okdp-server/internal/common/logging"
	"github.com/okdp/okdp-server/internal/config"
	"github.com/okdp/okdp-server/internal/model"
	"github.com/okdp/okdp-server/internal/utils"
)

var (
	instance *KubeClients
	once     sync.Once
)

type KubeClients struct {
	clients map[string]*KubeClient
}

type KubeClient struct {
	k8s.Client
	clusterID string
}

func GetClients() *KubeClients {
	once.Do(func() {
		clients := make(map[string]*KubeClient)
		clusters := config.GetAppConfig().Clusters
		for _, cluster := range clusters {
			log.Info("K8S Cluster configuration: %+v", cluster)
			config, err := clientcmd.BuildConfigFromFlags(cluster.Auth.Kubeconfig.APIServer, cluster.Auth.Kubeconfig.Path)
			if cluster.Auth.Kubeconfig.InsecureSkipTlsVerify {
				log.Warn("TLS verification is disabled for cluster ID: %s (%s).", cluster.ID, cluster.Env)
				config.Insecure = true
				config.CAFile = ""
				config.CAData = nil
			}

			if err != nil {
				log.Fatal("Error building kubeconfig for cluster ID '%s (%s)': %s", cluster.ID, cluster.Env, err.Error())
			}

			kubeClient, err := k8s.New(config, k8s.Options{
				Scheme: newScheme(),
			})

			if err != nil {
				log.Fatal("Error creating new k8s client for cluster ID '%s (%s)': %s", cluster.ID, cluster.Env, err.Error())
			}

			clients[utils.MapKey(cluster.ID)] = &KubeClient{kubeClient, cluster.ID}
		}
		instance = &KubeClients{clients: clients}
	})
	return instance
}

func (c KubeClients) GetClient(clusterID string) (*KubeClient, *model.ServerResponse) {
	client, found := c.clients[clusterID]
	if !found {
		return nil, model.ClusterNotFoundError(clusterID)
	}

	return client, nil
}

func newScheme() *apiruntime.Scheme {
	scheme := apiruntime.NewScheme()
	_ = sourcev1.AddToScheme(scheme)
	_ = kustomizev1.AddToScheme(scheme)
	_ = corev1.AddToScheme(scheme)
	_ = kubocdv1alpha1.AddToScheme(scheme)
	return scheme
}
