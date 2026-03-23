# Data Persistence Guide

## GuildStore Interface

Defined in `internal/repository/repository.go`. Data types are in `internal/model/guild.go`.

```go
type GuildStore interface {
    // Guild settings
    Get(guildID snowflake.ID) (*model.GuildSettings, error)
    Save(settings *model.GuildSettings) error
    IsModuleEnabled(guildID snowflake.ID, moduleID string) (bool, error)
    SetModuleEnabled(guildID snowflake.ID, moduleID string, enabled bool) error
    GetModuleSettings(guildID snowflake.ID, moduleID string) (string, error)
    SetModuleSettings(guildID snowflake.ID, moduleID string, settings string) error

    // Tickets
    CreateTicket(guildID snowflake.ID, number int, channelID, userID snowflake.ID, subject string) error
    GetTicketByChannel(channelID snowflake.ID) (*model.Ticket, error)
    CloseTicket(channelID snowflake.ID, closedBy snowflake.ID) error
    DeleteTicket(channelID snowflake.ID) error

    // RSS
    CreateRSSFeed(feed *model.RSSFeed) error
    DeleteRSSFeed(id int64, guildID snowflake.ID) error
    GetRSSFeeds(guildID snowflake.ID) ([]model.RSSFeed, error)
    GetAllRSSFeeds() ([]model.RSSFeed, error)
    CountRSSFeeds(guildID snowflake.ID) (int, error)
    IsItemSeen(feedID int64, itemHash string) (bool, error)
    MarkItemsSeen(feedID int64, itemHashes []string) error
    PruneSeenItems(olderThan time.Time) error

    // Lifecycle
    Close() error
}
```

All methods are safe for concurrent use.

## Data Types

Defined in `internal/model/guild.go`:

### GuildSettings

```go
type GuildSettings struct {
    GuildID        snowflake.ID          `json:"guild_id"`
    EnabledModules map[string]bool       `json:"enabled_modules"`
    ModuleSettings map[string]any        `json:"module_settings"`
}
```

| Field | Type | Description |
|-------|------|-------------|
| `GuildID` | `snowflake.ID` | Discord guild ID |
| `EnabledModules` | `map[string]bool` | Module ID → enabled state |
| `ModuleSettings` | `map[string]any` | Module ID → arbitrary settings (internal use by `Get`/`Save`) |

Default: when no record exists, returns empty maps with all modules disabled.

### Ticket
```go
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
```

### RSSFeed
```go
type RSSFeed struct {
    ID        int64
    GuildID   snowflake.ID
    URL       string
    ChannelID snowflake.ID
    Title     string
    AddedAt   time.Time
}
```

## SQLiteStore Implementation

Located in `internal/repository/sqlite.go`.

### Database Location
```
{DB_PATH}  (default: {DATA_DIR}/pedmin.db)
```
Override with the `storage.db_path` setting in `config.toml` or the `DB_PATH` environment variable.

### Schema

SQL migration files are in `migrations/` and loaded via `embed.FS`.

```sql
-- 001_guild_modules.sql: Core guild settings
CREATE TABLE guild_modules (
    guild_id   INTEGER NOT NULL,
    module_id  TEXT    NOT NULL,
    enabled    BOOLEAN NOT NULL DEFAULT 0,
    PRIMARY KEY (guild_id, module_id)
);

CREATE TABLE guild_module_settings (
    guild_id   INTEGER NOT NULL,
    module_id  TEXT    NOT NULL,
    settings   TEXT    NOT NULL DEFAULT '{}',
    PRIMARY KEY (guild_id, module_id)
);

-- 002_tickets.sql: Tickets
CREATE TABLE tickets (
    guild_id   INTEGER NOT NULL,
    number     INTEGER NOT NULL,
    channel_id INTEGER NOT NULL,
    user_id    INTEGER NOT NULL,
    subject    TEXT    NOT NULL DEFAULT '',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    closed_at  TIMESTAMP,
    closed_by  INTEGER,
    PRIMARY KEY (guild_id, number)
);
CREATE INDEX idx_tickets_channel ON tickets (channel_id);

-- 003_rss_feeds.sql: RSS feeds
CREATE TABLE rss_feeds (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    guild_id   INTEGER NOT NULL,
    url        TEXT    NOT NULL,
    channel_id INTEGER NOT NULL,
    title      TEXT    NOT NULL DEFAULT '',
    added_at   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(guild_id, url)
);
CREATE INDEX idx_rss_feeds_guild ON rss_feeds(guild_id);

CREATE TABLE rss_seen_items (
    feed_id    INTEGER NOT NULL REFERENCES rss_feeds(id) ON DELETE CASCADE,
    item_hash  TEXT    NOT NULL,
    seen_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (feed_id, item_hash)
);
```

### Configuration
- **WAL mode**: Concurrent readers with single writer
- **busy_timeout=5000**: 5 second wait on lock contention
- **Pure Go driver**: `modernc.org/sqlite` (no CGO)

### Performance
| Operation | Query | Complexity |
|-----------|-------|------------|
| `IsModuleEnabled` | Single-row SELECT by PK | O(1) |
| `SetModuleEnabled` | Single UPSERT | O(1) |
| `Get` | Two queries (modules + settings) | O(n modules) |
| `Save` | Transaction with UPSERTs | O(n modules) |
| `GetTicketByChannel` | Single-row SELECT by PK | O(1) |
| `GetRSSFeeds` | SELECT by guild_id | O(n feeds) |
| `IsItemSeen` | Single-row SELECT by PK | O(1) |

### Schema Migrations

Managed via `schema_migrations` table. SQL files are stored in `migrations/` and loaded via `embed.FS` at startup.
Add new migrations as `migrations/NNN_description.sql`. The version number is parsed from the filename prefix.

Applied automatically on `repository.NewSQLiteStore()`.

## Adding a New Store Implementation

1. Implement the `repository.GuildStore` interface (including `Close() error`)
2. Add a constructor function
3. Swap the initialization in `cmd/pedmin/main.go`

## ModuleSettings Usage

Use the generic helpers in `internal/repository/module_settings.go`:

```go
// internal/service/myfeature.go
func (s *MyFeatureService) LoadSettings(guildID snowflake.ID) (*model.MyFeatureSettings, error) {
    return repository.LoadModuleSettings(s.store, guildID, model.MyFeatureModuleID, func() *model.MyFeatureSettings {
        return &model.MyFeatureSettings{}
    })
}

func (s *MyFeatureService) SaveSettings(guildID snowflake.ID, settings *model.MyFeatureSettings) error {
    return repository.SaveModuleSettings(s.store, guildID, model.MyFeatureModuleID, settings)
}
```

The generic helpers handle JSON marshaling/unmarshaling automatically with a default factory fallback.
