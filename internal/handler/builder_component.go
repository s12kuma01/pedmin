// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package handler

import (
	"log/slog"
	"strconv"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/s12kuma01/pedmin/internal/model"
	"github.com/s12kuma01/pedmin/internal/view"
)

func (h *BuilderHandler) HandleComponent(e *events.ComponentInteractionCreate) {
	customID := e.Data.CustomID()
	_, rest, _ := strings.Cut(customID, ":")
	action, extra, _ := strings.Cut(rest, ":")

	guildID := e.GuildID()
	if guildID == nil {
		return
	}

	switch action {
	case "list":
		h.handleList(e)
	case "create_prompt":
		h.handleCreatePrompt(e)
	case "select":
		h.handleSelect(e, extra)
	case "add_text":
		h.handleAddTextPrompt(e, extra)
	case "add_section":
		h.handleAddSectionPrompt(e, extra)
	case "add_separator":
		h.handleAddSeparator(e, extra)
	case "sep_select":
		h.handleSepSelect(e, extra)
	case "add_media":
		h.handleAddMediaPrompt(e, extra)
	case "add_links":
		h.handleAddLinksPrompt(e, extra)
	case "manage":
		h.handleManage(e, extra)
	case "manage_select":
		h.handleManageSelect(e, extra)
	case "delete_comp":
		h.handleDeleteComponent(e, extra)
	case "move_up":
		h.handleMoveUp(e, extra)
	case "move_down":
		h.handleMoveDown(e, extra)
	case "preview":
		h.handlePreview(e, extra)
	case "deploy_prompt":
		h.handleDeployPrompt(e, extra)
	case "deploy_channel":
		h.handleDeployChannel(e, extra)
	case "deploy_confirm":
		h.handleDeployConfirm(e, extra)
	case "delete_panel":
		h.handleDeletePanel(e, extra)
	case "delete_confirm":
		h.handleDeleteConfirm(e, extra)
	case "rename":
		h.handleRenamePrompt(e, extra)
	case "back":
		h.handleList(e)
	}
}

func (h *BuilderHandler) handleList(e *events.ComponentInteractionCreate) {
	panels, err := h.service.GetPanels(*e.GuildID())
	if err != nil {
		h.logger.Error("failed to get panels", slog.Any("error", err))
		return
	}
	_ = e.UpdateMessage(view.BuilderListPanelUpdate(panels, len(panels)))
}

func (h *BuilderHandler) handleCreatePrompt(e *events.ComponentInteractionCreate) {
	_ = e.Modal(discord.ModalCreate{
		CustomID: model.BuilderModuleID + ":create_modal",
		Title:    "パネル作成",
		Components: []discord.LayoutComponent{
			discord.NewLabel("パネル名",
				discord.NewShortTextInput(model.BuilderModuleID+":panel_name").
					WithRequired(true).WithPlaceholder("ルール・ウェルカム等"),
			),
		},
	})
}

func (h *BuilderHandler) handleSelect(e *events.ComponentInteractionCreate, extra string) {
	panelID := extra
	if panelID == "" {
		data, ok := e.Data.(discord.StringSelectMenuInteractionData)
		if !ok || len(data.Values) == 0 {
			return
		}
		panelID = data.Values[0]
	}

	id, err := strconv.ParseInt(panelID, 10, 64)
	if err != nil {
		return
	}

	panel, err := h.service.GetPanel(id, *e.GuildID())
	if err != nil {
		h.logger.Error("failed to get panel", slog.Any("error", err))
		return
	}

	_ = e.UpdateMessage(view.BuilderEditPanel(panel))
}
