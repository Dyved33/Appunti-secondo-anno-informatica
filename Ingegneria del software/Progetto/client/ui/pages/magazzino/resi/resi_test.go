package resi

import (
	"strings"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/api"
)

func TestResi_InitAndUpdate(t *testing.T) {
	client := &api.APIClient{}
	m := New(client)

	cmd := m.Init()
	if cmd == nil {
		t.Error("Init() dovrebbe restituire un Cmd non-nil")
	}

	resiMock := []api.Reso{
		{
			ID: 1, OrdineUscitaID: 20, ProdottoID: 202, Quantita: 2,
			Motivo: "Prodotto difettoso", Data: time.Now(), CreatedByID: 3,
		},
	}

	res, _ := m.Update(ResiLoadedMsg{Resi: resiMock})
	mRes := res.(Model)

	if mRes.loading {
		t.Error("loading dovrebbe essere false dopo caricamento resi")
	}
	if len(mRes.resi) != 1 {
		t.Fatalf("Resi attesi 1, ottenuti %d", len(mRes.resi))
	}

	// Test espansione dettaglio con Enter
	resExpand, _ := mRes.Update(tea.KeyPressMsg{Text: "enter"})
	mExpand := resExpand.(Model)

	if !mExpand.dettagliEspansi[1] {
		t.Error("Il reso #1 dovrebbe essere espanso dopo Enter")
	}

	view := mExpand.View()
	if !strings.Contains(view.Content, "GESTIONE RESI") {
		t.Error("View() dovrebbe contenere il titolo GESTIONE RESI")
	}
	if !strings.Contains(view.Content, "Prodotto difettoso") {
		t.Error("View() dovrebbe mostrare il motivo del reso")
	}
	if !strings.Contains(view.Content, "DETTAGLI PRATICA RESO #1") {
		t.Error("View() espansa dovrebbe contenere l'intestazione della sotto-griglia del reso")
	}
}

func TestResi_FooterModificaStato(t *testing.T) {
	client := &api.APIClient{}
	m := New(client)

	resiMock := []api.Reso{
		{
			ID: 1, OrdineUscitaID: 20, ProdottoID: 202, Quantita: 2,
			Motivo: "Difettoso", Data: time.Now(), CreatedByID: 3,
		},
	}

	res, _ := m.Update(ResiLoadedMsg{Resi: resiMock})
	mRes := res.(Model)

	view := mRes.View()
	if !strings.Contains(view.Content, "[ctrl+e] Modifica Stato") {
		t.Error("View() dovrebbe mostrare la scorciatoia [ctrl+e] Modifica Stato nel footer")
	}
}

func TestResi_AperturaModificaStato(t *testing.T) {
	client := &api.APIClient{}
	m := New(client)

	resiMock := []api.Reso{
		{
			ID: 1, OrdineUscitaID: 20, ProdottoID: 202, Quantita: 2,
			Motivo: "Prodotto difettoso", Data: time.Now(), CreatedByID: 3,
			Stato: "da_verificare",
		},
	}

	res, _ := m.Update(ResiLoadedMsg{Resi: resiMock})
	mRes := res.(Model)

	// Apri il form di modifica stato con ctrl+e
	resEdit, _ := mRes.Update(tea.KeyPressMsg{Text: "ctrl+e"})
	mEdit := resEdit.(Model)

	if mEdit.modo != modoModificaStato {
		t.Error("Dopo ctrl+e il modo dovrebbe essere modoModificaStato")
	}

	if mEdit.idxStato != 0 {
		t.Errorf("idxStato dovrebbe essere 0 (da_verificare), ottenuto %d", mEdit.idxStato)
	}

	view := mEdit.View()
	if !strings.Contains(view.Content, "MODIFICA STATO RESO") {
		t.Error("View() dovrebbe mostrare il titolo MODIFICA STATO RESO")
	}
	if !strings.Contains(view.Content, "da_verificare") {
		t.Error("View() dovrebbe mostrare lo stato attuale da_verificare")
	}
}

func TestResi_CambioStatoConLeftRight(t *testing.T) {
	client := &api.APIClient{}
	m := New(client)

	resiMock := []api.Reso{
		{
			ID: 1, OrdineUscitaID: 20, ProdottoID: 202, Quantita: 2,
			Motivo: "Test", Data: time.Now(), CreatedByID: 3,
			Stato: "da_verificare",
		},
	}

	res, _ := m.Update(ResiLoadedMsg{Resi: resiMock})
	mRes := res.(Model)

	// Apri modifica
	resEdit, _ := mRes.Update(tea.KeyPressMsg{Text: "ctrl+e"})
	mEdit := resEdit.(Model)

	// Cambia stato con right
	resRight, _ := mEdit.Update(tea.KeyPressMsg{Text: "right"})
	mRight := resRight.(Model)

	if mRight.idxStato != 1 {
		t.Errorf("Dopo right idxStato dovrebbe essere 1 (accettato), ottenuto %d", mRight.idxStato)
	}

	// Cambia stato di nuovo con right
	resRight2, _ := mRight.Update(tea.KeyPressMsg{Text: "right"})
	mRight2 := resRight2.(Model)

	if mRight2.idxStato != 2 {
		t.Errorf("Dopo secondo right idxStato dovrebbe essere 2 (rifiutato), ottenuto %d", mRight2.idxStato)
	}

	// Torna indietro con left
	resLeft, _ := mRight2.Update(tea.KeyPressMsg{Text: "left"})
	mLeft := resLeft.(Model)

	if mLeft.idxStato != 1 {
		t.Errorf("Dopo left idxStato dovrebbe essere 1 (accettato), ottenuto %d", mLeft.idxStato)
	}
}

func TestResi_EscModificaStato(t *testing.T) {
	client := &api.APIClient{}
	m := New(client)

	resiMock := []api.Reso{
		{
			ID: 1, OrdineUscitaID: 20, ProdottoID: 202, Quantita: 2,
			Motivo: "Test", Data: time.Now(), CreatedByID: 3,
		},
	}

	res, _ := m.Update(ResiLoadedMsg{Resi: resiMock})
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
