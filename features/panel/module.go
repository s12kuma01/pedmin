package panel

import (
	"log/slog"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
	"github.com/s12kuma01/pedmin/config"
	"github.com/s12kuma01/pedmin/module"
)

const ModuleID = "panel"

type Panel struct {
	cfg     *config.Config
	pelican *PelicanClient
	logger  *slog.Logger
}

func New(cfg *config.Config, logger *slog.Logger) *Panel {
	return &Panel{
		cfg:     cfg,
		pelican: NewPelicanClient(cfg.PanelURL, cfg.PanelAPIKey),
		logger:  logger,
	}
}

func (p *Panel) Info() module.Info {
	return module.Info{
		ID:          ModuleID,
		Name:        "Panel",
		Description: "ゲームサーバー管理",
		AlwaysOn:    true,
	}
}

func (p *Panel) Commands() []discord.ApplicationCommandCreate {
	return []discord.ApplicationCommandCreate{
		discord.SlashCommandCreate{
			Name:        "panel",
			Description: "ゲームサーバーを管理する",
		},
	}
}

func (p *Panel) HandleSettingsComponent(_ *events.ComponentInteractionCreate) {}
func (p *Panel) SettingsPanel(_ snowflake.ID) []discord.LayoutComponent       { return nil }

func (p *Panel) isAllowed(userID snowflake.ID) bool {
	for _, id := range p.cfg.PanelAllowedUsers {
		if id == userID {
			return true
		}
	}
	return false
}
