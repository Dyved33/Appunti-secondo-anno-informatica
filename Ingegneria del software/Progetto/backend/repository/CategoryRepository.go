package Repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"gialosai/backend/domain"
)

type CategoryRepository struct {
	DB *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{DB: db}
}

func (r *CategoryRepository) Create(ctx context.Context, c *domain.Categoria) error {
	exec := GetExecutor(ctx, r.DB)
	query := `INSERT INTO Categorie (nome, descrizione) VALUES (?, ?)`
	res, err := exec.ExecContext(ctx, query, strings.TrimSpace(c.Nome), c.Descrizione)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err == nil {
		c.ID = int(id)
	}
	return nil
}

func (r *CategoryRepository) GetAll(ctx context.Context) ([]domain.Categoria, error) {
	exec := GetExecutor(ctx, r.DB)
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

func (r *CategoryRepository) GetByID(ctx context.Context, id int) (*domain.Categoria, error) {
	exec := GetExecutor(ctx, r.DB)
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

func (r *CategoryRepository) Update(ctx context.Context, c *domain.Categoria) error {
	exec := GetExecutor(ctx, r.DB)
	_, err := exec.ExecContext(ctx, `UPDATE Categorie SET nome = ?, descrizione = ? WHERE id = ?`, strings.TrimSpace(c.Nome), c.Descrizione, c.ID)
	return err
}

func (r *CategoryRepository) Delete(ctx context.Context, id int) error {
	exec := GetExecutor(ctx, r.DB)
	_, err := exec.ExecContext(ctx, `DELETE FROM Categorie WHERE id = ?`, id)
	return err
}