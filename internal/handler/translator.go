// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package handler

import (
	"log/slog"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
	"github.com/s12kuma01/pedmin/internal/model"
	"github.com/s12kuma01/pedmin/internal/module"
	"github.com/s12kuma01/pedmin/internal/service"
	"github.com/s12kuma01/pedmin/pkg/deepl"

	disgobot "github.com/disgoorg/disgo/bot"
)

// TranslatorBot defines the bot interface needed by the translator handler.
type TranslatorBot interface {
	IsModuleEnabled(guildID snowflake.ID, moduleID string) bool
}

// TranslatorHandler implements module.Module for the translator feature.
type TranslatorHandler struct {
	bot     TranslatorBot
	service *service.TranslatorService
	logger  *slog.Logger
}

// NewTranslatorHandler creates a new TranslatorHandler.
func NewTranslatorHandler(bot TranslatorBot, client *disgobot.Client, deeplAPIKey string, timeout time.Duration, logger *slog.Logger) *TranslatorHandler {
	deeplClient := deepl.NewTranslateClient(deeplAPIKey, timeout)
	svc := service.NewTranslatorService(client.Rest, deeplClient, logger)
	return &TranslatorHandler{
		bot:     bot,
		service: svc,
		logger:  logger,
	}
}

func (h *TranslatorHandler) Info() module.Info {
	return module.Info{
		ID:          model.TranslatorModuleID,
		Name:        "Translator",
		Description: "国旗リアクションでメッセージを翻訳",
		AlwaysOn:    false,
	}
}

func (h *TranslatorHandler) Commands() []discord.ApplicationCommandCreate {
	return nil
}

func (h *TranslatorHandler) HandleCommand(_ *events.ApplicationCommandInteractionCreate) {}
func (h *TranslatorHandler) HandleComponent(_ *events.ComponentInteractionCreate)        {}
func (h *TranslatorHandler) HandleModal(_ *events.ModalSubmitInteractionCreate)           {}

func (h *TranslatorHandler) SettingsPanel(_ snowflake.ID) []discord.LayoutComponent {
	return nil
}
