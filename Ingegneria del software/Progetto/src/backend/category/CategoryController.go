package category

import (
	"context"
	"github.com/tommasogarofani/ingsw_25-26_group-11/src/backend/domain"
)

type categoryService interface {
	CategoryService_Create(ctx context.Context, c *domain.Categoria) error
	CategoryService_GetAll(ctx context.Context) ([]domain.Categoria, error)
	CategoryService_GetByID(ctx context.Context, id int) (*domain.Categoria, error)
	CategoryService_Update(ctx context.Context, c *domain.Categoria) error
	CategoryService_Delete(ctx context.Context, id int) error
}

type CategoryController struct {
	Service categoryService
}

func NewCategoryController(svc categoryService) *CategoryController {
	return &CategoryController{Service: svc}
}

// CreateCategory gestisce la creazione di una nuova categoria
func (c *CategoryController) CategoryController_Create(ctx context.Context, input CategoryInputDTO) (CategoryOutputDTO, error) {
	cat := &domain.Categoria{
		Nome:        input.Nome,
		Descrizione: input.Descrizione,
	}

	// Chiama il Service per creare la categoria
	if err := c.Service.CategoryService_Create(ctx, cat); err != nil {
		return CategoryOutputDTO{}, err
	}

	// Formatta e restituisce il DTO di output
	return CategoryOutputDTO{ID: cat.ID, Nome: cat.Nome, Descrizione: cat.Descrizione}, nil
}

// GetCategories recupera tutte le categorie disponibili nel sistema, restituisce una slice di DTO con i dati delle categorie.
func (c *CategoryController) CategoryController_GetAll(ctx context.Context) ([]CategoryOutputDTO, error) {
	categories, err := c.Service.CategoryService_GetAll(ctx)
	if err != nil {
		return nil, err
	}

	// Trasforma i domain objects in DTO per la trasmissione al client
	result := make([]CategoryOutputDTO, 0, len(categories))
	for _, cat := range categories {
		result = append(result, CategoryOutputDTO{ID: cat.ID, Nome: cat.Nome, Descrizione: cat.Descrizione})
	}
	return result, nil
}

// GetCategoryByID recupera una singola categoria dal sistema, restituisce i dati della categoria in formato DTO.
func (c *CategoryController) CategoryController_GetByID(ctx context.Context, id int) (CategoryOutputDTO, error) {
	cat, err := c.Service.CategoryService_GetByID(ctx, id)
	if err != nil {
		return CategoryOutputDTO{}, err
	}

	return CategoryOutputDTO{ID: cat.ID, Nome: cat.Nome, Descrizione: cat.Descrizione}, nil
}

// UpdateCategory modifica una categoria esistente nel sistema, accetta l'ID della categoria e i dati di input (DTO) per l'aggiornamento.
func (c *CategoryController) CategoryController_Update(ctx context.Context, id int, input CategoryInputDTO) error {
	return c.Service.CategoryService_Update(ctx, &domain.Categoria{
		ID:          id,
		Nome:        input.Nome,
		Descrizione: input.Descrizione,
	})
}

// DeleteCategory rimuove una categoria dal sistema, accetta l'ID della categoria da eliminare.
func (c *CategoryController) CategoryController_Delete(ctx context.Context, id int) error {
	return c.Service.CategoryService_Delete(ctx, id)
}
