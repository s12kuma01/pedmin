// Package config loads application configuration from environment variables and TOML files.
package config

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/disgoorg/snowflake/v2"
)

type Config struct {
	// env (secrets)
	Token string
	AppID snowflake.ID

	// TOML (app settings)
	LavalinkHost     string
	LavalinkPassword string
	DataDir          string
	DBPath           string
	DefaultVolume    int
	AutoLeaveTimeout time.Duration
	PresenceInterval time.Duration
	LogLevel         slog.Level

	// Timeouts
	LavalinkTimeout         time.Duration
	LavalinkLoadTimeout     time.Duration
	HTTPClientTimeout       time.Duration
	PanelPowerActionTimeout time.Duration

	// RSS
	RSSPollInterval time.Duration
	RSSFeedTimeout  time.Duration

	// URL Tools
	XGDAPIKey      string
	VTAPIKey       string
	ShortenTimeout time.Duration
	ScanTimeout    time.Duration

	// Panel (Pelican)
	PanelURL          string
	PanelAPIKey       string
	PanelAllowedUsers []snowflake.ID

	// Embed Fix
	DeepLAPIKey string
}

func Load() (*Config, error) {
	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("DISCORD_TOKEN is required")
	}

	appIDStr := os.Getenv("DISCORD_APP_ID")
	if appIDStr == "" {
		return nil, fmt.Errorf("DISCORD_APP_ID is required")
	}
	appID, err := snowflake.Parse(appIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid DISCORD_APP_ID: %w", err)
	}

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./config.toml"
	}

	tc := defaultTOMLConfig()
	if _, err := toml.DecodeFile(configPath, &tc); err != nil {
		return nil, fmt.Errorf("failed to load config file %s: %w", configPath, err)
	}

	tc.fillDefaults()

	if err := tc.validate(); err != nil {
		return nil, err
	}

	allowedUsers := make([]snowflake.ID, len(tc.Panel.AllowedUsers))
	for i, id := range tc.Panel.AllowedUsers {
		allowedUsers[i] = snowflake.ID(id)
	}

	cfg := &Config{
		Token:            token,
		AppID:            appID,
		LavalinkHost:     tc.Lavalink.Host,
		LavalinkPassword: tc.Lavalink.Password,
		DataDir:          tc.Storage.DataDir,
		DBPath:           tc.Storage.DBPath,
		DefaultVolume:    tc.Player.DefaultVolume,
		AutoLeaveTimeout: time.Duration(tc.Player.AutoLeaveTimeout) * time.Second,
		PresenceInterval: time.Duration(tc.Presence.Interval) * time.Second,
		LogLevel:         parseSlogLevel(tc.LogLevel),

		LavalinkTimeout:         time.Duration(tc.Timeouts.Lavalink) * time.Second,
		LavalinkLoadTimeout:     time.Duration(tc.Timeouts.LavalinkLoad) * time.Second,
		HTTPClientTimeout:       time.Duration(tc.Timeouts.HTTPClient) * time.Second,
		PanelPowerActionTimeout: time.Duration(tc.Timeouts.PanelPowerAction) * time.Second,

		RSSPollInterval: time.Duration(tc.RSS.PollInterval) * time.Second,
		RSSFeedTimeout:  time.Duration(tc.RSS.FeedTimeout) * time.Second,

		XGDAPIKey:      os.Getenv("XGD_API_KEY"),
		VTAPIKey:       os.Getenv("VT_API_KEY"),
		ShortenTimeout: time.Duration(tc.URL.ShortenTimeout) * time.Second,
		ScanTimeout:    time.Duration(tc.URL.ScanTimeout) * time.Second,

		PanelURL:          tc.Panel.URL,
		PanelAPIKey:       os.Getenv("PANEL_API_KEY"),
		PanelAllowedUsers: allowedUsers,

		DeepLAPIKey: os.Getenv("DEEPL_API_KEY"),
	}

	return cfg, nil
}

func parseSlogLevel(s string) slog.Level {
	switch s {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
