package player

import (
	"fmt"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgolink/v3/disgolink"
)

func BuildQueueUI(queue *Queue, player disgolink.Player) discord.ContainerComponent {
	tracks := queue.Tracks()
	currentIdx := queue.CurrentIndex()

	if len(tracks) == 0 {
		return discord.NewContainer(
			discord.NewTextDisplay("### キュー"),
			discord.NewTextDisplay("キューは空です。"),
			discord.NewLargeSeparator(),
			discord.NewActionRow(
				discord.NewSecondaryButton("← 戻る", ModuleID+":back"),
				discord.NewDangerButton("キューをクリア", ModuleID+":clear_queue"),
			),
		).WithAccentColor(accentIdle)
	}

	var lines []string
	maxShow := 15
	start := 0
	if currentIdx > 5 {
		start = currentIdx - 5
	}
	end := start + maxShow
	if end > len(tracks) {
		end = len(tracks)
	}

	for i := start; i < end; i++ {
		t := tracks[i]
		prefix := "  "
		if i == currentIdx {
			prefix = "▶ "
		}
		lines = append(lines, fmt.Sprintf("%s**%d.** %s - %s `[%s]`",
			prefix, i+1, t.Info.Title, t.Info.Author, formatDuration(t.Info.Length)))
	}

	if end < len(tracks) {
		lines = append(lines, fmt.Sprintf("  *...他 %d曲*", len(tracks)-end))
	}

	components := []discord.ContainerSubComponent{
		discord.NewTextDisplay(fmt.Sprintf("### キュー (%d曲)", len(tracks))),
		discord.NewTextDisplay(strings.Join(lines, "\n")),
		discord.NewLargeSeparator(),
		discord.NewActionRow(
			discord.NewSecondaryButton("← 戻る", ModuleID+":back"),
			discord.NewDangerButton("キューをクリア", ModuleID+":clear_queue"),
		),
	}

	return discord.NewContainer(components...).WithAccentColor(accentPlaying)
}
