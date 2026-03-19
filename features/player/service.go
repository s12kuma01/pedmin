package player

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/disgolink/v3/lavalink"
	"github.com/disgoorg/snowflake/v2"
)

const lavalinkTimeout = 2 * time.Second

func lavalinkCtx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), lavalinkTimeout)
}

func (p *Player) handleSkip(e *events.ComponentInteractionCreate, guildID snowflake.ID) {
	player := p.lavalink.ExistingPlayer(guildID)
	if player == nil {
		_ = e.DeferUpdateMessage()
		return
	}

	queue := p.queues.Get(guildID)
	next, ok := queue.Next()
	if !ok {
		ctx, cancel := lavalinkCtx()
		defer cancel()
		_ = player.Update(ctx, lavalink.WithNullTrack())
		p.respondWithPlayerUpdate(e, player, guildID)
		return
	}

	ctx, cancel := lavalinkCtx()
	defer cancel()
	if err := player.Update(ctx, lavalink.WithTrack(next)); err != nil {
		p.logger.Error("failed to skip", slog.Any("error", err))
	}
	p.respondWithPlayerUpdate(e, player, guildID)
}

func (p *Player) handleStop(e *events.ComponentInteractionCreate, guildID snowflake.ID) {
	player := p.lavalink.ExistingPlayer(guildID)
	if player != nil {
		ctx, cancel := lavalinkCtx()
		_ = player.Destroy(ctx)
		cancel()
		p.lavalink.RemovePlayer(guildID)
	}
	p.queues.Delete(guildID)
	_ = e.Client().UpdateVoiceState(context.Background(), guildID, nil, false, false)

	queue := p.queues.Get(guildID)
	newPlayer := p.lavalink.Player(guildID)
	ui := BuildPlayerUI(newPlayer, queue)
	_ = e.UpdateMessage(discord.NewMessageUpdateV2([]discord.LayoutComponent{ui}))
}

func (p *Player) handleLoop(e *events.ComponentInteractionCreate, guildID snowflake.ID) {
	queue := p.queues.Get(guildID)
	queue.CycleLoop()

	player := p.lavalink.Player(guildID)
	p.respondWithPlayerUpdate(e, player, guildID)
}

func (p *Player) loadAndPlay(e *events.ModalSubmitInteractionCreate, guildID snowflake.ID, query string) {
	if !strings.HasPrefix(query, "http") {
		query = "ytsearch:" + query
	}

	node := p.lavalink.BestNode()
	if node == nil {
		p.logger.Error("no lavalink node available")
		return
	}

	loadCtx, loadCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer loadCancel()
	node.LoadTracksHandler(loadCtx, query, disgolink.NewResultHandler(
		func(track lavalink.Track) {
			p.playTrack(e, guildID, track)
		},
		func(playlist lavalink.Playlist) {
			if len(playlist.Tracks) == 0 {
				return
			}
			queue := p.queues.Get(guildID)
			queue.Add(playlist.Tracks...)
			queue.SetCurrent(0)

			_ = p.joinVoiceChannel(e.Client(), guildID, e.Member().User.ID)
			player := p.lavalink.Player(guildID)
			ctx, cancel := lavalinkCtx()
			_ = player.Update(ctx, lavalink.WithTrack(playlist.Tracks[0]))
			cancel()
		},
		func(tracks []lavalink.Track) {
			if len(tracks) == 0 {
				return
			}
			p.playTrack(e, guildID, tracks[0])
		},
		func() {
			p.logger.Info("no matches found", slog.String("query", query))
		},
		func(err error) {
			p.logger.Error("load failed", slog.Any("error", err))
		},
	))
}

func (p *Player) playTrack(e *events.ModalSubmitInteractionCreate, guildID snowflake.ID, track lavalink.Track) {
	queue := p.queues.Get(guildID)
	queue.Add(track)

	if queue.Len() == 1 {
		queue.SetCurrent(0)
	}

	_ = p.joinVoiceChannel(e.Client(), guildID, e.Member().User.ID)
	player := p.lavalink.Player(guildID)

	if player.Track() == nil {
		current, ok := queue.Current()
		if !ok {
			current = track
		}
		ctx, cancel := lavalinkCtx()
		_ = player.Update(ctx, lavalink.WithTrack(current))
		cancel()
	}
}

func (p *Player) respondWithPlayerUpdate(e *events.ComponentInteractionCreate, player disgolink.Player, guildID snowflake.ID) {
	queue := p.queues.Get(guildID)
	ui := BuildPlayerUI(player, queue)
	_ = e.UpdateMessage(discord.NewMessageUpdateV2([]discord.LayoutComponent{ui}))
}

func (p *Player) trackMessage(guildID, channelID, messageID snowflake.ID) {
	p.messages.Store(guildID, trackedMessage{
		channelID: channelID,
		messageID: messageID,
	})
}

func (p *Player) deleteTrackedMessage(guildID snowflake.ID) {
	val, ok := p.messages.LoadAndDelete(guildID)
	if !ok {
		return
	}
	tracked := val.(trackedMessage)
	if err := p.client.Rest.DeleteMessage(tracked.channelID, tracked.messageID); err != nil {
		p.logger.Warn("failed to delete tracked message", slog.Any("error", err))
	}
}

func (p *Player) updatePlayerMessage(player disgolink.Player) {
	guildID := player.GuildID()
	val, ok := p.messages.Load(guildID)
	if !ok {
		return
	}
	tracked := val.(trackedMessage)

	queue := p.queues.Get(guildID)
	ui := BuildPlayerUI(player, queue)
	if _, err := p.client.Rest.UpdateMessage(tracked.channelID, tracked.messageID, discord.NewMessageUpdateV2([]discord.LayoutComponent{ui})); err != nil {
		p.logger.Warn("failed to update player message", slog.Any("error", err))
		p.messages.Delete(guildID)
	}
}
