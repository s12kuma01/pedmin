# Data Persistence Guide

## GuildStore Interface

```go
type GuildStore interface {
    // Guild settings
    Get(guildID snowflake.ID) (*GuildSettings, error)
    Save(settings *GuildSettings) error
    IsModuleEnabled(guildID snowflake.ID, moduleID string) (bool, error)
    SetModuleEnabled(guildID snowflake.ID, moduleID string, enabled bool) error
    GetModuleSettings(guildID snowflake.ID, moduleID string) (string, error)
    SetModuleSettings(guildID snowflake.ID, moduleID string, settings string) error

    // Tickets
    CreateTicket(guildID snowflake.ID, number int, channelID, userID snowflake.ID, subject string) error
    GetTicketByChannel(channelID snowflake.ID) (*Ticket, error)
    CloseTicket(channelID snowflake.ID, closedBy snowflake.ID) error
    DeleteTicket(channelID snowflake.ID) error

    // RSS
    CreateRSSFeed(feed *RSSFeed) error
    DeleteRSSFeed(id int64, guildID snowflake.ID) error
    GetRSSFeeds(guildID snowflake.ID) ([]RSSFeed, error)
    GetAllRSSFeeds() ([]RSSFeed, error)
    CountRSSFeeds(guildID snowflake.ID) (int, error)
    IsItemSeen(feedID int64, itemHash string) (bool, error)
    MarkItemsSeen(feedID int64, itemHashes []string) error
    PruneSeenItems(olderThan time.Time) error

    // Lifecycle
    Close() error
}
```

All methods are safe for concurrent use.

## GuildSettings Structure

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

## Data Types

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

### Database Location
```
{DB_PATH}  (default: {DATA_DIR}/pedmin.db)
```
Override with the `storage.db_path` setting in `config.toml` or the `DB_PATH` environment variable.

### Schema

```sql
-- Migration 1: Core guild settings
CREATE TABLE schema_migrations (
    version    INTEGER PRIMARY KEY,
    applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

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

-- Migration 2: Tickets
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

-- Migration 3: RSS feeds
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

Managed via `schema_migrations` table. Add new migrations to the `migrations` slice in `sqlite_store.go`:

```go
var migrations = []struct {
    version int
    sql     string
}{
    {version: 1, sql: `CREATE TABLE ...`},
    {version: 2, sql: `CREATE TABLE tickets ...`},
    {version: 3, sql: `CREATE TABLE rss_feeds ...`},
}
```

Applied automatically on `NewSQLiteStore()`.

## Adding a New Store Implementation

1. Implement the `GuildStore` interface (including `Close() error`)
2. Add a constructor function
3. Swap the initialization in `main.go`

## ModuleSettings Usage

Modules use `GetModuleSettings`/`SetModuleSettings` to persist per-module configuration as JSON strings:

```go
// Load
data, _ := store.GetModuleSettings(guildID, "player")
var settings PlayerSettings
json.Unmarshal([]byte(data), &settings)

// Save
raw, _ := json.Marshal(settings)
store.SetModuleSettings(guildID, "player", string(raw))
```

### Per-Module Settings Pattern

Modules like `logger` and `ticket` wrap this in a `LoadSettings` helper:

```go
func LoadSettings(store GuildStore, guildID snowflake.ID) (*Settings, error) {
    data, err := store.GetModuleSettings(guildID, ModuleID)
    if err != nil {
        return nil, err
    }
    var s Settings
    if err := json.Unmarshal([]byte(data), &s); err != nil {
        return &Settings{}, nil
    }
    return &s, nil
}
```
