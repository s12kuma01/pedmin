// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package view

import (
	"fmt"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/s12kuma01/pedmin/internal/model"
)

// BuildQueueUI builds the queue display UI.
func BuildQueueUI(queue *model.Queue, player disgolink.Player) discord.ContainerComponent {
	tracks := queue.Tracks()
	currentIdx := queue.CurrentIndex()

	if len(tracks) == 0 {
		return discord.NewContainer(
			discord.NewTextDisplay("### キュー"),
			discord.NewTextDisplay("キューは空です。"),
			discord.NewLargeSeparator(),
			discord.NewActionRow(
				discord.NewSecondaryButton("\u2190 戻る", model.PlayerModuleID+":back"),
				discord.NewDangerButton("キューをクリア", model.PlayerModuleID+":clear_queue"),
			),
		)
	}

	var lines []string
	maxShow := 15
	start := 0
	// Keep the current track near the top (5th position) for context.
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
			prefix = "\u25b6 "
		}
		lines = append(lines, fmt.Sprintf("%s**%d.** %s - %s `[%s]`",
			prefix, i+1, t.Info.Title, t.Info.Author, FormatDuration(t.Info.Length)))
	}

	if end < len(tracks) {
		lines = append(lines, fmt.Sprintf("  *...他 %d曲*", len(tracks)-end))
	}

	components := []discord.ContainerSubComponent{
		discord.NewTextDisplay(fmt.Sprintf("### キュー (%d曲)", len(tracks))),
		discord.NewTextDisplay(strings.Join(lines, "\n")),
		discord.NewLargeSeparator(),
		discord.NewActionRow(
			discord.NewSecondaryButton("\u2190 戻る", model.PlayerModuleID+":back"),
			discord.NewDangerButton("キューをクリア", model.PlayerModuleID+":clear_queue"),
		),
	}

	return discord.NewContainer(components...)
}
