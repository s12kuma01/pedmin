// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/disgoorg/snowflake/v2"

	"github.com/s12kuma01/pedmin/internal/repository"
)

type voiceSession struct {
	GuildID   snowflake.ID
	UserID    snowflake.ID
	JoinedAt  time.Time
	LastAward time.Time
}

// OnVoiceStateUpdate handles a user joining or leaving a voice channel.
func (s *LevelingService) OnVoiceStateUpdate(guildID, channelID, userID snowflake.ID) {
	key := fmt.Sprintf("%d:%d", guildID, userID)

	if channelID == 0 {
		// User left voice
		if val, ok := s.voiceSessions.LoadAndDelete(key); ok {
			session := val.(*voiceSession)
			s.awardPendingVoiceXP(session)
		}
		return
	}

	// User joined or moved channels — create/update session
	if _, ok := s.voiceSessions.Load(key); !ok {
		now := time.Now()
		s.voiceSessions.Store(key, &voiceSession{
			GuildID:   guildID,
			UserID:    userID,
			JoinedAt:  now,
			LastAward: now,
		})
	}
}

// StartVoiceTicker starts the background goroutine that awards voice XP every minute.
func (s *LevelingService) StartVoiceTicker(ctx context.Context) {
	ctx, s.voiceCancel = context.WithCancel(ctx)
	go s.voiceTickLoop(ctx)
}

// StopVoiceTicker stops the voice XP ticker.
func (s *LevelingService) StopVoiceTicker() {
	if s.voiceCancel != nil {
		s.voiceCancel()
	}
}

// Shutdown stops the voice ticker and awards remaining voice XP.
func (s *LevelingService) Shutdown() {
	s.StopVoiceTicker()

	// Award remaining XP to all active sessions
	s.voiceSessions.Range(func(key, value any) bool {
		session := value.(*voiceSession)
		s.awardPendingVoiceXP(session)
		s.voiceSessions.Delete(key)
		return true
	})
}

func (s *LevelingService) voiceTickLoop(ctx context.Context) {
	s.logger.Info("voice xp ticker started")

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			s.logger.Info("voice xp ticker stopped")
			return
		case <-ticker.C:
			s.awardVoiceXP()
		}
	}
}

func (s *LevelingService) awardVoiceXP() {
	now := time.Now()
	var updates []repository.VoiceXPUpdate

	s.voiceSessions.Range(func(key, value any) bool {
		session := value.(*voiceSession)

		if now.Sub(session.LastAward) < 1*time.Minute {
			return true
		}

		settings, err := s.LoadSettings(session.GuildID)
		if err != nil || settings.VoiceXPPerMinute <= 0 {
			return true
		}

		minutes := int(now.Sub(session.LastAward).Minutes())
		if minutes < 1 {
			return true
		}

		xp := minutes * settings.VoiceXPPerMinute
		updates = append(updates, repository.VoiceXPUpdate{
			GuildID:  session.GuildID,
			UserID:   session.UserID,
			Minutes:  minutes,
			XPAmount: xp,
		})
		session.LastAward = now
		return true
	})

	if len(updates) == 0 {
		return
	}

	// Get old levels before batch update for level-up detection
	type oldLevelInfo struct {
		guildID snowflake.ID
		userID  snowflake.ID
		old     int
	}
	var levelChecks []oldLevelInfo
	for _, u := range updates {
		ux, err := s.store.GetUserXP(u.GuildID, u.UserID)
		if err != nil {
			continue
		}
		levelChecks = append(levelChecks, oldLevelInfo{u.GuildID, u.UserID, ux.Level})
	}

	if err := s.store.BatchAddVoiceXP(updates); err != nil {
		s.logger.Error("failed to batch add voice xp", slog.Any("error", err))
		return
	}

	// Check for level-ups
	for _, lc := range levelChecks {
		ux, err := s.store.GetUserXP(lc.guildID, lc.userID)
		if err != nil {
			continue
		}
		if ux.Level > lc.old {
			settings, err := s.LoadSettings(lc.guildID)
			if err != nil {
				continue
			}
			// For voice level-ups, use notification channel if set, otherwise skip
			var channelID snowflake.ID
			if settings.NotificationMode == "channel" && settings.NotificationChID != 0 {
				channelID = settings.NotificationChID
			}
			if channelID != 0 {
				s.onLevelUp(lc.guildID, lc.userID, channelID, ux.Level, settings)
			}
		}
	}
}

func (s *LevelingService) awardPendingVoiceXP(session *voiceSession) {
	elapsed := time.Since(session.LastAward)
	minutes := int(elapsed.Minutes())
	if minutes < 1 {
		return
	}

	settings, err := s.LoadSettings(session.GuildID)
	if err != nil || settings.VoiceXPPerMinute <= 0 {
		return
	}

	xp := minutes * settings.VoiceXPPerMinute
	_, _, err = s.store.AddXP(session.GuildID, session.UserID, xp, true)
	if err != nil {
		s.logger.Error("failed to award pending voice xp", slog.Any("error", err))
	}
}
