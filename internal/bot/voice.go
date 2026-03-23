// opyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package bot

import (
	"context"

	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
	"github.com/s12kuma01/pedmin/config"
	"github.com/s12kuma01/pedmin/internal/module"
)

// Bot's own voice state/server updates are forwarded to Lavalink so it can
// manage the audio connection. Other users' voice state changes are forwarded
// to modules (e.g. player auto-leave on empty VC).

func (b *Bot) onVoiceStateUpdate(e *events.GuildVoiceStateUpdate) {
	if e.VoiceState.UserID != b.Client.ApplicationID {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), config.DefaultLavalinkTimeout)
	defer cancel()
	b.Lavalink.OnVoiceStateUpdate(ctx, e.VoiceState.GuildID, e.VoiceState.ChannelID, e.VoiceState.SessionID)
}

func (b *Bot) onVoiceServerUpdate(e *events.VoiceServerUpdate) {
	ctx, cancel := context.WithTimeout(context.Background(), config.DefaultLavalinkTimeout)
	defer cancel()
	b.Lavalink.OnVoiceServerUpdate(ctx, e.GuildID, e.Token, *e.Endpoint)
}

func (b *Bot) onMemberVoiceStateUpdate(e *events.GuildVoiceStateUpdate) {
	if e.VoiceState.UserID == b.Client.ApplicationID {
		return
	}
	var channelID snowflake.ID
	if e.VoiceState.ChannelID != nil {
		channelID = *e.VoiceState.ChannelID
	}
	for _, m := range b.modules {
		if vsl, ok := m.(module.VoiceStateListener); ok {
			vsl.OnVoiceStateUpdate(e.VoiceState.GuildID, channelID, e.VoiceState.UserID)
		}
	}
}
