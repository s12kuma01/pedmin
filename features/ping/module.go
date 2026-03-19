package ping

import (
	"log/slog"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
	"github.com/s12kuma01/pedmin/module"
)

const ModuleID = "ping"

type Ping struct {
	logger *slog.Logger
}

func New(logger *slog.Logger) *Ping {
	return &Ping{logger: logger}
}

func (p *Ping) Info() module.Info {
	return module.Info{
		ID:          ModuleID,
		Name:        "Ping",
		Description: "Botの応答確認",
		AlwaysOn:    true,
	}
}

func (p *Ping) Commands() []discord.ApplicationCommandCreate {
	return []discord.ApplicationCommandCreate{
		discord.SlashCommandCreate{
			Name:        "ping",
			Description: "Botの応答速度を確認する",
		},
	}
}

func (p *Ping) HandleComponent(_ *events.ComponentInteractionCreate) {}
func (p *Ping) HandleModal(_ *events.ModalSubmitInteractionCreate)   {}
func (p *Ping) SettingsPanel(_ snowflake.ID) []discord.LayoutComponent {
	return nil
}
func (p *Ping) HandleSettingsComponent(_ *events.ComponentInteractionCreate) {}
