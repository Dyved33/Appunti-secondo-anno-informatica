package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// creo la struct del client api
type ApiClient struct {
	indirizzoServer string
	token           string
	httpClient      *http.Client
}

// funzione per inizializzare il client
// prende una stringa che contien l'url del server
// restituisce il puntatore al client
func New(url string) *ApiClient {
	return &ApiClient{
		indirizzoServer: url,
		httpClient:      &http.Client{},
	}
}

// getter del token di ApiClient
func (c *ApiClient) Token() string {
	return c.token
}

// funzione per il login
// riceve Utente e Password
// restituisce nil
func (c *ApiClient) Login(username, password string) error {

	payload, err := json.Marshal(map[string]string{
		"username": username,
		"password": password,
	})
	if err != nil {
		return fmt.Errorf("errore di serializzazione: %w", err)
	}
	resp, err := c.httpClient.Post(c.indirizzoServer+"/login", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("errore di rete: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("credenziali non valide o errore server (status %d)", resp.StatusCode)
	}

	var responseData struct {
		Token string `json:"token"` //cerca token nel json e lo collega alla variabile Token
	}

	if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
		return fmt.Errorf("Token non Valido: %s", err)
	}

	c.token = responseData.Token
	return nil
}
