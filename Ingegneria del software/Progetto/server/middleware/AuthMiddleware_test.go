package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tommasogarofani/ingsw_25-26_group-11/server/auth"
)

type mockSessionValidator struct {
	mock.Mock
}

func (m *mockSessionValidator) SessionRepository_IsValid(ctx context.Context, tokenHash string) (bool, error) {
	args := m.Called(ctx, tokenHash)
	return args.Bool(0), args.Error(1)
}

func TestAuthMiddleware(t *testing.T) {
	// Genera un token valido per i test
	tokenValido, err := auth.GenerateToken(123, "mario", "Amministratore")
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

func TestAuthMiddleware_WithSessionValidator(t *testing.T) {
	tokenValido, err := auth.GenerateToken(123, "mario", "Amministratore")
	assert.NoError(t, err)
	tokenHash := auth.HashToken(tokenValido)

	t.Run("sessione valida DB approvata", func(t *testing.T) {
		sv := new(mockSessionValidator)
		sv.On("SessionRepository_IsValid", mock.Anything, tokenHash).Return(true, nil)

		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		req := httptest.NewRequest("GET", "/api/v1/protetta", nil)
		req.Header.Set("Authorization", "Bearer "+tokenValido)
		rr := httptest.NewRecorder()

		handler := NewAuthMiddleware(sv)(nextHandler)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		sv.AssertExpectations(t)
	})

	t.Run("sessione non valida DB respinta 401", func(t *testing.T) {
		sv := new(mockSessionValidator)
		sv.On("SessionRepository_IsValid", mock.Anything, tokenHash).Return(false, nil)

		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		req := httptest.NewRequest("GET", "/api/v1/protetta", nil)
		req.Header.Set("Authorization", "Bearer "+tokenValido)
		rr := httptest.NewRecorder()

		handler := NewAuthMiddleware(sv)(nextHandler)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		sv.AssertExpectations(t)
	})
}