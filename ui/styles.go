package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/user/pspterm/config"
)

// Styles holds all lipgloss styles derived from config.Theme.
type Styles struct {
	// Category bar
	CatSelected lipgloss.Style
	CatAdjacent lipgloss.Style
	CatFar      lipgloss.Style
	CatHidden   lipgloss.Style

	// Item list
	ItemSelected lipgloss.Style
	ItemNormal   lipgloss.Style
	ItemTitle    lipgloss.Style

	// Layout
	Divider   lipgloss.Style
	Clock     lipgloss.Style
	StatusBar lipgloss.Style
	ErrorMsg  lipgloss.Style
}

// NewStyles creates Styles from a Theme.
func NewStyles(t config.Theme) Styles {
	accent := lipgloss.Color(t.AccentColor)
	dim := lipgloss.Color(t.DimColor)

	return Styles{
		// Selected category: bright PSP blue, bold — pops against the dark
		CatSelected: lipgloss.NewStyle().
			Foreground(accent).
			Bold(true).
			Padding(0, 2),

		// Adjacent: very dark navy — present but receding
		CatAdjacent: lipgloss.NewStyle().
			Foreground(dim).
			Padding(0, 1),

		// Far: barely visible against black terminal bg
		CatFar: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#071828")).
			Padding(0, 1),

		CatHidden: lipgloss.NewStyle(),

		// Selected item: accent blue, bold
		ItemSelected: lipgloss.NewStyle().
			Foreground(accent).
			Bold(true).
			PaddingLeft(2),

		// Normal items: very dark muted blue — barely readable, PSP-dim
		ItemNormal: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#1e4a66")).
			PaddingLeft(4),

		ItemTitle: lipgloss.NewStyle().
			Foreground(accent).
			Bold(true),

		// Divider: near-invisible dark seam
		Divider: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#091e2e")),

		// Clock: cool blue-white, like the PSP system font
		Clock: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#8ecce6")).
			Bold(true),

		// Status bar: barely there
		StatusBar: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#1a4060")).
			Faint(true),

		ErrorMsg: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ff4444")).
			Bold(true),
	}
}
