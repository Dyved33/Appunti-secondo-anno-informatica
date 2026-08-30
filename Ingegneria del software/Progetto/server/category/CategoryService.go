package category

import (
	"context"
	"errors"
	"github.com/tommasogarofani/ingsw_25-26_group-11/server/domain"
	"strings"
)

type categoryRepository interface {
	CategoryRepository_Create(ctx context.Context, c *domain.Categoria) error
	CategoryRepository_GetAll(ctx context.Context) ([]domain.Categoria, error)
	CategoryRepository_GetByID(ctx context.Context, id int) (*domain.Categoria, error)
	CategoryRepository_Update(ctx context.Context, c *domain.Categoria) error
	CategoryRepository_Delete(ctx context.Context, id int) error
}

type transactionRunner interface {
	WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

// CategoryService implementa la logica di business per la gestione delle categorie,
// validando i dati e orchestrando le operazioni sul repository all'interno di transazioni.
type CategoryService struct {
	Repo categoryRepository
	Tx   transactionRunner
}

// NewCategoryService crea un nuovo CategoryService con le dipendenze specificate.
func NewCategoryService(repo categoryRepository, tx transactionRunner) *CategoryService {
	return &CategoryService{Repo: repo, Tx: tx}
}

// CreateCategory crea una nuova categoria nel sistema
func (s *CategoryService) CategoryService_Create(ctx context.Context, c *domain.Categoria) error {
	// Trim dei spazi dal nome
	c.Nome = strings.TrimSpace(c.Nome)
	if c.Nome == "" {
		return errors.New("il nome della categoria non può essere vuoto")
	}

	// Esegue l'operazione in transazione per garantire coerenza
	return s.Tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		return s.Repo.CategoryRepository_Create(txCtx, c)
	})
}

// GetCategories recupera tutte le categorie disponibili nel sistema
func (s *CategoryService) CategoryService_GetAll(ctx context.Context) ([]domain.Categoria, error) {
	return s.Repo.CategoryRepository_GetAll(ctx)
}

// GetCategoryByID recupera una singola categoria dal sistema usando il suo ID
func (s *CategoryService) CategoryService_GetByID(ctx context.Context, id int) (*domain.Categoria, error) {
	return s.Repo.CategoryRepository_GetByID(ctx, id)
}

// UpdateCategory modifica una categoria esistente nel sistema
func (s *CategoryService) CategoryService_Update(ctx context.Context, c *domain.Categoria) error {
	// Trim dei spazi dal nome
	c.Nome = strings.TrimSpace(c.Nome)
	if c.Nome == "" {
		return errors.New("il nome della categoria non può essere vuoto")
	}

	// Esegue l'operazione in transazione per garantire coerenza
	return s.Tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		return s.Repo.CategoryRepository_Update(txCtx, c)
	})
}

// DeleteCategory rimuove una categoria dal sistema
func (s *CategoryService) CategoryService_Delete(ctx context.Context, id int) error {
	return s.Tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		return s.Repo.CategoryRepository_Delete(txCtx, id)
	})
}
