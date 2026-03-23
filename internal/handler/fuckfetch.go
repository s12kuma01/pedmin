package handler

import (
	"log/slog"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
	"github.com/s12kuma01/pedmin/internal/model"
	"github.com/s12kuma01/pedmin/internal/module"
	"github.com/s12kuma01/pedmin/internal/service"
	"github.com/s12kuma01/pedmin/internal/view"
)

// FuckfetchHandler implements module.Module for the fuckfetch feature.
type FuckfetchHandler struct {
	logger *slog.Logger
}

// NewFuckfetchHandler creates a new FuckfetchHandler.
func NewFuckfetchHandler(logger *slog.Logger) *FuckfetchHandler {
	return &FuckfetchHandler{logger: logger}
}

func (h *FuckfetchHandler) Info() module.Info {
	return module.Info{
		ID:          model.FuckfetchModuleID,
		Name:        "Fuckfetch",
		Description: "サーバーのシステム情報を表示",
		AlwaysOn:    true,
	}
}

func (h *FuckfetchHandler) Commands() []discord.ApplicationCommandCreate {
	return []discord.ApplicationCommandCreate{
		discord.SlashCommandCreate{
			Name:        "fuckfetch",
			Description: "サーバーマシンのシステム情報を表示する",
		},
	}
}

func (h *FuckfetchHandler) HandleCommand(e *events.ApplicationCommandInteractionCreate) {
	_ = e.DeferCreateMessage(true)

	info, err := service.GatherSystemInfo()
	if err != nil {
		h.logger.Error("failed to gather system info", slog.Any("error", err))
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

	ui := view.BuildFuckfetchOutput(info)
	_, _ = e.Client().Rest.UpdateInteractionResponse(
		e.Client().ApplicationID, e.Token(),
		discord.NewMessageUpdateV2([]discord.LayoutComponent{ui}),
	)
}

func (h *FuckfetchHandler) HandleComponent(_ *events.ComponentInteractionCreate) {}
func (h *FuckfetchHandler) HandleModal(_ *events.ModalSubmitInteractionCreate)   {}
func (h *FuckfetchHandler) SettingsPanel(_ snowflake.ID) []discord.LayoutComponent {
	return nil
}
