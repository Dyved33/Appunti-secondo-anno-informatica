package app

import (
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/api"
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/pages/login"
	style "github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/styles"

	tea "charm.land/bubbletea/v2"
)

// Helper per creare una mainLayout pronta all'uso nei test.
func nuovaFinestraTest() mainLayout {
	client := &api.APIClient{}
	f := New(client)
	// Assegna dimensioni ragionevoli per evitare altezze negative nel View.
	f.WindowSize.ChangeSize(120, 40)
	f.WindowSize.Border.SetSpacing(2, 2, 2, 2)
	return f
}

// ---------------------------------------------------------------------------
// Costruttore - New
// ---------------------------------------------------------------------------

// TestNew verifica che il costruttore New() restituisca una mainLayout
// con la pagina iniziale impostata a LOGIN.
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
// parseUserRole
// ---------------------------------------------------------------------------

// TestParseUserRole verifica la conversione da stringa a tipo enum userRole,
// inclusa la gestione del caso di default per ruoli non riconosciuti o vuoti.
func TestParseUserRole(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected userRole
	}{
		{"Magazziniere", "magazziniere", MAGAZZINIERE},
		{"Contabile", "contabile", CONTABILE},
		{"Admin", "admin", ADMIN},
		{"Ruolo non valido (fallback)", "sconosciuto", MAGAZZINIERE},
		{"Stringa vuota (fallback)", "", MAGAZZINIERE},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseUserRole(tt.input)
			if got != tt.expected {
				t.Errorf("parseUserRole(%q) = %v, atteso %v", tt.input, got, tt.expected)
			}
		})
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
// la pagina corrente passi da LOGIN al menu corrispondente al ruolo dell'utente.
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
// Update – Shortcut Keyboard (ctrl+l / Logout)
// ---------------------------------------------------------------------------

// TestUpdate_ShortcutCtrlL_ResetToLogin verifica che la combinazione ctrl+l
// quando l'utente è autenticato resetti lo stato e riporti alla schermata LOGIN.
func TestUpdate_ShortcutCtrlL_ResetToLogin(t *testing.T) {
	f := nuovaFinestraTest()
	f.currentPage = MAGMENU

	msg := tea.KeyPressMsg{Text: "ctrl+l"}
	res, cmd := f.Update(msg)
	fRes := res.(mainLayout)

	if fRes.currentPage != LOGIN {
		t.Errorf("dopo ctrl+l la pagina attesa è LOGIN (%d), ottenuta %d", LOGIN, fRes.currentPage)
	}
	if cmd == nil {
		t.Error("ctrl+l dovrebbe restituire il Cmd di Init() per re-inizializzare il sotto-modello Login")
	}
}

// TestUpdate_ShortcutCtrlL_IgnoratoInLogin verifica che ctrl+l non scatti
// se l'utente si trova già nella schermata di LOGIN.
func TestUpdate_ShortcutCtrlL_IgnoratoInLogin(t *testing.T) {
	f := nuovaFinestraTest()
	f.currentPage = LOGIN

	msg := tea.KeyPressMsg{Text: "ctrl+l"}
	res, _ := f.Update(msg)
	fRes := res.(mainLayout)

	if fRes.currentPage != LOGIN {
		t.Errorf("pagina corrente attesa LOGIN, ottenuta %d", fRes.currentPage)
	}
}

// ---------------------------------------------------------------------------
// Update – Propagazione messaggi
// ---------------------------------------------------------------------------

// TestUpdate_PropagaMessaggiALogin verifica che i messaggi di tastiera
// vengano inoltrati al sotto-modello Login quando la pagina corrente è LOGIN.
func TestUpdate_PropagaMessaggiALogin(t *testing.T) {
	f := nuovaFinestraTest()
	f.currentPage = LOGIN

	msg := tea.KeyPressMsg{Text: "a"}
	res, _ := f.Update(msg)
	fRes := res.(mainLayout)

	if fRes.currentPage != LOGIN {
		t.Errorf("pagina corrente attesa LOGIN, ottenuta %d", fRes.currentPage)
	}
}

// TestUpdate_PropagaMessaggiAMenu verifica che i messaggi
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
// Update – MouseClickMsg / Tab Navigation
// ---------------------------------------------------------------------------

// TestClickTabDopoLogin_SenzaResizeIntermedio (Test di regressione):
// Verifica che dopo il login le dimensioni del terminale vengano propagate ai layout
// e che le hitbox dei tab cliccabili rispondano correttamente ai click del mouse.
func TestClickTabDopoLogin_SenzaResizeIntermedio(t *testing.T) {
	const termW, termH = 150, 45

	f := nuovaFinestraTest()

	// 1. All'avvio bubbletea invia le dimensioni del terminale...
	res, _ := f.Update(tea.WindowSizeMsg{Width: termW, Height: termH})
	f = res.(mainLayout)

	// 2. ...poi l'utente effettua il login (nessun resize nel mezzo).
	res, _ = f.Update(login.LoginSubmitMsg{})
	f = res.(mainLayout)
	if f.currentPage != MAGMENU {
		t.Fatalf("attesa pagina MAGMENU, ottenuta %v", f.currentPage)
	}

	// 3. Il layout deve conoscere le dimensioni REALI del terminale.
	if w, h := f.MagMenu.Size(); w != termW || h != termH {
		t.Fatalf("layout con dimensioni (%d,%d), attese (%d,%d): dimensioni non propagate al login", w, h, termW, termH)
	}

	// 4/5. Per ogni tab: leggi la posizione dell'etichetta dalla View CORRENTE
	// e cliccaci con coordinate reali del terminale.
	casi := []struct {
		label string
		want  int // indice focus atteso in magazzino
	}{
		{"Resi", 3},
		{"Ordini In Entrata", 1},
		{"Scanner", 4},
		{"Giacenze", 0},
	}
	for _, c := range casi {
		lines := strings.Split(ansi.Strip(f.View().Content), "\n")
		x, y := -1, -1
		for ly, line := range lines {
			if idx := strings.Index(line, c.label); idx >= 0 && ly < 40 {
				x, y = ansi.StringWidth(line[:idx]), ly
				break
			}
		}
		if x < 0 {
			t.Fatalf("etichetta %q non trovata nella view", c.label)
		}

		msg := tea.MouseClickMsg{X: x + 1, Y: y, Button: tea.MouseLeft}
		res, _ := f.Update(msg)
		f = res.(mainLayout)
		if got := f.MagMenu.Attiva(); got != c.want {
			t.Errorf("click su %q a (%d,%d): attesa sotto-pagina %d, ottenuta %d", c.label, x, y, c.want, got)
		}
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
