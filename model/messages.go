package model

import "time"

// TickMsg is sent every ~33ms to drive animation.
type TickMsg time.Time

// ConfigReloadedMsg is sent after config is successfully reloaded.
type ConfigReloadedMsg struct{}

// ExecDoneMsg is sent when a tea.ExecProcess command finishes.
type ExecDoneMsg struct{ Err error }

// WhichResultMsg carries the result of a `which` lookup for command auto-complete.
type WhichResultMsg struct{ Path string }

// EditConfigDoneMsg is sent when the editor launched by editconfig exits.
type EditConfigDoneMsg struct{ Err error }
