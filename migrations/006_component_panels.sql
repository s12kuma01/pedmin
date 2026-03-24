CREATE TABLE IF NOT EXISTS component_panels (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    guild_id   INTEGER NOT NULL,
    name       TEXT    NOT NULL,
    components TEXT    NOT NULL DEFAULT '[]',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(guild_id, name)
);
CREATE INDEX IF NOT EXISTS idx_component_panels_guild ON component_panels(guild_id);
