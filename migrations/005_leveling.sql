CREATE TABLE IF NOT EXISTS user_xp (
    guild_id      INTEGER NOT NULL,
    user_id       INTEGER NOT NULL,
    total_xp      INTEGER NOT NULL DEFAULT 0,
    level         INTEGER NOT NULL DEFAULT 0,
    message_count INTEGER NOT NULL DEFAULT 0,
    voice_minutes INTEGER NOT NULL DEFAULT 0,
    last_xp_at    TIMESTAMP,
    PRIMARY KEY (guild_id, user_id)
);
CREATE INDEX IF NOT EXISTS idx_user_xp_ranking ON user_xp(guild_id, total_xp DESC);

CREATE TABLE IF NOT EXISTS level_role_rewards (
    id       INTEGER PRIMARY KEY AUTOINCREMENT,
    guild_id INTEGER NOT NULL,
    level    INTEGER NOT NULL,
    role_id  INTEGER NOT NULL,
    UNIQUE(guild_id, level, role_id)
);
CREATE INDEX IF NOT EXISTS idx_level_role_rewards_guild ON level_role_rewards(guild_id);
