package giacenze

import (
	"strings"
	"testing"
	"time"

	"github.com/tommasogarofani/ingsw_25-26_group-11/client/api"

	tea "charm.land/bubbletea/v2"
)

func TestGiacenze_InitAndUpdate(t *testing.T) {
	client := &api.APIClient{}
	m := New(client)

	cmd := m.Init()
	if cmd == nil {
		t.Error("Init() dovrebbe restituire un Cmd per il caricamento")
	}

	report := &api.InventoryReport{
		GeneratoIl:     time.Now(),
		TotaleQuantita: 150,
		ValoreTotale:   2500.50,
		Articoli: []api.InventoryItem{
			{
				ProdottoID: 1, Codice: "PRD-001", Nome: "Bulloni M8",
				Categoria: "Hardware", Quantita: 100, SogliaMin: 20, PrezzoVend: 1.50, ValoreRiga: 150.0, SottoSoglia: false,
			},
			{
				ProdottoID: 2, Codice: "PRD-002", Nome: "Dadi M8",
				Categoria: "Hardware", Quantita: 5, SogliaMin: 50, PrezzoVend: 0.50, ValoreRiga: 2.50, SottoSoglia: true,
			},
		},
		ArticoliSottoSoglia: []api.InventoryItem{
			{
				ProdottoID: 2, Codice: "PRD-002", Nome: "Dadi M8",
				Categoria: "Hardware", Quantita: 5, SogliaMin: 50, PrezzoVend: 0.50, ValoreRiga: 2.50, SottoSoglia: true,
			},
		},
	}

	// Test GiacenzeLoadedMsg
	res, _ := m.Update(GiacenzeLoadedMsg{Report: report})
	mRes := res.(Model)

	if mRes.loading {
		t.Error("loading dovrebbe essere false dopo GiacenzeLoadedMsg")
	}
	if len(mRes.report.Articoli) != 2 {
		t.Errorf("Articoli attesi 2, trovati %d", len(mRes.report.Articoli))
	}

	// Test tasto 's' per toggle filtro sotto soglia
	resFilter, _ := mRes.Update(tea.KeyPressMsg{Text: "s"})
	mFilter := resFilter.(Model)
	if !mFilter.filtroSottoSoglia {
		t.Error("filtroSottoSoglia dovrebbe essere true dopo il tasto 's'")
	}
	if len(mFilter.articoliFiltrati()) != 1 {
		t.Errorf("Articoli filtrati attesi 1, trovati %d", len(mFilter.articoliFiltrati()))
	}

	// Test View
	view := mFilter.View()
	if !strings.Contains(view.Content, "GIACENZE") {
		t.Error("View() dovrebbe contenere il titolo GIACENZE")
	}
	if !strings.Contains(view.Content, "Dadi M8") {
		t.Error("View() dovrebbe contenere l'articolo filtrato Dadi M8")
	}
}
