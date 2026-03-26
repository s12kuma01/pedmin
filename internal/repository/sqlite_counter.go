// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/disgoorg/snowflake/v2"

	"github.com/Sumire-Labs/pedmin/internal/model"
)

func (s *SQLiteStore) CreateCounter(counter *model.Counter) error {
	count, err := s.CountCounters(counter.GuildID)
	if err != nil {
		return fmt.Errorf("failed to check counter count: %w", err)
	}
	if count >= model.MaxCountersPerGuild {
		return fmt.Errorf("カウンター数が上限(%d)に達しています", model.MaxCountersPerGuild)
	}

	var existingID int64
	err = s.db.QueryRow(
		"SELECT id FROM counters WHERE guild_id = ? AND word = ? AND match_type = ?",
		int64(counter.GuildID), counter.Word, string(counter.MatchType),
	).Scan(&existingID)
	if err == nil {
		return fmt.Errorf("同じワードとマッチタイプのカウンターが既に登録されています")
	}
	if err != sql.ErrNoRows {
		return fmt.Errorf("failed to check existing counter: %w", err)
	}

	err = s.db.QueryRow(
		"INSERT INTO counters (guild_id, word, match_type) VALUES (?, ?, ?) RETURNING id, created_at",
		int64(counter.GuildID), counter.Word, string(counter.MatchType),
	).Scan(&counter.ID, &counter.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create counter: %w", err)
	}
	return nil
}

func (s *SQLiteStore) DeleteCounter(id int64, guildID snowflake.ID) error {
	if _, err := s.db.Exec("DELETE FROM counter_hits WHERE counter_id = ?", id); err != nil {
		return fmt.Errorf("failed to delete counter hits for counter %d: %w", id, err)
	}
	if _, err := s.db.Exec("DELETE FROM counters WHERE id = ? AND guild_id = ?", id, int64(guildID)); err != nil {
		return fmt.Errorf("failed to delete counter %d: %w", id, err)
	}
	return nil
}

func (s *SQLiteStore) GetCounters(guildID snowflake.ID) ([]model.Counter, error) {
	rows, err := s.db.Query(
		"SELECT id, guild_id, word, match_type, created_at FROM counters WHERE guild_id = ? ORDER BY created_at",
		int64(guildID),
	)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	return scanCounters(rows)
}

func (s *SQLiteStore) GetCounter(id int64, guildID snowflake.ID) (*model.Counter, error) {
	var c model.Counter
	var gid int64
	var mt string
	err := s.db.QueryRow(
		"SELECT id, guild_id, word, match_type, created_at FROM counters WHERE id = ? AND guild_id = ?",
		id, int64(guildID),
	).Scan(&c.ID, &gid, &c.Word, &mt, &c.CreatedAt)
	if err != nil {
		return nil, err
	}
	c.GuildID = snowflake.ID(gid)
	c.MatchType = model.MatchType(mt)
	return &c, nil
}

func (s *SQLiteStore) CountCounters(guildID snowflake.ID) (int, error) {
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM counters WHERE guild_id = ?", int64(guildID)).Scan(&count)
	return count, err
}

func (s *SQLiteStore) RecordHits(hits []CounterHit) error {
	if len(hits) == 0 {
		return nil
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	stmt, err := tx.Prepare("INSERT INTO counter_hits (counter_id, guild_id, user_id) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer func() { _ = stmt.Close() }()

	for _, h := range hits {
		if _, err := stmt.Exec(h.CounterID, int64(h.GuildID), int64(h.UserID)); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (s *SQLiteStore) GetCounterStats(guildID snowflake.ID, since *time.Time) ([]model.CounterStat, error) {
	var rows *sql.Rows
	var err error

	if since != nil {
		rows, err = s.db.Query(`
			SELECT c.id, c.word, c.match_type, COUNT(h.id) as hit_count
			FROM counters c
			LEFT JOIN counter_hits h ON h.counter_id = c.id AND h.hit_at >= ?
			WHERE c.guild_id = ?
			GROUP BY c.id
			ORDER BY hit_count DESC`,
			*since, int64(guildID),
		)
	} else {
		rows, err = s.db.Query(`
			SELECT c.id, c.word, c.match_type, COUNT(h.id) as hit_count
			FROM counters c
			LEFT JOIN counter_hits h ON h.counter_id = c.id
			WHERE c.guild_id = ?
			GROUP BY c.id
			ORDER BY hit_count DESC`,
			int64(guildID),
		)
	}
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var stats []model.CounterStat
	for rows.Next() {
		var st model.CounterStat
		var mt string
		if err := rows.Scan(&st.CounterID, &st.Word, &mt, &st.HitCount); err != nil {
			return nil, err
		}
		st.MatchType = model.MatchType(mt)
		stats = append(stats, st)
	}
	return stats, rows.Err()
}

func (s *SQLiteStore) GetCounterUserRanking(counterID int64, since *time.Time, limit int) ([]model.CounterUserRank, error) {
	var rows *sql.Rows
	var err error

	if since != nil {
		rows, err = s.db.Query(`
			SELECT user_id, COUNT(*) as hit_count
			FROM counter_hits
			WHERE counter_id = ? AND hit_at >= ?
			GROUP BY user_id
			ORDER BY hit_count DESC
			LIMIT ?`,
			counterID, *since, limit,
		)
	} else {
		rows, err = s.db.Query(`
			SELECT user_id, COUNT(*) as hit_count
			FROM counter_hits
			WHERE counter_id = ?
			GROUP BY user_id
			ORDER BY hit_count DESC
			LIMIT ?`,
			counterID, limit,
		)
	}
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var ranks []model.CounterUserRank
	for rows.Next() {
		var r model.CounterUserRank
		var uid int64
		if err := rows.Scan(&uid, &r.HitCount); err != nil {
			return nil, err
		}
		r.UserID = snowflake.ID(uid)
		ranks = append(ranks, r)
	}
	return ranks, rows.Err()
}

func scanCounters(rows *sql.Rows) ([]model.Counter, error) {
	var counters []model.Counter
	for rows.Next() {
		var c model.Counter
		var gid int64
		var mt string
		if err := rows.Scan(&c.ID, &gid, &c.Word, &mt, &c.CreatedAt); err != nil {
			return nil, err
		}
		c.GuildID = snowflake.ID(gid)
		c.MatchType = model.MatchType(mt)
		counters = append(counters, c)
	}
	return counters, rows.Err()
}
