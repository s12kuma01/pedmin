// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package service

import (
	"log/slog"

	disgobot "github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/snowflake/v2"

	"github.com/Sumire-Labs/pedmin/internal/model"
	"github.com/Sumire-Labs/pedmin/internal/repository"
)

// AutoroleService handles autorole business logic.
type AutoroleService struct {
	store  repository.GuildStore
	client *disgobot.Client
	logger *slog.Logger
}

// NewAutoroleService creates a new AutoroleService.
func NewAutoroleService(client *disgobot.Client, store repository.GuildStore, logger *slog.Logger) *AutoroleService {
	return &AutoroleService{
		store:  store,
		client: client,
		logger: logger,
	}
}

// LoadSettings loads autorole settings for a guild.
func (s *AutoroleService) LoadSettings(guildID snowflake.ID) (*model.AutoroleSettings, error) {
	return repository.LoadModuleSettings(s.store, guildID, model.AutoroleModuleID, model.DefaultAutoroleSettings)
}

// SaveSettings saves autorole settings for a guild.
func (s *AutoroleService) SaveSettings(guildID snowflake.ID, settings *model.AutoroleSettings) error {
	return repository.SaveModuleSettings(s.store, guildID, model.AutoroleModuleID, settings)
}

// AssignRole assigns the appropriate role to a member based on whether they are a bot.
func (s *AutoroleService) AssignRole(guildID, userID snowflake.ID, isBot bool) {
	settings, err := s.LoadSettings(guildID)
	if err != nil {
		s.logger.Error("failed to load autorole settings", slog.Any("error", err))
		return
	}

	var roleID snowflake.ID
	if isBot {
		roleID = settings.BotRoleID
	} else {
		roleID = settings.UserRoleID
	}

	if roleID == 0 {
		return
	}

	if err := s.client.Rest.AddMemberRole(guildID, userID, roleID); err != nil {
		s.logger.Error("failed to assign autorole",
			slog.Bool("is_bot", isBot),
			slog.Any("error", err),
		)
	}
}
