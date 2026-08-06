package style

import (
	"charm.land/lipgloss/v2"
)

var (
	// STILE DEI TAB
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

	// STILE DELLA FINESTRA
	FinestraStile = lipgloss.NewStyle().
			Border(windowBorder, true).
			BorderForegroundBlend(WindowBorderColors...).
			Bold(true)

	// STILE DEL HEADER
	HeaderStile = lipgloss.NewStyle().
			Foreground(HeaderColor).
			Width(90)

	// STILE DEL FOOTER
	FooterStile = lipgloss.NewStyle().
			Bold(true)

	// STILI PER LA SCHERMATA DI LOGIN
	LoginBoxStile = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder(), true).
			BorderForeground(LoginBorderColor).
			Padding(1, 2).
			Align(lipgloss.Center)

	LoginTitleStile = lipgloss.NewStyle().
			Foreground(LoginTitleColor).
			Bold(true).
			MarginBottom(1)

	LoginErrorStile = lipgloss.NewStyle().
			Foreground(ErrorColor).
			Bold(true).
			MarginTop(1)

	LoginSuccessStile = lipgloss.NewStyle().
			Foreground(SuccessColor).
			Bold(true).
			MarginTop(1)
)
