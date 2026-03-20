package fuckfetch

import (
	"log/slog"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

func (f *Fuckfetch) HandleCommand(e *events.ApplicationCommandInteractionCreate) {
	_ = e.DeferCreateMessage(true)

	info, err := GatherSystemInfo()
	if err != nil {
		f.logger.Error("failed to gather system info", slog.Any("error", err))
		_, _ = e.Client().Rest.UpdateInteractionResponse(
			e.Client().ApplicationID, e.Token(),
			discord.NewMessageUpdateV2([]discord.LayoutComponent{
				discord.NewContainer(
					discord.NewTextDisplay("❌ システム情報の取得に失敗: "+err.Error()),
				),
			}),
		)
		return
	}

	ui := BuildFuckfetchUI(info)
	_, _ = e.Client().Rest.UpdateInteractionResponse(
		e.Client().ApplicationID, e.Token(),
		discord.NewMessageUpdateV2([]discord.LayoutComponent{ui}),
	)
}
