// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/disgoorg/snowflake/v2"

	"github.com/s12kuma01/pedmin/internal/model"
)

func (s *SQLiteStore) CreatePanel(panel *model.ComponentPanel) error {
	count, err := s.CountPanels(panel.GuildID)
	if err != nil {
		return fmt.Errorf("failed to check panel count: %w", err)
	}
	if count >= model.MaxPanelsPerGuild {
		return fmt.Errorf("パネル数が上限(%d)に達しています", model.MaxPanelsPerGuild)
	}

	compsJSON, err := json.Marshal(panel.Components)
	if err != nil {
		return fmt.Errorf("failed to marshal components: %w", err)
	}

	err = s.db.QueryRow(
		"INSERT INTO component_panels (guild_id, name, components) VALUES (?, ?, ?) RETURNING id, created_at",
		int64(panel.GuildID), panel.Name, string(compsJSON),
	).Scan(&panel.ID, &panel.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create panel: %w", err)
	}
	return nil
}

func (s *SQLiteStore) UpdatePanel(panel *model.ComponentPanel) error {
	compsJSON, err := json.Marshal(panel.Components)
	if err != nil {
		return fmt.Errorf("failed to marshal components: %w", err)
	}

	_, err = s.db.Exec(
		"UPDATE component_panels SET name = ?, components = ? WHERE id = ? AND guild_id = ?",
		panel.Name, string(compsJSON), panel.ID, int64(panel.GuildID),
	)
	if err != nil {
		return fmt.Errorf("failed to update panel: %w", err)
	}
	return nil
}

func (s *SQLiteStore) DeletePanel(id int64, guildID snowflake.ID) error {
	_, err := s.db.Exec("DELETE FROM component_panels WHERE id = ? AND guild_id = ?", id, int64(guildID))
	if err != nil {
		return fmt.Errorf("failed to delete panel: %w", err)
	}
	return nil
}

func (s *SQLiteStore) GetPanel(id int64, guildID snowflake.ID) (*model.ComponentPanel, error) {
	var p model.ComponentPanel
	var gid int64
	var compsJSON string
	err := s.db.QueryRow(
		"SELECT id, guild_id, name, components, created_at FROM component_panels WHERE id = ? AND guild_id = ?",
		id, int64(guildID),
	).Scan(&p.ID, &gid, &p.Name, &compsJSON, &p.CreatedAt)
	if err != nil {
		return nil, err
	}
	p.GuildID = snowflake.ID(gid)
	if err := json.Unmarshal([]byte(compsJSON), &p.Components); err != nil {
		p.Components = nil
	}
	return &p, nil
}

func (s *SQLiteStore) GetPanels(guildID snowflake.ID) ([]model.ComponentPanel, error) {
	rows, err := s.db.Query(
		"SELECT id, guild_id, name, components, created_at FROM component_panels WHERE guild_id = ? ORDER BY created_at",
		int64(guildID),
	)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	return scanPanels(rows)
}

func (s *SQLiteStore) CountPanels(guildID snowflake.ID) (int, error) {
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM component_panels WHERE guild_id = ?", int64(guildID)).Scan(&count)
	return count, err
}

func scanPanels(rows *sql.Rows) ([]model.ComponentPanel, error) {
	var panels []model.ComponentPanel
	for rows.Next() {
		var p model.ComponentPanel
		var gid int64
		var compsJSON string
		if err := rows.Scan(&p.ID, &gid, &p.Name, &compsJSON, &p.CreatedAt); err != nil {
			return nil, err
		}
		p.GuildID = snowflake.ID(gid)
		if err := json.Unmarshal([]byte(compsJSON), &p.Components); err != nil {
			p.Components = nil
		}
		panels = append(panels, p)
	}
	return panels, rows.Err()
}
