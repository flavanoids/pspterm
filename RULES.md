# pspterm — Rules

Hard constraints for this project. These are non-negotiable.

---

## Code rules

### R1 — No import cycles
Package dependency order (arrows = "may import"):
```
main → model, config
model → ui, actions, config
ui → config
actions → config
```
Any other import direction is forbidden.

### R2 — No global mutable state
No `var` at package level that is mutated after init. Config, styles, and model state live inside `model.Model`. Spring state lives inside `model.Model`. No singletons.

### R3 — No `os.Exit` outside main
Only `main.go` may call `os.Exit`. All other packages return errors up the call stack.

### R4 — Config struct is the source of truth
Do not hardcode category names, item types, or paths in Go code. Everything user-configurable belongs in `config.Config` and the YAML file.

### R5 — `tea.ExecProcess` for interactive commands
Commands that need a real TTY (shell, vim, htop, etc.) must use `tea.ExecProcess`, not `exec.Command` with `Start()`. Background commands (xdg-open, etc.) may use `Start()` in a goroutine.

### R6 — Path expansion in one place
`~` and `$ENV_VAR` expansion for file paths happens only in `actions.expandPath()`. Do not expand paths elsewhere.

### R7 — One spring instance per model
The harmonica `Spring` is created once in `model.New()`. Do not create additional spring instances at animation time.

### R8 — Clock reads from TickMsg time
`model.now` must be set from `time.Time(msg)` in the `TickMsg` handler, not from `time.Now()` calls scattered elsewhere. This keeps time consistent per frame.

---

## UX rules

### U1 — Navigation resets item selection
Moving left or right across categories always resets `selectedItem` to 0. This matches PSP XMB behavior.

### U2 — Errors are transient
`model.errMsg` is cleared on the next navigation keypress. Errors do not persist across interactions.

### U3 — Config reload preserves valid selection
After `r` (reload), clamp `selectedCat` and `selectedItem` to valid ranges. Never leave them out of bounds.

### U4 — Minimum width enforced
If terminal width < 60 columns, show an error message instead of a garbled layout. The constant is `minWidth` in `model/model.go`.

### U5 — No waves
Wave background rendering is removed. Do not reintroduce it. The background is the terminal's own color (black).

---

## Contribution rules

### C1 — Read before editing
Always read the full target file before making changes. Never edit based on partial knowledge.

### C2 — Build must pass
`go build ./...` and `go vet ./...` must pass after every change. Do not leave the project in a broken state.

### C3 — No auto-commit
Do not commit or push changes without explicit user instruction.

### C4 — No new packages
The five existing packages (`main`, `config`, `model`, `ui`, `actions`) are sufficient. Do not introduce a sixth without a compelling architectural reason discussed with the user.

### C5 — Keep CLAUDE.md accurate
If a structural change is made (new field, new package, new API constraint), update `CLAUDE.md` to reflect it.
