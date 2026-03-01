# pspterm

A PSP XMB-style terminal launcher written in Go with [Bubble Tea](https://github.com/charmbracelet/bubbletea).

Navigate categories horizontally with spring-animated bounce, select items vertically, and launch commands, open directories, or visit URLs — all from a slick retro interface.

![pspterm](https://github.com/flavanoids/pspterm/raw/main/preview.png)

---

## Features

- **XMB navigation** — horizontal category bar with spring animation (harmonica)
- **Per-category accent colors** — override the theme for individual categories
- **Item descriptions** — optional subtitles shown below the selected item
- **Live config reload** — press `r` to reload `config.yaml` without restarting
- **In-TUI item manager** — add, edit, and delete items without touching YAML
- **`cd` on select** — directory items exit and change your shell's working directory (via a shell wrapper)
- **Editor integration** — open `config.yaml` in your preferred editor from within the UI

---

## Install

```sh
git clone https://github.com/flavanoids/pspterm.git
cd pspterm
go build -o pspterm .
```

Requires Go 1.21+.

---

## Shell wrapper

Directory items only work if you launch pspterm through this shell function. Add it to your `~/.bashrc` or `~/.zshrc`:

```sh
function psp() {
    local target="$(./pspterm)"
    [ -d "$target" ] && cd "$target"
}
```

Then run `psp` instead of `./pspterm`.

---

## Keybindings

| Key | Action |
|-----|--------|
| `←` / `→` | Switch category |
| `↑` / `↓` | Select item |
| `Enter` | Launch selected item |
| `r` | Reload config |
| `q` / `Ctrl+C` | Quit |

---

## Configuration

Config lives at `~/.config/pspterm/config.yaml` and is created automatically on first run. A fully-annotated reference is written to `~/.config/pspterm/config.yaml.example` on every startup.

```yaml
theme:
  accent_color: "#4fc8ff"   # selected category / item highlight
  dim_color:    "#1a3f5c"   # unselected categories

editor: ""  # preferred editor binary; falls back to $EDITOR then auto-detect

categories:
  - name: "Game"
    icon: "⊞"
    color: "#ff8800"        # optional per-category accent color
    items:
      - name: "Shell"
        description: "Open an interactive shell"   # optional subtitle
        type: command
        command: "$SHELL"

      - name: "Vim"
        type: command
        command: "vim"

  - name: "Files"
    icon: "⊡"
    items:
      - name: "Home"
        type: directory
        path: "~"

      - name: "GitHub"
        type: url
        url: "https://github.com"
```

### Item types

| Type | Field | Description |
|------|-------|-------------|
| `command` | `command:` | Runs a shell command (interactive apps like vim/htop work) |
| `directory` | `path:` | Exits pspterm and `cd`s to the path (needs shell wrapper) |
| `url` | `url:` | Opens in the default browser via `xdg-open` |
| `manager` | — | Opens the in-TUI item manager |
| `editconfig` | — | Opens `config.yaml` in your preferred editor |

---

## Project layout

```
main.go               — entry point
config/
  config.go           — Config/Theme/Category/Item structs + YAML load/save
  defaults.go         — LoadOrCreate, embedded default config YAML
model/
  model.go            — Bubble Tea model (Init/Update/View), spring animation
  editor.go           — In-TUI item manager state + logic
  keys.go             — KeyMap
  messages.go         — TickMsg, ExecDoneMsg, etc.
ui/
  styles.go           — lipgloss Styles
  xmb.go              — Horizontal category bar renderer
  items.go            — Vertical item list renderer
  editor.go           — Item manager UI
actions/
  executor.go         — Execute(item) → directory/command/url
```

---

## License

MIT
