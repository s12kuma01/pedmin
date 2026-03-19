package player

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgolink/v3/disgolink"
)

const accentPlaying = 0x00B894
const accentIdle = 0x636E72

func BuildPlayerUI(player disgolink.Player, queue *Queue) discord.ContainerComponent {
	track := player.Track()
	if track == nil {
		return buildIdleUI(queue)
	}

	info := track.Info

	return discord.NewContainer(
		discord.NewSection(
			discord.NewTextDisplay("### ▶️ 再生中"),
			discord.NewTextDisplay(fmt.Sprintf("**%s**\nby %s", info.Title, info.Author)),
		).WithAccessory(buildThumbnail(info)),
		discord.NewLargeSeparator(),
		discord.NewTextDisplay(fmt.Sprintf(
			"%s  %s / %s  |  %s",
			buildProgressBar(player.Position(), info.Length),
			formatDuration(player.Position()),
			formatDuration(info.Length),
			queue.LoopMode().String(),
		)),
		discord.NewLargeSeparator(),
		buildButtonRow(queue.LoopMode()),
	).WithAccentColor(accentPlaying)
}

func buildIdleUI(queue *Queue) discord.ContainerComponent {
	return discord.NewContainer(
		discord.NewTextDisplay("### Pedmin Player"),
		discord.NewTextDisplay("再生中の曲はありません。ボタンから曲を追加してください！"),
		discord.NewLargeSeparator(),
		buildButtonRow(queue.LoopMode()),
	).WithAccentColor(accentIdle)
}

func buildButtonRow(loopMode LoopMode) discord.ActionRowComponent {
	modeLabel := "モード"
	switch loopMode {
	case LoopTrack:
		modeLabel = "モード: トラック"
	case LoopQueue:
		modeLabel = "モード: キュー"
	}

	return discord.NewActionRow(
		discord.NewSecondaryButton("スキップ", ModuleID+":skip"),
		discord.NewDangerButton("停止", ModuleID+":stop"),
		discord.NewSecondaryButton("キュー", ModuleID+":queue"),
		discord.NewSuccessButton("追加", ModuleID+":add"),
		discord.NewSecondaryButton(modeLabel, ModuleID+":loop"),
	)
}
