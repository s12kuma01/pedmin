// opyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

// Package module defines the Module interface that all feature modules implement.
package module

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
)

type Info struct {
	ID          string
	Name        string
	Description string
	AlwaysOn    bool
}

type Module interface {
	Info() Info
	Commands() []discord.ApplicationCommandCreate
	HandleCommand(e *events.ApplicationCommandInteractionCreate)
	HandleComponent(e *events.ComponentInteractionCreate)
	HandleModal(e *events.ModalSubmitInteractionCreate)
	SettingsPanel(guildID snowflake.ID) []discord.LayoutComponent
}

// SettingsSummarizer is optionally implemented by modules to show
// a brief summary of their current settings in the main settings panel.
type SettingsSummarizer interface {
	SettingsSummary(guildID snowflake.ID) string
}

// VoiceStateListener is an optional interface that modules can implement
// to receive voice state updates for non-bot users.
type VoiceStateListener interface {
	OnVoiceStateUpdate(guildID, channelID, userID snowflake.ID)
}
