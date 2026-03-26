// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package view

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/Sumire-Labs/pedmin/internal/model"
)

// LevelingMultipliersTab builds the multipliers management tab.
func LevelingMultipliersTab(settings *model.LevelingSettings) discord.MessageUpdate {
	components := []discord.ContainerSubComponent{
		discord.NewTextDisplay("### XP倍率"),
		discord.NewSmallSeparator(),
	}

	roleText := "**ロール倍率:**\n"
	if len(settings.RoleMultipliers) == 0 {
		roleText += "なし\n"
	} else {
		for i, rm := range settings.RoleMultipliers {
			roleText += fmt.Sprintf("<@&%d> → %.1fx", rm.RoleID, rm.Multiplier)
			roleText += fmt.Sprintf("　[`削除`](%s:role_mult_remove:%d)\n", model.LevelingModuleID, i)
		}
	}

	chText := "**チャンネル倍率:**\n"
	if len(settings.ChannelMultipliers) == 0 {
		chText += "なし\n"
	} else {
		for i, cm := range settings.ChannelMultipliers {
			chText += fmt.Sprintf("<#%d> → %.1fx", cm.ChannelID, cm.Multiplier)
			chText += fmt.Sprintf("　[`削除`](%s:ch_mult_remove:%d)\n", model.LevelingModuleID, i)
		}
	}

	components = append(components, discord.NewTextDisplay(roleText+"\n"+chText))

	totalMults := len(settings.RoleMultipliers) + len(settings.ChannelMultipliers)
	addRoleBtn := discord.NewSecondaryButton("ロール倍率追加", model.LevelingModuleID+":role_mult_add")
	addChBtn := discord.NewSecondaryButton("チャンネル倍率追加", model.LevelingModuleID+":ch_mult_add")
	if totalMults >= model.MaxMultipliersPerGuild {
		addRoleBtn = addRoleBtn.AsDisabled()
		addChBtn = addChBtn.AsDisabled()
	}

	var removeButtons []discord.InteractiveComponent
	for i := range settings.RoleMultipliers {
		removeButtons = append(removeButtons, discord.NewDangerButton(
			fmt.Sprintf("ロール%d削除", i+1),
			fmt.Sprintf("%s:role_mult_remove:%d", model.LevelingModuleID, i),
		))
	}
	for i := range settings.ChannelMultipliers {
		removeButtons = append(removeButtons, discord.NewDangerButton(
			fmt.Sprintf("CH%d削除", i+1),
			fmt.Sprintf("%s:ch_mult_remove:%d", model.LevelingModuleID, i),
		))
	}

	components = append(components,
		discord.NewLargeSeparator(),
		discord.NewActionRow(addRoleBtn, addChBtn),
	)

	if len(removeButtons) > 0 {
		for i := 0; i < len(removeButtons); i += 5 {
			end := i + 5
			if end > len(removeButtons) {
				end = len(removeButtons)
			}
			components = append(components, discord.NewActionRow(removeButtons[i:end]...))
		}
	}

	components = append(components,
		discord.NewActionRow(
			discord.NewSecondaryButton("一般", model.LevelingModuleID+":settings_tab:general"),
			discord.NewSecondaryButton("報酬", model.LevelingModuleID+":settings_tab:rewards"),
			discord.NewPrimaryButton("倍率", model.LevelingModuleID+":settings_tab:multipliers").AsDisabled(),
		),
	)

	return discord.NewMessageUpdateV2([]discord.LayoutComponent{
		discord.NewContainer(components...),
	})
}
