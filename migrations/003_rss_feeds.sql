CREATE TABLE IF NOT EXISTS rss_feeds (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    guild_id   INTEGER NOT NULL,
    url        TEXT    NOT NULL,
    channel_id INTEGER NOT NULL,
    title      TEXT    NOT NULL DEFAULT '',
    added_at   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(guild_id, url)
);
CREATE INDEX IF NOT EXISTS idx_rss_feeds_guild ON rss_feeds(guild_id);

CREATE TABLE IF NOT EXISTS rss_seen_items (
    feed_id    INTEGER NOT NULL REFERENCES rss_feeds(id) ON DELETE CASCADE,
    item_hash  TEXT    NOT NULL,
    seen_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (feed_id, item_hash)
);
