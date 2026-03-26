// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package view

import (
	"fmt"
	"strconv"

	"github.com/disgoorg/disgo/discord"
	"github.com/Sumire-Labs/pedmin/internal/model"
)

// PlayerVolumePresets lists the available volume presets.
var PlayerVolumePresets = []int{10, 25, 50, 75, 100}

// BuildPlayerSettingsPanel builds the player volume settings panel.
func BuildPlayerSettingsPanel(currentVolume int) []discord.LayoutComponent {
	infoDisplay := discord.NewTextDisplay(fmt.Sprintf("**デフォルト音量:** %d%%", currentVolume))

	var options []discord.StringSelectMenuOption
	for _, v := range PlayerVolumePresets {
		opt := discord.StringSelectMenuOption{
			Label: fmt.Sprintf("%d%%", v),
			Value: strconv.Itoa(v),
		}
		if v == currentVolume {
			opt.Default = true
		}
		options = append(options, opt)
	}

	volumeSelect := discord.NewActionRow(
		discord.NewStringSelectMenu(model.PlayerModuleID+":volume", "デフォルト音量を選択...", options...),
	)

	return []discord.LayoutComponent{
		infoDisplay,
		volumeSelect,
	}
}
