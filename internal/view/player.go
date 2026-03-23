// opyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package view

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/s12kuma01/pedmin/internal/model"
)

// BuildPlayerUI builds the main player UI with track info and controls.
func BuildPlayerUI(player disgolink.Player, queue *model.Queue) discord.ContainerComponent {
	track := player.Track()
	if track == nil {
		return BuildIdleUI(queue)
	}

	info := track.Info

	components := []discord.ContainerSubComponent{
		discord.NewSection(
			discord.NewTextDisplay("### \u25b6\ufe0f 再生中"),
			discord.NewTextDisplay(fmt.Sprintf("**%s**\nby %s", info.Title, info.Author)),
		).WithAccessory(BuildThumbnail(info)),
		discord.NewLargeSeparator(),
		discord.NewTextDisplay(fmt.Sprintf(
			"%s  %s / %s  |  %s",
			BuildProgressBar(player.Position(), info.Length),
			FormatDuration(player.Position()),
			FormatDuration(info.Length),
			queue.LoopMode().String(),
		)),
		discord.NewLargeSeparator(),
	}
	components = append(components, BuildButtonRows(queue.LoopMode())...)
	return discord.NewContainer(components...)
}

// BuildIdleUI builds the idle player UI when no track is playing.
func BuildIdleUI(queue *model.Queue) discord.ContainerComponent {
	components := []discord.ContainerSubComponent{
		discord.NewTextDisplay("### Pedmin Player"),
		discord.NewTextDisplay("再生中の曲はありません。ボタンから曲を追加してください！"),
		discord.NewLargeSeparator(),
	}
	components = append(components, BuildButtonRows(queue.LoopMode())...)
	return discord.NewContainer(components...)
}

// BuildButtonRows builds the player control button rows.
func BuildButtonRows(loopMode model.LoopMode) []discord.ContainerSubComponent {
	modeLabel := "モード"
	switch loopMode {
	case model.LoopTrack:
		modeLabel = "モード: トラック"
	case model.LoopQueue:
		modeLabel = "モード: キュー"
	}

	return []discord.ContainerSubComponent{
		discord.NewActionRow(
			discord.NewSecondaryButton("\u23ea", model.PlayerModuleID+":seek_back"),
			discord.NewSecondaryButton("スキップ", model.PlayerModuleID+":skip"),
			discord.NewDangerButton("停止", model.PlayerModuleID+":stop"),
			discord.NewSecondaryButton("\u23e9", model.PlayerModuleID+":seek_forward"),
			discord.NewSuccessButton("追加", model.PlayerModuleID+":add"),
		),
		discord.NewActionRow(
			discord.NewSecondaryButton("\U0001f500 シャッフル", model.PlayerModuleID+":shuffle"),
			discord.NewSecondaryButton("キュー", model.PlayerModuleID+":queue"),
			discord.NewSecondaryButton(modeLabel, model.PlayerModuleID+":loop"),
		),
	}
}
