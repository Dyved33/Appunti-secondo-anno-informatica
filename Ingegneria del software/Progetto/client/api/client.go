package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// getter del token di ApiClient
func (c *APIClient) Token() string {
	return c.token
}

// getter del ruolo autenticato
func (c *APIClient) Role() string {
	return c.role
}

// doJSON è una funzione helper per inviare richieste HTTP con payload JSON e decodificare la risposta JSON.
func (c *APIClient) doJSON(method, path string, body any, out any) error {
	var payload *bytes.Reader
	if body != nil {
		raw, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("errore di serializzazione: %w", err)
		}
		payload = bytes.NewReader(raw)
	} else {
		payload = bytes.NewReader(nil)
	}

	baseURL := strings.TrimRight(c.indirizzoServer, "/")
	req, err := http.NewRequest(method, baseURL+path, payload)
	if err != nil {
		return fmt.Errorf("errore di richiesta: %w", err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("errore di rete: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		bodyBytes, _ := io.ReadAll(resp.Body)
		msg := strings.TrimSpace(string(bodyBytes))
		if msg != "" {
			return errors.New(msg)
		}
		return fmt.Errorf("errore server (status %d)", resp.StatusCode)
	}

	if out == nil {
		return nil
	}

	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return fmt.Errorf("errore di parsing risposta: %w", err)
	}

	return nil
}

// Login effettua il login al server API e memorizza il token di autenticazione.
func (c *APIClient) Login(username, password string) error {
	var responseData struct {
		Token string `json:"token"`
		Role  string `json:"role"`
	}

	if err := c.doJSON(http.MethodPost, "/api/v1/login", map[string]string{
		"username": username,
		"password": password,
	}, &responseData); err != nil {
		return err
	}

	c.token = responseData.Token
	c.role = strings.TrimSpace(responseData.Role)

	return nil
}
