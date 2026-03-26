// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package handler

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
	"github.com/Sumire-Labs/pedmin/internal/model"
	"github.com/Sumire-Labs/pedmin/internal/module"
	"github.com/Sumire-Labs/pedmin/internal/service"
	"github.com/Sumire-Labs/pedmin/internal/view"
)

// RSSHandler handles RSS feed management interactions.
type RSSHandler struct {
	service      *service.RSSService
	feedTimeout_ time.Duration
	logger       *slog.Logger
}

// NewRSSHandler creates a new RSSHandler.
func NewRSSHandler(svc *service.RSSService, feedTimeout time.Duration, logger *slog.Logger) *RSSHandler {
	return &RSSHandler{
		service:      svc,
		feedTimeout_: feedTimeout,
		logger:       logger,
	}
}

func (h *RSSHandler) Info() module.Info {
	return module.Info{
		ID:          model.RSSModuleID,
		Name:        "RSS",
		Description: "RSSフィード監視",
		AlwaysOn:    false,
	}
}

func (h *RSSHandler) Commands() []discord.ApplicationCommandCreate {
	return nil
}

func (h *RSSHandler) HandleCommand(_ *events.ApplicationCommandInteractionCreate) {}

func (h *RSSHandler) SettingsSummary(guildID snowflake.ID) string {
	count, err := h.service.CountFeeds(guildID)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("フィード: %d/%d件", count, model.MaxRSSFeedsPerGuild)
}

func (h *RSSHandler) SettingsPanel(guildID snowflake.ID) []discord.LayoutComponent {
	count, err := h.service.CountFeeds(guildID)
	if err != nil {
		h.logger.Error("failed to count rss feeds", slog.Any("error", err))
	}
	return view.RSSSettingsPanel(count)
}
