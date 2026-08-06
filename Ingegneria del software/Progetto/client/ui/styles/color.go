package style

import (
	"image/color"

	"charm.land/lipgloss/v2"
)

var (
	WindowBorderColors = []color.Color{
		lipgloss.Color("#41DCBD"),
		lipgloss.Color("#41AEDC"),
		lipgloss.Color("#4160DC"),
	}

	HeaderColor = lipgloss.Color("#16e7c4")

	TabActiveBorderColor    = lipgloss.Color("#5DC064")
	TabNonActiveBorderColor = lipgloss.Color("#5DB9C0")
	TabGapColor             = lipgloss.Color("#5DC095")

	TabTextColor       = lipgloss.Color("#F8F7F8")
	TabActiveTextColor = lipgloss.Color("#F8F7F8")

	LoginBorderColor = lipgloss.Color("#41AEDC")
	LoginTitleColor  = lipgloss.Color("#16e7c4")
	ErrorColor       = lipgloss.Color("#FF5555")
	SuccessColor     = lipgloss.Color("#5DC064")
)
