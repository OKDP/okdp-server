package repository

import (
	api "github.com/okdp/okdp-server/api/openapi/v3/_api"
)

type IDeploymentRepository interface {
	List() []*api.Deployment
	Get(string) *api.Deployment
	Update(string, api.Deployment) *api.Deployment
	Create(api.Deployment)
	Delete(string) *api.Deployment
}

type DeploymentRepository struct {
}

var deployments = []*api.Deployment{
	{
		ToDo1: "1",
		ToDo2: "ToDo2-2",
		ToDo3: "ToDo3-2",
		ToDo4: "ToDo4-2",
	},
}

func NewDeploymentRepository() (*DeploymentRepository, error) {
	return &DeploymentRepository{}, nil
}

func (r *DeploymentRepository) Get(id string) *api.Deployment {
	for _, deployment := range deployments {
		if deployment.ToDo1 == id {
			return deployment
		}
	}
	return nil
}

func (r *DeploymentRepository) List() []*api.Deployment {
	return deployments
}

func (r *DeploymentRepository) Create(deployment api.Deployment) {
	deployments = append(deployments, &deployment)
}

func (r *DeploymentRepository) Delete(id string) *api.Deployment {
	for i, deployment := range deployments {
		if deployment.ToDo1 == id {
			deployments = append(deployments[:i], (deployments)[i+1:]...)
			return &api.Deployment{}
		}
	}
	return nil
}

func (r DeploymentRepository) Update(id string, deploymentUpdate api.Deployment) *api.Deployment {
	for i, deployment := range deployments {
		if deployment.ToDo1 == id {
			deployments[i] = &deploymentUpdate
			return deployment
		}
	}
	return nil
}
