package bot

import (
	"context"
	"log/slog"

	"github.com/disgoorg/disgo"
	disgobot "github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/snowflake/v2"
	"github.com/s12kuma01/pedmin/config"
	"github.com/s12kuma01/pedmin/module"
	"github.com/s12kuma01/pedmin/store"
)

type Bot struct {
	Cfg            *config.Config
	Client         *disgobot.Client
	Lavalink       disgolink.Client
	Store          store.GuildStore
	Modules        map[string]module.Module
	Logger         *slog.Logger
	cancelPresence context.CancelFunc
}

func New(cfg *config.Config, guildStore store.GuildStore, logger *slog.Logger) (*Bot, error) {
	b := &Bot{
		Cfg:     cfg,
		Store:   guildStore,
		Modules: make(map[string]module.Module),
		Logger:  logger,
	}

	client, err := disgo.New(cfg.Token,
		disgobot.WithCacheConfigOpts(
			cache.WithCaches(cache.FlagVoiceStates),
		),
		disgobot.WithGatewayConfigOpts(
			gateway.WithIntents(
				gateway.IntentGuilds,
				gateway.IntentGuildVoiceStates,
				gateway.IntentGuildMessages,
			),
		),
		disgobot.WithEventListenerFunc(b.onCommandInteraction),
		disgobot.WithEventListenerFunc(b.onComponentInteraction),
		disgobot.WithEventListenerFunc(b.onModalSubmit),
		disgobot.WithEventListenerFunc(b.onVoiceStateUpdate),
		disgobot.WithEventListenerFunc(b.onVoiceServerUpdate),
	)
	if err != nil {
		return nil, err
	}
	b.Client = client

	b.Lavalink = disgolink.New(cfg.AppID)

	return b, nil
}

func (b *Bot) Register(m module.Module) {
	info := m.Info()
	b.Modules[info.ID] = m
	b.Logger.Info("registered module", slog.String("module", info.ID))
}

func (b *Bot) Start(ctx context.Context) error {
	if err := b.SyncCommands(); err != nil {
		return err
	}
	if err := b.Client.OpenGateway(ctx); err != nil {
		return err
	}

	presenceCtx, cancel := context.WithCancel(context.Background())
	b.cancelPresence = cancel
	go b.startPresenceUpdater(presenceCtx)

	return nil
}

func (b *Bot) Close(ctx context.Context) {
	if b.cancelPresence != nil {
		b.cancelPresence()
	}
	b.Client.Close(ctx)
	b.Lavalink.Close()
}

func (b *Bot) IsModuleEnabled(guildID snowflake.ID, moduleID string) bool {
	m, ok := b.Modules[moduleID]
	if !ok {
		return false
	}
	if m.Info().AlwaysOn {
		return true
	}
	enabled, err := b.Store.IsModuleEnabled(guildID, moduleID)
	if err != nil {
		b.Logger.Error("failed to check module enabled", slog.String("module", moduleID), slog.Any("error", err))
		return false
	}
	return enabled
}

func (b *Bot) GetModules() map[string]module.Module {
	return b.Modules
}

func (b *Bot) SetModuleEnabled(guildID snowflake.ID, moduleID string, enabled bool) error {
	return b.Store.SetModuleEnabled(guildID, moduleID, enabled)
}
