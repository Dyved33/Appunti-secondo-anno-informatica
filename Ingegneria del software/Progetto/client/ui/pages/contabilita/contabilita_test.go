package contabilita

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
)

func TestNew_StatoIniziale(t *testing.T) {
	c := New(nil)

	if got := c.Attiva(); got != int(SezioneOrdini) {
		t.Errorf("attesa sezione iniziale SezioneOrdini (%d), ottenuta %d", SezioneOrdini, got)
	}

	w, h := c.Size()
	if w != 0 || h != 0 {
		t.Errorf("attese dimensioni iniziali (0, 0), ottenute (%d, %d)", w, h)
	}
}

func TestInit(t *testing.T) {
	c := New(nil)
	cmd := c.Init()

	if cmd == nil {
		t.Error("Init() ha restituito un Cmd nil")
	}
}

func TestDimUpdate(t *testing.T) {
	c := New(nil)
	c.DimUpdate(tea.WindowSizeMsg{Width: 100, Height: 40})

	w, h := c.Size()
	if w != 100 || h != 40 {
		t.Errorf("attese dimensioni (100, 40), ottenute (%d, %d)", w, h)
	}
}

func TestSwitchTo_E_Routing(t *testing.T) {
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
}

func TestUpdate_InoltroMessaggi(t *testing.T) {
	c := New(nil)
	c.attiva = SezioneOrdini

	// Un evento tastiera standard deve essere inoltrato al sotto-modello attivo
	res, _ := c.Update(tea.KeyPressMsg{Text: "tab"})
	m := res.(ContLayout)

	if m.attiva != SezioneOrdini {
		t.Errorf("l'update non doveva cambiare la sezione attiva del layout padre, ottenuta %v", m.attiva)
	}
}

func TestView_E_Cursor(t *testing.T) {
	c := New(nil)
	c.DimUpdate(tea.WindowSizeMsg{Width: 120, Height: 30})

	v := c.View()

	// Verifica che le etichette dei tab siano incluse nella vista
	for _, label := range labels {
		if !strings.Contains(v.Content, label) {
			t.Errorf("View() non contiene l'etichetta del tab %q", label)
		}
	}

	// In assenza di form con focus attivo, il cursore deve essere nil
	if c.activeCursor() != nil {
		t.Error("activeCursor() atteso nil nella vista predefinita")
	}
}

func TestJoinWithSpacing(t *testing.T) {
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

func TestHitboxTab_ClickSempreCorretto(t *testing.T) {
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
