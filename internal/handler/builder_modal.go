// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package handler

import (
	"log/slog"
	"strconv"
	"strings"

	"github.com/disgoorg/disgo/events"
	"github.com/s12kuma01/pedmin/internal/model"
	"github.com/s12kuma01/pedmin/internal/ui"
	"github.com/s12kuma01/pedmin/internal/view"
)

func (h *BuilderHandler) HandleModal(e *events.ModalSubmitInteractionCreate) {
	customID := e.Data.CustomID
	_, rest, _ := strings.Cut(customID, ":")
	action, extra, _ := strings.Cut(rest, ":")

	guildID := e.GuildID()
	if guildID == nil {
		return
	}

	switch action {
	case "create_modal":
		h.handleCreateModal(e)
	case "text_modal":
		h.handleTextModal(e, extra)
	case "section_modal":
		h.handleSectionModal(e, extra)
	case "media_modal":
		h.handleMediaModal(e, extra)
	case "links_modal":
		h.handleLinksModal(e, extra)
	case "rename_modal":
		h.handleRenameModal(e, extra)
	}
}

func (h *BuilderHandler) handleCreateModal(e *events.ModalSubmitInteractionCreate) {
	name := strings.TrimSpace(e.Data.Text(model.BuilderModuleID + ":panel_name"))
	if name == "" {
		_ = e.CreateMessage(ui.EphemeralError("パネル名を入力してください。"))
		return
	}

	panel, err := h.service.CreatePanel(*e.GuildID(), name)
	if err != nil {
		h.logger.Error("failed to create panel", slog.Any("error", err))
		_ = e.CreateMessage(ui.EphemeralError("パネルの作成に失敗しました: " + err.Error()))
		return
	}

	_ = e.UpdateMessage(view.BuilderEditPanel(panel))
}

func (h *BuilderHandler) handleRenameModal(e *events.ModalSubmitInteractionCreate, panelID string) {
	newName := strings.TrimSpace(e.Data.Text(model.BuilderModuleID + ":new_name"))
	if newName == "" {
		_ = e.CreateMessage(ui.EphemeralError("パネル名を入力してください。"))
		return
	}

	id, err := strconv.ParseInt(panelID, 10, 64)
	if err != nil {
		return
	}

	panel, err := h.service.RenamePanel(id, *e.GuildID(), newName)
	if err != nil {
		h.logger.Error("failed to rename panel", slog.Any("error", err))
		_ = e.CreateMessage(ui.EphemeralError("名前変更に失敗しました: " + err.Error()))
		return
	}

	_ = e.UpdateMessage(view.BuilderEditPanel(panel))
}
