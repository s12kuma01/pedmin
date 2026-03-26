// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package handler

import (
	"log/slog"
	"strings"

	disgobot "github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
	"github.com/Sumire-Labs/pedmin/internal/model"
	"github.com/Sumire-Labs/pedmin/internal/module"
	"github.com/Sumire-Labs/pedmin/internal/service"
	"github.com/Sumire-Labs/pedmin/internal/view"
)

// EmbedFixBot is the interface the handler needs from the bot registry.
type EmbedFixBot interface {
	IsModuleEnabled(guildID snowflake.ID, moduleID string) bool
}

// EmbedFixHandler implements module.Module for the embedfix feature.
type EmbedFixHandler struct {
	bot     EmbedFixBot
	service *service.EmbedFixService
	logger  *slog.Logger
}

// NewEmbedFixHandler creates a new EmbedFixHandler.
func NewEmbedFixHandler(bot EmbedFixBot, svc *service.EmbedFixService, logger *slog.Logger) *EmbedFixHandler {
	return &EmbedFixHandler{
		bot:     bot,
		service: svc,
		logger:  logger,
	}
}

// SetupEmbedFixListeners registers the message listener on the Discord client.
func SetupEmbedFixListeners(client *disgobot.Client, h *EmbedFixHandler) {
	client.AddEventListeners(
		disgobot.NewListenerFunc(h.onMessageCreate),
	)
}

func (h *EmbedFixHandler) Info() module.Info {
	return module.Info{
		ID:          model.EmbedFixModuleID,
		Name:        "Embed Fix",
		Description: "SNSリンクのリッチ埋め込み表示",
		AlwaysOn:    false,
	}
}

func (h *EmbedFixHandler) Commands() []discord.ApplicationCommandCreate {
	return nil
}

func (h *EmbedFixHandler) HandleCommand(_ *events.ApplicationCommandInteractionCreate) {}

func (h *EmbedFixHandler) HandleComponent(e *events.ComponentInteractionCreate) {
	h.handleComponent(e)
}

func (h *EmbedFixHandler) HandleModal(_ *events.ModalSubmitInteractionCreate) {}

func (h *EmbedFixHandler) SettingsSummary(guildID snowflake.ID) string {
	settings, err := h.service.LoadSettings(guildID)
	if err != nil {
		return ""
	}
	var names []string
	for _, p := range model.AllPlatforms {
		if settings.IsPlatformEnabled(p.Key) {
			names = append(names, p.Label)
		}
	}
	if len(names) == 0 {
		return "全て無効"
	}
	return "対象: " + strings.Join(names, ", ")
}

func (h *EmbedFixHandler) SettingsPanel(guildID snowflake.ID) []discord.LayoutComponent {
	settings, err := h.service.LoadSettings(guildID)
	if err != nil {
		h.logger.Error("failed to load embedfix settings", slog.Any("error", err))
		settings = model.DefaultEmbedFixSettings()
	}
	return view.BuildEmbedFixSettingsPanel(settings)
}
