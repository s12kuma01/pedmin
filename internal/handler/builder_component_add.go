// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package handler

import (
	"log/slog"
	"strconv"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/s12kuma01/pedmin/internal/model"
	"github.com/s12kuma01/pedmin/internal/view"
)

func (h *BuilderHandler) handleAddTextPrompt(e *events.ComponentInteractionCreate, panelID string) {
	_ = e.Modal(discord.ModalCreate{
		CustomID: model.BuilderModuleID + ":text_modal:" + panelID,
		Title:    "テキスト追加",
		Components: []discord.LayoutComponent{
			discord.NewLabel("テキスト (Markdown対応)",
				discord.NewParagraphTextInput(model.BuilderModuleID+":text_content").
					WithRequired(true).WithPlaceholder("## タイトル\n本文テキスト..."),
			),
		},
	})
}

func (h *BuilderHandler) handleAddSectionPrompt(e *events.ComponentInteractionCreate, panelID string) {
	_ = e.Modal(discord.ModalCreate{
		CustomID: model.BuilderModuleID + ":section_modal:" + panelID,
		Title:    "セクション追加",
		Components: []discord.LayoutComponent{
			discord.NewLabel("テキスト1",
				discord.NewShortTextInput(model.BuilderModuleID+":section_text1").
					WithRequired(true),
			),
			discord.NewLabel("テキスト2 (任意)",
				discord.NewShortTextInput(model.BuilderModuleID+":section_text2").
					WithRequired(false),
			),
			discord.NewLabel("テキスト3 (任意)",
				discord.NewShortTextInput(model.BuilderModuleID+":section_text3").
					WithRequired(false),
			),
			discord.NewLabel("サムネイルURL (任意)",
				discord.NewShortTextInput(model.BuilderModuleID+":section_thumb").
					WithRequired(false).WithPlaceholder("https://..."),
			),
		},
	})
}

func (h *BuilderHandler) handleAddSeparator(e *events.ComponentInteractionCreate, panelID string) {
	options := []discord.StringSelectMenuOption{
		{Label: "小さいセパレータ", Value: "small"},
		{Label: "大きいセパレータ", Value: "large"},
	}
	_ = e.UpdateMessage(discord.NewMessageUpdateV2([]discord.LayoutComponent{
		discord.NewContainer(
			discord.NewTextDisplay("セパレータの種類を選択:"),
			discord.NewActionRow(
				discord.NewStringSelectMenu(model.BuilderModuleID+":sep_select:"+panelID, "種類を選択...", options...),
			),
		),
	}))
}

func (h *BuilderHandler) handleSepSelect(e *events.ComponentInteractionCreate, panelID string) {
	data, ok := e.Data.(discord.StringSelectMenuInteractionData)
	if !ok || len(data.Values) == 0 {
		return
	}

	id, err := strconv.ParseInt(panelID, 10, 64)
	if err != nil {
		return
	}

	comp := model.PanelComponent{
		Type:    model.PanelComponentSeparator,
		Spacing: data.Values[0],
	}

	panel, err := h.service.AddComponent(id, *e.GuildID(), comp)
	if err != nil {
		h.logger.Error("failed to add separator", slog.Any("error", err))
		return
	}

	_ = e.UpdateMessage(view.BuilderEditPanel(panel))
}

func (h *BuilderHandler) handleAddMediaPrompt(e *events.ComponentInteractionCreate, panelID string) {
	_ = e.Modal(discord.ModalCreate{
		CustomID: model.BuilderModuleID + ":media_modal:" + panelID,
		Title:    "画像追加",
		Components: []discord.LayoutComponent{
			discord.NewLabel("画像URL 1",
				discord.NewShortTextInput(model.BuilderModuleID+":media_url1").
					WithRequired(true).WithPlaceholder("https://..."),
			),
			discord.NewLabel("説明 1 (任意)",
				discord.NewShortTextInput(model.BuilderModuleID+":media_desc1").
					WithRequired(false),
			),
			discord.NewLabel("画像URL 2 (任意)",
				discord.NewShortTextInput(model.BuilderModuleID+":media_url2").
					WithRequired(false).WithPlaceholder("https://..."),
			),
			discord.NewLabel("説明 2 (任意)",
				discord.NewShortTextInput(model.BuilderModuleID+":media_desc2").
					WithRequired(false),
			),
		},
	})
}

func (h *BuilderHandler) handleAddLinksPrompt(e *events.ComponentInteractionCreate, panelID string) {
	_ = e.Modal(discord.ModalCreate{
		CustomID: model.BuilderModuleID + ":links_modal:" + panelID,
		Title:    "リンクボタン追加",
		Components: []discord.LayoutComponent{
			discord.NewLabel("ボタン (1行1個: ラベル|URL)",
				discord.NewParagraphTextInput(model.BuilderModuleID+":links_input").
					WithRequired(true).WithPlaceholder("公式サイト|https://example.com\nDiscord|https://discord.gg/..."),
			),
		},
	})
}
