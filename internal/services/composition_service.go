package services

import (
	api "github.com/okdp/okdp-server/api/openapi/v3/_api"
	"github.com/okdp/okdp-server/internal/repository"
)

type ICompositionService interface {
	List() []*api.Composition
	Get(string) *api.Composition
	Update(string, api.Composition) *api.Composition
	Create(api.Composition)
	Delete(string) *api.Composition
}

type CompositionService struct {
	compositionRepository *repository.CompositionRepository
}

func NewCompositionService() (*CompositionService, error) {
	compositionRepository, err := repository.NewCompositionRepository()
	if err != nil {
		return nil, err
	}
	return &CompositionService{
		compositionRepository: compositionRepository,
	}, nil
}

func (s CompositionService) Get(id string) *api.Composition {
	return s.compositionRepository.Get(id)
}

func (s CompositionService) List() []*api.Composition {
	return s.compositionRepository.List()
}

func (s CompositionService) Create(composition api.Composition) {
	s.compositionRepository.Create(composition)
}

func (s CompositionService) Delete(id string) *api.Composition {
	return s.compositionRepository.Delete(id)
}

func (s CompositionService) Update(id string, compositionUpdate api.Composition) *api.Composition {
	return s.compositionRepository.Update(id, compositionUpdate)
}
