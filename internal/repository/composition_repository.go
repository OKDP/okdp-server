package repository

import (
	api "github.com/okdp/okdp-server/api/openapi/v3/_api"
)

type ICompositionRepository interface {
	List() []*api.Composition
	Get(string) *api.Composition
	Update(string, api.Composition) *api.Composition
	Create(api.Composition)
	Delete(string) *api.Composition
}

type CompositionRepository struct {
}

var compositions = []*api.Composition{
	{
		ToDo1: "1",
		ToDo2: "ToDo2-1",
		ToDo3: "ToDo3-1",
		ToDo4: "ToDo4-1",
	},
}

func NewCompositionRepository() (*CompositionRepository, error) {
	return &CompositionRepository{}, nil
}

func (r *CompositionRepository) Get(id string) *api.Composition {
	for _, composition := range compositions {
		if composition.ToDo1 == id {
			return composition
		}
	}
	return nil
}

func (r *CompositionRepository) List() []*api.Composition {
	return compositions
}

func (r *CompositionRepository) Create(composition api.Composition) {
	compositions = append(compositions, &composition)
}

func (r *CompositionRepository) Delete(id string) *api.Composition {
	for i, composition := range compositions {
		if composition.ToDo1 == id {
			compositions = append(compositions[:i], (compositions)[i+1:]...)
			return &api.Composition{}
		}
	}
	return nil
}

func (r *CompositionRepository) Update(id string, compositionUpdate api.Composition) *api.Composition {
	for i, composition := range compositions {
		if composition.ToDo1 == id {
			compositions[i] = &compositionUpdate
			return composition
		}
	}
	return nil
}

