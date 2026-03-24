// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package handler

import (
	"fmt"
	"log/slog"
	"strings"

	disgobot "github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
	"github.com/s12kuma01/pedmin/internal/model"
	"github.com/s12kuma01/pedmin/internal/module"
	"github.com/s12kuma01/pedmin/internal/service"
	"github.com/s12kuma01/pedmin/internal/view"
)

// LoggerBot is the interface for checking module enabled status.
type LoggerBot interface {
	IsModuleEnabled(guildID snowflake.ID, moduleID string) bool
}

// LoggerHandler handles logger settings interactions and listens for guild events.
type LoggerHandler struct {
	bot     LoggerBot
	client  *disgobot.Client
	service *service.LoggerService
	logger  *slog.Logger
}

// NewLoggerHandler creates a new LoggerHandler.
func NewLoggerHandler(bot LoggerBot, client *disgobot.Client, svc *service.LoggerService, logger *slog.Logger) *LoggerHandler {
	return &LoggerHandler{
		bot:     bot,
		client:  client,
		service: svc,
		logger:  logger,
	}
}

func (h *LoggerHandler) Info() module.Info {
	return module.Info{
		ID:          model.LoggerModuleID,
		Name:        "Logger",
		Description: "サーバーイベントのログを記録",
		AlwaysOn:    false,
	}
}

func (h *LoggerHandler) Commands() []discord.ApplicationCommandCreate {
	return nil
}

func (h *LoggerHandler) HandleCommand(_ *events.ApplicationCommandInteractionCreate) {}

func (h *LoggerHandler) HandleModal(_ *events.ModalSubmitInteractionCreate) {}

func (h *LoggerHandler) SettingsSummary(guildID snowflake.ID) string {
	settings, err := h.service.LoadSettings(guildID)
	if err != nil {
		return ""
	}
	var parts []string
	if settings.ChannelID != 0 {
		parts = append(parts, fmt.Sprintf("ログ先: #%d", settings.ChannelID))
	}
	count := 0
	for _, enabled := range settings.Events {
		if enabled {
			count++
		}
	}
	if count > 0 {
		parts = append(parts, fmt.Sprintf("イベント: %d個", count))
	}
	if len(parts) == 0 {
		return "未設定"
	}
	return strings.Join(parts, ", ")
}

func (h *LoggerHandler) SettingsPanel(guildID snowflake.ID) []discord.LayoutComponent {
	settings, err := h.service.LoadSettings(guildID)
	if err != nil {
		h.logger.Error("failed to load logger settings", slog.Any("error", err))
		settings = &model.LoggerSettings{Events: make(map[string]bool)}
	}
	return view.LoggerSettingsPanel(settings)
}

// sendLog checks module enablement and event settings, then posts to the log channel.
func (h *LoggerHandler) sendLog(guildID snowflake.ID, event string, msg discord.MessageCreate) {
	if !h.bot.IsModuleEnabled(guildID, model.LoggerModuleID) {
		return
	}

	settings, err := h.service.LoadSettings(guildID)
	if err != nil {
		h.logger.Error("failed to load logger settings", slog.Any("error", err))
		return
	}

	if settings.ChannelID == 0 || !settings.IsEventEnabled(event) {
		return
	}

	if _, err := h.client.Rest.CreateMessage(settings.ChannelID, msg); err != nil {
		h.logger.Error("failed to send log message",
			slog.String("event", event),
			slog.Any("error", err),
		)
	}
}
