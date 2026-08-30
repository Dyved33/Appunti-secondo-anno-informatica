package fatture

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/tommasogarofani/ingsw_25-26_group-11/client/api"

	tea "charm.land/bubbletea/v2"
)

// TestNew_StatoIniziale verifica che il modello parta in stato di caricamento
// con il cursore azzerato.
func TestNew_StatoIniziale(t *testing.T) {
	m := New(nil)
	if !m.loading {
		t.Error("atteso loading=true all'avvio")
	}
	if m.cursore != 0 {
		t.Errorf("atteso cursore iniziale a 0, ottenuto %d", m.cursore)
	}
}

// TestUpdate_FattureLoadedMsg verifica che il caricamento popoli correttamente
// sia il registro delle fatture che quello delle note di credito.
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

// TestUpdate_ErroreMsg verifica che la gestione dell'errore asincrono interrompa
// lo stato di caricamento e salvi l'errore nel modello.
func TestUpdate_ErroreMsg(t *testing.T) {
	m := New(nil)
	res, _ := m.Update(ErroreMsg{Err: errors.New("timeout di rete")})
	m = res.(Model)

	if m.err == nil || m.loading {
		t.Error("atteso err impostato e loading=false")
	}
	if m.err.Error() != "timeout di rete" {
		t.Errorf("errore inatteso: %v", m.err)
	}
}

// TestUpdate_RicaricaConR verifica che la pressione del tasto 'r' imposti
// lo stato in caricamento e restituisca il comando per ricaricare le fatture.
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

// ── Navigazione e Tab ─────────────────────────────────────────────────────

// TestUpdate_NavigazioneCursore verifica il movimento del cursore con i tasti freccia e j/k.
func TestUpdate_NavigazioneCursore(t *testing.T) {
	m := caricaFattureDiProva(t, 3)

	res, _ := m.Update(tea.KeyPressMsg{Text: "down"})
	m = res.(Model)
	if m.cursore != 1 {
		t.Errorf("atteso cursore a 1 dopo 'down', ottenuto %d", m.cursore)
	}

	res, _ = m.Update(tea.KeyPressMsg{Text: "j"})
	m = res.(Model)
	if m.cursore != 2 {
		t.Errorf("atteso cursore a 2 dopo 'j', ottenuto %d", m.cursore)
	}

	res, _ = m.Update(tea.KeyPressMsg{Text: "up"})
	m = res.(Model)
	if m.cursore != 1 {
		t.Errorf("atteso cursore a 1 dopo 'up', ottenuto %d", m.cursore)
	}

	res, _ = m.Update(tea.KeyPressMsg{Text: "k"})
	m = res.(Model)
	if m.cursore != 0 {
		t.Errorf("atteso cursore a 0 dopo 'k', ottenuto %d", m.cursore)
	}
}

// TestUpdate_CursoreNonEsceDaiLimiti verifica che il cursore non vada mai
// sotto 0 né oltre la dimensione della lista visualizzata.
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

// TestUpdate_NavigazionePaginazione verifica la navigazione a blocchi di pagine (pgdown/pgup).
func TestUpdate_NavigazionePaginazione(t *testing.T) {
	m := caricaFattureDiProva(t, 25)

	res, _ := m.Update(tea.KeyPressMsg{Text: "pgdown"})
	m = res.(Model)
	if m.cursore <= 0 {
		t.Errorf("pgdown doveva avanzare il cursore, ottenuto %d", m.cursore)
	}

	res, _ = m.Update(tea.KeyPressMsg{Text: "pgdown"})
	res, _ = m.Update(tea.KeyPressMsg{Text: "pgdown"})
	m = res.(Model)
	if m.cursore != 24 {
		t.Errorf("atteso cursore al massimo valore (24), ottenuto %d", m.cursore)
	}

	res, _ = m.Update(tea.KeyPressMsg{Text: "pgup"})
	m = res.(Model)
	if m.cursore >= 24 {
		t.Errorf("pgup doveva retrocedere il cursore, ottenuto %d", m.cursore)
	}
}

// TestUpdate_CursoreResettatoAlRicaricamento verifica che la ricezione dei dati ricalibri il cursore.
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

// TestUpdate_CambioTab verifica l'alternanza tra la sezione Fatture e Note di Credito.
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

// ── Creazione Fattura e Nota di Credito ───────────────────────────────────

// TestUpdate_CreazioneFattura_FlussoCompleto verifica il ciclo completo di creazione
// di una nuova fattura da ordine evaso, compresi i controlli sui vincoli di sezione.
func TestUpdate_CreazioneFattura_FlussoCompleto(t *testing.T) {
	m := caricaFattureDiProva(t, 1)

	// Verifico che ctrl+n venga ignorato se ci troviamo nella sezione SEZ_NOTE_CREDITO
	mTab := m
	mTab.sezione = SEZ_NOTE_CREDITO
	mTab.mode = modeView
	resTab, cmdTab := mTab.Update(tea.KeyPressMsg{Text: "ctrl+n"})
	if cmdTab != nil || resTab.(Model).mode == modeCreateFattura {
		t.Error("ctrl+n non deve fare nulla se la sezione è SEZ_NOTE_CREDITO")
	}

	// 1. Avvio creazione con ctrl+n su SEZ_FATTURE
	res, cmd := m.Update(tea.KeyPressMsg{Text: "ctrl+n"})
	m = res.(Model)
	if m.mode != modeCreateFattura {
		t.Errorf("atteso modeCreateFattura, ottenuto %v", m.mode)
	}
	if cmd == nil {
		t.Error("atteso Cmd per caricamento ordini evasi")
	}

	// 2. Simulo la ricezione degli ordini evasi
	ordini := []api.OrdineUscita{
		{ID: 10, Numero: "ORD-1", Cliente: "Cliente A", Data: time.Now()},
		{ID: 11, Numero: "ORD-2", Cliente: "Cliente B", Data: time.Now()},
	}
	res, _ = m.Update(ordiniEvasiLoadedMsg{Ordini: ordini})
	m = res.(Model)
	if len(m.ordiniEvasi) != 2 {
		t.Fatalf("attesi 2 ordini evasi, ottenuti %d", len(m.ordiniEvasi))
	}

	// 3. Navigazione nella lista ordini evasi (down / up)
	res, _ = m.Update(tea.KeyPressMsg{Text: "down"})
	m = res.(Model)
	if m.fatturaCursore != 1 {
		t.Errorf("atteso fatturaCursore a 1, ottenuto %d", m.fatturaCursore)
	}

	res, _ = m.Update(tea.KeyPressMsg{Text: "up"})
	m = res.(Model)
	if m.fatturaCursore != 0 {
		t.Errorf("atteso fatturaCursore a 0, ottenuto %d", m.fatturaCursore)
	}

	// 4. Annullamento tramite ESC
	resEsc, _ := m.Update(tea.KeyPressMsg{Text: "esc"})
	if resEsc.(Model).mode != modeView {
		t.Error("esc doveva ripristinare modeView")
	}

	// 5. Invio conferma con ENTER
	res, cmd = m.Update(tea.KeyPressMsg{Text: "enter"})
	m = res.(Model)
	if !m.loading {
		t.Error("atteso loading=true durante l'invio della fattura")
	}
	if cmd == nil {
		t.Error("atteso creaFatturaCmd dopo enter")
	}

	// 6. Messaggio di fattura creata e ricaricamento registro
	res, cmd = m.Update(fatturaCreataMsg{})
	m = res.(Model)
	if m.mode != modeView {
		t.Error("atteso ritorno a modeView dopo fatturaCreataMsg")
	}
	if m.success == "" {
		t.Error("atteso messaggio di successo valorizzato")
	}
	if cmd == nil {
		t.Error("atteso Cmd per ricaricare le fatture dopo la creazione")
	}
}

// TestUpdate_CreazioneNotaCredito_FlussoCompleto verifica la creazione di una nota di credito.
func TestUpdate_CreazioneNotaCredito_FlussoCompleto(t *testing.T) {
	m := caricaFattureDiProvaConNoteCredito(t)

	// Avvio con ctrl+k
	res, cmd := m.Update(tea.KeyPressMsg{Text: "ctrl+k"})
	m = res.(Model)
	if m.mode != modeCreateNotaCredito {
		t.Errorf("atteso modeCreateNotaCredito, ottenuto %v", m.mode)
	}
	if cmd == nil {
		t.Error("atteso Cmd per caricamento ordini evasi")
	}

	// Popolamento ordini
	ordini := []api.OrdineUscita{{ID: 20, Numero: "ORD-NC", Cliente: "Cliente C", Data: time.Now()}}
	res, _ = m.Update(ordiniEvasiLoadedMsg{Ordini: ordini})
	m = res.(Model)

	// Annullamento esc
	mEsc, _ := m.Update(tea.KeyPressMsg{Text: "esc"})
	if mEsc.(Model).mode != modeView {
		t.Error("esc doveva ripristinare modeView da creazione nota credito")
	}

	// Invio con ENTER
	res, cmd = m.Update(tea.KeyPressMsg{Text: "enter"})
	m = res.(Model)
	if !m.loading || cmd == nil {
		t.Error("atteso loading e creaNotaCreditoCmd su enter")
	}

	// Ricezione esito positivo nota di credito
	res, cmd = m.Update(notaCreditoCreataMsg{})
	m = res.(Model)
	if m.mode != modeView || m.sezione != SEZ_NOTE_CREDITO {
		t.Error("atteso modeView e passaggio automatico alla sezione SEZ_NOTE_CREDITO")
	}
	if m.success == "" || cmd == nil {
		t.Error("atteso messaggio di successo e Cmd di ricaricamento")
	}
}

// TestUpdate_CreazioneConListaOrdiniVuota garantisce l'assenza di crash ed evita l'invio
// di comandi quando si preme enter su liste ordini non popolate.
func TestUpdate_CreazioneConListaOrdiniVuota(t *testing.T) {
	m := New(nil)
	m.mode = modeCreateFattura
	m.ordiniEvasi = nil

	res, cmd := m.Update(tea.KeyPressMsg{Text: "enter"})
	if cmd != nil {
		t.Error("enter con lista ordini vuota non deve generare comandi")
	}

	m.mode = modeCreateNotaCredito
	res, cmd = m.Update(tea.KeyPressMsg{Text: "enter"})
	if cmd != nil {
		t.Error("enter con lista ordini vuota per NC non deve generare comandi")
	}
	_ = res
}

// ── Eliminazione Fattura ──────────────────────────────────────────────────

// TestUpdate_EliminazioneFattura_FlussoCompleto testa la conferma e cancellazione di una fattura.
func TestUpdate_EliminazioneFattura_FlussoCompleto(t *testing.T) {
	m := caricaFattureDiProva(t, 2)
	m.fatture[0].ID = 99
	m.fatture[0].Numero = "FT-99"

	// 1. Inizio eliminazione con ctrl+d
	res, _ := m.Update(tea.KeyPressMsg{Text: "ctrl+d"})
	m = res.(Model)
	if m.mode != modeConfirmDeleteFattura {
		t.Errorf("atteso modeConfirmDeleteFattura, ottenuto %v", m.mode)
	}
	if m.confirmDeleteID != 99 || m.confirmDeleteNome != "FT-99" {
		t.Errorf("dati conferma eliminazione errati: ID=%d, Nome=%s", m.confirmDeleteID, m.confirmDeleteNome)
	}

	// 2. Annullamento con ESC
	resEsc, _ := m.Update(tea.KeyPressMsg{Text: "esc"})
	mEsc := resEsc.(Model)
	if mEsc.mode != modeView || mEsc.confirmDeleteID != 0 {
		t.Error("esc doveva azzerare lo stato di cancellazione e tornare a modeView")
	}

	// 3. Conferma con ENTER
	res, cmd := m.Update(tea.KeyPressMsg{Text: "enter"})
	m = res.(Model)
	if !m.loading || cmd == nil {
		t.Error("enter doveva avviare l'eliminazione async")
	}

	// 4. Esito eliminazione
	res, cmd = m.Update(fatturaEliminataMsg{})
	m = res.(Model)
	if m.mode != modeView || m.success == "" || cmd == nil {
		t.Error("atteso ripristino vista, messaggio successo e ricaricamento registro")
	}
}

// TestUpdate_EliminazioneFattura_EdgeCases verifica le condizioni limite di eliminazione.
func TestUpdate_EliminazioneFattura_EdgeCases(t *testing.T) {
	// Caso 1: ctrl+d su tab Note Credito -> ignorato
	mNC := caricaFattureDiProvaConNoteCredito(t)
	mNC.sezione = SEZ_NOTE_CREDITO
	res, cmd := mNC.Update(tea.KeyPressMsg{Text: "ctrl+d"})
	if cmd != nil || res.(Model).mode == modeConfirmDeleteFattura {
		t.Error("ctrl+d deve essere ignorato nella tab note di credito")
	}

	// Caso 2: ctrl+d con lista fatture vuota -> ignorato
	mViz := New(nil)
	mViz.sezione = SEZ_FATTURE
	res, cmd = mViz.Update(tea.KeyPressMsg{Text: "ctrl+d"})
	if cmd != nil || res.(Model).mode == modeConfirmDeleteFattura {
		t.Error("ctrl+d deve essere ignorato se non ci sono fatture")
	}
}

// ── Calcolo Fatturato e Rendering View ────────────────────────────────────

// TestFatturatoNetto_NoteCreditoSottratte verifica il calcolo algebrico del totale fatturato.
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

// TestView_StatiInterfaccia controlla la corretta generazione del testo ANSI/UI nelle varie schermate.
func TestView_StatiInterfaccia(t *testing.T) {
	m := New(nil)

	// 1. Caricamento in corso
	m.loading = true
	v := m.View()
	if !strings.Contains(v.Content, "Caricamento in corso...") {
		t.Errorf("vista loading errata: %q", v.Content)
	}

	// 2. Errore restituito dal backend
	m.loading = false
	m.err = errors.New("errore fatale API")
	v = m.View()
	if !strings.Contains(v.Content, "Errore: errore fatale API") {
		t.Errorf("vista errore errata: %q", v.Content)
	}

	// 3. Registro fatture vuoto
	m.err = nil
	v = m.View()
	if !strings.Contains(v.Content, "Nessuna fattura registrata.") {
		t.Errorf("vista registro vuoto errata: %q", v.Content)
	}

	// 4. Registro note di credito vuoto
	m.sezione = SEZ_NOTE_CREDITO
	v = m.View()
	if !strings.Contains(v.Content, "Nessuna nota di credito registrata.") {
		t.Errorf("vista note credito vuota errata: %q", v.Content)
	}

	// 5. Creazione Fattura (senza ordini evasi disponibili)
	m.mode = modeCreateFattura
	v = m.View()
	if !strings.Contains(v.Content, "SELEZIONA ORDINE DA FATTURARE") || !strings.Contains(v.Content, "Nessun ordine evaso disponibile.") {
		t.Errorf("vista creazione fattura errata: %q", v.Content)
	}

	// 6. Creazione Nota Credito (con ordini evasi disponibili)
	m.mode = modeCreateNotaCredito
	m.ordiniEvasi = []api.OrdineUscita{{ID: 1, Numero: "ORD-99", Cliente: "Test Client", Data: time.Now()}}
	v = m.View()
	if !strings.Contains(v.Content, "SELEZIONA ORDINE PER NOTA DI CREDITO") || !strings.Contains(v.Content, "ORD-99") {
		t.Errorf("vista creazione nota credito con ordini errata: %q", v.Content)
	}

	// 7. Modale conferma eliminazione
	m.mode = modeConfirmDeleteFattura
	m.confirmDeleteNome = "FT-2026-01"
	v = m.View()
	if !strings.Contains(v.Content, "Eliminare la fattura \"FT-2026-01\"?") {
		t.Errorf("vista conferma eliminazione errata: %q", v.Content)
	}

	// 8. Visualizzazione messaggio di esito favorevole
	m.mode = modeView
	m.success = "Operazione completata!"
	v = m.View()
	if !strings.Contains(v.Content, "Operazione completata!") {
		t.Errorf("messaggio di successo non visualizzato: %q", v.Content)
	}
}

// TestView_NotaCreditoEImporto verifica il rendering della valuta e del tipo di documento.
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

// TestUpdate_WindowSizeMsg accerta che la ridimensione finestra non lanci comandi pendenti.
func TestUpdate_WindowSizeMsg(t *testing.T) {
	m := New(nil)
	res, cmd := m.Update(tea.WindowSizeMsg{Width: 100, Height: 40})
	if cmd != nil {
		t.Error("WindowSizeMsg non deve restituire comandi")
	}
	_ = res
}

// ── Funzioni Helper e Comandi Asincroni API ───────────────────────────────

// TestHelper_TroncaStringa verifica la formattazione dei testi lunghi con ellipsis.
func TestHelper_TroncaStringa(t *testing.T) {
	if res := troncaStringa("ABC", 5); res != "ABC" {
		t.Errorf("atteso 'ABC', ottenuto %q", res)
	}
	if res := troncaStringa("Azienda Molto Lunga SPA", 10); res != "Azienda..." {
		t.Errorf("atteso 'Azienda...', ottenuto %q", res)
	}
	if res := troncaStringa("Testo", 2); res != "Te" {
		t.Errorf("atteso 'Te' per maxLen <= 3, ottenuto %q", res)
	}
}

// TestLoadFattureCmd_Successo simula il caricamento HTTP iniziale dal server backend mock.
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
		default:
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode([]map[string]any{})
		}
	}))
	defer srv.Close()

	m := New(api.New(srv.URL))

	msg := m.Init()()
	loaded, ok := msg.(FattureLoadedMsg)
	if !ok {
		t.Fatalf("atteso FattureLoadedMsg, ottenuto %T (%v)", msg, msg)
	}
	if len(loaded.Fatture) != 1 || len(loaded.NoteCredito) != 1 {
		t.Fatalf("conteggio documenti caricati errato: fatture=%d, noteCredito=%d", len(loaded.Fatture), len(loaded.NoteCredito))
	}

	fattura := loaded.Fatture[0]
	if fattura.NotaCredito || fattura.Anno != 2026 || fattura.Cliente != "Mario Rossi" {
		t.Errorf("campi fattura non corretti: %+v", fattura)
	}
}

// TestLoadFattureCmd_ErroreServer verifica la gestione delle risposte di errore HTTP 500.
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

// TestCmd_OperazioniAPIAsync verifica il corretto comportamento delle chiamate REST
// asincrone (GET, POST, DELETE).
func TestCmd_OperazioniAPIAsync(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		// Sostituire il controllo strict r.URL.Path == "/api/v1/orders" con HasPrefix
		case r.Method == "GET" && strings.HasPrefix(r.URL.Path, "/api/v1/orders"):
			json.NewEncoder(w).Encode([]map[string]any{{
				"id": 10, "numero": "ORD-10", "cliente": "Cliente Test", "stato": "evaso",
			}})
		case r.Method == "POST" && r.URL.Path == "/api/v1/orders/10/invoice":
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]any{"status": "ok"})
		case r.Method == "POST" && r.URL.Path == "/api/v1/orders/10/credit-note":
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]any{"status": "ok"})
		case r.Method == "DELETE" && r.URL.Path == "/api/v1/orders/invoices/99":
			w.WriteHeader(http.StatusNoContent)
		default:
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode([]map[string]any{})
		}
	}))
	defer srv.Close()

	m := New(api.New(srv.URL))

	// 1. loadOrdiniEvasiCmd
	msgOrdini := m.loadOrdiniEvasiCmd()()
	if loaded, ok := msgOrdini.(ordiniEvasiLoadedMsg); !ok || len(loaded.Ordini) != 1 {
		t.Fatalf("loadOrdiniEvasiCmd fallito: %v", msgOrdini)
	}

	// 2. creaFatturaCmd
	msgCrea := m.creaFatturaCmd(10)()
	if _, ok := msgCrea.(fatturaCreataMsg); !ok {
		t.Fatalf("creaFatturaCmd fallito: %v", msgCrea)
	}

	// 3. creaNotaCreditoCmd
	msgNC := m.creaNotaCreditoCmd(10)()
	if _, ok := msgNC.(notaCreditoCreataMsg); !ok {
		t.Fatalf("creaNotaCreditoCmd fallito: %v", msgNC)
	}

	// 4. eliminaFatturaCmd
	msgDel := m.eliminaFatturaCmd(99)()
	if _, ok := msgDel.(fatturaEliminataMsg); !ok {
		t.Fatalf("eliminaFatturaCmd fallito: %v", msgDel)
	}
}

// ── Test Helpers ──────────────────────────────────────────────────────────

func caricaFattureDiProva(t *testing.T, n int) Model {
	t.Helper()
	fatture := make([]api.Fattura, n)
	for i := range fatture {
		fatture[i] = api.Fattura{
			ID:      i + 1,
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

func caricaFattureDiProvaConNoteCredito(t *testing.T) Model {
	t.Helper()
	m := New(nil)
	res, _ := m.Update(FattureLoadedMsg{
		Fatture: []api.Fattura{
			{ID: 1, Numero: "FT-001", Anno: 2026, Cliente: "Cliente A", Importo: 100},
			{ID: 2, Numero: "FT-002", Anno: 2026, Cliente: "Cliente B", Importo: 200},
		},
		NoteCredito: []api.Fattura{
			{ID: 3, Numero: "NC-001", Anno: 2026, Cliente: "Cliente A", Importo: 50, NotaCredito: true},
		},
	})
	return res.(Model)
}
