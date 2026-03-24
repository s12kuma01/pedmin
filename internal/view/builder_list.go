// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package view

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/s12kuma01/pedmin/internal/model"
	"github.com/s12kuma01/pedmin/internal/ui"
)

// BuilderListPanel builds the panel list view.
func BuilderListPanel(panels []model.ComponentPanel, count int) discord.MessageCreate {
	return ui.EphemeralV2(builderListContainer(panels, count))
}

// BuilderListPanelUpdate builds the panel list as a message update.
func BuilderListPanelUpdate(panels []model.ComponentPanel, count int) discord.MessageUpdate {
	return discord.NewMessageUpdateV2([]discord.LayoutComponent{builderListContainer(panels, count)})
}

func builderListContainer(panels []model.ComponentPanel, count int) discord.ContainerComponent {
	components := []discord.ContainerSubComponent{
		discord.NewTextDisplay("## Component Builder"),
		discord.NewSmallSeparator(),
		discord.NewTextDisplay(fmt.Sprintf("**パネル:** %d/%d", count, model.MaxPanelsPerGuild)),
	}

	if len(panels) > 0 {
		var options []discord.StringSelectMenuOption
		for _, p := range panels {
			options = append(options, discord.StringSelectMenuOption{
				Label:       p.Name,
				Value:       fmt.Sprintf("%d", p.ID),
				Description: fmt.Sprintf("コンポーネント: %d個", len(p.Components)),
			})
		}
		components = append(components,
			discord.NewActionRow(
				discord.NewStringSelectMenu(model.BuilderModuleID+":select", "パネルを選択...", options...),
			),
		)
	}

	createBtn := discord.NewPrimaryButton("新規作成", model.BuilderModuleID+":create_prompt")
	if count >= model.MaxPanelsPerGuild {
		createBtn = createBtn.AsDisabled()
	}
	components = append(components, discord.NewActionRow(createBtn))

	return discord.NewContainer(components...)
}
