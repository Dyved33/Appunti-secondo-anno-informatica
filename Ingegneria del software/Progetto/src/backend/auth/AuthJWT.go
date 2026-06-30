package auth

import (
    "errors"
    "time"

    "github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("CAMBIA_QUESTA_CHIAVE")

type Claims struct {
    UserID   int    `json:"user_id"`
    Username string `json:"username"`
    jwt.RegisteredClaims
}

func GenerateToken(userID int, username string) (string, error) {
    claims := Claims{// Struct con i dati contenuti nel token
        UserID:   userID,
        Username: username,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
        },
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtSecret)// Firma digitalmente il token con la chiave segreta
}

func ValidateToken(tokenString string) (*Claims, error) {// Verifica che la firma sia autentica e che il token non sia scaduto
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