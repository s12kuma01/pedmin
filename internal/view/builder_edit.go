// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package view

import (
	"fmt"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/Sumire-Labs/pedmin/internal/model"
)

// BuilderEditPanel builds the edit mode view for a panel.
func BuilderEditPanel(panel *model.ComponentPanel) discord.MessageUpdate {
	pid := fmt.Sprintf("%d", panel.ID)

	components := []discord.ContainerSubComponent{
		discord.NewTextDisplay(fmt.Sprintf("### %s", panel.Name)),
		discord.NewSmallSeparator(),
		discord.NewTextDisplay(fmt.Sprintf("**コンポーネント:** %d/%d", len(panel.Components), model.MaxComponentsPerPanel)),
	}

	// Component listing
	if len(panel.Components) > 0 {
		var sb strings.Builder
		for i, comp := range panel.Components {
			fmt.Fprintf(&sb, "%d. %s: %s\n", i+1, ComponentTypeName(comp.Type), ComponentSummary(comp))
		}
		components = append(components, discord.NewTextDisplay(sb.String()))
	}

	components = append(components, discord.NewLargeSeparator())

	// Add buttons (row 1)
	addDisabled := len(panel.Components) >= model.MaxComponentsPerPanel
	textBtn := discord.NewSecondaryButton("テキスト", model.BuilderModuleID+":add_text:"+pid)
	sectionBtn := discord.NewSecondaryButton("セクション", model.BuilderModuleID+":add_section:"+pid)
	separatorBtn := discord.NewSecondaryButton("セパレータ", model.BuilderModuleID+":add_separator:"+pid)
	mediaBtn := discord.NewSecondaryButton("画像", model.BuilderModuleID+":add_media:"+pid)
	linksBtn := discord.NewSecondaryButton("リンクボタン", model.BuilderModuleID+":add_links:"+pid)
	if addDisabled {
		textBtn = textBtn.AsDisabled()
		sectionBtn = sectionBtn.AsDisabled()
		separatorBtn = separatorBtn.AsDisabled()
		mediaBtn = mediaBtn.AsDisabled()
		linksBtn = linksBtn.AsDisabled()
	}
	components = append(components,
		discord.NewActionRow(textBtn, sectionBtn, separatorBtn, mediaBtn, linksBtn),
	)

	// Action buttons (row 2)
	manageBtn := discord.NewSecondaryButton("管理", model.BuilderModuleID+":manage:"+pid)
	previewBtn := discord.NewSecondaryButton("プレビュー", model.BuilderModuleID+":preview:"+pid)
	if len(panel.Components) == 0 {
		manageBtn = manageBtn.AsDisabled()
		previewBtn = previewBtn.AsDisabled()
	}
	components = append(components,
		discord.NewActionRow(
			manageBtn,
			previewBtn,
			discord.NewSuccessButton("配信", model.BuilderModuleID+":deploy_prompt:"+pid),
			discord.NewSecondaryButton("名前変更", model.BuilderModuleID+":rename:"+pid),
			discord.NewDangerButton("削除", model.BuilderModuleID+":delete_panel:"+pid),
		),
	)

	// Back button
	components = append(components,
		discord.NewActionRow(
			discord.NewSecondaryButton("← 戻る", model.BuilderModuleID+":back"),
		),
	)

	return discord.NewMessageUpdateV2([]discord.LayoutComponent{
		discord.NewContainer(components...),
	})
}
