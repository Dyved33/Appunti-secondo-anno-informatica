package ui

import (
	"testing"

	style "github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/styles"

	tea "charm.land/bubbletea/v2"
)

// TestNuovaFinestra verifica che il costruttore principale NuovaFinestra
// crei una finestra impostando la pagina corrente iniziale a "menuprincipale".
func TestNuovaFinestra(t *testing.T) {
	f := NuovaFinestra()
	if f.paginaCorrente != "menuprincipale" {
		t.Errorf("expected initial page to be menuprincipale, got %s", f.paginaCorrente)
	}
}

// TestFinestraPrincipale_Init verifica che il metodo Init() del frame principale
// non esegua alcun comando asincrono iniziale (ritorna nil).
func TestFinestraPrincipale_Init(t *testing.T) {
	f := NuovaFinestra()
	cmd := f.Init()
	if cmd != nil {
		t.Error("expected Init to return nil cmd")
	}
}

// TestFinestraPrincipale_Update_WindowSizeMsg verifica che al resize della finestra
// le nuove dimensioni vengano propagate correttamente allo stato interno e al menu principale.
func TestFinestraPrincipale_Update_WindowSizeMsg(t *testing.T) {
	f := NuovaFinestra()
	msg := tea.WindowSizeMsg{Width: 100, Height: 50}
	res, cmd := f.Update(msg)
	if cmd != nil {
		t.Error("expected nil cmd for WindowSizeMsg")
	}
	fRes := res.(FinestraPrincipale)

	// Controlla che le nuove dimensioni della finestra siano registrate
	if style.ReadWidth(fRes.dimensioneFinestra) != 100 {
		t.Errorf("expected width 100, got %d", style.ReadWidth(fRes.dimensioneFinestra))
	}
	if style.ReadHeight(fRes.dimensioneFinestra) != 50 {
		t.Errorf("expected height 50, got %d", style.ReadHeight(fRes.dimensioneFinestra))
	}
}

// TestFinestraPrincipale_Update_KeyPressMsg_MainMenu verifica l'inoltro dei
// messaggi di tastiera alla sotto-schermata attiva (in questo caso il menu principale).
func TestFinestraPrincipale_Update_KeyPressMsg_MainMenu(t *testing.T) {
	f := NuovaFinestra()
	f.paginaCorrente = "menuprincipale"

	// Simula la pressione del tasto 'D'
	msgD := tea.KeyPressMsg{Text: "D"}
	res, _ := f.Update(msgD)
	fRes := res.(FinestraPrincipale)

	// Esegue il rendering grafico per verificare che la schermata si aggiorni senza panico
	view := fRes.View()
	if view.Content == "" {
		t.Error("expected non-empty View content")
	}
}

// TestFinestraPrincipale_View verifica che il layout complessivo dell'applicazione
// (Header + Logo + Menu + Spazi vuoti + Footer) venga generato correttamente.
func TestFinestraPrincipale_View(t *testing.T) {
	f := NuovaFinestra()
	f.dimensioneFinestra.ChangeSize(80, 24)
	view := f.View()
	if view.Content == "" {
		t.Error("expected non-empty View content")
	}
}
