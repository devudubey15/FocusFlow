package watcher

import (
	"os/exec"
	"strings"
)

type WindowInfo struct {
	AppName string
	Title   string
}

// GetActiveWindowX11 uses xdotool to get the current active window title and class (app name)
func GetActiveWindowX11() (*WindowInfo, error) {
	// Get window ID
	outId, err := exec.Command("xdotool", "getactivewindow").Output()
	if err != nil {
		return nil, err
	}
	windowId := strings.TrimSpace(string(outId))

	// Get window name (title)
	outTitle, err := exec.Command("xdotool", "getwindowname", windowId).Output()
	if err != nil {
		return nil, err
	}
	title := strings.TrimSpace(string(outTitle))

	// Get window class (app name)
	// xprop -id <id> WM_CLASS
	outClass, err := exec.Command("xprop", "-id", windowId, "WM_CLASS").Output()
	if err != nil {
		return nil, err
	}
	// Example output: WM_CLASS(STRING) = "code", "Code"
	classParts := strings.Split(string(outClass), "=")
	appName := "unknown"
	if len(classParts) > 1 {
		parts := strings.Split(classParts[1], ",")
		if len(parts) > 0 {
			appName = strings.Trim(strings.TrimSpace(parts[len(parts)-1]), "\"")
		}
	}

	return &WindowInfo{
		AppName: appName,
		Title:   title,
	}, nil
}
