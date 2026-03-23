package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/s12kuma01/pedmin/config"
	"github.com/s12kuma01/pedmin/internal/bot"
	"github.com/s12kuma01/pedmin/internal/client"
	"github.com/s12kuma01/pedmin/internal/handler"
	"github.com/s12kuma01/pedmin/internal/repository"
	"github.com/s12kuma01/pedmin/internal/service"
	"github.com/s12kuma01/pedmin/pkg/deepl"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: cfg.LogLevel}))

	guildStore, err := repository.NewSQLiteStore(cfg.DBPath)
	if err != nil {
		logger.Error("failed to create store", slog.Any("error", err))
		os.Exit(1)
	}

	b, err := bot.New(cfg, guildStore, logger)
	if err != nil {
		logger.Error("failed to create bot", slog.Any("error", err))
		os.Exit(1)
	}

	// --- Simple modules ---

	settingsHandler := handler.NewSettingsHandler(b, logger)
	b.Register(settingsHandler)

	avatarHandler := handler.NewAvatarHandler(logger)
	b.Register(avatarHandler)

	pingHandler := handler.NewPingHandler(logger)
	b.Register(pingHandler)

	fuckfetchHandler := handler.NewFuckfetchHandler(logger)
	b.Register(fuckfetchHandler)

	// --- Ticket ---

	ticketSvc := service.NewTicketService(b.Client, guildStore, logger)
	ticketHandler := handler.NewTicketHandler(b, ticketSvc, logger)
	b.Register(ticketHandler)

	// --- Logger ---

	loggerSvc := service.NewLoggerService(guildStore)
	loggerHandler := handler.NewLoggerHandler(b, b.Client, loggerSvc, logger)
	handler.SetupLoggerListeners(b.Client, loggerHandler)
	b.Register(loggerHandler)

	// --- Embed Fix ---

	twitterClient := client.NewFxTwitterClient(cfg.HTTPClientTimeout)
	redditClient := client.NewRedditClient(cfg.HTTPClientTimeout)
	tiktokClient := client.NewTikTokClient(cfg.HTTPClientTimeout)
	deeplClient := deepl.NewTranslateClient(cfg.DeepLAPIKey, cfg.HTTPClientTimeout)

	embedfixSvc := service.NewEmbedFixService(
		guildStore, twitterClient, redditClient, tiktokClient,
		deeplClient, b.Client, logger,
	)
	embedfixHandler := handler.NewEmbedFixHandler(b, embedfixSvc, logger)
	handler.SetupEmbedFixListeners(b.Client, embedfixHandler)
	b.Register(embedfixHandler)

	// --- Translator ---

	translatorHandler := handler.NewTranslatorHandler(b, b.Client, cfg.DeepLAPIKey, cfg.HTTPClientTimeout, logger)
	handler.SetupTranslatorListeners(b.Client, translatorHandler)
	b.Register(translatorHandler)

	// --- RSS ---

	rssSvc := service.NewRSSService(guildStore, b.Client, logger)
	rssHandler := handler.NewRSSHandler(rssSvc, cfg.RSSFeedTimeout, logger)
	b.Register(rssHandler)
	rssPoller := service.NewRSSPoller(b, rssSvc, guildStore, cfg.RSSPollInterval, logger)

	// --- Panel ---

	panelClient := client.NewPelicanClient(cfg.PanelURL, cfg.PanelAPIKey, cfg.PanelPowerActionTimeout)
	panelSvc := service.NewPanelService(panelClient)
	panelHandler := handler.NewPanelHandler(cfg, panelSvc, logger)
	b.Register(panelHandler)

	// --- URL ---

	urlHandler := handler.NewURLHandler(cfg, logger)
	b.Register(urlHandler)

	// --- Player ---

	playerSvc := service.NewPlayerService(
		b.Lavalink, b.Client,
		cfg.DefaultVolume, cfg.AutoLeaveTimeout,
		cfg.LavalinkTimeout, cfg.LavalinkLoadTimeout,
		guildStore, logger,
	)
	service.SetupPlayerListeners(b.Lavalink, playerSvc)
	playerHandler := handler.NewPlayerHandler(playerSvc, logger)
	b.Register(playerHandler)

	go func() {
		if err := service.ConnectNode(context.Background(), b.Lavalink, cfg.LavalinkHost, cfg.LavalinkPassword); err != nil {
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

	rssPoller.StartPoller(context.Background())

	// Graceful shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	logger.Info("shutting down...")
	playerHandler.Shutdown()
	rssPoller.StopPoller()
	b.Close(context.Background())
	if err := guildStore.Close(); err != nil {
		logger.Error("failed to close store", slog.Any("error", err))
	}
}
