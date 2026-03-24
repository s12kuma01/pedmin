// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package handler

import (
	"log/slog"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/omit"
	"github.com/disgoorg/snowflake/v2"
	"github.com/s12kuma01/pedmin/internal/model"
	"github.com/s12kuma01/pedmin/internal/module"
	"github.com/s12kuma01/pedmin/internal/ui"
)

// SettingsBot defines the bot interface needed by the settings handler.
type SettingsBot interface {
	IsModuleEnabled(guildID snowflake.ID, moduleID string) bool
	GetModules() map[string]module.Module
	SetModuleEnabled(guildID snowflake.ID, moduleID string, enabled bool) error
}

// SettingsHandler implements module.Module for the settings feature.
type SettingsHandler struct {
	bot    SettingsBot
	logger *slog.Logger
}

// NewSettingsHandler creates a new SettingsHandler.
func NewSettingsHandler(bot SettingsBot, logger *slog.Logger) *SettingsHandler {
	return &SettingsHandler{bot: bot, logger: logger}
}

func (h *SettingsHandler) Info() module.Info {
	return module.Info{
		ID:          model.SettingsModuleID,
		Name:        "設定",
		Description: "サーバー設定管理パネル",
		AlwaysOn:    true,
	}
}

func (h *SettingsHandler) Commands() []discord.ApplicationCommandCreate {
	manageGuild := discord.PermissionManageGuild
	return []discord.ApplicationCommandCreate{
		discord.SlashCommandCreate{
			Name:                     "settings",
			Description:              "サーバー設定パネルを開く",
			DefaultMemberPermissions: omit.New(&manageGuild),
		},
	}
}

func (h *SettingsHandler) HandleCommand(e *events.ApplicationCommandInteractionCreate) {
	guildID := e.GuildID()
	if guildID == nil {
		_ = e.CreateMessage(ui.ErrorMessage("設定はサーバー内でのみ使用できます。"))
		return
	}

	options := h.listModuleOptions(*guildID)
	_ = e.CreateMessage(ui.BuildMainPanel(options))
}

func (h *SettingsHandler) HandleComponent(e *events.ComponentInteractionCreate) {
	customID := e.Data.CustomID()
	_, action, _ := strings.Cut(customID, ":")

	guildID := e.GuildID()
	if guildID == nil {
		return
	}

	switch {
	case action == "select":
		data, ok := e.Data.(discord.StringSelectMenuInteractionData)
		if !ok || len(data.Values) == 0 {
			return
		}
		moduleID := data.Values[0]
		_ = e.UpdateMessage(h.buildModulePanel(*guildID, moduleID))

	case strings.HasPrefix(action, "toggle:"):
		moduleID := strings.TrimPrefix(action, "toggle:")
		enabled := h.bot.IsModuleEnabled(*guildID, moduleID)
		if err := h.bot.SetModuleEnabled(*guildID, moduleID, !enabled); err != nil {
			h.logger.Error("failed to toggle module", slog.Any("error", err))
		}
		_ = e.UpdateMessage(h.buildModulePanel(*guildID, moduleID))

	case action == "back":
		options := h.listModuleOptions(*guildID)
		_ = e.UpdateMessage(ui.BuildMainPanelUpdate(options))
	}
}

func (h *SettingsHandler) HandleModal(_ *events.ModalSubmitInteractionCreate) {}

func (h *SettingsHandler) SettingsPanel(_ snowflake.ID) []discord.LayoutComponent {
	return nil
}

func (h *SettingsHandler) listModuleOptions(guildID snowflake.ID) []ui.ModuleOption {
	modules := h.bot.GetModules()
	var options []ui.ModuleOption
	for _, m := range modules {
		info := m.Info()
		if info.AlwaysOn {
			continue
		}
		opt := ui.ModuleOption{
			ID:          info.ID,
			Name:        info.Name,
			Description: info.Description,
			Enabled:     h.bot.IsModuleEnabled(guildID, info.ID),
		}
		if summarizer, ok := m.(module.SettingsSummarizer); ok {
			opt.Summary = summarizer.SettingsSummary(guildID)
		}
		options = append(options, opt)
	}
	return options
}

func (h *SettingsHandler) buildModulePanel(guildID snowflake.ID, moduleID string) discord.MessageUpdate {
	modules := h.bot.GetModules()
	m, ok := modules[moduleID]
	if !ok {
		return ui.BuildModuleNotFound()
	}

	info := m.Info()
	enabled := h.bot.IsModuleEnabled(guildID, moduleID)
	settingsPanel := m.SettingsPanel(guildID)
	return ui.BuildModulePanel(info, enabled, settingsPanel)
}
