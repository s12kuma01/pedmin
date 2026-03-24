// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package model

// Module IDs used as CustomID prefixes and settings keys.
const (
	PlayerModuleID     = "player"
	EmbedFixModuleID   = "embedfix"
	TicketModuleID     = "ticket"
	LoggerModuleID     = "logger"
	RSSModuleID        = "rss"
	PanelModuleID      = "panel"
	URLModuleID        = "url"
	TranslatorModuleID = "translator"
	SettingsModuleID   = "settings"
	PingModuleID       = "ping"
	AvatarModuleID     = "avatar"
	FuckfetchModuleID  = "fuckfetch"
	CounterModuleID    = "counter"
	LevelingModuleID   = "leveling"
	AutoroleModuleID   = "autorole"
	BuilderModuleID    = "builder"

	// MaxRSSFeedsPerGuild limits the number of RSS feeds per guild.
	MaxRSSFeedsPerGuild = 10

	// MaxCountersPerGuild limits the number of word counters per guild.
	MaxCountersPerGuild = 10

	// MaxPanelsPerGuild limits the number of component panels per guild.
	MaxPanelsPerGuild = 10

	// MaxComponentsPerPanel limits the number of components per panel.
	MaxComponentsPerPanel = 20
)
