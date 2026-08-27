// Test di regressione: dopo il login i layout nascono SENZA aver mai ricevuto
// un WindowSizeMsg (le dimensioni arrivano all'avvio, prima dell'autenticazione).
// Se le dimensioni non vengono iniettate alla creazione, le hitbox dei tab
// cliccabili vengono misurate su un terminale fittizio 80x24 mentre la View è
// disegnata sul terminale reale: i click dell'utente cadono nel vuoto.
package app

import (
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/pages/login"

	tea "charm.land/bubbletea/v2"
)

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
	// (ciò che l'utente vede in quel momento: ogni cambio pagina può spostare
	// il centramento) e cliccaci con coordinate reali del terminale.
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
