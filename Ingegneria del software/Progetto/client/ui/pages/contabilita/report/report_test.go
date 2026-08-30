// Package report raccoglie i test d'unità della sotto-pagina Report della Contabilità.
// I test garantiscono la copertura completa del ciclo di vita del modello Bubble Tea,
// la gestione dei comandi asincroni di caricamento dati tramite API client, il comportamento
// della vista nei vari stati (loading, errore, lista vuota, griglia KPI) e la risposta all'input utente.
package report

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/tommasogarofani/ingsw_25-26_group-11/client/api"

	tea "charm.land/bubbletea/v2"
)

// errMock è una struttura ausiliaria per simulare errori di backend o di rete nei test.
type errMock struct{}

func (errMock) Error() string { return "errore di connessione al server" }

// TestNew_StatoIniziale verifica che il costruttore New inizializzi correttamente
// il modello nello stato di caricamento attivo (`loading = true`) e senza errori.
func TestNew_StatoIniziale(t *testing.T) {
	m := New(nil)

	if !m.loading {
		t.Error("atteso loading=true all'avvio del modello")
	}
	if m.err != nil {
		t.Errorf("atteso err=nil all'avvio, ottenuto: %v", m.err)
	}
	if len(m.kpi) != 0 {
		t.Errorf("attesa fetta kpi vuota all'avvio, ottenuti %d elementi", len(m.kpi))
	}

	v := m.View()
	if !strings.Contains(v.Content, "Calcolo metriche in corso...") {
		t.Errorf("attesa la vista di caricamento iniziale, ottenuta %q", v.Content)
	}
}

// TestInit verifica che il metodo Init() restituisca un tea.Cmd non nil per l'avvio del caricamento report.
func TestInit(t *testing.T) {
	m := New(nil)
	cmd := m.Init()
	if cmd == nil {
		t.Error("atteso tea.Cmd non nil ritornato da Init()")
	}
}

// TestUpdate_ReportLoadedMsg verifica la corretta transizione di stato alla ricezione
// del messaggio ReportLoadedMsg, con disattivazione del caricamento e popolamento dei KPI.
func TestUpdate_ReportLoadedMsg(t *testing.T) {
	m := New(nil)
	kpis := []api.KPI{
		{Etichetta: "Fatturato netto", Valore: "€ 12400.00"},
		{Etichetta: "Ordini evasi", Valore: "22"},
	}

	res, cmd := m.Update(ReportLoadedMsg{KPI: kpis})
	m = res.(Model)

	if m.loading {
		t.Error("atteso loading=false dopo la ricezione di ReportLoadedMsg")
	}
	if m.err != nil {
		t.Errorf("atteso err=nil dopo il caricamento, ottenuto: %v", m.err)
	}
	if len(m.kpi) != 2 {
		t.Fatalf("attesi 2 KPI nel modello, ottenuti %d", len(m.kpi))
	}
	if cmd != nil {
		t.Error("non atteso alcun comando in risposta a ReportLoadedMsg")
	}
}

// TestUpdate_ErroreMsg verifica che la ricezione di un ErroreMsg interrompa lo stato
// di caricamento e memorizzi l'errore per la visualizzazione.
func TestUpdate_ErroreMsg(t *testing.T) {
	m := New(nil)
	erroreSimulato := errMock{}

	res, cmd := m.Update(ErroreMsg{Err: erroreSimulato})
	m = res.(Model)

	if m.loading {
		t.Error("atteso loading=false dopo la ricezione di ErroreMsg")
	}
	if m.err == nil || m.err.Error() != "errore di connessione al server" {
		t.Errorf("atteso errore memorizzato nel modello, ottenuto: %v", m.err)
	}
	if cmd != nil {
		t.Error("non atteso alcun comando in risposta a ErroreMsg")
	}

	v := m.View()
	if !strings.Contains(v.Content, "Errore: errore di connessione al server") {
		t.Errorf("atteso messaggio di errore formattato nella vista, ottenuto %q", v.Content)
	}
}

// TestUpdate_WindowSizeMsg verifica che i messaggi di ridimensionamento finestra
// vengano gestiti senza alterare lo stato del modello né restituire comandi.
func TestUpdate_WindowSizeMsg(t *testing.T) {
	m := New(nil)
	res, cmd := m.Update(tea.WindowSizeMsg{Width: 100, Height: 40})
	m = res.(Model)

	if !m.loading {
		t.Error("WindowSizeMsg non deve modificare lo stato di loading")
	}
	if cmd != nil {
		t.Error("non atteso alcun comando per WindowSizeMsg")
	}
}

// TestUpdate_RicalcolaConR verifica che la pressione del tasto 'r' imposti `loading = true`
// e restituisca il comando per il ricaricamento delle metriche dal backend.
func TestUpdate_RicalcolaConR(t *testing.T) {
	m := New(nil)
	// Simuliamo un report già caricato
	res, _ := m.Update(ReportLoadedMsg{KPI: []api.KPI{{Etichetta: "Test", Valore: "10"}}})
	m = res.(Model)

	// Pressione del tasto "r"
	res, cmd := m.Update(tea.KeyPressMsg{Text: "r"})
	m = res.(Model)

	if !m.loading {
		t.Error("atteso loading=true dopo aver premuto 'r'")
	}
	if cmd == nil {
		t.Error("atteso un tea.Cmd di ricaricamento non nil dopo 'r'")
	}
}

// TestUpdate_TastiEMessaggiIgnorati controlla che tasti non gestiti o messaggi sconosciuti
// non provochino effetti collaterali né panico nel modello.
func TestUpdate_TastiEMessaggiIgnorati(t *testing.T) {
	m := New(nil)
	res, _ := m.Update(ReportLoadedMsg{KPI: []api.KPI{{Etichetta: "Fatturato", Valore: "100"}}})
	m = res.(Model)

	tastiInvalidi := []string{"a", "down", "enter", "tab", "esc"}
	for _, tasto := range tastiInvalidi {
		res, cmd := m.Update(tea.KeyPressMsg{Text: tasto})
		mod := res.(Model)
		if mod.loading {
			t.Errorf("tasto %q non gestito ha erroneamente attivato il loading", tasto)
		}
		if cmd != nil {
			t.Errorf("tasto %q non gestito ha restituito un comando inatteso", tasto)
		}
	}

	// Tipo di messaggio completamente sconosciuto
	type tipoMessaggioSconosciuto struct{}
	res, cmd := m.Update(tipoMessaggioSconosciuto{})
	mod := res.(Model)
	if mod.loading {
		t.Error("messaggio sconosciuto ha alterato lo stato del modello")
	}
	if cmd != nil {
		t.Error("messaggio sconosciuto ha restituito un comando non nil")
	}
}

// TestView_StatiDistinti testa singolarmente il rendering dei tre stati principali:
// caricamento, errore e report privo di metriche.
func TestView_StatiDistinti(t *testing.T) {
	// 1. Stato di caricamento
	mLoading := New(nil)
	if v := mLoading.View().Content; !strings.Contains(v, "Calcolo metriche in corso...") {
		t.Errorf("vista caricamento errata: %q", v)
	}

	// 2. Stato di errore
	mErr := New(nil)
	resErr, _ := mErr.Update(ErroreMsg{Err: errors.New("backend non raggiungibile")})
	mErr = resErr.(Model)
	if v := mErr.View().Content; !strings.Contains(v, "Errore: backend non raggiungibile") {
		t.Errorf("vista errore errata: %q", v)
	}

	// 3. Stato vuoto (nessun KPI)
	mEmpty := New(nil)
	resEmpty, _ := mEmpty.Update(ReportLoadedMsg{KPI: []api.KPI{}})
	mEmpty = resEmpty.(Model)
	vEmpty := mEmpty.View().Content
	if !strings.Contains(vEmpty, "Nessuna metrica disponibile.") {
		t.Errorf("vista vuota errata: %q", vEmpty)
	}
	if !strings.Contains(vEmpty, "[r] Ricalcola") {
		t.Errorf("footer mancante nella vista vuota: %q", vEmpty)
	}
}

// TestView_GrigliaKPI_PariEDispari verifica il corretto posizionamento nei box della griglia
// sia in caso di numero pari di KPI (coppie perfette) sia in caso di numero dispari (box singolo finale).
func TestView_GrigliaKPI_PariEDispari(t *testing.T) {
	//Caso 1: Numero dispari di KPI (3 elementi -> 1 coppia + 1 singolo finale)
	kpiDispari := []api.KPI{
		{Etichetta: "Fatturato netto", Valore: "€ 1000.00"},
		{Etichetta: "Fatture", Valore: "5"},
		{Etichetta: "Ordini evasi", Valore: "12"},
	}

	mDispari := New(nil)
	res, _ := mDispari.Update(ReportLoadedMsg{KPI: kpiDispari})
	mDispari = res.(Model)

	vDispari := mDispari.View().Content
	for _, k := range kpiDispari {
		if !strings.Contains(vDispari, strings.ToUpper(k.Etichetta)) || !strings.Contains(vDispari, k.Valore) {
			t.Errorf("vista dispari priva del KPI %q (%s):\n%s", k.Etichetta, k.Valore, vDispari)
		}
	}

	// Caso 2: Numero pari di KPI (4 elementi -> 2 coppie perfette)
	kpiPari := []api.KPI{
		{Etichetta: "Fatturato netto", Valore: "€ 2000.00"},
		{Etichetta: "Fatture", Valore: "10"},
		{Etichetta: "Ordini evasi", Valore: "20"},
		{Etichetta: "Note credito", Valore: "2"},
	}

	mPari := New(nil)
	res, _ = mPari.Update(ReportLoadedMsg{KPI: kpiPari})
	mPari = res.(Model)

	vPari := mPari.View().Content
	for _, k := range kpiPari {
		if !strings.Contains(vPari, strings.ToUpper(k.Etichetta)) || !strings.Contains(vPari, k.Valore) {
			t.Errorf("vista pari priva del KPI %q (%s):\n%s", k.Etichetta, k.Valore, vPari)
		}
	}
}

// TestLoadReportCmd_Successo esegue il comando restituito da Init() simulando un server HTTP backend.
// Verifica la corretta aggregazione dei dati (calcolo fatturato al netto delle note di credito).
func TestLoadReportCmd_Successo(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.URL.Path {
		case "/api/v1/orders/invoices":
			_ = json.NewEncoder(w).Encode([]map[string]any{{
				"id": 1, "numero": "FT-001", "ordine_uscita_id": 1,
				"totale": 150.00, "data": "2026-03-01T10:00:00Z",
			}})
		case "/api/v1/orders/credit-notes":
			_ = json.NewEncoder(w).Encode([]map[string]any{{
				"id": 1, "numero": "NC-001", "ordine_uscita_id": 1,
				"totale": 50.00, "data": "2026-04-01T10:00:00Z",
			}})
		case "/api/v1/orders":
			_ = json.NewEncoder(w).Encode([]map[string]any{
				{"id": 1, "numero": "ORD-001", "cliente": "Mario Rossi", "stato": "evaso"},
				{"id": 2, "numero": "ORD-002", "cliente": "Luigi Bianchi", "stato": "confermato"},
			})
		case "/api/v1/products":
			_ = json.NewEncoder(w).Encode([]map[string]any{})
		case "/api/v1/categories":
			_ = json.NewEncoder(w).Encode([]map[string]any{})
		default:
			t.Errorf("endpoint non previsto dal mock server: %s", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer srv.Close()

	m := New(api.New(srv.URL))

	cmd := m.Init()
	msg := cmd()

	loaded, ok := msg.(ReportLoadedMsg)
	if !ok {
		t.Fatalf("atteso ReportLoadedMsg da loadReportCmd, ottenuto %T (%v)", msg, msg)
	}

	valori := make(map[string]string, len(loaded.KPI))
	for _, k := range loaded.KPI {
		valori[k.Etichetta] = k.Valore
	}

	casiAttesi := map[string]string{
		"Fatturato netto": "€ 100.00", // 150 (fattura) - 50 (nota di credito)
		"Fatture":         "1",
		"Ordini evasi":    "1",
		"Note credito":    "1",
	}

	for etichetta, atteso := range casiAttesi {
		val, ok := valori[etichetta]
		if !ok {
			t.Errorf("metrica %q mancante nel risultato di GetReportKPI", etichetta)
			continue
		}
		if val != atteso {
			t.Errorf("KPI %q: atteso %q, ottenuto %q", etichetta, atteso, val)
		}
	}

	res, _ := m.Update(loaded)
	m = res.(Model)
	if v := m.View().Content; !strings.Contains(v, "FATTURATO NETTO") {
		t.Errorf("attesa l'etichetta FATTURATO NETTO nella vista finale, ottenuta %q", v)
	}
}

// TestLoadReportCmd_ErroreServer verifica che le risposte di errore dal server HTTP (es. 500)
// vengano correttamente intercettate e convertite in ErroreMsg.
func TestLoadReportCmd_ErroreServer(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	m := New(api.New(srv.URL))

	cmd := m.Init()
	msg := cmd()

	errMsg, ok := msg.(ErroreMsg)
	if !ok {
		t.Fatalf("atteso ErroreMsg a seguito di risposta HTTP 500, ottenuto %T", msg)
	}
	if errMsg.Err == nil {
		t.Error("atteso errore non nil all'interno di ErroreMsg")
	}

	res, _ := m.Update(errMsg)
	m = res.(Model)
	if m.loading || m.err == nil {
		t.Error("atteso loading=false ed errore registrato nel modello a seguito del fallimento")
	}
}