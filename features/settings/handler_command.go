package settings

import (
	"github.com/disgoorg/disgo/events"
	"github.com/s12kuma01/pedmin/ui"
)

func (s *Settings) HandleCommand(e *events.ApplicationCommandInteractionCreate) {
	guildID := e.GuildID()
	if guildID == nil {
		_ = e.CreateMessage(ui.ErrorMessage("設定はサーバー内でのみ使用できます。"))
		return
	}

	options := s.listModuleOptions(*guildID)
	_ = e.CreateMessage(ui.BuildMainPanel(options))
}
