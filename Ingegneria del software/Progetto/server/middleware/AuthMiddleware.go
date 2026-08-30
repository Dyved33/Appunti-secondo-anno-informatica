// Package middleware contiene middleware HTTP usati dal server.
// In particolare AuthMiddleware protegge le rotte verificate l'autenticazione
// tramite token JWT e la validità della sessione lato server.
package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/tommasogarofani/ingsw_25-26_group-11/server/auth"
)

type contextKey string

// UserIDKey è la chiave usata per salvare l'ID dell'utente autenticato
// nel contesto della richiesta.
const UserIDKey contextKey = "userID"

// SessionValidator verifica che una sessione (token hash) sia ancora valida
// nel repository delle sessioni persistente.
type SessionValidator interface {
	SessionRepository_IsValid(ctx context.Context, tokenHash string) (bool, error)
}

// NewAuthMiddleware costruisce un middleware di autenticazione. Se viene
// fornito un SessionValidator non nil, oltre alla validità del token JWT
// viene verificata anche la sessione lato server; in caso di token o sessione
// non validi restituisce 401 Unauthorized.
func NewAuthMiddleware(sv SessionValidator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization") // Recupera il token JWT dall'intestazione Authorization della richiesta HTTP
			if header == "" {
				http.Error(w, "token assente", http.StatusUnauthorized) // = error 401 Unauthorized
				return
			}

			tokenString := strings.TrimPrefix(header, "Bearer ") // Rimuove il prefisso "Bearer " dal token JWT per ottenere il token effettivo
			claims, err := auth.ValidateToken(tokenString)
			if err != nil {
				http.Error(w, "token invalido", http.StatusUnauthorized) // = error 401 Unauthorized
				return
			}

			if sv != nil {
				tokenHash := auth.HashToken(tokenString)
				valid, err := sv.SessionRepository_IsValid(r.Context(), tokenHash)
				if err != nil || !valid {
					http.Error(w, "sessione invalida o scaduta", http.StatusUnauthorized)
					return
				}
			}

			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// AuthMiddleware applica l'autenticazione basata solo sul token JWT,
// senza verifica della sessione lato server.
func AuthMiddleware(next http.Handler) http.Handler { // L'http.Handler rappresenta il prossimo handler da chiamare se l'utente è autenticato
	return NewAuthMiddleware(nil)(next)
}
