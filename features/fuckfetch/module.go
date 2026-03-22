package fuckfetch

import (
	"log/slog"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
	"github.com/s12kuma01/pedmin/module"
)

const ModuleID = "fuckfetch"

type Fuckfetch struct {
	logger *slog.Logger
}

func New(logger *slog.Logger) *Fuckfetch {
	return &Fuckfetch{logger: logger}
}

func (f *Fuckfetch) Info() module.Info {
	return module.Info{
		ID:          ModuleID,
		Name:        "Fuckfetch",
		Description: "サーバーのシステム情報を表示",
		AlwaysOn:    true,
	}
}

func (f *Fuckfetch) Commands() []discord.ApplicationCommandCreate {
	return []discord.ApplicationCommandCreate{
		discord.SlashCommandCreate{
			Name:        "fuckfetch",
			Description: "サーバーマシンのシステム情報を表示する",
		},
	}
}

func (f *Fuckfetch) HandleComponent(_ *events.ComponentInteractionCreate) {}
func (f *Fuckfetch) HandleModal(_ *events.ModalSubmitInteractionCreate)   {}
func (f *Fuckfetch) SettingsPanel(_ snowflake.ID) []discord.LayoutComponent {
	return nil
}
