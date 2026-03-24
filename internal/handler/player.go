// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package handler

import (
	"fmt"
	"log/slog"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgolink/v3/lavalink"
	"github.com/disgoorg/snowflake/v2"
	"github.com/s12kuma01/pedmin/internal/model"
	"github.com/s12kuma01/pedmin/internal/module"
	"github.com/s12kuma01/pedmin/internal/service"
	"github.com/s12kuma01/pedmin/internal/ui"
	"github.com/s12kuma01/pedmin/internal/view"
)

// PlayerHandler implements module.Module for the player feature.
type PlayerHandler struct {
	service *service.PlayerService
	logger  *slog.Logger
}

// NewPlayerHandler creates a new PlayerHandler.
func NewPlayerHandler(svc *service.PlayerService, logger *slog.Logger) *PlayerHandler {
	return &PlayerHandler{
		service: svc,
		logger:  logger,
	}
}

func (h *PlayerHandler) Info() module.Info {
	return module.Info{
		ID:          model.PlayerModuleID,
		Name:        "ミュージックプレイヤー",
		Description: "様々なソースから音楽を再生するミュージックプレイヤー",
		AlwaysOn:    false,
	}
}

func (h *PlayerHandler) Commands() []discord.ApplicationCommandCreate {
	return []discord.ApplicationCommandCreate{
		discord.SlashCommandCreate{
			Name:        "player",
			Description: "ミュージックプレイヤーを表示",
		},
	}
}

func (h *PlayerHandler) HandleCommand(e *events.ApplicationCommandInteractionCreate) {
	guildID := e.GuildID()
	if guildID == nil {
		_ = e.CreateMessage(ui.EphemeralError("このコマンドはサーバー内でのみ使用できます。"))
		return
	}

	if err := h.service.JoinVoiceChannel(*guildID, e.Member().User.ID); err != nil {
		_ = e.CreateMessage(ui.EphemeralError("ボイスチャンネルに接続してからコマンドを実行してください。"))
		return
	}

	h.service.DeleteTrackedMessage(*guildID)

	player := h.service.Lavalink().Player(*guildID)
	volume := h.service.GetDefaultVolume(*guildID)
	if player.Volume() != volume {
		ctx, cancel := h.service.LavalinkCtx()
		_ = player.Update(ctx, lavalink.WithVolume(volume))
		cancel()
	}
	playerUI := h.service.BuildPlayerUI(*guildID)

	_ = e.CreateMessage(discord.NewMessageCreateV2(playerUI))

	msg, err := e.Client().Rest.GetInteractionResponse(e.Client().ApplicationID, e.Token())
	if err == nil {
		h.service.TrackMessage(*guildID, msg.ChannelID, msg.ID)
	}
}

func (h *PlayerHandler) HandleComponent(e *events.ComponentInteractionCreate) {
	h.handleComponent(e)
}

func (h *PlayerHandler) HandleModal(e *events.ModalSubmitInteractionCreate) {
	h.handleModal(e)
}

func (h *PlayerHandler) SettingsSummary(guildID snowflake.ID) string {
	vol := h.service.GetDefaultVolume(guildID)
	return fmt.Sprintf("デフォルト音量: %d%%", vol)
}

func (h *PlayerHandler) SettingsPanel(guildID snowflake.ID) []discord.LayoutComponent {
	vol := h.service.GetDefaultVolume(guildID)
	return view.BuildPlayerSettingsPanel(vol)
}

// OnVoiceStateUpdate implements module.VoiceStateListener.
func (h *PlayerHandler) OnVoiceStateUpdate(guildID, channelID, userID snowflake.ID) {
	h.service.OnVoiceStateUpdate(guildID, channelID, userID)
}

// Shutdown stops all background goroutines. Call during graceful shutdown.
func (h *PlayerHandler) Shutdown() {
	h.service.Shutdown()
}
