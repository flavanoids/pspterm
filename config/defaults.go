package config

import (
	"os"
)

// templateYAML is the canonical, fully-commented reference for every supported
// config option.  It is written to config.yaml.example on every startup so
// users always have an up-to-date reference alongside their live config.
const templateYAML = `# ──────────────────────────────────────────────────────────────
# pspterm — configuration reference  (~/.config/pspterm/config.yaml)
# ──────────────────────────────────────────────────────────────
#
# Edit this file in any text editor.
# Press  r  inside pspterm to hot-reload without restarting.
# Select  Settings › Edit Config  to open it in your chosen editor.
#
# STRUCTURE
#   theme       — colour palette
#   editor      — preferred editor binary (optional)
#   categories  — list of XMB categories, each with items
#
# ITEM TYPES
#   command     runs a shell command (interactive apps like vim/htop work)
#   directory   exits pspterm and cd's to the path (needs shell wrapper)
#   url         opens in default browser via xdg-open
#   manager     opens the in-TUI item manager (built-in, no extra fields)
#   editconfig  opens config.yaml in your preferred editor (built-in)
#
# SHELL WRAPPER (required for "directory" items to actually cd)
#   Add to ~/.bashrc or ~/.zshrc:
#     function psp() {
#         local target="$(./pspterm)"
#         [ -d "$target" ] && cd "$target"
#     }
# ──────────────────────────────────────────────────────────────

# ── Theme ──────────────────────────────────────────────────────
theme:
  accent_color: "#4fc8ff"   # selected category / item highlight (PSP cold blue)
  dim_color:    "#1a3f5c"   # unselected categories — dark navy

# ── Editor ─────────────────────────────────────────────────────
# Binary used by  Settings › Edit Config.
# Leave blank to auto-detect: checks $EDITOR, $VISUAL, then tries
# nano → vim → vi → micro → emacs in order.
editor: ""

# ── Categories ─────────────────────────────────────────────────
# Each category appears as a node in the horizontal XMB bar.
# 'icon' accepts any Unicode character (emoji, box-drawing, etc.)
#
# Item fields (only fill in the field matching 'type'):
#   type: command   →  command: "shell command here"
#   type: directory →  path:    "~/some/path"
#   type: url       →  url:     "https://example.com"
#   type: manager   →  (no extra fields)
#   type: editconfig→  (no extra fields)
#
# Optional per-item description (subtitle shown below selected item):
#   description: "..."
#
# Optional per-category accent color (overrides theme accent):
#   color: "#4fc8ff"   # hex color, e.g. "#ff8800"
#
# Auto-scan applications from XDG .desktop files:
#   scan: true   # replaces items with discovered apps from ~/.local/share/applications
#                # and /usr/share/applications (plus $XDG_DATA_DIRS)

categories:
  - name: "Terminal"
    icon: "⌘"
    items:
      - name: "Shell"
        type: command
        command: "$SHELL"

      - name: "Vim"
        type: command
        command: "vim"

      - name: "Htop"
        type: command
        command: "htop"

  - name: "Files"
    icon: "⌂"
    items:
      - name: "Home"
        type: directory
        path: "~"
      - name: "Documents"
        type: directory
        path: "~/Documents"
      - name: "Downloads"
        type: directory
        path: "~/Downloads"
      - name: "Config"
        type: directory
        path: "~/.config"
      - name: "Pictures"
        type: directory
        path: "~/Pictures"
      - name: "Music"
        type: directory
        path: "~/Music"
      - name: "Videos"
        type: directory
        path: "~/Videos"

  - name: "Applications"
    icon: "⊞"
    scan: true

  - name: "Network"
    icon: "⊕"
    items:
      - name: "Ping Google"
        type: command
        command: "ping -c 4 google.com"
      - name: "IP Info"
        type: command
        command: "ip addr"

  - name: "Settings"
    icon: "⚙"
    items:
      - name: "Manage Items"
        type: manager
      - name: "Edit Config"
        type: editconfig
      - name: "Home"
        type: directory
        path: "~"
`

// defaultConfigYAML is what gets written as the user's live config on first run.
// It's intentionally leaner than templateYAML — just working defaults, no wall of comments.
const defaultConfigYAML = `# pspterm config — edit freely, press r in the UI to reload
# Full reference: ~/.config/pspterm/config.yaml.example

theme:
  accent_color: "#4fc8ff"
  dim_color:    "#1a3f5c"

# Preferred editor for 'Edit Config'.  Leave blank to auto-detect.
editor: ""

categories:
  - name: "Terminal"
    icon: "⌘"
    items:
      - name: "Shell"
        type: command
        command: "$SHELL"
      - name: "Vim"
        type: command
        command: "vim"
      - name: "Htop"
        type: command
        command: "htop"

  - name: "Files"
    icon: "⌂"
    items:
      - name: "Home"
        type: directory
        path: "~"
      - name: "Documents"
        type: directory
        path: "~/Documents"
      - name: "Downloads"
        type: directory
        path: "~/Downloads"
      - name: "Config"
        type: directory
        path: "~/.config"
      - name: "Pictures"
        type: directory
        path: "~/Pictures"
      - name: "Music"
        type: directory
        path: "~/Music"
      - name: "Videos"
        type: directory
        path: "~/Videos"

  - name: "Applications"
    icon: "⊞"
    scan: true

  - name: "Network"
    icon: "⊕"
    items:
      - name: "Ping Google"
        type: command
        command: "ping -c 4 google.com"
      - name: "IP Info"
        type: command
        command: "ip addr"

  - name: "Settings"
    icon: "⚙"
    items:
      - name: "Manage Items"
        type: manager
      - name: "Edit Config"
        type: editconfig
      - name: "Home"
        type: directory
        path: "~"
`

// LoadOrCreate loads the config, creating defaults if it doesn't exist.
// It also always writes the fully-annotated template to config.yaml.example
// so users have a current reference even after customising their config.
func LoadOrCreate() (Config, error) {
	if err := os.MkdirAll(ConfigDir(), 0755); err != nil {
		return Config{}, err
	}

	// Write (or refresh) the reference template
	_ = os.WriteFile(ExampleConfigPath(), []byte(templateYAML), 0644)

	if _, err := os.Stat(ConfigPath()); os.IsNotExist(err) {
		if err := os.WriteFile(ConfigPath(), []byte(defaultConfigYAML), 0644); err != nil {
			return Config{}, err
		}
	}
	return Load()
}
