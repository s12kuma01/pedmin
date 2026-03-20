package settings

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

func (s *Settings) HandleCommand(e *events.ApplicationCommandInteractionCreate) {
	guildID := e.GuildID()
	if guildID == nil {
		_ = e.CreateMessage(ephemeralV2(
			discord.NewContainer(
				discord.NewTextDisplay("設定はサーバー内でのみ使用できます。"),
			),
		))
		return
	}

	options := s.listModuleOptions(*guildID)
	_ = e.CreateMessage(BuildMainPanel(options))
}
