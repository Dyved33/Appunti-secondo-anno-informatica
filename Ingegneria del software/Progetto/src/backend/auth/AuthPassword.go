package auth

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {// Trasforma la password in chiaro in un hash sicuro utilizzando bcrypt
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(bytes), err
}

func CheckPassword(hash, password string) error {// Verifica se la password in chiaro corrisponde all'hash memorizzato, restituendo un errore se non corrispondono
    return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}