package translator

import (
	"log/slog"
	"time"

	disgobot "github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
	"github.com/s12kuma01/pedmin/deepl"
	"github.com/s12kuma01/pedmin/module"
)

const ModuleID = "translator"

type Bot interface {
	IsModuleEnabled(guildID snowflake.ID, moduleID string) bool
}

type Translator struct {
	bot             Bot
	client          *disgobot.Client
	translateClient *deepl.TranslateClient
	logger          *slog.Logger
}

func New(bot Bot, client *disgobot.Client, deeplAPIKey string, timeout time.Duration, logger *slog.Logger) *Translator {
	return &Translator{
		bot:             bot,
		client:          client,
		translateClient: deepl.NewTranslateClient(deeplAPIKey, timeout),
		logger:          logger,
	}
}

func (t *Translator) Info() module.Info {
	return module.Info{
		ID:          ModuleID,
		Name:        "Translator",
		Description: "国旗リアクションでメッセージを翻訳",
		AlwaysOn:    false,
	}
}

func (t *Translator) Commands() []discord.ApplicationCommandCreate {
	return nil
}

func (t *Translator) HandleCommand(_ *events.ApplicationCommandInteractionCreate) {}

func (t *Translator) HandleComponent(_ *events.ComponentInteractionCreate) {}

func (t *Translator) HandleModal(_ *events.ModalSubmitInteractionCreate) {}

func (t *Translator) SettingsPanel(_ snowflake.ID) []discord.LayoutComponent {
	return nil
}
