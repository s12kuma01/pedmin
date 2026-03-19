package player

import (
	"context"
	"log/slog"

	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/disgolink/v3/lavalink"
)

func SetupListeners(link disgolink.Client, p *Player) {
	link.AddListeners(
		disgolink.NewListenerFunc(p.onTrackStart),
		disgolink.NewListenerFunc(p.onTrackEnd),
		disgolink.NewListenerFunc(p.onTrackException),
		disgolink.NewListenerFunc(p.onTrackStuck),
		disgolink.NewListenerFunc(p.onWebSocketClosed),
	)
}

func (p *Player) onTrackStart(player disgolink.Player, event lavalink.TrackStartEvent) {
	p.logger.Info("track started",
		slog.String("title", event.Track.Info.Title),
		slog.Int64("guild", int64(event.GuildID())),
	)
	p.updatePlayerMessage(player)
}

func (p *Player) onTrackEnd(player disgolink.Player, event lavalink.TrackEndEvent) {
	if event.Reason != lavalink.TrackEndReasonFinished && event.Reason != lavalink.TrackEndReasonLoadFailed {
		return
	}

	queue := p.queues.Get(event.GuildID())
	next, ok := queue.Next()
	if !ok {
		p.updatePlayerMessage(player)
		return
	}

	ctx, cancel := lavalinkCtx()
	defer cancel()
	if err := player.Update(ctx, lavalink.WithTrack(next)); err != nil {
		p.logger.Error("failed to play next track", slog.Any("error", err))
	}
}

func (p *Player) onTrackException(player disgolink.Player, event lavalink.TrackExceptionEvent) {
	p.logger.Error("track exception",
		slog.String("title", event.Track.Info.Title),
		slog.String("message", event.Exception.Message),
	)
}

func (p *Player) onTrackStuck(player disgolink.Player, event lavalink.TrackStuckEvent) {
	p.logger.Warn("track stuck",
		slog.String("title", event.Track.Info.Title),
	)

	queue := p.queues.Get(event.GuildID())
	next, ok := queue.Next()
	if !ok {
		return
	}
	ctx, cancel := lavalinkCtx()
	defer cancel()
	_ = player.Update(ctx, lavalink.WithTrack(next))
}

func (p *Player) onWebSocketClosed(_ disgolink.Player, event lavalink.WebSocketClosedEvent) {
	p.logger.Warn("websocket closed",
		slog.Int("code", event.Code),
		slog.String("reason", event.Reason),
		slog.Bool("by_remote", event.ByRemote),
	)
}

func ConnectNode(ctx context.Context, link disgolink.Client, host, password string) error {
	_, err := link.AddNode(ctx, disgolink.NodeConfig{
		Name:     "main",
		Address:  host,
		Password: password,
	})
	return err
}
