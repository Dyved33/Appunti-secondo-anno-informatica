package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
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
	expectedRole := "admin"

	// Avvia un server di test locale per intercettare le chiamate del client
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verifica che l'endpoint invocato sia corretto (/api/v1/login)
		if r.URL.Path != "/api/v1/login" {
			t.Errorf("expected path /api/v1/login, got %s", r.URL.Path)
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
		json.NewEncoder(w).Encode(map[string]string{"token": expectedToken, "role": expectedRole})
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
	if client.role != expectedRole {
		t.Errorf("expected role %s, got %s", expectedRole, client.role)
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

func TestAPIClient_CreaOrdini(t *testing.T) {
	tests := []struct {
		name        string
		crea        func(*APIClient) error
		wantPath    string
		wantPayload map[string]any
	}{
		{
			name:     "uscita",
			wantPath: "/api/v1/orders/crea-uscita",
			wantPayload: map[string]any{
				"numero": "OU-001", "cliente": "Cliente Srl", "note": "urgente",
				"righe": []any{map[string]any{"prodotto_id": float64(7), "quantita": float64(2), "prezzo_unitario": 1.5}},
			},
			crea: func(c *APIClient) error {
				return c.CreaOrdineUscita(ClientOrderDTO{
					Numero: "OU-001", Cliente: "Cliente Srl", Note: "urgente",
					Righe: []ClientOrderLineDTO{{ProdottoID: 7, Quantita: 2, PrezzoUnitario: 1.5}},
				})
			},
		},
		{
			name:     "entrata",
			wantPath: "/api/v1/orders/crea-entrata",
			wantPayload: map[string]any{
				"numero": "OE-001", "fornitore": "Fornitore Srl", "note": "consegna parziale",
				"righe": []any{map[string]any{"prodotto_id": float64(4), "quantita": float64(5), "prezzo_unitario": 0.8}},
			},
			crea: func(c *APIClient) error {
				return c.CreaOrdineEntrata(ClientOrderDTO{
					Numero: "OE-001", Fornitore: "Fornitore Srl", Note: "consegna parziale",
					Righe: []ClientOrderLineDTO{{ProdottoID: 4, Quantita: 5, PrezzoUnitario: 0.8}},
				})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost || r.URL.Path != tt.wantPath {
					t.Errorf("attesa POST %s, ottenuta %s %s", tt.wantPath, r.Method, r.URL.Path)
				}
				if got := r.Header.Get("Authorization"); got != "Bearer token-di-test" {
					t.Errorf("header Authorization inatteso: %q", got)
				}

				var payload map[string]any
				if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
					t.Fatalf("body JSON non valido: %v", err)
				}
				if !reflect.DeepEqual(payload, tt.wantPayload) {
					t.Errorf("payload inatteso:\n ottenuto: %#v\n atteso: %#v", payload, tt.wantPayload)
				}
				w.WriteHeader(http.StatusCreated)
			}))
			defer server.Close()

			client := New(server.URL)
			client.token = "token-di-test"
			if err := tt.crea(client); err != nil {
				t.Fatalf("creazione ordine fallita: %v", err)
			}
		})
	}
}
