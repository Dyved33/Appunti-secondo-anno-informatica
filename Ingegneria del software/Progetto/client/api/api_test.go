package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// setupTestServer crea un server HTTP di mock e il relativo client API configurato.
func setupTestServer(handler http.HandlerFunc) (*httptest.Server, *APIClient) {
	server := httptest.NewServer(handler)
	client := New(server.URL)
	return server, client
}

// TestNewApiClient verifica che il costruttore New configuri correttamente l'URL e l'httpClient.
func TestNewApiClient(t *testing.T) {
	t.Parallel()
	url := "http://localhost:8080"
	client := New(url)

	if client.indirizzoServer != url {
		t.Errorf("atteso url %s, ottenuto %s", url, client.indirizzoServer)
	}
	if client.httpClient == nil {
		t.Error("atteso httpClient non nil")
	}
	if client.Token() != "" {
		t.Errorf("atteso token iniziale vuoto, ottenuto %s", client.Token())
	}
	if client.Role() != "" {
		t.Errorf("atteso ruolo iniziale vuoto, ottenuto %s", client.Role())
	}
}

// TestApiClient_AuthAndHeaders verifica il login con credenziali corrette, errate,
// la memorizzazione del token Bearer e il ruolo utente associato.
func TestApiClient_AuthAndHeaders(t *testing.T) {
	t.Parallel()

	t.Run("Login e invio Token Bearer", func(t *testing.T) {
		expectedToken := "jwt-secret-token"
		expectedRole := "admin"

		server, client := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/api/v1/login" {
				if r.Method != http.MethodPost {
					t.Errorf("metodo atteso POST, ottenuto %s", r.Method)
				}
				var req map[string]string
				_ = json.NewDecoder(r.Body).Decode(&req)
				if req["username"] != "admin" || req["password"] != "pass" {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode(map[string]string{
					"token": expectedToken,
					"role":  expectedRole,
				})
				return
			}

			// Endpoint protetto: verifica invio del token Bearer
			if r.URL.Path == "/api/v1/protected" {
				authHeader := r.Header.Get("Authorization")
				if authHeader != "Bearer "+expectedToken {
					t.Errorf("header Authorization inatteso: ottenuto %s", authHeader)
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"status":"ok"}`))
				return
			}

			http.NotFound(w, r)
		})
		defer server.Close()

		// 1. Verifica Login
		err := client.Login("admin", "pass")
		if err != nil {
			t.Fatalf("Login fallito: %v", err)
		}
		if client.Token() != expectedToken {
			t.Errorf("atteso token %s, ottenuto %s", expectedToken, client.Token())
		}
		if client.Role() != expectedRole {
			t.Errorf("atteso ruolo %s, ottenuto %s", expectedRole, client.Role())
		}

		// 2. Verifica che le chiamate successive usino il token salvato
		var out struct{ Status string }
		err = client.doJSON(http.MethodGet, "/api/v1/protected", nil, &out)
		if err != nil {
			t.Fatalf("richiesta protetta fallita: %v", err)
		}
	})

	t.Run("Login Fallito - Credenziali Errate", func(t *testing.T) {
		server, client := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte("credenziali non valide"))
		})
		defer server.Close()

		err := client.Login("utente", "errata")
		if err == nil {
			t.Fatal("atteso errore per credenziali errate, ottenuto nil")
		}
		if !strings.Contains(err.Error(), "credenziali non valide") {
			t.Errorf("messaggio di errore inatteso: %v", err)
		}
	})
}

// TestAPIClient_Prodotti verifica tutte le operazioni sui prodotti del catalogo:
// recupero del catalogo, ricerca per codice e barcode, creazione, eliminazione e aggiornamento prezzo.
func TestAPIClient_Prodotti(t *testing.T) {
	t.Parallel()

	t.Run("GetCatalogo e GetProductByCode", func(t *testing.T) {
		server, client := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/api/v1/products":
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode([]serverProductDTO{
					{ID: 1, Codice: "P01", Nome: "Laptop", CategoriaID: 10, PrezzoAcq: 600.0, PrezzoVend: 999.99},
				})
			case "/api/v1/categories":
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode([]serverCategoryDTO{
					{ID: 10, Nome: "Informatica"},
				})
			default:
				http.NotFound(w, r)
			}
		})
		defer server.Close()

		prodotti, categorie, err := client.GetCatalogo()
		if err != nil {
			t.Fatalf("GetCatalogo fallito: %v", err)
		}
		if len(prodotti) != 1 || prodotti[0].Categoria != "Informatica" {
			t.Errorf("mappatura catalogo errata: %+v", prodotti)
		}
		if len(categorie) != 1 || categorie[0] != "Informatica" {
			t.Errorf("categorie inattese: %+v", categorie)
		}

		// Test GetProductByCode trovato
		prod, err := client.GetProductByCode("P01")
		if err != nil {
			t.Fatalf("GetProductByCode fallito: %v", err)
		}
		if prod.Nome != "Laptop" || prod.PrezzoAcquisto != 600.0 || prod.PrezzoVendita != 999.99 {
			t.Errorf("prodotto inatteso: %+v", prod)
		}

		// Test GetProductByCode non trovato
		_, err = client.GetProductByCode("NON_ESISTE")
		if err == nil {
			t.Error("atteso errore per codice inesistente")
		}
	})

	t.Run("GetCatalogo Errore Categorie", func(t *testing.T) {
		server, client := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/api/v1/products" {
				_ = json.NewEncoder(w).Encode([]serverProductDTO{})
				return
			}
			if r.URL.Path == "/api/v1/categories" {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			http.NotFound(w, r)
		})
		defer server.Close()

		_, _, err := client.GetCatalogo()
		if err == nil {
			t.Fatal("atteso errore quando getCategorie fallisce")
		}
	})

	t.Run("GetProductByBarcode Successo ed Errore", func(t *testing.T) {
		server, client := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/api/v1/products/barcode/12345678" {
				_ = json.NewEncoder(w).Encode(serverProductDTO{
					ID: 5, Codice: "BAR01", Nome: "Mouse", CategoriaID: 1, PrezzoAcq: 10.0, PrezzoVend: 25.0,
				})
				return
			}
			if r.URL.Path == "/api/v1/categories" {
				_ = json.NewEncoder(w).Encode([]serverCategoryDTO{{ID: 1, Nome: "Accessori"}})
				return
			}
			http.NotFound(w, r)
		})
		defer server.Close()

		p, err := client.GetProductByBarcode("12345678")
		if err != nil {
			t.Fatalf("GetProductByBarcode fallito: %v", err)
		}
		if p.ID != 5 || p.Categoria != "Accessori" || p.PrezzoVendita != 25.0 {
			t.Errorf("dati prodotto errati: %+v", p)
		}

		// Errore 404 barcode inesistente
		_, err = client.GetProductByBarcode("99999999")
		if err == nil {
			t.Error("atteso errore per barcode inesistente")
		}
	})

	t.Run("GetProductByBarcode Errore Categorie", func(t *testing.T) {
		server, client := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/api/v1/products/barcode/12345678" {
				_ = json.NewEncoder(w).Encode(serverProductDTO{ID: 5, CategoriaID: 1})
				return
			}
			if r.URL.Path == "/api/v1/categories" {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			http.NotFound(w, r)
		})
		defer server.Close()

		_, err := client.GetProductByBarcode("12345678")
		if err == nil {
			t.Fatal("atteso errore quando getCategorie fallisce in GetProductByBarcode")
		}
	})

	t.Run("AggiornaPrezzo Successo ed Errore", func(t *testing.T) {
		server, client := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPut && r.URL.Path == "/api/v1/products/10/price" {
				var body map[string]float64
				_ = json.NewDecoder(r.Body).Decode(&body)
				if body["prezzo_vendita"] == 199.99 {
					w.WriteHeader(http.StatusOK)
					return
				}
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			http.NotFound(w, r)
		})
		defer server.Close()

		err := client.AggiornaPrezzo(10, 199.99)
		if err != nil {
			t.Fatalf("AggiornaPrezzo fallito: %v", err)
		}

		err = client.AggiornaPrezzo(999, 10.0)
		if err == nil {
			t.Fatal("atteso errore per prodotto non trovato")
		}
	})

	t.Run("CreaProdotto ed EliminaProdotto", func(t *testing.T) {
		server, client := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost && r.URL.Path == "/api/v1/products/" {
				var input ClientProductDTO
				_ = json.NewDecoder(r.Body).Decode(&input)
				_ = json.NewEncoder(w).Encode(serverProductDTO{
					ID: 99, Codice: input.Codice, Nome: input.Nome, CategoriaID: input.CategoriaID,
					PrezzoAcq: input.PrezzoAcq, PrezzoVend: input.PrezzoVend,
				})
				return
			}
			if r.Method == http.MethodDelete && r.URL.Path == "/api/v1/products/99" {
				w.WriteHeader(http.StatusOK)
				return
			}
			if r.URL.Path == "/api/v1/categories" {
				_ = json.NewEncoder(w).Encode([]serverCategoryDTO{{ID: 2, Nome: "Periferiche"}})
				return
			}
			http.NotFound(w, r)
		})
		defer server.Close()

		prodotto, err := client.CreaProdotto(ClientProductDTO{Codice: "NEW", Nome: "Tastiera", CategoriaID: 2, PrezzoVend: 49.99})
		if err != nil {
			t.Fatalf("CreaProdotto fallito: %v", err)
		}
		if prodotto.ID != 99 || prodotto.Categoria != "Periferiche" || prodotto.PrezzoVendita != 49.99 {
			t.Errorf("prodotto creato non valido: %+v", prodotto)
		}

		err = client.EliminaProdotto(99)
		if err != nil {
			t.Fatalf("EliminaProdotto fallito: %v", err)
		}
	})

	t.Run("CreaProdotto Errore Categorie", func(t *testing.T) {
		server, client := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost && r.URL.Path == "/api/v1/products/" {
				_ = json.NewEncoder(w).Encode(serverProductDTO{ID: 1, CategoriaID: 2})
				return
			}
			if r.URL.Path == "/api/v1/categories" {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			http.NotFound(w, r)
		})
		defer server.Close()

		_, err := client.CreaProdotto(ClientProductDTO{Codice: "P1", Nome: "Test", CategoriaID: 2})
		if err == nil {
			t.Fatal("atteso errore quando getCategorie fallisce in CreaProdotto")
		}
	})
}

// TestAPIClient_Categorie verifica la gestione delle categorie: recupero ID-nome, creazione ed eliminazione.
func TestAPIClient_Categorie(t *testing.T) {
	t.Parallel()

	server, client := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/v1/categories":
			_ = json.NewEncoder(w).Encode([]serverCategoryDTO{
				{ID: 1, Nome: "Ferramenta", Descrizione: "Viti e bulloni"},
				{ID: 2, Nome: "Elettronica", Descrizione: "Circuiti"},
			})
		case r.Method == http.MethodPost && r.URL.Path == "/api/v1/categories/":
			var payload map[string]string
			_ = json.NewDecoder(r.Body).Decode(&payload)
			if payload["nome"] == "NuovaCat" {
				w.WriteHeader(http.StatusCreated)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
		case r.Method == http.MethodDelete && r.URL.Path == "/api/v1/categories/1":
			w.WriteHeader(http.StatusOK)
		default:
			http.NotFound(w, r)
		}
	})
	defer server.Close()

	// GetCategorieIDs
	idMap, catList, err := client.GetCategorieIDs()
	if err != nil {
		t.Fatalf("GetCategorieIDs fallito: %v", err)
	}
	if len(catList) != 2 || idMap[1] != "Ferramenta" || idMap[2] != "Elettronica" {
		t.Errorf("mappatura categorie errata: %+v, %+v", idMap, catList)
	}

	// CreaCategoria
	err = client.CreaCategoria("NuovaCat", "Descrizione test")
	if err != nil {
		t.Fatalf("CreaCategoria fallito: %v", err)
	}

	// EliminaCategoria
	err = client.EliminaCategoria(1)
	if err != nil {
		t.Fatalf("EliminaCategoria fallito: %v", err)
	}
}

// TestAPIClient_GiacenzeReport verifica il recupero del report di inventario.
func TestAPIClient_GiacenzeReport(t *testing.T) {
	t.Parallel()

	t.Run("GetReportGiacenze Successo", func(t *testing.T) {
		server, client := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/api/v1/products/inventory-report" {
				_ = json.NewEncoder(w).Encode(InventoryReport{
					TotaleQuantita: 150,
					ValoreTotale:   3000.0,
					Articoli: []InventoryItem{
						{ProdottoID: 1, Codice: "PRD-1", Nome: "Bullone", Quantita: 100, SogliaMin: 20},
					},
					ArticoliSottoSoglia: []InventoryItem{},
				})
				return
			}
			http.NotFound(w, r)
		})
		defer server.Close()

		rep, err := client.GetReportGiacenze()
		if err != nil {
			t.Fatalf("GetReportGiacenze fallito: %v", err)
		}
		if rep.TotaleQuantita != 150 || rep.ValoreTotale != 3000.0 || len(rep.Articoli) != 1 {
			t.Errorf("report giacenze inatteso: %+v", rep)
		}
	})

	t.Run("GetReportGiacenze Errore", func(t *testing.T) {
		server, client := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		})
		defer server.Close()

		_, err := client.GetReportGiacenze()
		if err == nil {
			t.Fatal("atteso errore su report giacenze non disponibile")
		}
	})
}

// TestAPIClient_Utenti verifica la gestione degli utenti: ricerca, creazione,
// aggiornamento, blocco/sblocco ed eliminazione.
func TestAPIClient_Utenti(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		action     func(c *APIClient) error
		serverFunc http.HandlerFunc
	}{
		{
			name: "GetUsersByUsername Successo",
			action: func(c *APIClient) error {
				u, err := c.GetUsersByUsername("mario")
				if err != nil || u.Username != "mario" {
					t.Fatalf("GetUsersByUsername errato: %v, user: %+v", err, u)
				}
				return nil
			},
			serverFunc: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodGet && r.URL.Path == "/api/v1/users/mario" {
					_ = json.NewEncoder(w).Encode(serverUserDTO{ID: 1, Username: "mario", Ruolo: "user"})
				} else {
					http.NotFound(w, r)
				}
			},
		},
		{
			name: "CreateUser Successo",
			action: func(c *APIClient) error {
				u, err := c.CreateUser(UserCreateInput{Username: "luigi", Ruolo: "admin"})
				if err != nil || u.Username != "luigi" {
					t.Fatalf("CreateUser errato: %v, user: %+v", err, u)
				}
				return nil
			},
			serverFunc: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPost && r.URL.Path == "/api/v1/users" {
					_ = json.NewEncoder(w).Encode(serverUserDTO{ID: 2, Username: "luigi", Ruolo: "admin"})
				} else {
					http.NotFound(w, r)
				}
			},
		},
		{
			name: "UpdateUser e ToggleUserBlock Successo",
			action: func(c *APIClient) error {
				u, err := c.ToggleUserBlock("mario", true)
				if err != nil || !u.Bloccato {
					t.Fatalf("ToggleUserBlock errato: %v, user: %+v", err, u)
				}
				return nil
			},
			serverFunc: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPut && r.URL.Path == "/api/v1/users/mario" {
					var input UserUpdateInput
					_ = json.NewDecoder(r.Body).Decode(&input)
					_ = json.NewEncoder(w).Encode(serverUserDTO{
						ID:       1,
						Username: "mario",
						Ruolo:    "user",
						Bloccato: *input.Bloccato,
					})
				} else {
					http.NotFound(w, r)
				}
			},
		},
		{
			name: "DeleteUser Successo",
			action: func(c *APIClient) error {
				return c.DeleteUser("mario")
			},
			serverFunc: func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodDelete && r.URL.Path == "/api/v1/users/mario" {
					w.WriteHeader(http.StatusOK)
				} else {
					http.NotFound(w, r)
				}
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server, client := setupTestServer(tt.serverFunc)
			defer server.Close()

			if err := tt.action(client); err != nil {
				t.Errorf("operazione fallita inattesa: %v", err)
			}
		})
	}
}

// TestAPIClient_OrdiniEFatture verifica le chiamate relative a fatture, note di credito,
// ordini di entrata, ordini di uscita, resi, DDT e statistiche KPI.
func TestAPIClient_OrdiniEFatture(t *testing.T) {
	t.Parallel()

	t.Run("GetFatture e GetReportKPI", func(t *testing.T) {
		server, client := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/api/v1/orders/invoices":
				_ = json.NewEncoder(w).Encode([]serverFatturaDTO{
					{ID: 101, OrdineUscitaID: 1, Totale: 150.0, Data: time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)},
				})
			case "/api/v1/orders/credit-notes":
				_ = json.NewEncoder(w).Encode([]serverFatturaDTO{
					{ID: 102, OrdineUscitaID: 1, Totale: 50.0, Data: time.Date(2026, 2, 20, 0, 0, 0, 0, time.UTC)},
				})
			case "/api/v1/orders":
				_ = json.NewEncoder(w).Encode([]serverOrderDTO{
					{ID: 1, Cliente: "ACME Corp", Stato: "evaso", Righe: []serverOrderLineDTO{{ProdottoID: 5, Quantita: 2, PrezzoUnitario: 75.0}}},
					{ID: 2, Cliente: "Second Corp", Stato: "confermato", Righe: []serverOrderLineDTO{{ProdottoID: 999, Quantita: 1, PrezzoUnitario: 10.0}}},
				})
			case "/api/v1/products":
				_ = json.NewEncoder(w).Encode([]serverProductDTO{{ID: 5, Nome: "Prodotto X"}})
			case "/api/v1/categories":
				_ = json.NewEncoder(w).Encode([]serverCategoryDTO{})
			default:
				http.NotFound(w, r)
			}
		})
		defer server.Close()

		fatture, err := client.GetFatture()
		if err != nil {
			t.Fatalf("GetFatture fallito: %v", err)
		}
		if len(fatture) != 2 || fatture[0].Cliente != "ACME Corp" || !fatture[1].NotaCredito {
			t.Errorf("fatture inattese: %+v", fatture)
		}

		ordini, err := client.GetOrdini()
		if err != nil {
			t.Fatalf("GetOrdini fallito: %v", err)
		}
		if len(ordini) != 2 || ordini[0].Prodotto != "Prodotto X" || ordini[1].Prodotto != "Prodotto 999" {
			t.Errorf("ordini inattesi: %+v", ordini)
		}

		kpis, err := client.GetReportKPI()
		if err != nil {
			t.Fatalf("GetReportKPI fallito: %v", err)
		}

		expected := map[string]string{
			"Fatturato netto": "€ 100.00",
			"Fatture":         "1",
			"Ordini evasi":    "1",
			"Note credito":    "1",
		}
		for _, kpi := range kpis {
			if val, ok := expected[kpi.Etichetta]; !ok || val != kpi.Valore {
				t.Errorf("KPI %s inatteso: ottenuto %s, atteso %s", kpi.Etichetta, kpi.Valore, val)
			}
		}
	})

	t.Run("GetOrdiniEntrata e GetOrdiniUscita", func(t *testing.T) {
		server, client := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/api/v1/orders/entrata":
				_ = json.NewEncoder(w).Encode([]serverOrderDTO{
					{ID: 1, Numero: "E-01", Fornitore: "Fornitore A", Righe: []serverOrderLineDTO{
						{ID: 10, ProdottoID: 1, Quantita: 50, PrezzoUnitario: 10.0},
					}},
				})
			case "/api/v1/orders/uscita":
				ddtID := 77
				_ = json.NewEncoder(w).Encode([]serverOrderDTO{
					{ID: 2, Numero: "U-01", Cliente: "Cliente B", Stato: "confermato", DDTID: &ddtID, Righe: []serverOrderLineDTO{
						{ID: 20, ProdottoID: 1, Quantita: 10, PrezzoUnitario: 15.0},
					}},
				})
			case "/api/v1/products":
				_ = json.NewEncoder(w).Encode([]serverProductDTO{
					{ID: 1, Codice: "P1", Nome: "Prodotto 1", CategoriaID: 1},
				})
			case "/api/v1/categories":
				_ = json.NewEncoder(w).Encode([]serverCategoryDTO{{ID: 1, Nome: "Cat 1"}})
			default:
				http.NotFound(w, r)
			}
		})
		defer server.Close()

		entrate, err := client.GetOrdiniEntrata()
		if err != nil {
			t.Fatalf("GetOrdiniEntrata fallito: %v", err)
		}
		if len(entrate) != 1 || len(entrate[0].Righe) != 1 || entrate[0].Righe[0].Prodotto == nil || entrate[0].Righe[0].Prodotto.Nome != "Prodotto 1" {
			t.Errorf("ordini entrata inattesi: %+v", entrate)
		}

		uscite, err := client.GetOrdiniUscita()
		if err != nil {
			t.Fatalf("GetOrdiniUscita fallito: %v", err)
		}
		if len(uscite) != 1 || *uscite[0].DDTID != 77 || uscite[0].Righe[0].Prodotto == nil {
			t.Errorf("ordini uscita inattesi: %+v", uscite)
		}
	})

	t.Run("GetOrdiniEvasiLeggero", func(t *testing.T) {
		server, client := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/api/v1/orders/uscita" {
				_ = json.NewEncoder(w).Encode([]serverOrderDTO{
					{ID: 1, Numero: "U1", Cliente: "C1", Stato: "evaso"},
					{ID: 2, Numero: "U2", Cliente: "C2", Stato: "confermato"},
				})
				return
			}
			http.NotFound(w, r)
		})
		defer server.Close()

		evasi, err := client.GetOrdiniEvasiLeggero()
		if err != nil {
			t.Fatalf("GetOrdiniEvasiLeggero fallito: %v", err)
		}
		if len(evasi) != 1 || evasi[0].ID != 1 {
			t.Errorf("atteso solo ordine evaso 1, ottenuti: %+v", evasi)
		}
	})

	t.Run("Errori Di Chiamata In GetFatture, GetOrdini, GetReportKPI", func(t *testing.T) {
		server, client := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		})
		defer server.Close()

		if _, err := client.GetFatture(); err == nil {
			t.Error("atteso errore da GetFatture con server 500")
		}
		if _, err := client.GetOrdini(); err == nil {
			t.Error("atteso errore da GetOrdini con server 500")
		}
		if _, err := client.GetReportKPI(); err == nil {
			t.Error("atteso errore da GetReportKPI con server 500")
		}
		if _, err := client.GetOrdiniEntrata(); err == nil {
			t.Error("atteso errore da GetOrdiniEntrata con server 500")
		}
		if _, err := client.GetOrdiniUscita(); err == nil {
			t.Error("atteso errore da GetOrdiniUscita con server 500")
		}
		if _, err := client.GetResi(); err == nil {
			t.Error("atteso errore da GetResi con server 500")
		}
	})

	t.Run("CRUD Ordini, Resi, DDT e Fatture", func(t *testing.T) {
		server, client := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
			switch {
			case r.Method == http.MethodPost && r.URL.Path == "/api/v1/orders/crea-uscita":
				w.WriteHeader(http.StatusCreated)
			case r.Method == http.MethodDelete && r.URL.Path == "/api/v1/orders/delete-uscita/10":
				w.WriteHeader(http.StatusOK)
			case r.Method == http.MethodPost && r.URL.Path == "/api/v1/orders/crea-entrata":
				w.WriteHeader(http.StatusCreated)
			case r.Method == http.MethodDelete && r.URL.Path == "/api/v1/orders/delete-entrata/20":
				w.WriteHeader(http.StatusOK)
			case r.Method == http.MethodPut && r.URL.Path == "/api/v1/orders/update-ordine-uscita/10":
				w.WriteHeader(http.StatusOK)
			case r.Method == http.MethodGet && r.URL.Path == "/api/v1/orders/resi":
				_ = json.NewEncoder(w).Encode([]serverResoDTO{
					{ID: 1, OrdineUscitaID: 10, ProdottoID: 5, Quantita: 2, Motivo: "difettoso", Stato: "da_verificare"},
				})
			case r.Method == http.MethodPost && r.URL.Path == "/api/v1/orders/crea-reso":
				w.WriteHeader(http.StatusCreated)
			case r.Method == http.MethodPut && r.URL.Path == "/api/v1/orders/update-resi":
				w.WriteHeader(http.StatusOK)
			case r.Method == http.MethodDelete && r.URL.Path == "/api/v1/orders/delete-reso/1":
				w.WriteHeader(http.StatusOK)
			case r.Method == http.MethodPost && r.URL.Path == "/api/v1/orders/10/create-ddt":
				_ = json.NewEncoder(w).Encode(serverDDTDTO{
					ID:             100,
					Numero:         "DDT-100",
					OrdineUscitaID: 10,
					Mittente:       "Azienda",
					Destinatario:   "Cliente",
					PathPDF:        "/pdf/ddt100.pdf",
				})
			case r.Method == http.MethodPost && r.URL.Path == "/api/v1/orders/create-fattura/10":
				w.WriteHeader(http.StatusCreated)
			case r.Method == http.MethodDelete && r.URL.Path == "/api/v1/orders/delete-fattura/101":
				w.WriteHeader(http.StatusOK)
			case r.Method == http.MethodPost && r.URL.Path == "/api/v1/orders/create-credit-notes/10":
				w.WriteHeader(http.StatusCreated)
			default:
				http.NotFound(w, r)
			}
		})
		defer server.Close()

		// CreaOrdineUscita e DeleteOrdineUscita
		if err := client.CreaOrdineUscita(ClientOrderDTO{Numero: "U1"}); err != nil {
			t.Errorf("CreaOrdineUscita fallito: %v", err)
		}
		if err := client.DeleteOrdineUscita(10); err != nil {
			t.Errorf("DeleteOrdineUscita fallito: %v", err)
		}

		// CreaOrdineEntrata e DeleteOrdineEntrata
		if err := client.CreaOrdineEntrata(ClientOrderDTO{Numero: "E1"}); err != nil {
			t.Errorf("CreaOrdineEntrata fallito: %v", err)
		}
		if err := client.DeleteOrdineEntrata(20); err != nil {
			t.Errorf("DeleteOrdineEntrata fallito: %v", err)
		}

		// AggiornaOrdineUscita
		if err := client.AggiornaOrdineUscita(10, ClientOrderDTO{Numero: "U1-mod"}); err != nil {
			t.Errorf("AggiornaOrdineUscita fallito: %v", err)
		}

		// Resi (GetResi, CreaReso, AggiornaStatoReso, DeleteReso)
		resi, err := client.GetResi()
		if err != nil || len(resi) != 1 || resi[0].Stato != StatoInAttesa {
			t.Errorf("GetResi fallito o inatteso: %v, resi: %+v", err, resi)
		}
		if err := client.CreaReso(ClientResoDTO{OrdineUscitaID: 10, ProdottoID: 5, Quantita: 2}); err != nil {
			t.Errorf("CreaReso fallito: %v", err)
		}
		if err := client.AggiornaStatoReso(1, "accettato"); err != nil {
			t.Errorf("AggiornaStatoReso fallito: %v", err)
		}
		if err := client.DeleteReso(1); err != nil {
			t.Errorf("DeleteReso fallito: %v", err)
		}

		// DDT
		ddt, err := client.GetDDT(10, ClientDDTDTO{Mittente: "Azienda", Destinatario: "Cliente"})
		if err != nil || ddt.ID != 100 || ddt.Numero != "DDT-100" {
			t.Errorf("GetDDT fallito o inatteso: %v, ddt: %+v", err, ddt)
		}

		// Fatture e note credito
		if err := client.CreaFattura(10); err != nil {
			t.Errorf("CreaFattura fallito: %v", err)
		}
		if err := client.EliminaFattura(101); err != nil {
			t.Errorf("EliminaFattura fallito: %v", err)
		}
		if err := client.CreaNotaCredito(10); err != nil {
			t.Errorf("CreaNotaCredito fallito: %v", err)
		}
	})
}

// TestAPIClient_GestioneErrori verifica la corretta gestione dei codici di errore HTTP,
// dei messaggi di errore restituiti dal server, dei payload JSON non validi e dei problemi di rete.
func TestAPIClient_GestioneErrori(t *testing.T) {
	t.Parallel()

	t.Run("Errore HTTP con messaggio personalizzato", func(t *testing.T) {
		server, client := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("Parametro ID non valido"))
		})
		defer server.Close()

		err := client.EliminaProdotto(0)
		if err == nil || err.Error() != "Parametro ID non valido" {
			t.Errorf("errore atteso 'Parametro ID non valido', ottenuto: %v", err)
		}
	})

	t.Run("Errore HTTP con body vuoto", func(t *testing.T) {
		server, client := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		})
		defer server.Close()

		err := client.EliminaCategoria(1)
		if err == nil || err.Error() != "errore server (status 500)" {
			t.Errorf("errore di fallback inatteso: %v", err)
		}
	})

	t.Run("Risposta JSON Malformata dal Server", func(t *testing.T) {
		server, client := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{json-non-valido}`))
		})
		defer server.Close()

		_, err := client.GetUsersByUsername("mario")
		if err == nil || !strings.Contains(err.Error(), "errore di parsing risposta") {
			t.Errorf("atteso errore di parsing JSON, ottenuto: %v", err)
		}
	})

	t.Run("Server Irraggiungibile o Errore di Rete", func(t *testing.T) {
		client := New("http://127.0.0.1:59999") // Porta non in ascolto

		err := client.Login("admin", "password")
		if err == nil || !strings.Contains(err.Error(), "errore di rete") {
			t.Errorf("atteso errore di rete, ottenuto: %v", err)
		}
	})
}
