package settings

import (
	"log/slog"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
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

func (s *Settings) HandleComponent(e *events.ComponentInteractionCreate) {
	customID := e.Data.CustomID()
	_, action, _ := strings.Cut(customID, ":")

	guildID := e.GuildID()
	if guildID == nil {
		return
	}

	switch {
	case action == "select":
		data, ok := e.Data.(discord.StringSelectMenuInteractionData)
		if !ok || len(data.Values) == 0 {
			return
		}
		moduleID := data.Values[0]
		_ = e.UpdateMessage(s.buildModulePanel(*guildID, moduleID))

	case strings.HasPrefix(action, "toggle:"):
		moduleID := strings.TrimPrefix(action, "toggle:")
		enabled := s.bot.IsModuleEnabled(*guildID, moduleID)
		if err := s.bot.SetModuleEnabled(*guildID, moduleID, !enabled); err != nil {
			s.logger.Error("failed to toggle module", slog.Any("error", err))
		}
		_ = e.UpdateMessage(s.buildModulePanel(*guildID, moduleID))

	case action == "back":
		options := s.listModuleOptions(*guildID)
		_ = e.UpdateMessage(BuildMainPanelUpdate(options))
	}
}

func (s *Settings) listModuleOptions(guildID snowflake.ID) []ModuleOption {
	modules := s.bot.GetModules()
	var options []ModuleOption
	for _, m := range modules {
		info := m.Info()
		if info.AlwaysOn {
			continue
		}
		options = append(options, ModuleOption{
			ID:          info.ID,
			Name:        info.Name,
			Description: info.Description,
			Enabled:     s.bot.IsModuleEnabled(guildID, info.ID),
		})
	}
	return options
}

func (s *Settings) buildModulePanel(guildID snowflake.ID, moduleID string) discord.MessageUpdate {
	modules := s.bot.GetModules()
	m, ok := modules[moduleID]
	if !ok {
		return BuildModuleNotFound()
	}

	info := m.Info()
	enabled := s.bot.IsModuleEnabled(guildID, moduleID)
	settingsPanel := m.SettingsPanel(guildID)
	return BuildModulePanel(info, enabled, settingsPanel)
}
