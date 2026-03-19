# Data Persistence Guide

## GuildStore Interface

```go
type GuildStore interface {
    Get(guildID snowflake.ID) (*GuildSettings, error)
    Save(settings *GuildSettings) error
    IsModuleEnabled(guildID snowflake.ID, moduleID string) (bool, error)
    SetModuleEnabled(guildID snowflake.ID, moduleID string, enabled bool) error
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

## SQLiteStore Implementation

### Database Location
```
{DB_PATH}  (default: {DATA_DIR}/pedmin.db)
```
Override with the `DB_PATH` environment variable.

### Schema

```sql
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

### Schema Migrations

Managed via `schema_migrations` table. Add new migrations to the `migrations` slice in `sqlite_store.go`:

```go
var migrations = []struct {
    version int
    sql     string
}{
    {version: 1, sql: `CREATE TABLE ...`},
    {version: 2, sql: `ALTER TABLE ...`},  // new migrations here
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
