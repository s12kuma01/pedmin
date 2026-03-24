// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package handler

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"

	disgobot "github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
	"github.com/s12kuma01/pedmin/internal/model"
	"github.com/s12kuma01/pedmin/internal/module"
	"github.com/s12kuma01/pedmin/internal/service"
	"github.com/s12kuma01/pedmin/internal/ui"
	"github.com/s12kuma01/pedmin/internal/view"
)

// LevelingBot is the interface the handler needs from the bot registry.
type LevelingBot interface {
	IsModuleEnabled(guildID snowflake.ID, moduleID string) bool
}

// LevelingHandler implements module.Module for the leveling feature.
type LevelingHandler struct {
	bot     LevelingBot
	service *service.LevelingService
	logger  *slog.Logger
}

// NewLevelingHandler creates a new LevelingHandler.
func NewLevelingHandler(bot LevelingBot, svc *service.LevelingService, logger *slog.Logger) *LevelingHandler {
	return &LevelingHandler{
		bot:     bot,
		service: svc,
		logger:  logger,
	}
}

// SetupLevelingListeners registers the message listener on the Discord client.
func SetupLevelingListeners(client *disgobot.Client, h *LevelingHandler) {
	client.AddEventListeners(
		disgobot.NewListenerFunc(h.onMessageCreate),
	)
}

// Shutdown stops the voice XP ticker.
func (h *LevelingHandler) Shutdown() {
	h.service.Shutdown()
}

func (h *LevelingHandler) Info() module.Info {
	return module.Info{
		ID:          model.LevelingModuleID,
		Name:        "レベリング",
		Description: "XPレベリングシステム",
		AlwaysOn:    false,
	}
}

func (h *LevelingHandler) Commands() []discord.ApplicationCommandCreate {
	return []discord.ApplicationCommandCreate{
		discord.SlashCommandCreate{
			Name:        "rank",
			Description: "ランクカードを表示する",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionUser{
					Name:        "user",
					Description: "対象ユーザー（省略時は自分）",
					Required:    false,
				},
			},
		},
		discord.SlashCommandCreate{
			Name:        "leaderboard",
			Description: "レベルランキングを表示する",
		},
	}
}

func (h *LevelingHandler) HandleCommand(e *events.ApplicationCommandInteractionCreate) {
	data := e.SlashCommandInteractionData()

	switch data.CommandName() {
	case "rank":
		h.handleRank(e, data)
	case "leaderboard":
		h.handleLeaderboard(e)
	}
}

func (h *LevelingHandler) handleRank(e *events.ApplicationCommandInteractionCreate, data discord.SlashCommandInteractionData) {
	targetUser := e.User().ID
	if user, ok := data.OptUser("user"); ok {
		targetUser = user.ID
	}

	guildID := e.GuildID()
	if guildID == nil {
		return
	}

	_ = e.DeferCreateMessage(false)

	cardPNG, err := h.service.GenerateRankCard(context.Background(), *guildID, targetUser)
	if err != nil {
		h.logger.Error("failed to generate rank card", slog.Any("error", err))
		_, _ = e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(), discord.MessageUpdate{
			Content: ptrTo("ランクカードの生成に失敗しました。"),
		})
		return
	}

	_, _ = e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(), discord.MessageUpdate{
		Files: []*discord.File{
			{
				Name:   "rank.png",
				Reader: bytes.NewReader(cardPNG),
			},
		},
	})
}

func (h *LevelingHandler) handleLeaderboard(e *events.ApplicationCommandInteractionCreate) {
	guildID := e.GuildID()
	if guildID == nil {
		return
	}

	entries, err := h.service.GetLeaderboard(*guildID, 10, 0)
	if err != nil {
		h.logger.Error("failed to get leaderboard", slog.Any("error", err))
		_ = e.CreateMessage(ui.EphemeralError("ランキングの取得に失敗しました。"))
		return
	}

	// Count total pages (approximate)
	totalPages := 1
	if len(entries) == 10 {
		totalPages = 10 // Estimate; will be refined on navigation
	}

	_ = e.CreateMessage(view.LevelingLeaderboard(entries, 0, totalPages))
}

func (h *LevelingHandler) SettingsSummary(guildID snowflake.ID) string {
	settings, err := h.service.LoadSettings(guildID)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("XP: %d-%d, CD: %ds", settings.MinXP, settings.MaxXP, settings.CooldownSeconds)
}

func (h *LevelingHandler) SettingsPanel(guildID snowflake.ID) []discord.LayoutComponent {
	settings, err := h.service.LoadSettings(guildID)
	if err != nil {
		h.logger.Error("failed to load leveling settings", slog.Any("error", err))
		settings = model.DefaultLevelingSettings()
	}
	rewardCount, err := h.service.CountRoleRewards(guildID)
	if err != nil {
		h.logger.Error("failed to count role rewards", slog.Any("error", err))
	}
	return view.LevelingSettingsPanel(settings, rewardCount)
}

// OnVoiceStateUpdate implements module.VoiceStateListener.
func (h *LevelingHandler) OnVoiceStateUpdate(guildID, channelID, userID snowflake.ID) {
	if !h.bot.IsModuleEnabled(guildID, model.LevelingModuleID) {
		return
	}
	h.service.OnVoiceStateUpdate(guildID, channelID, userID)
}

func ptrTo[T any](v T) *T {
	return &v
}
