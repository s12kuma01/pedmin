package ticket

import (
	"github.com/disgoorg/snowflake/v2"
	"github.com/s12kuma01/pedmin/store"
)

type TicketSettings struct {
	CategoryID    snowflake.ID `json:"category_id"`
	LogChannelID  snowflake.ID `json:"log_channel_id"`
	SupportRoleID snowflake.ID `json:"support_role_id"`
	NextNumber    int          `json:"next_number"`
}

func LoadSettings(guildStore store.GuildStore, guildID snowflake.ID) (*TicketSettings, error) {
	return store.LoadModuleSettings(guildStore, guildID, ModuleID, func() *TicketSettings {
		return &TicketSettings{}
	})
}

func SaveSettings(guildStore store.GuildStore, guildID snowflake.ID, s *TicketSettings) error {
	return store.SaveModuleSettings(guildStore, guildID, ModuleID, s)
}
