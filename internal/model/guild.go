package model

import (
	"errors"
	"time"

	"github.com/disgoorg/snowflake/v2"
)

var ErrDuplicateFeed = errors.New("duplicate feed")

type GuildSettings struct {
	GuildID        snowflake.ID    `json:"guild_id"`
	EnabledModules map[string]bool `json:"enabled_modules"`
	ModuleSettings map[string]any  `json:"module_settings"`
}

type Ticket struct {
	GuildID   snowflake.ID
	Number    int
	ChannelID snowflake.ID
	UserID    snowflake.ID
	Subject   string
	CreatedAt time.Time
	ClosedAt  *time.Time
	ClosedBy  *snowflake.ID
}

type RSSFeed struct {
	ID        int64
	GuildID   snowflake.ID
	URL       string
	ChannelID snowflake.ID
	Title     string
	AddedAt   time.Time
}
