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
// TEST LOGIN
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

func TestIntegration_Login_WrongPassword(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	r := NewRouter(db)

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

	body := `{"username":"admin","password":"wrongpassword"}`
	req := httptest.NewRequest("POST", "/api/v1/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

// ============================================================
// TEST TOKEN
// ============================================================

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
// TEST PRODUCTS
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

func TestIntegration_CreateProduct_Success(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	r := NewRouter(db)
	token := setupAuthMock(mock)

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO Prodotti").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO Giacenze").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	body := `{"codice":"PRD-001","barcode":"1234567890128","nome":"Nuovo Prodotto","descrizione":"Test","prezzo_acquisto":10,"prezzo_vendita":15,"soglia_minima":5,"categoria_id":1}`
	req := httptest.NewRequest("POST", "/api/v1/products", bytes.NewBufferString(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	require.Equal(t, http.StatusOK, rr.Code, "body: %s", rr.Body.String())
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

func TestIntegration_DeleteProduct_InvalidID(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	r := NewRouter(db)
	token := setupAuthMock(mock)

	req := httptest.NewRequest("DELETE", "/api/v1/products/abc", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestIntegration_GetProductByBarcode_NotFound(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	r := NewRouter(db)
	token := setupAuthMock(mock)

	rows := sqlmock.NewRows([]string{
		"id", "codice", "barcode", "nome", "descrizione",
		"categoria_id", "prezzo_acquisto", "prezzo_vendita", "soglia_minima",
		"categoria_nome", "categoria_descrizione",
	})
	mock.ExpectQuery("SELECT .+ FROM Prodotti p").
		WithArgs("9999999999999").
		WillReturnRows(rows)

	req := httptest.NewRequest("GET", "/api/v1/products/barcode/9999999999999", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNotFound, rr.Code)
}

// ============================================================
// TEST ORDERS
// ============================================================

func TestIntegration_GetOrdiniUscita_Success(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	r := NewRouter(db)
	token := setupAuthMock(mock)

	rows := sqlmock.NewRows([]string{
		"id", "numero", "cliente", "data", "stato", "note", "created_by", "ddt_id",
	})
	mock.ExpectQuery("SELECT id, numero, cliente, data, stato, note, created_by, ddt_id FROM OrdiniUscita").
		WillReturnRows(rows)

	req := httptest.NewRequest("GET", "/api/v1/orders", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestIntegration_GetOrdiniEntrata_Success(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	r := NewRouter(db)
	token := setupAuthMock(mock)

	rows := sqlmock.NewRows([]string{
		"id", "numero", "fornitore", "data", "note", "created_by",
	})
	mock.ExpectQuery("SELECT id, numero, fornitore, data, note, created_by FROM OrdiniEntrata").
		WillReturnRows(rows)

	req := httptest.NewRequest("GET", "/api/v1/orders/entrata", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestIntegration_GetInvoices_Success(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	r := NewRouter(db)
	token := setupAuthMock(mock)

	rows := sqlmock.NewRows([]string{
		"id", "numero", "anno", "tipo", "ordine_uscita_id",
		"imponibile", "iva", "totale", "data", "path_pdf",
	})
	mock.ExpectQuery("SELECT id, numero, anno, tipo, ordine_uscita_id").
		WillReturnRows(rows)

	req := httptest.NewRequest("GET", "/api/v1/orders/invoices", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	require.Equal(t, http.StatusOK, rr.Code, "body: %s", rr.Body.String())
}

func TestIntegration_GetCreditNotes_Success(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	r := NewRouter(db)
	token := setupAuthMock(mock)

	rows := sqlmock.NewRows([]string{
		"id", "numero", "anno", "tipo", "ordine_uscita_id",
		"imponibile", "iva", "totale", "data", "path_pdf",
	}).AddRow(1, "NC-001/2026", 2026, "nota_credito", 1, 100.0, 22.0, 122.0, time.Now(), "")
	mock.ExpectQuery("SELECT id, numero, anno, tipo, ordine_uscita_id").
		WillReturnRows(rows)

	req := httptest.NewRequest("GET", "/api/v1/orders/credit-notes", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code, "body: %s", rr.Body.String())
	var result []map[string]interface{}
	err := json.NewDecoder(rr.Body).Decode(&result)
	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, "nota_credito", result[0]["tipo"])
}

func TestIntegration_CreateCreditNote_InvalidID(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	r := NewRouter(db)
	token := setupAuthMock(mock)

	req := httptest.NewRequest("POST", "/api/v1/orders/create-credit-notes/abc", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestIntegration_DeleteFattura_InvalidID(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	r := NewRouter(db)
	token := setupAuthMock(mock)

	req := httptest.NewRequest("DELETE", "/api/v1/orders/delete-fattura/abc", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestIntegration_DeleteEntrata_InvalidID(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	r := NewRouter(db)
	token := setupAuthMock(mock)

	req := httptest.NewRequest("DELETE", "/api/v1/orders/delete-entrata/abc", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestIntegration_DeleteEntrata_Success(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	r := NewRouter(db)
	token := setupAuthMock(mock)

	mock.ExpectQuery("SELECT id, numero, fornitore, data, note, created_by FROM OrdiniEntrata WHERE id").
		WithArgs(5).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "numero", "fornitore", "data", "note", "created_by",
		}).AddRow(5, "OE-005", "Fornitore", time.Now(), "", 1))
	mock.ExpectQuery("SELECT id, ordine_entrata_id, prodotto_id, quantita, prezzo_unitario FROM RigheOrdineEntrata WHERE ordine_entrata_id").
		WithArgs(5).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "ordine_entrata_id", "prodotto_id", "quantita", "prezzo_unitario",
		}))
	mock.ExpectExec("DELETE FROM OrdiniEntrata WHERE id").
		WithArgs(5).
		WillReturnResult(sqlmock.NewResult(0, 1))

	req := httptest.NewRequest("DELETE", "/api/v1/orders/delete-entrata/5", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNoContent, rr.Code, "body: %s", rr.Body.String())
}

func TestIntegration_DeleteEntrata_ControllerError(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	r := NewRouter(db)
	token := setupAuthMock(mock)

	mock.ExpectQuery("SELECT id, numero, fornitore, data, note, created_by FROM OrdiniEntrata WHERE id").
		WithArgs(999).
		WillReturnError(fmt.Errorf("ordine non trovato"))

	req := httptest.NewRequest("DELETE", "/api/v1/orders/delete-entrata/999", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

// ============================================================
// TEST CATEGORIES
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

func TestIntegration_CreateCategory_Success(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	r := NewRouter(db)
	token := setupAuthMock(mock)

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO Categorie \\(nome, descrizione\\)").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	body := `{"nome":"Bevande","descrizione":"Bibite varie"}`
	req := httptest.NewRequest("POST", "/api/v1/categories", bytes.NewBufferString(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	require.Equal(t, http.StatusOK, rr.Code, "body: %s", rr.Body.String())
}

func TestIntegration_DeleteCategory_Success(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	r := NewRouter(db)
	token := setupAuthMock(mock)

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM Categorie WHERE id").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	req := httptest.NewRequest("DELETE", "/api/v1/categories/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	require.Equal(t, http.StatusNoContent, rr.Code, "body: %s", rr.Body.String())
}

func TestIntegration_DeleteCategory_InvalidID(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	r := NewRouter(db)
	token := setupAuthMock(mock)

	req := httptest.NewRequest("DELETE", "/api/v1/categories/abc", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

// ============================================================
// TEST ROUTE ESISTENZA
// ============================================================

// ============================================================
// TEST PRODUCTS - ERROR PATHS & UNTESTED HANDLERS
// ============================================================

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

func TestIntegration_CreateProduct_BadBody(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	r := NewRouter(db)
	token := setupAuthMock(mock)

	req := httptest.NewRequest("POST", "/api/v1/products", bytes.NewBufferString("not-json"))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestIntegration_CreateProduct_ControllerError(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	r := NewRouter(db)
	token := setupAuthMock(mock)

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO Prodotti").WillReturnError(fmt.Errorf("unique violation"))
	mock.ExpectRollback()

	body := `{"codice":"PRD-001","barcode":"1234567890128","nome":"Nuovo","descrizione":"Test","prezzo_acquisto":10,"prezzo_vendita":15,"soglia_minima":5,"categoria_id":1}`
	req := httptest.NewRequest("POST", "/api/v1/products", bytes.NewBufferString(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
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

	// GetByID outside tx (from UpdateSellPrice controller)
	mock.ExpectQuery("SELECT p\\.id, p\\.codice").WithArgs(1).WillReturnRows(makeRow())
	// BeginTx for UpdateListinoPrices
	mock.ExpectBegin()
	// GetByID inside tx
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

func TestIntegration_UpdatePrice_InvalidID(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	r := NewRouter(db)
	token := setupAuthMock(mock)

	body := `{"prezzo_acquisto":10.0,"prezzo_vendita":25.0}`
	req := httptest.NewRequest("PUT", "/api/v1/products/abc/price", bytes.NewBufferString(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestIntegration_UpdatePrice_BadBody(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	r := NewRouter(db)
	token := setupAuthMock(mock)

	req := httptest.NewRequest("PUT", "/api/v1/products/1/price", bytes.NewBufferString("not-json"))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestIntegration_DeleteProduct_ControllerError(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	r := NewRouter(db)
	token := setupAuthMock(mock)

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM Giacenze WHERE prodotto_id").WillReturnError(fmt.Errorf("fk constraint"))
	mock.ExpectRollback()

	req := httptest.NewRequest("DELETE", "/api/v1/products/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestIntegration_GetProductByBarcode_Success(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	r := NewRouter(db)
	token := setupAuthMock(mock)

	rows := sqlmock.NewRows([]string{
		"id", "codice", "barcode", "nome", "descrizione",
		"categoria_id", "prezzo_acquisto", "prezzo_vendita", "soglia_minima",
	}).AddRow(1, "PRD-001", "1234567890128", "Prodotto Test", "Desc",
		1, 10.0, 15.0, 5)
	mock.ExpectQuery("SELECT .+ FROM Prodotti p WHERE p\\.barcode").WithArgs("1234567890128").
		WillReturnRows(rows)

	req := httptest.NewRequest("GET", "/api/v1/products/barcode/1234567890128", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	require.Equal(t, http.StatusOK, rr.Code, "body: %s", rr.Body.String())
	assert.Contains(t, rr.Header().Get("Content-Type"), "application/json")

	var result map[string]interface{}
	err := json.NewDecoder(rr.Body).Decode(&result)
	require.NoError(t, err)
	assert.Equal(t, "Prodotto Test", result["nome"])
}

func TestIntegration_InventoryReport_Success(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	r := NewRouter(db)
	token := setupAuthMock(mock)

	prodRows := sqlmock.NewRows([]string{
		"id", "codice", "barcode", "nome", "descrizione",
		"categoria_id", "prezzo_acquisto", "prezzo_vendita", "soglia_minima",
		"categoria_nome", "categoria_descrizione",
	}).AddRow(1, "P001", "1234567890128", "Vite", "Desc",
		sql.NullInt64{Int64: 1, Valid: true}, 10.0, 15.0, 5,
		sql.NullString{String: "Alimentari", Valid: true}, sql.NullString{String: "Desc", Valid: true})
	mock.ExpectQuery("SELECT p\\.id, p\\.codice").WillReturnRows(prodRows)

	giacRows := sqlmock.NewRows([]string{
		"g.id", "g.prodotto_id", "g.quantita", "g.updated_at",
		"p.nome", "p.soglia_minima",
	}).AddRow(1, 1, 30, time.Now(), "Vite", 100)
	mock.ExpectQuery("SELECT g\\.id, g\\.prodotto_id").WillReturnRows(giacRows)

	req := httptest.NewRequest("GET", "/api/v1/products/inventory-report", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	require.Equal(t, http.StatusOK, rr.Code, "body: %s", rr.Body.String())
}

func TestIntegration_InventoryReport_Error(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	r := NewRouter(db)
	token := setupAuthMock(mock)

	mock.ExpectQuery("SELECT p\\.id, p\\.codice").WillReturnError(fmt.Errorf("db error"))

	req := httptest.NewRequest("GET", "/api/v1/products/inventory-report", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

// ============================================================
// TEST ORDERS - UNTESTED HANDLERS
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

func TestIntegration_CreaOrdineUscita_BadBody(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	r := NewRouter(db)
	token := setupAuthMock(mock)

	req := httptest.NewRequest("POST", "/api/v1/orders/crea-uscita", bytes.NewBufferString("not-json"))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestIntegration_CreaOrdineEntrata_Success(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	r := NewRouter(db)
	token := setupAuthMock(mock)

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO OrdiniEntrata").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO RigheOrdineEntrata").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO Giacenze").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	body := `{"numero":"OE-001","fornitore":"Fornitore SRL","note":"test","righe":[{"prodotto_id":1,"quantita":10,"prezzo_unitario":3.0}]}`
	req := httptest.NewRequest("POST", "/api/v1/orders/crea-entrata", bytes.NewBufferString(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	require.Equal(t, http.StatusCreated, rr.Code, "body: %s", rr.Body.String())
}

func TestIntegration_CreaOrdineEntrata_BadBody(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	r := NewRouter(db)
	token := setupAuthMock(mock)

	req := httptest.NewRequest("POST", "/api/v1/orders/crea-entrata", bytes.NewBufferString("not-json"))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestIntegration_GetOrdiniUscita_Error(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	r := NewRouter(db)
	token := setupAuthMock(mock)

	mock.ExpectQuery("SELECT id, numero, cliente, data, stato, note, created_by, ddt_id FROM OrdiniUscita").
		WillReturnError(fmt.Errorf("db error"))

	req := httptest.NewRequest("GET", "/api/v1/orders", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestIntegration_GetOrdiniUscita_Success_WithDetail(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	r := NewRouter(db)
	token := setupAuthMock(mock)

	// GetOrdiniUscita
	mock.ExpectQuery("SELECT id, numero, cliente, data, stato, note, created_by, ddt_id FROM OrdiniUscita").
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "numero", "cliente", "data", "stato", "note", "created_by", "ddt_id",
		}).AddRow(1, "OU-001", "Acme", time.Now(), "confermato", "", 1, nil))

	// GetOrdiniUscitaById (detail)
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

	req := httptest.NewRequest("GET", "/api/v1/orders", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	require.Equal(t, http.StatusOK, rr.Code, "body: %s", rr.Body.String())
}

func TestIntegration_GetOrdiniUscita_Error_WithDetail(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	r := NewRouter(db)
	token := setupAuthMock(mock)

	mock.ExpectQuery("SELECT id, numero, cliente, data, stato, note, created_by, ddt_id FROM OrdiniUscita").
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "numero", "cliente", "data", "stato", "note", "created_by", "ddt_id",
		}).AddRow(1, "OU-001", "Acme", time.Now(), "confermato", "", 1, nil))

	mock.ExpectQuery("SELECT id, numero, cliente, data, stato, note, created_by, ddt_id FROM OrdiniUscita WHERE id").
		WithArgs(1).
		WillReturnError(fmt.Errorf("db error"))

	req := httptest.NewRequest("GET", "/api/v1/orders", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestIntegration_GetOrdiniUscitaViaEndpoint_Success(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	r := NewRouter(db)
	token := setupAuthMock(mock)

	mock.ExpectQuery("SELECT id, numero, cliente, data, stato, note, created_by, ddt_id FROM OrdiniUscita").
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "numero", "cliente", "data", "stato", "note", "created_by", "ddt_id",
		}))

	req := httptest.NewRequest("GET", "/api/v1/orders/uscita", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestIntegration_GetOrdiniEntrata_Error(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	r := NewRouter(db)
	token := setupAuthMock(mock)

	mock.ExpectQuery("SELECT id, numero, fornitore, data, note, created_by FROM OrdiniEntrata").
		WillReturnError(fmt.Errorf("db error"))

	req := httptest.NewRequest("GET", "/api/v1/orders/entrata", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestIntegration_UpdateOrdineUscita_Success(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	r := NewRouter(db)
	token := setupAuthMock(mock)

	// GetOrdiniUscitaById (outside tx, from controller)
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

	// BeginTx for UpdateOrdineUscita service
	mock.ExpectBegin()
	giacQuery := mock.ExpectQuery("SELECT id, prodotto_id, quantita, updated_at FROM Giacenze WHERE prodotto_id").
		WithArgs(1)
	giacQuery.WillReturnRows(sqlmock.NewRows([]string{"id", "prodotto_id", "quantita", "updated_at"}).
		AddRow(1, 1, 100, time.Now()))
	mock.ExpectExec("INSERT INTO Giacenze").WillReturnResult(sqlmock.NewResult(0, 1))
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

func TestIntegration_UpdateOrdineUscita_InvalidID(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	r := NewRouter(db)
	token := setupAuthMock(mock)

	body := `{"numero":"OU-001"}`
	req := httptest.NewRequest("PUT", "/api/v1/orders/update-ordine-uscita/abc", bytes.NewBufferString(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestIntegration_UpdateOrdineUscita_BadBody(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	r := NewRouter(db)
	token := setupAuthMock(mock)

	req := httptest.NewRequest("PUT", "/api/v1/orders/update-ordine-uscita/1", bytes.NewBufferString("not-json"))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestIntegration_GetResi_Error(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	r := NewRouter(db)
	token := setupAuthMock(mock)

	mock.ExpectQuery("SELECT id, ordine_uscita_id").WillReturnError(fmt.Errorf("db error"))

	req := httptest.NewRequest("GET", "/api/v1/orders/resi", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestIntegration_GetResi_Success(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	r := NewRouter(db)
	token := setupAuthMock(mock)

	mock.ExpectQuery("SELECT id, ordine_uscita_id").
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "ordine_uscita_id", "prodotto_id", "quantita",
			"motivo", "stato", "data", "created_by",
		}).AddRow(1, 1, 1, 5, "difettoso", "da_verificare", time.Now(), 1))

	req := httptest.NewRequest("GET", "/api/v1/orders/resi", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	require.Equal(t, http.StatusOK, rr.Code, "body: %s", rr.Body.String())
}

func TestIntegration_UpdateResi_Success(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	r := NewRouter(db)
	token := setupAuthMock(mock)

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT id, ordine_uscita_id, prodotto_id, quantita, motivo, stato, data, created_by FROM Resi WHERE id").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "ordine_uscita_id", "prodotto_id", "quantita", "motivo", "stato", "data", "created_by",
		}).AddRow(1, 1, 1, 5, "difettoso", "da_verificare", time.Now(), 1))
	mock.ExpectExec("INSERT INTO Giacenze").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("UPDATE Resi SET stato").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	body := `{"id":1,"stato":"da_verificare"}`
	req := httptest.NewRequest("PUT", "/api/v1/orders/update-resi", bytes.NewBufferString(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNoContent, rr.Code, "body: %s", rr.Body.String())
}

func TestIntegration_UpdateResi_BadBody(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	r := NewRouter(db)
	token := setupAuthMock(mock)

	req := httptest.NewRequest("PUT", "/api/v1/orders/update-resi", bytes.NewBufferString("not-json"))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestIntegration_CreateCreditNote_Success(t *testing.T) {
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

	req := httptest.NewRequest("POST", "/api/v1/orders/create-credit-notes/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	require.Equal(t, http.StatusOK, rr.Code, "body: %s", rr.Body.String())
}

func TestIntegration_CreateCreditNote_ControllerError(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	r := NewRouter(db)
	token := setupAuthMock(mock)

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT id, numero, cliente, data, stato, note, created_by, ddt_id FROM OrdiniUscita WHERE id").
		WillReturnError(fmt.Errorf("ordine non trovato"))
	mock.ExpectRollback()

	req := httptest.NewRequest("POST", "/api/v1/orders/create-credit-notes/999", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
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

func TestIntegration_CreateDDT_InvalidID(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	r := NewRouter(db)
	token := setupAuthMock(mock)

	body := `{"mittente":"SRL","destinatario":"SRL"}`
	req := httptest.NewRequest("POST", "/api/v1/orders/abc/create-ddt", bytes.NewBufferString(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestIntegration_CreateDDT_BadBody(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	r := NewRouter(db)
	token := setupAuthMock(mock)

	req := httptest.NewRequest("POST", "/api/v1/orders/1/create-ddt", bytes.NewBufferString("not-json"))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

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

func TestIntegration_CreateFattura_InvalidID(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	r := NewRouter(db)
	token := setupAuthMock(mock)

	req := httptest.NewRequest("POST", "/api/v1/orders/create-fattura/abc", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestIntegration_CreateFattura_ControllerError(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	r := NewRouter(db)
	token := setupAuthMock(mock)

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT id, numero, cliente, data, stato, note, created_by, ddt_id FROM OrdiniUscita WHERE id").
		WillReturnError(fmt.Errorf("ordine non trovato"))
	mock.ExpectRollback()

	req := httptest.NewRequest("POST", "/api/v1/orders/create-fattura/999", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestIntegration_DeleteFattura_Success(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	r := NewRouter(db)
	token := setupAuthMock(mock)

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM Fatture WHERE id").
		WithArgs(5).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	req := httptest.NewRequest("DELETE", "/api/v1/orders/delete-fattura/5", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNoContent, rr.Code)
}

func TestIntegration_DeleteFattura_ControllerError(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	r := NewRouter(db)
	token := setupAuthMock(mock)

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM Fatture WHERE id").
		WillReturnError(fmt.Errorf("fk constraint"))
	mock.ExpectRollback()

	req := httptest.NewRequest("DELETE", "/api/v1/orders/delete-fattura/5", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestIntegration_GetInvoices_Error(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	r := NewRouter(db)
	token := setupAuthMock(mock)

	mock.ExpectQuery("SELECT id, numero, anno, tipo, ordine_uscita_id").
		WillReturnError(fmt.Errorf("db error"))

	req := httptest.NewRequest("GET", "/api/v1/orders/invoices", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestIntegration_GetCreditNotes_Error(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	r := NewRouter(db)
	token := setupAuthMock(mock)

	mock.ExpectQuery("SELECT id, numero, anno, tipo, ordine_uscita_id").
		WillReturnError(fmt.Errorf("db error"))

	req := httptest.NewRequest("GET", "/api/v1/orders/credit-notes", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

// ============================================================
// TEST CATEGORIES - ERROR PATHS
// ============================================================

func TestIntegration_GetCategories_Error(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	r := NewRouter(db)
	token := setupAuthMock(mock)

	mock.ExpectQuery("SELECT id, nome, descrizione FROM Categorie ORDER BY nome").
		WillReturnError(fmt.Errorf("db error"))

	req := httptest.NewRequest("GET", "/api/v1/categories", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestIntegration_CreateCategory_BadBody(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	r := NewRouter(db)
	token := setupAuthMock(mock)

	req := httptest.NewRequest("POST", "/api/v1/categories", bytes.NewBufferString("not-json"))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestIntegration_CreateCategory_ControllerError(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	r := NewRouter(db)
	token := setupAuthMock(mock)

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO Categorie \\(nome, descrizione\\)").
		WillReturnError(fmt.Errorf("duplicate entry"))
	mock.ExpectRollback()

	body := `{"nome":"Alimentari","descrizione":"Duplicate"}`
	req := httptest.NewRequest("POST", "/api/v1/categories", bytes.NewBufferString(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestIntegration_DeleteCategory_ControllerError(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	r := NewRouter(db)
	token := setupAuthMock(mock)

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM Categorie WHERE id").
		WillReturnError(fmt.Errorf("fk constraint"))
	mock.ExpectRollback()

	req := httptest.NewRequest("DELETE", "/api/v1/categories/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

// ============================================================
// TEST USERS - ENTIRE CRUD SECTION
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

func TestIntegration_CreateUser_BadBody(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	r := NewRouter(db)
	token := setupAuthMock(mock)

	req := httptest.NewRequest("POST", "/api/v1/users", bytes.NewBufferString("not-json"))
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestIntegration_CreateUser_DBError(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	r := NewRouter(db)
	token := setupAuthMock(mock)

	mock.ExpectExec("INSERT INTO Utenti").WillReturnError(fmt.Errorf("duplicate username"))

	body := `{"username":"esistente","password":"secret","ruolo":"admin"}`
	req := httptest.NewRequest("POST", "/api/v1/users", bytes.NewBufferString(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestIntegration_GetUser_Success(t *testing.T) {
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

	req := httptest.NewRequest("GET", "/api/v1/users/mario", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	require.Equal(t, http.StatusOK, rr.Code, "body: %s", rr.Body.String())

	var resp map[string]interface{}
	err := json.NewDecoder(rr.Body).Decode(&resp)
	require.NoError(t, err)
	assert.Equal(t, "mario", resp["Username"])
}

func TestIntegration_GetUser_NotFound(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	r := NewRouter(db)
	token := setupAuthMock(mock)

	mock.ExpectQuery("SELECT id, username, password_hash, ruolo, tentativi_falliti, bloccato, created_at").
		WithArgs("inesistente").
		WillReturnError(fmt.Errorf("credenziali errate"))

	req := httptest.NewRequest("GET", "/api/v1/users/inesistente", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestIntegration_UpdateUser_Success(t *testing.T) {
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

	mock.ExpectExec("UPDATE Utenti SET username").
		WithArgs("mario", sqlmock.AnyArg(), "contabile", 0, false, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	body := `{"ruolo":"contabile"}`
	req := httptest.NewRequest("PUT", "/api/v1/users/mario", bytes.NewBufferString(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	require.Equal(t, http.StatusOK, rr.Code, "body: %s", rr.Body.String())
}

func TestIntegration_UpdateUser_NotFound(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	r := NewRouter(db)
	token := setupAuthMock(mock)

	mock.ExpectQuery("SELECT id, username, password_hash, ruolo, tentativi_falliti, bloccato, created_at").
		WithArgs("inesistente").
		WillReturnError(fmt.Errorf("credenziali errate"))

	body := `{"ruolo":"admin"}`
	req := httptest.NewRequest("PUT", "/api/v1/users/inesistente", bytes.NewBufferString(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestIntegration_UpdateUser_BadBody(t *testing.T) {
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

	req := httptest.NewRequest("PUT", "/api/v1/users/mario", bytes.NewBufferString("not-json"))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
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

func TestIntegration_UpdateUser_PasswordChange(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	r := NewRouter(db)
	token := setupAuthMock(mock)

	rows := sqlmock.NewRows([]string{
		"id", "username", "password_hash", "ruolo",
		"tentativi_falliti", "bloccato", "created_at",
	}).AddRow(1, "mario", "oldhash", "magazziniere", 0, false, time.Now())
	mock.ExpectQuery("SELECT id, username, password_hash, ruolo, tentativi_falliti, bloccato, created_at").
		WithArgs("mario").
		WillReturnRows(rows)

	mock.ExpectExec("UPDATE Utenti SET username").
		WillReturnResult(sqlmock.NewResult(0, 1))

	body := `{"password":"newpassword123"}`
	req := httptest.NewRequest("PUT", "/api/v1/users/mario", bytes.NewBufferString(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	require.Equal(t, http.StatusOK, rr.Code, "body: %s", rr.Body.String())
}

func TestIntegration_UpdateUser_DBError(t *testing.T) {
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

	mock.ExpectExec("UPDATE Utenti SET username").WillReturnError(fmt.Errorf("db error"))

	body := `{"ruolo":"admin"}`
	req := httptest.NewRequest("PUT", "/api/v1/users/mario", bytes.NewBufferString(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
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

func TestIntegration_DeleteUser_NotFound(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	r := NewRouter(db)
	token := setupAuthMock(mock)

	mock.ExpectQuery("SELECT id, username, password_hash, ruolo, tentativi_falliti, bloccato, created_at").
		WithArgs("inesistente").
		WillReturnError(fmt.Errorf("credenziali errate"))

	req := httptest.NewRequest("DELETE", "/api/v1/users/inesistente", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestIntegration_DeleteUser_DeleteError(t *testing.T) {
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

	mock.ExpectExec("DELETE FROM Utenti WHERE id").WillReturnError(fmt.Errorf("fk constraint"))

	req := httptest.NewRequest("DELETE", "/api/v1/users/mario", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestIntegration_RoutesExist(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	r := NewRouter(db)
	token := setupAuthMock(mock)

	mock.ExpectQuery("SELECT p\\.id, p\\.codice").
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "codice", "barcode", "nome", "descrizione",
			"categoria_id", "prezzo_acquisto", "prezzo_vendita", "soglia_minima",
			"categoria_nome", "categoria_descrizione",
		}))

	mock.ExpectQuery("SELECT id, numero, cliente, data, stato, note, created_by, ddt_id FROM OrdiniUscita").
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "numero", "cliente", "data", "stato", "note", "created_by", "ddt_id",
		}))

	mock.ExpectQuery("SELECT id, numero, fornitore, data, note, created_by FROM OrdiniEntrata").
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "numero", "fornitore", "data", "note", "created_by",
		}))

	mock.ExpectQuery("SELECT id, nome, descrizione FROM Categorie ORDER BY nome").
		WillReturnRows(sqlmock.NewRows([]string{"id", "nome", "descrizione"}))

	mock.ExpectQuery("SELECT id, numero, anno, tipo, ordine_uscita_id").
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "numero", "anno", "tipo", "ordine_uscita_id",
			"imponibile", "iva", "totale", "data", "path_pdf",
		}))

	mock.ExpectQuery("SELECT id, numero, anno, tipo, ordine_uscita_id").
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "numero", "anno", "tipo", "ordine_uscita_id",
			"imponibile", "iva", "totale", "data", "path_pdf",
		}))

	mock.ExpectQuery("SELECT id, ordine_uscita_id").
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "ordine_uscita_id", "prodotto_id", "quantita",
			"motivo", "stato", "data", "created_by",
		}))

	routes := []struct {
		method string
		path   string
	}{
		{"GET", "/api/v1/products"},
		{"GET", "/api/v1/orders"},
		{"GET", "/api/v1/orders/uscita"},
		{"GET", "/api/v1/orders/entrata"},
		{"GET", "/api/v1/orders/invoices"},
		{"GET", "/api/v1/orders/credit-notes"},
		{"GET", "/api/v1/orders/resi"},
		{"GET", "/api/v1/categories"},
	}

	for _, rt := range routes {
		t.Run(fmt.Sprintf("%s %s should not 404/405", rt.method, rt.path), func(t *testing.T) {
			req := httptest.NewRequest(rt.method, rt.path, nil)
			req.Header.Set("Authorization", "Bearer "+token)
			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)
			assert.NotEqual(t, http.StatusNotFound, rr.Code, "route should be registered")
			assert.NotEqual(t, http.StatusMethodNotAllowed, rr.Code, "method should be allowed")
		})
	}
}
