package model

import "github.com/charmbracelet/bubbles/key"

// KeyMap defines all key bindings for pspterm.
type KeyMap struct {
	Left   key.Binding
	Right  key.Binding
	Up     key.Binding
	Down   key.Binding
	Select key.Binding
	Reload key.Binding
	Quit   key.Binding
}

// DefaultKeyMap returns the standard PSP-style key map.
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Left: key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp("←/h", "prev category"),
		),
		Right: key.NewBinding(
			key.WithKeys("right", "l"),
			key.WithHelp("→/l", "next category"),
		),
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "item up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "item down"),
		),
		Select: key.NewBinding(
			key.WithKeys("enter", " "),
			key.WithHelp("enter", "select"),
		),
		Reload: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "reload config"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "esc"),
			key.WithHelp("q/esc", "quit"),
		),
	}
}

// ShortHelp returns bindings for the compact help view.
func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Left, k.Right, k.Up, k.Down, k.Select, k.Reload, k.Quit}
}

// FullHelp returns grouped bindings for the full help view.
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Left, k.Right},
		{k.Up, k.Down},
		{k.Select, k.Reload, k.Quit},
	}
}
