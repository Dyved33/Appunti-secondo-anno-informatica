package ui

import (
	components "github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/components"
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/pages/login"
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/pages/mainmenu"
	style "github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/styles"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type FinestraPrincipale struct {
	paginaCorrente          string
	dimensioneFinestra      style.WindowSize
	schermataMenuPrincipale mainmenu.MenuPrincipale
	schermataDiLogin        login.Login
}

func NuovaFinestra() FinestraPrincipale {
	return FinestraPrincipale{
		paginaCorrente:          "menuprincipale",
		dimensioneFinestra:      style.NewWindowSize(0, 0),
		schermataMenuPrincipale: mainmenu.New(),
	}
}

func (f FinestraPrincipale) Init() tea.Cmd {
	return nil
}

func (f FinestraPrincipale) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	//CHIAMO IL CAMBIO DIMENSIONI PER LA FINESTRA
	case tea.WindowSizeMsg:
		f.dimensioneFinestra.ChangeSize(msg.Width, msg.Height)
		f.dimensioneFinestra.Border.SetSpacing(2, 2, 2, 2)
		f.schermataMenuPrincipale.AggiornaDimensioni(f.dimensioneFinestra)

		//INVIO I PULSANTI ALLA PAGINA CHE E ATTIVA
	case tea.KeyPressMsg:
		switch f.paginaCorrente {
		case "menuprincipale":
			nuovoMenu, cmd := f.schermataMenuPrincipale.Update(msg)
			f.schermataMenuPrincipale = nuovoMenu.(mainmenu.MenuPrincipale)
			return f, cmd

		case "login":
			nuovoLogin, cmd := f.schermataDiLogin.Update(msg)
			f.schermataDiLogin = nuovoLogin.(login.Login)
			return f, cmd
		}
	}

	return f, nil
}

func (f FinestraPrincipale) View() tea.View {

	//DISEGNO HEADER DEL PROGRAMMA
	preHeader := ""
	postHeader := ""
	bloccoHeader := lipgloss.JoinHorizontal(lipgloss.Left,
		preHeader,
		components.Logo,
		postHeader)
	bloccoHeader = lipgloss.JoinVertical(lipgloss.Center,
		bloccoHeader,
		components.Intestazione)

	var bloccoMenu string = ""
	//DISEGNO IL MENU
	switch f.paginaCorrente {
	case "menuprincipale":
		bloccoMenu = f.schermataMenuPrincipale.View().Content
	}

	//FOOTER MENU
	bloccoFooter := style.FooterStile.Render(components.Footer)

	//ORGANIZZO LA SCHERMATA
	schermataFinale := lipgloss.JoinVertical(lipgloss.Center,
		bloccoHeader, "",
		bloccoMenu, "")

	//CALCOLO DELLA SCHERMATA IN VERTICALE
	ingombroVerticale := lipgloss.Height(schermataFinale)
	spazioVerticaleRimanente := style.ReadHeight(f.dimensioneFinestra) - ingombroVerticale - 4
	spazioVuoto := lipgloss.NewStyle().Height(spazioVerticaleRimanente).Render("")
	schermataFinale = lipgloss.JoinVertical(lipgloss.Center,
		schermataFinale,
		spazioVuoto,
		bloccoFooter)

	style.FinestraStile = style.FinestraStile.Height(style.ReadHeight(f.dimensioneFinestra)).Width(style.ReadWidth(f.dimensioneFinestra))

	altscreen := tea.NewView(style.FinestraStile.Render(schermataFinale))
	altscreen.AltScreen = true

	return altscreen
}
