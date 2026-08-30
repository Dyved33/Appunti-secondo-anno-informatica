package contabilita

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/api"
)

// TestNew_StatoIniziale verifica che New inizializzi la sezione attiva su SezioneOrdini
// e che le dimensioni iniziali siano impostate a (0, 0).
func TestNew_StatoIniziale(t *testing.T) {
	t.Parallel()

	client := api.New("http://localhost:8080")
	c := New(client)

	if got := c.Attiva(); got != int(SezioneOrdini) {
		t.Errorf("attesa sezione iniziale SezioneOrdini (%d), ottenuta %d", SezioneOrdini, got)
	}

	w, h := c.Size()
	if w != 0 || h != 0 {
		t.Errorf("attese dimensioni iniziali (0, 0), ottenute (%d, %d)", w, h)
	}
}

// TestInit verifica che Init() deleghi l'avvio alla sotto-pagina attiva ordini.
func TestInit(t *testing.T) {
	t.Parallel()

	c := New(nil)
	cmd := c.Init()

	if cmd == nil {
		t.Error("Init() ha restituito un Cmd nil")
	}
}

// TestDimUpdate verifica che DimUpdate aggiorni le dimensioni del terminale
// e le propaghi a tutte le 4 sotto-pagine di contabilità.
func TestDimUpdate(t *testing.T) {
	t.Parallel()

	c := New(nil)
	c.DimUpdate(tea.WindowSizeMsg{Width: 100, Height: 40})

	w, h := c.Size()
	if w != 100 || h != 40 {
		t.Errorf("attese dimensioni (100, 40), ottenute (%d, %d)", w, h)
	}
}

// TestSwitchTo_E_Routing verifica che switchTo cambi correttamente la sezione attiva
// per ciascuna delle 4 sezioni e gestisca eventuali indici non validi.
func TestSwitchTo_E_Routing(t *testing.T) {
	t.Parallel()

	tests := []struct {
		dest     focus
		expected focus
	}{
		{SezioneFatture, SezioneFatture},
		{SezioneCatalogo, SezioneCatalogo},
		{SezioneReport, SezioneReport},
		{SezioneOrdini, SezioneOrdini},
	}

	for _, tt := range tests {
		c := New(nil)
		m, cmd := c.switchTo(tt.dest)
		updated := m.(ContLayout)

		if updated.attiva != tt.expected {
			t.Errorf("switchTo(%v): attesa sezione %v, ottenuta %v", tt.dest, tt.expected, updated.attiva)
		}
		if cmd == nil {
			t.Errorf("switchTo(%v): atteso Cmd non nil per l'init della sotto-pagina", tt.dest)
		}
	}

	// Caso indice non valido
	c := New(nil)
	m, cmd := c.switchTo(focus(999))
	updated := m.(ContLayout)
	if updated.attiva != focus(999) {
		t.Errorf("attesa attiva=999, ottenuta %v", updated.attiva)
	}
	if cmd != nil {
		t.Errorf("atteso cmd nil per focus non valido, ottenuto %v", cmd)
	}
}

// TestUpdate_InoltroMessaggi verifica l'inoltro dei messaggi alla sotto-pagina attiva per tutte le sezioni.
func TestUpdate_InoltroMessaggi(t *testing.T) {
	t.Parallel()

	for att := SezioneOrdini; att <= SezioneReport; att++ {
		c := New(nil)
		c.attiva = att

		res, _ := c.Update(tea.KeyPressMsg{Text: "x"})
		m := res.(ContLayout)

		if m.attiva != att {
			t.Errorf("Update non doveva cambiare la sezione attiva del layout padre, ottenuta %v", m.attiva)
		}
	}
}

// TestView_E_Cursor verifica il rendering della vista e la propagazione del cursore.
func TestView_E_Cursor(t *testing.T) {
	t.Parallel()

	for att := SezioneOrdini; att <= SezioneReport; att++ {
		c := New(nil)
		c.DimUpdate(tea.WindowSizeMsg{Width: 120, Height: 30})
		c.attiva = att

		v := c.View()
		if v.Content == "" {
			t.Errorf("View() con attiva=%v ha prodotto contenuto vuoto", att)
		}

		// Verifica che l'etichetta del tab sia inclusa
		if !strings.Contains(v.Content, labels[att]) {
			t.Errorf("View() non contiene l'etichetta del tab %q", labels[att])
		}

		_ = c.activeCursor()
	}

	// Caso attiva non valida per activeCursor
	cInvalid := New(nil)
	cInvalid.attiva = focus(999)
	if cursor := cInvalid.activeCursor(); cursor != nil {
		t.Errorf("activeCursor() con indice non valido atteso nil, ottenuto %+v", cursor)
	}
}

// TestHitboxes verifica il getter Hitboxes.
func TestHitboxes(t *testing.T) {
	t.Parallel()

	c := New(nil)
	c.DimUpdate(tea.WindowSizeMsg{Width: 120, Height: 40})

	hitboxes := c.Hitboxes()
	if len(hitboxes) != len(labels) {
		t.Errorf("Hitboxes() attese %d, ottenute %d", len(labels), len(hitboxes))
	}
}

// TestJoinWithSpacing verifica la funzione joinWithSpacing con input vuoto, singolo e multiplo.
func TestJoinWithSpacing(t *testing.T) {
	t.Parallel()

	if got := joinWithSpacing([]string{}, " "); got != nil {
		t.Errorf("atteso nil per slice vuota, ottenuto %v", got)
	}

	singolo := joinWithSpacing([]string{"A"}, " ")
	if len(singolo) != 1 || singolo[0] != "A" {
		t.Errorf("risultato inatteso per singolo elemento: %v", singolo)
	}

	multiplo := joinWithSpacing([]string{"A", "B", "C"}, "  ")
	expected := []string{"A", "  ", "B", "  ", "C"}
	if len(multiplo) != len(expected) {
		t.Fatalf("lunghezza inattesa: attesa %d, ottenuta %d", len(expected), len(multiplo))
	}
	for i, v := range multiplo {
		if v != expected[i] {
			t.Errorf("indice %d: atteso %q, ottenuto %q", i, expected[i], v)
		}
	}
}

// TestHitboxTab_ClickSempreCorretto verifica che per tutte le risoluzioni terminale
// e per ciascuna sezione attiva, il click al centro di ciascun tab cambi correttamente sezione.
func TestHitboxTab_ClickSempreCorretto(t *testing.T) {
	t.Parallel()

	sizes := []struct{ w, h int }{{80, 24}, {100, 30}, {120, 40}, {160, 50}}

	for _, s := range sizes {
		for att := SezioneOrdini; att <= SezioneReport; att++ {
			c := New(nil)
			c.DimUpdate(tea.WindowSizeMsg{Width: s.w, Height: s.h})
			c.attiva = att

			rects := c.hitboxes()
			if len(rects) != len(labels) {
				t.Fatalf("%dx%d attiva=%v: attese %d hitbox, ottenute %d", s.w, s.h, att, len(labels), len(rects))
			}

			// Un click fuori da qualsiasi area non deve cambiare sezione.
			res, _ := c.Update(tea.MouseClickMsg{X: 0, Y: 0, Button: tea.MouseLeft})
			if got := res.(ContLayout).attiva; got != att {
				t.Errorf("%dx%d attiva=%v: click fuori area ha cambiato sezione in %v", s.w, s.h, att, got)
			}

			// Il click al centro di ciascun tab deve selezionarlo.
			for i, r := range rects {
				msg := tea.MouseClickMsg{X: (r.X0 + r.X1) / 2, Y: (r.Y0 + r.Y1) / 2, Button: tea.MouseLeft}
				res, _ := c.Update(msg)
				if got := res.(ContLayout).attiva; got != focus(i) {
					t.Errorf("%dx%d attiva=%v: click sul tab %q ha selezionato %v", s.w, s.h, att, labels[i], got)
				}
			}
		}
	}
}
