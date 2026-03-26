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

func (s *SQLiteStore) CreateTicket(guildID snowflake.ID, number int, channelID, userID snowflake.ID, subject string) error {
	_, err := s.db.Exec(
		"INSERT INTO tickets (guild_id, number, channel_id, user_id, subject) VALUES (?, ?, ?, ?, ?)",
		int64(guildID), number, int64(channelID), int64(userID), subject,
	)
	if err != nil {
		return fmt.Errorf("failed to create ticket #%d: %w", number, err)
	}
	return nil
}

func (s *SQLiteStore) GetTicketByChannel(channelID snowflake.ID) (*model.Ticket, error) {
	var t model.Ticket
	var guildID, chID, userID int64
	var closedAt sql.NullTime
	var closedBy sql.NullInt64

	err := s.db.QueryRow(
		"SELECT guild_id, number, channel_id, user_id, subject, created_at, closed_at, closed_by FROM tickets WHERE channel_id = ?",
		int64(channelID),
	).Scan(&guildID, &t.Number, &chID, &userID, &t.Subject, &t.CreatedAt, &closedAt, &closedBy)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	t.GuildID = snowflake.ID(guildID)
	t.ChannelID = snowflake.ID(chID)
	t.UserID = snowflake.ID(userID)
	if closedAt.Valid {
		t.ClosedAt = &closedAt.Time
	}
	if closedBy.Valid {
		id := snowflake.ID(closedBy.Int64)
		t.ClosedBy = &id
	}
	return &t, nil
}

func (s *SQLiteStore) CloseTicket(channelID snowflake.ID, closedBy snowflake.ID) error {
	_, err := s.db.Exec(
		"UPDATE tickets SET closed_at = ?, closed_by = ? WHERE channel_id = ?",
		time.Now().UTC(), int64(closedBy), int64(channelID),
	)
	if err != nil {
		return fmt.Errorf("failed to close ticket in channel %d: %w", channelID, err)
	}
	return nil
}

func (s *SQLiteStore) DeleteTicket(channelID snowflake.ID) error {
	_, err := s.db.Exec("DELETE FROM tickets WHERE channel_id = ?", int64(channelID))
	if err != nil {
		return fmt.Errorf("failed to delete ticket in channel %d: %w", channelID, err)
	}
	return nil
}
