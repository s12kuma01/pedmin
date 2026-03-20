package store

import "fmt"

var migrations = []struct {
	version int
	sql     string
}{
	{
		version: 1,
		sql: `
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
		`,
	},
	{
		version: 2,
		sql: `
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
		`,
	},
	{
		version: 3,
		sql: `
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
		`,
	},
}

func (s *SQLiteStore) migrate() error {
	_, err := s.db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (
		version    INTEGER PRIMARY KEY,
		applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	for _, m := range migrations {
		var exists int
		err := s.db.QueryRow("SELECT COUNT(*) FROM schema_migrations WHERE version = ?", m.version).Scan(&exists)
		if err != nil {
			return fmt.Errorf("failed to check migration %d: %w", m.version, err)
		}
		if exists > 0 {
			continue
		}

		tx, err := s.db.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin transaction for migration %d: %w", m.version, err)
		}

		if _, err := tx.Exec(m.sql); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("failed to execute migration %d: %w", m.version, err)
		}

		if _, err := tx.Exec("INSERT INTO schema_migrations (version) VALUES (?)", m.version); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("failed to record migration %d: %w", m.version, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit migration %d: %w", m.version, err)
		}
	}
	return nil
}
