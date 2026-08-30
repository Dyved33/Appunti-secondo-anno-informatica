package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// jwtSecret è la chiave segreta usata per firmare e verificare i token JWT.
// In produzione deve essere sostituita con un valore sicuro e riservato.
var jwtSecret = []byte("CAMBIA_QUESTA_CHIAVE")

// Claims rappresenta i dati incorporati in un token JWT: l'identità dell'utente
// (ID, username, ruolo) più i claim registrati standard (ad es. la scadenza).
type Claims struct {
	// UserID è l'identificativo numerico dell'utente.
	UserID int `json:"user_id"`
	// Username è il nome dell'utente.
	Username string `json:"username"`
	// Role è il ruolo dell'utente (es. admin, cliente).
	Role string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken crea un nuovo token JWT firmato con HS256 e valido 24 ore,
// contenente ID, username e ruolo dell'utente.
func GenerateToken(userID int, username string, role string) (string, error) {
	claims := Claims{ // Struct con i dati contenuti nel token
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret) // Firma digitalmente il token con la chiave segreta
}

// ValidateToken verifica che la firma del token sia autentica e che il token
// non sia scaduto. Restituisce i Claims estratti se il token è valido,
// altrimenti un errore.
func ValidateToken(tokenString string) (*Claims, error) { // Verifica che la firma sia autentica e che il token non sia scaduto
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("token invalido")
	}
	return claims, nil
}

//N.B. Chiunque può leggere il contenuto del token, ma solo chi ha la chiave segreta può modificarlo

// HashToken calcola l'hash SHA-256 del token e lo restituisce come stringa
// esadecimale. Viene usato per registrare la sessione senza salvare il token in chiaro.
func HashToken(tokenString string) string {
	hash := sha256.Sum256([]byte(tokenString))
	return hex.EncodeToString(hash[:])
}
