CREATE TABLE IF NOT EXISTS guild_modules (
    guild_id   INTEGER NOT NULL,
    module_id  TEXT    NOT NULL,
    enabled    BOOLEAN NOT NULL DEFAULT 0,
    PRIMARY KEY (guild_id, module_id)
);

CREATE TABLE IF NOT EXISTS guild_module_settings (
    guild_id   INTEGER NOT NULL,
    module_id  TEXT    NOT NULL,
    settings   TEXT    NOT NULL DEFAULT '{}',
    PRIMARY KEY (guild_id, module_id)
);
