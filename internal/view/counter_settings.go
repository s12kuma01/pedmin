// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package view

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/Sumire-Labs/pedmin/internal/model"
)

// CounterSettingsPanel builds the counter settings panel components.
func CounterSettingsPanel(counterCount int) []discord.LayoutComponent {
	infoDisplay := discord.NewTextDisplay(
		fmt.Sprintf("**登録カウンター:** %d/%d", counterCount, model.MaxCountersPerGuild),
	)

	addBtn := discord.NewPrimaryButton("カウンター追加", model.CounterModuleID+":add_prompt")
	if counterCount >= model.MaxCountersPerGuild {
		addBtn = addBtn.AsDisabled()
	}

	manageBtn := discord.NewSecondaryButton("カウンター管理", model.CounterModuleID+":manage")
	if counterCount == 0 {
		manageBtn = manageBtn.AsDisabled()
	}

	statsBtn := discord.NewSecondaryButton("統計", model.CounterModuleID+":stats")
	if counterCount == 0 {
		statsBtn = statsBtn.AsDisabled()
	}

	actionRow := discord.NewActionRow(addBtn, manageBtn, statsBtn)

	return []discord.LayoutComponent{
		infoDisplay,
		actionRow,
	}
}
