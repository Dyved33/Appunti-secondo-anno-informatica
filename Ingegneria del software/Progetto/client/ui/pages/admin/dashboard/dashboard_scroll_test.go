package dashboard

import (
	"strings"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/tommasogarofani/ingsw_25-26_group-11/client/api"
)

// nuoviReportConScorte genera un report giacenze con n articoli sotto soglia,
// così la lista scorte critiche diventa più lunga dello schermo.
func nuoviReportConScorte(n int) *api.InventoryReport {
	items := make([]api.InventoryItem, 0, n)
	for i := range n {
		items = append(items, api.InventoryItem{
			ProdottoID: i + 1,
			Codice:     "SKU-" + string(rune('A'+i%26)),
			Nome:       "Articolo critico numero",
			Quantita:   1,
			SogliaMin:  5,
		})
	}
	return &api.InventoryReport{
		GeneratoIl:          time.Date(2026, 8, 24, 0, 0, 0, 0, time.UTC),
		TotaleQuantita:      100,
		ValoreTotale:        1000,
		ArticoliSottoSoglia: items,
	}
}

// TestDashboard_RotellinaScorrePagina verifica che su un terminale basso la
// dashboard risponda alla rotellina mostrando porzioni successive del corpo.
func TestDashboard_RotellinaScorrePagina(t *testing.T) {
	m := New(nil)

	// Terminale piccolo: il contenuto eccede e serve lo scorrimento.
	res, _ := m.Update(tea.WindowSizeMsg{Width: 110, Height: 30})
	m = res.(Model)

	res, _ = m.Update(GiacenzeCaricateMsg{Report: nuoviReportConScorte(40)})
	m = res.(Model)

	v1 := m.View().Content

	if !strings.Contains(v1, "Righe") || !strings.Contains(v1, "per scorrere") {
		t.Fatalf("atteso indicatore di scorrimento con contenuto che eccede, ottenuto:\n%s", v1)
	}

	// Ruota la rotellina verso il basso finché la finestra si sposta.
	primaRiga := m.scrollOffset
	for range 4 {
		res, _ = m.Update(tea.MouseWheelMsg{Button: tea.MouseWheelDown})
		m = res.(Model)
	}
	if m.scrollOffset == primaRiga {
		t.Fatal("la rotellina non ha spostato scrollOffset")
	}

	v2 := m.View().Content
	if v2 == v1 {
		t.Error("la vista non è cambiata dopo lo scorrimento")
	}
	if !strings.Contains(v2, "di ") {
		t.Errorf("atteso indicatore 'X–Y di Z' anche dopo lo scorrimento:\n%s", v2)
	}
}

// TestDashboard_SenzaOverflowNessunIndicatore verifica che su terminali alti
// la dashboard resti interamente visibile senza indicatori superflui.
func TestDashboard_SenzaOverflowNessunIndicatore(t *testing.T) {
	m := New(nil)

	res, _ := m.Update(tea.WindowSizeMsg{Width: 110, Height: 60})
	m = res.(Model)

	res, _ = m.Update(GiacenzeCaricateMsg{Report: nuoviReportConScorte(3)})
	m = res.(Model)

	v := m.View().Content
	if strings.Contains(v, "per scorrere") {
		t.Errorf("indicatore di scorrimento presente pur entrando tutto nello schermo:\n%s", v)
	}
	if !strings.Contains(v, "ctrl+l") {
		t.Error("footer scorciatoie assente nella vista senza overflow")
	}
}
