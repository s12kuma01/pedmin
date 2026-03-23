// opyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

// Package config loads application configuration from environment variables.
package config

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/disgoorg/snowflake/v2"
)

// Config holds all application configuration loaded from environment variables.
type Config struct {
	// Discord
	Token string
	AppID snowflake.ID

	// Lavalink
	LavalinkHost     string
	LavalinkPassword string

	// Storage
	DBPath string

	// API Keys
	DeepLAPIKey string
	XGDAPIKey   string
	VTAPIKey    string
	PanelAPIKey string

	// Panel (Pelican)
	PanelURL          string
	PanelAllowedUsers []snowflake.ID

	// Logging
	LogLevel slog.Level
}

// Load reads configuration from environment variables.
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

	dataDir := envOrDefault("DATA_DIR", "./data")
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = filepath.Join(dataDir, "pedmin.db")
	}

	cfg := &Config{
		Token:            token,
		AppID:            appID,
		LavalinkHost:     envOrDefault("LAVALINK_HOST", "lavalink:2333"),
		LavalinkPassword: envOrDefault("LAVALINK_PASSWORD", "youshallnotpass"),
		DBPath:           dbPath,
		DeepLAPIKey:      os.Getenv("DEEPL_API_KEY"),
		XGDAPIKey:        os.Getenv("XGD_API_KEY"),
		VTAPIKey:         os.Getenv("VT_API_KEY"),
		PanelAPIKey:      os.Getenv("PANEL_API_KEY"),
		PanelURL:         os.Getenv("PANEL_URL"),
		PanelAllowedUsers: parseSnowflakeList(os.Getenv("PANEL_ALLOWED_USERS")),
		LogLevel:         parseSlogLevel(envOrDefault("LOG_LEVEL", "info")),
	}

	return cfg, nil
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func parseSlogLevel(s string) slog.Level {
	switch strings.ToLower(s) {
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

func parseSnowflakeList(s string) []snowflake.ID {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	var ids []snowflake.ID
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		n, err := strconv.ParseInt(p, 10, 64)
		if err != nil {
			continue
		}
		ids = append(ids, snowflake.ID(n))
	}
	return ids
}
