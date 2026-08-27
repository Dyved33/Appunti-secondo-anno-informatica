package app

import (
	"testing"

	"github.com/tommasogarofani/ingsw_25-26_group-11/client/api"
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/pages/login"
	style "github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/styles"

	tea "charm.land/bubbletea/v2"
)

// helper per creare una FinestraPrincipale pronta all'uso nei test.
func nuovaFinestraTest() mainLayout {
	client := &api.APIClient{}
	f := New(client)
	// Assegna dimensioni ragionevoli per evitare altezze negative nel View.
	f.WindowSize.ChangeSize(120, 40)
	f.WindowSize.Border.SetSpacing(2, 2, 2, 2)
	return f
}

// ---------------------------------------------------------------------------
// Costruttore
// ---------------------------------------------------------------------------

// TestNew verifica che il costruttore New() restituisca una FinestraPrincipale
// con la pagina iniziale impostata a LOGIN (non a MENUPRINCIPALE).
func TestNew(t *testing.T) {
	client := &api.APIClient{}
	f := New(client)

	if f.currentPage != LOGIN {
		t.Errorf("pagina corrente attesa LOGIN (%d), ottenuta %d", LOGIN, f.currentPage)
	}
}

// TestNew_DimensioneInizialeZero verifica che la dimensione iniziale della
// finestra sia (0,0) poiché nessun WindowSizeMsg è stato ancora ricevuto.
func TestNew_DimensioneInizialeZero(t *testing.T) {
	client := &api.APIClient{}
	f := New(client)

	if w := style.ReadWidth(f.WindowSize); w != 0 {
		t.Errorf("larghezza iniziale attesa 0, ottenuta %d", w)
	}
	if h := style.ReadHeight(f.WindowSize); h != 0 {
		t.Errorf("altezza iniziale attesa 0, ottenuta %d", h)
	}
}

// ---------------------------------------------------------------------------
// Init
// ---------------------------------------------------------------------------

// TestInit verifica che Init() restituisca un Cmd non-nil, dato che delega
// all'Init() del sotto-modello Login che mette il focus sul primo input.
func TestInit(t *testing.T) {
	f := nuovaFinestraTest()
	cmd := f.Init()

	if cmd == nil {
		t.Error("Init() dovrebbe restituire un Cmd non-nil (focus sul primo input di login)")
	}
}

// ---------------------------------------------------------------------------
// Update – WindowSizeMsg
// ---------------------------------------------------------------------------

// TestUpdate_WindowSizeMsg verifica che le dimensioni della finestra vengano
// aggiornate correttamente alla ricezione di un WindowSizeMsg.
func TestUpdate_WindowSizeMsg(t *testing.T) {
	f := nuovaFinestraTest()

	msg := tea.WindowSizeMsg{Width: 200, Height: 60}
	res, _ := f.Update(msg)
	fRes := res.(mainLayout)

	// Le dimensioni lette tengono conto del border spacing (2+2 su ogni asse).
	wantW := 200 - 2 - 2
	wantH := 60 - 2 - 2
	gotW, gotH := style.GetInnerSpace(fRes.WindowSize)
	if gotW != wantW {
		t.Errorf("larghezza attesa %d, ottenuta %d", wantW, gotW)
	}
	if gotH != wantH {
		t.Errorf("altezza attesa %d, ottenuta %d", wantH, gotH)
	}
}

// TestUpdate_WindowSizeMsg_PaginaNonCambia verifica che il resize della
// finestra non cambi la pagina corrente.
func TestUpdate_WindowSizeMsg_PaginaNonCambia(t *testing.T) {
	f := nuovaFinestraTest()

	msg := tea.WindowSizeMsg{Width: 100, Height: 50}
	res, _ := f.Update(msg)
	fRes := res.(mainLayout)

	if fRes.currentPage != LOGIN {
		t.Errorf("pagina corrente attesa LOGIN, ottenuta %d", fRes.currentPage)
	}
}

// ---------------------------------------------------------------------------
// Update – LoginSubmitMsg (transizione di pagina)
// ---------------------------------------------------------------------------

// TestUpdate_LoginSubmitMsg verifica che alla ricezione di un LoginSubmitMsg
// la pagina corrente passi da LOGIN al menu corrispondente il ruolo dell'utente.
func TestUpdate_LoginSubmitMsg(t *testing.T) {
	f := nuovaFinestraTest()

	if f.currentPage != LOGIN {
		t.Fatalf("prerequisito fallito: pagina iniziale attesa LOGIN, ottenuta %d", f.currentPage)
	}

	res, cmd := f.Update(login.LoginSubmitMsg{})
	fRes := res.(mainLayout)

	switch fRes.UserRole {
	case MAGAZZINIERE:
		if fRes.currentPage != MAGMENU {
			t.Errorf("pagina corrente attesa MAGMENU (%d), ottenuta %d", MAGMENU, fRes.currentPage)
		}
	case CONTABILE:
		if fRes.currentPage != CONTMENU {
			t.Errorf("pagina corrente attesa CONTMENU (%d), ottenuta %d", CONTMENU, fRes.currentPage)
		}
	case ADMIN:
		if fRes.currentPage != ADMINMENU {
			t.Errorf("pagina corrente attesa ADMINMENU (%d), ottenuta %d", ADMINMENU, fRes.currentPage)
		}
	default:
		t.Errorf("ruolo utente sconosciuto: %v", f.UserRole)
	}

	// Il login deve restituire il Cmd di Init del sotto-menu attivo
	if cmd == nil {
		t.Errorf("LoginSubmitMsg ci si aspettava un Cmd da Init(), ottenuto nil")
	}
}

// ---------------------------------------------------------------------------
// Update – propagazione messaggi alla schermata Login
// ---------------------------------------------------------------------------

// TestUpdate_PropagaMessaggiALogin verifica che i messaggi di tastiera
// vengano inoltrati al sotto-modello Login quando la pagina corrente è LOGIN.
func TestUpdate_PropagaMessaggiALogin(t *testing.T) {
	f := nuovaFinestraTest()
	f.currentPage = LOGIN

	// Simula la pressione di un tasto qualsiasi.
	msg := tea.KeyPressMsg{Text: "a"}
	res, _ := f.Update(msg)
	fRes := res.(mainLayout)

	// Verifica che il modello restituito sia ancora nella pagina LOGIN.
	if fRes.currentPage != LOGIN {
		t.Errorf("pagina corrente attesa LOGIN, ottenuta %d", fRes.currentPage)
	}
}

// ---------------------------------------------------------------------------
// Update – propagazione messaggi alla schermata MenuPrincipale
// ---------------------------------------------------------------------------

// TestUpdate_PropagaMessaggiAMenuPrincipale verifica che i messaggi
// di tastiera vengano inoltrati al sotto-modello di menu specifico quando la
// pagina corrente è ...MENU.
func TestUpdate_PropagaMessaggiAMenu(t *testing.T) {
	f := nuovaFinestraTest()
	f.currentPage = MAGMENU

	msg := tea.KeyPressMsg{Text: "j"}
	res, _ := f.Update(msg)
	fRes := res.(mainLayout)

	if fRes.currentPage != MAGMENU {
		t.Errorf("pagina corrente attesa MAGMENU, ottenuta %d", fRes.currentPage)
	}
}

// ---------------------------------------------------------------------------
// Update – messaggio sconosciuto senza pagina valida
// ---------------------------------------------------------------------------

// TestUpdate_PaginaNonGestita verifica che un messaggio generico con una
// pagina non gestita restituisca il modello invariato e cmd nil.
func TestUpdate_PaginaNonGestita(t *testing.T) {
	f := nuovaFinestraTest()
	f.currentPage = 99 // valore non gestito

	res, cmd := f.Update(tea.KeyPressMsg{Text: "x"})
	fRes := res.(mainLayout)

	if fRes.currentPage != 99 {
		t.Errorf("pagina corrente attesa 99, ottenuta %d", fRes.currentPage)
	}
	if cmd != nil {
		t.Error("cmd atteso nil per pagina non gestita")
	}
}

// ---------------------------------------------------------------------------
// View
// ---------------------------------------------------------------------------

// TestView_PaginaLogin verifica che il rendering della vista con pagina LOGIN
// produca un contenuto non vuoto.
func TestView_PaginaLogin(t *testing.T) {
	f := nuovaFinestraTest()
	f.currentPage = LOGIN

	view := f.View()
	if view.Content == "" {
		t.Error("View con pagina LOGIN dovrebbe produrre contenuto non vuoto")
	}
}

// TestView_PaginaMenuPrincipale verifica che il rendering della vista con
// pagina MAGMENU produca un contenuto non vuoto.
func TestView_PaginaMenuPrincipale(t *testing.T) {
	f := nuovaFinestraTest()
	f.currentPage = MAGMENU

	view := f.View()
	if view.Content == "" {
		t.Error("View con pagina MAGMENU dovrebbe produrre contenuto non vuoto")
	}
}

// TestView_AltScreenAbilitato verifica che la View restituisca una vista
// con AltScreen abilitato (l'app occupa l'intero terminale).
func TestView_AltScreenAbilitato(t *testing.T) {
	f := nuovaFinestraTest()

	view := f.View()
	if !view.AltScreen {
		t.Error("AltScreen dovrebbe essere abilitato nella View")
	}
}
