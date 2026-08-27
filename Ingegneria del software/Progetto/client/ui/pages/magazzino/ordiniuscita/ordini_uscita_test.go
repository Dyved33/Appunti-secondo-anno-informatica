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
