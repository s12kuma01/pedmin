// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package handler

import (
	"context"
	"log/slog"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/s12kuma01/pedmin/config"
	"github.com/s12kuma01/pedmin/internal/ui"
	"github.com/s12kuma01/pedmin/internal/view"
)

func (h *PanelHandler) HandleComponent(e *events.ComponentInteractionCreate) {
	if !h.isAllowed(e.User().ID) {
		_ = e.CreateMessage(ui.ErrorMessage("このコマンドを使用する権限がありません。"))
		return
	}

	customID := e.Data.CustomID()
	_, rest, _ := strings.Cut(customID, ":")
	action, extra, _ := strings.Cut(rest, ":")

	switch action {
	case "select":
		h.panelHandleSelect(e)
	case "power_start":
		h.panelHandlePower(e, extra, "start")
	case "power_restart":
		h.panelHandlePower(e, extra, "restart")
	case "power_stop":
		h.panelHandlePower(e, extra, "stop")
	case "refresh":
		h.panelHandleRefresh(e, extra)
	case "back":
		h.panelHandleBack(e)
	case "refresh_list":
		h.panelHandleBack(e)
	case "console":
		h.panelHandleConsolePrompt(e, extra)
	}
}

func (h *PanelHandler) panelHandleSelect(e *events.ComponentInteractionCreate) {
	data, ok := e.Data.(discord.StringSelectMenuInteractionData)
	if !ok || len(data.Values) == 0 {
		return
	}
	identifier := data.Values[0]

	_ = e.DeferUpdateMessage()

	ctx, cancel := context.WithTimeout(context.Background(), config.DefaultHTTPClientTimeout)
	defer cancel()

	server, res, err := h.service.GetServerDetail(ctx, identifier)
	if err != nil {
		h.logger.Error("failed to get server detail", slog.Any("error", err))
		_, _ = e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(), view.PanelErrorPanel(err.Error()))
		return
	}

	_, _ = e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(), view.PanelServerDetail(*server, res))
}

func (h *PanelHandler) panelHandlePower(e *events.ComponentInteractionCreate, identifier, signal string) {
	_ = e.DeferUpdateMessage()

	ctx, cancel := context.WithTimeout(context.Background(), config.DefaultPanelPowerTimeout)
	defer cancel()

	server, res, err := h.service.PowerAction(ctx, identifier, signal)
	if err != nil {
		h.logger.Error("failed to send power action", slog.String("signal", signal), slog.Any("error", err))
		_, _ = e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(), view.PanelErrorPanel(err.Error()))
		return
	}

	_, _ = e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(), view.PanelServerDetail(*server, res))
}

func (h *PanelHandler) panelHandleRefresh(e *events.ComponentInteractionCreate, identifier string) {
	_ = e.DeferUpdateMessage()

	ctx, cancel := context.WithTimeout(context.Background(), config.DefaultHTTPClientTimeout)
	defer cancel()

	server, res, err := h.service.GetServerDetail(ctx, identifier)
	if err != nil {
		h.logger.Error("failed to refresh server detail", slog.Any("error", err))
		_, _ = e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(), view.PanelErrorPanel(err.Error()))
		return
	}

	_, _ = e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(), view.PanelServerDetail(*server, res))
}

func (h *PanelHandler) panelHandleBack(e *events.ComponentInteractionCreate) {
	_ = e.DeferUpdateMessage()

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

func (h *PanelHandler) panelHandleConsolePrompt(e *events.ComponentInteractionCreate, identifier string) {
	_ = e.Modal(view.PanelConsoleModal(identifier))
}
