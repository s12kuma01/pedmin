// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package model

import (
	"time"

	"github.com/disgoorg/snowflake/v2"
)

// MaxLevel is the maximum attainable level.
const MaxLevel = 100

// MaxRoleRewardsPerGuild limits role reward entries per guild.
const MaxRoleRewardsPerGuild = 25

// MaxMultipliersPerGuild limits role/channel multiplier entries per guild.
const MaxMultipliersPerGuild = 10

// XPForLevel returns the XP required to complete level n (from n to n+1).
// Formula from Noctaly: 5n² + 20n + 100.
func XPForLevel(n int) int {
	return 5*n*n + 20*n + 100
}

// TotalXPForLevel returns cumulative XP needed to reach level n from 0.
func TotalXPForLevel(n int) int {
	total := 0
	for i := 0; i < n; i++ {
		total += XPForLevel(i)
	}
	return total
}

// LevelFromTotalXP returns the level and remaining XP within that level for a given total XP.
func LevelFromTotalXP(totalXP int) (level int, remaining int) {
	for level = 0; level < MaxLevel; level++ {
		needed := XPForLevel(level)
		if totalXP < needed {
			return level, totalXP
		}
		totalXP -= needed
	}
	return MaxLevel, 0
}

// LevelingSettings holds per-guild leveling configuration (stored as JSON).
type LevelingSettings struct {
	MinXP              int                 `json:"min_xp"`
	MaxXP              int                 `json:"max_xp"`
	CooldownSeconds    int                 `json:"cooldown_seconds"`
	NotificationMode   string              `json:"notification_mode"`
	NotificationChID   snowflake.ID        `json:"notification_ch_id"`
	VoiceXPPerMinute   int                 `json:"voice_xp_per_minute"`
	RoleMultipliers    []RoleMultiplier    `json:"role_multipliers"`
	ChannelMultipliers []ChannelMultiplier `json:"channel_multipliers"`
}

// RoleMultiplier maps a role to an XP multiplier.
type RoleMultiplier struct {
	RoleID     snowflake.ID `json:"role_id"`
	Multiplier float64      `json:"multiplier"`
}

// ChannelMultiplier maps a channel to an XP multiplier.
type ChannelMultiplier struct {
	ChannelID  snowflake.ID `json:"channel_id"`
	Multiplier float64      `json:"multiplier"`
}

// DefaultLevelingSettings returns the default settings.
func DefaultLevelingSettings() *LevelingSettings {
	return &LevelingSettings{
		MinXP:            15,
		MaxXP:            25,
		CooldownSeconds:  60,
		NotificationMode: "same",
		VoiceXPPerMinute: 5,
	}
}

// NotificationModeLabel returns the Japanese label for a notification mode.
func NotificationModeLabel(mode string) string {
	switch mode {
	case "same":
		return "同じチャンネル"
	case "channel":
		return "指定チャンネル"
	case "off":
		return "オフ"
	default:
		return mode
	}
}

// UserXP represents a user's XP record in a guild.
type UserXP struct {
	GuildID      snowflake.ID
	UserID       snowflake.ID
	TotalXP      int
	Level        int
	MessageCount int
	VoiceMinutes int
	LastXPAt     *time.Time
}

// CurrentXP returns the XP within the current level.
func (u *UserXP) CurrentXP() int {
	_, remaining := LevelFromTotalXP(u.TotalXP)
	return remaining
}

// NeededXP returns the XP needed to reach the next level.
func (u *UserXP) NeededXP() int {
	return XPForLevel(u.Level)
}

// LevelRoleReward maps a level to a role.
type LevelRoleReward struct {
	ID      int64
	GuildID snowflake.ID
	Level   int
	RoleID  snowflake.ID
}

// LeaderboardEntry holds data for the leaderboard display.
type LeaderboardEntry struct {
	UserID  snowflake.ID
	Level   int
	TotalXP int
	Rank    int
}
