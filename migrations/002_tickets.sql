CREATE TABLE IF NOT EXISTS tickets (
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
CREATE INDEX IF NOT EXISTS idx_tickets_channel ON tickets (channel_id);
