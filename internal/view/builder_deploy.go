// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package view

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
	"github.com/Sumire-Labs/pedmin/internal/model"
	"github.com/Sumire-Labs/pedmin/internal/ui"
)

// BuilderDeployPrompt builds the channel selection prompt.
func BuilderDeployPrompt(panelID int64) discord.MessageCreate {
	pid := fmt.Sprintf("%d", panelID)
	return ui.EphemeralV2(
		discord.NewContainer(
			discord.NewTextDisplay("パネルを配信するチャンネルを選択してください:"),
			discord.NewActionRow(
				discord.NewChannelSelectMenu(model.BuilderModuleID+":deploy_channel:"+pid, "チャンネルを選択...").
					WithChannelTypes(discord.ChannelTypeGuildText),
			),
		),
	)
}

// BuilderDeployConfirm builds the deploy confirmation.
func BuilderDeployConfirm(panelID int64, channelID snowflake.ID) discord.MessageUpdate {
	pid := fmt.Sprintf("%d", panelID)
	return discord.NewMessageUpdateV2([]discord.LayoutComponent{
		discord.NewContainer(
			discord.NewTextDisplay(fmt.Sprintf("配信先: <#%d>", channelID)),
			discord.NewActionRow(
				discord.NewSuccessButton("配信する", fmt.Sprintf("%s:deploy_confirm:%s:%d", model.BuilderModuleID, pid, channelID)),
				discord.NewSecondaryButton("キャンセル", model.BuilderModuleID+":back"),
			),
		),
	})
}

// BuilderDeleteConfirm builds the delete confirmation.
func BuilderDeleteConfirm(panelID int64, panelName string) discord.MessageUpdate {
	pid := fmt.Sprintf("%d", panelID)
	return discord.NewMessageUpdateV2([]discord.LayoutComponent{
		discord.NewContainer(
			discord.NewTextDisplay(fmt.Sprintf("**%s** を削除しますか？", panelName)),
			discord.NewActionRow(
				discord.NewDangerButton("削除する", model.BuilderModuleID+":delete_confirm:"+pid),
				discord.NewSecondaryButton("キャンセル", fmt.Sprintf("%s:select:%s", model.BuilderModuleID, pid)),
			),
		),
	})
}
