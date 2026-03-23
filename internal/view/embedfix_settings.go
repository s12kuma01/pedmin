// opyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package view

import (
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/s12kuma01/pedmin/internal/model"
)

// BuildEmbedFixSettingsPanel builds the embedfix platform toggle settings panel.
func BuildEmbedFixSettingsPanel(settings *model.EmbedFixSettings) []discord.LayoutComponent {
	var enabledNames []string
	for _, p := range model.AllPlatforms {
		if settings.IsPlatformEnabled(p.Key) {
			enabledNames = append(enabledNames, p.Label)
		}
	}
	statusText := "なし"
	if len(enabledNames) > 0 {
		statusText = strings.Join(enabledNames, ", ")
	}

	infoDisplay := discord.NewTextDisplay("**埋め込み対象:** " + statusText)

	var options []discord.StringSelectMenuOption
	for _, p := range model.AllPlatforms {
		opt := discord.StringSelectMenuOption{
			Label: p.Label,
			Value: string(p.Key),
		}
		if settings.IsPlatformEnabled(p.Key) {
			opt.Default = true
		}
		options = append(options, opt)
	}

	platformSelect := discord.NewActionRow(
		discord.NewStringSelectMenu(model.EmbedFixModuleID+":platforms", "埋め込み対象を選択...", options...).
			WithMinValues(0).
			WithMaxValues(len(model.AllPlatforms)),
	)

	return []discord.LayoutComponent{
		infoDisplay,
		platformSelect,
	}
}
