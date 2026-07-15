package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRouter(t *testing.T) {
	db, _, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	r := NewRouter(db)
	require.NotNil(t, r)

	// Verifica che un endpoint pubblico (login) risponda
	// In questo caso, inviando un body vuoto o malformato deve restituire 400 Bad Request
	req := httptest.NewRequest("POST", "/api/v1/login", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}
