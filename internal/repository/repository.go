// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

// Package repository defines the GuildStore persistence interface and generic helpers.
package repository

import (
	"encoding/json"
	"time"

	"github.com/disgoorg/snowflake/v2"

	"github.com/s12kuma01/pedmin/internal/model"
)

// SettingsStore handles guild-level module configuration.
type SettingsStore interface {
	Get(guildID snowflake.ID) (*model.GuildSettings, error)
	Save(settings *model.GuildSettings) error
	IsModuleEnabled(guildID snowflake.ID, moduleID string) (bool, error)
	SetModuleEnabled(guildID snowflake.ID, moduleID string, enabled bool) error
	GetModuleSettings(guildID snowflake.ID, moduleID string) (string, error)
	SetModuleSettings(guildID snowflake.ID, moduleID string, settings string) error
}

// TicketStore handles ticket persistence.
type TicketStore interface {
	CreateTicket(guildID snowflake.ID, number int, channelID, userID snowflake.ID, subject string) error
	GetTicketByChannel(channelID snowflake.ID) (*model.Ticket, error)
	CloseTicket(channelID snowflake.ID, closedBy snowflake.ID) error
	DeleteTicket(channelID snowflake.ID) error
}

// RSSStore handles RSS feed and seen-item persistence.
type RSSStore interface {
	CreateRSSFeed(feed *model.RSSFeed) error
	DeleteRSSFeed(id int64, guildID snowflake.ID) error
	GetRSSFeeds(guildID snowflake.ID) ([]model.RSSFeed, error)
	GetAllRSSFeeds() ([]model.RSSFeed, error)
	CountRSSFeeds(guildID snowflake.ID) (int, error)
	IsItemSeen(feedID int64, itemHash string) (bool, error)
	MarkItemsSeen(feedID int64, itemHashes []string) error
	PruneSeenItems(olderThan time.Time) error
}

// CounterStore handles word counter persistence.
type CounterStore interface {
	CreateCounter(counter *model.Counter) error
	DeleteCounter(id int64, guildID snowflake.ID) error
	GetCounters(guildID snowflake.ID) ([]model.Counter, error)
	GetCounter(id int64, guildID snowflake.ID) (*model.Counter, error)
	CountCounters(guildID snowflake.ID) (int, error)
	RecordHits(hits []CounterHit) error
	GetCounterStats(guildID snowflake.ID, since *time.Time) ([]model.CounterStat, error)
	GetCounterUserRanking(counterID int64, since *time.Time, limit int) ([]model.CounterUserRank, error)
}

// CounterHit is a single hit event for batch recording.
type CounterHit struct {
	CounterID int64
	GuildID   snowflake.ID
	UserID    snowflake.ID
}

// LevelingStore handles user XP and role reward persistence.
type LevelingStore interface {
	GetUserXP(guildID, userID snowflake.ID) (*model.UserXP, error)
	AddXP(guildID, userID snowflake.ID, amount int, isVoice bool) (*model.UserXP, int, error)
	GetLeaderboard(guildID snowflake.ID, limit, offset int) ([]model.LeaderboardEntry, error)
	GetUserRank(guildID, userID snowflake.ID) (int, error)
	GetRoleRewards(guildID snowflake.ID) ([]model.LevelRoleReward, error)
	AddRoleReward(guildID snowflake.ID, level int, roleID snowflake.ID) error
	RemoveRoleReward(id int64, guildID snowflake.ID) error
	CountRoleRewards(guildID snowflake.ID) (int, error)
	BatchAddVoiceXP(updates []VoiceXPUpdate) error
}

// VoiceXPUpdate is a single voice XP grant for batch processing.
type VoiceXPUpdate struct {
	GuildID  snowflake.ID
	UserID   snowflake.ID
	Minutes  int
	XPAmount int
}

// PanelBuilderStore handles component panel persistence.
type PanelBuilderStore interface {
	CreatePanel(panel *model.ComponentPanel) error
	UpdatePanel(panel *model.ComponentPanel) error
	DeletePanel(id int64, guildID snowflake.ID) error
	GetPanel(id int64, guildID snowflake.ID) (*model.ComponentPanel, error)
	GetPanels(guildID snowflake.ID) ([]model.ComponentPanel, error)
	CountPanels(guildID snowflake.ID) (int, error)
}

// GuildStore is the composite persistence interface embedding all domain stores.
type GuildStore interface {
	SettingsStore
	TicketStore
	RSSStore
	CounterStore
	LevelingStore
	PanelBuilderStore
	Close() error
}

// LoadModuleSettings loads and unmarshals module-specific settings from the store.
// If the data is missing or invalid, defaultFn provides the fallback value.
func LoadModuleSettings[T any](gs GuildStore, guildID snowflake.ID, moduleID string, defaultFn func() *T) (*T, error) {
	data, err := gs.GetModuleSettings(guildID, moduleID)
	if err != nil {
		return nil, err
	}
	var s T
	if err := json.Unmarshal([]byte(data), &s); err != nil {
		return defaultFn(), nil
	}
	return &s, nil
}

// SaveModuleSettings marshals and persists module-specific settings to the store.
func SaveModuleSettings[T any](gs GuildStore, guildID snowflake.ID, moduleID string, settings *T) error {
	data, err := json.Marshal(settings)
	if err != nil {
		return err
	}
	return gs.SetModuleSettings(guildID, moduleID, string(data))
}
