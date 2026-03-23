// opyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
	"github.com/s12kuma01/pedmin/internal/model"
	"github.com/s12kuma01/pedmin/internal/view"
)

// OnVoiceStateUpdate checks if the bot's voice channel is empty and starts/cancels the auto-leave timer.
func (s *PlayerService) OnVoiceStateUpdate(guildID, channelID, userID snowflake.ID) {
	botVoiceState, ok := s.client.Caches.VoiceState(guildID, s.client.ApplicationID)
	if !ok || botVoiceState.ChannelID == nil {
		return
	}
	botChannelID := *botVoiceState.ChannelID

	memberCount := 0
	for vs := range s.client.Caches.VoiceStates(guildID) {
		if vs.ChannelID != nil && *vs.ChannelID == botChannelID && vs.UserID != s.client.ApplicationID {
			memberCount++
		}
	}

	if memberCount == 0 {
		s.startAutoLeaveTimer(guildID)
	} else {
		s.cancelAutoLeaveTimer(guildID)
	}
}

func (s *PlayerService) startAutoLeaveTimer(guildID snowflake.ID) {
	s.cancelAutoLeaveTimer(guildID)

	if s.autoLeaveTimeout == 0 {
		return
	}

	timer := time.AfterFunc(s.autoLeaveTimeout, func() {
		s.logger.Info("auto-leaving voice channel due to inactivity", slog.Any("guild", guildID))
		s.leaveTimers.Delete(guildID)
		s.StopProgressTicker(guildID)

		if player := s.lavalink.ExistingPlayer(guildID); player != nil {
			ctx, cancel := s.LavalinkCtx()
			_ = player.Destroy(ctx)
			cancel()
			s.lavalink.RemovePlayer(guildID)
		}

		_ = s.client.UpdateVoiceState(context.Background(), guildID, nil, false, false)
		s.queues.Delete(guildID)

		val, ok := s.messages.Load(guildID)
		if ok {
			tracked, ok := val.(model.TrackedMessage)
			if !ok {
				return
			}
			newPlayer := s.lavalink.Player(guildID)
			queue := s.queues.Get(guildID)
			ui := view.BuildPlayerUI(newPlayer, queue)
			if _, err := s.client.Rest.UpdateMessage(tracked.ChannelID, tracked.MessageID, discord.NewMessageUpdateV2([]discord.LayoutComponent{ui})); err != nil {
				s.logger.Warn("failed to update player message on auto-leave", slog.Any("error", err))
				s.messages.Delete(guildID)
			}
		}
	})

	s.leaveTimers.Store(guildID, timer)
}

func (s *PlayerService) cancelAutoLeaveTimer(guildID snowflake.ID) {
	val, ok := s.leaveTimers.LoadAndDelete(guildID)
	if !ok {
		return
	}
	timer, ok := val.(*time.Timer)
	if !ok {
		return
	}
	timer.Stop()
}
