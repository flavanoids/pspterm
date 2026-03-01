package model

import (
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/user/pspterm/config"
	"github.com/user/pspterm/ui"
)

// viewModel builds a ui.EditorView from the current editor state + config.
func (e *EditorState) viewModel(cfg config.Config) ui.EditorView {
	ev := ui.EditorView{
		SubMode:    int(e.subMode),
		ListIdx:    e.listIdx,
		IsEditing:  e.isEditing,
		NameView:   e.nameInput.View(),
		DescView:   e.descInput.View(),
		ValueView:  e.valueInput.View(),
		CatSel:     e.catSel,
		TypeSel:    e.typeSel,
		FocusField: int(e.focusField),
		Suggestion: e.suggestion,
		StatusMsg:  e.statusMsg,
		ItemTypes:  itemTypes,
		Categories: cfg.Categories,
	}
	for _, fi := range e.flatItems {
		ev.FlatCatIdxs = append(ev.FlatCatIdxs, fi.catIdx)
		ev.FlatItemIdxs = append(ev.FlatItemIdxs, fi.itemIdx)
	}
	return ev
}

type editorSubMode int

const (
	editorList    editorSubMode = iota
	editorForm    editorSubMode = iota
	editorConfirm editorSubMode = iota
)

type formField int

const (
	fieldName  formField = iota
	fieldDesc  formField = iota
	fieldCat   formField = iota
	fieldType  formField = iota
	fieldValue formField = iota
	fieldCount formField = iota // sentinel — number of fields
)

var itemTypes = []string{"command", "directory", "url"}

// flatItem maps a list position to its category + item indices.
type flatItem struct {
	catIdx  int
	itemIdx int
}

// EditorState holds all state for the in-TUI item manager.
type EditorState struct {
	subMode editorSubMode

	// List view
	flatItems []flatItem
	listIdx   int

	// Form view
	isEditing   bool
	editCatIdx  int
	editItemIdx int
	nameInput   textinput.Model
	descInput   textinput.Model
	valueInput  textinput.Model
	catSel      int // selected category index
	typeSel     int // index into itemTypes
	focusField  formField
	suggestion  string // path from `which`
	statusMsg   string
}

func newEditorState(cfg config.Config) EditorState {
	nameIn := textinput.New()
	nameIn.Placeholder = "e.g. htop"
	nameIn.CharLimit = 64

	descIn := textinput.New()
	descIn.Placeholder = "optional description"
	descIn.CharLimit = 128

	valIn := textinput.New()
	valIn.Placeholder = "command, path, or url"
	valIn.CharLimit = 256

	e := EditorState{
		subMode:    editorList,
		nameInput:  nameIn,
		descInput:  descIn,
		valueInput: valIn,
	}
	e.buildFlatItems(cfg)
	return e
}

func (e *EditorState) buildFlatItems(cfg config.Config) {
	e.flatItems = nil
	for ci, cat := range cfg.Categories {
		for ii := range cat.Items {
			e.flatItems = append(e.flatItems, flatItem{ci, ii})
		}
	}
}

func (e *EditorState) startAdd(cfg config.Config, defaultCatIdx int) {
	e.subMode = editorForm
	e.isEditing = false
	e.catSel = defaultCatIdx
	e.typeSel = 0
	e.focusField = fieldName
	e.nameInput.Reset()
	e.nameInput.Focus()
	e.descInput.Reset()
	e.descInput.Blur()
	e.valueInput.Reset()
	e.valueInput.Blur()
	e.suggestion = ""
	e.statusMsg = ""
}

func (e *EditorState) startEdit(cfg config.Config) {
	if len(e.flatItems) == 0 {
		return
	}
	fi := e.flatItems[e.listIdx]
	item := cfg.Categories[fi.catIdx].Items[fi.itemIdx]

	e.subMode = editorForm
	e.isEditing = true
	e.editCatIdx = fi.catIdx
	e.editItemIdx = fi.itemIdx
	e.catSel = fi.catIdx
	e.focusField = fieldName
	e.suggestion = ""
	e.statusMsg = ""

	e.typeSel = 0
	for i, t := range itemTypes {
		if t == item.Type {
			e.typeSel = i
			break
		}
	}

	val := item.Command
	if item.Type == "directory" {
		val = item.Path
	} else if item.Type == "url" {
		val = item.URL
	}

	e.nameInput.Reset()
	e.nameInput.SetValue(item.Name)
	e.nameInput.Focus()
	e.descInput.Reset()
	e.descInput.SetValue(item.Description)
	e.descInput.Blur()
	e.valueInput.Reset()
	e.valueInput.SetValue(val)
	e.valueInput.Blur()
}

func (e *EditorState) setFocus() {
	e.nameInput.Blur()
	e.descInput.Blur()
	e.valueInput.Blur()
	switch e.focusField {
	case fieldName:
		e.nameInput.Focus()
	case fieldDesc:
		e.descInput.Focus()
	case fieldValue:
		e.valueInput.Focus()
	}
}

func (e *EditorState) formItem() config.Item {
	name := strings.TrimSpace(e.nameInput.Value())
	desc := strings.TrimSpace(e.descInput.Value())
	val := strings.TrimSpace(e.valueInput.Value())
	t := itemTypes[e.typeSel]
	item := config.Item{Name: name, Description: desc, Type: t}
	switch t {
	case "command":
		item.Command = val
	case "directory":
		item.Path = val
	case "url":
		item.URL = val
	}
	return item
}

// whichCmd runs `which <name>` asynchronously and returns a WhichResultMsg.
func whichCmd(name string) tea.Cmd {
	return func() tea.Msg {
		name = strings.TrimSpace(name)
		if name == "" || strings.ContainsAny(name, " /") {
			return WhichResultMsg{}
		}
		out, err := exec.Command("which", name).Output()
		if err != nil {
			return WhichResultMsg{}
		}
		return WhichResultMsg{Path: strings.TrimSpace(string(out))}
	}
}

// updateEditor handles all messages when the app is in editor mode.
func (m Model) updateEditor(msg tea.Msg) (tea.Model, tea.Cmd) {
	e := &m.editor

	// WhichResultMsg can arrive in any sub-mode
	if wr, ok := msg.(WhichResultMsg); ok {
		e.suggestion = wr.Path
		return m, nil
	}

	keyMsg, isKey := msg.(tea.KeyMsg)

	// Esc is universal: back up one level
	if isKey && keyMsg.String() == "esc" {
		switch e.subMode {
		case editorList:
			m.appMode = appModeXMB
		case editorForm, editorConfirm:
			e.subMode = editorList
			e.statusMsg = ""
		}
		return m, nil
	}

	switch e.subMode {
	case editorList:
		return m.updateEditorList(msg)
	case editorForm:
		return m.updateEditorForm(msg)
	case editorConfirm:
		return m.updateEditorConfirm(msg)
	}

	return m, nil
}

func (m Model) updateEditorList(msg tea.Msg) (tea.Model, tea.Cmd) {
	e := &m.editor
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	switch keyMsg.String() {
	case "up", "k":
		if e.listIdx > 0 {
			e.listIdx--
		}
	case "down", "j":
		if e.listIdx < len(e.flatItems)-1 {
			e.listIdx++
		}
	case "a":
		// Default new item to current category if possible
		defCat := 0
		if len(e.flatItems) > 0 {
			defCat = e.flatItems[e.listIdx].catIdx
		}
		e.startAdd(m.cfg, defCat)
	case "e", "enter":
		if len(e.flatItems) > 0 {
			e.startEdit(m.cfg)
		}
	case "d":
		if len(e.flatItems) > 0 {
			e.subMode = editorConfirm
		}
	}
	return m, nil
}

func (m Model) updateEditorForm(msg tea.Msg) (tea.Model, tea.Cmd) {
	e := &m.editor

	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "tab":
			// On value field: Tab auto-completes suggestion if available, else advances
			if e.focusField == fieldValue && e.suggestion != "" {
				e.valueInput.SetValue(e.suggestion)
				e.suggestion = ""
				return m, nil
			}
			e.focusField = (e.focusField + 1) % fieldCount
			e.setFocus()
			if e.focusField == fieldValue && itemTypes[e.typeSel] == "command" {
				return m, whichCmd(e.nameInput.Value())
			}
			return m, nil

		case "shift+tab":
			e.focusField = (e.focusField + fieldCount - 1) % fieldCount
			e.setFocus()
			return m, nil

		case "enter":
			if e.focusField == fieldValue {
				return m.saveEditorForm()
			}
			// Advance field
			e.focusField = (e.focusField + 1) % fieldCount
			e.setFocus()
			if e.focusField == fieldValue && itemTypes[e.typeSel] == "command" {
				return m, whichCmd(e.nameInput.Value())
			}
			return m, nil

		case "left":
			switch e.focusField {
			case fieldCat:
				if e.catSel > 0 {
					e.catSel--
				}
				return m, nil
			case fieldType:
				if e.typeSel > 0 {
					e.typeSel--
				}
				e.suggestion = ""
				return m, nil
			}

		case "right":
			switch e.focusField {
			case fieldCat:
				if e.catSel < len(m.cfg.Categories)-1 {
					e.catSel++
				}
				return m, nil
			case fieldType:
				if e.typeSel < len(itemTypes)-1 {
					e.typeSel++
				}
				e.suggestion = ""
				return m, nil
			}
		}
	}

	// Forward to focused textinput
	var cmd tea.Cmd
	switch e.focusField {
	case fieldName:
		e.nameInput, cmd = e.nameInput.Update(msg)
	case fieldDesc:
		e.descInput, cmd = e.descInput.Update(msg)
	case fieldValue:
		e.valueInput, cmd = e.valueInput.Update(msg)
	}
	return m, cmd
}

func (m Model) updateEditorConfirm(msg tea.Msg) (tea.Model, tea.Cmd) {
	e := &m.editor
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	switch keyMsg.String() {
	case "y", "enter":
		return m.deleteSelectedItem()
	case "n":
		e.subMode = editorList
	}
	return m, nil
}

func (m Model) saveEditorForm() (tea.Model, tea.Cmd) {
	e := &m.editor
	item := e.formItem()

	if item.Name == "" {
		e.statusMsg = "Name cannot be empty"
		return m, nil
	}

	cfg := m.cfg

	if e.isEditing {
		// Remove from old location
		oldCat := &cfg.Categories[e.editCatIdx]
		oldCat.Items = append(oldCat.Items[:e.editItemIdx], oldCat.Items[e.editItemIdx+1:]...)
		// Add to (possibly different) target category
		cfg.Categories[e.catSel].Items = append(cfg.Categories[e.catSel].Items, item)
	} else {
		cfg.Categories[e.catSel].Items = append(cfg.Categories[e.catSel].Items, item)
	}

	if err := config.Save(cfg); err != nil {
		e.statusMsg = "Save failed: " + err.Error()
		return m, nil
	}

	m.cfg = cfg
	m.styles = ui.NewStyles(cfg.Theme)
	e.buildFlatItems(cfg)
	e.subMode = editorList
	e.statusMsg = "Saved!"

	if e.listIdx >= len(e.flatItems) {
		e.listIdx = max(0, len(e.flatItems)-1)
	}
	return m, nil
}

func (m Model) deleteSelectedItem() (tea.Model, tea.Cmd) {
	e := &m.editor
	if len(e.flatItems) == 0 {
		e.subMode = editorList
		return m, nil
	}

	fi := e.flatItems[e.listIdx]
	cfg := m.cfg
	cat := &cfg.Categories[fi.catIdx]
	cat.Items = append(cat.Items[:fi.itemIdx], cat.Items[fi.itemIdx+1:]...)

	if err := config.Save(cfg); err != nil {
		e.statusMsg = "Delete failed: " + err.Error()
		e.subMode = editorList
		return m, nil
	}

	m.cfg = cfg
	m.styles = ui.NewStyles(cfg.Theme)
	e.buildFlatItems(cfg)
	e.subMode = editorList
	e.statusMsg = "Deleted."

	if e.listIdx >= len(e.flatItems) {
		e.listIdx = max(0, len(e.flatItems)-1)
	}
	return m, nil
}
