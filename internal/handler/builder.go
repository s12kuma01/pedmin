// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package handler

import (
	"log/slog"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
	"github.com/Sumire-Labs/pedmin/internal/model"
	"github.com/Sumire-Labs/pedmin/internal/module"
	"github.com/Sumire-Labs/pedmin/internal/service"
)

// BuilderHandler implements module.Module for the component builder.
type BuilderHandler struct {
	service *service.BuilderService
	logger  *slog.Logger
}

// NewBuilderHandler creates a new BuilderHandler.
func NewBuilderHandler(svc *service.BuilderService, logger *slog.Logger) *BuilderHandler {
	return &BuilderHandler{
		service: svc,
		logger:  logger,
	}
}

func (h *BuilderHandler) Info() module.Info {
	return module.Info{
		ID:          model.BuilderModuleID,
		Name:        "Component Builder",
		Description: "Components V2 パネルビルダー",
		AlwaysOn:    true,
	}
}

func (h *BuilderHandler) Commands() []discord.ApplicationCommandCreate {
	return nil
}

func (h *BuilderHandler) HandleCommand(_ *events.ApplicationCommandInteractionCreate) {}

func (h *BuilderHandler) SettingsPanel(_ snowflake.ID) []discord.LayoutComponent {
	return nil
}
