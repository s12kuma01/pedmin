// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package handler

import (
	"fmt"
	"log/slog"

	disgobot "github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
	"github.com/Sumire-Labs/pedmin/internal/model"
	"github.com/Sumire-Labs/pedmin/internal/module"
	"github.com/Sumire-Labs/pedmin/internal/service"
	"github.com/Sumire-Labs/pedmin/internal/view"
)

// CounterBot is the interface the handler needs from the bot registry.
type CounterBot interface {
	IsModuleEnabled(guildID snowflake.ID, moduleID string) bool
}

// CounterHandler implements module.Module for the counter feature.
type CounterHandler struct {
	bot     CounterBot
	service *service.CounterService
	logger  *slog.Logger
}

// NewCounterHandler creates a new CounterHandler.
func NewCounterHandler(bot CounterBot, svc *service.CounterService, logger *slog.Logger) *CounterHandler {
	return &CounterHandler{
		bot:     bot,
		service: svc,
		logger:  logger,
	}
}

// SetupCounterListeners registers the message listener on the Discord client.
func SetupCounterListeners(client *disgobot.Client, h *CounterHandler) {
	client.AddEventListeners(
		disgobot.NewListenerFunc(h.onMessageCreate),
	)
}

func (h *CounterHandler) Info() module.Info {
	return module.Info{
		ID:          model.CounterModuleID,
		Name:        "Counter",
		Description: "ワードカウンター",
		AlwaysOn:    false,
	}
}

func (h *CounterHandler) Commands() []discord.ApplicationCommandCreate {
	return nil
}

func (h *CounterHandler) HandleCommand(_ *events.ApplicationCommandInteractionCreate) {}

func (h *CounterHandler) SettingsSummary(guildID snowflake.ID) string {
	count, err := h.service.CountCounters(guildID)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("カウンター: %d/%d", count, model.MaxCountersPerGuild)
}

func (h *CounterHandler) SettingsPanel(guildID snowflake.ID) []discord.LayoutComponent {
	count, err := h.service.CountCounters(guildID)
	if err != nil {
		h.logger.Error("failed to count counters", slog.Any("error", err))
	}
	return view.CounterSettingsPanel(count)
}

func (h *CounterHandler) onMessageCreate(e *events.GuildMessageCreate) {
	if e.Message.Author.Bot {
		return
	}
	if !h.bot.IsModuleEnabled(e.GuildID, model.CounterModuleID) {
		return
	}

	h.service.ProcessMessage(e.GuildID, e.Message.Author.ID, e.Message.Content)
}
