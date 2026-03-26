// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package view

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/Sumire-Labs/pedmin/internal/model"
)

// RSSSettingsPanel builds the RSS settings panel components.
func RSSSettingsPanel(feedCount int) []discord.LayoutComponent {
	infoDisplay := discord.NewTextDisplay(
		fmt.Sprintf("**登録フィード:** %d/%d", feedCount, model.MaxRSSFeedsPerGuild),
	)

	manageBtn := discord.NewSecondaryButton("フィード管理", model.RSSModuleID+":manage")
	if feedCount == 0 {
		manageBtn = manageBtn.AsDisabled()
	}

	actionRow := discord.NewActionRow(
		discord.NewPrimaryButton("フィード追加", model.RSSModuleID+":add_prompt"),
		manageBtn,
	)

	return []discord.LayoutComponent{
		infoDisplay,
		actionRow,
	}
}
