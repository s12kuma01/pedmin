package view

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/s12kuma01/pedmin/internal/model"
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
