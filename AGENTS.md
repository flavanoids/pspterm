# pspterm — Agent Guide

Guidance for AI agents (Claude Code or otherwise) working on this codebase.

## Understand before changing
Read the affected file fully before editing. The codebase is small — read all relevant files in parallel before proposing a plan.

## Critical constraints

### No import cycles
`actions` → `config` only. It must never import `model` or `ui`.
`model` → `config`, `ui`, `actions`. No circular deps.
`ui` → `config` only.

### Spring parameters live in one place
`model.New()` in `model/model.go`. Do not duplicate spring construction elsewhere.

### Config is user-owned
`~/.config/pspterm/config.yaml` belongs to the user. Code changes must not overwrite it. Only `config/defaults.go` is authoritative for the initial template.

### The eval trick is intentional
`main.go` writes the exit path to real stdout while the TUI draws to `/dev/tty`. Do not "fix" this by merging the output streams.

## Common tasks

### Adding a new item type
1. Add the type constant/string handling in `actions/executor.go` `Execute()` switch.
2. Add the field to `config.Item` struct in `config/config.go`.
3. Document the new type in `config/defaults.go` header comment.
4. Update `CLAUDE.md` item type list.

### Adding a new UI element
1. Add the lipgloss style to `ui/styles.go` `Styles` struct and `NewStyles()`.
2. Render in `model/model.go` `View()`.
3. Adjust `usedLines` constant in `View()` if the element takes a fixed line.

### Changing the color scheme
Edit `config/defaults.go` for the two config-driven colors (`accent_color`, `dim_color`).
Edit hardcoded colors directly in `ui/styles.go` `NewStyles()` for structural colors (divider, clock, far categories, normal items).

### Tuning animation feel
Only touch `harmonica.NewSpring(...)` args in `model.New()` and `tickInterval` in `model/model.go`.
- `angFreq` 10–18: sweet spot for XMB feel.
- `dampingRatio` 0.45–0.7: below 1.0 = underdamped (bouncy). Do not go below 0.4 or it oscillates too long.

## What agents should NOT do
- Add `wave_enabled` or wave rendering back — removed intentionally.
- Add global state or `init()` functions.
- Use `os.Exit` outside of `main.go`.
- Skip reading a file before editing it.
- Create new packages beyond the existing five (`main`, `config`, `model`, `ui`, `actions`).
- Auto-commit or auto-push changes.
