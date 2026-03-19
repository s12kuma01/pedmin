package avatar

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

func (a *Avatar) HandleCommand(e *events.ApplicationCommandInteractionCreate) {
	data := e.SlashCommandInteractionData()

	var user discord.User
	var member *discord.ResolvedMember

	if optUser, ok := data.OptUser("user"); ok {
		user = optUser
		if m, ok := data.OptMember("user"); ok {
			member = &m
		}
	} else {
		user = e.User()
		if m := e.Member(); m != nil {
			member = m
		}
	}

	guildID := e.GuildID()
	ui := BuildAvatarUI(user, member, guildID)

	_ = e.CreateMessage(discord.NewMessageCreateV2(ui))
}
