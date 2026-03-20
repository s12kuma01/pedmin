package player

import (
	"context"
	"log/slog"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
)

// OnVoiceStateUpdate implements module.VoiceStateListener.
func (p *Player) OnVoiceStateUpdate(guildID, channelID, userID snowflake.ID) {
	botVoiceState, ok := p.client.Caches.VoiceState(guildID, p.client.ApplicationID)
	if !ok || botVoiceState.ChannelID == nil {
		return
	}
	botChannelID := *botVoiceState.ChannelID

	memberCount := 0
	for vs := range p.client.Caches.VoiceStates(guildID) {
		if vs.ChannelID != nil && *vs.ChannelID == botChannelID && vs.UserID != p.client.ApplicationID {
			memberCount++
		}
	}

	if memberCount == 0 {
		p.startAutoLeaveTimer(guildID)
	} else {
		p.cancelAutoLeaveTimer(guildID)
	}
}

func (p *Player) startAutoLeaveTimer(guildID snowflake.ID) {
	p.cancelAutoLeaveTimer(guildID)

	if p.autoLeaveTimeout == 0 {
		return
	}

	timer := time.AfterFunc(p.autoLeaveTimeout, func() {
		p.logger.Info("auto-leaving voice channel due to inactivity", slog.Any("guild", guildID))
		p.leaveTimers.Delete(guildID)

		if player := p.lavalink.ExistingPlayer(guildID); player != nil {
			ctx, cancel := lavalinkCtx()
			_ = player.Destroy(ctx)
			cancel()
			p.lavalink.RemovePlayer(guildID)
		}

		_ = p.client.UpdateVoiceState(context.Background(), guildID, nil, false, false)
		p.queues.Delete(guildID)

		val, ok := p.messages.Load(guildID)
		if ok {
			tracked, ok := val.(trackedMessage)
			if !ok {
				return
			}
			newPlayer := p.lavalink.Player(guildID)
			queue := p.queues.Get(guildID)
			ui := BuildPlayerUI(newPlayer, queue)
			if _, err := p.client.Rest.UpdateMessage(tracked.channelID, tracked.messageID, discord.NewMessageUpdateV2([]discord.LayoutComponent{ui})); err != nil {
				p.logger.Warn("failed to update player message on auto-leave", slog.Any("error", err))
				p.messages.Delete(guildID)
			}
		}
	})

	p.leaveTimers.Store(guildID, timer)
}

func (p *Player) cancelAutoLeaveTimer(guildID snowflake.ID) {
	val, ok := p.leaveTimers.LoadAndDelete(guildID)
	if !ok {
		return
	}
	timer, ok := val.(*time.Timer)
	if !ok {
		return
	}
	timer.Stop()
}
