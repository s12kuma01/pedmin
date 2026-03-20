# Data Persistence Guide

## GuildStore Interface

```go
type GuildStore interface {
    // Guild settings
    Get(guildID snowflake.ID) (*GuildSettings, error)
    Save(settings *GuildSettings) error
    IsModuleEnabled(guildID snowflake.ID, moduleID string) (bool, error)
    SetModuleEnabled(guildID snowflake.ID, moduleID string, enabled bool) error
    GetModuleSettings(guildID snowflake.ID, moduleID string) (map[string]any, error)
    SetModuleSettings(guildID snowflake.ID, moduleID string, settings map[string]any) error

    // Tickets
    CreateTicket(ticket *Ticket) error
    GetTicketByChannel(channelID snowflake.ID) (*Ticket, error)
    CloseTicket(channelID snowflake.ID) error
    DeleteTicket(channelID snowflake.ID) error

    // RSS
    CreateRSSFeed(feed *RSSFeed) error
    DeleteRSSFeed(id int64) error
    GetRSSFeeds(guildID snowflake.ID) ([]RSSFeed, error)
    GetAllRSSFeeds() ([]RSSFeed, error)
    CountRSSFeeds(guildID snowflake.ID) (int, error)
    IsItemSeen(feedID int64, guid string) (bool, error)
    MarkItemsSeen(feedID int64, guids []string) error
    PruneSeenItems(feedID int64, keep int) error

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
| `ModuleSettings` | `map[string]any` | Module ID → arbitrary settings |

Default: when no record exists, returns empty maps with all modules disabled.

## Data Types

### Ticket
```go
type Ticket struct {
    GuildID   snowflake.ID
    ChannelID snowflake.ID
    UserID    snowflake.ID
    ClosedAt  *time.Time
}
```

### RSSFeed
```go
type RSSFeed struct {
    ID        int64
    GuildID   snowflake.ID
    ChannelID snowflake.ID
    URL       string
}
```

## SQLiteStore Implementation

### Database Location
```
{DB_PATH}  (default: {DATA_DIR}/pedmin.db)
```
Override with the `storage.dbPath` setting in `config.cue` or the `DB_PATH` environment variable.

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
    channel_id INTEGER NOT NULL PRIMARY KEY,
    user_id    INTEGER NOT NULL,
    closed_at  TIMESTAMP
);
CREATE INDEX idx_tickets_channel ON tickets(channel_id);

-- Migration 3: RSS feeds
CREATE TABLE rss_feeds (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    guild_id   INTEGER NOT NULL,
    channel_id INTEGER NOT NULL,
    url        TEXT    NOT NULL
);

CREATE TABLE rss_seen_items (
    feed_id    INTEGER NOT NULL REFERENCES rss_feeds(id) ON DELETE CASCADE,
    guid       TEXT    NOT NULL,
    seen_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (feed_id, guid)
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

```go
settings, _ := store.Get(guildID)
if val, ok := settings.ModuleSettings["player"]; ok {
    playerSettings := val.(map[string]any)
}

settings.ModuleSettings["player"] = map[string]any{
    "default_volume": 50,
}
store.Save(settings)
```

Note: JSON serialization means numbers become `float64`, nested objects become `map[string]any`.

### Per-Module Settings Pattern

Modules like `logger` and `ticket` use `GetModuleSettings`/`SetModuleSettings` with typed settings structs:

```go
func LoadSettings(store GuildStore, guildID snowflake.ID) (*Settings, error) {
    raw, err := store.GetModuleSettings(guildID, ModuleID)
    if err != nil {
        return nil, err
    }
    // Marshal raw map to JSON, then unmarshal to typed struct
    var s Settings
    data, _ := json.Marshal(raw)
    json.Unmarshal(data, &s)
    return &s, nil
}
```
