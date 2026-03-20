package url

import (
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

func (u *URL) HandleComponent(e *events.ComponentInteractionCreate) {
	customID := e.Data.CustomID()
	_, action, _ := strings.Cut(customID, ":")

	switch action {
	case "shorten":
		_ = e.Modal(discord.ModalCreate{
			CustomID: ModuleID + ":shorten_modal",
			Title:    "URL短縮",
			Components: []discord.LayoutComponent{
				discord.NewLabel("URL",
					discord.NewShortTextInput(ModuleID+":url").
						WithRequired(true).
						WithPlaceholder("https://example.com/long/path"),
				),
			},
		})

	case "check":
		_ = e.Modal(discord.ModalCreate{
			CustomID: ModuleID + ":check_modal",
			Title:    "URLチェッカー",
			Components: []discord.LayoutComponent{
				discord.NewLabel("URL",
					discord.NewShortTextInput(ModuleID+":url").
						WithRequired(true).
						WithPlaceholder("https://example.com"),
				),
			},
		})

	case "back":
		hasXGD := u.cfg.XGDAPIKey != ""
		hasVT := u.cfg.VTAPIKey != ""

		msg := BuildMainPanel(hasXGD, hasVT)
		_, _ = e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(),
			discord.NewMessageUpdateV2(msg.Components))
	}
}
