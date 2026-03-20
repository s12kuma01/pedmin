package ticket

import "github.com/disgoorg/snowflake/v2"

// UpdateCategory updates the ticket category setting for a guild.
func (t *Ticket) UpdateCategory(guildID, categoryID snowflake.ID) error {
	settings, err := LoadSettings(t.store, guildID)
	if err != nil {
		return err
	}
	settings.CategoryID = categoryID
	return SaveSettings(t.store, guildID, settings)
}

// UpdateLogChannel updates the log channel setting for a guild.
func (t *Ticket) UpdateLogChannel(guildID, channelID snowflake.ID) error {
	settings, err := LoadSettings(t.store, guildID)
	if err != nil {
		return err
	}
	settings.LogChannelID = channelID
	return SaveSettings(t.store, guildID, settings)
}

// UpdateSupportRole updates the support role setting for a guild.
func (t *Ticket) UpdateSupportRole(guildID, roleID snowflake.ID) error {
	settings, err := LoadSettings(t.store, guildID)
	if err != nil {
		return err
	}
	settings.SupportRoleID = roleID
	return SaveSettings(t.store, guildID, settings)
}
