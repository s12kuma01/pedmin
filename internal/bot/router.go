// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package bot

import (
	"log/slog"
	"strings"

	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
	"github.com/s12kuma01/pedmin/internal/ui"
)

func (b *Bot) isModuleDisabledForGuild(guildID *snowflake.ID, moduleID string) bool {
	return guildID != nil && !b.IsModuleEnabled(*guildID, moduleID)
}

func (b *Bot) onCommandInteraction(e *events.ApplicationCommandInteractionCreate) {
	cmdName := e.SlashCommandInteractionData().CommandName()

	for _, m := range b.modules {
		for _, cmd := range m.Commands() {
			if cmd.CommandName() == cmdName {
				if b.isModuleDisabledForGuild(e.GuildID(), m.Info().ID) {
					_ = e.CreateMessage(ui.ErrorMessage("このモジュールは現在無効です。"))
					return
				}
				m.HandleCommand(e)
				return
			}
		}
	}
	b.Logger.Warn("unhandled command", slog.String("command", cmdName))
}

func (b *Bot) onComponentInteraction(e *events.ComponentInteractionCreate) {
	customID := e.Data.CustomID()
	moduleID, _, _ := strings.Cut(customID, ":")

	m, ok := b.modules[moduleID]
	if !ok {
		b.Logger.Warn("unhandled component", slog.String("custom_id", customID))
		return
	}

	if b.isModuleDisabledForGuild(e.GuildID(), m.Info().ID) {
		_ = e.CreateMessage(ui.ErrorMessage("このモジュールは現在無効です。"))
		return
	}

	m.HandleComponent(e)
}

func (b *Bot) onModalSubmit(e *events.ModalSubmitInteractionCreate) {
	customID := e.Data.CustomID
	moduleID, _, _ := strings.Cut(customID, ":")

	m, ok := b.modules[moduleID]
	if !ok {
		b.Logger.Warn("unhandled modal", slog.String("custom_id", customID))
		return
	}

	if b.isModuleDisabledForGuild(e.GuildID(), m.Info().ID) {
		_ = e.CreateMessage(ui.ErrorMessage("このモジュールは現在無効です。"))
		return
	}

	m.HandleModal(e)
}
