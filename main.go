package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	xterm "github.com/charmbracelet/x/term"
	"github.com/user/pspterm/config"
	"github.com/user/pspterm/model"
)

func main() {
	cfg, err := config.LoadOrCreate()
	if err != nil {
		fmt.Fprintf(os.Stderr, "pspterm: failed to load config: %v\n", err)
		os.Exit(1)
	}

	m := model.New(cfg)

	opts := []tea.ProgramOption{tea.WithAltScreen()}

	// Open /dev/tty explicitly so the TUI renders to the real terminal device
	// even when stdout is captured (e.g. inside a shell function: target=$(pspterm)).
	// This also makes it work correctly in SSH sessions where stdin/stdout may
	// be redirected.
	tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err == nil {
		// Both input and output go through the raw terminal device.
		// The TUI draws on /dev/tty; the exit path is printed to real stdout.
		opts = append(opts, tea.WithInput(tty), tea.WithOutput(tty))
		defer tty.Close()
	} else {
		// /dev/tty not available — fall back to stdin/stdout if they are a tty.
		if !xterm.IsTerminal(os.Stdin.Fd()) {
			fmt.Fprintln(os.Stderr, "pspterm: no terminal available (not a tty)")
			os.Exit(1)
		}
		// stdin is already a tty; bubbletea will use it by default.
	}

	p := tea.NewProgram(m, opts...)

	finalModel, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "pspterm: %v\n", err)
		os.Exit(1)
	}

	// Eval trick: if the user selected a directory item, print its path to
	// stdout so the shell wrapper can cd to it.
	if fm, ok := finalModel.(model.Model); ok {
		if path := fm.ExitPath(); path != "" {
			fmt.Println(path)
		}
	}
}
