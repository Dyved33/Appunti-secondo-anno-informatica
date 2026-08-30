package product

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tommasogarofani/ingsw_25-26_group-11/server/domain"
)

func newMockProductRepo(t *testing.T) (*ProductRepository, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	return NewProductRepository(db), mock, func() { db.Close() }
}

func TestProductRepository_GetAll(t *testing.T) {
	repo, mock, closeDB := newMockProductRepo(t)
	defer closeDB()

	rows := sqlmock.NewRows([]string{
		"id", "codice", "barcode", "nome", "descrizione", "categoria_id",
		"prezzo_acquisto", "prezzo_vendita", "soglia_minima", "cat_nome", "cat_descrizione",
	}).
		AddRow(1, "P001", "8001234567890", "Vite M6", "Vite a brugola", 10, 0.10, 0.25, 100, "Ferramenta", "Minuteria metallica").
		AddRow(2, "P002", "8009876543210", "Bullone M8", "", nil, 0.20, 0.40, 50, nil, nil)

	mock.ExpectQuery(`SELECT p.id, p.codice, p.barcode, p.nome, p.descrizione, p.categoria_id`).
		WillReturnRows(rows)

	got, err := repo.ProductRepository_GetAllProducts(context.Background())
	require.NoError(t, err)
	require.Len(t, got, 2)

	assert.Equal(t, "Vite M6", got[0].Nome)
	require.NotNil(t, got[0].Categoria, "il prodotto con categoria_id valido deve avere Categoria popolata")
	assert.Equal(t, "Ferramenta", got[0].Categoria.Nome)

	assert.Nil(t, got[1].Categoria, "il prodotto senza categoria (NULL) non deve avere Categoria popolata")
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestProductRepository_GetByID(t *testing.T) {
	tests := []struct {
		name      string
		id        int
		mockSetup func(mock sqlmock.Sqlmock)
		wantErr   error
	}{
		{
			name: "prodotto trovato senza categoria",
			id:   5,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "codice", "barcode", "nome", "descrizione", "categoria_id",
					"prezzo_acquisto", "prezzo_vendita", "soglia_minima", "cat_nome", "cat_descrizione",
				}).AddRow(5, "P005", "8005555555555", "Martello", "", nil, 5.0, 9.99, 10, nil, nil)
				mock.ExpectQuery(`SELECT p.id, p.codice, p.barcode, p.nome, p.descrizione, p.categoria_id`).
					WithArgs(5).
					WillReturnRows(rows)
			},
		},
		{
			name: "prodotto non trovato",
			id:   999,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT p.id, p.codice, p.barcode, p.nome, p.descrizione, p.categoria_id`).
					WithArgs(999).
					WillReturnError(sql.ErrNoRows)
			},
			wantErr: domain.ErrProdottoNonTrovato,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, mock, closeDB := newMockProductRepo(t)
			defer closeDB()
			tt.mockSetup(mock)

			got, err := repo.ProductRepository_GetByID(context.Background(), tt.id)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, got)
			} else {
				require.NoError(t, err)
				require.NotNil(t, got)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestProductRepository_GetByBarcode(t *testing.T) {
	repo, mock, closeDB := newMockProductRepo(t)
	defer closeDB()

	rows := sqlmock.NewRows([]string{
		"id", "codice", "barcode", "nome", "descrizione", "categoria_id",
		"prezzo_acquisto", "prezzo_vendita", "soglia_minima",
	}).AddRow(1, "P001", "8001234567890", "Vite M6", "", 10, 0.10, 0.25, 100)

	mock.ExpectQuery(`SELECT p.id, p.codice, p.barcode, p.nome, p.descrizione, p.categoria_id, p.prezzo_acquisto, p.prezzo_vendita, p.soglia_minima\s+FROM Prodotti p WHERE p.barcode = \?`).
		WithArgs("8001234567890").
		WillReturnRows(rows)

	got, err := repo.ProductRepository_GetByBarcode(context.Background(), "8001234567890")
	require.NoError(t, err)
	assert.Equal(t, "Vite M6", got.Nome)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestProductRepository_Create(t *testing.T) {
	repo, mock, closeDB := newMockProductRepo(t)
	defer closeDB()

	mock.ExpectExec(`INSERT INTO Prodotti`).
		WithArgs("P010", "8001112223334", "Chiave inglese", "", 0, 3.5, 7.0, 5).
		WillReturnResult(sqlmock.NewResult(10, 1))
	mock.ExpectExec(`INSERT INTO Giacenze \(prodotto_id, quantita, updated_at\) VALUES \(\?, 0, \?\)`).
		WithArgs(10, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	p := &domain.Prodotto{
		Codice:     "P010",
		Barcode:    "8001112223334",
		Nome:       "Chiave inglese",
		PrezzoAcq:  3.5,
		PrezzoVend: 7.0,
		SogliaMin:  5,
	}
	err := repo.ProductRepository_Create(context.Background(), p)
	require.NoError(t, err)
	assert.Equal(t, 10, p.ID)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestProductRepository_Update(t *testing.T) {
	repo, mock, closeDB := newMockProductRepo(t)
	defer closeDB()

	mock.ExpectExec(`UPDATE Prodotti SET codice = \?, barcode = \?, nome = \?, descrizione = \?, categoria_id = \?, prezzo_acquisto = \?, prezzo_vendita = \?, soglia_minima = \? WHERE id = \?`).
		WithArgs("P010", "8001112223334", "Chiave inglese", "", 0, 3.5, 7.0, 5, 10).
		WillReturnResult(sqlmock.NewResult(0, 1))

	p := &domain.Prodotto{
		ID:         10,
		Codice:     "P010",
		Barcode:    "8001112223334",
		Nome:       "Chiave inglese",
		PrezzoAcq:  3.5,
		PrezzoVend: 7.0,
		SogliaMin:  5,
	}
	err := repo.ProductRepository_Update(context.Background(), p)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestProductRepository_Delete(t *testing.T) {
	repo, mock, closeDB := newMockProductRepo(t)
	defer closeDB()

	mock.ExpectExec(`DELETE FROM Giacenze WHERE prodotto_id = \?`).
		WithArgs(3).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(`DELETE FROM Prodotti WHERE id = \?`).
		WithArgs(3).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.ProductRepository_Delete(context.Background(), 3)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

//===============================

func TestProductRepository_GetGiacenzaById(t *testing.T) {
	t.Run("giacenza esistente", func(t *testing.T) {
		repo, mock, closeDB := newMockProductRepo(t)
		defer closeDB()

		rows := sqlmock.NewRows([]string{"id", "prodotto_id", "quantita", "updated_at"}).
			AddRow(1, 7, 42, time.Now())
		mock.ExpectQuery(`SELECT id, prodotto_id, quantita, updated_at FROM Giacenze WHERE prodotto_id = \?`).
			WithArgs(7).
			WillReturnRows(rows)

		g, err := repo.ProductRepository_GetGiacenzaById(context.Background(), 7)
		require.NoError(t, err)
		assert.Equal(t, 42, g.Quantita)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("giacenza assente restituisce quantita zero senza errore", func(t *testing.T) {
		repo, mock, closeDB := newMockProductRepo(t)
		defer closeDB()

		mock.ExpectQuery(`SELECT id, prodotto_id, quantita, updated_at FROM Giacenze WHERE prodotto_id = \?`).
			WithArgs(99).
			WillReturnError(sql.ErrNoRows)

		g, err := repo.ProductRepository_GetGiacenzaById(context.Background(), 99)
		require.NoError(t, err)
		require.NotNil(t, g)
		assert.Equal(t, 0, g.Quantita)
		assert.Equal(t, 99, g.ProdottoID)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("errore generico del database viene propagato", func(t *testing.T) {
		repo, mock, closeDB := newMockProductRepo(t)
		defer closeDB()

		mock.ExpectQuery(`SELECT id, prodotto_id, quantita, updated_at FROM Giacenze WHERE prodotto_id = \?`).
			WithArgs(7).
			WillReturnError(errors.New("db irraggiungibile"))

		_, err := repo.ProductRepository_GetGiacenzaById(context.Background(), 7)
		assert.Error(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestProductRepository_GetAllGiacenze(t *testing.T) {
	repo, mock, closeDB := newMockProductRepo(t)
	defer closeDB()

	rows := sqlmock.NewRows([]string{"id", "prodotto_id", "quantita", "updated_at", "nome", "soglia_minima"}).
		AddRow(1, 1, 5, time.Now(), "Vite M6", 100).
		AddRow(2, 2, 200, time.Now(), "Bullone M8", 50)

	mock.ExpectQuery(`SELECT g.id, g.prodotto_id, g.quantita, g.updated_at`).
		WillReturnRows(rows)

	got, err := repo.ProductRepository_GetAllGiacenze(context.Background())
	require.NoError(t, err)
	require.Len(t, got, 2)
	assert.Equal(t, "Vite M6", got[0].Prodotto.Nome)
	assert.Equal(t, 100, got[0].Prodotto.SogliaMin)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestProductRepository_UpdateGiacenza(t *testing.T) {
	repo, mock, closeDB := newMockProductRepo(t)
	defer closeDB()

	mock.ExpectExec(`INSERT INTO Giacenze \(prodotto_id, quantita, updated_at\) VALUES \(\?, \?, \?\) ON DUPLICATE KEY UPDATE quantita = quantita \+ VALUES\(quantita\), updated_at = VALUES\(updated_at\)`).
		WithArgs(1, 150, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(0, 1))

	g := &domain.Giacenza{
		ProdottoID: 1,
		Quantita:   150,
	}
	err := repo.ProductRepository_UpdateGiacenza(context.Background(), g)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}
