package view

import "fmt"

// PanelFormatMBToHuman formats megabytes to a human-readable string.
func PanelFormatMBToHuman(mb int) string {
	if mb >= 1024 {
		return fmt.Sprintf("%.1f GB", float64(mb)/1024)
	}
	return fmt.Sprintf("%d MB", mb)
}

// PanelStatusEmoji returns the status emoji for a server state.
func PanelStatusEmoji(state string) string {
	switch state {
	case "running":
		return "🟢"
	case "starting", "stopping":
		return "🟡"
	case "offline":
		return "🔴"
	default:
		return "⚪"
	}
}
