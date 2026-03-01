package model

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/harmonica"
	"github.com/charmbracelet/lipgloss"
	"github.com/user/pspterm/actions"
	"github.com/user/pspterm/config"
	"github.com/user/pspterm/ui"
)

const tickInterval = 16 * time.Millisecond // ~60 fps

const minWidth = 60

type appMode int

const (
	appModeXMB    appMode = iota
	appModeEditor appMode = iota
)

// Model is the Bubble Tea application model.
type Model struct {
	cfg          config.Config
	selectedCat  int
	selectedItem int

	// Spring animation state
	animPos float64
	animVel float64
	spring  harmonica.Spring

	now    time.Time
	width  int
	height int

	exitPath string // set on directory selection
	errMsg   string // transient error message

	keys   KeyMap
	help   help.Model
	styles ui.Styles

	appMode appMode
	editor  EditorState
}

// New creates a new Model from the given config.
func New(cfg config.Config) Model {
	return Model{
		cfg:    cfg,
		spring: harmonica.NewSpring(harmonica.FPS(60), 14.0, 0.55),
		now:    time.Now(),
		keys:   DefaultKeyMap(),
		help:   help.New(),
		styles: ui.NewStyles(cfg.Theme),
	}
}

// ExitPath returns the path to cd to on exit (empty if not a directory selection).
func (m Model) ExitPath() string { return m.exitPath }

// Init starts the tick loop.
func (m Model) Init() tea.Cmd {
	return tick()
}

func tick() tea.Cmd {
	return tea.Tick(tickInterval, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

// Update processes messages and returns the next model + command.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Always handle window resize and ticks regardless of mode
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.help.Width = msg.Width
		return m, nil

	case TickMsg:
		m.now = time.Time(msg)
		target := float64(m.selectedCat)
		m.animPos, m.animVel = m.spring.Update(m.animPos, m.animVel, target)
		return m, tick()
	}

	if m.appMode == appModeEditor {
		return m.updateEditor(msg)
	}

	switch msg := msg.(type) {
	case ExecDoneMsg:
		if msg.Err != nil {
			m.errMsg = fmt.Sprintf("exec error: %v", msg.Err)
		}
		return m, nil

	case EditConfigDoneMsg:
		if msg.Err != nil {
			m.errMsg = fmt.Sprintf("editor error: %v", msg.Err)
		} else {
			// Auto-reload so changes are visible immediately
			if cfg, err := config.Load(); err == nil {
				m.cfg = cfg
				m.styles = ui.NewStyles(cfg.Theme)
				clampSelections(&m)
			}
		}
		return m, nil

	case ConfigReloadedMsg:
		return m, nil

	case tea.KeyMsg:
		return m.handleKey(msg)
	}

	return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	cats := m.cfg.Categories

	switch {
	case key.Matches(msg, m.keys.Quit):
		return m, tea.Quit

	case key.Matches(msg, m.keys.Left):
		m.errMsg = ""
		if m.selectedCat > 0 {
			m.selectedCat--
			m.selectedItem = 0
		}

	case key.Matches(msg, m.keys.Right):
		m.errMsg = ""
		if m.selectedCat < len(cats)-1 {
			m.selectedCat++
			m.selectedItem = 0
		}

	case key.Matches(msg, m.keys.Up):
		m.errMsg = ""
		if m.selectedItem > 0 {
			m.selectedItem--
		}

	case key.Matches(msg, m.keys.Down):
		m.errMsg = ""
		if len(cats) > 0 && m.selectedItem < len(cats[m.selectedCat].Items)-1 {
			m.selectedItem++
		}

	case key.Matches(msg, m.keys.Select):
		m.errMsg = ""
		return m.executeSelected()

	case key.Matches(msg, m.keys.Reload):
		cfg, err := config.Load()
		if err != nil {
			m.errMsg = fmt.Sprintf("reload error: %v", err)
		} else {
			m.cfg = cfg
			m.styles = ui.NewStyles(cfg.Theme)
			clampSelections(&m)
		}
	}

	return m, nil
}

func (m Model) executeSelected() (tea.Model, tea.Cmd) {
	cats := m.cfg.Categories
	if len(cats) == 0 || len(cats[m.selectedCat].Items) == 0 {
		return m, nil
	}

	item := cats[m.selectedCat].Items[m.selectedItem]

	// Built-in: open in-TUI item manager
	if item.Type == "manager" {
		m.appMode = appModeEditor
		m.editor = newEditorState(m.cfg)
		return m, nil
	}

	// Built-in: open config.yaml in the user's preferred editor
	if item.Type == "editconfig" {
		editor := resolveEditor(m.cfg)
		c := exec.Command("sh", "-c", editor+" "+config.ConfigPath())
		return m, tea.ExecProcess(c, func(err error) tea.Msg {
			return EditConfigDoneMsg{Err: err}
		})
	}

	cmd, exitPath, err := actions.Execute(item, func(e error) tea.Msg {
		return ExecDoneMsg{Err: e}
	})
	if err != nil {
		m.errMsg = err.Error()
		return m, nil
	}
	if exitPath != "" {
		m.exitPath = exitPath
		return m, tea.Quit
	}
	if cmd != nil {
		return m, cmd
	}
	return m, nil
}

// View renders the full TUI.
func (m Model) View() string {
	if m.width < minWidth {
		return m.styles.ErrorMsg.Render(
			fmt.Sprintf("Terminal too narrow (%d cols). Please widen to at least %d.", m.width, minWidth),
		)
	}

	if m.appMode == appModeEditor {
		return ui.RenderEditor(m.editor.viewModel(m.cfg), m.width, m.height, m.styles)
	}

	cats := m.cfg.Categories
	if len(cats) == 0 {
		return renderEmptyState(m)
	}

	var sb strings.Builder

	// Top bar: clock right-aligned (PSP top-right corner)
	clockStr := m.now.Format("Mon 02 Jan  15:04")
	clockRendered := m.styles.Clock.Render(clockStr)
	clockPad := m.width - len(clockStr)
	if clockPad > 0 {
		sb.WriteString(strings.Repeat(" ", clockPad))
	}
	sb.WriteString(clockRendered)
	sb.WriteRune('\n')

	// Category bar
	sb.WriteString(ui.RenderCategoryBar(cats, m.animPos, m.width, m.styles))
	sb.WriteRune('\n')

	// Divider
	sb.WriteString(m.styles.Divider.Render(strings.Repeat("─", m.width)))
	sb.WriteRune('\n')

	// Item list area — fill remaining height
	// clock(1) + catbar(1) + divider(1) + divider(1) + statusbar(1) = 5
	const usedLines = 5
	itemAreaHeight := m.height - usedLines
	if itemAreaHeight < 3 {
		itemAreaHeight = 3
	}

	cat := cats[m.selectedCat]
	itemBlock := ui.RenderItems(cat, m.selectedItem, m.width, m.styles)
	itemLines := strings.Split(itemBlock, "\n")

	// Top-pad the items to vertically center the block
	topPad := (itemAreaHeight - len(itemLines)) / 2
	if topPad < 0 {
		topPad = 0
	}
	for i := 0; i < topPad; i++ {
		sb.WriteRune('\n')
	}
	sb.WriteString(itemBlock)
	// Fill remaining lines
	rendered := topPad + len(itemLines)
	for i := rendered; i < itemAreaHeight; i++ {
		sb.WriteRune('\n')
	}

	// Divider
	sb.WriteString(m.styles.Divider.Render(strings.Repeat("─", m.width)))
	sb.WriteRune('\n')

	// Status / error bar
	if m.errMsg != "" {
		sb.WriteString(m.styles.ErrorMsg.Render(m.errMsg))
	} else {
		sb.WriteString(m.styles.StatusBar.Render(renderHelp(m)))
	}

	return sb.String()
}

func renderEmptyState(m Model) string {
	var sb strings.Builder
	topPad := (m.height - 6) / 2
	if topPad < 0 {
		topPad = 0
	}
	for i := 0; i < topPad; i++ {
		sb.WriteRune('\n')
	}
	center := lipgloss.NewStyle().Width(m.width).Align(lipgloss.Center)
	sb.WriteString(center.Render(m.styles.CatSelected.Render("◈ pspterm ◈")))
	sb.WriteString("\n\n")
	sb.WriteString(center.Render(m.styles.ItemNormal.Render("No categories configured.")))
	sb.WriteString("\n\n")
	sb.WriteString(center.Render(m.styles.ItemNormal.Render("Edit  ~/.config/pspterm/config.yaml")))
	sb.WriteString("\n")
	sb.WriteString(center.Render(m.styles.ItemNormal.Render("Then press  r  to reload, or  q  to quit.")))
	return sb.String()
}

func renderHelp(m Model) string {
	bindings := []string{
		"←/→ category",
		"↑/↓ item",
		"enter select",
		"r reload",
		"q quit",
	}
	return strings.Join(bindings, "  ")
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// resolveEditor returns the editor binary to use for editconfig.
// Priority: config.Editor → $EDITOR → $VISUAL → first found in well-known list.
func resolveEditor(cfg config.Config) string {
	if cfg.Editor != "" {
		return cfg.Editor
	}
	if e := os.Getenv("EDITOR"); e != "" {
		return e
	}
	if e := os.Getenv("VISUAL"); e != "" {
		return e
	}
	for _, e := range []string{"nano", "vim", "vi", "micro", "emacs"} {
		if _, err := exec.LookPath(e); err == nil {
			return e
		}
	}
	return "vi" // POSIX fallback
}

// clampSelections ensures selectedCat/selectedItem are within bounds after a config reload.
func clampSelections(m *Model) {
	if m.selectedCat >= len(m.cfg.Categories) {
		m.selectedCat = max(0, len(m.cfg.Categories)-1)
	}
	if len(m.cfg.Categories) > 0 && m.selectedItem >= len(m.cfg.Categories[m.selectedCat].Items) {
		m.selectedItem = max(0, len(m.cfg.Categories[m.selectedCat].Items)-1)
	}
}
