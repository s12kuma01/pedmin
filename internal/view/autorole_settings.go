// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package view

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/Sumire-Labs/pedmin/internal/model"
)

// AutoroleSettingsPanel builds the autorole settings panel components.
func AutoroleSettingsPanel(settings *model.AutoroleSettings) []discord.LayoutComponent {
	userRole := "未設定"
	if settings.UserRoleID != 0 {
		userRole = fmt.Sprintf("<@&%d>", settings.UserRoleID)
	}

	botRole := "未設定"
	if settings.BotRoleID != 0 {
		botRole = fmt.Sprintf("<@&%d>", settings.BotRoleID)
	}

	infoDisplay := discord.NewTextDisplay(fmt.Sprintf(
		"**ユーザー用ロール:** %s\n**Bot用ロール:** %s",
		userRole, botRole,
	))

	userSelect := discord.NewActionRow(
		discord.NewRoleSelectMenu(model.AutoroleModuleID+":user_role", "ユーザー用ロールを選択..."),
	)

	botSelect := discord.NewActionRow(
		discord.NewRoleSelectMenu(model.AutoroleModuleID+":bot_role", "Bot用ロールを選択..."),
	)

	clearRow := discord.NewActionRow(
		discord.NewSecondaryButton("ユーザー用をリセット", model.AutoroleModuleID+":clear_user"),
		discord.NewSecondaryButton("Bot用をリセット", model.AutoroleModuleID+":clear_bot"),
	)

	return []discord.LayoutComponent{
		infoDisplay,
		userSelect,
		botSelect,
		clearRow,
	}
}
