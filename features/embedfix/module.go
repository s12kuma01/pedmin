package embedfix

import (
	"log/slog"
	"strings"
	"time"

	disgobot "github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
	"github.com/s12kuma01/pedmin/deepl"
	"github.com/s12kuma01/pedmin/module"
	"github.com/s12kuma01/pedmin/store"
)

const ModuleID = "embedfix"

type Bot interface {
	IsModuleEnabled(guildID snowflake.ID, moduleID string) bool
}

type EmbedFix struct {
	bot             Bot
	client          *disgobot.Client
	store           store.GuildStore
	twitterClient   *FxTwitterClient
	redditClient    *RedditClient
	tiktokClient    *TikTokClient
	translateClient *deepl.TranslateClient
	logger          *slog.Logger
}

func New(bot Bot, client *disgobot.Client, deeplAPIKey string, timeout time.Duration, guildStore store.GuildStore, logger *slog.Logger) *EmbedFix {
	return &EmbedFix{
		bot:             bot,
		client:          client,
		store:           guildStore,
		twitterClient:   NewFxTwitterClient(timeout),
		redditClient:    NewRedditClient(timeout),
		tiktokClient:    NewTikTokClient(timeout),
		translateClient: deepl.NewTranslateClient(deeplAPIKey, timeout),
		logger:          logger,
	}
}

func (ef *EmbedFix) Info() module.Info {
	return module.Info{
		ID:          ModuleID,
		Name:        "Embed Fix",
		Description: "SNSリンクのリッチ埋め込み表示",
		AlwaysOn:    false,
	}
}

func (ef *EmbedFix) Commands() []discord.ApplicationCommandCreate {
	return nil
}

func (ef *EmbedFix) HandleCommand(_ *events.ApplicationCommandInteractionCreate) {}

func (ef *EmbedFix) HandleComponent(e *events.ComponentInteractionCreate) {
	ef.handleComponent(e)
}

func (ef *EmbedFix) HandleModal(_ *events.ModalSubmitInteractionCreate) {}

func (ef *EmbedFix) SettingsSummary(guildID snowflake.ID) string {
	settings, err := LoadSettings(ef.store, guildID)
	if err != nil {
		return ""
	}
	var names []string
	for _, p := range AllPlatforms {
		if settings.IsPlatformEnabled(p.Key) {
			names = append(names, p.Label)
		}
	}
	if len(names) == 0 {
		return "全て無効"
	}
	return "対象: " + strings.Join(names, ", ")
}

func (ef *EmbedFix) SettingsPanel(guildID snowflake.ID) []discord.LayoutComponent {
	settings, err := LoadSettings(ef.store, guildID)
	if err != nil {
		ef.logger.Error("failed to load embedfix settings", slog.Any("error", err))
		settings = defaultSettings()
	}
	return BuildSettingsPanel(settings)
}
