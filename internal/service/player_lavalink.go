package service

import (
	"context"
	"log/slog"

	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/disgolink/v3/lavalink"
)

// SetupPlayerListeners registers Lavalink event listeners on the player service.
func SetupPlayerListeners(link disgolink.Client, svc *PlayerService) {
	link.AddListeners(
		disgolink.NewListenerFunc(svc.onTrackStart),
		disgolink.NewListenerFunc(svc.onTrackEnd),
		disgolink.NewListenerFunc(svc.onTrackException),
		disgolink.NewListenerFunc(svc.onTrackStuck),
		disgolink.NewListenerFunc(svc.onWebSocketClosed),
	)
}

// ConnectNode connects to a Lavalink node.
func ConnectNode(ctx context.Context, link disgolink.Client, host, password string) error {
	_, err := link.AddNode(ctx, disgolink.NodeConfig{
		Name:     "main",
		Address:  host,
		Password: password,
	})
	return err
}

func (s *PlayerService) onTrackStart(player disgolink.Player, event lavalink.TrackStartEvent) {
	s.logger.Info("track started",
		slog.String("title", event.Track.Info.Title),
		slog.Int64("guild", int64(event.GuildID())),
	)
	s.StartProgressTicker(event.GuildID())
	s.UpdateTrackedPlayer(player)
}

func (s *PlayerService) onTrackEnd(player disgolink.Player, event lavalink.TrackEndEvent) {
	if event.Reason != lavalink.TrackEndReasonFinished && event.Reason != lavalink.TrackEndReasonLoadFailed {
		return
	}

	queue := s.queues.Get(event.GuildID())
	next, ok := queue.Next()
	if !ok {
		s.StopProgressTicker(event.GuildID())
		s.UpdateTrackedPlayer(player)
		return
	}

	ctx, cancel := s.LavalinkCtx()
	defer cancel()
	if err := player.Update(ctx, lavalink.WithTrack(next)); err != nil {
		s.logger.Error("failed to play next track", slog.Any("error", err))
	}
}

func (s *PlayerService) onTrackException(player disgolink.Player, event lavalink.TrackExceptionEvent) {
	s.logger.Error("track exception",
		slog.String("title", event.Track.Info.Title),
		slog.String("message", event.Exception.Message),
	)
}

func (s *PlayerService) onTrackStuck(player disgolink.Player, event lavalink.TrackStuckEvent) {
	s.logger.Warn("track stuck",
		slog.String("title", event.Track.Info.Title),
	)

	queue := s.queues.Get(event.GuildID())
	next, ok := queue.Next()
	if !ok {
		return
	}
	ctx, cancel := s.LavalinkCtx()
	defer cancel()
	_ = player.Update(ctx, lavalink.WithTrack(next))
}

func (s *PlayerService) onWebSocketClosed(_ disgolink.Player, event lavalink.WebSocketClosedEvent) {
	s.logger.Warn("websocket closed",
		slog.Int("code", event.Code),
		slog.String("reason", event.Reason),
		slog.Bool("by_remote", event.ByRemote),
	)
}
