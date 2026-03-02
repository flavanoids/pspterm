# pspterm — Skills & Recipes

Concrete how-to patterns for working on this project.

---

## Build & verify

```sh
# Build
go build -o pspterm .

# Build + run immediately
go build -o pspterm . && ./pspterm

# Vet only (no binary)
go vet ./...
```

---

## Hot-reload config
Press `r` while pspterm is running. No restart needed.
The config is re-read from `~/.config/pspterm/config.yaml` and styles recomputed.

---

## Reset config to defaults
```sh
rm ~/.config/pspterm/config.yaml
# next run recreates it from config/defaults.go
```

---

## Add a new category item type

**1. Extend the struct** (`config/config.go`):
```go
type Item struct {
    Name    string `yaml:"name"`
    Type    string `yaml:"type"`
    Path    string `yaml:"path"`
    Command string `yaml:"command"`
    URL     string `yaml:"url"`
    // Add your new field here
    MyField string `yaml:"my_field"`
}
```

**2. Handle in executor** (`actions/executor.go`):
```go
case "mytype":
    // implement action
    return nil, "", nil
```

---

## Add a new keybinding

In `model/keys.go`, add to `KeyMap` struct:
```go
MyAction key.Binding
```

In `DefaultKeyMap()`:
```go
MyAction: key.NewBinding(
    key.WithKeys("m"),
    key.WithHelp("m", "my action"),
),
```

Handle in `model/model.go` `handleKey()`:
```go
case key.Matches(msg, m.keys.MyAction):
    // ...
```

---

## Add a new UI style

In `ui/styles.go`:
```go
// Struct field
MyStyle lipgloss.Style

// In NewStyles():
MyStyle: lipgloss.NewStyle().Foreground(accent).Bold(true),
```

Use in `model/model.go` `View()`:
```go
sb.WriteString(m.styles.MyStyle.Render("text"))
```

---

## Tune spring feel

In `model/model.go` `New()`:
```go
spring: harmonica.NewSpring(harmonica.FPS(60), angFreq, dampingRatio),
```

| Feel | angFreq | dampingRatio |
|------|---------|--------------|
| Slow, floaty | 6 | 0.7 |
| PSP-authentic | 14 | 0.55 |
| Snappy, tight | 20 | 0.75 |
| Very bouncy | 12 | 0.42 |

---

## Add a new theme color to config

**1.** Add field to `config.Theme` in `config/config.go`:
```go
type Theme struct {
    AccentColor  string `yaml:"accent_color"`
    DimColor     string `yaml:"dim_color"`
    MyColor      string `yaml:"my_color"`
}
```

**2.** Use in `ui/styles.go` `NewStyles(t config.Theme)`:
```go
myColor := lipgloss.Color(t.MyColor)
```

**3.** Add default value in `config/defaults.go` YAML template.

---

## Shell wrapper setup (one-time)

Add to `~/.bashrc` or `~/.zshrc`:
```sh
function psp() {
    local target
    target="$(/path/to/pspterm)"
    [ -d "$target" ] && cd "$target"
}
```

The TUI writes to `/dev/tty` directly; the exit path goes to stdout for the wrapper to capture.

---

## Understand the clock

- `model.now` (`time.Time`) is updated on every `TickMsg`.
- `TickMsg` is `type TickMsg time.Time` — cast with `time.Time(msg)`.
- Rendered in `View()` as `m.now.Format("15:04")` — 24h HH:MM, right-aligned.
- Style: `m.styles.Clock` (`ui/styles.go`).
