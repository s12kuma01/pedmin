// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package handler

import (
	"log/slog"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
	"github.com/s12kuma01/pedmin/internal/model"
	"github.com/s12kuma01/pedmin/internal/module"
	"github.com/s12kuma01/pedmin/internal/view"
)

// AvatarHandler implements module.Module for the avatar feature.
type AvatarHandler struct {
	logger *slog.Logger
}

// NewAvatarHandler creates a new AvatarHandler.
func NewAvatarHandler(logger *slog.Logger) *AvatarHandler {
	return &AvatarHandler{logger: logger}
}

func (h *AvatarHandler) Info() module.Info {
	return module.Info{
		ID:          model.AvatarModuleID,
		Name:        "アバター",
		Description: "ユーザーのアバターを表示",
		AlwaysOn:    true,
	}
}

func (h *AvatarHandler) Commands() []discord.ApplicationCommandCreate {
	return []discord.ApplicationCommandCreate{
		discord.SlashCommandCreate{
			Name:        "avatar",
			Description: "ユーザーのアバターを表示する",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionUser{
					Name:        "user",
					Description: "アバターを表示するユーザー（省略時は自分）",
					Required:    false,
				},
			},
		},
	}
}

func (h *AvatarHandler) HandleCommand(e *events.ApplicationCommandInteractionCreate) {
	data := e.SlashCommandInteractionData()

	var user discord.User
	var member *discord.ResolvedMember

	if optUser, ok := data.OptUser("user"); ok {
		user = optUser
		if m, ok := data.OptMember("user"); ok {
			member = &m
		}
	} else {
		user = e.User()
		if m := e.Member(); m != nil {
			member = m
		}
	}

	guildID := e.GuildID()
	ui := view.BuildAvatarGallery(user, member, guildID)

	_ = e.CreateMessage(discord.NewMessageCreateV2(ui))
}

func (h *AvatarHandler) HandleComponent(_ *events.ComponentInteractionCreate) {}
func (h *AvatarHandler) HandleModal(_ *events.ModalSubmitInteractionCreate)   {}
func (h *AvatarHandler) SettingsPanel(_ snowflake.ID) []discord.LayoutComponent {
	return nil
}
