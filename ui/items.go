package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/user/pspterm/config"
)

// RenderItems renders the vertical sub-item list for a category.
// Returns a block of text centered in `width` columns.
func RenderItems(cat config.Category, selected int, width int, s Styles) string {
	center := lipgloss.NewStyle().Width(width).Align(lipgloss.Center)

	if len(cat.Items) == 0 {
		return center.Render("(no items)")
	}

	// Derive effective accent (per-category color overrides theme)
	titleStyle := s.ItemTitle
	selectedStyle := s.ItemSelected
	if cat.Color != "" {
		accent := lipgloss.Color(cat.Color)
		titleStyle = lipgloss.NewStyle().Inherit(s.ItemTitle).Foreground(accent)
		selectedStyle = lipgloss.NewStyle().Inherit(s.ItemSelected).Foreground(accent)
	}

	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#446680")).
		Faint(true).
		Italic(true)

	var lines []string

	// Category title
	title := cat.Icon + "  " + strings.ToUpper(cat.Name)
	lines = append(lines, center.Render(titleStyle.Render(title)))
	lines = append(lines, center.Render(s.Divider.Render(strings.Repeat("─", len([]rune(cat.Name))+4))))

	// Items
	for i, item := range cat.Items {
		if i == selected {
			lines = append(lines, center.Render(selectedStyle.Render("▶  "+item.Name)))
			if item.Description != "" {
				lines = append(lines, center.Render(descStyle.Render("  "+item.Description)))
			}
		} else {
			lines = append(lines, center.Render(s.ItemNormal.Render("   "+item.Name)))
		}
	}

	return strings.Join(lines, "\n")
}
