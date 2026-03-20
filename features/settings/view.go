package settings

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
)

func (s *Settings) mainPanel(guildID snowflake.ID) discord.MessageCreate {
	return ephemeralV2(s.buildMainContainer(guildID))
}

func (s *Settings) mainPanelUpdate(guildID snowflake.ID) discord.MessageUpdate {
	return discord.NewMessageUpdateV2([]discord.LayoutComponent{s.buildMainContainer(guildID)})
}

func (s *Settings) buildMainContainer(guildID snowflake.ID) discord.ContainerComponent {
	modules := s.bot.GetModules()
	var options []discord.StringSelectMenuOption
	for _, m := range modules {
		info := m.Info()
		if info.AlwaysOn {
			continue
		}
		status := "❌"
		if s.bot.IsModuleEnabled(guildID, info.ID) {
			status = "✅"
		}
		options = append(options, discord.StringSelectMenuOption{
			Label:       fmt.Sprintf("%s %s", status, info.Name),
			Value:       info.ID,
			Description: info.Description,
		})
	}

	components := []discord.ContainerSubComponent{
		discord.NewTextDisplay("## Pedmin Settings"),
		discord.NewLargeSeparator(),
	}

	if len(options) > 0 {
		components = append(components,
			discord.NewTextDisplay("設定するモジュールを選択してください:"),
			discord.NewActionRow(
				discord.NewStringSelectMenu(ModuleID+":select", "モジュールを選択...", options...),
			),
		)
	} else {
		components = append(components,
			discord.NewTextDisplay("設定可能なモジュールがありません。"),
		)
	}

	return discord.NewContainer(components...)
}

func (s *Settings) modulePanel(guildID snowflake.ID, moduleID string) discord.MessageUpdate {
	modules := s.bot.GetModules()
	m, ok := modules[moduleID]
	if !ok {
		return discord.NewMessageUpdateV2([]discord.LayoutComponent{
			discord.NewContainer(
				discord.NewTextDisplay("モジュールが見つかりません。"),
			),
		})
	}

	info := m.Info()
	enabled := s.bot.IsModuleEnabled(guildID, moduleID)

	statusText := "無効"
	toggleLabel := "有効にする"
	toggleStyle := discord.ButtonStyleSuccess
	if enabled {
		statusText = "有効"
		toggleLabel = "無効にする"
		toggleStyle = discord.ButtonStyleDanger
	}

	components := []discord.ContainerSubComponent{
		discord.NewTextDisplay(fmt.Sprintf("## %s", info.Name)),
		discord.NewTextDisplay(info.Description),
		discord.NewLargeSeparator(),
		discord.NewTextDisplay(fmt.Sprintf("**ステータス:** %s", statusText)),
	}

	settingsPanel := m.SettingsPanel(guildID)
	if len(settingsPanel) > 0 {
		components = append(components, discord.NewLargeSeparator())
		for _, lc := range settingsPanel {
			if sub, ok := lc.(discord.ContainerSubComponent); ok {
				components = append(components, sub)
			}
		}
	}

	components = append(components,
		discord.NewLargeSeparator(),
		discord.NewActionRow(
			discord.NewButton(toggleStyle, toggleLabel, fmt.Sprintf("%s:toggle:%s", ModuleID, moduleID), "", 0),
			discord.NewSecondaryButton("← 戻る", ModuleID+":back"),
		),
	)

	return discord.NewMessageUpdateV2([]discord.LayoutComponent{
		discord.NewContainer(components...),
	})
}

func ephemeralV2(components ...discord.LayoutComponent) discord.MessageCreate {
	return discord.NewMessageCreateV2(components...).WithEphemeral(true)
}
