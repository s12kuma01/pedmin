package settings

import (
	"log/slog"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
	"github.com/s12kuma01/pedmin/module"
)

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
		opt := ModuleOption{
			ID:          info.ID,
			Name:        info.Name,
			Description: info.Description,
			Enabled:     s.bot.IsModuleEnabled(guildID, info.ID),
		}
		if summarizer, ok := m.(module.SettingsSummarizer); ok {
			opt.Summary = summarizer.SettingsSummary(guildID)
		}
		options = append(options, opt)
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
