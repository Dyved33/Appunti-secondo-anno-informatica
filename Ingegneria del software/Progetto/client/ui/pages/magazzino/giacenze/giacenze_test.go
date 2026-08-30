package giacenze

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/tommasogarofani/ingsw_25-26_group-11/client/api"
	style "github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/styles"
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/tabnav"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
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

// TestGiacenze_FooterSempreVisibile verifica la regressione del "footer invisibile":
// qualunque sia l'altezza del terminale, le scorciatoie di navigazione devono restare
// presenti nella vista composta (il numero di righe della tabella si adegua allo schermo).
func TestGiacenze_FooterSempreVisibile(t *testing.T) {
	client := &api.APIClient{}
	m := New(client)

	report := &api.InventoryReport{
		GeneratoIl: time.Now(),
		Articoli:   make([]api.InventoryItem, 26),
	}
	for i := range report.Articoli {
		report.Articoli[i] = api.InventoryItem{
			ProdottoID: i + 1,
			Codice:     fmt.Sprintf("PRD-%03d", i+1),
			Nome:       fmt.Sprintf("Articolo %d", i+1),
			Categoria:  "Cat",
			Quantita:   10,
			SogliaMin:  5,
		}
	}

	for _, winH := range []int{30, 40, 50} {
		upd, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: winH})
		m = upd.(Model)
		upd, _ = m.Update(GiacenzeLoadedMsg{Report: report})
		m = upd.(Model)

		view := m.View()
		if !strings.Contains(view.Content, "[pgup/pgdown o rotellina] Scorri") {
			t.Errorf("winH=%d: le scorciatoie del footer non sono presenti nella vista", winH)
		}

		tabs := lipgloss.JoinHorizontal(lipgloss.Top,
			style.ActiveTab.Render("GIACENZE"), "  ", style.TabGap.Render("ORDINI IN ENTRATA"))
		composed := tabnav.PlacePage(120, winH,
			style.FinestraStile.Render(tabs+"\n\n"+view.Content))
		if strings.Count(composed, "\n")+1 > winH {
			t.Errorf("winH=%d: la pagina composta occupa %d righe (footer fuori schermo)",
				winH, strings.Count(composed, "\n")+1)
		}
	}

	// Lo scroll con la rotellina deve poter raggiungere l'ultimo articolo.
	upd, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 30})
	m = upd.(Model)
	for i := 0; i < 40; i++ {
		upd, _ = m.Update(tea.MouseWheelMsg{Button: tea.MouseWheelDown})
		m = upd.(Model)
	}
	if m.cursore != len(report.Articoli)-1 {
		t.Errorf("Dopo lo scroll il cursore dovrebbe essere sull'ultimo articolo (%d), trovato %d",
			len(report.Articoli)-1, m.cursore)
	}
}

func TestGiacenze_Errore(t *testing.T) {
	m := New(&api.APIClient{})
	res, _ := m.Update(ErroreMsg{Err: fmt.Errorf("connection refused")})
	mErr := res.(Model)
	if !strings.Contains(mErr.View().Content, "connection refused") {
		t.Error("La View deve mostrare il testo dell'errore")
	}
}
