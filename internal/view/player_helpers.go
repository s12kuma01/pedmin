// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package view

import (
	"fmt"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgolink/v3/lavalink"
	"github.com/Sumire-Labs/pedmin/internal/ui"
)

// BuildThumbnail builds a thumbnail component from track info artwork.
func BuildThumbnail(info lavalink.TrackInfo) discord.ThumbnailComponent {
	if info.ArtworkURL != nil && *info.ArtworkURL != "" {
		return discord.NewThumbnail(*info.ArtworkURL)
	}
	return discord.NewThumbnail("https://cdn.discordapp.com/embed/avatars/0.png")
}

// BuildProgressBar builds a text-based progress bar for the player.
func BuildProgressBar(position, length lavalink.Duration) string {
	if length <= 0 {
		return ui.BuildBar(100, 20, false)
	}
	percent := float64(position) / float64(length) * 100
	return ui.BuildBar(percent, 20, false)
}

// FormatDuration formats a lavalink duration as h:mm:ss or m:ss.
func FormatDuration(d lavalink.Duration) string {
	dur := time.Duration(d) * time.Millisecond
	hours := int(dur.Hours())
	minutes := int(dur.Minutes()) % 60
	seconds := int(dur.Seconds()) % 60

	if hours > 0 {
		return fmt.Sprintf("%d:%02d:%02d", hours, minutes, seconds)
	}
	return fmt.Sprintf("%d:%02d", minutes, seconds)
}
