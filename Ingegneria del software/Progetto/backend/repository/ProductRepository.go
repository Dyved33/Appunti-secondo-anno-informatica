package Repository

import (
	"context"
	"database/sql"
	"errors"
	"gialosai/backend/domain"
	"time"
)

type ProductRepository struct {
	DB *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{DB: db}
}

func (r *ProductRepository) GetAll(ctx context.Context) ([]domain.Prodotto, error) {
	exec := GetExecutor(ctx, r.DB)
	query := `SELECT p.id, p.codice, p.barcode, p.nome, p.descrizione, p.categoria_id, p.prezzo_acquisto, p.prezzo_vendita, p.soglia_minima,
	                 c.nome, c.descrizione
	          FROM Prodotti p
	          LEFT JOIN Categorie c ON p.categoria_id = c.id`

	rows, err := exec.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []domain.Prodotto
	for rows.Next() {
		var p domain.Prodotto
		var catID sql.NullInt64
		var catNome, catDesc sql.NullString
		err := rows.Scan(
			&p.ID, &p.Codice, &p.Barcode, &p.Nome, &p.Descrizione, &catID, &p.PrezzoAcq, &p.PrezzoVend, &p.SogliaMin,
			&catNome, &catDesc,
		)
		if err != nil {
			return nil, err
		}
		if catID.Valid {
			p.CategoriaID = int(catID.Int64)
			p.Categoria = &domain.Categoria{
				ID:          p.CategoriaID,
				Nome:        catNome.String,
				Descrizione: catDesc.String,
			}
		}
		products = append(products, p)
	}
	return products, nil
}

func (r *ProductRepository) GetByID(ctx context.Context, id int) (*domain.Prodotto, error) {
	exec := GetExecutor(ctx, r.DB)
	query := `SELECT p.id, p.codice, p.barcode, p.nome, p.descrizione, p.categoria_id, p.prezzo_acquisto, p.prezzo_vendita, p.soglia_minima,
	                 c.nome, c.descrizione
	          FROM Prodotti p
	          LEFT JOIN Categorie c ON p.categoria_id = c.id
	          WHERE p.id = ?`

	var p domain.Prodotto
	var catID sql.NullInt64
	var catNome, catDesc sql.NullString
	err := exec.QueryRowContext(ctx, query, id).Scan(
		&p.ID, &p.Codice, &p.Barcode, &p.Nome, &p.Descrizione, &catID, &p.PrezzoAcq, &p.PrezzoVend, &p.SogliaMin,
		&catNome, &catDesc,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrProdottoNonTrovato
		}
		return nil, err
	}
	if catID.Valid {
		p.CategoriaID = int(catID.Int64)
		p.Categoria = &domain.Categoria{
			ID:          p.CategoriaID,
			Nome:        catNome.String,
			Descrizione: catDesc.String,
		}
	}
	return &p, nil
}

func (r *ProductRepository) GetByBarcode(ctx context.Context, barcode string) (*domain.Prodotto, error) {
	exec := GetExecutor(ctx, r.DB)
	query := `SELECT p.id, p.codice, p.barcode, p.nome, p.descrizione, p.categoria_id, p.prezzo_acquisto, p.prezzo_vendita, p.soglia_minima
	          FROM Prodotti p WHERE p.barcode = ?`

	var p domain.Prodotto
	err := exec.QueryRowContext(ctx, query, barcode).Scan(
		&p.ID, &p.Codice, &p.Barcode, &p.Nome, &p.Descrizione, &p.CategoriaID, &p.PrezzoAcq, &p.PrezzoVend, &p.SogliaMin,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrProdottoNonTrovato
		}
		return nil, err
	}
	return &p, nil
}

func (r *ProductRepository) Create(ctx context.Context, p *domain.Prodotto) error {
	exec := GetExecutor(ctx, r.DB)
	query := `INSERT INTO Prodotti (codice, barcode, nome, descrizione, categoria_id, prezzo_acquisto, prezzo_vendita, soglia_minima)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	res, err := exec.ExecContext(ctx, query, p.Codice, p.Barcode, p.Nome, p.Descrizione, p.CategoriaID, p.PrezzoAcq, p.PrezzoVend, p.SogliaMin)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err == nil {
		p.ID = int(id)
	}

	// Initialize empty stock
	queryG := `INSERT INTO Giacenze (prodotto_id, quantita, updated_at) VALUES (?, 0, ?)`
	_, err = exec.ExecContext(ctx, queryG, p.ID, time.Now())
	return err
}

func (r *ProductRepository) Update(ctx context.Context, p *domain.Prodotto) error {
	exec := GetExecutor(ctx, r.DB)
	query := `UPDATE Prodotti SET codice = ?, barcode = ?, nome = ?, descrizione = ?, categoria_id = ?, prezzo_acquisto = ?, prezzo_vendita = ?, soglia_minima = ?
	          WHERE id = ?`
	_, err := exec.ExecContext(ctx, query, p.Codice, p.Barcode, p.Nome, p.Descrizione, p.CategoriaID, p.PrezzoAcq, p.PrezzoVend, p.SogliaMin, p.ID)
	return err
}

func (r *ProductRepository) Delete(ctx context.Context, id int) error {
	exec := GetExecutor(ctx, r.DB)
	// First delete stock
	_, _ = exec.ExecContext(ctx, "DELETE FROM Giacenze WHERE prodotto_id = ?", id)
	query := `DELETE FROM Prodotti WHERE id = ?`
	_, err := exec.ExecContext(ctx, query, id)
	return err
}

// Giacenza Management
func (r *ProductRepository) GetGiacenza(ctx context.Context, prodID int) (*domain.Giacenza, error) {
	exec := GetExecutor(ctx, r.DB)
	query := `SELECT id, prodotto_id, quantita, updated_at FROM Giacenze WHERE prodotto_id = ?`
	var g domain.Giacenza
	err := exec.QueryRowContext(ctx, query, prodID).Scan(&g.ID, &g.ProdottoID, &g.Quantita, &g.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &domain.Giacenza{ProdottoID: prodID, Quantita: 0, UpdatedAt: time.Now()}, nil
		}
		return nil, err
	}
	return &g, nil
}

func (r *ProductRepository) UpdateGiacenza(ctx context.Context, prodID int, delta int) error {
	exec := GetExecutor(ctx, r.DB)
	query := `INSERT INTO Giacenze (prodotto_id, quantita, updated_at) VALUES (?, ?, ?)
	          ON CONFLICT(prodotto_id) DO UPDATE SET quantita = quantita + excluded.quantita, updated_at = excluded.updated_at`
	_, err := exec.ExecContext(ctx, query, prodID, delta, time.Now())
	return err
}

func (r *ProductRepository) GetAllGiacenze(ctx context.Context) ([]domain.Giacenza, error) {
	exec := GetExecutor(ctx, r.DB)
	query := `SELECT g.id, g.prodotto_id, g.quantita, g.updated_at,
	                 p.nome, p.soglia_minima
	          FROM Giacenze g
	          JOIN Prodotti p ON g.prodotto_id = p.id`

	rows, err := exec.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.Giacenza
	for rows.Next() {
		var g domain.Giacenza
		var p domain.Prodotto
		err := rows.Scan(&g.ID, &g.ProdottoID, &g.Quantita, &g.UpdatedAt, &p.Nome, &p.SogliaMin)
		if err != nil {
			return nil, err
		}
		g.Prodotto = &p
		list = append(list, g)
	}
	return list, nil
}
