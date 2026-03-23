// opyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package model

import "github.com/disgoorg/snowflake/v2"

// TicketSettings holds per-guild ticket configuration.
type TicketSettings struct {
	CategoryID    snowflake.ID `json:"category_id"`
	LogChannelID  snowflake.ID `json:"log_channel_id"`
	SupportRoleID snowflake.ID `json:"support_role_id"`
	NextNumber    int          `json:"next_number"`
}
