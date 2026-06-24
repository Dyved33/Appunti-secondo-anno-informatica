package category

import (
	"context"
	"gialosai/backend/domain"
)

// Controller rappresenta il layer di controllo/handler per le operazioni su categorie.
// Trasla le richieste in ingresso (DTO) alle chiamate al Service e formatta le risposte.
type Controller struct {
	Service *Service
}

// NewController crea una nuova istanza del controller per le categorie.
func NewController(svc *Service) *Controller {
	return &Controller{Service: svc}
}

// CreateCategory gestisce la creazione di una nuova categoria.
// Accetta i dati di input (DTO), chiama il Service e restituisce i dati della categoria creata.
func (c *Controller) CreateCategory(ctx context.Context, input CategoryInputDTO) (CategoryOutputDTO, error) {
	cat := &domain.Categoria{
		Nome:        input.Nome,
		Descrizione: input.Descrizione,
	}

	// Chiama il Service per creare la categoria
	if err := c.Service.CreateCategory(ctx, cat); err != nil {
		return CategoryOutputDTO{}, err
	}

	// Formatta e restituisce il DTO di output
	return CategoryOutputDTO{ID: cat.ID, Nome: cat.Nome, Descrizione: cat.Descrizione}, nil
}

// GetCategories recupera tutte le categorie disponibili nel sistema.
// Restituisce una slice di DTO con i dati delle categorie.
func (c *Controller) GetCategories(ctx context.Context) ([]CategoryOutputDTO, error) {
	categories, err := c.Service.GetCategories(ctx)
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

// GetCategoryByID recupera una singola categoria dal sistema.
// Restituisce i dati della categoria in formato DTO.
func (c *Controller) GetCategoryByID(ctx context.Context, id int) (CategoryOutputDTO, error) {
	cat, err := c.Service.GetCategoryByID(ctx, id)
	if err != nil {
		return CategoryOutputDTO{}, err
	}

	return CategoryOutputDTO{ID: cat.ID, Nome: cat.Nome, Descrizione: cat.Descrizione}, nil
}

// UpdateCategory modifica una categoria esistente nel sistema.
// Accetta l'ID della categoria e i dati di input (DTO) per l'aggiornamento.
func (c *Controller) UpdateCategory(ctx context.Context, id int, input CategoryInputDTO) error {
	return c.Service.UpdateCategory(ctx, &domain.Categoria{
		ID:          id,
		Nome:        input.Nome,
		Descrizione: input.Descrizione,
	})
}

// DeleteCategory rimuove una categoria dal sistema.
// Accetta l'ID della categoria da eliminare.
func (c *Controller) DeleteCategory(ctx context.Context, id int) error {
	return c.Service.DeleteCategory(ctx, id)
}
