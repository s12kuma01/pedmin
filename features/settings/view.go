package settings

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/s12kuma01/pedmin/module"
	"github.com/s12kuma01/pedmin/ui"
)

// ModuleOption holds display data for a module in the settings panel.
type ModuleOption struct {
	ID          string
	Name        string
	Description string
	Enabled     bool
	Summary     string // optional: brief settings summary (e.g. "ログ先: #logs")
}

// BuildMainPanel builds the initial settings panel message.
func BuildMainPanel(options []ModuleOption) discord.MessageCreate {
	return ui.EphemeralV2(buildMainContainer(options))
}

// BuildMainPanelUpdate builds the settings panel as a message update.
func BuildMainPanelUpdate(options []ModuleOption) discord.MessageUpdate {
	return discord.NewMessageUpdateV2([]discord.LayoutComponent{buildMainContainer(options)})
}

func buildMainContainer(options []ModuleOption) discord.ContainerComponent {
	var selectOptions []discord.StringSelectMenuOption
	for _, opt := range options {
		status := "❌"
		if opt.Enabled {
			status = "✅"
		}
		desc := opt.Description
		if opt.Summary != "" {
			desc = opt.Summary
		}
		selectOptions = append(selectOptions, discord.StringSelectMenuOption{
			Label:       fmt.Sprintf("%s %s", status, opt.Name),
			Value:       opt.ID,
			Description: desc,
		})
	}

	components := []discord.ContainerSubComponent{
		discord.NewTextDisplay("## Pedmin Settings"),
		discord.NewLargeSeparator(),
	}

	if len(selectOptions) > 0 {
		components = append(components,
			discord.NewTextDisplay("設定するモジュールを選択してください:"),
			discord.NewActionRow(
				discord.NewStringSelectMenu(ModuleID+":select", "モジュールを選択...", selectOptions...),
			),
		)
	} else {
		components = append(components,
			discord.NewTextDisplay("設定可能なモジュールがありません。"),
		)
	}

	return discord.NewContainer(components...)
}

// BuildModulePanel builds the module detail panel as a message update.
func BuildModulePanel(info module.Info, enabled bool, settingsPanel []discord.LayoutComponent) discord.MessageUpdate {
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
			discord.NewButton(toggleStyle, toggleLabel, fmt.Sprintf("%s:toggle:%s", ModuleID, info.ID), "", 0),
			discord.NewSecondaryButton("← 戻る", ModuleID+":back"),
		),
	)

	return discord.NewMessageUpdateV2([]discord.LayoutComponent{
		discord.NewContainer(components...),
	})
}

// BuildModuleNotFound builds an error panel when a module is not found.
func BuildModuleNotFound() discord.MessageUpdate {
	return discord.NewMessageUpdateV2([]discord.LayoutComponent{
		discord.NewContainer(
			discord.NewTextDisplay("モジュールが見つかりません。"),
		),
	})
}
