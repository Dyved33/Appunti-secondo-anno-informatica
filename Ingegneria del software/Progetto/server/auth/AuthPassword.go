package auth

import "golang.org/x/crypto/bcrypt"

// HashPassword trasforma una password in chiaro in un hash sicuro utilizzando
// bcrypt con il costo predefinito. Restituisce l'hash come stringa oppure un
// errore se la generazione fallisce.
func HashPassword(password string) (string, error) { // Trasforma la password in chiaro in un hash sicuro utilizzando bcrypt
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword verifica se la password in chiaro corrisponde all'hash bcrypt
// memorizzato. Restituisce nil se corrispondono, altrimenti un errore.
func CheckPassword(hash, password string) error { // Verifica se la password in chiaro corrisponde all'hash memorizzato, restituendo un errore se non corrispondono
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
