// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package view

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/s12kuma01/pedmin/internal/model"
	"github.com/s12kuma01/pedmin/internal/ui"
)

// BuilderManagePanel builds the component management panel with a select menu.
func BuilderManagePanel(panel *model.ComponentPanel) discord.MessageCreate {
	pid := fmt.Sprintf("%d", panel.ID)

	if len(panel.Components) == 0 {
		return ui.EphemeralV2(BuilderErrorContainer("コンポーネントがありません。"))
	}

	var options []discord.StringSelectMenuOption
	for i, comp := range panel.Components {
		options = append(options, discord.StringSelectMenuOption{
			Label:       fmt.Sprintf("%d. %s", i+1, ComponentTypeName(comp.Type)),
			Value:       fmt.Sprintf("%d", i),
			Description: truncateBuilderStr(ComponentSummary(comp), 100),
		})
	}

	return ui.EphemeralV2(
		discord.NewContainer(
			discord.NewTextDisplay(fmt.Sprintf("### %s — コンポーネント管理", panel.Name)),
			discord.NewSmallSeparator(),
			discord.NewActionRow(
				discord.NewStringSelectMenu(model.BuilderModuleID+":manage_select:"+pid, "コンポーネントを選択...", options...),
			),
		),
	)
}

// BuilderComponentDetail builds the detail view for a component with move/delete buttons.
func BuilderComponentDetail(panel *model.ComponentPanel, index int) discord.MessageUpdate {
	pid := fmt.Sprintf("%d", panel.ID)
	comp := panel.Components[index]

	upBtn := discord.NewSecondaryButton("↑ 上に移動", fmt.Sprintf("%s:move_up:%s:%d", model.BuilderModuleID, pid, index))
	if index == 0 {
		upBtn = upBtn.AsDisabled()
	}
	downBtn := discord.NewSecondaryButton("↓ 下に移動", fmt.Sprintf("%s:move_down:%s:%d", model.BuilderModuleID, pid, index))
	if index >= len(panel.Components)-1 {
		downBtn = downBtn.AsDisabled()
	}
	deleteBtn := discord.NewDangerButton("削除", fmt.Sprintf("%s:delete_comp:%s:%d", model.BuilderModuleID, pid, index))

	return discord.NewMessageUpdateV2([]discord.LayoutComponent{
		discord.NewContainer(
			discord.NewTextDisplay(fmt.Sprintf("### コンポーネント %d", index+1)),
			discord.NewSmallSeparator(),
			discord.NewTextDisplay(fmt.Sprintf("**タイプ:** %s\n**内容:** %s", ComponentTypeName(comp.Type), ComponentSummary(comp))),
			discord.NewLargeSeparator(),
			discord.NewActionRow(upBtn, downBtn, deleteBtn),
		),
	})
}
