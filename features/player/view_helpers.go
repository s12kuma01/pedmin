package player

import (
	"fmt"
	"strings"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgolink/v3/lavalink"
)

func buildThumbnail(info lavalink.TrackInfo) discord.ThumbnailComponent {
	if info.ArtworkURL != nil && *info.ArtworkURL != "" {
		return discord.NewThumbnail(*info.ArtworkURL)
	}
	return discord.NewThumbnail("https://cdn.discordapp.com/embed/avatars/0.png")
}

func buildProgressBar(position, length lavalink.Duration) string {
	total := 20
	if length <= 0 {
		return strings.Repeat("▬", total)
	}
	filled := int(float64(position) / float64(length) * float64(total))
	if filled > total {
		filled = total
	}
	return strings.Repeat("▓", filled) + strings.Repeat("░", total-filled)
}

func formatDuration(d lavalink.Duration) string {
	dur := time.Duration(d) * time.Millisecond
	hours := int(dur.Hours())
	minutes := int(dur.Minutes()) % 60
	seconds := int(dur.Seconds()) % 60

	if hours > 0 {
		return fmt.Sprintf("%d:%02d:%02d", hours, minutes, seconds)
	}
	return fmt.Sprintf("%d:%02d", minutes, seconds)
}

func ephemeralV2Error(text string) discord.MessageCreate {
	return discord.NewMessageCreateV2(
		discord.NewContainer(
			discord.NewTextDisplay(fmt.Sprintf("❌ %s", text)),
		),
	).WithEphemeral(true)
}
