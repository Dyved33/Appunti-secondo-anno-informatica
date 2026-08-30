// Package router — test di integrazione end-to-end.
//
// Questi test NON testano singole unità (repository/service isolati), ma
// l'intera catena che il client TUI (frontend) attraversa quando chiama
// l'API: richiesta HTTP grezza -> router -> middleware di autenticazione
// -> controller -> service -> repository -> DB (qui mockato con sqlmock).
//
// L'obiettivo è verificare il CONTRATTO tra frontend e backend: path,
// metodo HTTP, status code e forma del body JSON restituito — esattamente
// ciò che il client si aspetta di ricevere quando parla con il server reale.
//
// Il set è stato ridotto a un numero limitato di scenari rappresentativi
// (uno per comportamento significativo), evitando di ripetere varianti
// quasi identiche dello stesso endpoint.
package router

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tommasogarofani/ingsw_25-26_group-11/server/auth"
)

func newMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	t.Helper()
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	return db, mock
}

func setupAuthMock(mock sqlmock.Sqlmock) string {
	token, _ := auth.GenerateToken(1, "admin", "admin")
	tokenHash := auth.HashToken(token)

	rows := sqlmock.NewRows([]string{"count"}).AddRow(1)
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM Sessioni").
		WithArgs(tokenHash).
		WillReturnRows(rows)

	return token
}

func setupLoginMock(mock sqlmock.Sqlmock) {
	mock.ExpectBegin()
	rows := sqlmock.NewRows([]string{
		"id", "username", "password_hash", "ruolo",
		"tentativi_falliti", "bloccato", "created_at",
	}).AddRow(1, "admin", "$2a$10$BA/ajnx3BUxjfhUj7cu11.DRz.vSPhOhkxtR/KMQVfHJsxsfr3gL.", "admin", 0, false, time.Now())
	mock.ExpectQuery("SELECT .+ FROM Utenti WHERE username").
		WithArgs("admin").
		WillReturnRows(rows)
	mock.ExpectExec("UPDATE Utenti SET").
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	mock.ExpectExec("INSERT INTO Sessioni").
		WillReturnResult(sqlmock.NewResult(1, 1))
}

// ============================================================
// AUTENTICAZIONE — contratto login/protezione endpoint che il
// client usa per ottenere e validare il token.
// ============================================================

func TestIntegration_Login_Success(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	r := NewRouter(db)
	setupLoginMock(mock)

	body := `{"username":"admin","password":"password123"}`
	req := httptest.NewRequest("POST", "/api/v1/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code, "body: %s", rr.Body.String())
	assert.Contains(t, rr.Header().Get("Content-Type"), "application/json")

	var resp map[string]interface{}
	err := json.NewDecoder(rr.Body).Decode(&resp)
	require.NoError(t, err)
	assert.NotEmpty(t, resp["token"])
	assert.Equal(t, "admin", resp["role"])
}

func TestIntegration_Login_BadRequest(t *testing.T) {
	db, _ := newMockDB(t)
	defer db.Close()

	r := NewRouter(db)

	req := httptest.NewRequest("POST", "/api/v1/login", bytes.NewBufferString("not-json"))
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestIntegration_Protected_NoToken(t *testing.T) {
	db, _ := newMockDB(t)
	defer db.Close()

	r := NewRouter(db)

	endpoints := []struct {
		method string
		path   string
	}{
		{"GET", "/api/v1/products"},
		{"GET", "/api/v1/orders"},
		{"GET", "/api/v1/orders/entrata"},
		{"GET", "/api/v1/orders/invoices"},
		{"GET", "/api/v1/orders/credit-notes"},
		{"GET", "/api/v1/orders/resi"},
		{"GET", "/api/v1/categories"},
	}

	for _, ep := range endpoints {
		t.Run(fmt.Sprintf("%s %s senza token", ep.method, ep.path), func(t *testing.T) {
			req := httptest.NewRequest(ep.method, ep.path, nil)
			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)
			assert.Equal(t, http.StatusUnauthorized, rr.Code)
		})
	}
}

func TestIntegration_Protected_InvalidToken(t *testing.T) {
	db, _ := newMockDB(t)
	defer db.Close()

	r := NewRouter(db)

	req := httptest.NewRequest("GET", "/api/v1/products", nil)
	req.Header.Set("Authorization", "Bearer invalid-token-string")
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

// ============================================================
// PRODOTTI — lettura, errore a valle, ed eliminazione in
// transazione (comportamento rappresentativo delle write path).
// ============================================================

func TestIntegration_GetProducts_Success(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	r := NewRouter(db)
	token := setupAuthMock(mock)

	rows := sqlmock.NewRows([]string{
		"id", "codice", "barcode", "nome", "descrizione",
		"categoria_id", "prezzo_acquisto", "prezzo_vendita", "soglia_minima",
		"categoria_nome", "categoria_descrizione",
	}).AddRow(1, "PRD-001", "1234567890128", "Prodotto Test", "Desc",
		sql.NullInt64{Int64: 1, Valid: true}, 10.0, 15.0, 5,
		sql.NullString{String: "Alimentari", Valid: true}, sql.NullString{String: "Desc Cat", Valid: true})
	mock.ExpectQuery("SELECT p\\.id, p\\.codice").
		WillReturnRows(rows)

	req := httptest.NewRequest("GET", "/api/v1/products", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code, "body: %s", rr.Body.String())
	assert.Contains(t, rr.Header().Get("Content-Type"), "application/json")

	var result []map[string]interface{}
	err := json.NewDecoder(rr.Body).Decode(&result)
	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, "Prodotto Test", result[0]["nome"])
}

func TestIntegration_GetProducts_Error(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	r := NewRouter(db)
	token := setupAuthMock(mock)

	mock.ExpectQuery("SELECT p\\.id, p\\.codice").WillReturnError(fmt.Errorf("db down"))

	req := httptest.NewRequest("GET", "/api/v1/products", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestIntegration_DeleteProduct_Success(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	r := NewRouter(db)
	token := setupAuthMock(mock)

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM Giacenze WHERE prodotto_id").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("DELETE FROM Prodotti WHERE id").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	req := httptest.NewRequest("DELETE", "/api/v1/products/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	require.Equal(t, http.StatusNoContent, rr.Code, "body: %s", rr.Body.String())
}

func TestIntegration_UpdatePrice_Success(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	r := NewRouter(db)
	token := setupAuthMock(mock)

	productCols := []string{
		"id", "codice", "barcode", "nome", "descrizione",
		"categoria_id", "prezzo_acquisto", "prezzo_vendita", "soglia_minima",
		"categoria_nome", "categoria_descrizione",
	}
	makeRow := func() *sqlmock.Rows {
		return sqlmock.NewRows(productCols).AddRow(1, "PRD-001", "1234567890128", "Prodotto", "Desc",
			sql.NullInt64{Int64: 1, Valid: true}, 10.0, 15.0, 5,
			sql.NullString{String: "Alimentari", Valid: true}, sql.NullString{String: "", Valid: true})
	}

	// GetByID fuori dalla tx (dal controller UpdateSellPrice)
	mock.ExpectQuery("SELECT p\\.id, p\\.codice").WithArgs(1).WillReturnRows(makeRow())
	// BeginTx per UpdateListinoPrices
	mock.ExpectBegin()
	// GetByID dentro la tx
	mock.ExpectQuery("SELECT p\\.id, p\\.codice").WithArgs(1).WillReturnRows(makeRow())
	mock.ExpectExec("UPDATE Prodotti SET").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	body := `{"prezzo_acquisto":10.0,"prezzo_vendita":25.0}`
	req := httptest.NewRequest("PUT", "/api/v1/products/1/price", bytes.NewBufferString(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNoContent, rr.Code, "body: %s", rr.Body.String())
}

// ============================================================
// ORDINI — logica di business più delicata: creazione,
// modifica di stato e cancellazione con effetti sulla giacenza.
// ============================================================

func TestIntegration_CreaOrdineUscita_Success(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	r := NewRouter(db)
	token := setupAuthMock(mock)

	mock.ExpectBegin()
	giacQuery := mock.ExpectQuery("SELECT id, prodotto_id, quantita, updated_at FROM Giacenze WHERE prodotto_id").
		WithArgs(1)
	giacQuery.WillReturnRows(sqlmock.NewRows([]string{"id", "prodotto_id", "quantita", "updated_at"}).
		AddRow(1, 1, 100, time.Now()))
	mock.ExpectExec("INSERT INTO Giacenze").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("INSERT INTO OrdiniUscita").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO RigheOrdineUscita").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	body := `{"numero":"OU-001","cliente":"Acme","note":"test","righe":[{"prodotto_id":1,"quantita":10,"prezzo_unitario":5.0}]}`
	req := httptest.NewRequest("POST", "/api/v1/orders/crea-uscita", bytes.NewBufferString(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	require.Equal(t, http.StatusCreated, rr.Code, "body: %s", rr.Body.String())
}

func TestIntegration_UpdateOrdineUscita_Success(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	r := NewRouter(db)
	token := setupAuthMock(mock)

	// GetOrdiniUscitaById (fuori tx, dal controller)
	mock.ExpectQuery("SELECT id, numero, cliente, data, stato, note, created_by, ddt_id FROM OrdiniUscita WHERE id").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "numero", "cliente", "data", "stato", "note", "created_by", "ddt_id",
		}).AddRow(1, "OU-001", "Acme", time.Now(), "confermato", "", 1, nil))
	mock.ExpectQuery("SELECT id, ordine_uscita_id, prodotto_id, quantita, prezzo_unitario FROM RigheOrdineUscita WHERE ordine_uscita_id").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "ordine_uscita_id", "prodotto_id", "quantita", "prezzo_unitario",
		}))

	// BeginTx per UpdateOrdineUscita service
	mock.ExpectBegin()
	// Il service ricarica l'ordine dentro la transazione per determinare la transizione di stato
	mock.ExpectQuery("SELECT id, numero, cliente, data, stato, note, created_by, ddt_id FROM OrdiniUscita WHERE id").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "numero", "cliente", "data", "stato", "note", "created_by", "ddt_id",
		}).AddRow(1, "OU-001", "Acme", time.Now(), "confermato", "", 1, nil))
	mock.ExpectQuery("SELECT id, ordine_uscita_id, prodotto_id, quantita, prezzo_unitario FROM RigheOrdineUscita WHERE ordine_uscita_id").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "ordine_uscita_id", "prodotto_id", "quantita", "prezzo_unitario",
		}).AddRow(1, 1, 1, 10, 5.0))
	// stato corrente "confermato" -> nuovo "evaso": nessuna variazione di giacenza
	mock.ExpectExec("UPDATE OrdiniUscita SET").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	body := `{"numero":"OU-001","cliente":"Acme","note":"test","stato":"evaso","righe":[{"prodotto_id":1,"quantita":10,"prezzo_unitario":5.0}]}`
	req := httptest.NewRequest("PUT", "/api/v1/orders/update-ordine-uscita/1", bytes.NewBufferString(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNoContent, rr.Code, "body: %s", rr.Body.String())
}

func TestIntegration_DeleteUscita_RipristinaGiacenza(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	r := NewRouter(db)
	token := setupAuthMock(mock)

	// Il controller carica l'ordine con le sue righe (fuori transazione)
	mock.ExpectQuery("SELECT id, numero, cliente, data, stato, note, created_by, ddt_id FROM OrdiniUscita WHERE id").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "numero", "cliente", "data", "stato", "note", "created_by", "ddt_id",
		}).AddRow(1, "OU-001", "Acme", time.Now(), "confermato", "", 1, nil))
	mock.ExpectQuery("SELECT id, ordine_uscita_id, prodotto_id, quantita, prezzo_unitario FROM RigheOrdineUscita WHERE ordine_uscita_id").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "ordine_uscita_id", "prodotto_id", "quantita", "prezzo_unitario",
		}).AddRow(1, 1, 1, 10, 5.0))

	// Il service elimina in transazione ripristinando le giacenze (incremento)
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO Giacenze").
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("DELETE FROM OrdiniUscita WHERE id").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	req := httptest.NewRequest("DELETE", "/api/v1/orders/delete-uscita/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNoContent, rr.Code, "body: %s", rr.Body.String())
}

// ============================================================
// DOCUMENTI — generazione fattura e DDT a partire da un ordine
// (join di più tabelle + scrittura, percorso usato dal client
// per i pulsanti "genera fattura" / "genera DDT").
// ============================================================

func TestIntegration_CreateFattura_Success(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	r := NewRouter(db)
	token := setupAuthMock(mock)

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT id, numero, cliente, data, stato, note, created_by, ddt_id FROM OrdiniUscita WHERE id").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "numero", "cliente", "data", "stato", "note", "created_by", "ddt_id",
		}).AddRow(1, "OU-001", "Acme", time.Now(), "evaso", "", 1, nil))
	mock.ExpectQuery("SELECT id, ordine_uscita_id, prodotto_id, quantita, prezzo_unitario FROM RigheOrdineUscita WHERE ordine_uscita_id").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "ordine_uscita_id", "prodotto_id", "quantita", "prezzo_unitario",
		}))
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM Fatture WHERE anno").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	mock.ExpectExec("INSERT INTO Fatture").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	req := httptest.NewRequest("POST", "/api/v1/orders/create-fattura/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	require.Equal(t, http.StatusOK, rr.Code, "body: %s", rr.Body.String())
}

func TestIntegration_CreateDDT_Success(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	r := NewRouter(db)
	token := setupAuthMock(mock)

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT id, numero, cliente, data, stato, note, created_by, ddt_id FROM OrdiniUscita WHERE id").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "numero", "cliente", "data", "stato", "note", "created_by", "ddt_id",
		}).AddRow(1, "OU-001", "Acme", time.Now(), "confermato", "", 1, nil))
	mock.ExpectQuery("SELECT id, ordine_uscita_id, prodotto_id, quantita, prezzo_unitario FROM RigheOrdineUscita WHERE ordine_uscita_id").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "ordine_uscita_id", "prodotto_id", "quantita", "prezzo_unitario",
		}))
	mock.ExpectQuery("SELECT id, numero, ordine_uscita_id, data, mittente, destinatario, path_pdf FROM DDT WHERE numero LIKE").
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "numero", "ordine_uscita_id", "data", "mittente", "destinatario", "path_pdf",
		}))
	mock.ExpectExec("INSERT INTO DDT").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("UPDATE OrdiniUscita SET").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	body := `{"mittente":"Mittente SRL","destinatario":"Dest SRL"}`
	req := httptest.NewRequest("POST", "/api/v1/orders/1/create-ddt", bytes.NewBufferString(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	require.Equal(t, http.StatusOK, rr.Code, "body: %s", rr.Body.String())
}

// ============================================================
// CATEGORIE
// ============================================================

func TestIntegration_GetCategories_Success(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	r := NewRouter(db)
	token := setupAuthMock(mock)

	rows := sqlmock.NewRows([]string{"id", "nome", "descrizione"}).
		AddRow(1, "Alimentari", "Prodotti alimentari").
		AddRow(2, "Non Alimentari", "")
	mock.ExpectQuery("SELECT id, nome, descrizione FROM Categorie ORDER BY nome").
		WillReturnRows(rows)

	req := httptest.NewRequest("GET", "/api/v1/categories", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code, "body: %s", rr.Body.String())

	var result []map[string]interface{}
	err := json.NewDecoder(rr.Body).Decode(&result)
	require.NoError(t, err)
	require.Len(t, result, 2)
	assert.Equal(t, "Alimentari", result[0]["nome"])
}

// ============================================================
// UTENTI — ciclo di vita minimo (creazione, blocco/sblocco,
// eliminazione) usato dalla sezione amministrativa del client.
// ============================================================

func TestIntegration_CreateUser_Success(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	r := NewRouter(db)
	token := setupAuthMock(mock)

	mock.ExpectExec("INSERT INTO Utenti").
		WithArgs("nuovo", sqlmock.AnyArg(), "magazziniere", 0, false, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(10, 1))

	body := `{"username":"nuovo","password":"secret","ruolo":"magazziniere"}`
	req := httptest.NewRequest("POST", "/api/v1/users", bytes.NewBufferString(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	require.Equal(t, http.StatusOK, rr.Code, "body: %s", rr.Body.String())

	var resp map[string]interface{}
	err := json.NewDecoder(rr.Body).Decode(&resp)
	require.NoError(t, err)
	assert.Equal(t, "nuovo", resp["Username"])
}

func TestIntegration_UpdateUser_BlockUnblock(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	r := NewRouter(db)
	token := setupAuthMock(mock)

	rows := sqlmock.NewRows([]string{
		"id", "username", "password_hash", "ruolo",
		"tentativi_falliti", "bloccato", "created_at",
	}).AddRow(1, "mario", "hash123", "magazziniere", 3, true, time.Now())
	mock.ExpectQuery("SELECT id, username, password_hash, ruolo, tentativi_falliti, bloccato, created_at").
		WithArgs("mario").
		WillReturnRows(rows)

	mock.ExpectExec("UPDATE Utenti SET username").
		WithArgs("mario", "hash123", "magazziniere", 5, false, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	body := `{"bloccato":false,"tentativi_falliti":5}`
	req := httptest.NewRequest("PUT", "/api/v1/users/mario", bytes.NewBufferString(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	require.Equal(t, http.StatusOK, rr.Code, "body: %s", rr.Body.String())

	var resp map[string]interface{}
	err := json.NewDecoder(rr.Body).Decode(&resp)
	require.NoError(t, err)
	assert.Equal(t, false, resp["Bloccato"])
	assert.Equal(t, float64(5), resp["TentativiFalliti"])
}

func TestIntegration_DeleteUser_Success(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	r := NewRouter(db)
	token := setupAuthMock(mock)

	rows := sqlmock.NewRows([]string{
		"id", "username", "password_hash", "ruolo",
		"tentativi_falliti", "bloccato", "created_at",
	}).AddRow(1, "mario", "hash123", "magazziniere", 0, false, time.Now())
	mock.ExpectQuery("SELECT id, username, password_hash, ruolo, tentativi_falliti, bloccato, created_at").
		WithArgs("mario").
		WillReturnRows(rows)

	mock.ExpectExec("DELETE FROM Utenti WHERE id").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	req := httptest.NewRequest("DELETE", "/api/v1/users/mario", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNoContent, rr.Code)
}
