// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package service

import (
	"context"
	"fmt"
	"log/slog"
	"crypto/rand"
	"math/big"
	"sync"
	"time"

	disgobot "github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"

	"github.com/Sumire-Labs/pedmin/internal/model"
	"github.com/Sumire-Labs/pedmin/internal/repository"
)

// LevelingService handles XP processing, leveling logic, and settings management.
type LevelingService struct {
	store     repository.GuildStore
	client    *disgobot.Client
	logger    *slog.Logger
	cooldowns sync.Map // key: "guildID:userID" → time.Time

	// Voice tracking fields are in leveling_voice.go
	voiceSessions sync.Map // key: "guildID:userID" → *voiceSession
	voiceCancel   context.CancelFunc
}

// NewLevelingService creates a new LevelingService.
func NewLevelingService(client *disgobot.Client, store repository.GuildStore, logger *slog.Logger) *LevelingService {
	return &LevelingService{
		store:  store,
		client: client,
		logger: logger,
	}
}

// LoadSettings loads leveling settings for a guild.
func (s *LevelingService) LoadSettings(guildID snowflake.ID) (*model.LevelingSettings, error) {
	return repository.LoadModuleSettings(s.store, guildID, model.LevelingModuleID, model.DefaultLevelingSettings)
}

// SaveSettings saves leveling settings for a guild.
func (s *LevelingService) SaveSettings(guildID snowflake.ID, settings *model.LevelingSettings) error {
	return repository.SaveModuleSettings(s.store, guildID, model.LevelingModuleID, settings)
}

// ProcessMessage handles XP gain from a message.
func (s *LevelingService) ProcessMessage(guildID, userID, channelID snowflake.ID, memberRoles []snowflake.ID) {
	key := fmt.Sprintf("%d:%d", guildID, userID)

	// Check in-memory cooldown
	if val, ok := s.cooldowns.Load(key); ok {
		if lastXP, ok := val.(time.Time); ok {
			settings, err := s.LoadSettings(guildID)
			if err != nil {
				s.logger.Error("failed to load leveling settings", slog.Any("error", err))
				return
			}
			if time.Since(lastXP) < time.Duration(settings.CooldownSeconds)*time.Second {
				return
			}
		}
	}

	settings, err := s.LoadSettings(guildID)
	if err != nil {
		s.logger.Error("failed to load leveling settings", slog.Any("error", err))
		return
	}

	// Generate random XP
	baseXP := settings.MinXP
	if settings.MaxXP > settings.MinXP {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(settings.MaxXP-settings.MinXP+1)))
		if err == nil {
			baseXP = settings.MinXP + int(n.Int64())
		}
	}

	// Apply multipliers
	multiplier := highestRoleMultiplier(memberRoles, settings) * channelMultiplier(channelID, settings)
	finalXP := int(float64(baseXP) * multiplier)
	if finalXP < 1 {
		finalXP = 1
	}

	// Add XP
	ux, oldLevel, err := s.store.AddXP(guildID, userID, finalXP, false)
	if err != nil {
		s.logger.Error("failed to add xp", slog.Any("error", err))
		return
	}

	// Record cooldown
	s.cooldowns.Store(key, time.Now())

	// Check for level-up
	if ux.Level > oldLevel {
		s.onLevelUp(guildID, userID, channelID, ux.Level, settings)
	}
}

func (s *LevelingService) onLevelUp(guildID, userID, channelID snowflake.ID, newLevel int, settings *model.LevelingSettings) {
	// Assign role rewards
	rewards, err := s.store.GetRoleRewards(guildID)
	if err != nil {
		s.logger.Error("failed to get role rewards", slog.Any("error", err))
	}

	for _, r := range rewards {
		if r.Level <= newLevel {
			if err := s.client.Rest.AddMemberRole(guildID, userID, r.RoleID); err != nil {
				s.logger.Error("failed to assign level role",
					slog.Int("level", r.Level),
					slog.Any("error", err),
				)
			}
		}
	}

	// Send level-up notification
	if settings.NotificationMode == "off" {
		return
	}

	targetChannel := channelID
	if settings.NotificationMode == "channel" && settings.NotificationChID != 0 {
		targetChannel = settings.NotificationChID
	}

	text := fmt.Sprintf("<@%d> がレベル **%d** に到達しました!", userID, newLevel)

	// Mention newly earned roles for this exact level
	for _, r := range rewards {
		if r.Level == newLevel {
			text += fmt.Sprintf("\n<@&%d> を獲得!", r.RoleID)
		}
	}

	_, err = s.client.Rest.CreateMessage(targetChannel, discord.NewMessageCreate().
		WithContent(text).
		WithAllowedMentions(&discord.AllowedMentions{
			Users: []snowflake.ID{userID},
		}))
	if err != nil {
		s.logger.Error("failed to send level-up notification", slog.Any("error", err))
	}
}

// GetUserXP returns a user's XP record.
func (s *LevelingService) GetUserXP(guildID, userID snowflake.ID) (*model.UserXP, error) {
	return s.store.GetUserXP(guildID, userID)
}

// GetLeaderboard returns the leaderboard for a guild.
func (s *LevelingService) GetLeaderboard(guildID snowflake.ID, limit, offset int) ([]model.LeaderboardEntry, error) {
	return s.store.GetLeaderboard(guildID, limit, offset)
}

// GetUserRank returns a user's rank in the guild.
func (s *LevelingService) GetUserRank(guildID, userID snowflake.ID) (int, error) {
	return s.store.GetUserRank(guildID, userID)
}

// GetRoleRewards returns all role rewards for a guild.
func (s *LevelingService) GetRoleRewards(guildID snowflake.ID) ([]model.LevelRoleReward, error) {
	return s.store.GetRoleRewards(guildID)
}

// AddRoleReward adds a role reward.
func (s *LevelingService) AddRoleReward(guildID snowflake.ID, level int, roleID snowflake.ID) error {
	return s.store.AddRoleReward(guildID, level, roleID)
}

// RemoveRoleReward removes a role reward.
func (s *LevelingService) RemoveRoleReward(id int64, guildID snowflake.ID) error {
	return s.store.RemoveRoleReward(id, guildID)
}

// CountRoleRewards returns the number of role rewards for a guild.
func (s *LevelingService) CountRoleRewards(guildID snowflake.ID) (int, error) {
	return s.store.CountRoleRewards(guildID)
}

func highestRoleMultiplier(memberRoles []snowflake.ID, settings *model.LevelingSettings) float64 {
	highest := 1.0
	for _, rm := range settings.RoleMultipliers {
		for _, mr := range memberRoles {
			if rm.RoleID == mr && rm.Multiplier > highest {
				highest = rm.Multiplier
			}
		}
	}
	return highest
}

func channelMultiplier(channelID snowflake.ID, settings *model.LevelingSettings) float64 {
	for _, cm := range settings.ChannelMultipliers {
		if cm.ChannelID == channelID {
			return cm.Multiplier
		}
	}
	return 1.0
}
