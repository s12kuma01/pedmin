package settings

import (
	"log/slog"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/omit"
	"github.com/disgoorg/snowflake/v2"
	"github.com/s12kuma01/pedmin/module"
)

const ModuleID = "settings"

type Bot interface {
	IsModuleEnabled(guildID snowflake.ID, moduleID string) bool
	GetModules() map[string]module.Module
	SetModuleEnabled(guildID snowflake.ID, moduleID string, enabled bool) error
}

type Settings struct {
	bot    Bot
	logger *slog.Logger
}

func New(bot Bot, logger *slog.Logger) *Settings {
	return &Settings{bot: bot, logger: logger}
}

func (s *Settings) Info() module.Info {
	return module.Info{
		ID:          ModuleID,
		Name:        "設定",
		Description: "サーバー設定管理パネル",
		AlwaysOn:    true,
	}
}

func (s *Settings) Commands() []discord.ApplicationCommandCreate {
	manageGuild := discord.PermissionManageGuild
	return []discord.ApplicationCommandCreate{
		discord.SlashCommandCreate{
			Name:                     "settings",
			Description:              "サーバー設定パネルを開く",
			DefaultMemberPermissions: omit.New(&manageGuild),
		},
	}
}

func (s *Settings) HandleModal(_ *events.ModalSubmitInteractionCreate) {}

func (s *Settings) SettingsPanel(_ snowflake.ID) []discord.LayoutComponent {
	return nil
}
