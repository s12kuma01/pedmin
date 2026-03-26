// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package handler

import (
	"log/slog"
	"strconv"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/Sumire-Labs/pedmin/internal/view"
)

func (h *BuilderHandler) handleManage(e *events.ComponentInteractionCreate, panelID string) {
	id, err := strconv.ParseInt(panelID, 10, 64)
	if err != nil {
		return
	}

	panel, err := h.service.GetPanel(id, *e.GuildID())
	if err != nil {
		h.logger.Error("failed to get panel", slog.Any("error", err))
		return
	}

	_ = e.CreateMessage(view.BuilderManagePanel(panel))
}

func (h *BuilderHandler) handleManageSelect(e *events.ComponentInteractionCreate, panelID string) {
	data, ok := e.Data.(discord.StringSelectMenuInteractionData)
	if !ok || len(data.Values) == 0 {
		return
	}

	id, err := strconv.ParseInt(panelID, 10, 64)
	if err != nil {
		return
	}

	index, err := strconv.Atoi(data.Values[0])
	if err != nil {
		return
	}

	panel, err := h.service.GetPanel(id, *e.GuildID())
	if err != nil || index >= len(panel.Components) {
		return
	}

	_ = e.UpdateMessage(view.BuilderComponentDetail(panel, index))
}

func (h *BuilderHandler) handleDeleteComponent(e *events.ComponentInteractionCreate, extra string) {
	panelID, indexStr, _ := strings.Cut(extra, ":")
	id, err := strconv.ParseInt(panelID, 10, 64)
	if err != nil {
		return
	}
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		return
	}

	panel, err := h.service.RemoveComponent(id, *e.GuildID(), index)
	if err != nil {
		h.logger.Error("failed to remove component", slog.Any("error", err))
		return
	}

	if len(panel.Components) == 0 {
		_ = e.UpdateMessage(view.BuilderEditPanel(panel))
		return
	}

	msg := view.BuilderManagePanel(panel)
	_ = e.UpdateMessage(discord.NewMessageUpdateV2(msg.Components))
}

func (h *BuilderHandler) handleMoveUp(e *events.ComponentInteractionCreate, extra string) {
	h.handleMove(e, extra, -1)
}

func (h *BuilderHandler) handleMoveDown(e *events.ComponentInteractionCreate, extra string) {
	h.handleMove(e, extra, 1)
}

func (h *BuilderHandler) handleMove(e *events.ComponentInteractionCreate, extra string, direction int) {
	panelID, indexStr, _ := strings.Cut(extra, ":")
	id, err := strconv.ParseInt(panelID, 10, 64)
	if err != nil {
		return
	}
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		return
	}

	panel, err := h.service.MoveComponent(id, *e.GuildID(), index, index+direction)
	if err != nil {
		h.logger.Error("failed to move component", slog.Any("error", err))
		return
	}

	_ = e.UpdateMessage(view.BuilderComponentDetail(panel, index+direction))
}
