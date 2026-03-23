package model

import "github.com/disgoorg/snowflake/v2"

// TicketSettings holds per-guild ticket configuration.
type TicketSettings struct {
	CategoryID    snowflake.ID `json:"category_id"`
	LogChannelID  snowflake.ID `json:"log_channel_id"`
	SupportRoleID snowflake.ID `json:"support_role_id"`
	NextNumber    int          `json:"next_number"`
}
