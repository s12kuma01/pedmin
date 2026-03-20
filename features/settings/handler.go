package settings

import (
	"log/slog"
	"strings"

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

	_ = e.CreateMessage(s.mainPanel(*guildID))
}

func (s *Settings) HandleComponent(e *events.ComponentInteractionCreate) {
	customID := e.Data.CustomID()
	_, action, _ := strings.Cut(customID, ":")

	guildID := e.GuildID()
	if guildID == nil {
		return
	}

	switch {
	case action == "select":
		data := e.Data.(discord.StringSelectMenuInteractionData)
		if len(data.Values) == 0 {
			return
		}
		moduleID := data.Values[0]
		_ = e.UpdateMessage(s.modulePanel(*guildID, moduleID))

	case strings.HasPrefix(action, "toggle:"):
		moduleID := strings.TrimPrefix(action, "toggle:")
		enabled := s.bot.IsModuleEnabled(*guildID, moduleID)
		if err := s.bot.SetModuleEnabled(*guildID, moduleID, !enabled); err != nil {
			s.logger.Error("failed to toggle module", slog.Any("error", err))
		}
		_ = e.UpdateMessage(s.modulePanel(*guildID, moduleID))

	case action == "back":
		_ = e.UpdateMessage(s.mainPanelUpdate(*guildID))
	}
}
