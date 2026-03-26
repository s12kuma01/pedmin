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

func (s *SQLiteStore) GetUserXP(guildID, userID snowflake.ID) (*model.UserXP, error) {
	var ux model.UserXP
	var gid, uid int64
	var lastXPAt sql.NullTime
	err := s.db.QueryRow(
		"SELECT guild_id, user_id, total_xp, level, message_count, voice_minutes, last_xp_at FROM user_xp WHERE guild_id = ? AND user_id = ?",
		int64(guildID), int64(userID),
	).Scan(&gid, &uid, &ux.TotalXP, &ux.Level, &ux.MessageCount, &ux.VoiceMinutes, &lastXPAt)
	if err == sql.ErrNoRows {
		return &model.UserXP{GuildID: guildID, UserID: userID}, nil
	}
	if err != nil {
		return nil, err
	}
	ux.GuildID = snowflake.ID(gid)
	ux.UserID = snowflake.ID(uid)
	if lastXPAt.Valid {
		ux.LastXPAt = &lastXPAt.Time
	}
	return &ux, nil
}

// AddXP adds XP to a user and returns the updated UserXP and the old level.
func (s *SQLiteStore) AddXP(guildID, userID snowflake.ID, amount int, isVoice bool) (*model.UserXP, int, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, 0, err
	}
	defer func() { _ = tx.Rollback() }()

	var oldTotalXP, oldMsgCount, oldVoiceMin int
	var lastXPAt sql.NullTime
	err = tx.QueryRow(
		"SELECT total_xp, message_count, voice_minutes, last_xp_at FROM user_xp WHERE guild_id = ? AND user_id = ?",
		int64(guildID), int64(userID),
	).Scan(&oldTotalXP, &oldMsgCount, &oldVoiceMin, &lastXPAt)
	if err != nil && err != sql.ErrNoRows {
		return nil, 0, err
	}

	oldLevel, _ := model.LevelFromTotalXP(oldTotalXP)
	newTotalXP := oldTotalXP + amount
	newLevel, _ := model.LevelFromTotalXP(newTotalXP)

	newMsgCount := oldMsgCount
	newVoiceMin := oldVoiceMin
	if isVoice {
		newVoiceMin++
	} else {
		newMsgCount++
	}

	now := time.Now()
	_, err = tx.Exec(`
		INSERT INTO user_xp (guild_id, user_id, total_xp, level, message_count, voice_minutes, last_xp_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(guild_id, user_id) DO UPDATE SET
			total_xp = ?, level = ?, message_count = ?, voice_minutes = ?, last_xp_at = ?`,
		int64(guildID), int64(userID), newTotalXP, newLevel, newMsgCount, newVoiceMin, now,
		newTotalXP, newLevel, newMsgCount, newVoiceMin, now,
	)
	if err != nil {
		return nil, 0, err
	}

	if err := tx.Commit(); err != nil {
		return nil, 0, err
	}

	ux := &model.UserXP{
		GuildID:      guildID,
		UserID:       userID,
		TotalXP:      newTotalXP,
		Level:        newLevel,
		MessageCount: newMsgCount,
		VoiceMinutes: newVoiceMin,
		LastXPAt:     &now,
	}
	return ux, oldLevel, nil
}

func (s *SQLiteStore) GetLeaderboard(guildID snowflake.ID, limit, offset int) ([]model.LeaderboardEntry, error) {
	rows, err := s.db.Query(`
		SELECT user_id, level, total_xp
		FROM user_xp
		WHERE guild_id = ? AND total_xp > 0
		ORDER BY total_xp DESC
		LIMIT ? OFFSET ?`,
		int64(guildID), limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var entries []model.LeaderboardEntry
	rank := offset + 1
	for rows.Next() {
		var e model.LeaderboardEntry
		var uid int64
		if err := rows.Scan(&uid, &e.Level, &e.TotalXP); err != nil {
			return nil, err
		}
		e.UserID = snowflake.ID(uid)
		e.Rank = rank
		rank++
		entries = append(entries, e)
	}
	return entries, rows.Err()
}

func (s *SQLiteStore) GetUserRank(guildID, userID snowflake.ID) (int, error) {
	var rank int
	err := s.db.QueryRow(`
		SELECT COUNT(*) + 1 FROM user_xp
		WHERE guild_id = ? AND total_xp > (
			SELECT COALESCE(total_xp, 0) FROM user_xp WHERE guild_id = ? AND user_id = ?
		)`,
		int64(guildID), int64(guildID), int64(userID),
	).Scan(&rank)
	if err != nil {
		return 0, err
	}
	return rank, nil
}

func (s *SQLiteStore) GetRoleRewards(guildID snowflake.ID) ([]model.LevelRoleReward, error) {
	rows, err := s.db.Query(
		"SELECT id, guild_id, level, role_id FROM level_role_rewards WHERE guild_id = ? ORDER BY level",
		int64(guildID),
	)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var rewards []model.LevelRoleReward
	for rows.Next() {
		var r model.LevelRoleReward
		var gid, rid int64
		if err := rows.Scan(&r.ID, &gid, &r.Level, &rid); err != nil {
			return nil, err
		}
		r.GuildID = snowflake.ID(gid)
		r.RoleID = snowflake.ID(rid)
		rewards = append(rewards, r)
	}
	return rewards, rows.Err()
}

func (s *SQLiteStore) AddRoleReward(guildID snowflake.ID, level int, roleID snowflake.ID) error {
	count, err := s.CountRoleRewards(guildID)
	if err != nil {
		return fmt.Errorf("failed to check reward count: %w", err)
	}
	if count >= model.MaxRoleRewardsPerGuild {
		return fmt.Errorf("報酬数が上限(%d)に達しています", model.MaxRoleRewardsPerGuild)
	}

	_, err = s.db.Exec(
		"INSERT OR IGNORE INTO level_role_rewards (guild_id, level, role_id) VALUES (?, ?, ?)",
		int64(guildID), level, int64(roleID),
	)
	if err != nil {
		return fmt.Errorf("failed to add role reward: %w", err)
	}
	return nil
}

func (s *SQLiteStore) RemoveRoleReward(id int64, guildID snowflake.ID) error {
	_, err := s.db.Exec("DELETE FROM level_role_rewards WHERE id = ? AND guild_id = ?", id, int64(guildID))
	if err != nil {
		return fmt.Errorf("failed to remove role reward: %w", err)
	}
	return nil
}

func (s *SQLiteStore) CountRoleRewards(guildID snowflake.ID) (int, error) {
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM level_role_rewards WHERE guild_id = ?", int64(guildID)).Scan(&count)
	return count, err
}

func (s *SQLiteStore) BatchAddVoiceXP(updates []VoiceXPUpdate) error {
	if len(updates) == 0 {
		return nil
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	stmt, err := tx.Prepare(`
		INSERT INTO user_xp (guild_id, user_id, total_xp, level, message_count, voice_minutes, last_xp_at)
		VALUES (?, ?, ?, ?, 0, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(guild_id, user_id) DO UPDATE SET
			total_xp = total_xp + ?,
			voice_minutes = voice_minutes + ?,
			last_xp_at = CURRENT_TIMESTAMP`)
	if err != nil {
		return err
	}
	defer func() { _ = stmt.Close() }()

	for _, u := range updates {
		newLevel, _ := model.LevelFromTotalXP(u.XPAmount)
		if _, err := stmt.Exec(
			int64(u.GuildID), int64(u.UserID), u.XPAmount, newLevel, u.Minutes,
			u.XPAmount, u.Minutes,
		); err != nil {
			return err
		}
	}
	return tx.Commit()
}
