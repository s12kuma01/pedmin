package rss

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	disgobot "github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
	"github.com/s12kuma01/pedmin/module"
	"github.com/s12kuma01/pedmin/store"
)

const ModuleID = "rss"

type Bot interface {
	IsModuleEnabled(guildID snowflake.ID, moduleID string) bool
}

type RSS struct {
	bot          Bot
	client       *disgobot.Client
	store        store.GuildStore
	logger       *slog.Logger
	cancel       context.CancelFunc
	pollInterval time.Duration
	feedTimeout  time.Duration
}

func New(bot Bot, client *disgobot.Client, guildStore store.GuildStore, pollInterval, feedTimeout time.Duration, logger *slog.Logger) *RSS {
	return &RSS{
		bot:          bot,
		client:       client,
		store:        guildStore,
		pollInterval: pollInterval,
		feedTimeout:  feedTimeout,
		logger:       logger,
	}
}

func (r *RSS) Info() module.Info {
	return module.Info{
		ID:          ModuleID,
		Name:        "RSS",
		Description: "RSSフィード監視",
		AlwaysOn:    false,
	}
}

func (r *RSS) Commands() []discord.ApplicationCommandCreate {
	return nil
}

func (r *RSS) HandleCommand(_ *events.ApplicationCommandInteractionCreate) {}

func (r *RSS) HandleComponent(e *events.ComponentInteractionCreate) {
	r.handleComponent(e)
}

func (r *RSS) HandleModal(e *events.ModalSubmitInteractionCreate) {
	r.handleModal(e)
}

func (r *RSS) SettingsSummary(guildID snowflake.ID) string {
	count, err := r.store.CountRSSFeeds(guildID)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("フィード: %d/%d件", count, MaxFeedsPerGuild)
}

func (r *RSS) SettingsPanel(guildID snowflake.ID) []discord.LayoutComponent {
	count, err := r.store.CountRSSFeeds(guildID)
	if err != nil {
		r.logger.Error("failed to count rss feeds", slog.Any("error", err))
	}
	return BuildSettingsPanel(count)
}
