package ticket

import (
	"fmt"
	"log/slog"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
)

func (t *Ticket) createTicket(guildID, userID snowflake.ID, subject, description string) (snowflake.ID, int, error) {
	settings, err := LoadSettings(t.store, guildID)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to load settings: %w", err)
	}

	settings.NextNumber++
	number := settings.NextNumber

	if err := SaveSettings(t.store, guildID, settings); err != nil {
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
			UserID: t.client.ApplicationID,
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

	ch, err := t.client.Rest.CreateGuildChannel(guildID, create)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to create channel: %w", err)
	}

	channelID := ch.ID()

	if err := t.store.CreateTicket(guildID, number, channelID, userID, subject); err != nil {
		return 0, 0, fmt.Errorf("failed to save ticket: %w", err)
	}

	msg := BuildTicketInfo(number, userID, subject, description)
	if _, err := t.client.Rest.CreateMessage(channelID, msg); err != nil {
		t.logger.Error("failed to send ticket info", slog.Any("error", err))
	}

	return channelID, number, nil
}

func (t *Ticket) archiveTicket(e *events.ComponentInteractionCreate, guildID snowflake.ID) {
	channelID := e.Channel().ID()

	ticket, err := t.store.GetTicketByChannel(channelID)
	if err != nil || ticket == nil {
		t.logger.Error("failed to get ticket", slog.Any("error", err))
		return
	}

	_ = e.DeferUpdateMessage()

	// Revoke the creator's view permission
	err = t.client.Rest.UpdatePermissionOverwrite(channelID, ticket.UserID, discord.MemberPermissionOverwriteUpdate{
		Deny: ptrPermissions(discord.PermissionViewChannel),
	})
	if err != nil {
		t.logger.Error("failed to update permissions", slog.Any("error", err))
	}

	if err := t.store.CloseTicket(channelID, e.User().ID); err != nil {
		t.logger.Error("failed to close ticket in db", slog.Any("error", err))
	}

	archiveUI := BuildArchiveInfo(ticket.Number, ticket.UserID, ticket.Subject)
	_, _ = e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(), discord.NewMessageUpdateV2(archiveUI.Components))

	// Send transcript to log channel
	t.sendTranscriptLog(guildID, ticket)
}

func (t *Ticket) deleteTicket(e *events.ComponentInteractionCreate, guildID snowflake.ID) {
	channelID := e.Channel().ID()

	ticket, err := t.store.GetTicketByChannel(channelID)
	if err != nil || ticket == nil {
		t.logger.Error("failed to get ticket", slog.Any("error", err))
		return
	}

	_ = e.DeferUpdateMessage()

	t.sendTicketLog(guildID, ticket)

	if err := t.store.DeleteTicket(channelID); err != nil {
		t.logger.Error("failed to delete ticket from db", slog.Any("error", err))
	}

	_ = t.client.Rest.DeleteChannel(channelID)
}

func ptrPermissions(p discord.Permissions) *discord.Permissions {
	return &p
}
