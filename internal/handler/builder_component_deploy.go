// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package handler

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
	"github.com/Sumire-Labs/pedmin/internal/model"
	"github.com/Sumire-Labs/pedmin/internal/ui"
	"github.com/Sumire-Labs/pedmin/internal/view"
)

func (h *BuilderHandler) handlePreview(e *events.ComponentInteractionCreate, panelID string) {
	id, err := strconv.ParseInt(panelID, 10, 64)
	if err != nil {
		return
	}

	panel, err := h.service.GetPanel(id, *e.GuildID())
	if err != nil {
		h.logger.Error("failed to get panel", slog.Any("error", err))
		return
	}

	_ = e.CreateMessage(h.service.PreviewPanel(panel))
}

func (h *BuilderHandler) handleDeployPrompt(e *events.ComponentInteractionCreate, panelID string) {
	id, err := strconv.ParseInt(panelID, 10, 64)
	if err != nil {
		return
	}
	_ = e.CreateMessage(view.BuilderDeployPrompt(id))
}

func (h *BuilderHandler) handleDeployChannel(e *events.ComponentInteractionCreate, panelID string) {
	data, ok := e.Data.(discord.ChannelSelectMenuInteractionData)
	if !ok || len(data.Values) == 0 {
		return
	}

	id, err := strconv.ParseInt(panelID, 10, 64)
	if err != nil {
		return
	}

	_ = e.UpdateMessage(view.BuilderDeployConfirm(id, data.Values[0]))
}

func (h *BuilderHandler) handleDeployConfirm(e *events.ComponentInteractionCreate, extra string) {
	panelID, channelIDStr, _ := strings.Cut(extra, ":")
	id, err := strconv.ParseInt(panelID, 10, 64)
	if err != nil {
		return
	}
	channelID, err := snowflake.Parse(channelIDStr)
	if err != nil {
		return
	}

	panel, err := h.service.GetPanel(id, *e.GuildID())
	if err != nil {
		h.logger.Error("failed to get panel", slog.Any("error", err))
		return
	}

	if err := h.service.DeployPanel(panel, channelID); err != nil {
		h.logger.Error("failed to deploy panel", slog.Any("error", err))
		_ = e.UpdateMessage(discord.NewMessageUpdateV2([]discord.LayoutComponent{
			view.BuilderErrorContainer(fmt.Sprintf("配信に失敗しました: %s", err.Error())),
		}))
		return
	}

	_ = e.UpdateMessage(discord.NewMessageUpdateV2([]discord.LayoutComponent{
		discord.NewContainer(
			discord.NewTextDisplay(fmt.Sprintf("**%s** を <#%d> に配信しました。", panel.Name, channelID)),
		),
	}))
}

func (h *BuilderHandler) handleDeletePanel(e *events.ComponentInteractionCreate, panelID string) {
	id, err := strconv.ParseInt(panelID, 10, 64)
	if err != nil {
		return
	}

	panel, err := h.service.GetPanel(id, *e.GuildID())
	if err != nil {
		return
	}

	_ = e.UpdateMessage(view.BuilderDeleteConfirm(id, panel.Name))
}

func (h *BuilderHandler) handleDeleteConfirm(e *events.ComponentInteractionCreate, panelID string) {
	id, err := strconv.ParseInt(panelID, 10, 64)
	if err != nil {
		return
	}

	if err := h.service.DeletePanel(id, *e.GuildID()); err != nil {
		h.logger.Error("failed to delete panel", slog.Any("error", err))
		_ = e.CreateMessage(ui.EphemeralError("パネルの削除に失敗しました。"))
		return
	}

	h.handleList(e)
}

func (h *BuilderHandler) handleRenamePrompt(e *events.ComponentInteractionCreate, panelID string) {
	_ = e.Modal(discord.ModalCreate{
		CustomID: model.BuilderModuleID + ":rename_modal:" + panelID,
		Title:    "パネル名変更",
		Components: []discord.LayoutComponent{
			discord.NewLabel("新しいパネル名",
				discord.NewShortTextInput(model.BuilderModuleID+":new_name").
					WithRequired(true),
			),
		},
	})
}
