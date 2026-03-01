package actions

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/user/pspterm/config"
)

// Execute dispatches an item's action.
// onExecDone is called with the error (if any) when a command finishes.
// Returns (cmd, exitPath, err):
//   - cmd != nil: a tea.Cmd to run (for ExecProcess)
//   - exitPath != "": the model should store and quit
//   - err != nil: display error
func Execute(item config.Item, onExecDone func(error) tea.Msg) (tea.Cmd, string, error) {
	switch item.Type {
	case "directory":
		path, err := expandPath(item.Path)
		if err != nil {
			return nil, "", fmt.Errorf("invalid path %q: %w", item.Path, err)
		}
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return nil, "", fmt.Errorf("directory does not exist: %s", path)
		}
		return nil, path, nil

	case "command":
		if item.Command == "" {
			return nil, "", fmt.Errorf("command is empty")
		}
		c := exec.Command("sh", "-c", item.Command)
		cmd := tea.ExecProcess(c, onExecDone)
		return cmd, "", nil

	case "url":
		if item.URL == "" {
			return nil, "", fmt.Errorf("url is empty")
		}
		go exec.Command("xdg-open", item.URL).Start() //nolint:errcheck
		return nil, "", nil

	default:
		return nil, "", fmt.Errorf("unknown item type: %q", item.Type)
	}
}

// expandPath expands ~ and environment variables in a path.
func expandPath(path string) (string, error) {
	path = os.ExpandEnv(path)
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		path = filepath.Join(home, path[2:])
	} else if path == "~" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		path = home
	}
	return filepath.Clean(path), nil
}
