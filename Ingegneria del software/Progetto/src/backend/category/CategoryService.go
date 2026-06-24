package category

import (
	"context"
	"errors"
	"strings"
	"gialosai/backend/db"
	"gialosai/backend/domain"
)

// Service contiene la logica di business per la gestione delle categorie.
// Si interfaccia con il Repository per l'accesso ai dati e gestisce le regole di dominio.
type Service struct {
	Repo *Repository
	Tx   *db.TransactionHelper
}

// NewService crea una nuova istanza del service per le categorie.
func NewService(repo *Repository, tx *db.TransactionHelper) *Service {
	return &Service{Repo: repo, Tx: tx}
}

// CreateCategory crea una nuova categoria nel sistema.
// Valida che il nome non sia vuoto, trimma gli spazi e persiste nel database entro una transazione.
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

// GetCategories recupera tutte le categorie disponibili nel sistema.
// Restituisce una slice di categorie ordinate alfabeticamente per nome.
func (s *Service) GetCategories(ctx context.Context) ([]domain.Categoria, error) {
	return s.Repo.GetAll(ctx)
}

// GetCategoryByID recupera una singola categoria dal sistema usando il suo ID.
// Restituisce un errore se la categoria non esiste.
func (s *Service) GetCategoryByID(ctx context.Context, id int) (*domain.Categoria, error) {
	return s.Repo.GetByID(ctx, id)
}

// UpdateCategory modifica una categoria esistente nel sistema.
// Valida che il nome non sia vuoto e persiste le modifiche nel database entro una transazione.
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

// DeleteCategory rimuove una categoria dal sistema.
// Restituisce un errore se l'operazione fallisce nel database.
func (s *Service) DeleteCategory(ctx context.Context, id int) error {
	return s.Tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		return s.Repo.Delete(txCtx, id)
	})
}
