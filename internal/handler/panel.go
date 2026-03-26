// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package handler

import (
	"context"
	"log/slog"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
	"github.com/Sumire-Labs/pedmin/config"
	"github.com/Sumire-Labs/pedmin/internal/model"
	"github.com/Sumire-Labs/pedmin/internal/module"
	"github.com/Sumire-Labs/pedmin/internal/service"
	"github.com/Sumire-Labs/pedmin/internal/ui"
	"github.com/Sumire-Labs/pedmin/internal/view"
)

// PanelHandler handles /panel command and component interactions.
type PanelHandler struct {
	cfg     *config.Config
	service *service.PanelService
	logger  *slog.Logger
}

// NewPanelHandler creates a new PanelHandler.
func NewPanelHandler(cfg *config.Config, svc *service.PanelService, logger *slog.Logger) *PanelHandler {
	return &PanelHandler{
		cfg:     cfg,
		service: svc,
		logger:  logger,
	}
}

func (h *PanelHandler) Info() module.Info {
	return module.Info{
		ID:          model.PanelModuleID,
		Name:        "Panel",
		Description: "ゲームサーバー管理",
		AlwaysOn:    true,
	}
}

func (h *PanelHandler) Commands() []discord.ApplicationCommandCreate {
	return []discord.ApplicationCommandCreate{
		discord.SlashCommandCreate{
			Name:        "panel",
			Description: "ゲームサーバーを管理する",
		},
	}
}

func (h *PanelHandler) SettingsPanel(_ snowflake.ID) []discord.LayoutComponent { return nil }

func (h *PanelHandler) HandleCommand(e *events.ApplicationCommandInteractionCreate) {
	userID := e.User().ID

	if !h.isAllowed(userID) {
		_ = e.CreateMessage(ui.ErrorMessage("このコマンドを使用する権限がありません。"))
		return
	}

	if h.cfg.PanelURL == "" || h.cfg.PanelAPIKey == "" {
		_ = e.CreateMessage(ui.ErrorMessage("パネルが設定されていません。"))
		return
	}

	_ = e.DeferCreateMessage(false)

	ctx, cancel := context.WithTimeout(context.Background(), config.DefaultHTTPClientTimeout)
	defer cancel()

	servers, err := h.service.ListServersWithStatus(ctx)
	if err != nil {
		h.logger.Error("failed to list servers", slog.Any("error", err))
		_, _ = e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(), view.PanelErrorPanel(err.Error()))
		return
	}

	msg := view.PanelServerList(servers)
	_, _ = e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(), discord.NewMessageUpdateV2(msg.Components))
}

func (h *PanelHandler) isAllowed(userID snowflake.ID) bool {
	for _, id := range h.cfg.PanelAllowedUsers {
		if id == userID {
			return true
		}
	}
	return false
}
