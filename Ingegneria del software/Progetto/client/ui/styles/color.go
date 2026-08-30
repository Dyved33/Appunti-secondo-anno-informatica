package style

import (
	"image/color"

	"charm.land/lipgloss/v2"
)

var (
	// WindowBorderColors è la palette di colori sfumati per il bordo della finestra.
	WindowBorderColors = []color.Color{
		lipgloss.Color("#41DCBD"),
		lipgloss.Color("#41AEDC"),
		lipgloss.Color("#4160DC"),
	}

	// HeaderColor è il colore del testo dell'intestazione.
	HeaderColor = lipgloss.Color("#16e7c4")

	// TabActiveBorderColor è il colore del bordo del tab attivo.
	TabActiveBorderColor = lipgloss.Color("#5DC064")
	// TabNonActiveBorderColor è il colore del bordo dei tab inattivi.
	TabNonActiveBorderColor = lipgloss.Color("#5DB9C0")
	// TabGapColor è il colore degli spazi tra i tab.
	TabGapColor = lipgloss.Color("#5DC095")

	// TabTextColor è il colore del testo dei tab inattivi.
	TabTextColor = lipgloss.Color("#F8F7F8")
	// TabActiveTextColor è il colore del testo del tab attivo.
	TabActiveTextColor = lipgloss.Color("#F8F7F8")

	// LoginBorderColor è il colore del bordo della schermata di login.
	LoginBorderColor = lipgloss.Color("#41AEDC")
	// LoginTitleColor è il colore del titolo della schermata di login.
	LoginTitleColor = lipgloss.Color("#16e7c4")
	// ErrorColor è il colore dei messaggi di errore.
	ErrorColor = lipgloss.Color("#FF5555")
	// SuccessColor è il colore dei messaggi di successo.
	SuccessColor = lipgloss.Color("#5DC064")
)
