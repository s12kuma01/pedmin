package player

import (
	"context"
	"errors"

	disgobot "github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/snowflake/v2"
)

var errNotInVoiceChannel = errors.New("user not in voice channel")

func (p *Player) joinVoiceChannel(client *disgobot.Client, guildID, userID snowflake.ID) error {
	voiceState, ok := client.Caches.VoiceState(guildID, userID)
	if !ok || voiceState.ChannelID == nil {
		return errNotInVoiceChannel
	}
	return client.UpdateVoiceState(context.Background(), guildID, voiceState.ChannelID, false, true)
}
