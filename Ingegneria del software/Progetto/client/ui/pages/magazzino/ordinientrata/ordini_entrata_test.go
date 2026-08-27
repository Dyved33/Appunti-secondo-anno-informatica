package ordinientrata

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

// nuoviOrdiniMock genera n ordini di prova senza alcuno stato (non esiste più nel backend).
func nuoviOrdiniMock(n int) []api.OrdineEntrata {
	ordini := make([]api.OrdineEntrata, 0, n)
	for i := 0; i < n; i++ {
		ordini = append(ordini, api.OrdineEntrata{
			ID:        100 + i,
			Numero:    fmt.Sprintf("ENT-2026-%04d", i+1),
			Fornitore: fmt.Sprintf("Fornitore %d", i%3+1),
			Data:      time.Now(),
			Note:      "Consegna prevista lunedì",
			Righe: []api.RigaOrdineEntrata{
				{
					ID: i + 1, OrdineEntrataID: 100 + i, ProdottoID: 101 + i,
					Quantita: 10 * (i + 1), PrezzoUnitario: 2.5,
					Prodotto: &api.Prodotto{Codice: "PRD-101", Nome: "Tasselli 10mm"},
				},
			},
		})
	}
	return ordini
}

func TestOrdiniEntrata_InitAndUpdate(t *testing.T) {
	client := &api.APIClient{}
	m := New(client)

	cmd := m.Init()
	if cmd == nil {
		t.Error("Init() dovrebbe restituire un Cmd non-nil")
	}

	ordiniMock := nuoviOrdiniMock(1)

	res, _ := m.Update(OrdiniEntrataLoadedMsg{Ordini: ordiniMock})
	mRes := res.(Model)

	if mRes.loading {
		t.Error("loading dovrebbe essere false dopo caricamento ordini")
	}
	if len(mRes.ordini) != 1 {
		t.Fatalf("Ordini attesi 1, ottenuti %d", len(mRes.ordini))
	}

	// Test espansione dettaglio ordine con Enter
	resExpand, _ := mRes.Update(tea.KeyPressMsg{Text: "enter"})
	mExpand := resExpand.(Model)

	if !mExpand.dettagliEspansi[100] {
		t.Error("L'ordine #100 dovrebbe essere espanso dopo Enter")
	}

	view := mExpand.View()
	if !strings.Contains(view.Content, "ENT-2026-0001") {
		t.Error("View() dovrebbe mostrare il numero ordine ENT-2026-0001")
	}
	if !strings.Contains(view.Content, "Tasselli 10mm") {
		t.Error("View() espansa dovrebbe mostrare il dettaglio prodotto Tasselli 10mm")
	}
}

func TestOrdiniEntrata_SenzaColonnaStato(t *testing.T) {
	client := &api.APIClient{}
	m := New(client)
	res, _ := m.Update(OrdiniEntrataLoadedMsg{Ordini: nuoviOrdiniMock(2)})
	m = res.(Model)

	view := m.View().Content

	if strings.Contains(view, "STATO") {
		t.Error("La tabella non deve più contenere la colonna STATO per gli ordini in entrata")
	}
	if strings.Contains(view, "CONFERMATI") || strings.Contains(view, "EVASI") || strings.Contains(view, "ANNULLATI") {
		t.Error("I KPI legati agli stati non devono più comparire: gli ordini in entrata non hanno stati")
	}
	for _, atteso := range []string{"TOTALE ORDINI", "FORNITORI", "ARTICOLI", "VALORE TOTALE"} {
		if !strings.Contains(view, atteso) {
			t.Errorf("Manca il blocco KPI %q", atteso)
		}
	}
	// Fornitore distinto: 1 (mock con indice i%3 → due ordini → fornitori 1 e 2)
	if !strings.Contains(view, "Fornitore 1") {
		t.Error("La tabella dovrebbe mostrare la colonna FORNITORE")
	}
}

func TestOrdiniEntrata_ScorrimentoTastiera(t *testing.T) {
	client := &api.APIClient{}
	m := New(client)
	m.winW, m.winH = 120, 30 // terminale piccolo: poche righe visibili

	const numOrdini = 40
	res, _ := m.Update(OrdiniEntrataLoadedMsg{Ordini: nuoviOrdiniMock(numOrdini)})
	m = res.(Model)

	vis := m.righeVisibili()
	if vis >= numOrdini {
		t.Fatalf("Con winH=30 dovrebbero essere visibili poche righe, trovate %d", vis)
	}

	// PgDown avanza di una pagina intera
	res, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyPgDown})
	m = res.(Model)
	if m.scrollOffset != vis {
		t.Errorf("Dopo PgDown scrollOffset dovrebbe essere %d, ottenuto %d", vis, m.scrollOffset)
	}

	// Home torna all'inizio, End salta alla fine
	res, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyHome})
	m = res.(Model)
	if m.scrollOffset != 0 {
		t.Errorf("Dopo Home scrollOffset dovrebbe essere 0, ottenuto %d", m.scrollOffset)
	}
	res, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyEnd})
	m = res.(Model)
	maxOffset := len(m.listaRighe()) - vis
	if maxOffset < 0 {
		maxOffset = 0
	}
	if m.scrollOffset != maxOffset {
		t.Errorf("Dopo End scrollOffset dovrebbe essere %d, ottenuto %d", maxOffset, m.scrollOffset)
	}

	// L'indicatore "Ordini X–Y di Z" compare quando non tutto è visibile
	view := m.View().Content
	if !strings.Contains(view, "di 40") {
		t.Error("Il footer dovrebbe contenere l'indicatore 'Ordini X–Y di 40'")
	}
}

func TestOrdiniEntrata_ScorrimentoRotellina(t *testing.T) {
	client := &api.APIClient{}
	m := New(client)
	m.winH = 30
	res, _ := m.Update(OrdiniEntrataLoadedMsg{Ordini: nuoviOrdiniMock(40)})
	m = res.(Model)

	res, _ = m.Update(tea.MouseWheelMsg{Button: tea.MouseWheelDown})
	m = res.(Model)
	if m.scrollOffset != passoRotellina {
		t.Errorf("Un notch di rotellina in giù dovrebbe scorrere di %d, offset=%d", passoRotellina, m.scrollOffset)
	}

	res, _ = m.Update(tea.MouseWheelMsg{Button: tea.MouseWheelUp})
	m = res.(Model)
	if m.scrollOffset != 0 {
		t.Errorf("Rotellina in su dovrebbe tornare a 0, offset=%d", m.scrollOffset)
	}

	// Sforamento verso l'alto: resta a 0
	for i := 0; i < 5; i++ {
		res, _ = m.Update(tea.MouseWheelMsg{Button: tea.MouseWheelUp})
		m = res.(Model)
	}
	if m.scrollOffset != 0 {
		t.Errorf("Lo scorrimento non deve andare sotto zero, offset=%d", m.scrollOffset)
	}
}

func TestOrdiniEntrata_CursoreSegueFinestra(t *testing.T) {
	client := &api.APIClient{}
	m := New(client)
	m.winH = 26 // ~pochi ordini visibili
	res, _ := m.Update(OrdiniEntrataLoadedMsg{Ordini: nuoviOrdiniMock(40)})
	m = res.(Model)

	// Scendendo col cursore oltre la finestra, lo scroll segue il cursore
	passi := m.righeVisibili() + 3
	for i := 0; i < passi; i++ {
		res, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
		m = res.(Model)
	}
	if m.cursore != passi {
		t.Fatalf("cursore atteso %d, ottenuto %d", passi, m.cursore)
	}
	if m.cursore < m.scrollOffset || m.cursore >= m.scrollOffset+m.righeVisibili() {
		t.Errorf("Il cursore %d è fuori dalla finestra [%d, %d)",
			m.cursore, m.scrollOffset, m.scrollOffset+m.righeVisibili())
	}

	// L'ordine evidenziato deve comparire nella vista
	view := m.View().Content
	numeroCorrente := fmt.Sprintf("ENT-2026-%04d", passi+1)
	if !strings.Contains(view, numeroCorrente) {
		t.Errorf("L'ordine evidenziato %q dovrebbe essere visibile a schermo", numeroCorrente)
	}
}

func TestOrdiniEntrata_KPIAggregati(t *testing.T) {
	ordini := []api.OrdineEntrata{
		{
			ID: 1, Numero: "ENT-001", Fornitore: "Alfa", Data: time.Now(),
			Righe: []api.RigaOrdineEntrata{{Quantita: 4, PrezzoUnitario: 10}},
		},
		{
			ID: 2, Numero: "ENT-002", Fornitore: "Beta", Data: time.Now(),
			Righe: []api.RigaOrdineEntrata{
				{Quantita: 6, PrezzoUnitario: 5},
				{Quantita: 2, PrezzoUnitario: 20},
			},
		},
		{
			ID: 3, Numero: "ENT-003", Fornitore: "Alfa", Data: time.Now(),
			Righe: nil,
		},
	}

	fornitori, articoli, valore := statisticheEntrata(ordini)
	if fornitori != 2 {
		t.Errorf("Fornitori distinti attesi 2, ottenuti %d", fornitori)
	}
	if articoli != 12 {
		t.Errorf("Articoli attesi 12, ottenuti %d", articoli)
	}
	if valore != 110 {
		t.Errorf("Valore atteso 110, ottenuto %.2f", valore)
	}
}

// TestOrdiniEntrata_ComposizioneEntraNelloSchermo verifica la regressione degli
// "ordini invisibili": componendo la pagina esattamente come fa l'app
// (tab + corpo incorniciato, impaginati da tabnav.PlacePage) il risultato deve
// stare dentro le righe del terminale per ogni altezza ragionevole.
func TestOrdiniEntrata_ComposizioneEntraNelloSchermo(t *testing.T) {
	tabs := lipgloss.JoinHorizontal(lipgloss.Top,
		style.ActiveTab.Render("GIACENZE"), "  ", style.TabGap.Render("ORDINI IN ENTRATA"))

	for _, winH := range []int{28, 30, 40, 50} {
		m := New(nil)
		m.winW = 120

		res, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: winH})
		m = res.(Model)

		res, _ = m.Update(OrdiniEntrataLoadedMsg{Ordini: nuoviOrdiniMock(60)})
		m = res.(Model)

		composed := tabnav.PlacePage(120, winH,
			style.FinestraStile.Render(tabs+"\n\n"+m.View().Content))
		altezza := strings.Count(composed, "\n") + 1

		if altezza > winH {
			t.Errorf("winH=%d: la pagina composta occupa %d righe e fuoriesce dal terminale", winH, altezza)
		}

		view := m.View().Content
		if !strings.Contains(view, "di 60") {
			t.Errorf("winH=%d: atteso indicatore di scorrimento 'Ordini X–Y di 60'", winH)
		}
	}
}
