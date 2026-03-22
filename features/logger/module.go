package logger

import (
	"fmt"
	"log/slog"
	"strings"

	disgobot "github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
	"github.com/s12kuma01/pedmin/module"
	"github.com/s12kuma01/pedmin/store"
)

const ModuleID = "logger"

type Bot interface {
	IsModuleEnabled(guildID snowflake.ID, moduleID string) bool
}

type Logger struct {
	bot    Bot
	client *disgobot.Client
	store  store.GuildStore
	logger *slog.Logger
}

func New(bot Bot, client *disgobot.Client, guildStore store.GuildStore, logger *slog.Logger) *Logger {
	return &Logger{
		bot:    bot,
		client: client,
		store:  guildStore,
		logger: logger,
	}
}

func (l *Logger) Info() module.Info {
	return module.Info{
		ID:          ModuleID,
		Name:        "Logger",
		Description: "サーバーイベントのログを記録",
		AlwaysOn:    false,
	}
}

func (l *Logger) Commands() []discord.ApplicationCommandCreate {
	return nil
}

func (l *Logger) HandleCommand(_ *events.ApplicationCommandInteractionCreate) {}

func (l *Logger) HandleComponent(e *events.ComponentInteractionCreate) {
	l.handleComponent(e)
}

func (l *Logger) HandleModal(_ *events.ModalSubmitInteractionCreate) {}

func (l *Logger) SettingsSummary(guildID snowflake.ID) string {
	settings, err := LoadSettings(l.store, guildID)
	if err != nil {
		return ""
	}
	var parts []string
	if settings.ChannelID != 0 {
		parts = append(parts, fmt.Sprintf("ログ先: #%d", settings.ChannelID))
	}
	count := 0
	for _, enabled := range settings.Events {
		if enabled {
			count++
		}
	}
	if count > 0 {
		parts = append(parts, fmt.Sprintf("イベント: %d個", count))
	}
	if len(parts) == 0 {
		return "未設定"
	}
	return strings.Join(parts, ", ")
}

func (l *Logger) SettingsPanel(guildID snowflake.ID) []discord.LayoutComponent {
	settings, err := LoadSettings(l.store, guildID)
	if err != nil {
		l.logger.Error("failed to load logger settings", slog.Any("error", err))
		settings = &LoggerSettings{Events: make(map[string]bool)}
	}
	return BuildSettingsPanel(settings)
}
