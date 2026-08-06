package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestNewApiClient verifica che la funzione di inizializzazione NewApiClient
// configuri correttamente l'indirizzo del server e l'http.Client interno.
func TestNewApiClient(t *testing.T) {
	url := "http://localhost:8080"
	client := New(url)

	// Controlla che l'URL del server sia stato memorizzato correttamente
	if client.indirizzoServer != url {
		t.Errorf("expected url %s, got %s", url, client.indirizzoServer)
	}

	// Assicura che l'http.Client sia stato istanziato e non sia nil
	if client.httpClient == nil {
		t.Error("expected httpClient to be initialized")
	}
}

// TestApiClient_Login_Success simula un caso di successo per il metodo Login.
// Viene configurato un server HTTP di test che risponde con un token JWT valido.
func TestApiClient_Login_Success(t *testing.T) {
	expectedToken := "test-token-123"

	// Avvia un server di test locale per intercettare le chiamate del client
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verifica che l'endpoint invocato sia corretto (/login)
		if r.URL.Path != "/login" {
			t.Errorf("expected path /login, got %s", r.URL.Path)
		}

		// Il client deve effettuare una richiesta di tipo POST
		if r.Method != http.MethodPost {
			t.Errorf("expected POST method, got %s", r.Method)
		}

		// Decodifica il body per verificare che username e password inviati siano corretti
		var reqBody map[string]string
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}
		if reqBody["username"] != "admin" || reqBody["password"] != "password" {
			t.Errorf("unexpected credentials: %+v", reqBody)
		}

		// Invia una risposta JSON di successo contenente il token di sessione
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"token": expectedToken})
	}))
	defer server.Close()

	// Inizializza il client collegandolo al server di test
	client := New(server.URL)
	err := client.Login("admin", "password")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Controlla che il token sia stato salvato nello stato dell'ApiClient
	if client.token != expectedToken {
		t.Errorf("expected token %s, got %s", expectedToken, client.token)
	}
}

// TestApiClient_Login_Failure verifica il comportamento dell'ApiClient
// in caso di errore restituito dal server (es. credenziali errate o 401 Unauthorized).
func TestApiClient_Login_Failure(t *testing.T) {
	// Avvia un server di test locale che simula un fallimento dell'autenticazione
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("invalid credentials"))
	}))
	defer server.Close()

	// Inizializza il client e prova ad accedere
	client := New(server.URL)
	err := client.Login("admin", "wrong")

	// Ci aspettiamo che l'errore non sia nil
	if err == nil {
		t.Error("expected error, got nil")
	}
}
