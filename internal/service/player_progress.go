package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/snowflake/v2"
	"github.com/s12kuma01/pedmin/internal/model"
	"github.com/s12kuma01/pedmin/internal/view"
)

const progressTickInterval = 10 * time.Second

// StartProgressTicker starts a per-guild ticker that periodically updates
// the player message to refresh the progress bar. Calling it when a ticker
// already exists for the guild cancels the old one first.
func (s *PlayerService) StartProgressTicker(guildID snowflake.ID) {
	s.StopProgressTicker(guildID)

	ctx, cancel := context.WithCancel(context.Background())
	s.progressTickers.Store(guildID, cancel)

	go s.progressTickLoop(ctx, guildID)
}

// StopProgressTicker stops the progress ticker for a guild.
func (s *PlayerService) StopProgressTicker(guildID snowflake.ID) {
	val, ok := s.progressTickers.LoadAndDelete(guildID)
	if !ok {
		return
	}
	cancel, ok := val.(context.CancelFunc)
	if !ok {
		return
	}
	cancel()
}

// StopAllProgressTickers cancels all active tickers. Called during shutdown.
func (s *PlayerService) StopAllProgressTickers() {
	s.progressTickers.Range(func(key, value any) bool {
		if cancel, ok := value.(context.CancelFunc); ok {
			cancel()
		}
		s.progressTickers.Delete(key)
		return true
	})
}

func (s *PlayerService) progressTickLoop(ctx context.Context, guildID snowflake.ID) {
	ticker := time.NewTicker(progressTickInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			player := s.lavalink.ExistingPlayer(guildID)
			if player == nil || player.Track() == nil || player.Paused() {
				continue
			}
			s.UpdateTrackedPlayer(player)
		}
	}
}

// TrackMessage stores the tracked message for a guild.
func (s *PlayerService) TrackMessage(guildID, channelID, messageID snowflake.ID) {
	s.messages.Store(guildID, model.TrackedMessage{
		ChannelID: channelID,
		MessageID: messageID,
	})
}

// DeleteTrackedMessage deletes the tracked player message for a guild.
func (s *PlayerService) DeleteTrackedMessage(guildID snowflake.ID) {
	val, ok := s.messages.LoadAndDelete(guildID)
	if !ok {
		return
	}
	tracked, ok := val.(model.TrackedMessage)
	if !ok {
		return
	}
	if err := s.client.Rest.DeleteMessage(tracked.ChannelID, tracked.MessageID); err != nil {
		s.logger.Warn("failed to delete tracked message", slog.Any("error", err))
	}
}

// UpdateTrackedPlayer updates the tracked player message with current state.
func (s *PlayerService) UpdateTrackedPlayer(player disgolink.Player) {
	guildID := player.GuildID()
	val, ok := s.messages.Load(guildID)
	if !ok {
		return
	}
	tracked, ok := val.(model.TrackedMessage)
	if !ok {
		return
	}

	queue := s.queues.Get(guildID)
	ui := view.BuildPlayerUI(player, queue)
	if _, err := s.client.Rest.UpdateMessage(tracked.ChannelID, tracked.MessageID, discord.NewMessageUpdateV2([]discord.LayoutComponent{ui})); err != nil {
		s.logger.Warn("failed to update player message", slog.Any("error", err))
		s.messages.Delete(guildID)
	}
}
