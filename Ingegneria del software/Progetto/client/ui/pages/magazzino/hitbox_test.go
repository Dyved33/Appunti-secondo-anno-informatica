// Test di regressione per le hitbox dei tab: verificano che il click funzioni
// per OGNI tab, in OGNI stato attivo e su PIÙ dimensioni di terminale,
// usando la stessa composizione condivisa con app.go (tabnav.PlacePage).
package magazzino

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/api"
)

func TestHitboxTab_ClickSempreCorretto(t *testing.T) {
	sizes := []struct{ w, h int }{{80, 24}, {100, 30}, {120, 40}, {160, 50}}

	for _, s := range sizes {
		for att := GIACENZE; att <= SCANNER; att++ {
			m := New(&api.APIClient{})
			m.DimUpdate(tea.WindowSizeMsg{Width: s.w, Height: s.h})
			m.attiva = att

			rects := m.hitboxes()
			if len(rects) != len(labels) {
				t.Fatalf("%dx%d attiva=%v: attese %d hitbox, ottenute %d", s.w, s.h, att, len(labels), len(rects))
			}

			// Un click fuori da qualsiasi area non deve cambiare pagina.
			res, _ := m.Update(tea.MouseClickMsg{X: 0, Y: 0, Button: tea.MouseLeft})
			if got := res.(MagLayout).attiva; got != att {
				t.Errorf("%dx%d attiva=%v: click fuori area ha cambiato pagina in %v", s.w, s.h, att, got)
			}

			// Il click al centro di ciascun tab deve selezionarlo.
			for i, r := range rects {
				msg := tea.MouseClickMsg{X: (r.X0 + r.X1) / 2, Y: (r.Y0 + r.Y1) / 2, Button: tea.MouseLeft}
				res, _ := m.Update(msg)
				if got := res.(MagLayout).attiva; got != focus(i) {
					t.Errorf("%dx%d attiva=%v: click sul tab %q ha selezionato %v", s.w, s.h, att, labels[i], got)
				}
			}
		}
	}
}
