package ticket

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

const ModuleID = "ticket"

type Bot interface {
	IsModuleEnabled(guildID snowflake.ID, moduleID string) bool
}

type Ticket struct {
	bot    Bot
	client *disgobot.Client
	store  store.GuildStore
	logger *slog.Logger
}

func New(bot Bot, client *disgobot.Client, guildStore store.GuildStore, logger *slog.Logger) *Ticket {
	return &Ticket{
		bot:    bot,
		client: client,
		store:  guildStore,
		logger: logger,
	}
}

func (t *Ticket) Info() module.Info {
	return module.Info{
		ID:          ModuleID,
		Name:        "チケット",
		Description: "サポートチケットシステム",
		AlwaysOn:    false,
	}
}

func (t *Ticket) Commands() []discord.ApplicationCommandCreate {
	return nil
}

func (t *Ticket) HandleCommand(_ *events.ApplicationCommandInteractionCreate) {}

func (t *Ticket) HandleComponent(e *events.ComponentInteractionCreate) {
	t.handleComponent(e)
}

func (t *Ticket) HandleModal(e *events.ModalSubmitInteractionCreate) {
	t.handleModal(e)
}

func (t *Ticket) SettingsSummary(guildID snowflake.ID) string {
	settings, err := LoadSettings(t.store, guildID)
	if err != nil {
		return ""
	}
	var parts []string
	if settings.CategoryID != 0 {
		parts = append(parts, fmt.Sprintf("カテゴリ: #%d", settings.CategoryID))
	}
	if settings.LogChannelID != 0 {
		parts = append(parts, fmt.Sprintf("ログ: #%d", settings.LogChannelID))
	}
	if len(parts) == 0 {
		return "未設定"
	}
	return strings.Join(parts, ", ")
}

func (t *Ticket) SettingsPanel(guildID snowflake.ID) []discord.LayoutComponent {
	settings, err := LoadSettings(t.store, guildID)
	if err != nil {
		t.logger.Error("failed to load ticket settings", slog.Any("error", err))
		settings = &TicketSettings{}
	}
	return BuildSettingsPanel(settings)
}
