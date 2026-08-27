package magazzino

import (
	"strings"
	"testing"
	"time"

	ordinientrata "github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/pages/magazzino/ordinientrata"

	tea "charm.land/bubbletea/v2"

	"github.com/tommasogarofani/ingsw_25-26_group-11/client/api"
)

// nuoviOrdiniMock genera n ordini in entrata di prova.
func nuoviOrdiniMock(n int) []api.OrdineEntrata {
	ordini := make([]api.OrdineEntrata, 0, n)
	for i := range n {
		ordini = append(ordini, api.OrdineEntrata{
			ID:        i + 1,
			Numero:    "ORD-2026-" + string(rune('A'+i)),
			Fornitore: "Fornitore Test",
			Data:      time.Date(2026, 8, 24, 0, 0, 0, 0, time.UTC),
			Righe:     []api.RigaOrdineEntrata{{ProdottoID: 1, Quantita: 10, PrezzoUnitario: 1.0}},
		})
	}
	return ordini
}

// TestMagLayout_RotellinaScorrePagina simula la catena completa:
// layout -> UpdateActive -> pagina Ordini in Entrata, inviando eventi
// rotellina e verificando che la vista mostri righe diverse.
func TestMagLayout_RotellinaScorrePagina(t *testing.T) {
	m := New(nil)
	m.attiva = ORDINIENTRATA

	// Dimensiona il terminale come quello reale.
	res, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	m = res.(MagLayout)

	// Carica molti ordini per garantire lo scorrimento.
	res, _ = m.Update(ordinientrata.OrdiniEntrataLoadedMsg{Ordini: nuoviOrdiniMock(60)})
	m = res.(MagLayout)

	v1 := m.View().Content
	if !strings.Contains(v1, "ORD-2026-A") {
		t.Fatalf("primo ordine non visibile nella vista iniziale")
	}

	// Ruota la rotellina in giù tre volte.
	for range 3 {
		res, _ = m.Update(tea.MouseWheelMsg{Button: tea.MouseWheelDown})
		m = res.(MagLayout)
	}

	v2 := m.View().Content
	if strings.Contains(v2, "ORD-2026-A") {
		t.Error("dopo la rotazione della rotellina la vista è invariata: la pagina non scorre")
	}
	if !strings.Contains(v2, "di 60") {
		t.Errorf("atteso indicatore 'di 60' nel footer, ottenuto:\n%s", v2[len(v2)-200:])
	}
}
