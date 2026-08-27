// Package report raccoglie i test d'unità della sotto-pagina Report.
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

// TestNew_StatoIniziale verifica che il report parta in stato di caricamento.
func TestNew_StatoIniziale(t *testing.T) {
	m := New(nil)
	if !m.loading {
		t.Error("atteso loading=true all'avvio")
	}
	if v := m.View(); !strings.Contains(v.Content, "Calcolo metriche") {
		t.Errorf("attesa la view di caricamento, ottenuta %q", v.Content)
	}
}

// TestUpdate_ReportLoadedMsg verifica che il caricamento popoli le metriche KPI.
func TestUpdate_ReportLoadedMsg(t *testing.T) {
	m := New(nil)
	res, _ := m.Update(ReportLoadedMsg{KPI: []api.KPI{
		{Etichetta: "Fatturato netto", Valore: "€ 12400.00"},
		{Etichetta: "Ordini evasi", Valore: "22"},
	}})
	m = res.(Model)

	if m.loading {
		t.Error("atteso loading=false dopo il caricamento")
	}
	if len(m.kpi) != 2 {
		t.Errorf("attesi 2 KPI, ottenuti %d", len(m.kpi))
	}
}

// TestUpdate_ErroreMsg verifica che un errore async venga registrato correttamente.
func TestUpdate_ErroreMsg(t *testing.T) {
	m := New(nil)
	res, _ := m.Update(ErroreMsg{Err: errors.New("dati non disponibili")})
	m = res.(Model)

	if m.err == nil || m.loading {
		t.Error("atteso err impostato e loading=false")
	}
	if v := m.View(); !strings.Contains(v.Content, "dati non disponibili") {
		t.Errorf("attesa la view di errore contenente il messaggio, ottenuta %q", v.Content)
	}
}

// TestUpdate_RicalcolaConR verifica che il tasto "r" rimetta il report in
// caricamento e restituisca un Cmd per ricalcolare le metriche.
func TestUpdate_RicalcolaConR(t *testing.T) {
	m := New(nil)
	res, _ := m.Update(ReportLoadedMsg{KPI: []api.KPI{{Etichetta: "X", Valore: "1"}}})
	m = res.(Model)

	res, cmd := m.Update(tea.KeyPressMsg{Text: "r"})
	m = res.(Model)

	if !m.loading {
		t.Error("atteso loading=true dopo 'r'")
	}
	if cmd == nil {
		t.Error("atteso un Cmd di ricaricamento dopo 'r'")
	}
}

// TestUpdate_TastoIgnorato verifica che un tasto non gestito lasci il modello
// invariato: la pagina è di sola lettura e non ha navigazione interna.
func TestUpdate_TastoIgnorato(t *testing.T) {
	m := New(nil)
	res, _ := m.Update(ReportLoadedMsg{KPI: []api.KPI{{Etichetta: "X", Valore: "1"}}})
	m = res.(Model)

	res, cmd := m.Update(tea.KeyPressMsg{Text: "down"})
	m = res.(Model)

	if m.loading {
		t.Error("un tasto non gestito non deve rimettere la pagina in caricamento")
	}
	if cmd != nil {
		t.Error("non atteso alcun Cmd per un tasto non gestito")
	}
}

// TestView_MostraEtichetteKPI verifica che la view contenga etichetta e valore
// di ciascun KPI caricato.
func TestView_MostraEtichetteKPI(t *testing.T) {
	m := New(nil)
	res, _ := m.Update(ReportLoadedMsg{KPI: []api.KPI{{Etichetta: "Ordini evasi", Valore: "22"}}})
	m = res.(Model)

	v := m.View()
	if !strings.Contains(v.Content, "ORDINI EVASI") || !strings.Contains(v.Content, "22") {
		t.Errorf("attesa l'etichetta e il valore del KPI nella view, ottenuta %q", v.Content)
	}
}

// TestView_NessunaMetrica verifica che un report senza KPI non produca una view
// vuota né un panic.
func TestView_NessunaMetrica(t *testing.T) {
	m := New(nil)
	res, _ := m.Update(ReportLoadedMsg{KPI: nil})
	m = res.(Model)

	if v := m.View(); !strings.Contains(v.Content, "Nessuna metrica") {
		t.Errorf("atteso il messaggio di report vuoto, ottenuta %q", v.Content)
	}
}

// TestLoadReportCmd_Successo esegue il Cmd di Init() contro un backend di prova
// e verifica le metriche calcolate da api.GetReportKPI: il fatturato è al netto
// delle note di credito e queste ultime non contano come fatture.
func TestLoadReportCmd_Successo(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.URL.Path {
		case "/api/v1/orders/invoices":
			json.NewEncoder(w).Encode([]map[string]any{{
				"id": 1, "numero": "FT-001", "ordine_uscita_id": 1,
				"totale": 120.00, "data": "2026-03-01T10:00:00Z",
			}})
		case "/api/v1/orders/credit-notes":
			json.NewEncoder(w).Encode([]map[string]any{{
				"id": 2, "numero": "NC-001", "ordine_uscita_id": 1,
				"totale": 20.00, "data": "2026-04-01T10:00:00Z",
			}})
		case "/api/v1/orders":
			json.NewEncoder(w).Encode([]map[string]any{
				{"id": 1, "numero": "ORD-001", "cliente": "Mario Rossi", "stato": "evaso"},
				{"id": 2, "numero": "ORD-002", "cliente": "Luigi Bianchi", "stato": "confermato"},
			})
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
	loaded, ok := msg.(ReportLoadedMsg)
	if !ok {
		t.Fatalf("atteso ReportLoadedMsg, ottenuto %T (%v)", msg, msg)
	}

	valori := make(map[string]string, len(loaded.KPI))
	for _, k := range loaded.KPI {
		valori[k.Etichetta] = k.Valore
	}

	casi := map[string]string{
		"Fatturato netto": "€ 100.00", // 120 di fatture - 20 di nota di credito
		"Fatture":         "1",
		"Ordini evasi":    "1",
		"Note credito":    "1",
	}
	for etichetta, atteso := range casi {
		if valori[etichetta] != atteso {
			t.Errorf("KPI %q: atteso %q, ottenuto %q", etichetta, atteso, valori[etichetta])
		}
	}

	res, _ := m.Update(msg)
	m = res.(Model)
	if v := m.View(); !strings.Contains(v.Content, "FATTURATO NETTO") {
		t.Errorf("attesa la metrica caricata nella view, ottenuta %q", v.Content)
	}
}

// TestLoadReportCmd_ErroreServer verifica che un errore HTTP diventi un
// ErroreMsg e chiuda lo stato di caricamento.
func TestLoadReportCmd_ErroreServer(t *testing.T) {
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
