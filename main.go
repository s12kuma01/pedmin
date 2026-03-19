package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/s12kuma01/pedmin/bot"
	"github.com/s12kuma01/pedmin/config"
	"github.com/s12kuma01/pedmin/features/avatar"
	"github.com/s12kuma01/pedmin/features/player"
	"github.com/s12kuma01/pedmin/features/settings"
	"github.com/s12kuma01/pedmin/store"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	cfg, err := config.Load()
	if err != nil {
		logger.Error("failed to load config", slog.Any("error", err))
		os.Exit(1)
	}

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

	// Register modules
	settingsModule := settings.New(b, logger)
	b.Register(settingsModule)

	avatarModule := avatar.New(logger)
	b.Register(avatarModule)

	playerModule := player.New(b.Lavalink, b.Client, logger)
	player.SetupListeners(b.Lavalink, playerModule)
	b.Register(playerModule)

	// Connect to Lavalink
	go func() {
		if err := player.ConnectNode(context.Background(), b.Lavalink, cfg.LavalinkHost, cfg.LavalinkPassword); err != nil {
			logger.Error("failed to connect to lavalink", slog.Any("error", err))
		} else {
			logger.Info("connected to lavalink", slog.String("host", cfg.LavalinkHost))
		}
	}()

	// Start the bot
	if err := b.Start(context.Background()); err != nil {
		logger.Error("failed to start bot", slog.Any("error", err))
		os.Exit(1)
	}
	logger.Info("pedmin is online")

	// Graceful shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	logger.Info("shutting down...")
	b.Close(context.Background())
	_ = guildStore.Close()
}
