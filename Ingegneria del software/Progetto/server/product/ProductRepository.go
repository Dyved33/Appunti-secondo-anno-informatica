package product

import (
	"context"
	"database/sql"
	"errors"
	"github.com/tommasogarofani/ingsw_25-26_group-11/server/domain"
	"time"
)

type ProductRepository struct {
	DB *sql.DB
}


func NewProductRepository(database *sql.DB) *ProductRepository {
	return &ProductRepository{DB: database}
}

// GetAll recupera tutti i prodotti dal database con i dati delle categorie associate via LEFT JOIN.
func (r *ProductRepository) ProductRepository_GetAllProducts(ctx context.Context) ([]domain.Prodotto, error) {
	exec := domain.GetExecutor(ctx, r.DB)
	query := `SELECT p.id, p.codice, p.barcode, p.nome, p.descrizione, p.categoria_id, p.prezzo_acquisto, p.prezzo_vendita, p.soglia_minima,
	                 c.nome, c.descrizione
	          FROM Prodotti p
	          LEFT JOIN Categorie c 
			  ON p.categoria_id = c.id`

	rows, err := exec.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()// la presenza di defer indica che viene ritardata la chiusura delle righe in attesa della terminazione della funzione 

	// questo codice serve a leggere tutte le righe restituira dalla query e trasformarle in una slice di oggetti Prodotto
	var products []domain.Prodotto// vettore vuoto di prodotti
	for rows.Next() {// passa alla riga successiva e se questa esiste restituisce true, altrimenti false
		var p domain.Prodotto

		//tipi speciali quando una colonna SQL può contenere NULL
		var catID sql.NullInt64
		var catNome, catDesc sql.NullString

		err := rows.Scan(// scan copia i valori della riga corrente nelle variabili passate 
			&p.ID, &p.Codice, &p.Barcode, &p.Nome, &p.Descrizione, &catID, &p.PrezzoAcq, &p.PrezzoVend, &p.SogliaMin,
			&catNome, &catDesc,
		)
		if err != nil {
			return nil, err
		}

		// Popola i dati della categoria se presenti
		if catID.Valid {
			p.CategoriaID = int(catID.Int64)// int64 -> int

			p.Categoria = &domain.Categoria{// crea un nuovo oggetto Categoria e lo assegna al campo Categoria del prodotto
				ID:          p.CategoriaID,
				Nome:        catNome.String,
				Descrizione: catDesc.String,
			}
		}
		products = append(products, p)
	}
	if err := rows.Err(); err != nil {
 		return nil, err
 	}
	return products, nil
}

// GetByID recupera un singolo prodotto dal database usando il suo ID.
func (r *ProductRepository) ProductRepository_GetByID(ctx context.Context, id int) (*domain.Prodotto, error) {
	exec := domain.GetExecutor(ctx, r.DB)
	query := `SELECT p.id, p.codice, p.barcode, p.nome, p.descrizione, p.categoria_id, p.prezzo_acquisto, p.prezzo_vendita, p.soglia_minima,
	                 c.nome, c.descrizione
	          FROM Prodotti p
	          LEFT JOIN Categorie c 
			  ON p.categoria_id = c.id
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

// GetByBarcode recupera un prodotto dal database usando il suo barcode (EAN-13).
func (r *ProductRepository) ProductRepository_GetByBarcode(ctx context.Context, barcode string) (*domain.Prodotto, error) {
	exec := domain.GetExecutor(ctx, r.DB)
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

// Create inserisce un nuovo prodotto nel database e inizializza la giacenza a zero.
func (r *ProductRepository) ProductRepository_Create(ctx context.Context, p *domain.Prodotto) error {
	exec := domain.GetExecutor(ctx, r.DB)
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

	// Inizializza la giacenza a zero per il nuovo prodotto
	queryG := `INSERT INTO Giacenze (prodotto_id, quantita, updated_at) VALUES (?, 0, ?)`
	_, err = exec.ExecContext(ctx, queryG, p.ID, time.Now())
	return err
}

// Update modifica i dati di un prodotto esistente nel database.
func (r *ProductRepository) ProductRepository_Update(ctx context.Context, p *domain.Prodotto) error {
	exec := domain.GetExecutor(ctx, r.DB)
	query := `UPDATE Prodotti SET codice = ?, barcode = ?, nome = ?, descrizione = ?, categoria_id = ?, prezzo_acquisto = ?, prezzo_vendita = ?, soglia_minima = ?
	          WHERE id = ?`
	_, err := exec.ExecContext(ctx, query, p.Codice, p.Barcode, p.Nome, p.Descrizione, p.CategoriaID, p.PrezzoAcq, p.PrezzoVend, p.SogliaMin, p.ID)
	return err
}

// Delete rimuove un prodotto dal database e elimina anche la sua giacenza.
func (r *ProductRepository) ProductRepository_Delete(ctx context.Context, id int) error {
	exec := domain.GetExecutor(ctx, r.DB)
	// Elimina prima la giacenza associata
	if _, err := exec.ExecContext(ctx, "DELETE FROM Giacenze WHERE prodotto_id = ?", id); err != nil {
 		return err
 	}
	query := `DELETE FROM Prodotti WHERE id = ?`
	_, err := exec.ExecContext(ctx, query, id)
	return err
}

// GetGiacenza recupera la giacenza di un prodotto dal database. Se la giacenza non esiste, restituisce una giacenza vuota (quantità 0).
func (r *ProductRepository) ProductRepository_GetGiacenzaById(ctx context.Context, prodID int) (*domain.Giacenza, error) {
	exec := domain.GetExecutor(ctx, r.DB)
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

// UpdateGiacenza modifica la giacenza di un prodotto aggiungendo o sottraendo una quantità 
func (r *ProductRepository) ProductRepository_UpdateGiacenza(ctx context.Context, g *domain.Giacenza) error {
	exec := domain.GetExecutor(ctx, r.DB)
	query := `INSERT INTO Giacenze (prodotto_id, quantita, updated_at) VALUES (?, ?, ?)
	          ON DUPLICATE KEY UPDATE quantita = quantita + VALUES(quantita), updated_at = VALUES(updated_at)`
	_, err := exec.ExecContext(ctx, query, g.ProdottoID, g.Quantita, time.Now())
	return err
}

// GetAllGiacenze recupera tutte le giacenze dal database insieme al nome e alla soglia minima del prodotto.
func (r *ProductRepository) ProductRepository_GetAllGiacenze(ctx context.Context) ([]domain.Giacenza, error) {
	exec := domain.GetExecutor(ctx, r.DB)
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
	if err := rows.Err(); err != nil {
 		return nil, err
 	}
	return list, nil
}
