package embedfix

import (
	"context"
	"log/slog"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/s12kuma01/pedmin/ui"
)

func (ef *EmbedFix) handleComponent(e *events.ComponentInteractionCreate) {
	customID := e.Data.CustomID()
	_, rest, _ := strings.Cut(customID, ":")
	action, rest, _ := strings.Cut(rest, ":")

	switch action {
	case "translate":
		ef.handleTranslate(e, rest)
	case "platforms":
		ef.handlePlatformSettings(e)
	}
}

func (ef *EmbedFix) handleTranslate(e *events.ComponentInteractionCreate, rest string) {
	_ = e.DeferUpdateMessage()

	if !ef.translateClient.IsAvailable() {
		ef.respondTranslateError(e, "翻訳APIキーが設定されていないため、翻訳できません。")
		return
	}

	platform, params, _ := strings.Cut(rest, ":")
	ctx := context.Background()

	ui, err := ef.translateContent(ctx, platform, params)
	if err != nil {
		ef.respondTranslateError(e, "翻訳に失敗しました。")
		return
	}

	_, _ = e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(),
		discord.NewMessageUpdateV2(ui))
}

func (ef *EmbedFix) handlePlatformSettings(e *events.ComponentInteractionCreate) {
	guildID := e.GuildID()
	if guildID == nil {
		return
	}

	data, ok := e.Data.(discord.StringSelectMenuInteractionData)
	if !ok {
		return
	}

	settings, err := LoadSettings(ef.store, *guildID)
	if err != nil {
		ef.logger.Error("failed to load embedfix settings", slog.Any("error", err))
		return
	}

	// Disable all, then enable selected
	for k := range settings.Platforms {
		settings.Platforms[k] = false
	}
	for _, v := range data.Values {
		settings.Platforms[Platform(v)] = true
	}

	if err := SaveSettings(ef.store, *guildID, settings); err != nil {
		ef.logger.Error("failed to save embedfix settings", slog.Any("error", err))
	}

	settingsUI := BuildSettingsPanel(settings)
	enabled := ef.bot.IsModuleEnabled(*guildID, ModuleID)
	_ = e.UpdateMessage(ui.BuildModulePanel(ef.Info(), enabled, settingsUI))
}

func (ef *EmbedFix) respondTranslateError(e *events.ComponentInteractionCreate, msg string) {
	_, _ = e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(),
		discord.NewMessageUpdateV2([]discord.LayoutComponent{
			discord.NewContainer(
				discord.NewTextDisplay(msg),
			),
		}))
}
