// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package handler

import (
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

// AutoroleBot is the interface the handler needs from the bot registry.
type AutoroleBot interface {
	IsModuleEnabled(guildID snowflake.ID, moduleID string) bool
}

// AutoroleHandler implements module.Module for the autorole feature.
type AutoroleHandler struct {
	bot     AutoroleBot
	service *service.AutoroleService
	logger  *slog.Logger
}

// NewAutoroleHandler creates a new AutoroleHandler.
func NewAutoroleHandler(bot AutoroleBot, svc *service.AutoroleService, logger *slog.Logger) *AutoroleHandler {
	return &AutoroleHandler{
		bot:     bot,
		service: svc,
		logger:  logger,
	}
}

// SetupAutoroleListeners registers the member join listener on the Discord client.
func SetupAutoroleListeners(client *disgobot.Client, h *AutoroleHandler) {
	client.AddEventListeners(
		disgobot.NewListenerFunc(h.onMemberJoin),
	)
}

func (h *AutoroleHandler) Info() module.Info {
	return module.Info{
		ID:          model.AutoroleModuleID,
		Name:        "Autorole",
		Description: "参加時の自動ロール付与",
		AlwaysOn:    false,
	}
}

func (h *AutoroleHandler) Commands() []discord.ApplicationCommandCreate {
	return nil
}

func (h *AutoroleHandler) HandleCommand(_ *events.ApplicationCommandInteractionCreate) {}

func (h *AutoroleHandler) HandleModal(_ *events.ModalSubmitInteractionCreate) {}

func (h *AutoroleHandler) HandleComponent(e *events.ComponentInteractionCreate) {
	h.handleComponent(e)
}

func (h *AutoroleHandler) SettingsSummary(guildID snowflake.ID) string {
	settings, err := h.service.LoadSettings(guildID)
	if err != nil {
		return ""
	}
	var parts []string
	if settings.UserRoleID != 0 {
		parts = append(parts, "ユーザー: 設定済")
	}
	if settings.BotRoleID != 0 {
		parts = append(parts, "Bot: 設定済")
	}
	if len(parts) == 0 {
		return "未設定"
	}
	return joinParts(parts)
}

func (h *AutoroleHandler) SettingsPanel(guildID snowflake.ID) []discord.LayoutComponent {
	settings, err := h.service.LoadSettings(guildID)
	if err != nil {
		h.logger.Error("failed to load autorole settings", slog.Any("error", err))
		settings = model.DefaultAutoroleSettings()
	}
	return view.AutoroleSettingsPanel(settings)
}

func (h *AutoroleHandler) onMemberJoin(e *events.GuildMemberJoin) {
	if !h.bot.IsModuleEnabled(e.GuildID, model.AutoroleModuleID) {
		return
	}
	h.service.AssignRole(e.GuildID, e.Member.User.ID, e.Member.User.Bot)
}

func (h *AutoroleHandler) handleComponent(e *events.ComponentInteractionCreate) {
	customID := e.Data.CustomID()
	guildID := e.GuildID()
	if guildID == nil {
		return
	}

	settings, err := h.service.LoadSettings(*guildID)
	if err != nil {
		h.logger.Error("failed to load autorole settings", slog.Any("error", err))
		return
	}

	switch customID {
	case model.AutoroleModuleID + ":user_role":
		data, ok := e.Data.(discord.RoleSelectMenuInteractionData)
		if !ok || len(data.Values) == 0 {
			return
		}
		settings.UserRoleID = data.Values[0]

	case model.AutoroleModuleID + ":bot_role":
		data, ok := e.Data.(discord.RoleSelectMenuInteractionData)
		if !ok || len(data.Values) == 0 {
			return
		}
		settings.BotRoleID = data.Values[0]

	case model.AutoroleModuleID + ":clear_user":
		settings.UserRoleID = 0

	case model.AutoroleModuleID + ":clear_bot":
		settings.BotRoleID = 0

	default:
		return
	}

	if err := h.service.SaveSettings(*guildID, settings); err != nil {
		h.logger.Error("failed to save autorole settings", slog.Any("error", err))
	}

	settingsUI := view.AutoroleSettingsPanel(settings)
	enabled := h.bot.IsModuleEnabled(*guildID, model.AutoroleModuleID)
	_ = e.UpdateMessage(ui.BuildModulePanel(h.Info(), enabled, settingsUI))
}

func joinParts(parts []string) string {
	result := parts[0]
	for _, p := range parts[1:] {
		result += ", " + p
	}
	return result
}
