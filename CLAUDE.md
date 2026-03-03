# pspterm — Claude Code Instructions

## What this is
PSP XMB terminal launcher replica written in Go with Bubble Tea. Navigates categories horizontally with spring-animated XMB-style bounce; items listed vertically beneath. Selecting a `directory` item exits and `cd`s via a shell wrapper eval trick.

## Build & run
```sh
go build -o pspterm .
./pspterm
```

## Shell wrapper (required for cd)
```sh
function psp() {
    local target
    target="$(./pspterm)"
    [ -d "$target" ] && cd "$target"
}
```

## Package layout
```
main.go               — entry point, /dev/tty setup, eval trick stdout
config/
  config.go           — Config/Theme/Category/Item structs + YAML load/save
  defaults.go         — LoadOrCreate, embedded default config YAML
  scan.go             — ScanApplications() for XDG .desktop discovery
model/
  model.go            — tea.Model (Init/Update/View), spring step, clock
  keys.go             — KeyMap (key.NewBinding)
  messages.go         — TickMsg, ConfigReloadedMsg, ExecDoneMsg
ui/
  styles.go           — lipgloss Styles struct + NewStyles(theme)
  xmb.go              — horizontal category bar renderer (spring-interpolated)
  items.go            — vertical item list renderer
actions/
  executor.go         — Execute(item, onExecDone) → directory/command/url
```

## Key API facts
- `key.Matches(msg, binding)` is a standalone function in bubbles v1.0.0 — not a method on binding.
- `harmonica.NewSpring(deltaTime, angFreq, dampingRatio)` — call `spring.Update(pos, vel, target)` each tick.
- `tea.ExecProcess(cmd, callback)` suspends the TUI for interactive subprocesses.
- `p.Run()` returns `(tea.Model, error)` — type-assert to `model.Model` to read `ExitPath()`.

## Import cycle rule
`actions` must NOT import `model`. The `ExecDoneMsg` callback is injected from `model.executeSelected()`.

## Animation tuning
Spring lives in `model.New()`: `harmonica.NewSpring(harmonica.FPS(60), 14.0, 0.55)`.
- Increase `angFreq` → faster snap.
- Decrease `dampingRatio` (keep > 0) → more bounce/overshoot.
- Tick rate is 60 fps (`16ms`).

## Config file
`~/.config/pspterm/config.yaml` — auto-created on first run from `config/defaults.go`.
Press `r` in the UI to hot-reload without restarting.

## Application scanning
Categories with `scan: true` are populated from XDG .desktop files at load time.
Searches `~/.local/share/applications`, `/usr/share/applications`, and `$XDG_DATA_DIRS`.
Scan categories are read-only in the item manager.

## No test suite yet
Add tests under `model/` or `config/` using standard `testing` package. No special build tags needed.
