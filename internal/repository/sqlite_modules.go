// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/disgoorg/snowflake/v2"

	"github.com/Sumire-Labs/pedmin/internal/model"
)

func (s *SQLiteStore) Get(guildID snowflake.ID) (*model.GuildSettings, error) {
	settings := &model.GuildSettings{
		GuildID:        guildID,
		EnabledModules: make(map[string]bool),
		ModuleSettings: make(map[string]any),
	}

	rows, err := s.db.Query("SELECT module_id, enabled FROM guild_modules WHERE guild_id = ?", int64(guildID))
	if err != nil {
		return nil, fmt.Errorf("failed to query guild modules: %w", err)
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var moduleID string
		var enabled bool
		if err := rows.Scan(&moduleID, &enabled); err != nil {
			return nil, fmt.Errorf("failed to scan guild module: %w", err)
		}
		settings.EnabledModules[moduleID] = enabled
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate guild modules: %w", err)
	}

	settingsRows, err := s.db.Query("SELECT module_id, settings FROM guild_module_settings WHERE guild_id = ?", int64(guildID))
	if err != nil {
		return nil, fmt.Errorf("failed to query guild module settings: %w", err)
	}
	defer func() { _ = settingsRows.Close() }()

	for settingsRows.Next() {
		var moduleID string
		var settingsJSON string
		if err := settingsRows.Scan(&moduleID, &settingsJSON); err != nil {
			return nil, fmt.Errorf("failed to scan guild module settings: %w", err)
		}
		var parsed any
		if err := json.Unmarshal([]byte(settingsJSON), &parsed); err != nil {
			return nil, fmt.Errorf("failed to parse module settings for %s: %w", moduleID, err)
		}
		settings.ModuleSettings[moduleID] = parsed
	}
	if err := settingsRows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate guild module settings: %w", err)
	}

	return settings, nil
}

func (s *SQLiteStore) Save(settings *model.GuildSettings) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	gid := int64(settings.GuildID)

	for moduleID, enabled := range settings.EnabledModules {
		_, err := tx.Exec(
			"INSERT INTO guild_modules (guild_id, module_id, enabled) VALUES (?, ?, ?) ON CONFLICT(guild_id, module_id) DO UPDATE SET enabled = excluded.enabled",
			gid, moduleID, enabled,
		)
		if err != nil {
			return fmt.Errorf("failed to upsert module %s: %w", moduleID, err)
		}
	}

	for moduleID, modSettings := range settings.ModuleSettings {
		settingsJSON, err := json.Marshal(modSettings)
		if err != nil {
			return fmt.Errorf("failed to marshal settings for %s: %w", moduleID, err)
		}
		_, err = tx.Exec(
			"INSERT INTO guild_module_settings (guild_id, module_id, settings) VALUES (?, ?, ?) ON CONFLICT(guild_id, module_id) DO UPDATE SET settings = excluded.settings",
			gid, moduleID, string(settingsJSON),
		)
		if err != nil {
			return fmt.Errorf("failed to upsert settings for %s: %w", moduleID, err)
		}
	}

	return tx.Commit()
}

func (s *SQLiteStore) IsModuleEnabled(guildID snowflake.ID, moduleID string) (bool, error) {
	var enabled bool
	err := s.db.QueryRow(
		"SELECT enabled FROM guild_modules WHERE guild_id = ? AND module_id = ?",
		int64(guildID), moduleID,
	).Scan(&enabled)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to check module enabled: %w", err)
	}
	return enabled, nil
}

func (s *SQLiteStore) SetModuleEnabled(guildID snowflake.ID, moduleID string, enabled bool) error {
	_, err := s.db.Exec(
		"INSERT INTO guild_modules (guild_id, module_id, enabled) VALUES (?, ?, ?) ON CONFLICT(guild_id, module_id) DO UPDATE SET enabled = excluded.enabled",
		int64(guildID), moduleID, enabled,
	)
	if err != nil {
		return fmt.Errorf("failed to set module enabled: %w", err)
	}
	return nil
}

func (s *SQLiteStore) GetModuleSettings(guildID snowflake.ID, moduleID string) (string, error) {
	var settings string
	err := s.db.QueryRow(
		"SELECT settings FROM guild_module_settings WHERE guild_id = ? AND module_id = ?",
		int64(guildID), moduleID,
	).Scan(&settings)
	if err == sql.ErrNoRows {
		return "{}", nil
	}
	return settings, err
}

func (s *SQLiteStore) SetModuleSettings(guildID snowflake.ID, moduleID string, settings string) error {
	_, err := s.db.Exec(
		`INSERT INTO guild_module_settings (guild_id, module_id, settings) VALUES (?, ?, ?)
		 ON CONFLICT(guild_id, module_id) DO UPDATE SET settings = excluded.settings`,
		int64(guildID), moduleID, settings,
	)
	if err != nil {
		return fmt.Errorf("failed to set module settings for %s: %w", moduleID, err)
	}
	return nil
}
