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

// Tab è lo stile di base di un tab inattivo nella barra di navigazione.
var (
	Tab = lipgloss.NewStyle().
		Border(tabBorder, true).
		BorderForeground(TabGapColor).
		Foreground(TabTextColor)

	// ActiveTab è lo stile del tab attualmente selezionato, con bordo evidenziato e testo in grassetto.
	ActiveTab = Tab.Border(activeTabBorder, true).
			BorderForeground(TabActiveBorderColor).
			Bold(true)

	// TabGap è lo stile dello spazio di separazione tra i tab.
	TabGap = Tab.
		BorderTop(false).
		BorderLeft(false).
		BorderRight(false).
		BorderForeground(TabNonActiveBorderColor)

	// FinestraStile è lo stile della cornice che delimita il contenuto principale della finestra.
	FinestraStile = lipgloss.NewStyle().
			Border(windowBorder, true).
			BorderForegroundBlend(WindowBorderColors...).
			Bold(true)

	// FinestraInternaStile è lo stile del contenuto interno alla finestra principale.
	FinestraInternaStile = lipgloss.NewStyle().
				Bold(true)

	// HeaderStile è lo stile dell'intestazione con larghezza fissa.
	HeaderStile = lipgloss.NewStyle().
			Foreground(HeaderColor).
			Width(90)

	// FreeHeaderStile è lo stile dell'intestazione senza vincolo di larghezza.
	FreeHeaderStile = lipgloss.NewStyle().
			Foreground(HeaderColor)

	// FooterStile è lo stile del piè di pagina con testo in grassetto.
	FooterStile = lipgloss.NewStyle().
			Bold(true)

	// TitleStile è lo stile dei titoli, con colore dedicato e margine inferiore.
	TitleStile = lipgloss.NewStyle().
			Foreground(LoginTitleColor).
			Bold(true).
			MarginBottom(1)

	// SubTitleStile è lo stile dei sottotitoli, con colore dedicato e margine inferiore.
	SubTitleStile = lipgloss.NewStyle().
			Foreground(LoginTitleColor).
			Bold(true).
			MarginBottom(1)

	// ErrorStile è lo stile dei messaggi di errore.
	ErrorStile = lipgloss.NewStyle().
			Foreground(ErrorColor).
			Bold(true)

	// SuccessStile è lo stile dei messaggi di successo.
	SuccessStile = lipgloss.NewStyle().
			Foreground(SuccessColor).
			Bold(true)

	// InfoStile è lo stile dei messaggi informativi.
	InfoStile = lipgloss.NewStyle().
			Bold(true)

	// FocusStile è lo stile che evidenzia l'elemento correntemente focalizzato.
	FocusStile = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("229")).
			Background(lipgloss.Color("57"))

	// BoxStile è lo stile dei box per le schede.
	BoxStile = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder(), true).
			BorderForeground(LoginBorderColor).
			Padding(0, 2).
			Width(38).
			Height(4)

	// LoginBoxStile è lo stile del box contenente la schermata di login.
	LoginBoxStile = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder(), true).
			BorderForeground(LoginBorderColor).
			Padding(1, 2).
			Align(lipgloss.Center)
)
