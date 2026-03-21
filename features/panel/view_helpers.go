package panel

import "fmt"

func formatMBToHuman(mb int) string {
	if mb >= 1024 {
		return fmt.Sprintf("%.1f GB", float64(mb)/1024)
	}
	return fmt.Sprintf("%d MB", mb)
}

func statusEmoji(state string) string {
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
