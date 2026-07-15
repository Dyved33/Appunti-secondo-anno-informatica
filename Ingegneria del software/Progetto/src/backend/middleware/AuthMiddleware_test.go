package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tommasogarofani/ingsw_25-26_group-11/src/backend/auth"
)

func TestAuthMiddleware(t *testing.T) {
	// Genera un token valido per i test
	tokenValido, err := auth.GenerateToken(123, "mario")
	assert.NoError(t, err)

	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
		verifyContext  bool
	}{
		{
			name:           "Token assente",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Token invalido",
			authHeader:     "Bearer token-non-valido",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Token valido",
			authHeader:     "Bearer " + tokenValido,
			expectedStatus: http.StatusOK,
			verifyContext:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Handler dummy che verifica il contesto se richiesto
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tc.verifyContext {
					userID, ok := r.Context().Value(UserIDKey).(int)
					assert.True(t, ok)
					assert.Equal(t, 123, userID)
				}
				w.WriteHeader(http.StatusOK)
			})

			// Crea la richiesta HTTP
			req := httptest.NewRequest("GET", "/api/v1/protetta", nil)
			if tc.authHeader != "" {
				req.Header.Set("Authorization", tc.authHeader)
			}

			rr := httptest.NewRecorder()
			handler := AuthMiddleware(nextHandler)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code)
		})
	}
}
