package ui

import (
	"math"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/user/pspterm/config"
)

const slotWidth = 18 // horizontal spacing between category slots

// RenderCategoryBar renders the XMB horizontal category bar.
// animPos is the spring-interpolated position (float64 of selected index).
func RenderCategoryBar(cats []config.Category, animPos float64, width int, s Styles) string {
	type slot struct {
		x    int
		text string
	}

	var slots []slot

	for i, cat := range cats {
		relPos := float64(i) - animPos
		label := cat.Icon + " " + cat.Name

		var styled string
		absRel := math.Abs(relPos)
		switch {
		case absRel < 0.5:
			accent := s.CatSelected.GetForeground()
			if cat.Color != "" {
				accent = lipgloss.Color(cat.Color)
			}
			styled = lipgloss.NewStyle().Inherit(s.CatSelected).Foreground(accent).Render(strings.ToUpper(label))
		case absRel < 1.5:
			styled = s.CatAdjacent.Render(label)
		case absRel < 2.5:
			styled = s.CatFar.Render(label)
		default:
			continue // hidden
		}

		labelW := ansi.StringWidth(styled)
		x := width/2 + int(math.Round(relPos*slotWidth)) - labelW/2
		slots = append(slots, slot{x: x, text: styled})
	}

	// Sort left-to-right
	sort.Slice(slots, func(i, j int) bool { return slots[i].x < slots[j].x })

	// Build the line by walking character positions
	var sb strings.Builder
	cursor := 0
	for _, sl := range slots {
		if sl.x < cursor {
			// overlap — skip to avoid garbled output
			continue
		}
		if sl.x > cursor {
			gap := sl.x - cursor
			if gap > 0 {
				sb.WriteString(strings.Repeat(" ", gap))
				cursor += gap
			}
		}
		sb.WriteString(sl.text)
		cursor += ansi.StringWidth(sl.text)
	}

	// Trim or pad to width
	result := sb.String()
	resultW := ansi.StringWidth(result)
	if resultW < width {
		result += strings.Repeat(" ", width-resultW)
	}
	return result
}

// RenderDotIndicator renders a centered row of dots — one per category.
// The active category is a filled dot (•); others are hollow (·).
func RenderDotIndicator(cats []config.Category, selected int, width int, s Styles) string {
	if len(cats) <= 1 {
		return strings.Repeat(" ", width)
	}
	var parts []string
	for i := range cats {
		if i == selected {
			parts = append(parts, s.DotActive.Render("•"))
		} else {
			parts = append(parts, s.DotInactive.Render("·"))
		}
	}
	line := strings.Join(parts, " ")
	return lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(line)
}
