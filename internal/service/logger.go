// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package service

import (
	"github.com/disgoorg/snowflake/v2"
	"github.com/s12kuma01/pedmin/internal/model"
	"github.com/s12kuma01/pedmin/internal/repository"
)

// LoggerService handles logger settings persistence.
type LoggerService struct {
	store repository.GuildStore
}

// NewLoggerService creates a new LoggerService.
func NewLoggerService(store repository.GuildStore) *LoggerService {
	return &LoggerService{store: store}
}

// LoadSettings loads the logger settings for a guild.
func (s *LoggerService) LoadSettings(guildID snowflake.ID) (*model.LoggerSettings, error) {
	settings, err := repository.LoadModuleSettings(s.store, guildID, model.LoggerModuleID, func() *model.LoggerSettings {
		return &model.LoggerSettings{Events: make(map[string]bool)}
	})
	if err != nil {
		return nil, err
	}
	if settings.Events == nil {
		settings.Events = make(map[string]bool)
	}
	return settings, nil
}

// SaveSettings saves the logger settings for a guild.
func (s *LoggerService) SaveSettings(guildID snowflake.ID, settings *model.LoggerSettings) error {
	return repository.SaveModuleSettings(s.store, guildID, model.LoggerModuleID, settings)
}
