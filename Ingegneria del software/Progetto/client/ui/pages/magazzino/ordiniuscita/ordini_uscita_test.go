package ordiniuscita

import (
	"strings"
	"testing"
	"time"

	"github.com/tommasogarofani/ingsw_25-26_group-11/client/api"

	tea "charm.land/bubbletea/v2"
)

func TestOrdiniUscita_InitAndUpdate(t *testing.T) {
	client := &api.APIClient{}
	m := New(client)

	cmd := m.Init()
	if cmd == nil {
		t.Error("Init() dovrebbe restituire un Cmd non-nil")
	}

	ddtID := 88
	ordiniMock := []api.OrdineUscita{
		{
			ID: 20, Numero: "USC-2026-0005", Cliente: "Cliente Beta SRL",
			Data: time.Now(), Stato: api.StatoEvaso, DDTID: &ddtID, Note: "Spedito via corriere",
			Righe: []api.RigaOrdineUscita{
				{
					ID: 2, OrdineUscitaID: 20, ProdottoID: 202, Quantita: 10, PrezzoUnitario: 99.00,
					Prodotto: &api.Prodotto{Codice: "PRD-202", Nome: "Trapano Elettrico"},
				},
			},
		},
	}

	res, _ := m.Update(OrdiniUscitaLoadedMsg{Ordini: ordiniMock})
	mRes := res.(Model)

	if mRes.loading {
		t.Error("loading dovrebbe essere false dopo caricamento ordini")
	}
	if len(mRes.ordini) != 1 {
		t.Fatalf("Ordini attesi 1, ottenuti %d", len(mRes.ordini))
	}

	// Test espansione dettaglio con Enter
	resExpand, _ := mRes.Update(tea.KeyPressMsg{Text: "enter"})
	mExpand := resExpand.(Model)

	if !mExpand.dettagliEspansi[20] {
		t.Error("L'ordine #20 dovrebbe essere espanso dopo Enter")
	}

	view := mExpand.View()
	if !strings.Contains(view.Content, "USC-2026-0005") {
		t.Error("View() dovrebbe mostrare il numero ordine USC-2026-0005")
	}
	if !strings.Contains(view.Content, "DDT #88") {
		t.Error("View() dovrebbe mostrare il riferimento DDT #88")
	}
	if !strings.Contains(view.Content, "Trapano Elettrico") {
		t.Error("View() espansa dovrebbe mostrare il dettaglio prodotto Trapano Elettrico")
	}
}

func TestOrdiniUscita_FooterModifica(t *testing.T) {
	client := &api.APIClient{}
	m := New(client)

	ordiniMock := []api.OrdineUscita{
		{
			ID: 1, Numero: "USC-2026-0001", Cliente: "Cliente Alpha",
			Data: time.Now(), Stato: api.StatoConfermato,
			Righe: []api.RigaOrdineUscita{
				{ID: 1, OrdineUscitaID: 1, ProdottoID: 10, Quantita: 5, PrezzoUnitario: 20.00},
			},
		},
	}

	res, _ := m.Update(OrdiniUscitaLoadedMsg{Ordini: ordiniMock})
	mRes := res.(Model)

	view := mRes.View()
	if !strings.Contains(view.Content, "[ctrl+e] Modifica") {
		t.Error("View() dovrebbe mostrare la scorciatoia [ctrl+e] Modifica nel footer")
	}
}

func TestOrdiniUscita_AperturaModifica(t *testing.T) {
	client := &api.APIClient{}
	m := New(client)

	ordiniMock := []api.OrdineUscita{
		{
			ID: 1, Numero: "USC-2026-0001", Cliente: "Cliente Alpha",
			Data: time.Now(), Stato: api.StatoConfermato, Note: "Note test",
			Righe: []api.RigaOrdineUscita{
				{ID: 1, OrdineUscitaID: 1, ProdottoID: 10, Quantita: 5, PrezzoUnitario: 20.00},
			},
		},
	}

	res, _ := m.Update(OrdiniUscitaLoadedMsg{Ordini: ordiniMock})
	mRes := res.(Model)

	// Apri il form di modifica con ctrl+e
	resEdit, _ := mRes.Update(tea.KeyPressMsg{Text: "ctrl+e"})
	mEdit := resEdit.(Model)

	if mEdit.modo != modoModifica {
		t.Error("Dopo ctrl+e il modo dovrebbe essere modoModifica")
	}

	view := mEdit.View()
	if !strings.Contains(view.Content, "MODIFICA ORDINE USCITA") {
		t.Error("View() dovrebbe mostrare il titolo MODIFICA ORDINE USCITA")
	}
	if !strings.Contains(view.Content, "USC-2026-0001") {
		t.Error("View() dovrebbe mostrare il numero ordine pre-popolato")
	}
	if !strings.Contains(view.Content, "Cliente Alpha") {
		t.Error("View() dovrebbe mostrare il cliente pre-popolato")
	}
	if !strings.Contains(view.Content, "Note test") {
		t.Error("View() dovrebbe mostrare le note pre-popolate")
	}
}

func TestOrdiniUscita_EscModifica(t *testing.T) {
	client := &api.APIClient{}
	m := New(client)

	ordiniMock := []api.OrdineUscita{
		{
			ID: 1, Numero: "USC-2026-0001", Cliente: "Cliente Alpha",
			Data: time.Now(), Stato: api.StatoConfermato,
			Righe: []api.RigaOrdineUscita{
				{ID: 1, OrdineUscitaID: 1, ProdottoID: 10, Quantita: 5, PrezzoUnitario: 20.00},
			},
		},
	}

	res, _ := m.Update(OrdiniUscitaLoadedMsg{Ordini: ordiniMock})
	mRes := res.(Model)

	// Apri modifica
	resEdit, _ := mRes.Update(tea.KeyPressMsg{Text: "ctrl+e"})
	mEdit := resEdit.(Model)

	// Chiudi con esc
	resEsc, _ := mEdit.Update(tea.KeyPressMsg{Text: "esc"})
	mEsc := resEsc.(Model)

	if mEsc.modo != modoLista {
		t.Error("Dopo esc il modo dovrebbe tornare a modoLista")
	}
}

// Test navigazione tra i campi del form e salvataggio
func TestOrdiniUscita_NavigazioneFormESalvataggio(t *testing.T) {
	client := &api.APIClient{}
	m := New(client)

	ordiniMock := []api.OrdineUscita{
		{
			ID: 1, Numero: "USC-001", Cliente: "Cliente Test",
			Data: time.Now(), Stato: api.StatoConfermato,
			Righe: []api.RigaOrdineUscita{{ProdottoID: 5, Quantita: 2, PrezzoUnitario: 10.0}},
		},
	}

	res, _ := m.Update(OrdiniUscitaLoadedMsg{Ordini: ordiniMock})
	mRes := res.(Model)

	// Entra in modifica (focusNumero)
	resEdit, _ := mRes.Update(tea.KeyPressMsg{Text: "ctrl+e"})
	mEdit := resEdit.(Model)

	// Premi Tab per passare a focusCliente, poi focusNote, poi focusStato
	for i := 0; i < 3; i++ {
		resEdit, _ = mEdit.Update(tea.KeyPressMsg{Text: "tab"})
		mEdit = resEdit.(Model)
	}

	if mEdit.focusMod != focusStato {
		t.Errorf("Dopo 3 Tab il focus dovrebbe essere su focusStato, ottenuto: %v", mEdit.focusMod)
	}

	// Cambia stato da 'confermato' a 'evaso' usando la freccia destra
	resRight, _ := mEdit.Update(tea.KeyPressMsg{Text: "right"})
	mRight := resRight.(Model)
	if mRight.idxStato != 1 {
		t.Errorf("Indice stato atteso 1 (evaso), ottenuto: %d", mRight.idxStato)
	}

	// Test esecuzione del comando al ritorno del messaggio OrdineAggiornatoMsg
	resUpdated, cmd := mRight.Update(OrdineAggiornatoMsg{})
	mUpdated := resUpdated.(Model)

	if mUpdated.modo != modoLista {
		t.Error("Dopo OrdineAggiornatoMsg il modo deve tornare a modoLista")
	}
	if cmd == nil {
		t.Error("Dopo OrdineAggiornatoMsg deve essere restituito il Cmd di ricaricamento")
	}
}

// Test gestione della vista DDT
func TestOrdiniUscita_GestioneVistaDDT(t *testing.T) {
	client := &api.APIClient{}
	m := New(client)

	ddtMock := api.DDT{
		ID:           99,
		Numero:       "DDT-2026-99",
		Data:         time.Now(),
		Mittente:     "Magazzino Centrale",
		Destinatario: "Cliente SRL",
		PathPDF:      "/tmp/ddt_99.pdf",
	}

	// Simulazione ricezione DDTCreatedMsg
	res, _ := m.Update(DDTCreatedMsg{DDT: ddtMock})
	mDDT := res.(Model)

	if mDDT.view != DDT {
		t.Error("La vista corrente dovrebbe essere DDT")
	}

	view := mDDT.View()
	if !strings.Contains(view.Content, "DOCUMENTO DI TRASPORTO (DDT #99)") {
		t.Error("View() dovrebbe renderizzare l'intestazione del DDT #99")
	}
	if !strings.Contains(view.Content, "DDT-2026-99") {
		t.Error("View() dovrebbe mostrare il numero del DDT")
	}

	// Uscita dalla vista DDT con ESC
	resEsc, _ := mDDT.Update(tea.KeyPressMsg{Text: "esc"})
	mEsc := resEsc.(Model)

	if mEsc.view != MAIN {
		t.Error("Premendo esc dalla vista DDT la vista deve tornare a MAIN")
	}
}
