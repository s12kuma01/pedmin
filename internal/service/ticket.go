// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package service

import (
	"fmt"
	"log/slog"

	disgobot "github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
	"github.com/Sumire-Labs/pedmin/internal/model"
	"github.com/Sumire-Labs/pedmin/internal/repository"
	"github.com/Sumire-Labs/pedmin/internal/view"
)

// TicketService handles ticket creation, archival, and deletion logic.
type TicketService struct {
	client *disgobot.Client
	store  repository.GuildStore
	logger *slog.Logger
}

// NewTicketService creates a new TicketService.
func NewTicketService(client *disgobot.Client, store repository.GuildStore, logger *slog.Logger) *TicketService {
	return &TicketService{
		client: client,
		store:  store,
		logger: logger,
	}
}

// CreateTicket creates a new ticket channel and records it in the store.
func (s *TicketService) CreateTicket(guildID, userID snowflake.ID, subject, description string) (snowflake.ID, int, error) {
	settings, err := s.LoadSettings(guildID)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to load settings: %w", err)
	}

	settings.NextNumber++
	number := settings.NextNumber

	if err := s.SaveSettings(guildID, settings); err != nil {
		return 0, 0, fmt.Errorf("failed to save settings: %w", err)
	}

	channelName := fmt.Sprintf("ticket-%04d", number)

	overwrites := []discord.PermissionOverwrite{
		discord.RolePermissionOverwrite{
			RoleID: guildID,
			Deny:   discord.PermissionViewChannel,
		},
		discord.MemberPermissionOverwrite{
			UserID: userID,
			Allow:  discord.PermissionViewChannel | discord.PermissionSendMessages | discord.PermissionReadMessageHistory,
		},
		discord.MemberPermissionOverwrite{
			UserID: s.client.ApplicationID,
			Allow:  discord.PermissionViewChannel | discord.PermissionSendMessages | discord.PermissionManageChannels,
		},
	}

	if settings.SupportRoleID != 0 {
		overwrites = append(overwrites, discord.RolePermissionOverwrite{
			RoleID: settings.SupportRoleID,
			Allow:  discord.PermissionViewChannel | discord.PermissionSendMessages | discord.PermissionReadMessageHistory,
		})
	}

	create := discord.GuildTextChannelCreate{
		Name:                 channelName,
		PermissionOverwrites: overwrites,
	}
	if settings.CategoryID != 0 {
		create.ParentID = settings.CategoryID
	}

	ch, err := s.client.Rest.CreateGuildChannel(guildID, create)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to create channel: %w", err)
	}

	channelID := ch.ID()

	if err := s.store.CreateTicket(guildID, number, channelID, userID, subject); err != nil {
		return 0, 0, fmt.Errorf("failed to save ticket: %w", err)
	}

	msg := view.TicketInfo(number, userID, subject, description)
	if _, err := s.client.Rest.CreateMessage(channelID, msg); err != nil {
		s.logger.Error("failed to send ticket info", slog.Any("error", err))
	}

	return channelID, number, nil
}

// ArchiveTicket closes a ticket by revoking the creator's access and sending a transcript.
func (s *TicketService) ArchiveTicket(guildID, channelID, closedBy snowflake.ID) (*model.Ticket, error) {
	ticket, err := s.store.GetTicketByChannel(channelID)
	if err != nil || ticket == nil {
		return nil, fmt.Errorf("failed to get ticket: %w", err)
	}

	// Revoke the creator's view permission
	err = s.client.Rest.UpdatePermissionOverwrite(channelID, ticket.UserID, discord.MemberPermissionOverwriteUpdate{
		Deny: ptrPermissions(discord.PermissionViewChannel),
	})
	if err != nil {
		s.logger.Error("failed to update permissions", slog.Any("error", err))
	}

	if err := s.store.CloseTicket(channelID, closedBy); err != nil {
		s.logger.Error("failed to close ticket in db", slog.Any("error", err))
	}

	// Send transcript to log channel
	s.sendTranscriptLog(guildID, ticket)

	return ticket, nil
}

// DeleteTicket sends a log and then deletes the ticket channel.
func (s *TicketService) DeleteTicket(guildID, channelID snowflake.ID) (*model.Ticket, error) {
	ticket, err := s.store.GetTicketByChannel(channelID)
	if err != nil || ticket == nil {
		return nil, fmt.Errorf("failed to get ticket: %w", err)
	}

	s.sendTicketLog(guildID, ticket)

	if err := s.store.DeleteTicket(channelID); err != nil {
		s.logger.Error("failed to delete ticket from db", slog.Any("error", err))
	}

	_ = s.client.Rest.DeleteChannel(channelID)

	return ticket, nil
}

// DeployPanel sends the ticket panel message to the specified channel.
func (s *TicketService) DeployPanel(channelID snowflake.ID) error {
	panel := view.TicketPanel()
	_, err := s.client.Rest.CreateMessage(channelID, panel)
	return err
}

func ptrPermissions(p discord.Permissions) *discord.Permissions {
	return &p
}
