package category

import (
	"context"
	"database/sql"
	"errors"
	"github.com/tommasogarofani/ingsw_25-26_group-11/server/domain"
)

// CategoryRepository fornisce l'accesso persistente (SQL) alle categorie di prodotti.
type CategoryRepository struct {
	DB *sql.DB
}

// NewCategoryRepository crea un nuovo CategoryRepository legato al database specificato.
func NewCategoryRepository(database *sql.DB) *CategoryRepository {
	return &CategoryRepository{DB: database}
}

// CategoryRepository_Create inserisce una nuova categoria nel database.
func (r *CategoryRepository) CategoryRepository_Create(ctx context.Context, c *domain.Categoria) error {
	exec := domain.GetExecutor(ctx, r.DB)
	query := `INSERT INTO Categorie (nome, descrizione) VALUES (?, ?)`
	res, err := exec.ExecContext(ctx, query, c.Nome, c.Descrizione)
	if err != nil {
		return err
	}
	// Recupera l'ID auto-generato dal database e lo assegna alla categoria
	id, err := res.LastInsertId()
	if err == nil {
		c.ID = int(id)
	}
	return nil
}

// CategoryRepository_GetAll recupera tutte le categorie dal database, ordinate per nome.
func (r *CategoryRepository) CategoryRepository_GetAll(ctx context.Context) ([]domain.Categoria, error) {
	exec := domain.GetExecutor(ctx, r.DB)
	rows, err := exec.QueryContext(ctx, `SELECT id, nome, descrizione FROM Categorie ORDER BY nome`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []domain.Categoria
	for rows.Next() {
		var c domain.Categoria
		if err := rows.Scan(&c.ID, &c.Nome, &c.Descrizione); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}

// CategoryRepository_GetByID recupera una singola categoria dal database usando il suo ID.
func (r *CategoryRepository) CategoryRepository_GetByID(ctx context.Context, id int) (*domain.Categoria, error) {
	exec := domain.GetExecutor(ctx, r.DB)
	var c domain.Categoria
	err := exec.QueryRowContext(ctx, `SELECT id, nome, descrizione FROM Categorie WHERE id = ?`, id).Scan(&c.ID, &c.Nome, &c.Descrizione)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrCategoriaNonTrovata
		}
		return nil, err
	}
	return &c, nil
}

// CategoryRepository_Update modifica i dati di una categoria esistente nel database.
func (r *CategoryRepository) CategoryRepository_Update(ctx context.Context, c *domain.Categoria) error {
	exec := domain.GetExecutor(ctx, r.DB)
	_, err := exec.ExecContext(ctx, `UPDATE Categorie SET nome = ?, descrizione = ? WHERE id = ?`, c.Nome, c.Descrizione, c.ID)
	return err
}

// CategoryRepository_Delete rimuove una categoria dal database usando il suo ID.
func (r *CategoryRepository) CategoryRepository_Delete(ctx context.Context, id int) error {
	exec := domain.GetExecutor(ctx, r.DB)
	_, err := exec.ExecContext(ctx, `DELETE FROM Categorie WHERE id = ?`, id)
	return err
}
