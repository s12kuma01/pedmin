package panel

import (
	"fmt"
	"math"
	"strings"
	"time"
)

func formatBytes(bytes int64) string {
	if bytes <= 0 {
		return "0 B"
	}
	units := []string{"B", "KB", "MB", "GB", "TB"}
	i := int(math.Log(float64(bytes)) / math.Log(1024))
	if i >= len(units) {
		i = len(units) - 1
	}
	val := float64(bytes) / math.Pow(1024, float64(i))
	return fmt.Sprintf("%.1f %s", val, units[i])
}

func buildBar(percent float64, total int) string {
	filled := int(percent / 100 * float64(total))
	if filled > total {
		filled = total
	}
	if filled < 0 {
		filled = 0
	}
	return strings.Repeat("🟢", filled) + strings.Repeat("⚫", total-filled)
}

func formatUptime(ms int64) string {
	d := time.Duration(ms) * time.Millisecond
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60

	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
	}
	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
	return fmt.Sprintf("%dm", minutes)
}

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
