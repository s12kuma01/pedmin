// Package store defines the GuildStore persistence interface and data types.
package store

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

type GuildStore interface {
	Get(guildID snowflake.ID) (*GuildSettings, error)
	Save(settings *GuildSettings) error
	IsModuleEnabled(guildID snowflake.ID, moduleID string) (bool, error)
	SetModuleEnabled(guildID snowflake.ID, moduleID string, enabled bool) error
	GetModuleSettings(guildID snowflake.ID, moduleID string) (string, error)
	SetModuleSettings(guildID snowflake.ID, moduleID string, settings string) error
	CreateTicket(guildID snowflake.ID, number int, channelID, userID snowflake.ID, subject string) error
	GetTicketByChannel(channelID snowflake.ID) (*Ticket, error)
	CloseTicket(channelID snowflake.ID, closedBy snowflake.ID) error
	DeleteTicket(channelID snowflake.ID) error
	CreateRSSFeed(feed *RSSFeed) error
	DeleteRSSFeed(id int64, guildID snowflake.ID) error
	GetRSSFeeds(guildID snowflake.ID) ([]RSSFeed, error)
	GetAllRSSFeeds() ([]RSSFeed, error)
	CountRSSFeeds(guildID snowflake.ID) (int, error)
	IsItemSeen(feedID int64, itemHash string) (bool, error)
	MarkItemsSeen(feedID int64, itemHashes []string) error
	PruneSeenItems(olderThan time.Time) error
	Close() error
}
