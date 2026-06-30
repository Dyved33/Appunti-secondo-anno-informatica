package middleware

// Il middleware serve a controllare se l'utente è autenticato prima di accedere alle rotte protette. 
// Se l'utente non è autenticato, viene restituito un errore 401 Unauthorized. 
// Se l'utente è autenticato, il middleware aggiunge l'ID dell'utente al contesto della richiesta e passa la richiesta al prossimo handler.

import (
    "context"
    "net/http"
    "strings"

    "github.com/tommasogarofani/ingsw_25-26_group-11/src/backend/auth"
)

type contextKey string

const UserIDKey contextKey = "userID"

func AuthMiddleware(next http.Handler) http.Handler {// L'http.Handler rappresenta il prossimo handler da chiamare se l'utente è autenticato
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        header := r.Header.Get("Authorization")// Recupera il token JWT dall'intestazione Authorization della richiesta HTTP
        if header == "" {
            http.Error(w, "token assente", http.StatusUnauthorized)// = error 401 Unauthorized
            return
        }

        tokenString := strings.TrimPrefix(header, "Bearer ")// Rimuove il prefisso "Bearer " dal token JWT per ottenere il token effettivo
        claims, err := auth.ValidateToken(tokenString)
        if err != nil {
            http.Error(w, "token invalido", http.StatusUnauthorized)// = error 401 Unauthorized
            return
        }

        ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}