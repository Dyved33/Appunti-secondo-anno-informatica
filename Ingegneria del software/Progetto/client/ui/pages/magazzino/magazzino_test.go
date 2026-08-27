// Package magazzino_test raccoglie i test d'unità per il router ed il layout delle pagine del magazzino.
package magazzino

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/api"
)

// TestMagLayout_Init verifica che il costruttore New() ed il metodo Init() del layout Magazzino
// impostino la sotto-pagina di default a GIACENZE ed eseguano l'inizializzazione senza sollevare eccezioni.
func TestMagLayout_Init(t *testing.T) {
	client := &api.APIClient{}
	m := New(client)

	// Verifichiamo che la sezione iniziale di default sia GIACENZE
	if m.attiva != GIACENZE {
		t.Errorf("Stato iniziale non corretto: ottenuto %v, atteso %v", m.attiva, GIACENZE)
	}

	// Verifica che Init() risponda senza panico
	cmd := m.Init()
	_ = cmd
}

// TestMagLayout_Update_TabSwitching verifica la commutazione corretta della sotto-pagina attiva
// al click del mouse sui tab della barra di navigazione, oltre al messaggio speciale "SCANNER".
func TestMagLayout_Update_TabSwitching(t *testing.T) {
	const winW, winH = 120, 40

	tests := []struct {
		name      string
		tab       focus // tab da cliccare; -1 per usare il messaggio tastiera
		keyMsg    tea.KeyPressMsg
		wantFocus focus
	}{
		{
			name:      "Click sul tab Ordini In Entrata",
			tab:       ORDINIENTRATA,
			wantFocus: ORDINIENTRATA,
		},
		{
			name:      "Click sul tab Ordini In Uscita",
			tab:       ORDINIUSCITA,
			wantFocus: ORDINIUSCITA,
		},
		{
			name:      "Click sul tab Resi",
			tab:       RESI,
			wantFocus: RESI,
		},
		{
			name:      "Messaggio Scanner (trigger hardware)",
			tab:       -1,
			keyMsg:    tea.KeyPressMsg{Text: "SCANNER"},
			wantFocus: SCANNER,
		},
		{
			name:      "Click sul tab Giacenze",
			tab:       GIACENZE,
			wantFocus: GIACENZE,
		},
	}

	client := &api.APIClient{}
	m := New(client)
	m.DimUpdate(tea.WindowSizeMsg{Width: winW, Height: winH})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var msg tea.Msg
			if tt.tab >= 0 {
				// Centro dell'area cliccabile del tab indicato, calcolata
				// dal layout stesso (stessa strada usata in produzione).
				r := m.hitboxes()[tt.tab]
				msg = tea.MouseClickMsg{X: (r.X0 + r.X1) / 2, Y: (r.Y0 + r.Y1) / 2, Button: tea.MouseLeft}
			} else {
				msg = tt.keyMsg
			}

			updatedModel, _ := m.Update(msg)
			typedModel, ok := updatedModel.(MagLayout)
			if !ok {
				t.Fatalf("Update() non ha restituito un tipo MagLayout")
			}

			if typedModel.attiva != tt.wantFocus {
				t.Errorf("Focus errato dopo %v: ottenuto %v, atteso %v", msg, typedModel.attiva, tt.wantFocus)
			}

			// Sincronizziamo m per il test successivo in sequenza
			m = typedModel
		})
	}
}

// TestMagLayout_View verifica che la funzione View() generi correttamente l'interfaccia utente
// ed includa tutte le etichette dei tab previste.
func TestMagLayout_View(t *testing.T) {
	client := &api.APIClient{}
	m := New(client)

	view := m.View()

	// Verifichiamo che la View generata contenga le etichette dei Tab attesi
	for _, label := range labels {
		if !strings.Contains(view.Content, label) {
			t.Errorf("View() non contiene l'etichetta del tab %q", label)
		}
	}
}
