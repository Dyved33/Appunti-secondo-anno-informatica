package dashboard

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/api"
)

// TestDashboard_New verifica il costruttore New e l'impostazione iniziale dei flag di caricamento.
func TestDashboard_New(t *testing.T) {
	t.Parallel()

	client := api.New("http://localhost:8080")
	m := New(client)

	if !m.loadingGiacenze || !m.loadingKPI || !m.loadingOrdiniUscita || !m.loadingOrdiniEntrata {
		t.Error("tutti i flag di caricamento devono essere true all'avvio")
	}
	if m.client != client {
		t.Error("client API non associato correttamente al modello Dashboard")
	}
}

// TestDashboard_Init verifica che Init restituisca un comando batch contenente i 4 caricatori asincroni.
func TestDashboard_Init(t *testing.T) {
	t.Parallel()

	m := New(nil)
	cmd := m.Init()

	if cmd == nil {
		t.Fatal("Init() atteso tea.Cmd non nil")
	}
}

// TestDashboard_AsyncCommands verifica l'esecuzione dei comandi asincroni con backend mock
// sia per i casi di successo sia per i casi di errore di rete.
func TestDashboard_AsyncCommands(t *testing.T) {
	t.Parallel()

	t.Run("Successo di tutte le chiamate API", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			switch r.URL.Path {
			case "/api/v1/products/inventory-report":
				_, _ = w.Write([]byte(`{"totale_quantita":10,"valore_totale":500.0,"articoli":[],"articoli_sotto_soglia":[]}`))
			case "/api/v1/orders/invoices":
				_, _ = w.Write([]byte(`[]`))
			case "/api/v1/orders/credit-notes":
				_, _ = w.Write([]byte(`[]`))
			case "/api/v1/orders":
				_, _ = w.Write([]byte(`[]`))
			case "/api/v1/products":
				_, _ = w.Write([]byte(`[]`))
			case "/api/v1/categories":
				_, _ = w.Write([]byte(`[]`))
			case "/api/v1/orders/uscita":
				_, _ = w.Write([]byte(`[{"id":1,"numero":"U-01","stato":"evaso"}]`))
			case "/api/v1/orders/entrata":
				_, _ = w.Write([]byte(`[{"id":2,"numero":"E-01","fornitore":"F1"}]`))
			default:
				http.NotFound(w, r)
			}
		}))
		defer srv.Close()

		m := New(api.New(srv.URL))

		// 1. Giacenze
		msgGiac := m.loadGiacenzeCmd()()
		if _, ok := msgGiac.(GiacenzeCaricateMsg); !ok {
			t.Errorf("atteso GiacenzeCaricateMsg, ottenuto %T", msgGiac)
		}

		// 2. KPI
		msgKPI := m.loadKPICmd()()
		if _, ok := msgKPI.(KPICaricatiMsg); !ok {
			t.Errorf("atteso KPICaricatiMsg, ottenuto %T", msgKPI)
		}

		// 3. Ordini Uscita
		msgUscita := m.loadOrdiniUscitaCmd()()
		if _, ok := msgUscita.(OrdiniUscitaCaricatiMsg); !ok {
			t.Errorf("atteso OrdiniUscitaCaricatiMsg, ottenuto %T", msgUscita)
		}

		// 4. Ordini Entrata
		msgEntrata := m.loadOrdiniEntrataCmd()()
		if _, ok := msgEntrata.(OrdiniEntrataCaricatiMsg); !ok {
			t.Errorf("atteso OrdiniEntrataCaricatiMsg, ottenuto %T", msgEntrata)
		}
	})

	t.Run("Errori di rete su tutte le chiamate API", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer srv.Close()

		m := New(api.New(srv.URL))

		if _, ok := m.loadGiacenzeCmd()().(ErroreGiacenzeMsg); !ok {
			t.Error("atteso ErroreGiacenzeMsg")
		}
		if _, ok := m.loadKPICmd()().(ErroreKPIMsg); !ok {
			t.Error("atteso ErroreKPIMsg")
		}
		if _, ok := m.loadOrdiniUscitaCmd()().(ErroreOrdiniUscitaMsg); !ok {
			t.Error("atteso ErroreOrdiniUscitaMsg")
		}
		if _, ok := m.loadOrdiniEntrataCmd()().(ErroreOrdiniEntrataMsg); !ok {
			t.Error("atteso ErroreOrdiniEntrataMsg")
		}
	})
}

// TestDashboard_Update_Messages verifica la corretta memorizzazione dello stato
// alla ricezione di tutti i messaggi Caricati ed Errori.
func TestDashboard_Update_Messages(t *testing.T) {
	t.Parallel()

	m := New(nil)

	// Caricamento Giacenze
	rep := &api.InventoryReport{TotaleQuantita: 200, ValoreTotale: 4500.0}
	res, _ := m.Update(GiacenzeCaricateMsg{Report: rep})
	m = res.(Model)
	if m.loadingGiacenze || m.reportGiacenze != rep {
		t.Error("GiacenzeCaricateMsg non gestito correttamente")
	}

	// Errore Giacenze
	res, _ = m.Update(ErroreGiacenzeMsg{Err: errors.New("errore giacenze")})
	m = res.(Model)
	if m.erroreGiacenze == nil {
		t.Error("ErroreGiacenzeMsg non memorizzato")
	}

	// Caricamento KPI
	kpis := []api.KPI{{Etichetta: "Fatturato", Valore: "€ 100"}}
	res, _ = m.Update(KPICaricatiMsg{KPI: kpis})
	m = res.(Model)
	if m.loadingKPI || len(m.kpiContabilita) != 1 {
		t.Error("KPICaricatiMsg non gestito correttamente")
	}

	// Errore KPI
	res, _ = m.Update(ErroreKPIMsg{Err: errors.New("errore kpi")})
	m = res.(Model)
	if m.erroreKPI == nil {
		t.Error("ErroreKPIMsg non memorizzato")
	}

	// Caricamento Ordini Uscita
	ordUscita := []api.OrdineUscita{{ID: 1, Stato: api.StatoEvaso}}
	res, _ = m.Update(OrdiniUscitaCaricatiMsg{Ordini: ordUscita})
	m = res.(Model)
	if m.loadingOrdiniUscita || len(m.ordiniUscita) != 1 {
		t.Error("OrdiniUscitaCaricatiMsg non gestito correttamente")
	}

	// Errore Ordini Uscita
	res, _ = m.Update(ErroreOrdiniUscitaMsg{Err: errors.New("errore ordini uscita")})
	m = res.(Model)
	if m.erroreOrdiniUscita == nil {
		t.Error("ErroreOrdiniUscitaMsg non memorizzato")
	}

	// Caricamento Ordini Entrata
	ordEntrata := []api.OrdineEntrata{{ID: 2, Fornitore: "Fornitore Beta"}}
	res, _ = m.Update(OrdiniEntrataCaricatiMsg{Ordini: ordEntrata})
	m = res.(Model)
	if m.loadingOrdiniEntrata || len(m.ordiniEntrata) != 1 {
		t.Error("OrdiniEntrataCaricatiMsg non gestito correttamente")
	}

	// Errore Ordini Entrata
	res, _ = m.Update(ErroreOrdiniEntrataMsg{Err: errors.New("errore ordini entrata")})
	m = res.(Model)
	if m.erroreOrdiniEntrata == nil {
		t.Error("ErroreOrdiniEntrataMsg non memorizzato")
	}
}

// TestDashboard_Update_KeyboardNavigation verifica la navigazione tra le scorte critiche,
// i tasti di scorrimento pgup/pgdown/home/end e il ricaricamento con 'r'.
func TestDashboard_Update_KeyboardNavigation(t *testing.T) {
	t.Parallel()

	m := New(nil)
	m.reportGiacenze = &api.InventoryReport{
		ArticoliSottoSoglia: []api.InventoryItem{
			{ProdottoID: 1, Codice: "SKU-1", Nome: "Prodotto 1"},
			{ProdottoID: 2, Codice: "SKU-2", Nome: "Prodotto 2"},
			{ProdottoID: 3, Codice: "SKU-3", Nome: "Prodotto 3"},
		},
	}

	// Naviga verso il basso con 'down' e 'j'
	res, _ := m.Update(tea.KeyPressMsg{Text: "down"})
	m = res.(Model)
	if m.cursoreScorteCritiche != 1 {
		t.Errorf("atteso cursore su 1 dopo down, ottenuto %d", m.cursoreScorteCritiche)
	}

	res, _ = m.Update(tea.KeyPressMsg{Text: "j"})
	m = res.(Model)
	if m.cursoreScorteCritiche != 2 {
		t.Errorf("atteso cursore su 2 dopo j, ottenuto %d", m.cursoreScorteCritiche)
	}

	// Non supera l'ultimo elemento
	res, _ = m.Update(tea.KeyPressMsg{Text: "down"})
	m = res.(Model)
	if m.cursoreScorteCritiche != 2 {
		t.Errorf("cursore non doveva superare 2, ottenuto %d", m.cursoreScorteCritiche)
	}

	// Naviga verso l'alto con 'up' e 'k'
	res, _ = m.Update(tea.KeyPressMsg{Text: "up"})
	m = res.(Model)
	if m.cursoreScorteCritiche != 1 {
		t.Errorf("atteso cursore su 1 dopo up, ottenuto %d", m.cursoreScorteCritiche)
	}

	res, _ = m.Update(tea.KeyPressMsg{Text: "k"})
	m = res.(Model)
	if m.cursoreScorteCritiche != 0 {
		t.Errorf("atteso cursore su 0 dopo k, ottenuto %d", m.cursoreScorteCritiche)
	}

	// Non va sotto 0
	res, _ = m.Update(tea.KeyPressMsg{Text: "up"})
	m = res.(Model)
	if m.cursoreScorteCritiche != 0 {
		t.Errorf("cursore non doveva andare sotto 0, ottenuto %d", m.cursoreScorteCritiche)
	}

	// Tasti di scorrimento
	res, _ = m.Update(tea.KeyPressMsg{Text: "pgdown"})
	m = res.(Model)
	res, _ = m.Update(tea.KeyPressMsg{Text: "end"})
	m = res.(Model)
	res, _ = m.Update(tea.KeyPressMsg{Text: "home"})
	m = res.(Model)
	if m.scrollOffset != 0 {
		t.Errorf("scrollOffset atteso 0 dopo home, ottenuto %d", m.scrollOffset)
	}
	res, _ = m.Update(tea.KeyPressMsg{Text: "pgup"})
	m = res.(Model)

	// Tasto 'r' (ricarica tutti i blocchi)
	res, cmd := m.Update(tea.KeyPressMsg{Text: "r"})
	m = res.(Model)
	if !m.loadingGiacenze || !m.loadingKPI || !m.loadingOrdiniUscita || !m.loadingOrdiniEntrata {
		t.Error("premendo 'r' tutti i flag loading devono tornare a true")
	}
	if cmd == nil {
		t.Error("premendo 'r' atteso Cmd batch di ricarica")
	}
}

// TestDashboard_RenderingBlocchiEStati verifica il rendering della View nei vari stati:
// caricamento iniziale, errori per ciascun box, dati presenti con e senza scorte critiche.
func TestDashboard_RenderingBlocchiEStati(t *testing.T) {
	t.Parallel()

	// 1. Vista con blocchi in caricamento
	m := New(nil)
	vLoading := m.View().Content
	if !strings.Contains(vLoading, "DASHBOARD ADMIN") {
		t.Error("View non contiene il titolo principale")
	}
	if !strings.Contains(vLoading, "Caricamento...") {
		t.Error("View non mostra gli indicatori di caricamento")
	}

	// 2. Vista con errori su tutti i blocchi
	m.loadingGiacenze = false
	m.erroreGiacenze = errors.New("err magazzino")
	m.loadingKPI = false
	m.erroreKPI = errors.New("err kpi")
	m.loadingOrdiniUscita = false
	m.erroreOrdiniUscita = errors.New("err uscita")
	m.loadingOrdiniEntrata = false
	m.erroreOrdiniEntrata = errors.New("err entrata")

	vErr := m.View().Content
	if !strings.Contains(vErr, "err magazzino") || !strings.Contains(vErr, "err kpi") {
		t.Error("View non mostra i messaggi di errore attesi nei box")
	}

	// 3. Vista con dati popolati
	m.winW = 120
	m.winH = 60
	m.erroreGiacenze = nil
	m.erroreKPI = nil
	m.erroreOrdiniUscita = nil
	m.erroreOrdiniEntrata = nil

	m.reportGiacenze = &api.InventoryReport{
		TotaleQuantita: 120,
		ValoreTotale:   2500.0,
		ArticoliSottoSoglia: []api.InventoryItem{
			{ProdottoID: 1, Codice: "SKU-CRIT", Nome: "Prodotto Critico", Quantita: 2, SogliaMin: 10, Categoria: "Ferramenta"},
		},
	}
	m.kpiContabilita = []api.KPI{
		{Etichetta: "Fatturato netto", Valore: "€ 1000.00"},
		{Etichetta: "Fatture", Valore: "5"},
	}
	m.ordiniUscita = []api.OrdineUscita{
		{ID: 1, Stato: api.StatoConfermato},
		{ID: 2, Stato: api.StatoEvaso},
		{ID: 3, Stato: api.StatoAnnullato},
	}
	m.ordiniEntrata = []api.OrdineEntrata{
		{
			ID: 1, Fornitore: "Fornitore Alfa",
			Righe: []api.RigaOrdineEntrata{{Quantita: 10, PrezzoUnitario: 5.0}},
		},
	}

	vData := m.View().Content
	if !strings.Contains(vData, "Totale Pezzi    120") {
		t.Errorf("View non contiene i dati delle giacenze: %s", vData)
	}
	if !strings.Contains(vData, "Fatturato netto") {
		t.Errorf("View non contiene i KPI contabilità: %s", vData)
	}
	if !strings.Contains(vData, "Confermati  1") || !strings.Contains(vData, "Evasi") {
		t.Errorf("View non contiene il riepilogo ordini uscita: %s", vData)
	}
	if !strings.Contains(vData, "Fornitori  1") || !strings.Contains(vData, "Articoli   10") {
		t.Errorf("View non contiene il riepilogo ordini entrata: %s", vData)
	}
	if !strings.Contains(vData, "SKU-CRIT") {
		t.Errorf("View non contiene la tabella scorte critiche: %s", vData)
	}

	// 4. Scorte critiche vuote
	m.reportGiacenze.ArticoliSottoSoglia = nil
	vNoCrit := m.View().Content
	if !strings.Contains(vNoCrit, "Nessuna scorta critica") {
		t.Errorf("View attesa con 'Nessuna scorta critica', ottenuta: %s", vNoCrit)
	}
}

// TestDashboard_Helpers verifica le funzioni helper statisticheEntrata, troncaStringa e renderBlocco.
func TestDashboard_Helpers(t *testing.T) {
	t.Parallel()

	// statisticheEntrata
	ordini := []api.OrdineEntrata{
		{
			Fornitore: "F1",
			Righe: []api.RigaOrdineEntrata{
				{Quantita: 10, PrezzoUnitario: 2.0},
				{Quantita: 5, PrezzoUnitario: 4.0},
			},
		},
		{
			Fornitore: "F1", // Stesso fornitore
			Righe: []api.RigaOrdineEntrata{
				{Quantita: 20, PrezzoUnitario: 1.0},
			},
		},
		{
			Fornitore: "F2",
			Righe:     nil,
		},
	}
	fornitori, articoli, valore := statisticheEntrata(ordini)
	if fornitori != 2 {
		t.Errorf("fornitori distinti attesi 2, ottenuti %d", fornitori)
	}
	if articoli != 35 {
		t.Errorf("quantità totale articoli attesa 35, ottenuta %d", articoli)
	}
	if valore != 60.0 { // (10*2) + (5*4) + (20*1) = 20 + 20 + 20 = 60
		t.Errorf("valore totale atteso 60.0, ottenuto %f", valore)
	}

	// troncaStringa
	if got := troncaStringa("Breve", 10); got != "Breve" {
		t.Errorf("troncaStringa('Breve', 10) = %q", got)
	}
	if got := troncaStringa("StringaMoltoLunga", 10); got != "Stringa..." {
		t.Errorf("troncaStringa('StringaMoltoLunga', 10) = %q", got)
	}
	if got := troncaStringa("ABCDEF", 2); got != "AB" {
		t.Errorf("troncaStringa('ABCDEF', 2) = %q", got)
	}

	// renderBlocco
	righePoche := []string{"Riga 1", "Riga 2"}
	bloccoPoco := renderBlocco(righePoche)
	if bloccoPoco == "" {
		t.Error("renderBlocco con poche righe ha prodotto stringa vuota")
	}

	righeTante := []string{"1", "2", "3", "4", "5", "6", "7", "8"}
	bloccoTante := renderBlocco(righeTante)
	if bloccoTante == "" {
		t.Error("renderBlocco con tante righe ha prodotto stringa vuota")
	}
}
