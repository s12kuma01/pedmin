// Package config loads application configuration from environment variables and CUE files.
package config

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"github.com/disgoorg/snowflake/v2"
)

type Config struct {
	// env (secrets)
	Token string
	AppID snowflake.ID

	// CUE (app settings)
	LavalinkHost     string
	LavalinkPassword string
	DataDir          string
	DBPath           string
	DefaultVolume    int
	AutoLeaveTimeout time.Duration
	PresenceInterval time.Duration
	LogLevel         slog.Level

	// Panel (Pelican)
	PanelURL          string
	PanelAPIKey       string
	PanelAllowedUsers []snowflake.ID

	// URL Tools
	XGDAPIKey string
	VTAPIKey  string
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
		configPath = "./config.cue"
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", configPath, err)
	}

	ctx := cuecontext.New()
	value := ctx.CompileBytes(data)
	if err := value.Validate(cue.Concrete(true)); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	cfg := &Config{
		Token: token,
		AppID: appID,
	}

	if err := lookupString(value, "lavalink.host", &cfg.LavalinkHost); err != nil {
		return nil, err
	}
	if err := lookupString(value, "lavalink.password", &cfg.LavalinkPassword); err != nil {
		return nil, err
	}
	if err := lookupString(value, "storage.dataDir", &cfg.DataDir); err != nil {
		return nil, err
	}
	if err := lookupString(value, "storage.dbPath", &cfg.DBPath); err != nil {
		return nil, err
	}

	var defaultVolume int
	if err := lookupInt(value, "player.defaultVolume", &defaultVolume); err != nil {
		return nil, err
	}
	cfg.DefaultVolume = defaultVolume

	var autoLeaveSeconds int
	if err := lookupInt(value, "player.autoLeaveTimeout", &autoLeaveSeconds); err != nil {
		return nil, err
	}
	cfg.AutoLeaveTimeout = time.Duration(autoLeaveSeconds) * time.Second

	var presenceSeconds int
	if err := lookupInt(value, "presence.interval", &presenceSeconds); err != nil {
		return nil, err
	}
	cfg.PresenceInterval = time.Duration(presenceSeconds) * time.Second

	var logLevelStr string
	if err := lookupString(value, "logLevel", &logLevelStr); err != nil {
		return nil, err
	}
	cfg.LogLevel = parseSlogLevel(logLevelStr)

	// URL Tools (optional)
	cfg.XGDAPIKey = os.Getenv("XGD_API_KEY")
	cfg.VTAPIKey = os.Getenv("VT_API_KEY")

	// Panel (optional)
	cfg.PanelAPIKey = os.Getenv("PANEL_API_KEY")
	_ = lookupString(value, "panel.url", &cfg.PanelURL)
	cfg.PanelAllowedUsers = lookupSnowflakeList(value, "panel.allowedUsers")

	return cfg, nil
}

func lookupString(v cue.Value, path string, dst *string) error {
	val := v.LookupPath(cue.ParsePath(path))
	s, err := val.String()
	if err != nil {
		return fmt.Errorf("config %s: %w", path, err)
	}
	*dst = s
	return nil
}

func lookupInt(v cue.Value, path string, dst *int) error {
	val := v.LookupPath(cue.ParsePath(path))
	n, err := val.Int64()
	if err != nil {
		return fmt.Errorf("config %s: %w", path, err)
	}
	*dst = int(n)
	return nil
}

func lookupSnowflakeList(v cue.Value, path string) []snowflake.ID {
	iter, err := v.LookupPath(cue.ParsePath(path)).List()
	if err != nil {
		return nil
	}
	var ids []snowflake.ID
	for iter.Next() {
		n, err := iter.Value().Int64()
		if err == nil {
			ids = append(ids, snowflake.ID(n))
		}
	}
	return ids
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
