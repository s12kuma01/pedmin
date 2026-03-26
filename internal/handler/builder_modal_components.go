// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package handler

import (
	"log/slog"
	"strconv"
	"strings"

	"github.com/disgoorg/disgo/events"
	"github.com/Sumire-Labs/pedmin/internal/model"
	"github.com/Sumire-Labs/pedmin/internal/ui"
	"github.com/Sumire-Labs/pedmin/internal/view"
)

func (h *BuilderHandler) handleTextModal(e *events.ModalSubmitInteractionCreate, panelID string) {
	content := strings.TrimSpace(e.Data.Text(model.BuilderModuleID + ":text_content"))
	if content == "" {
		_ = e.CreateMessage(ui.EphemeralError("テキストを入力してください。"))
		return
	}

	id, err := strconv.ParseInt(panelID, 10, 64)
	if err != nil {
		return
	}

	comp := model.PanelComponent{
		Type:    model.PanelComponentText,
		Content: content,
	}

	panel, err := h.service.AddComponent(id, *e.GuildID(), comp)
	if err != nil {
		h.logger.Error("failed to add text", slog.Any("error", err))
		_ = e.CreateMessage(ui.EphemeralError("追加に失敗しました: " + err.Error()))
		return
	}

	_ = e.UpdateMessage(view.BuilderEditPanel(panel))
}

func (h *BuilderHandler) handleSectionModal(e *events.ModalSubmitInteractionCreate, panelID string) {
	text1 := strings.TrimSpace(e.Data.Text(model.BuilderModuleID + ":section_text1"))
	text2 := strings.TrimSpace(e.Data.Text(model.BuilderModuleID + ":section_text2"))
	text3 := strings.TrimSpace(e.Data.Text(model.BuilderModuleID + ":section_text3"))
	thumb := strings.TrimSpace(e.Data.Text(model.BuilderModuleID + ":section_thumb"))

	if text1 == "" {
		_ = e.CreateMessage(ui.EphemeralError("テキスト1は必須です。"))
		return
	}

	var texts []string
	texts = append(texts, text1)
	if text2 != "" {
		texts = append(texts, text2)
	}
	if text3 != "" {
		texts = append(texts, text3)
	}

	id, err := strconv.ParseInt(panelID, 10, 64)
	if err != nil {
		return
	}

	comp := model.PanelComponent{
		Type:         model.PanelComponentSection,
		Texts:        texts,
		ThumbnailURL: thumb,
	}

	panel, err := h.service.AddComponent(id, *e.GuildID(), comp)
	if err != nil {
		h.logger.Error("failed to add section", slog.Any("error", err))
		_ = e.CreateMessage(ui.EphemeralError("追加に失敗しました: " + err.Error()))
		return
	}

	_ = e.UpdateMessage(view.BuilderEditPanel(panel))
}

func (h *BuilderHandler) handleMediaModal(e *events.ModalSubmitInteractionCreate, panelID string) {
	url1 := strings.TrimSpace(e.Data.Text(model.BuilderModuleID + ":media_url1"))
	desc1 := strings.TrimSpace(e.Data.Text(model.BuilderModuleID + ":media_desc1"))
	url2 := strings.TrimSpace(e.Data.Text(model.BuilderModuleID + ":media_url2"))
	desc2 := strings.TrimSpace(e.Data.Text(model.BuilderModuleID + ":media_desc2"))

	if url1 == "" {
		_ = e.CreateMessage(ui.EphemeralError("画像URLを入力してください。"))
		return
	}

	var items []model.PanelMediaItem
	items = append(items, model.PanelMediaItem{URL: url1, Description: desc1})
	if url2 != "" {
		items = append(items, model.PanelMediaItem{URL: url2, Description: desc2})
	}

	id, err := strconv.ParseInt(panelID, 10, 64)
	if err != nil {
		return
	}

	comp := model.PanelComponent{
		Type:  model.PanelComponentMedia,
		Items: items,
	}

	panel, err := h.service.AddComponent(id, *e.GuildID(), comp)
	if err != nil {
		h.logger.Error("failed to add media", slog.Any("error", err))
		_ = e.CreateMessage(ui.EphemeralError("追加に失敗しました: " + err.Error()))
		return
	}

	_ = e.UpdateMessage(view.BuilderEditPanel(panel))
}

func (h *BuilderHandler) handleLinksModal(e *events.ModalSubmitInteractionCreate, panelID string) {
	input := strings.TrimSpace(e.Data.Text(model.BuilderModuleID + ":links_input"))
	if input == "" {
		_ = e.CreateMessage(ui.EphemeralError("ボタン情報を入力してください。"))
		return
	}

	var buttons []model.PanelLinkButton
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "|", 2)
		if len(parts) != 2 {
			continue
		}
		label := strings.TrimSpace(parts[0])
		url := strings.TrimSpace(parts[1])
		if label == "" || url == "" {
			continue
		}
		buttons = append(buttons, model.PanelLinkButton{Label: label, URL: url})
		if len(buttons) >= 5 {
			break
		}
	}

	if len(buttons) == 0 {
		_ = e.CreateMessage(ui.EphemeralError("有効なボタンがありません。形式: ラベル|URL"))
		return
	}

	id, err := strconv.ParseInt(panelID, 10, 64)
	if err != nil {
		return
	}

	comp := model.PanelComponent{
		Type:    model.PanelComponentLinks,
		Buttons: buttons,
	}

	panel, err := h.service.AddComponent(id, *e.GuildID(), comp)
	if err != nil {
		h.logger.Error("failed to add links", slog.Any("error", err))
		_ = e.CreateMessage(ui.EphemeralError("追加に失敗しました: " + err.Error()))
		return
	}

	_ = e.UpdateMessage(view.BuilderEditPanel(panel))
}
