// Package fatture raccoglie i test d'unità della sotto-pagina Fatture.
package fatture

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/tommasogarofani/ingsw_25-26_group-11/client/api"

	tea "charm.land/bubbletea/v2"
)

// TestNew_StatoIniziale verifica che il modello parta in stato di caricamento
// con il cursore a zero.
func TestNew_StatoIniziale(t *testing.T) {
	m := New(nil)
	if !m.loading {
		t.Error("atteso loading=true all'avvio")
	}
	if m.cursore != 0 {
		t.Errorf("atteso cursore iniziale a 0, ottenuto %d", m.cursore)
	}
}

// TestUpdate_FattureLoadedMsg verifica che il caricamento popoli il registro fatture.
func TestUpdate_FattureLoadedMsg(t *testing.T) {
	m := New(nil)
	res, _ := m.Update(FattureLoadedMsg{
		Fatture: []api.Fattura{
			{Numero: "001", Anno: 2026, Cliente: "Mario Rossi", Importo: 120.50},
		},
		NoteCredito: []api.Fattura{
			{Numero: "002", Anno: 2026, Cliente: "Luigi Bianchi", Importo: 80, NotaCredito: true},
		},
	})
	m = res.(Model)

	if m.loading {
		t.Error("atteso loading=false dopo il caricamento")
	}
	if len(m.fatture) != 1 {
		t.Errorf("attesa 1 fattura, ottenute %d", len(m.fatture))
	}
	if len(m.noteCredito) != 1 {
		t.Errorf("attesa 1 nota di credito, ottenute %d", len(m.noteCredito))
	}
}

// TestUpdate_CursoreResettatoAlRicaricamento verifica che un nuovo registro non
// lasci il cursore puntato su una riga che ora non esiste più.
func TestUpdate_CursoreResettatoAlRicaricamento(t *testing.T) {
	m := caricaFattureDiProva(t, 3)

	res, _ := m.Update(tea.KeyPressMsg{Text: "down"})
	m = res.(Model)
	if m.cursore != 1 {
		t.Fatalf("precondizione fallita: cursore atteso a 1, ottenuto %d", m.cursore)
	}

	res, _ = m.Update(FattureLoadedMsg{Fatture: []api.Fattura{{Numero: "001", Anno: 2026}}})
	m = res.(Model)
	if m.cursore != 0 {
		t.Errorf("atteso reset del cursore a 0 dopo un nuovo caricamento, ottenuto %d", m.cursore)
	}
}

// TestUpdate_NavigazioneCursore verifica lo spostamento del cursore con up/down.
func TestUpdate_NavigazioneCursore(t *testing.T) {
	m := caricaFattureDiProva(t, 3)

	res, _ := m.Update(tea.KeyPressMsg{Text: "down"})
	m = res.(Model)
	if m.cursore != 1 {
		t.Errorf("atteso cursore a 1 dopo 'down', ottenuto %d", m.cursore)
	}

	res, _ = m.Update(tea.KeyPressMsg{Text: "up"})
	m = res.(Model)
	if m.cursore != 0 {
		t.Errorf("atteso cursore a 0 dopo 'up', ottenuto %d", m.cursore)
	}
}

// TestUpdate_CursoreNonEsceDaiLimiti verifica che il cursore non vada mai
// sotto 0 né sopra l'ultimo indice disponibile.
func TestUpdate_CursoreNonEsceDaiLimiti(t *testing.T) {
	m := caricaFattureDiProva(t, 1)

	res, _ := m.Update(tea.KeyPressMsg{Text: "up"})
	m = res.(Model)
	if m.cursore != 0 {
		t.Errorf("il cursore non deve andare sotto 0, ottenuto %d", m.cursore)
	}

	res, _ = m.Update(tea.KeyPressMsg{Text: "down"})
	m = res.(Model)
	if m.cursore != 0 {
		t.Errorf("con una sola fattura il cursore deve restare a 0, ottenuto %d", m.cursore)
	}
}

// TestUpdate_ErroreMsg verifica che un errore async venga registrato correttamente.
func TestUpdate_ErroreMsg(t *testing.T) {
	m := New(nil)
	res, _ := m.Update(ErroreMsg{Err: errors.New("timeout")})
	m = res.(Model)

	if m.err == nil || m.loading {
		t.Error("atteso err impostato e loading=false")
	}
}

// TestUpdate_RicaricaConR verifica che "r" rimetta la pagina in caricamento e
// restituisca il Cmd di ricaricamento del registro.
func TestUpdate_RicaricaConR(t *testing.T) {
	m := caricaFattureDiProva(t, 1)

	res, cmd := m.Update(tea.KeyPressMsg{Text: "r"})
	m = res.(Model)

	if !m.loading {
		t.Error("atteso loading=true dopo 'r'")
	}
	if cmd == nil {
		t.Error("atteso un Cmd di ricaricamento dopo 'r'")
	}
}

// TestFatturatoNetto_NoteCreditoSottratte verifica che le note di credito vengano
// sottratte dal fatturato, con lo stesso segno usato da api.GetReportKPI.
func TestFatturatoNetto_NoteCreditoSottratte(t *testing.T) {
	m := New(nil)
	res, _ := m.Update(FattureLoadedMsg{
		Fatture:     []api.Fattura{{Numero: "001", Anno: 2026, Importo: 120}},
		NoteCredito: []api.Fattura{{Numero: "002", Anno: 2026, Importo: 20, NotaCredito: true}},
	})
	m = res.(Model)

	netto := m.fatturatoNetto()
	if netto != 100 {
		t.Errorf("atteso fatturato netto 100 (120 - 20), ottenuto %v", netto)
	}
}

// TestUpdate_CambioTab verifica che il tab alterni tra Fatture e Note di Credito.
func TestUpdate_CambioTab(t *testing.T) {
	m := caricaFattureDiProvaConNoteCredito(t)

	if m.sezione != SEZ_FATTURE {
		t.Error("la sezione iniziale deve essere SEZ_FATTURE")
	}

	res, _ := m.Update(tea.KeyPressMsg{Text: "tab"})
	m = res.(Model)
	if m.sezione != SEZ_NOTE_CREDITO {
		t.Error("dopo un tab la sezione deve essere SEZ_NOTE_CREDITO")
	}

	res, _ = m.Update(tea.KeyPressMsg{Text: "tab"})
	m = res.(Model)
	if m.sezione != SEZ_FATTURE {
		t.Error("dopo due tab la sezione deve tornare a SEZ_FATTURE")
	}
}

// TestView_NotaCreditoEImporto verifica che la riga mostri sia l'importo sia il
// tipo di documento nella tab Note di Credito.
func TestView_NotaCreditoEImporto(t *testing.T) {
	m := New(nil)
	res, _ := m.Update(FattureLoadedMsg{
		NoteCredito: []api.Fattura{
			{Numero: "001", Anno: 2026, Cliente: "Mario Rossi", Importo: 120.50, NotaCredito: true},
		},
	})
	m = res.(Model)
	m.sezione = SEZ_NOTE_CREDITO

	v := m.View()
	if !strings.Contains(v.Content, "nota di credito") {
		t.Errorf("attesa l'etichetta 'nota di credito' nella view, ottenuta %q", v.Content)
	}
	if !strings.Contains(v.Content, "€ 120.50") {
		t.Errorf("atteso l'importo del documento nella view, ottenuta %q", v.Content)
	}
}

// TestView_RegistroVuoto verifica che un registro senza documenti non produca
// una view vuota né un panic.
func TestView_RegistroVuoto(t *testing.T) {
	m := New(nil)
	res, _ := m.Update(FattureLoadedMsg{})
	m = res.(Model)

	if v := m.View(); !strings.Contains(v.Content, "Nessuna fattura") {
		t.Errorf("atteso il messaggio di registro vuoto, ottenuta %q", v.Content)
	}
}

// TestLoadFattureCmd_Successo esegue il Cmd di Init() contro un backend di prova
// e verifica che fatture e note di credito vengano separate correttamente.
func TestLoadFattureCmd_Successo(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.URL.Path {
		case "/api/v1/orders/invoices":
			json.NewEncoder(w).Encode([]map[string]any{{
				"id":               1,
				"numero":           "FT-001",
				"ordine_uscita_id": 1,
				"totale":           120.50,
				"data":             "2026-03-01T10:00:00Z",
			}})
		case "/api/v1/orders/credit-notes":
			json.NewEncoder(w).Encode([]map[string]any{{
				"id":               2,
				"numero":           "NC-001",
				"ordine_uscita_id": 1,
				"totale":           20.00,
				"data":             "2026-04-01T10:00:00Z",
			}})
		case "/api/v1/orders":
			json.NewEncoder(w).Encode([]map[string]any{{
				"id":      1,
				"numero":  "ORD-001",
				"cliente": "Mario Rossi",
				"stato":   "evaso",
			}})
		case "/api/v1/products":
			json.NewEncoder(w).Encode([]map[string]any{})
		case "/api/v1/categories":
			json.NewEncoder(w).Encode([]map[string]any{})
		default:
			t.Errorf("endpoint non previsto dal test: %s", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer srv.Close()

	m := New(api.New(srv.URL))

	msg := m.Init()()
	loaded, ok := msg.(FattureLoadedMsg)
	if !ok {
		t.Fatalf("atteso FattureLoadedMsg, ottenuto %T (%v)", msg, msg)
	}
	if len(loaded.Fatture) != 1 {
		t.Fatalf("attesa 1 fattura, ottenute %d", len(loaded.Fatture))
	}
	if len(loaded.NoteCredito) != 1 {
		t.Fatalf("attesa 1 nota di credito, ottenute %d", len(loaded.NoteCredito))
	}

	fattura := loaded.Fatture[0]
	if fattura.NotaCredito {
		t.Error("il primo elemento doveva essere una fattura")
	}
	if fattura.Anno != 2026 {
		t.Errorf("atteso anno 2026 ricavato dalla data del documento, ottenuto %d", fattura.Anno)
	}
	if fattura.Cliente != "Mario Rossi" {
		t.Errorf("atteso il cliente risolto dall'ordine di origine, ottenuto %q", fattura.Cliente)
	}

	if !loaded.NoteCredito[0].NotaCredito {
		t.Error("il secondo elemento doveva essere una nota di credito")
	}
}

// TestLoadFattureCmd_ErroreServer verifica che un errore HTTP diventi un
// ErroreMsg e chiuda lo stato di caricamento.
func TestLoadFattureCmd_ErroreServer(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	m := New(api.New(srv.URL))

	msg := m.Init()()
	if _, ok := msg.(ErroreMsg); !ok {
		t.Fatalf("atteso ErroreMsg su risposta 500, ottenuto %T", msg)
	}

	res, _ := m.Update(msg)
	m = res.(Model)
	if m.loading || m.err == nil {
		t.Error("atteso loading=false ed errore memorizzato dopo il fallimento del caricamento")
	}
}

// caricaFattureDiProva restituisce un Model già popolato con n fatture generiche,
// utile per i test che devono partire da uno stato non vuoto.
func caricaFattureDiProva(t *testing.T, n int) Model {
	t.Helper()
	fatture := make([]api.Fattura, n)
	for i := range fatture {
		fatture[i] = api.Fattura{
			Numero:  "FT-" + strconv.Itoa(i+1),
			Anno:    2026,
			Cliente: "Cliente di prova",
			Importo: 100,
		}
	}

	m := New(nil)
	res, _ := m.Update(FattureLoadedMsg{Fatture: fatture})
	return res.(Model)
}

// caricaFattureDiProvaConNoteCredito restituisce un Model con sia fatture che note di credito.
func caricaFattureDiProvaConNoteCredito(t *testing.T) Model {
	t.Helper()
	m := New(nil)
	res, _ := m.Update(FattureLoadedMsg{
		Fatture: []api.Fattura{
			{Numero: "FT-001", Anno: 2026, Cliente: "Cliente A", Importo: 100},
			{Numero: "FT-002", Anno: 2026, Cliente: "Cliente B", Importo: 200},
		},
		NoteCredito: []api.Fattura{
			{Numero: "NC-001", Anno: 2026, Cliente: "Cliente A", Importo: 50, NotaCredito: true},
		},
	})
	return res.(Model)
}
