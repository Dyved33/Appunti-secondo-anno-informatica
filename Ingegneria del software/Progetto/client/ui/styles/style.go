package style

import (
	"charm.land/lipgloss/v2"
)

// RighePerPagina è il numero massimo di righe visibili in una tabella paginata.
const RighePerPagina = 20

// PaginaCalcola aggiorna l'offset di paginazione in base alla posizione del cursore.
func PaginaCalcola(cursore, offset, totale int) int {
	if totale == 0 {
		return 0
	}
	// Scorri verso il basso
	if cursore >= offset+RighePerPagina {
		return cursore - RighePerPagina + 1
	}
	// Scorri verso l'alto
	if cursore < offset {
		return cursore
	}
	return offset
}

var (
	Tab = lipgloss.NewStyle().
		Border(tabBorder, true).
		BorderForeground(TabGapColor).
		Foreground(TabTextColor)

	ActiveTab = Tab.Border(activeTabBorder, true).
			BorderForeground(TabActiveBorderColor).
			Bold(true)

	TabGap = Tab.
		BorderTop(false).
		BorderLeft(false).
		BorderRight(false).
		BorderForeground(TabNonActiveBorderColor)

	FinestraStile = lipgloss.NewStyle().
			Border(windowBorder, true).
			BorderForegroundBlend(WindowBorderColors...).
			Bold(true)

	FinestraInternaStile = lipgloss.NewStyle().
				Bold(true)

	HeaderStile = lipgloss.NewStyle().
			Foreground(HeaderColor).
			Width(90)

	FreeHeaderStile = lipgloss.NewStyle().
			Foreground(HeaderColor)

	FooterStile = lipgloss.NewStyle().
			Bold(true)

	TitleStile = lipgloss.NewStyle().
			Foreground(LoginTitleColor).
			Bold(true).
			MarginBottom(1)

	SubTitleStile = lipgloss.NewStyle().
			Foreground(LoginTitleColor).
			Bold(true).
			MarginBottom(1)

	ErrorStile = lipgloss.NewStyle().
			Foreground(ErrorColor).
			Bold(true)

	SuccessStile = lipgloss.NewStyle().
			Foreground(SuccessColor).
			Bold(true)

	InfoStile = lipgloss.NewStyle().
			Bold(true)

	FocusStile = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("229")).
			Background(lipgloss.Color("57"))

	// STILE DEI BOX PER LE SCHEDE
	BoxStile = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder(), true).
			BorderForeground(LoginBorderColor).
			Padding(0, 2).
			Width(38).
			Height(4)

	// BOX PER LA SCHERMATA DI LOGIN
	LoginBoxStile = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder(), true).
			BorderForeground(LoginBorderColor).
			Padding(1, 2).
			Align(lipgloss.Center)
)
