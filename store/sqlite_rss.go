package store

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/disgoorg/snowflake/v2"
)

func (s *SQLiteStore) CreateRSSFeed(feed *RSSFeed) error {
	err := s.db.QueryRow(
		"INSERT INTO rss_feeds (guild_id, url, channel_id, title) VALUES (?, ?, ?, ?) RETURNING id",
		int64(feed.GuildID), feed.URL, int64(feed.ChannelID), feed.Title,
	).Scan(&feed.ID)
	if err != nil && strings.Contains(err.Error(), "UNIQUE constraint") {
		return fmt.Errorf("%w: %s", ErrDuplicateFeed, feed.URL)
	}
	return err
}

func (s *SQLiteStore) DeleteRSSFeed(id int64, guildID snowflake.ID) error {
	// Delete seen items first (SQLite foreign key support varies)
	_, _ = s.db.Exec("DELETE FROM rss_seen_items WHERE feed_id = ?", id)
	_, err := s.db.Exec("DELETE FROM rss_feeds WHERE id = ? AND guild_id = ?", id, int64(guildID))
	return err
}

func (s *SQLiteStore) GetRSSFeeds(guildID snowflake.ID) ([]RSSFeed, error) {
	rows, err := s.db.Query(
		"SELECT id, guild_id, url, channel_id, title, added_at FROM rss_feeds WHERE guild_id = ? ORDER BY added_at",
		int64(guildID),
	)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	return scanRSSFeeds(rows)
}

func (s *SQLiteStore) GetAllRSSFeeds() ([]RSSFeed, error) {
	rows, err := s.db.Query("SELECT id, guild_id, url, channel_id, title, added_at FROM rss_feeds ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	return scanRSSFeeds(rows)
}

func scanRSSFeeds(rows *sql.Rows) ([]RSSFeed, error) {
	var feeds []RSSFeed
	for rows.Next() {
		var f RSSFeed
		var gid, chid int64
		if err := rows.Scan(&f.ID, &gid, &f.URL, &chid, &f.Title, &f.AddedAt); err != nil {
			return nil, err
		}
		f.GuildID = snowflake.ID(gid)
		f.ChannelID = snowflake.ID(chid)
		feeds = append(feeds, f)
	}
	return feeds, rows.Err()
}

func (s *SQLiteStore) CountRSSFeeds(guildID snowflake.ID) (int, error) {
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM rss_feeds WHERE guild_id = ?", int64(guildID)).Scan(&count)
	return count, err
}

func (s *SQLiteStore) IsItemSeen(feedID int64, itemHash string) (bool, error) {
	var exists int
	err := s.db.QueryRow(
		"SELECT 1 FROM rss_seen_items WHERE feed_id = ? AND item_hash = ?",
		feedID, itemHash,
	).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *SQLiteStore) MarkItemsSeen(feedID int64, itemHashes []string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	stmt, err := tx.Prepare("INSERT OR IGNORE INTO rss_seen_items (feed_id, item_hash) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer func() { _ = stmt.Close() }()

	for _, h := range itemHashes {
		if _, err := stmt.Exec(feedID, h); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (s *SQLiteStore) PruneSeenItems(olderThan time.Time) error {
	_, err := s.db.Exec("DELETE FROM rss_seen_items WHERE seen_at < ?", olderThan)
	return err
}
