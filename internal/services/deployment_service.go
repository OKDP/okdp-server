package services

import (
	api "github.com/okdp/okdp-server/api/openapi/v3/_api"
	"github.com/okdp/okdp-server/internal/repository"
)

type IDeploymentRepository interface {
	List() []*api.Deployment
	Get(string) *api.Deployment
	Update(string, api.Deployment) *api.Deployment
	Create(api.Deployment)
	Delete(string) *api.Deployment
}

type DeploymentService struct {
	deploymentRepository *repository.DeploymentRepository
}

func NewDeploymentService() (*DeploymentService, error) {
	deploymentRepository, err := repository.NewDeploymentRepository()
	if err != nil {
		return nil, err
	}
	return &DeploymentService{
		deploymentRepository: deploymentRepository,
	}, nil
}

func (s DeploymentService) Get(id string) *api.Deployment {
	return s.deploymentRepository.Get(id)
}

func (s DeploymentService) List() []*api.Deployment {
	return s.deploymentRepository.List()
}

func (s DeploymentService) Create(deployment api.Deployment) {
	s.deploymentRepository.Create(deployment)
}

func (s DeploymentService) Delete(id string) *api.Deployment {
	return s.deploymentRepository.Delete(id)
}

func (s DeploymentService) Update(id string, deploymentUpdate api.Deployment) *api.Deployment {
	return s.deploymentRepository.Update(id, deploymentUpdate)
}
