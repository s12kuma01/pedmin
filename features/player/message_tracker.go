package player

import (
	"log/slog"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/snowflake/v2"
)

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
	tracked, ok := val.(trackedMessage)
	if !ok {
		return
	}
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
	tracked, ok := val.(trackedMessage)
	if !ok {
		return
	}

	queue := p.queues.Get(guildID)
	ui := BuildPlayerUI(player, queue)
	if _, err := p.client.Rest.UpdateMessage(tracked.channelID, tracked.messageID, discord.NewMessageUpdateV2([]discord.LayoutComponent{ui})); err != nil {
		p.logger.Warn("failed to update player message", slog.Any("error", err))
		p.messages.Delete(guildID)
	}
}
