package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/user/pspterm/config"
)

// EditorSubMode mirrors the model's editorSubMode for rendering.
const (
	EditorSubList    = 0
	EditorSubForm    = 1
	EditorSubConfirm = 2
)

// FormField mirrors model's formField constants.
const (
	FormFieldName  = 0
	FormFieldDesc  = 1
	FormFieldCat   = 2
	FormFieldType  = 3
	FormFieldValue = 4
)

// EditorView is a data-transfer object built by the model for UI rendering.
type EditorView struct {
	SubMode int

	// List view
	FlatCatIdxs  []int // catIdx for each flat item (parallel slices)
	FlatItemIdxs []int // itemIdx for each flat item
	ListIdx      int

	// Form view
	IsEditing  bool
	NameView   string // rendered textinput.View()
	DescView   string // rendered textinput.View()
	ValueView  string // rendered textinput.View()
	CatSel     int
	TypeSel    int
	FocusField int
	Suggestion string
	StatusMsg  string
	ItemTypes  []string

	Categories []config.Category
}

// RenderEditor renders the full editor screen.
func RenderEditor(ev EditorView, width, height int, s Styles) string {
	switch ev.SubMode {
	case EditorSubForm:
		return renderEditorForm(ev, width, height, s)
	case EditorSubConfirm:
		return renderEditorConfirm(ev, width, height, s)
	default:
		return renderEditorList(ev, width, height, s)
	}
}

// ── List view ────────────────────────────────────────────────────────────────

func renderEditorList(ev EditorView, width, height int, s Styles) string {
	accent := s.ItemTitle.GetForeground()
	dim := s.ItemNormal.GetForeground()

	header := lipgloss.NewStyle().
		Foreground(accent).
		Bold(true).
		Width(width).
		Align(lipgloss.Center)
	catLabel := lipgloss.NewStyle().
		Foreground(accent).
		Bold(true)
	selected := lipgloss.NewStyle().
		Foreground(accent).
		Bold(true)
	normal := lipgloss.NewStyle().
		Foreground(dim)
	typeTag := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#1a4060")).
		Faint(true)
	divStyle := s.Divider

	var lines []string
	lines = append(lines, header.Render("─── Item Manager ───"))
	lines = append(lines, "")

	if len(ev.FlatCatIdxs) == 0 {
		lines = append(lines, lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(
			normal.Render("(no items — press 'a' to add one)"),
		))
	} else {
		lastCat := -1
		for i, ci := range ev.FlatCatIdxs {
			ii := ev.FlatItemIdxs[i]
			cat := ev.Categories[ci]

			// Category header when it changes
			if ci != lastCat {
				if lastCat != -1 {
					lines = append(lines, "")
				}
				lines = append(lines, "  "+catLabel.Render(cat.Icon+"  "+strings.ToUpper(cat.Name)))
				lines = append(lines, "  "+divStyle.Render(strings.Repeat("─", len([]rune(cat.Name))+4)))
				lastCat = ci
			}

			item := cat.Items[ii]
			tag := typeTag.Render("[" + item.Type + "]")

			// Pad name + tag to fixed width so tags align
			nameWidth := width - 10
			if nameWidth < 10 {
				nameWidth = 10
			}
			namePart := item.Name
			if len([]rune(namePart)) > nameWidth-12 {
				namePart = string([]rune(namePart)[:nameWidth-15]) + "..."
			}
			row := fmt.Sprintf("%-*s %s", nameWidth-8, namePart, tag)

			if i == ev.ListIdx {
				lines = append(lines, "  "+selected.Render("▶ "+row))
			} else {
				lines = append(lines, "    "+normal.Render(row))
			}
		}
	}

	// Fill to height, then add footer
	helpBar := s.StatusBar.Render("a add  e edit  d delete  esc back")
	if ev.StatusMsg != "" {
		helpBar = s.ItemTitle.Render(ev.StatusMsg)
	}

	// Build final output: pad with blank lines then footer
	contentLines := len(lines)
	footerLines := 2 // blank + help
	for i := contentLines; i < height-footerLines; i++ {
		lines = append(lines, "")
	}
	lines = append(lines, "")
	lines = append(lines, helpBar)

	return strings.Join(lines, "\n")
}

// ── Form view ────────────────────────────────────────────────────────────────

func renderEditorForm(ev EditorView, width, height int, s Styles) string {
	accent := s.ItemTitle.GetForeground()

	label := lipgloss.NewStyle().
		Foreground(accent).
		Bold(true).
		Width(12)
	focused := lipgloss.NewStyle().
		Foreground(accent)
	unfocused := lipgloss.NewStyle().
		Foreground(s.ItemNormal.GetForeground())
	hint := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#1a4060")).
		Faint(true)

	title := "Add Item"
	if ev.IsEditing {
		title = "Edit Item"
	}
	header := lipgloss.NewStyle().
		Foreground(accent).
		Bold(true).
		Width(width).
		Align(lipgloss.Center)

	var lines []string
	lines = append(lines, header.Render("─── "+title+" ───"))
	lines = append(lines, "")

	// Name field
	nameLabel := label.Render("Name:")
	if ev.FocusField == FormFieldName {
		nameLabel = focused.Render(label.Render("Name:"))
	}
	lines = append(lines, "  "+nameLabel+"  "+ev.NameView)
	lines = append(lines, "")

	// Description field
	dLabel := label.Render("Desc:")
	if ev.FocusField == FormFieldDesc {
		dLabel = focused.Render(label.Render("Desc:"))
	}
	lines = append(lines, "  "+dLabel+"  "+ev.DescView)
	lines = append(lines, "")

	// Category picker
	catLabel := label.Render("Category:")
	if ev.FocusField == FormFieldCat {
		catLabel = focused.Render(label.Render("Category:"))
	}
	catPicker := renderPicker(ev.Categories, func(i int) string {
		return ev.Categories[i].Name
	}, ev.CatSel, ev.FocusField == FormFieldCat, focused, unfocused)
	lines = append(lines, "  "+catLabel+"  "+catPicker)
	lines = append(lines, "")

	// Type picker
	typeLabel := label.Render("Type:")
	if ev.FocusField == FormFieldType {
		typeLabel = focused.Render(label.Render("Type:"))
	}
	typePicker := renderStringPicker(ev.ItemTypes, ev.TypeSel, ev.FocusField == FormFieldType, focused, unfocused)
	lines = append(lines, "  "+typeLabel+"  "+typePicker)
	lines = append(lines, "")

	// Value field label changes by type
	valueLabel := "Command:"
	if len(ev.ItemTypes) > ev.TypeSel {
		switch ev.ItemTypes[ev.TypeSel] {
		case "directory":
			valueLabel = "Path:"
		case "url":
			valueLabel = "URL:"
		}
	}
	vl := label.Render(valueLabel)
	if ev.FocusField == FormFieldValue {
		vl = focused.Render(label.Render(valueLabel))
	}
	lines = append(lines, "  "+vl+"  "+ev.ValueView)

	// Auto-complete suggestion
	if ev.Suggestion != "" && ev.FocusField == FormFieldValue {
		lines = append(lines, "  "+strings.Repeat(" ", 14)+hint.Render("Tab to fill: "+ev.Suggestion))
	} else {
		lines = append(lines, "")
	}

	// Status / error
	statusLine := ""
	if ev.StatusMsg != "" {
		statusLine = s.ErrorMsg.Render(ev.StatusMsg)
	}

	// Help bar
	helpBar := s.StatusBar.Render("tab next field  enter save  esc cancel")

	// Pad and append footer
	contentLines := len(lines)
	footerLines := 3 // status + blank + help
	for i := contentLines; i < height-footerLines; i++ {
		lines = append(lines, "")
	}
	lines = append(lines, statusLine)
	lines = append(lines, "")
	lines = append(lines, helpBar)

	return strings.Join(lines, "\n")
}

// renderPicker renders a ← item1 · item2 · item3 → style picker for categories.
func renderPicker(cats []config.Category, label func(int) string, sel int, active bool, focused, unfocused lipgloss.Style) string {
	var parts []string
	for i := range cats {
		if i == sel {
			if active {
				parts = append(parts, focused.Bold(true).Render(label(i)))
			} else {
				parts = append(parts, unfocused.Bold(true).Render(label(i)))
			}
		} else {
			parts = append(parts, unfocused.Faint(true).Render(label(i)))
		}
	}
	arrow := unfocused.Faint(true).Render
	if active {
		arrow = focused.Render
	}
	return arrow("← ") + strings.Join(parts, unfocused.Faint(true).Render(" · ")) + arrow(" →")
}

// renderStringPicker renders a ← a · b · c → picker for a string slice.
func renderStringPicker(items []string, sel int, active bool, focused, unfocused lipgloss.Style) string {
	var parts []string
	for i, s := range items {
		if i == sel {
			if active {
				parts = append(parts, focused.Bold(true).Render(s))
			} else {
				parts = append(parts, unfocused.Bold(true).Render(s))
			}
		} else {
			parts = append(parts, unfocused.Faint(true).Render(s))
		}
	}
	arrow := unfocused.Faint(true).Render
	if active {
		arrow = focused.Render
	}
	return arrow("← ") + strings.Join(parts, unfocused.Faint(true).Render(" · ")) + arrow(" →")
}

// ── Confirm view ─────────────────────────────────────────────────────────────

func renderEditorConfirm(ev EditorView, width, height int, s Styles) string {
	if len(ev.FlatCatIdxs) == 0 {
		return ""
	}
	fi := ev.FlatCatIdxs[ev.ListIdx]
	fii := ev.FlatItemIdxs[ev.ListIdx]
	itemName := ev.Categories[fi].Items[fii].Name

	center := lipgloss.NewStyle().Width(width).Align(lipgloss.Center)
	msg := s.ErrorMsg.Render(fmt.Sprintf("Delete %q?  [y] yes   [n] no", itemName))

	lines := make([]string, height/2)
	lines = append(lines, center.Render(msg))
	return strings.Join(lines, "\n")
}
