package avatar

import (
	"log/slog"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
	"github.com/s12kuma01/pedmin/module"
)

const ModuleID = "avatar"

type Avatar struct {
	logger *slog.Logger
}

func New(logger *slog.Logger) *Avatar {
	return &Avatar{logger: logger}
}

func (a *Avatar) Info() module.Info {
	return module.Info{
		ID:          ModuleID,
		Name:        "アバター",
		Description: "ユーザーのアバターを表示",
		AlwaysOn:    true,
	}
}

func (a *Avatar) Commands() []discord.ApplicationCommandCreate {
	return []discord.ApplicationCommandCreate{
		discord.SlashCommandCreate{
			Name:        "avatar",
			Description: "ユーザーのアバターを表示する",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionUser{
					Name:        "user",
					Description: "アバターを表示するユーザー（省略時は自分）",
					Required:    false,
				},
			},
		},
	}
}

func (a *Avatar) HandleComponent(_ *events.ComponentInteractionCreate) {}
func (a *Avatar) HandleModal(_ *events.ModalSubmitInteractionCreate)   {}
func (a *Avatar) SettingsPanel(_ snowflake.ID) []discord.LayoutComponent {
	return nil
}
