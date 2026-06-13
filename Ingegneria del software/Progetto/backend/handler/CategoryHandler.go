package Handler

import (
	"context"

	"gialosai/backend/DTO"
	"gialosai/backend/domain"
	"gialosai/backend/service"
)

type CategoryHandler struct {
	Service *Service.CategoryService
}

func NewCategoryHandler(svc *Service.CategoryService) *CategoryHandler {
	return &CategoryHandler{Service: svc}
}

func (h *CategoryHandler) CreateCategory(ctx context.Context, input DTO.CategoryInputDTO) (DTO.CategoryOutputDTO, error) {
	c := &domain.Categoria{
		Nome:        input.Nome,
		Descrizione: input.Descrizione,
	}

	if err := h.Service.CreateCategory(ctx, c); err != nil {
		return DTO.CategoryOutputDTO{}, err
	}

	return DTO.CategoryOutputDTO{ID: c.ID, Nome: c.Nome, Descrizione: c.Descrizione}, nil
}

func (h *CategoryHandler) GetCategories(ctx context.Context) ([]DTO.CategoryOutputDTO, error) {
	categories, err := h.Service.GetCategories(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]DTO.CategoryOutputDTO, 0, len(categories))
	for _, c := range categories {
		result = append(result, DTO.CategoryOutputDTO{ID: c.ID, Nome: c.Nome, Descrizione: c.Descrizione})
	}
	return result, nil
}

func (h *CategoryHandler) GetCategoryByID(ctx context.Context, id int) (DTO.CategoryOutputDTO, error) {
	c, err := h.Service.GetCategoryByID(ctx, id)
	if err != nil {
		return DTO.CategoryOutputDTO{}, err
	}

	return DTO.CategoryOutputDTO{ID: c.ID, Nome: c.Nome, Descrizione: c.Descrizione}, nil
}

func (h *CategoryHandler) UpdateCategory(ctx context.Context, id int, input DTO.CategoryInputDTO) error {
	return h.Service.UpdateCategory(ctx, &domain.Categoria{
		ID:          id,
		Nome:        input.Nome,
		Descrizione: input.Descrizione,
	})
}

func (h *CategoryHandler) DeleteCategory(ctx context.Context, id int) error {
	return h.Service.DeleteCategory(ctx, id)
}