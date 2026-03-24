// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package view

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
	"github.com/s12kuma01/pedmin/internal/model"
)

// TicketSettingsPanel builds the ticket settings panel components.
func TicketSettingsPanel(settings *model.TicketSettings) []discord.LayoutComponent {
	categoryText := "未設定"
	if settings.CategoryID != 0 {
		categoryText = fmt.Sprintf("<#%d>", settings.CategoryID)
	}
	logText := "未設定"
	if settings.LogChannelID != 0 {
		logText = fmt.Sprintf("<#%d>", settings.LogChannelID)
	}
	roleText := "未設定"
	if settings.SupportRoleID != 0 {
		roleText = fmt.Sprintf("<@&%d>", settings.SupportRoleID)
	}

	infoDisplay := discord.NewTextDisplay(fmt.Sprintf(
		"**カテゴリ:** %s\n**ログチャンネル:** %s\n**サポートロール:** %s",
		categoryText, logText, roleText,
	))

	categorySelect := discord.NewActionRow(
		discord.NewChannelSelectMenu(model.TicketModuleID+":category", "カテゴリを選択...").
			WithChannelTypes(discord.ChannelTypeGuildCategory),
	)

	buttons := discord.NewActionRow(
		discord.NewSecondaryButton("ログ設定", model.TicketModuleID+":log_prompt"),
		discord.NewSecondaryButton("ロール設定", model.TicketModuleID+":role_prompt"),
		discord.NewSuccessButton("パネル設置", model.TicketModuleID+":deploy_prompt"),
	)

	return []discord.LayoutComponent{
		infoDisplay,
		categorySelect,
		buttons,
	}
}

// TicketDeployConfirm builds the deploy confirmation panel.
func TicketDeployConfirm(channelID snowflake.ID) discord.MessageUpdate {
	return discord.NewMessageUpdateV2([]discord.LayoutComponent{
		discord.NewContainer(
			discord.NewTextDisplay("パネルを設置するチャンネル: <#"+channelID.String()+">"),
			discord.NewActionRow(
				discord.NewSuccessButton("設置する", model.TicketModuleID+":deploy_confirm:"+channelID.String()),
				discord.NewSecondaryButton("キャンセル", model.TicketModuleID+":deploy_cancel"),
			),
		),
	})
}
