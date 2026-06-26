package category

import (
	"context"
	"errors"
	"strings"
	"gialosai/backend/domain"
)

type Service struct {
	Repo *Repository
	Tx   *domain.TransactionHelper
}

func NewService(repo *Repository, tx *domain.TransactionHelper) *Service {
	return &Service{Repo: repo, Tx: tx}
}

// CreateCategory crea una nuova categoria nel sistema
func (s *Service) CreateCategory(ctx context.Context, c *domain.Categoria) error {
	// Trim dei spazi dal nome
	c.Nome = strings.TrimSpace(c.Nome)
	if c.Nome == "" {
		return errors.New("il nome della categoria non può essere vuoto")
	}

	// Esegue l'operazione in transazione per garantire coerenza
	return s.Tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		return s.Repo.Create(txCtx, c)
	})
}

// GetCategories recupera tutte le categorie disponibili nel sistema
func (s *Service) GetCategories(ctx context.Context) ([]domain.Categoria, error) {
	return s.Repo.GetAll(ctx)
}

// GetCategoryByID recupera una singola categoria dal sistema usando il suo ID
func (s *Service) GetCategoryByID(ctx context.Context, id int) (*domain.Categoria, error) {
	return s.Repo.GetByID(ctx, id)
}

// UpdateCategory modifica una categoria esistente nel sistema
func (s *Service) UpdateCategory(ctx context.Context, c *domain.Categoria) error {
	// Trim dei spazi dal nome
	c.Nome = strings.TrimSpace(c.Nome)
	if c.Nome == "" {
		return errors.New("il nome della categoria non può essere vuoto")
	}

	// Esegue l'operazione in transazione per garantire coerenza
	return s.Tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		return s.Repo.Update(txCtx, c)
	})
}

// DeleteCategory rimuove una categoria dal sistema
func (s *Service) DeleteCategory(ctx context.Context, id int) error {
	return s.Tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		return s.Repo.Delete(txCtx, id)
	})
}
