package player

import (
	"fmt"
	"strconv"

	"github.com/disgoorg/disgo/discord"
)

var volumePresets = []int{10, 25, 50, 75, 100}

func BuildSettingsPanel(currentVolume int) []discord.LayoutComponent {
	infoDisplay := discord.NewTextDisplay(fmt.Sprintf("**デフォルト音量:** %d%%", currentVolume))

	var options []discord.StringSelectMenuOption
	for _, v := range volumePresets {
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
		discord.NewStringSelectMenu(ModuleID+":volume", "デフォルト音量を選択...", options...),
	)

	return []discord.LayoutComponent{
		infoDisplay,
		volumeSelect,
	}
}
