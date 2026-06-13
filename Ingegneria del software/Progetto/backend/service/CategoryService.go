package Service

import (
	"context"
	"errors"
	"strings"

	"gialosai/backend/domain"
	"gialosai/backend/repository"
)

type CategoryService struct {
	Repo *Repository.CategoryRepository
	Tx   *Repository.TransactionHelper
}

func NewCategoryService(repo *Repository.CategoryRepository, tx *Repository.TransactionHelper) *CategoryService {
	return &CategoryService{Repo: repo, Tx: tx}
}

func (s *CategoryService) CreateCategory(ctx context.Context, c *domain.Categoria) error {
	c.Nome = strings.TrimSpace(c.Nome)
	if c.Nome == "" {
		return errors.New("il nome della categoria non può essere vuoto")
	}

	return s.Tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		return s.Repo.Create(txCtx, c)
	})
}

func (s *CategoryService) GetCategories(ctx context.Context) ([]domain.Categoria, error) {
	return s.Repo.GetAll(ctx)
}

func (s *CategoryService) GetCategoryByID(ctx context.Context, id int) (*domain.Categoria, error) {
	return s.Repo.GetByID(ctx, id)
}

func (s *CategoryService) UpdateCategory(ctx context.Context, c *domain.Categoria) error {
	c.Nome = strings.TrimSpace(c.Nome)
	if c.Nome == "" {
		return errors.New("il nome della categoria non può essere vuoto")
	}

	return s.Tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		return s.Repo.Update(txCtx, c)
	})
}

func (s *CategoryService) DeleteCategory(ctx context.Context, id int) error {
	return s.Tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		return s.Repo.Delete(txCtx, id)
	})
}