package view

import (
	"fmt"
	"strings"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgolink/v3/lavalink"
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
	total := 20
	if length <= 0 {
		return strings.Repeat("\u25ac", total)
	}
	filled := int(float64(position) / float64(length) * float64(total))
	if filled > total {
		filled = total
	}
	return strings.Repeat("\u2593", filled) + strings.Repeat("\u2591", total-filled)
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
