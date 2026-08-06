package mainmenu

import (
	style "github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/styles"

	tea "charm.land/bubbletea/v2"
	lip "charm.land/lipgloss/v2"
)

type MenuPrincipale struct {
	menu               []string
	scelta             int
	larghezzaTotale    int
	larghezzaContenuto int
}

func New() MenuPrincipale {
	return MenuPrincipale{
		menu: []string{"LOGIN", "MAGAZZINO", "CONTABILITA", "AMMINISTRAZIONE"},
	}
}

func (m MenuPrincipale) Init() tea.Cmd {
	return nil
}

func (m MenuPrincipale) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyPressMsg:

		switch msg.String() {

		//COMANDI ALLA PRESSIONE DEI TASTI
		case "A":
			if m.scelta > 0 {
				m.scelta--
			} else {
				m.scelta = len(m.menu)
			}

		case "D":
			if m.scelta < len(m.menu)-1 {
				m.scelta++
			} else {
				m.scelta = 0
			}

		case "esc":
			return m, tea.Quit

		case "enter":

		}
	}
	return m, nil
}

func (m MenuPrincipale) View() tea.View {

	//VARIABILE DI CONTENUTO DELLE PAGINE E UNA STRINGA
	var contenutoPagina []string

	//LOGICA DELLE TAB ATTIVA/DISATTIVA
	for elemento, scelto := range m.menu {

		if m.scelta == elemento {
			contenutoPagina = append(contenutoPagina, style.ActiveTab.Render(scelto))

		} else {
			contenutoPagina = append(contenutoPagina, style.Tab.Render(scelto))
		}
	}

	//METTE IN FILA TUTTE LE TAB
	pagina := lip.JoinHorizontal(lip.Bottom, contenutoPagina...)

	//CALCOLA L'INGOMBRO CHE OCCUPANO LE TAB
	ingombro := lip.Width(pagina)
	spazioRimanente := m.larghezzaContenuto - ingombro
	testoRimanente := style.TabGap.Width(spazioRimanente).Render("")
	paginaFinale := lip.JoinHorizontal(lip.Bottom, pagina, testoRimanente)
	return tea.NewView(paginaFinale)
}

// CALCOLA LE DIMENSIONI SENZA IL BORDO ESTERNO
func (m *MenuPrincipale) AggiornaDimensioni(size style.WindowSize) {

	m.larghezzaTotale = style.ReadWidth(size)
	m.larghezzaContenuto, _ = style.GetInnerSpace(size)
}
