// opyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package handler

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
	"github.com/s12kuma01/pedmin/internal/model"
	"github.com/s12kuma01/pedmin/internal/module"
	"github.com/s12kuma01/pedmin/internal/service"
	"github.com/s12kuma01/pedmin/internal/view"
)

// TicketBot is the interface for checking module enabled status.
type TicketBot interface {
	IsModuleEnabled(guildID snowflake.ID, moduleID string) bool
}

// TicketHandler handles ticket system interactions.
type TicketHandler struct {
	bot     TicketBot
	service *service.TicketService
	logger  *slog.Logger
}

// NewTicketHandler creates a new TicketHandler.
func NewTicketHandler(bot TicketBot, svc *service.TicketService, logger *slog.Logger) *TicketHandler {
	return &TicketHandler{
		bot:     bot,
		service: svc,
		logger:  logger,
	}
}

func (h *TicketHandler) Info() module.Info {
	return module.Info{
		ID:          model.TicketModuleID,
		Name:        "チケット",
		Description: "サポートチケットシステム",
		AlwaysOn:    false,
	}
}

func (h *TicketHandler) Commands() []discord.ApplicationCommandCreate {
	return nil
}

func (h *TicketHandler) HandleCommand(_ *events.ApplicationCommandInteractionCreate) {}

func (h *TicketHandler) SettingsSummary(guildID snowflake.ID) string {
	settings, err := h.service.LoadSettings(guildID)
	if err != nil {
		return ""
	}
	var parts []string
	if settings.CategoryID != 0 {
		parts = append(parts, fmt.Sprintf("カテゴリ: #%d", settings.CategoryID))
	}
	if settings.LogChannelID != 0 {
		parts = append(parts, fmt.Sprintf("ログ: #%d", settings.LogChannelID))
	}
	if len(parts) == 0 {
		return "未設定"
	}
	return strings.Join(parts, ", ")
}

func (h *TicketHandler) SettingsPanel(guildID snowflake.ID) []discord.LayoutComponent {
	settings, err := h.service.LoadSettings(guildID)
	if err != nil {
		h.logger.Error("failed to load ticket settings", slog.Any("error", err))
		settings = &model.TicketSettings{}
	}
	return view.TicketSettingsPanel(settings)
}
