package embedfix

import (
	"log/slog"
	"time"

	disgobot "github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
	"github.com/s12kuma01/pedmin/module"
)

const ModuleID = "embedfix"

// TODO: Add embed replacement support for other SNS platforms (Reddit, TikTok, Instagram)

type Bot interface {
	IsModuleEnabled(guildID snowflake.ID, moduleID string) bool
}

type EmbedFix struct {
	bot             Bot
	client          *disgobot.Client
	fxClient        *FxTwitterClient
	translateClient *TranslateClient
	logger          *slog.Logger
}

func New(bot Bot, client *disgobot.Client, deeplAPIKey string, timeout time.Duration, logger *slog.Logger) *EmbedFix {
	return &EmbedFix{
		bot:             bot,
		client:          client,
		fxClient:        NewFxTwitterClient(timeout),
		translateClient: NewTranslateClient(deeplAPIKey, timeout),
		logger:          logger,
	}
}

func (ef *EmbedFix) Info() module.Info {
	return module.Info{
		ID:          ModuleID,
		Name:        "Embed Fix",
		Description: "X/Twitterリンクのリッチ埋め込み表示",
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

func (ef *EmbedFix) SettingsPanel(_ snowflake.ID) []discord.LayoutComponent {
	return nil
}

func (ef *EmbedFix) HandleSettingsComponent(_ *events.ComponentInteractionCreate) {}
