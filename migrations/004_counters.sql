CREATE TABLE IF NOT EXISTS counters (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    guild_id   INTEGER NOT NULL,
    word       TEXT    NOT NULL,
    match_type TEXT    NOT NULL DEFAULT 'partial',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(guild_id, word, match_type)
);
CREATE INDEX IF NOT EXISTS idx_counters_guild ON counters(guild_id);

CREATE TABLE IF NOT EXISTS counter_hits (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    counter_id INTEGER NOT NULL REFERENCES counters(id) ON DELETE CASCADE,
    guild_id   INTEGER NOT NULL,
    user_id    INTEGER NOT NULL,
    hit_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_counter_hits_counter ON counter_hits(counter_id);
CREATE INDEX IF NOT EXISTS idx_counter_hits_guild_time ON counter_hits(guild_id, hit_at);
