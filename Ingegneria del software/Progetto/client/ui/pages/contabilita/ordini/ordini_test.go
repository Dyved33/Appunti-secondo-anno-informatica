// Package ordini raccoglie i test d'unità della sotto-pagina Ordini & Listini.
package ordini

import (
	"strings"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/api"
)

// nuoviDatiMock costruisce un messaggio OrdiniLoadedMsg con dati di prova per
// le tre sezioni (entrata, uscita, resi).
func nuoviDatiMock() OrdiniLoadedMsg {
	return OrdiniLoadedMsg{
		Entrata: []api.OrdineEntrata{
			{ID: 1, Numero: "ENT-001", Fornitore: "Fornitore Alfa", Data: time.Date(2026, 8, 24, 0, 0, 0, 0, time.UTC),
				Righe: []api.RigaOrdineEntrata{{ProdottoID: 7, Quantita: 100, PrezzoUnitario: 1.5}}},
		},
		Uscita: []api.OrdineUscita{
			{ID: 2, Numero: "USC-001", Cliente: "Mario Rossi", Data: time.Date(2026, 8, 24, 0, 0, 0, 0, time.UTC),
				Stato: api.StatoConfermato,
				Righe: []api.RigaOrdineUscita{{ProdottoID: 7, Quantita: 50, PrezzoUnitario: 2.0}}},
			{ID: 3, Numero: "USC-002", Cliente: "Luigi Verdi", Data: time.Date(2026, 8, 25, 0, 0, 0, 0, time.UTC),
				Stato: api.StatoEvaso},
		},
		Resi: []api.Reso{
			{ID: 10, OrdineUscitaID: 2, ProdottoID: 7, Quantita: 5, Motivo: "prodotto danneggiato",
				Stato: api.StatoAccettato, Data: time.Date(2026, 8, 26, 0, 0, 0, 0, time.UTC)},
		},
	}
}

// caricaDati restituisce un modello già popolato con i dati mock.
func caricaDati(t *testing.T) Model {
	t.Helper()
	m := New(nil)
	res, _ := m.Update(nuoviDatiMock())
	caricato, ok := res.(Model)
	if !ok {
		t.Fatalf("atteso Model da Update(OrdiniLoadedMsg), ottenuto %T", res)
	}
	if caricato.loading {
		t.Fatal("atteso loading=false dopo il caricamento")
	}
	return caricato
}

func TestNew_StatoIniziale(t *testing.T) {
	m := New(nil)

	if !m.loading {
		t.Error("atteso loading=true all'avvio")
	}
	if m.sezione != ORDINI_ENTRATA {
		t.Errorf("attesa sezione iniziale ORDINI_ENTRATA, ottenuta %d", m.sezione)
	}
}

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

func TestUpdate_ErroreMsg(t *testing.T) {
	m := New(nil)

	res, _ := m.Update(ErroreMsg{Err: errMock{}})
	m = res.(Model)

	if m.loading {
		t.Error("atteso loading=false dopo ErroreMsg")
	}
	if v := m.View().Content; !strings.Contains(v, "Errore:") {
		t.Errorf("atteso messaggio di errore nella vista, ottenuto %q", v)
	}
}

type errMock struct{}

func (errMock) Error() string { return "connessione rifiutata" }

func TestUpdate_CicloSezioniConTab(t *testing.T) {
	m := caricaDati(t)

	passi := []struct {
		attesa      sezioneAttiva
		descrizione string
	}{
		{ORDINI_USCITA, "primo tab"},
		{RESI, "secondo tab"},
		{ORDINI_ENTRATA, "terzo tab (ritorno)"},
	}

	for _, p := range passi {
		res, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
		m = res.(Model)
		if m.sezione != p.attesa {
			t.Errorf("%s: attesa sezione %d, ottenuta %d", p.descrizione, p.attesa, m.sezione)
		}
	}
}

func TestUpdate_NavigazioneCursoreIndipendente(t *testing.T) {
	m := caricaDati(t)

	// Nella sezione entrata (un solo ordine) il cursore non può scendere.
	res, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	m = res.(Model)
	if m.cursoreEntrata != 0 {
		t.Fatalf("cursore entrate spostato oltre l'ultimo ordine: %d", m.cursoreEntrata)
	}

	// Passo all'uscita e scendo: cursoreUscita avanza senza toccare cursoreEntrata.
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
}

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

func TestView_ResiMostraMotivoEspanso(t *testing.T) {
	m := caricaDati(t)

	// Due tab: entrata -> uscita -> resi.
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
	if !m.espansiResi[10] {
		t.Error("atteso reso #10 segnato come espanso")
	}
}

func TestParseRigheInput(t *testing.T) {
	// Caso valido singolo
	righe, err := parseRigheInput("7:100:1.50")
	if err != nil {
		t.Fatalf("errore inatteso su input valido: %v", err)
	}
	if len(righe) != 1 || righe[0].ProdottoID != 7 || righe[0].Quantita != 100 || righe[0].PrezzoUnitario != 1.50 {
		t.Errorf("riga parsata non corretta: %+v", righe)
	}

	// Caso valido multiplo
	righe, err = parseRigheInput(" 1:10:5.00 , 2:20:10.50 ")
	if err != nil {
		t.Fatalf("errore inatteso su input multiplo valido: %v", err)
	}
	if len(righe) != 2 {
		t.Fatalf("attese 2 righe, ottenute %d", len(righe))
	}

	// Casi non validi
	invalidInputs := []string{
		"",
		"   ",
		"1:10",
		"abc:10:5.0",
		"1:abc:5.0",
		"1:10:abc",
		"-1:10:5.0",
		"1:-5:5.0",
		"1:10:-2.0",
	}
	for _, in := range invalidInputs {
		if _, err := parseRigheInput(in); err == nil {
			t.Errorf("atteso errore per input non valido %q, ottenuto nil", in)
		}
	}
}

func TestCreazioneOrdineEntrata_AvvioECancellazione(t *testing.T) {
	m := caricaDati(t)

	// Ctrl+N nella tab entrata attiva il form di creazione entrata
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

	// Esc annulla e resetta il form
	res, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyEscape})
	m = res.(Model)

	if m.creazioneEntrata || m.creazioneUscita || m.creazioneReso {
		t.Errorf("atteso reset del form dopo Esc, ottenuto entrata=%v uscita=%v reso=%v",
			m.creazioneEntrata, m.creazioneUscita, m.creazioneReso)
	}
}

func TestCreazioneOrdineUscita_AvvioENavigazione(t *testing.T) {
	m := caricaDati(t)

	// Passa a tab Uscita
	res, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	m = res.(Model)

	// Ctrl+N nella tab uscita attiva il form di creazione uscita
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
	if !strings.Contains(v, "Cliente:") {
		t.Errorf("campo Cliente non presente nella vista: %s", v)
	}

	// Navigazione campi con Tab / Down
	if m.focusCampo != FocusControparte {
		t.Errorf("atteso focus iniziale su FocusControparte, ottenuto %d", m.focusCampo)
	}

	res, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	m = res.(Model)
	if m.focusCampo != FocusNumero {
		t.Errorf("atteso focus su FocusNumero dopo tab, ottenuto %d", m.focusCampo)
	}

	res, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	m = res.(Model)
	if m.focusCampo != FocusNote {
		t.Errorf("atteso focus su FocusNote dopo tab, ottenuto %d", m.focusCampo)
	}

	res, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	m = res.(Model)
	if m.focusCampo != FocusRighe {
		t.Errorf("atteso focus su FocusRighe dopo tab, ottenuto %d", m.focusCampo)
	}

	// Navigazione all'indietro con Up
	res, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyUp})
	m = res.(Model)
	if m.focusCampo != FocusNote {
		t.Errorf("atteso focus su FocusNote dopo up, ottenuto %d", m.focusCampo)
	}
}

func TestCreazioneReso_AvvioENavigazione(t *testing.T) {
	m := caricaDati(t)

	// Passa alla tab Resi (due Tab: Entrata -> Uscita -> Resi)
	res, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	m = res.(Model)
	res, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	m = res.(Model)

	// Ctrl+N nella tab Resi attiva il form di creazione reso
	res, _ = m.Update(tea.KeyPressMsg{Code: 'n', Mod: tea.ModCtrl})
	m = res.(Model)

	if !m.creazioneReso || m.creazioneEntrata || m.creazioneUscita {
		t.Errorf("attesa creazioneReso=true, ottenuto entrata=%v uscita=%v reso=%v",
			m.creazioneEntrata, m.creazioneUscita, m.creazioneReso)
	}

	v := m.View().Content
	if !strings.Contains(v, "CREAZIONE NUOVO RESO") {
		t.Errorf("titolo form reso non presente nella vista: %s", v)
	}
	if !strings.Contains(v, "ID Ordine Uscita:") {
		t.Errorf("campo ID Ordine Uscita non presente nella vista: %s", v)
	}

	// Focus iniziale su FocusResoOrdineID
	if m.focusCampo != FocusResoOrdineID {
		t.Errorf("atteso focus iniziale su FocusResoOrdineID, ottenuto %d", m.focusCampo)
	}

	// Navigazione campi con Tab
	res, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	m = res.(Model)
	if m.focusCampo != FocusResoProdottoID {
		t.Errorf("atteso focus su FocusResoProdottoID dopo tab, ottenuto %d", m.focusCampo)
	}

	res, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	m = res.(Model)
	if m.focusCampo != FocusResoQuantita {
		t.Errorf("atteso focus su FocusResoQuantita dopo tab, ottenuto %d", m.focusCampo)
	}

	res, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	m = res.(Model)
	if m.focusCampo != FocusResoMotivo {
		t.Errorf("atteso focus su FocusResoMotivo dopo tab, ottenuto %d", m.focusCampo)
	}

	// Wrap-around del focus da FocusResoMotivo a FocusResoOrdineID
	res, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	m = res.(Model)
	if m.focusCampo != FocusResoOrdineID {
		t.Errorf("atteso wrap-around a FocusResoOrdineID dopo tab, ottenuto %d", m.focusCampo)
	}

	// Navigazione all'indietro con Up (da FocusResoOrdineID a FocusResoMotivo)
	res, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyUp})
	m = res.(Model)
	if m.focusCampo != FocusResoMotivo {
		t.Errorf("atteso focus su FocusResoMotivo dopo up da ID ordine, ottenuto %d", m.focusCampo)
	}

	// Esc annulla e resetta il form reso
	res, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyEscape})
	m = res.(Model)

	if m.creazioneReso {
		t.Errorf("atteso reset del form reso dopo Esc, ottenuto creazioneReso=true")
	}
}

func TestCreazioneOrdine_ValidazioneCampi(t *testing.T) {
	m := caricaDati(t)

	// Avvia form entrata
	res, _ := m.Update(tea.KeyPressMsg{Code: 'n', Mod: tea.ModCtrl})
	m = res.(Model)

	// Salva a campi vuoti -> errore validazione controparte
	res, _ = m.Update(tea.KeyPressMsg{Code: 's', Mod: tea.ModCtrl})
	m = res.(Model)
	if m.err == nil || !strings.Contains(m.err.Error(), "fornitore") {
		t.Errorf("atteso errore validazione fornitore, ottenuto: %v", m.err)
	}

	// Imposta controparte
	m.formControparte.SetValue("Fornitore Test")
	res, _ = m.Update(tea.KeyPressMsg{Code: 's', Mod: tea.ModCtrl})
	m = res.(Model)
	if m.err == nil || !strings.Contains(m.err.Error(), "numero") {
		t.Errorf("atteso errore validazione numero, ottenuto: %v", m.err)
	}

	// Imposta numero d'ordine
	m.formNumero.SetValue("ORD-999")
	res, _ = m.Update(tea.KeyPressMsg{Code: 's', Mod: tea.ModCtrl})
	m = res.(Model)
	if m.err == nil || !strings.Contains(m.err.Error(), "riga") {
		t.Errorf("atteso errore validazione righe, ottenuto: %v", m.err)
	}
}

func TestCreazioneReso_ValidazioneCampi(t *testing.T) {
	m := caricaDati(t)

	// Passa a tab Resi e avvia form
	res, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	m = res.(Model)
	res, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	m = res.(Model)
	res, _ = m.Update(tea.KeyPressMsg{Code: 'n', Mod: tea.ModCtrl})
	m = res.(Model)

	// Tenta salvataggio con campi vuoti -> errore ID ordine
	res, _ = m.Update(tea.KeyPressMsg{Code: 's', Mod: tea.ModCtrl})
	m = res.(Model)
	if m.err == nil || !strings.Contains(m.err.Error(), "ID ordine") {
		t.Errorf("atteso errore validazione ID ordine, ottenuto: %v", m.err)
	}

	// Imposta ID ordine valido
	m.formResoOrdineID.SetValue("2")
	res, _ = m.Update(tea.KeyPressMsg{Code: 's', Mod: tea.ModCtrl})
	m = res.(Model)
	if m.err == nil || !strings.Contains(m.err.Error(), "ID prodotto") {
		t.Errorf("atteso errore validazione ID prodotto, ottenuto: %v", m.err)
	}

	// Imposta ID prodotto valido
	m.formResoProdottoID.SetValue("7")
	res, _ = m.Update(tea.KeyPressMsg{Code: 's', Mod: tea.ModCtrl})
	m = res.(Model)
	if m.err == nil || !strings.Contains(m.err.Error(), "quantità") {
		t.Errorf("atteso errore validazione quantità, ottenuto: %v", m.err)
	}

	// Imposta quantità non valida (<= 0)
	m.formResoQuantita.SetValue("0")
	res, _ = m.Update(tea.KeyPressMsg{Code: 's', Mod: tea.ModCtrl})
	m = res.(Model)
	if m.err == nil || !strings.Contains(m.err.Error(), "quantità") {
		t.Errorf("atteso errore validazione quantità > 0, ottenuto: %v", m.err)
	}
}

func TestCreazioneOrdine_SuccessoReset(t *testing.T) {
	m := caricaDati(t)

	// Avvia form entrata
	res, _ := m.Update(tea.KeyPressMsg{Code: 'n', Mod: tea.ModCtrl})
	m = res.(Model)

	// Notifica di ordine creato con successo
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

func TestCreazioneReso_SuccessoReset(t *testing.T) {
	m := caricaDati(t)

	// Avvia form reso
	res, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	m = res.(Model)
	res, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	m = res.(Model)
	res, _ = m.Update(tea.KeyPressMsg{Code: 'n', Mod: tea.ModCtrl})
	m = res.(Model)

	// Notifica di reso creato con successo
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

func TestEliminazione_FlussoConfermaEAnnullamento(t *testing.T) {
	m := caricaDati(t)

	// Premendo ctrl+d si attiva la modalità di conferma eliminazione
	res, _ := m.Update(tea.KeyPressMsg{Code: 'd', Mod: tea.ModCtrl})
	m = res.(Model)

	if !m.confermaEliminazione {
		t.Errorf("attesa confermaEliminazione=true dopo ctrl+d")
	}
	if m.idDaEliminare != 1 {
		t.Errorf("atteso idDaEliminare=1, ottenuto %d", m.idDaEliminare)
	}

	// Annullamento con tasto 'n'
	res, _ = m.Update(tea.KeyPressMsg{Code: 'n'})
	m = res.(Model)

	if m.confermaEliminazione {
		t.Errorf("atteso annullamento della conferma dopo la pressione di 'n'")
	}
}

func TestEliminazione_ConfermaEGenerazioneMsg(t *testing.T) {
	m := caricaDati(t)

	// Avvia eliminazione su Ordine Entrata #1
	res, _ := m.Update(tea.KeyPressMsg{Code: 'd', Mod: tea.ModCtrl})
	m = res.(Model)

	// Conferma con Enter
	res, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	m = res.(Model)

	if m.confermaEliminazione {
		t.Errorf("atteso reset di confermaEliminazione dopo l'invio")
	}
	if !m.loading {
		t.Errorf("atteso loading=true durante l'elaborazione dell'eliminazione")
	}
	if cmd == nil {
		t.Errorf("atteso tea.Cmd non nil per l'invio dell'eliminazione al backend")
	}
}

func TestUpdate_MessaggiEliminazioneCompletata(t *testing.T) {
	m := caricaDati(t)

	// Simulazione ricezione messaggio di eliminazione avvenuta
	res, cmd := m.Update(OrdineEntrataEliminatoMsg{})
	m = res.(Model)

	if !m.loading {
		t.Errorf("atteso loading=true per ricaricare gli ordini dopo l'eliminazione")
	}
	if cmd == nil {
		t.Errorf("atteso tea.Cmd di ricaricamento non nil")
	}
}

func TestUpdate_PaginazioneERicarica(t *testing.T) {
	m := caricaDati(t)

	// Test tasto r
	res, cmd := m.Update(tea.KeyPressMsg{Code: 'r'})
	m = res.(Model)
	if !m.loading {
		t.Errorf("atteso loading=true dopo aver premuto 'r'")
	}
	if cmd == nil {
		t.Errorf("atteso comando di ricarica dati non nil")
	}

	// Test PageDown / PageUp
	res, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyPgDown})
	m = res.(Model)
	if m.cursoreEntrata < 0 || m.cursoreEntrata >= len(m.ordiniEntrata) {
		t.Errorf("cursore fuori dai limiti dopo PgDown: %d", m.cursoreEntrata)
	}
}
