package service

import (
	"github.com/disgoorg/snowflake/v2"
	"github.com/s12kuma01/pedmin/internal/model"
	"github.com/s12kuma01/pedmin/internal/repository"
)

// LoadSettings loads the ticket settings for a guild.
func (s *TicketService) LoadSettings(guildID snowflake.ID) (*model.TicketSettings, error) {
	return repository.LoadModuleSettings(s.store, guildID, model.TicketModuleID, func() *model.TicketSettings {
		return &model.TicketSettings{}
	})
}

// SaveSettings saves the ticket settings for a guild.
func (s *TicketService) SaveSettings(guildID snowflake.ID, settings *model.TicketSettings) error {
	return repository.SaveModuleSettings(s.store, guildID, model.TicketModuleID, settings)
}

// UpdateCategory updates the ticket category setting for a guild.
func (s *TicketService) UpdateCategory(guildID, categoryID snowflake.ID) error {
	settings, err := s.LoadSettings(guildID)
	if err != nil {
		return err
	}
	settings.CategoryID = categoryID
	return s.SaveSettings(guildID, settings)
}

// UpdateLogChannel updates the log channel setting for a guild.
func (s *TicketService) UpdateLogChannel(guildID, channelID snowflake.ID) error {
	settings, err := s.LoadSettings(guildID)
	if err != nil {
		return err
	}
	settings.LogChannelID = channelID
	return s.SaveSettings(guildID, settings)
}

// UpdateSupportRole updates the support role setting for a guild.
func (s *TicketService) UpdateSupportRole(guildID, roleID snowflake.ID) error {
	settings, err := s.LoadSettings(guildID)
	if err != nil {
		return err
	}
	settings.SupportRoleID = roleID
	return s.SaveSettings(guildID, settings)
}
