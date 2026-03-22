package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/s12kuma01/pedmin/bot"
	"github.com/s12kuma01/pedmin/config"
	"github.com/s12kuma01/pedmin/features/avatar"
	"github.com/s12kuma01/pedmin/features/embedfix"
	"github.com/s12kuma01/pedmin/features/fuckfetch"
	loggermod "github.com/s12kuma01/pedmin/features/logger"
	panelmod "github.com/s12kuma01/pedmin/features/panel"
	"github.com/s12kuma01/pedmin/features/ping"
	"github.com/s12kuma01/pedmin/features/player"
	rssmod "github.com/s12kuma01/pedmin/features/rss"
	"github.com/s12kuma01/pedmin/features/settings"
	ticketmod "github.com/s12kuma01/pedmin/features/ticket"
	urlmod "github.com/s12kuma01/pedmin/features/url"
	"github.com/s12kuma01/pedmin/store"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: cfg.LogLevel}))

	guildStore, err := store.NewSQLiteStore(cfg.DBPath)
	if err != nil {
		logger.Error("failed to create store", slog.Any("error", err))
		os.Exit(1)
	}

	b, err := bot.New(cfg, guildStore, logger)
	if err != nil {
		logger.Error("failed to create bot", slog.Any("error", err))
		os.Exit(1)
	}

	settingsModule := settings.New(b, logger)
	b.Register(settingsModule)

	avatarModule := avatar.New(logger)
	b.Register(avatarModule)

	pingModule := ping.New(logger)
	b.Register(pingModule)

	fuckfetchModule := fuckfetch.New(logger)
	b.Register(fuckfetchModule)

	ticketModule := ticketmod.New(b, b.Client, guildStore, logger)
	b.Register(ticketModule)

	loggerModule := loggermod.New(b, b.Client, guildStore, logger)
	loggermod.SetupListeners(b.Client, loggerModule)
	b.Register(loggerModule)

	embedfixModule := embedfix.New(b, b.Client, cfg.DeepLAPIKey, cfg.MetaAccessToken, cfg.HTTPClientTimeout, guildStore, logger)
	embedfix.SetupListeners(b.Client, embedfixModule)
	b.Register(embedfixModule)

	rssModule := rssmod.New(b, b.Client, guildStore, cfg.RSSPollInterval, cfg.RSSFeedTimeout, logger)
	b.Register(rssModule)

	panelModule := panelmod.New(cfg, logger)
	b.Register(panelModule)

	urlModule := urlmod.New(cfg, logger)
	b.Register(urlModule)

	playerModule := player.New(b.Lavalink, b.Client, cfg.DefaultVolume, cfg.AutoLeaveTimeout, cfg.LavalinkTimeout, cfg.LavalinkLoadTimeout, guildStore, logger)
	player.SetupListeners(b.Lavalink, playerModule)
	b.Register(playerModule)

	go func() {
		if err := player.ConnectNode(context.Background(), b.Lavalink, cfg.LavalinkHost, cfg.LavalinkPassword); err != nil {
			logger.Error("failed to connect to lavalink", slog.Any("error", err))
		} else {
			logger.Info("connected to lavalink", slog.String("host", cfg.LavalinkHost))
		}
	}()

	if err := b.Start(context.Background()); err != nil {
		logger.Error("failed to start bot", slog.Any("error", err))
		os.Exit(1)
	}
	logger.Info("pedmin is online")

	rssModule.StartPoller(context.Background())

	// Graceful shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	logger.Info("shutting down...")
	rssModule.StopPoller()
	b.Close(context.Background())
	if err := guildStore.Close(); err != nil {
		logger.Error("failed to close store", slog.Any("error", err))
	}
}
