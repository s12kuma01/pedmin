// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package ui

import (
	"fmt"
	"math"
	"strings"
	"time"
)

// FormatBytes formats a byte count into a human-readable string (e.g. "1.5 GB").
func FormatBytes(bytes uint64) string {
	if bytes == 0 {
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

// block characters from full (█, index 8) to thinnest (▏, index 1).
var blocks = [9]rune{' ', '▏', '▎', '▍', '▌', '▋', '▊', '▉', '█'}

// BuildBarRaw builds a progress bar string without surrounding backticks.
func BuildBarRaw(percent float64, total int) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}

	// Each cell has 8 sub-steps.
	steps := percent / 100 * float64(total) * 8
	fullCells := int(steps) / 8
	partialIdx := int(steps) % 8

	var sb strings.Builder
	for i := 0; i < total; i++ {
		if i < fullCells {
			sb.WriteRune(blocks[8]) // █
		} else if i == fullCells {
			sb.WriteRune(blocks[partialIdx])
		} else {
			sb.WriteRune(' ')
		}
	}
	return sb.String()
}

// BuildBar builds a progress bar using Unicode block characters, wrapped in backticks.
func BuildBar(percent float64, total int, showPercent bool) string {
	bar := "`" + BuildBarRaw(percent, total) + "`"
	if showPercent {
		bar += fmt.Sprintf(" %.1f%%", percent)
	}
	return bar
}

// FormatUptime formats a duration into a human-readable uptime string (e.g. "3d 2h 15m").
func FormatUptime(d time.Duration) string {
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
