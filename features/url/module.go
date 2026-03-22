package url

import (
	"log/slog"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
	"github.com/s12kuma01/pedmin/config"
	"github.com/s12kuma01/pedmin/module"
)

const ModuleID = "url"

type URL struct {
	cfg    *config.Config
	client *URLClient
	logger *slog.Logger
}

func New(cfg *config.Config, logger *slog.Logger) *URL {
	return &URL{
		cfg:    cfg,
		client: NewURLClient(cfg.XGDAPIKey, cfg.VTAPIKey, cfg.HTTPClientTimeout),
		logger: logger,
	}
}

func (u *URL) Info() module.Info {
	return module.Info{
		ID:          ModuleID,
		Name:        "URL Tools",
		Description: "URL短縮・安全チェック",
		AlwaysOn:    true,
	}
}

func (u *URL) Commands() []discord.ApplicationCommandCreate {
	return []discord.ApplicationCommandCreate{
		discord.SlashCommandCreate{
			Name:        "url",
			Description: "URLツール（短縮・安全チェック）",
		},
	}
}

func (u *URL) SettingsPanel(_ snowflake.ID) []discord.LayoutComponent { return nil }
