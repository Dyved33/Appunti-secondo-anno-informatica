package scanner

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/tommasogarofani/ingsw_25-26_group-11/client/api"

	tea "charm.land/bubbletea/v2"
)

func TestScanner_InputAndUpdate(t *testing.T) {
	client := &api.APIClient{}
	m := New(client)

	// Digitazione del codice "PRD-001"
	inputs := []string{"P", "R", "D", "-", "0", "0", "1"}
	var modelCurrent tea.Model = m
	for _, ch := range inputs {
		modelCurrent, _ = modelCurrent.Update(tea.KeyPressMsg{Text: ch})
	}
	mInput := modelCurrent.(Model)

	if mInput.input != "PRD-001" {
		t.Errorf("Input atteso 'PRD-001', ottenuto '%s'", mInput.input)
	}

	// Test ProdottoScansionatoMsg
	prodMock := api.Prodotto{
		ID: 1, Codice: "PRD-001", Nome: "Bulloni M8", Categoria: "Hardware",
		PrezzoAcquisto: 1.0, PrezzoVendita: 1.50, Attivo: true,
	}

	resScan, _ := mInput.Update(ProdottoScansionatoMsg{Prodotto: prodMock})
	mScan := resScan.(Model)

	if mScan.prodottoTrovato == nil {
		t.Fatal("prodottoTrovato non dovrebbe essere nil dopo la scansione")
	}
	if mScan.prodottoTrovato.Codice != "PRD-001" {
		t.Errorf("Codice atteso 'PRD-001', ottenuto '%s'", mScan.prodottoTrovato.Codice)
	}
	if len(mScan.storicoScansioni) != 1 {
		t.Errorf("Storico atteso 1 elemento, trovati %d", len(mScan.storicoScansioni))
	}

	view := mScan.View()
	if !strings.Contains(view.Content, "Bulloni M8") {
		t.Error("View() dovrebbe contenere il nome del prodotto trovato Bulloni M8")
	}
}

// TestScanner_RicercaBarcodeConFallback testa scansionaCmd: ricerca prima per barcode,
// poi fallback per codice prodotto, infine errore se nessuna ricerca ha successo.
func TestScanner_RicercaBarcodeConFallback(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/products/barcode/8001234567897":
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"id": 8, "codice": "PROD008", "barcode": "8001234567897",
				"nome": "Risma Carta A4", "descrizione": "Carta bianca",
				"categoria_id": 4, "prezzo_acquisto": 2.9,
				"prezzo_vendita": 5.2, "soglia_minima": 50,
			})
		case "/api/v1/products":
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode([]map[string]any{
				{
					"id": 1, "codice": "PROD001", "barcode": "",
					"nome": "Bulloni M8", "descrizione": "",
					"categoria_id": 2, "prezzo_acquisto": 0.5,
					"prezzo_vendita": 1.0, "soglia_minima": 10,
				},
			})
		case "/api/v1/categories":
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode([]map[string]any{
				{"id": 2, "nome": "Hardware", "descrizione": ""},
				{"id": 4, "nome": "Cancelleria", "descrizione": ""},
			})
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	m := New(api.New(srv.URL))

	// Ricerca per barcode riuscita
	msg := m.scansionaCmd("8001234567897")()
	scan, ok := msg.(ProdottoScansionatoMsg)
	if !ok {
		t.Fatalf("atteso ProdottoScansionatoMsg, ottenuto %T", msg)
	}
	if scan.Prodotto.Codice != "PROD008" {
		t.Errorf("codice atteso 'PROD008', ottenuto '%s'", scan.Prodotto.Codice)
	}
	if scan.Prodotto.Categoria != "Cancelleria" {
		t.Errorf("categoria attesa 'Cancelleria', ottenuta '%s'", scan.Prodotto.Categoria)
	}

	// Fallback per codice prodotto quando il barcode non esiste
	msg = m.scansionaCmd("PROD001")()
	scan, ok = msg.(ProdottoScansionatoMsg)
	if !ok {
		t.Fatalf("atteso ProdottoScansionatoMsg da fallback, ottenuto %T", msg)
	}
	if scan.Prodotto.Nome != "Bulloni M8" {
		t.Errorf("nome atteso 'Bulloni M8', ottenuto '%s'", scan.Prodotto.Nome)
	}
	if scan.Prodotto.Categoria != "Hardware" {
		t.Errorf("categoria attesa 'Hardware', ottenuta '%s'", scan.Prodotto.Categoria)
	}

	// Nessun risultato: errore
	msg = m.scansionaCmd("INESISTENTE")()
	if _, ok := msg.(ErroreMsg); !ok {
		t.Fatalf("atteso ErroreMsg, ottenuto %T", msg)
	}
}
