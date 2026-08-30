// Package ordini raccoglie i test d'unità della sotto-pagina Ordini & Listini.
// I test coprono il ciclo di vita del modello Bubble Tea, le interazioni utente (navigazione,
// scorrimento, form di creazione e invio), la validazione dei dati e la resa grafica (View).
package ordini

import (
	"strings"
	"testing"
	"time"

	"charm.land/lipgloss/v2"
	tea "charm.land/bubbletea/v2"
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/api"
	style "github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/styles"
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/tabnav"
)

// errMock è una struttura ausiliaria per simulare errori di rete o di backend durante i test.
type errMock struct{}

func (errMock) Error() string { return "connessione rifiutata dal server" }

// nuoviDatiMock costruisce un messaggio OrdiniLoadedMsg popolato con un set completo di dati
// di prova per ciascuna delle tre sezioni (ordini di entrata, ordini di uscita e resi).
func nuoviDatiMock() OrdiniLoadedMsg {
	return OrdiniLoadedMsg{
		Entrata: []api.OrdineEntrata{
			{
				ID:        1,
				Numero:    "ENT-001",
				Fornitore: "Fornitore Alfa",
				Data:      time.Date(2026, 8, 24, 0, 0, 0, 0, time.UTC),
				Righe: []api.RigaOrdineEntrata{
					{ProdottoID: 7, Quantita: 100, PrezzoUnitario: 1.5},
				},
				Note: "Nota ordine entrata 1",
			},
		},
		Uscita: []api.OrdineUscita{
			{
				ID:      2,
				Numero:  "USC-001",
				Cliente: "Mario Rossi",
				Data:    time.Date(2026, 8, 24, 0, 0, 0, 0, time.UTC),
				Stato:   api.StatoConfermato,
				Righe:   []api.RigaOrdineUscita{{ProdottoID: 7, Quantita: 50, PrezzoUnitario: 2.0}},
				Note:    "Consegna urgente",
			},
			{
				ID:      3,
				Numero:  "USC-002",
				Cliente: "Luigi Verdi",
				Data:    time.Date(2026, 8, 25, 0, 0, 0, 0, time.UTC),
				Stato:   api.StatoEvaso,
			},
		},
		Resi: []api.Reso{
			{
				ID:             10,
				OrdineUscitaID: 2,
				ProdottoID:     7,
				Quantita:       5,
				Motivo:         "prodotto danneggiato",
				Stato:          api.StatoAccettato,
				Data:           time.Date(2026, 8, 26, 0, 0, 0, 0, time.UTC),
			},
		},
	}
}

// caricaDati inizializza un nuovo Model e vi applica un messaggio di caricamento dati mock.
// Restituisce un Model pronto all'uso con lo stato `loading` impostato a false.
func caricaDati(t *testing.T) Model {
	t.Helper()
	m := New(nil)
	res, _ := m.Update(nuoviDatiMock())
	caricato, ok := res.(Model)
	if !ok {
		t.Fatalf("atteso Model da Update(OrdiniLoadedMsg), ottenuto %T", res)
	}
	if caricato.loading {
		t.Fatal("atteso loading=false dopo il caricamento dei dati")
	}
	return caricato
}

// TestNew_StatoIniziale verifica che la funzione di costruzione New restituisca un modello
// correttamente inizializzato con caricamento attivo e la prima sezione selezionata.
func TestNew_StatoIniziale(t *testing.T) {
	m := New(nil)

	if !m.loading {
		t.Error("atteso loading=true all'avvio del modello")
	}
	if m.sezione != ORDINI_ENTRATA {
		t.Errorf("attesa sezione iniziale ORDINI_ENTRATA (%d), ottenuta %d", ORDINI_ENTRATA, m.sezione)
	}
	if m.espansiEntrata == nil || m.espansiUscita == nil || m.espansiResi == nil {
		t.Error("attesa inizializzazione delle mappe di espansione dettagli")
	}
}

// TestInit verifica che il metodo Init() restituisca un tea.Cmd non nil per il caricamento iniziale dei dati.
func TestInit(t *testing.T) {
	m := New(nil)
	cmd := m.Init()
	if cmd == nil {
		t.Error("atteso tea.Cmd non nil ritornato da Init()")
	}
}

// TestUpdate_OrdiniLoadedMsg verifica la corretta elaborazione del messaggio OrdiniLoadedMsg,
// accertando la popolazione degli array e la disattivazione dello stato di loading.
func TestUpdate_OrdiniLoadedMsg(t *testing.T) {
	m := caricaDati(t)

	if len(m.ordiniEntrata) != 1 || len(m.ordiniUscita) != 2 || len(m.resi) != 1 {
		t.Fatalf("dati non distribuiti correttamente: entrata=%d uscita=%d resi=%d",
			len(m.ordiniEntrata), len(m.ordiniUscita), len(m.resi))
	}

	v := m.View().Content
	if !strings.Contains(v, "GESTIONE ORDINI") {
		t.Errorf("atteso titolo GESTIONE ORDINI, ottenuto %q", v)
	}
	if !strings.Contains(v, "Fornitore Alfa") {
		t.Errorf("attesa tabella entrata popolata con 'Fornitore Alfa', ottenuta %q", v)
	}
}

// TestUpdate_ErroreMsg verifica la corretta gestione del messaggio di errore, assicurandosi
// che interrompa lo stato di caricamento e mostri il testo dell'errore nella vista.
func TestUpdate_ErroreMsg(t *testing.T) {
	m := New(nil)

	res, _ := m.Update(ErroreMsg{Err: errMock{}})
	m = res.(Model)

	if m.loading {
		t.Error("atteso loading=false dopo la ricezione di ErroreMsg")
	}
	if m.err == nil || m.err.Error() != "connessione rifiutata dal server" {
		t.Errorf("atteso errore memorizzato nel modello, ottenuto: %v", m.err)
	}
	if v := m.View().Content; !strings.Contains(v, "Errore: connessione rifiutata dal server") {
		t.Errorf("atteso messaggio di errore formattato nella vista, ottenuto %q", v)
	}
}

// TestUpdate_CicloSezioniConTab verifica la rotazione ciclica delle sezioni/tab tramite il tasto Tab.
func TestUpdate_CicloSezioniConTab(t *testing.T) {
	m := caricaDati(t)

	passi := []struct {
		attesa      sezioneAttiva
		descrizione string
	}{
		{ORDINI_USCITA, "primo tab (Ordini Uscita)"},
		{RESI, "secondo tab (Resi)"},
		{ORDINI_ENTRATA, "terzo tab (ritorno ad Ordini Entrata)"},
	}

	for _, p := range passi {
		res, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
		m = res.(Model)
		if m.sezione != p.attesa {
			t.Errorf("%s: attesa sezione %d, ottenuta %d", p.descrizione, p.attesa, m.sezione)
		}
	}
}

// TestUpdate_NavigazioneCursoreIndipendente verifica che i cursori di selezione rimangano
// indipendenti tra le varie schede e non sforino i limiti della lista.
func TestUpdate_NavigazioneCursoreIndipendente(t *testing.T) {
	m := caricaDati(t)

	// Nella sezione entrata (un solo ordine) il cursore verso il basso non deve superare l'indice 0.
	res, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	m = res.(Model)
	if m.cursoreEntrata != 0 {
		t.Fatalf("cursore entrate spostato oltre l'ultimo ordine: %d", m.cursoreEntrata)
	}

	// Passaggio al tab Uscita e scorrimento in basso.
	res, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	m = res.(Model)
	res, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	m = res.(Model)

	if m.cursoreUscita != 1 {
		t.Errorf("atteso cursore uscita su 1, ottenuto %d", m.cursoreUscita)
	}
	if m.cursoreEntrata != 0 {
		t.Errorf("atteso cursore entrata invariato a 0, ottenuto %d", m.cursoreEntrata)
	}

	// Tentativo di scorrere oltre l'ultimo elemento in Uscita
	res, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	m = res.(Model)
	if m.cursoreUscita != 1 {
		t.Errorf("cursore uscita sfora il limite massimo: %d", m.cursoreUscita)
	}
}

// TestToggleEspansione_TutteLeSezioni controlla la funzionalità di espansione/compressione dei dettagli
// per tutte e tre le tipologie di tabella (Entrata, Uscita, Resi).
func TestToggleEspansione_TutteLeSezioni(t *testing.T) {
	m := caricaDati(t)

	// 1. Espansione ordine in entrata (ID 1)
	res, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	m = res.(Model)
	if !m.espansiEntrata[1] {
		t.Error("atteso ordine entrata #1 espanso")
	}

	// 2. Passaggio a Uscita ed espansione ordine in uscita (ID 2)
	res, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	m = res.(Model)
	res, _ = m.Update(tea.KeyPressMsg{Code: tea.KeySpace}) // Proviamo anche il tasto spazio
	m = res.(Model)
	if !m.espansiUscita[2] {
		t.Error("atteso ordine uscita #2 espanso tramite tasto Space")
	}

	// 3. Passaggio a Resi ed espansione reso (ID 10)
	res, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	m = res.(Model)
	res, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	m = res.(Model)
	if !m.espansiResi[10] {
		t.Error("atteso reso #10 espanso")
	}

	// Richiudiamo il reso
	res, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	m = res.(Model)
	if m.espansiResi[10] {
		t.Error("atteso reso #10 me-compress / chiuso")
	}
}

// TestView_UscitaMostraStatiEClienti garantisce che la vista della tabella degli ordini in uscita
// visualizzi correttamente i dettagli dei clienti, l'intestazione e i badge di stato.
func TestView_UscitaMostraStatiEClienti(t *testing.T) {
	m := caricaDati(t)

	res, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	m = res.(Model)

	v := m.View().Content
	for _, atteso := range []string{"Mario Rossi", "Luigi Verdi", "STATO"} {
		if !strings.Contains(v, atteso) {
			t.Errorf("vista uscita senza %q:\n%s", atteso, v)
		}
	}
}

// TestView_ResiMostraMotivoEspanso controlla che i motivi dei resi non siano visibili inizialmente,
// ma appaiano soltanto dopo aver premuto il tasto di espansione.
func TestView_ResiMostraMotivoEspanso(t *testing.T) {
	m := caricaDati(t)

	// Spostamento nella tab Resi.
	res, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	m = res.(Model)
	res, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	m = res.(Model)

	if v := m.View().Content; strings.Contains(v, "prodotto danneggiato") {
		t.Error("motivo del reso visibile prima dell'espansione del dettaglio")
	}

	res, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	m = res.(Model)

	v := m.View().Content
	if !strings.Contains(v, "prodotto danneggiato") {
		t.Errorf("atteso dettaglio reso espanso con motivo, ottenuto:\n%s", v)
	}
}

// TestParseRigheInput testa approfonditamente la funzione di parsing delle righe d'ordine,
// inclusi casi validi con spaziature ed errori di formattazione o valori negativi/zero.
func TestParseRigheInput(t *testing.T) {
	// Caso valido singolo
	righe, err := parseRigheInput("7:100:1.50")
	if err != nil {
		t.Fatalf("errore inatteso su input valido: %v", err)
	}
	if len(righe) != 1 || righe[0].ProdottoID != 7 || righe[0].Quantita != 100 || righe[0].PrezzoUnitario != 1.50 {
		t.Errorf("riga parsata non corretta: %+v", righe)
	}

	// Caso valido multiplo con spaziature irregolari
	righe, err = parseRigheInput(" 1:10:5.00 , 2:20:10.50 ")
	if err != nil {
		t.Fatalf("errore inatteso su input multiplo valido: %v", err)
	}
	if len(righe) != 2 {
		t.Fatalf("attese 2 righe, ottenute %d", len(righe))
	}

	// Matrice di test per casi di errata formattazione
	invalidInputs := []struct {
		input       string
		description string
	}{
		{"", "stringa vuota"},
		{"   ", "soli spazi"},
		{"1:10", "manca il prezzo"},
		{"abc:10:5.0", "ID prodotto non numerico"},
		{"1:abc:5.0", "quantità non numerica"},
		{"1:10:abc", "prezzo non numerico"},
		{"-1:10:5.0", "ID prodotto negativo"},
		{"0:10:5.0", "ID prodotto pari a zero"},
		{"1:-5:5.0", "quantità negativa"},
		{"1:0:5.0", "quantità pari a zero"},
		{"1:10:-2.0", "prezzo negativo"},
		{"1:10:5.0:extra", "troppi campi separati da due punti"},
	}

	for _, tt := range invalidInputs {
		t.Run(tt.description, func(t *testing.T) {
			if _, err := parseRigheInput(tt.input); err == nil {
				t.Errorf("atteso errore per input non valido %q (%s), ottenuto nil", tt.input, tt.description)
			}
		})
	}
}

// TestCreazioneOrdineEntrata_AvvioECancellazione assicura che il form di creazione ordine in entrata
// si apra correttamente con Ctrl+N e venga completamente resettato premendo Esc.
func TestCreazioneOrdineEntrata_AvvioECancellazione(t *testing.T) {
	m := caricaDati(t)

	res, _ := m.Update(tea.KeyPressMsg{Code: 'n', Mod: tea.ModCtrl})
	m = res.(Model)

	if !m.creazioneEntrata || m.creazioneUscita || m.creazioneReso {
		t.Errorf("attesa creazioneEntrata=true, ottenuto entrata=%v uscita=%v reso=%v",
			m.creazioneEntrata, m.creazioneUscita, m.creazioneReso)
	}

	v := m.View().Content
	if !strings.Contains(v, "CREAZIONE ORDINE IN ENTRATA") {
		t.Errorf("titolo form entrata non presente nella vista: %s", v)
	}
	if !strings.Contains(v, "Fornitore:") {
		t.Errorf("campo Fornitore non presente nella vista: %s", v)
	}

	// Esc annulla ed esegue il reset completo del form.
	res, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyEscape})
	m = res.(Model)

	if m.creazioneEntrata || m.creazioneUscita || m.creazioneReso {
		t.Errorf("atteso reset del form dopo Esc, ottenuto entrata=%v uscita=%v reso=%v",
			m.creazioneEntrata, m.creazioneUscita, m.creazioneReso)
	}
}

// TestCreazioneOrdineUscita_AvvioENavigazione verifica l'avvio del form ordini in uscita,
// la navigazione tra i campi tramite Tab/Down e il movimento a ritroso con Up/Shift+Tab.
func TestCreazioneOrdineUscita_AvvioENavigazione(t *testing.T) {
	m := caricaDati(t)

	// Passa a tab Uscita
	res, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	m = res.(Model)

	// Ctrl+N nella tab uscita attiva il form relativo
	res, _ = m.Update(tea.KeyPressMsg{Code: 'n', Mod: tea.ModCtrl})
	m = res.(Model)

	if !m.creazioneUscita || m.creazioneEntrata || m.creazioneReso {
		t.Errorf("attesa creazioneUscita=true, ottenuto entrata=%v uscita=%v reso=%v",
			m.creazioneEntrata, m.creazioneUscita, m.creazioneReso)
	}

	v := m.View().Content
	if !strings.Contains(v, "CREAZIONE ORDINE IN USCITA") {
		t.Errorf("titolo form uscita non presente nella vista: %s", v)
	}

	// Verifiche avanzamento focus
	if m.focusCampo != FocusControparte {
		t.Errorf("atteso focus iniziale su FocusControparte, ottenuto %d", m.focusCampo)
	}

	res, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	m = res.(Model)
	if m.focusCampo != FocusNumero {
		t.Errorf("atteso focus su FocusNumero dopo tab, ottenuto %d", m.focusCampo)
	}

	res, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	m = res.(Model)
	if m.focusCampo != FocusNote {
		t.Errorf("atteso focus su FocusNote dopo down, ottenuto %d", m.focusCampo)
	}

	// Navigazione a ritroso con Shift+Tab e Up
	res, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyUp})
	m = res.(Model)
	if m.focusCampo != FocusNumero {
		t.Errorf("atteso focus su FocusNumero dopo up, ottenuto %d", m.focusCampo)
	}

	res, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyTab, Mod: tea.ModShift})
	m = res.(Model)
	if m.focusCampo != FocusControparte {
		t.Errorf("atteso focus su FocusControparte dopo shift+tab, ottenuto %d", m.focusCampo)
	}
}

// TestCreazioneReso_AvvioENavigazione verifica il ciclo completo di navigazione e wrap-around
// del focus all'interno del form di creazione reso.
func TestCreazioneReso_AvvioENavigazione(t *testing.T) {
	m := caricaDati(t)

	// Selezione tab Resi
	res, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	m = res.(Model)
	res, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	m = res.(Model)

	// Apertura form di inserimento reso
	res, _ = m.Update(tea.KeyPressMsg{Code: 'n', Mod: tea.ModCtrl})
	m = res.(Model)

	if !m.creazioneReso {
		t.Fatalf("attesa creazioneReso=true")
	}

	// Scorrimento fino all'ultimo campo
	campi := []focusForm{FocusResoProdottoID, FocusResoQuantita, FocusResoMotivo}
	for _, c := range campi {
		res, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
		m = res.(Model)
		if m.focusCampo != c {
			t.Errorf("atteso focus su %d, ottenuto %d", c, m.focusCampo)
		}
	}

	// Wrap-around del focus dal campo Motivo a ID Ordine
	res, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	m = res.(Model)
	if m.focusCampo != FocusResoOrdineID {
		t.Errorf("atteso wrap-around a FocusResoOrdineID, ottenuto %d", m.focusCampo)
	}

	// Wrap-around a ritroso con Up
	res, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyUp})
	m = res.(Model)
	if m.focusCampo != FocusResoMotivo {
		t.Errorf("atteso wrap-around a ritroso su FocusResoMotivo, ottenuto %d", m.focusCampo)
	}
}

// TestCreazioneOrdine_ValidazioneCampi verifica progressivamente i messaggi di errore restituiti
// quando si tenta di salvare un ordine con i campi vuoti o incompleti.
func TestCreazioneOrdine_ValidazioneCampi(t *testing.T) {
	m := caricaDati(t)

	// Avvio form ordine entrata
	res, _ := m.Update(tea.KeyPressMsg{Code: 'n', Mod: tea.ModCtrl})
	m = res.(Model)

	// 1. Errore per fornitore assente
	res, _ = m.Update(tea.KeyPressMsg{Code: 's', Mod: tea.ModCtrl})
	m = res.(Model)
	if m.err == nil || !strings.Contains(m.err.Error(), "fornitore") {
		t.Errorf("atteso errore validazione fornitore, ottenuto: %v", m.err)
	}

	// 2. Errore per numero ordine assente
	m.formControparte.SetValue("Fornitore Test")
	res, _ = m.Update(tea.KeyPressMsg{Code: 's', Mod: tea.ModCtrl})
	m = res.(Model)
	if m.err == nil || !strings.Contains(m.err.Error(), "numero") {
		t.Errorf("atteso errore validazione numero, ottenuto: %v", m.err)
	}

	// 3. Errore per righe prodotti assenti
	m.formNumero.SetValue("ORD-999")
	res, _ = m.Update(tea.KeyPressMsg{Code: 's', Mod: tea.ModCtrl})
	m = res.(Model)
	if m.err == nil || !strings.Contains(m.err.Error(), "riga") {
		t.Errorf("atteso errore validazione righe, ottenuto: %v", m.err)
	}
}

// TestCreazioneReso_ValidazioneCampi verifica i controlli sulle componenti numeriche del form dei resi.
func TestCreazioneReso_ValidazioneCampi(t *testing.T) {
	m := caricaDati(t)

	// Passaggio alla tab Resi e apertura form
	res, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	m = res.(Model)
	res, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	m = res.(Model)
	res, _ = m.Update(tea.KeyPressMsg{Code: 'n', Mod: tea.ModCtrl})
	m = res.(Model)

	// 1. ID Ordine non valido o mancante
	res, _ = m.Update(tea.KeyPressMsg{Code: 's', Mod: tea.ModCtrl})
	m = res.(Model)
	if m.err == nil || !strings.Contains(m.err.Error(), "ID ordine") {
		t.Errorf("atteso errore validazione ID ordine, ottenuto: %v", m.err)
	}

	// 2. ID Prodotto non valido o mancante
	m.formResoOrdineID.SetValue("2")
	res, _ = m.Update(tea.KeyPressMsg{Code: 's', Mod: tea.ModCtrl})
	m = res.(Model)
	if m.err == nil || !strings.Contains(m.err.Error(), "ID prodotto") {
		t.Errorf("atteso errore validazione ID prodotto, ottenuto: %v", m.err)
	}

	// 3. Quantità non valida o mancante
	m.formResoProdottoID.SetValue("7")
	res, _ = m.Update(tea.KeyPressMsg{Code: 's', Mod: tea.ModCtrl})
	m = res.(Model)
	if m.err == nil || !strings.Contains(m.err.Error(), "quantità") {
		t.Errorf("atteso errore validazione quantità, ottenuto: %v", m.err)
	}

	// 4. Quantità non maggiore di zero
	m.formResoQuantita.SetValue("0")
	res, _ = m.Update(tea.KeyPressMsg{Code: 's', Mod: tea.ModCtrl})
	m = res.(Model)
	if m.err == nil || !strings.Contains(m.err.Error(), "quantità") {
		t.Errorf("atteso errore quantita > 0, ottenuto: %v", m.err)
	}
}

// TestSalvaOrdine_Successo controlla il salvataggio corretto degli ordini (Entrata e Uscita)
// quando tutti i campi compilati risultano validi.
func TestSalvaOrdine_Successo(t *testing.T) {
	m := caricaDati(t)

	// Form Entrata valido
	res, _ := m.Update(tea.KeyPressMsg{Code: 'n', Mod: tea.ModCtrl})
	m = res.(Model)
	m.formControparte.SetValue("Fornitore Spa")
	m.formNumero.SetValue("ENT-100")
	m.formRighe.SetValue("1:10:2.50")

	mRes, cmd := m.salvaOrdine()
	mod := mRes.(Model)
	if mod.err != nil {
		t.Errorf("errore inatteso durante il salvataggio ordine entrata: %v", mod.err)
	}
	if !mod.loading {
		t.Error("atteso loading=true durante l'invio al backend")
	}
	if cmd == nil {
		t.Error("atteso tea.Cmd non nil per creaOrdineEntrataCmd")
	}

	// Form Uscita valido
	m2 := caricaDati(t)
	res, _ = m2.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	m2 = res.(Model)
	res, _ = m2.Update(tea.KeyPressMsg{Code: 'n', Mod: tea.ModCtrl})
	m2 = res.(Model)
	m2.formControparte.SetValue("Cliente Srl")
	m2.formNumero.SetValue("USC-100")
	m2.formRighe.SetValue("2:5:10.00")

	mRes2, cmd2 := m2.salvaOrdine()
	mod2 := mRes2.(Model)
	if mod2.err != nil {
		t.Errorf("errore inatteso durante il salvataggio ordine uscita: %v", mod2.err)
	}
	if cmd2 == nil {
		t.Error("atteso tea.Cmd non nil per creaOrdineUscitaCmd")
	}
}

// TestSalvaReso_Successo controlla l'esito positivo del metodo salvaReso con dati corretti.
func TestSalvaReso_Successo(t *testing.T) {
	m := caricaDati(t)
	m.avviaCreazioneReso()
	m.formResoOrdineID.SetValue("10")
	m.formResoProdottoID.SetValue("3")
	m.formResoQuantita.SetValue("2")
	m.formResoMotivo.SetValue("Difettoso")

	mRes, cmd := m.salvaReso()
	mod := mRes.(Model)
	if mod.err != nil {
		t.Errorf("errore inatteso durante salvaReso: %v", mod.err)
	}
	if !mod.loading {
		t.Error("atteso loading=true dopo l'invio del reso")
	}
	if cmd == nil {
		t.Error("atteso tea.Cmd non nil ritornato da salvaReso")
	}
}

// TestForm_InvioConEnter controlla che premendo Enter nei form di creazione si avanzi tra i campi
// e che sul campo finale si scateni l'azione di salvataggio.
func TestForm_InvioConEnter(t *testing.T) {
	m := caricaDati(t)

	// Apertura form ordini
	res, _ := m.Update(tea.KeyPressMsg{Code: 'n', Mod: tea.ModCtrl})
	m = res.(Model)

	// Inseriamo dati incompleti per verificare l'invio finale
	m.formControparte.SetValue("Test")
	res, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyEnter}) // Da Controparte a Numero
	m = res.(Model)
	if m.focusCampo != FocusNumero {
		t.Errorf("atteso focus su FocusNumero, ottenuto %d", m.focusCampo)
	}

	res, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyEnter}) // Da Numero a Note
	m = res.(Model)
	if m.focusCampo != FocusNote {
		t.Errorf("atteso focus su FocusNote, ottenuto %d", m.focusCampo)
	}

	res, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyEnter}) // Da Note a Righe
	m = res.(Model)
	if m.focusCampo != FocusRighe {
		t.Errorf("atteso focus su FocusRighe, ottenuto %d", m.focusCampo)
	}

	// Su FocusRighe, Enter scatena il tentativo di salvataggio (che produrrà errore per numero vuoto)
	res, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	m = res.(Model)
	if m.err == nil {
		t.Error("atteso errore di validazione dal salvataggio con Enter su ultimo campo")
	}
}

// TestCreazioneOrdine_SuccessoReset garantisce la chiusura automatica del form e l'avvio del ricaricamento
// alla ricezione dei messaggi di avvenuta creazione.
func TestCreazioneOrdine_SuccessoReset(t *testing.T) {
	m := caricaDati(t)

	res, _ := m.Update(tea.KeyPressMsg{Code: 'n', Mod: tea.ModCtrl})
	m = res.(Model)

	res, cmd := m.Update(OrdineEntrataCreatoMsg{})
	m = res.(Model)

	if m.creazioneEntrata || m.creazioneUscita || m.creazioneReso {
		t.Errorf("atteso form chiuso dopo OrdineEntrataCreatoMsg")
	}
	if !m.loading {
		t.Errorf("atteso loading=true per ricaricare gli ordini")
	}
	if cmd == nil {
		t.Errorf("atteso cmd di ricaricamento non nil")
	}
}

// TestCreazioneReso_SuccessoReset garantisce il ripristino dello stato a seguito del messaggio ResoCreatoMsg.
func TestCreazioneReso_SuccessoReset(t *testing.T) {
	m := caricaDati(t)

	res, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	m = res.(Model)
	res, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	m = res.(Model)
	res, _ = m.Update(tea.KeyPressMsg{Code: 'n', Mod: tea.ModCtrl})
	m = res.(Model)

	res, cmd := m.Update(ResoCreatoMsg{})
	m = res.(Model)

	if m.creazioneReso {
		t.Errorf("atteso form reso chiuso dopo ResoCreatoMsg")
	}
	if !m.loading {
		t.Errorf("atteso loading=true per ricaricare i dati")
	}
	if cmd == nil {
		t.Errorf("atteso cmd di ricaricamento non nil")
	}
}

// TestEliminazione_FlussoCompletoETutteLeSezioni controlla la procedura di eliminazione guidata
// per tutte le schede (Entrata, Uscita, Resi), verificando conferme e annullamenti tramite Esc/N.
func TestEliminazione_FlussoCompletoETutteLeSezioni(t *testing.T) {
	// 1. Eliminazione Ordine Entrata annullata con tasto 'n'
	m := caricaDati(t)
	res, _ := m.Update(tea.KeyPressMsg{Code: 'd', Mod: tea.ModCtrl})
	m = res.(Model)
	if !m.confermaEliminazione || m.idDaEliminare != 1 {
		t.Errorf("attesa conferma eliminazione attivata su ID 1")
	}

	res, _ = m.Update(tea.KeyPressMsg{Code: 'n'})
	m = res.(Model)
	if m.confermaEliminazione {
		t.Error("atteso annullamento della conferma con 'n'")
	}

	// 2. Eliminazione Ordine Uscita annullata con 'esc'
	res, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	m = res.(Model)
	res, _ = m.Update(tea.KeyPressMsg{Code: 'd', Mod: tea.ModCtrl})
	m = res.(Model)
	if !m.confermaEliminazione || m.idDaEliminare != 2 {
		t.Errorf("attesa conferma eliminazione attivata su ordine uscita #2")
	}

	res, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyEscape})
	m = res.(Model)
	if m.confermaEliminazione {
		t.Error("atteso annullamento della conferma con 'esc'")
	}

	// 3. Eliminazione Reso confermata con Enter ('y')
	res, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	m = res.(Model) // Passa alla tab Resi
	res, _ = m.Update(tea.KeyPressMsg{Code: 'd', Mod: tea.ModCtrl})
	m = res.(Model)
	if !m.confermaEliminazione || m.idDaEliminare != 10 {
		t.Errorf("attesa conferma eliminazione attivata su reso #10")
	}

	res, cmd := m.Update(tea.KeyPressMsg{Code: 'y'})
	m = res.(Model)
	if m.confermaEliminazione {
		t.Error("attesa uscita da stato di conferma dopo 'y'")
	}
	if !m.loading {
		t.Error("atteso loading=true durante l'elaborazione dell'eliminazione")
	}
	if cmd == nil {
		t.Error("atteso comando di eliminazione reso non nil")
	}
}

// TestUpdate_MessaggiEliminazioneCompletata verifica la reazione del modello ai messaggi
// di avvenuta eliminazione per tutte le entità.
func TestUpdate_MessaggiEliminazioneCompletata(t *testing.T) {
	m := caricaDati(t)

	messaggi := []tea.Msg{
		OrdineEntrataEliminatoMsg{},
		OrdineUscitaEliminatoMsg{},
		ResoEliminatoMsg{},
	}

	for _, msg := range messaggi {
		res, cmd := m.Update(msg)
		mod := res.(Model)
		if !mod.loading {
			t.Errorf("atteso loading=true dopo la ricezione di %T", msg)
		}
		if cmd == nil {
			t.Errorf("atteso tea.Cmd di ricaricamento dati per %T", msg)
		}
	}
}

// TestUpdate_PaginazioneERicarica controlla l'effetto del tasto 'r' di ricarica
// e i comandi di paginazione rapida PgDown / PgUp.
func TestUpdate_PaginazioneERicarica(t *testing.T) {
	m := caricaDati(t)

	// Tasto 'r' per la ricarica manuale dei dati
	res, cmd := m.Update(tea.KeyPressMsg{Code: 'r'})
	m = res.(Model)
	if !m.loading {
		t.Errorf("atteso loading=true dopo aver premuto 'r'")
	}
	if cmd == nil {
		t.Errorf("atteso comando di ricarica dati non nil")
	}

	// Test scorrimento pagina
	res, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyPgDown})
	m = res.(Model)
	if m.cursoreEntrata < 0 || m.cursoreEntrata >= len(m.ordiniEntrata) {
		t.Errorf("cursore fuori dai limiti dopo PgDown: %d", m.cursoreEntrata)
	}

	res, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyPgUp})
	m = res.(Model)
	if m.cursoreEntrata != 0 {
		t.Errorf("atteso cursore su 0 dopo PgUp, ottenuto %d", m.cursoreEntrata)
	}
}

// TestScorrimentoEWidgetRose controlla lo scorrimento tramite rotellina del mouse e la conformità delle dimensioni
// del layout complessivo integrato con le macro-componenti dell'applicazione.
func TestScorrimentoEWidgetRose(t *testing.T) {
	var entrate []api.OrdineEntrata
	for i := 1; i <= 10; i++ {
		entrate = append(entrate, api.OrdineEntrata{
			ID:        i,
			Numero:    "ENT-00" + string(rune('0'+i)),
			Fornitore: "Fornitore X",
			Data:      time.Date(2026, 8, 24, 0, 0, 0, 0, time.UTC),
			Righe:     []api.RigaOrdineEntrata{{ProdottoID: 7, Quantita: 100, PrezzoUnitario: 1.5}},
		})
	}

	m := New(nil)
	res, _ := m.Update(OrdiniLoadedMsg{Entrata: entrate, Uscita: nil, Resi: nil})
	m = res.(Model)
	m.winH = 40
	m.winW = 120

	in := m.cursoreEntrata
	res, _ = m.Update(tea.MouseWheelMsg{Button: tea.MouseWheelDown, X: 1, Y: 1})
	m = res.(Model)
	if m.cursoreEntrata <= in {
		t.Errorf("atteso avanzamento del cursore con la rotellina (da %d a %d)", in, m.cursoreEntrata)
	}
	if !strings.Contains(m.View().Content, "Scorri") {
		t.Error("atteso suggerimento di scorrimento nel footer")
	}

	// Verifica che la pagina formattata rispetti l'altezza dello schermo disponibile.
	for _, winH := range []int{40, 50} {
		upd, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: winH})
		m = upd.(Model)
		view := m.View().Content
		if !strings.Contains(view, "rotellina") {
			t.Errorf("winH=%d: le scorciatoie del footer non sono presenti nella vista", winH)
		}
		tabs := lipgloss.JoinHorizontal(lipgloss.Top,
			style.ActiveTab.Render("ORDINI"), "  ", style.TabGap.Render("CATALOGO"))
		composed := tabnav.PlacePage(120, winH,
			style.FinestraStile.Render(tabs+"\n\n"+view))
		if strings.Count(composed, "\n")+1 > winH {
			t.Errorf("winH=%d: la pagina composta occupa %d righe (footer fuori schermo)",
				winH, strings.Count(composed, "\n")+1)
		}
	}
}

// TestTroncaStringa verifica il corretto troncamento delle stringhe lunghe con omissione "...".
func TestTroncaStringa(t *testing.T) {
	str := "Questa è una stringa molto lunga"

	if res := troncaStringa(str, 50); res != str {
		t.Errorf("attesa stringa inalterata, ottenuta %q", res)
	}

	if res := troncaStringa(str, 10); res != "Questa ..." {
		t.Errorf("atteso troncamento 'Questa ...', ottenuto %q", res)
	}

	if res := troncaStringa("Abcde", 3); res != "Abc" {
		t.Errorf("attesa stringa troncata senza punti per maxLen<=3, ottenuta %q", res)
	}
}

// TestRenderingStato verifica che la stilizzazione visiva delle opzioni di stato dell'ordine
// e del reso copra tutti i rami previsti dalle enumeration API.
func TestRenderingStato(t *testing.T) {
	statiOrdine := []api.StatoOrdine{
		api.StatoConfermato,
		api.StatoEvaso,
		api.StatoAnnullato,
		api.StatoOrdine("Sconosciuto"),
	}

	for _, s := range statiOrdine {
		res := renderingStatoOrdine(s, string(s))
		if res == "" {
			t.Errorf("attesa stringa formattata non vuota per stato ordine %s", s)
		}
	}

	statiReso := []api.StatoReso{
		api.StatoInAttesa,
		api.StatoAccettato,
		api.StatoRifiutato,
		api.StatoReso("Sconosciuto"),
	}

	for _, s := range statiReso {
		res := renderingStatoReso(s, string(s))
		if res == "" {
			t.Errorf("attesa stringa formattata non vuota per stato reso %s", s)
		}
	}
}

