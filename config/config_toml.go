// opyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package config

import (
	"fmt"
	"path/filepath"
)

type tomlConfig struct {
	Lavalink tomlLavalink `toml:"lavalink"`
	Storage  tomlStorage  `toml:"storage"`
	Player   tomlPlayer   `toml:"player"`
	Presence tomlPresence `toml:"presence"`
	Timeouts tomlTimeouts `toml:"timeouts"`
	RSS      tomlRSS      `toml:"rss"`
	URL      tomlURL      `toml:"url"`
	LogLevel string       `toml:"log_level"`
	Panel    tomlPanel    `toml:"panel"`
}

type tomlLavalink struct {
	Host     string `toml:"host"`
	Password string `toml:"password"`
}

type tomlStorage struct {
	DataDir string `toml:"data_dir"`
	DBPath  string `toml:"db_path"`
}

type tomlPlayer struct {
	DefaultVolume    int `toml:"default_volume"`
	AutoLeaveTimeout int `toml:"auto_leave_timeout"`
}

type tomlPresence struct {
	Interval int `toml:"interval"`
}

type tomlTimeouts struct {
	Lavalink         int `toml:"lavalink"`
	LavalinkLoad     int `toml:"lavalink_load"`
	HTTPClient       int `toml:"http_client"`
	PanelPowerAction int `toml:"panel_power_action"`
}

type tomlRSS struct {
	PollInterval int `toml:"poll_interval"`
	FeedTimeout  int `toml:"feed_timeout"`
}

type tomlURL struct {
	ShortenTimeout int `toml:"shorten_timeout"`
	ScanTimeout    int `toml:"scan_timeout"`
}

type tomlPanel struct {
	URL          string  `toml:"url"`
	AllowedUsers []int64 `toml:"allowed_users"`
}

func defaultTOMLConfig() tomlConfig {
	return tomlConfig{
		Lavalink: tomlLavalink{
			Host:     "lavalink:2333",
			Password: "youshallnotpass",
		},
		Storage: tomlStorage{
			DataDir: "./data",
			DBPath:  "",
		},
		Player: tomlPlayer{
			DefaultVolume:    50,
			AutoLeaveTimeout: 180,
		},
		Presence: tomlPresence{
			Interval: 30,
		},
		Timeouts: tomlTimeouts{
			Lavalink:         2,
			LavalinkLoad:     10,
			HTTPClient:       10,
			PanelPowerAction: 15,
		},
		RSS: tomlRSS{
			PollInterval: 300,
			FeedTimeout:  30,
		},
		URL: tomlURL{
			ShortenTimeout: 10,
			ScanTimeout:    15,
		},
		LogLevel: "info",
		Panel: tomlPanel{
			URL:          "",
			AllowedUsers: nil,
		},
	}
}

func (tc *tomlConfig) fillDefaults() {
	if tc.Storage.DBPath == "" {
		tc.Storage.DBPath = filepath.Join(tc.Storage.DataDir, "pedmin.db")
	}
}

func (tc *tomlConfig) validate() error {
	if err := checkRange("player.default_volume", tc.Player.DefaultVolume, 0, 200); err != nil {
		return err
	}
	if err := checkMin("player.auto_leave_timeout", tc.Player.AutoLeaveTimeout, 0); err != nil {
		return err
	}
	if err := checkRange("presence.interval", tc.Presence.Interval, 5, 600); err != nil {
		return err
	}
	if err := checkRange("timeouts.lavalink", tc.Timeouts.Lavalink, 1, 30); err != nil {
		return err
	}
	if err := checkRange("timeouts.lavalink_load", tc.Timeouts.LavalinkLoad, 5, 60); err != nil {
		return err
	}
	if err := checkRange("timeouts.http_client", tc.Timeouts.HTTPClient, 5, 60); err != nil {
		return err
	}
	if err := checkRange("timeouts.panel_power_action", tc.Timeouts.PanelPowerAction, 5, 60); err != nil {
		return err
	}
	if err := checkRange("rss.poll_interval", tc.RSS.PollInterval, 60, 3600); err != nil {
		return err
	}
	if err := checkRange("rss.feed_timeout", tc.RSS.FeedTimeout, 10, 120); err != nil {
		return err
	}
	if err := checkRange("url.shorten_timeout", tc.URL.ShortenTimeout, 5, 60); err != nil {
		return err
	}
	if err := checkRange("url.scan_timeout", tc.URL.ScanTimeout, 5, 120); err != nil {
		return err
	}

	switch tc.LogLevel {
	case "debug", "info", "warn", "error":
	default:
		return fmt.Errorf("config log_level: must be one of debug, info, warn, error (got %q)", tc.LogLevel)
	}

	return nil
}

func checkRange(name string, val, min, max int) error {
	if val < min || val > max {
		return fmt.Errorf("config %s: must be between %d and %d (got %d)", name, min, max, val)
	}
	return nil
}

func checkMin(name string, val, min int) error {
	if val < min {
		return fmt.Errorf("config %s: must be >= %d (got %d)", name, min, val)
	}
	return nil
}
