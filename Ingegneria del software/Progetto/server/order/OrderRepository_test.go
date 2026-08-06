package order

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tommasogarofani/ingsw_25-26_group-11/server/domain"
)

func newMockOrderRepo(t *testing.T) (*OrderRepository, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	return NewOrderRepository(db), mock, func() { db.Close() }
}

func TestOrderRepository_CreateOrdineEntrata(t *testing.T) {
	repo, mock, closeDB := newMockOrderRepo(t)
	defer closeDB()

	now := time.Now()
	mock.ExpectExec(`INSERT INTO OrdiniEntrata`).
		WithArgs("OE-001", "ACME Srl", now, "confermato", "", 1).
		WillReturnResult(sqlmock.NewResult(10, 1))
	mock.ExpectExec(`INSERT INTO RigheOrdineEntrata`).
		WithArgs(10, 5, 100, 0.5).
		WillReturnResult(sqlmock.NewResult(1, 1))

	o := &domain.OrdineEntrata{
		Numero:      "OE-001",
		Fornitore:   "ACME Srl",
		Data:        now,
		Stato:       domain.StatoOrdine("confermato"),
		CreatedByID: 1,
		Righe: []domain.RigaOrdineEntrata{
			{ProdottoID: 5, Quantita: 100, PrezzoUnitario: 0.5},
		},
	}

	err := repo.OrderRepository_CreateOrdineEntrata(context.Background(), o)
	require.NoError(t, err)
	assert.Equal(t, 10, o.ID)
	assert.Equal(t, 1, o.Righe[0].ID)
	assert.Equal(t, 10, o.Righe[0].OrdineEntrataID)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestOrderRepository_GetOrdineEntrataByID(t *testing.T) {
	t.Run("ordine non trovato", func(t *testing.T) {
		repo, mock, closeDB := newMockOrderRepo(t)
		defer closeDB()

		mock.ExpectQuery(`SELECT id, numero, fornitore, data, stato, note, created_by FROM OrdiniEntrata WHERE id = \?`).
			WithArgs(99).
			WillReturnError(sql.ErrNoRows)

		_, err := repo.OrderRepository_GetOrdineEntrataByID(context.Background(), 99)
		assert.Error(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ordine trovato con righe", func(t *testing.T) {
		repo, mock, closeDB := newMockOrderRepo(t)
		defer closeDB()

		headerRows := sqlmock.NewRows([]string{"id", "numero", "fornitore", "data", "stato", "note", "created_by"}).
			AddRow(10, "OE-001", "ACME Srl", time.Now(), "confermato", "", 1)
		mock.ExpectQuery(`SELECT id, numero, fornitore, data, stato, note, created_by FROM OrdiniEntrata WHERE id = \?`).
			WithArgs(10).
			WillReturnRows(headerRows)

		lineRows := sqlmock.NewRows([]string{"id", "ordine_entrata_id", "prodotto_id", "quantita", "prezzo_unitario"}).
			AddRow(1, 10, 5, 100, 0.5)
		mock.ExpectQuery(`SELECT id, ordine_entrata_id, prodotto_id, quantita, prezzo_unitario FROM RigheOrdineEntrata WHERE ordine_entrata_id = \?`).
			WithArgs(10).
			WillReturnRows(lineRows)

		o, err := repo.OrderRepository_GetOrdineEntrataByID(context.Background(), 10)
		require.NoError(t, err)
		assert.Equal(t, domain.StatoOrdine("confermato"), o.Stato)
		require.Len(t, o.Righe, 1)
		assert.Equal(t, 100, o.Righe[0].Quantita)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestOrderRepository_CreateOrdineUscita(t *testing.T) {
	repo, mock, closeDB := newMockOrderRepo(t)
	defer closeDB()

	now := time.Now()
	mock.ExpectExec(`INSERT INTO OrdiniUscita`).
		WithArgs("OU-001", "Cliente Spa", now, "confermato", "", 1, nil).
		WillReturnResult(sqlmock.NewResult(20, 1))
	mock.ExpectExec(`INSERT INTO RigheOrdineUscita`).
		WithArgs(20, 5, 10, 1.0).
		WillReturnResult(sqlmock.NewResult(2, 1))

	o := &domain.OrdineUscita{
		Numero:      "OU-001",
		Cliente:     "Cliente Spa",
		Data:        now,
		Stato:       domain.StatoOrdine("confermato"),
		CreatedByID: 1,
		Righe: []domain.RigaOrdineUscita{
			{ProdottoID: 5, Quantita: 10, PrezzoUnitario: 1.0},
		},
	}

	err := repo.OrderRepository_CreateOrdineUscita(context.Background(), o)
	require.NoError(t, err)
	assert.Equal(t, 20, o.ID)
	require.NoError(t, mock.ExpectationsWereMet())
}
 
func TestOrderRepository_GetOrdineUscitaByID(t *testing.T) {
	repo, mock, closeDB := newMockOrderRepo(t)
	defer closeDB()

	ddtID := 7
	headerRows := sqlmock.NewRows([]string{"id", "numero", "cliente", "data", "stato", "note", "created_by", "ddt_id"}).
		AddRow(20, "OU-001", "Cliente Spa", time.Now(), "evaso", "", 1, ddtID)
	mock.ExpectQuery(`SELECT id, numero, cliente, data, stato, note, created_by, ddt_id FROM OrdiniUscita WHERE id = \?`).
		WithArgs(20).
		WillReturnRows(headerRows)

	lineRows := sqlmock.NewRows([]string{"id", "ordine_uscita_id", "prodotto_id", "quantita", "prezzo_unitario"}).
		AddRow(2, 20, 5, 10, 1.0)
	mock.ExpectQuery(`SELECT id, ordine_uscita_id, prodotto_id, quantita, prezzo_unitario FROM RigheOrdineUscita WHERE ordine_uscita_id = \?`).
		WithArgs(20).
		WillReturnRows(lineRows)

	o, err := repo.OrderRepository_GetOrdineUscitaByID(context.Background(), 20)
	require.NoError(t, err)
	require.NotNil(t, o.DDTID)
	assert.Equal(t, 7, *o.DDTID)
	require.Len(t, o.Righe, 1)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestOrderRepository_GetOrdiniUscitaByStato(t *testing.T) {
	repo, mock, closeDB := newMockOrderRepo(t)
	defer closeDB()

	rows := sqlmock.NewRows([]string{"id", "numero", "cliente", "data", "stato", "note", "created_by", "ddt_id"}).
		AddRow(1, "OU-001", "Cliente A", time.Now(), "confermato", "", 1, nil).
		AddRow(2, "OU-002", "Cliente B", time.Now(), "confermato", "", 1, nil)

	mock.ExpectQuery(`SELECT id, numero, cliente, data, stato, note, created_by, ddt_id FROM OrdiniUscita WHERE stato = \?`).
		WithArgs("confermato").
		WillReturnRows(rows)

	list, err := repo.OrderRepository_GetOrdiniUscitaByStato(context.Background(), domain.StatoOrdine("confermato"))
	require.NoError(t, err)
	require.Len(t, list, 2)
	assert.Nil(t, list[0].DDTID)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestOrderRepository_UpdateOrdineUscita(t *testing.T) {
	repo, mock, closeDB := newMockOrderRepo(t)
	defer closeDB()

	mock.ExpectExec(`UPDATE OrdiniUscita SET stato = \?, ddt_id = \? WHERE id = \?`).
		WithArgs("evaso", 7, 20).
		WillReturnResult(sqlmock.NewResult(0, 1))

	ddtID := 7
	err := repo.OrderRepository_UpdateOrdineUscita(context.Background(), &domain.OrdineUscita{ID: 20, Stato: domain.StatoOrdine("evaso"), DDTID: &ddtID})
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestOrderRepository_CreateDDT(t *testing.T) {
	repo, mock, closeDB := newMockOrderRepo(t)
	defer closeDB()

	now := time.Now()
	mock.ExpectExec(`INSERT INTO DDT`).
		WithArgs("DDT-2026-0001", 20, now, "Magazzino Centrale", "Cliente Spa", "./output/pdf/ddt_20.pdf").
		WillReturnResult(sqlmock.NewResult(5, 1))

	d := &domain.DDT{
		Numero:         "DDT-2026-0001",
		OrdineUscitaID: 20,
		Data:           now,
		Mittente:       "Magazzino Centrale",
		Destinatario:   "Cliente Spa",
		PathPDF:        "./output/pdf/ddt_20.pdf",
	}

	err := repo.OrderRepository_CreateDDT(context.Background(), d)
	require.NoError(t, err)
	assert.Equal(t, 5, d.ID)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestOrderRepository_CountDDTByYear(t *testing.T) {
	repo, mock, closeDB := newMockOrderRepo(t)
	defer closeDB()

	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM DDT WHERE numero LIKE \?`).
		WithArgs("DDT-2026-%").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(4))

	count, err := repo.OrderRepository_CountDDTByYear(context.Background(), 2026)
	require.NoError(t, err)
	assert.Equal(t, 4, count)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestOrderRepository_GetDDTByOrdineUscita(t *testing.T) {
	t.Run("nessun DDT associato restituisce nil senza errore", func(t *testing.T) {
		repo, mock, closeDB := newMockOrderRepo(t)
		defer closeDB()

		mock.ExpectQuery(`SELECT id, numero, ordine_uscita_id, data, mittente, destinatario, path_pdf FROM DDT WHERE ordine_uscita_id = \?`).
			WithArgs(20).
			WillReturnError(sql.ErrNoRows)

		d, err := repo.OrderRepository_GetDDTByOrdineUscita(context.Background(), 20)
		require.NoError(t, err)
		assert.Nil(t, d)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DDT trovato", func(t *testing.T) {
		repo, mock, closeDB := newMockOrderRepo(t)
		defer closeDB()

		now := time.Now()
		rows := sqlmock.NewRows([]string{"id", "numero", "ordine_uscita_id", "data", "mittente", "destinatario", "path_pdf"}).
			AddRow(5, "DDT-2026-0001", 20, now, "Magazzino Centrale", "Cliente Spa", "./output/pdf/ddt_20.pdf")
		mock.ExpectQuery(`SELECT id, numero, ordine_uscita_id, data, mittente, destinatario, path_pdf FROM DDT WHERE ordine_uscita_id = \?`).
			WithArgs(20).
			WillReturnRows(rows)

		d, err := repo.OrderRepository_GetDDTByOrdineUscita(context.Background(), 20)
		require.NoError(t, err)
		require.NotNil(t, d)
		assert.Equal(t, "DDT-2026-0001", d.Numero)
		assert.Equal(t, 20, d.OrdineUscitaID)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestOrderRepository_GetOrdiniUscita(t *testing.T) {
	repo, mock, closeDB := newMockOrderRepo(t)
	defer closeDB()

	ddtID := 7
	rows := sqlmock.NewRows([]string{"id", "numero", "cliente", "data", "stato", "note", "created_by", "ddt_id"}).
		AddRow(1, "OU-001", "Cliente A", time.Now(), "confermato", "", 1, nil).
		AddRow(2, "OU-002", "Cliente B", time.Now(), "evaso", "", 1, ddtID)

	mock.ExpectQuery(`SELECT id, numero, cliente, data, stato, note, created_by, ddt_id FROM OrdiniUscita`).
		WillReturnRows(rows)

	list, err := repo.OrderRepository_GetOrdiniUscita(context.Background())
	require.NoError(t, err)
	require.Len(t, list, 2)
	assert.Nil(t, list[0].DDTID)
	require.NotNil(t, list[1].DDTID)
	assert.Equal(t, 7, *list[1].DDTID)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestOrderRepository_CreateFattura(t *testing.T) {
	repo, mock, closeDB := newMockOrderRepo(t)
	defer closeDB()

	now := time.Now()
	mock.ExpectExec(`INSERT INTO Fatture`).
		WithArgs("FT-2026-0001", 2026, "vendita", 20, 100.0, 22.0, 122.0, now, "./output/pdf/fattura_20.pdf").
		WillReturnResult(sqlmock.NewResult(8, 1))

	f := &domain.Fattura{
		Numero:         "FT-2026-0001",
		Anno:           2026,
		Tipo:           domain.TipoFattura("vendita"),
		OrdineUscitaID: 20,
		Imponibile:     100.0,
		IVA:            22.0,
		Totale:         122.0,
		Data:           now,
		PathPDF:        "./output/pdf/fattura_20.pdf",
	}

	err := repo.OrderRepository_CreateFattura(context.Background(), f)
	require.NoError(t, err)
	assert.Equal(t, 8, f.ID)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestOrderRepository_CreateReso(t *testing.T) {
	repo, mock, closeDB := newMockOrderRepo(t)
	defer closeDB()

	now := time.Now()
	mock.ExpectExec(`INSERT INTO Resi`).
		WithArgs(20, 5, 2, "Prodotto difettoso", "da_verificare", now, 1).
		WillReturnResult(sqlmock.NewResult(3, 1))

	re := &domain.Reso{
		OrdineUscitaID: 20,
		ProdottoID:     5,
		Quantita:       2,
		Motivo:         "Prodotto difettoso",
		Data:           now,
		CreatedByID:    1,
	}

	err := repo.OrderRepository_CreateReso(context.Background(), re, "da_verificare")
	require.NoError(t, err)
	assert.Equal(t, 3, re.ID)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestOrderRepository_GetResoByID(t *testing.T) {
	repo, mock, closeDB := newMockOrderRepo(t)
	defer closeDB()

	rows := sqlmock.NewRows([]string{"id", "ordine_uscita_id", "prodotto_id", "quantita", "motivo", "stato", "data", "created_by"}).
		AddRow(3, 20, 5, 2, "Prodotto difettoso", "da_verificare", time.Now(), 1)
	mock.ExpectQuery(`SELECT id, ordine_uscita_id, prodotto_id, quantita, motivo, stato, data, created_by FROM Resi WHERE id = \?`).
		WithArgs(3).
		WillReturnRows(rows)

	re, stato, err := repo.OrderRepository_GetResoByID(context.Background(), 3)
	require.NoError(t, err)
	assert.Equal(t, "da_verificare", stato)
	assert.Equal(t, 2, re.Quantita)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestOrderRepository_UpdateResoStato(t *testing.T) {
	repo, mock, closeDB := newMockOrderRepo(t)
	defer closeDB()

	mock.ExpectExec(`UPDATE Resi SET stato = \? WHERE id = \?`).
		WithArgs("approvato", 3).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.OrderRepository_UpdateResoStato(context.Background(), 3, "approvato")
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestOrderRepository_CountFattureByYearAndType(t *testing.T) {
	repo, mock, closeDB := newMockOrderRepo(t)
	defer closeDB()

	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM Fatture WHERE anno = \? AND tipo = \?`).
		WithArgs(2026, "vendita").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(3))

	count, err := repo.OrderRepository_CountFattureByYearAndType(context.Background(), 2026, domain.TipoFattura("vendita"))
	require.NoError(t, err)
	assert.Equal(t, 3, count)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestOrderRepository_GetAllFatture(t *testing.T) {
	repo, mock, closeDB := newMockOrderRepo(t)
	defer closeDB()

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "numero", "anno", "tipo", "ordine_uscita_id", "imponibile", "iva", "totale", "data", "path_pdf"}).
		AddRow(8, "FT-2026-0001", 2026, "vendita", 20, 100.0, 22.0, 122.0, now, "./output/pdf/fattura_20.pdf")

	mock.ExpectQuery(`SELECT id, numero, anno, tipo, ordine_uscita_id, imponibile, iva, totale, data, path_pdf FROM Fatture`).
		WillReturnRows(rows)

	list, err := repo.OrderRepository_GetAllFatture(context.Background())
	require.NoError(t, err)
	require.Len(t, list, 1)
	assert.Equal(t, domain.TipoFattura("vendita"), list[0].Tipo)
	assert.Equal(t, 122.0, list[0].Totale)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestOrderRepository_DeleteFattura(t *testing.T) {
	repo, mock, closeDB := newMockOrderRepo(t)
	defer closeDB()

	mock.ExpectExec(`DELETE FROM Fatture WHERE id = \?`).
		WithArgs(8).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.OrderRepository_DeleteFattura(context.Background(), 8)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestOrderRepository_GetAllResi(t *testing.T) {
	repo, mock, closeDB := newMockOrderRepo(t)
	defer closeDB()

	rows := sqlmock.NewRows([]string{"id", "ordine_uscita_id", "prodotto_id", "quantita", "motivo", "stato", "data", "created_by"}).
		AddRow(3, 20, 5, 2, "Prodotto difettoso", "da_verificare", "2026-06-01", 1)

	mock.ExpectQuery(`SELECT id, ordine_uscita_id, prodotto_id, quantita, motivo, stato, data, created_by FROM Resi`).
		WillReturnRows(rows)

	list, err := repo.OrderRepository_GetAllResi(context.Background())
	require.NoError(t, err)
	require.Len(t, list, 1)
	assert.Equal(t, 3, list[0]["id"])
	assert.Equal(t, "da_verificare", list[0]["stato"])
	require.NoError(t, mock.ExpectationsWereMet())
}