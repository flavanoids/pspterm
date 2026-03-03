package config

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// desktopExecCleaner strips Exec field codes (%f, %F, %u, %U, etc.) for terminal launch.
var desktopExecCleaner = regexp.MustCompile(`%\w`)

// ScanApplications discovers applications from XDG .desktop files.
// Searches ~/.local/share/applications and /usr/share/applications (plus $XDG_DATA_DIRS).
// Returns items suitable for a category; skips NoDisplay=true and non-Application types.
func ScanApplications() []Item {
	dirs := xdgApplicationDirs()
	seen := make(map[string]bool) // dedupe by normalized name
	var items []Item

	for _, dir := range dirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, e := range entries {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".desktop") {
				continue
			}
			path := filepath.Join(dir, e.Name())
			item, ok := parseDesktopFile(path)
			if !ok || seen[item.Name] {
				continue
			}
			seen[item.Name] = true
			items = append(items, item)
		}
	}
	return items
}

func xdgApplicationDirs() []string {
	home, _ := os.UserHomeDir()
	dirs := []string{
		filepath.Join(home, ".local", "share", "applications"),
		"/usr/share/applications",
	}
	if xdg := os.Getenv("XDG_DATA_DIRS"); xdg != "" {
		for _, d := range filepath.SplitList(xdg) {
			dirs = append(dirs, filepath.Join(d, "applications"))
		}
	}
	return dirs
}

func parseDesktopFile(path string) (Item, bool) {
	f, err := os.Open(path)
	if err != nil {
		return Item{}, false
	}
	defer f.Close()

	var name, exec string
	var noDisplay bool
	var typ string

	scanner := bufio.NewScanner(f)
	inDesktopEntry := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "[") {
			inDesktopEntry = line == "[Desktop Entry]"
			continue
		}
		if !inDesktopEntry {
			continue
		}
		k, v, ok := parseDesktopKey(line)
		if !ok {
			continue
		}
		switch k {
		case "Name":
			if name == "" {
				name = v
			}
		case "Exec":
			exec = cleanExec(v)
		case "NoDisplay":
			noDisplay = v == "true" || v == "1"
		case "Type":
			typ = v
		}
	}

	if noDisplay || typ != "Application" || name == "" || exec == "" {
		return Item{}, false
	}
	return Item{
		Name:    name,
		Type:    "command",
		Command: exec,
	}, true
}

func parseDesktopKey(line string) (key, value string, ok bool) {
	i := strings.Index(line, "=")
	if i < 0 {
		return "", "", false
	}
	key = strings.TrimSpace(line[:i])
	value = strings.TrimSpace(line[i+1:])
	// Unescape per desktop spec: \\ first, then \s \n \t
	value = strings.ReplaceAll(value, `\\`, `\`)
	value = strings.ReplaceAll(value, `\s`, " ")
	value = strings.ReplaceAll(value, `\n`, "\n")
	value = strings.ReplaceAll(value, `\t`, "\t")
	return key, value, true
}

func cleanExec(exec string) string {
	return strings.TrimSpace(desktopExecCleaner.ReplaceAllString(exec, ""))
}
