// opyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package handler

import (
	"log/slog"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
	"github.com/s12kuma01/pedmin/config"
	"github.com/s12kuma01/pedmin/internal/client"
	"github.com/s12kuma01/pedmin/internal/model"
	"github.com/s12kuma01/pedmin/internal/module"
	"github.com/s12kuma01/pedmin/internal/service"
	"github.com/s12kuma01/pedmin/internal/view"
)

// URLHandler implements module.Module for the URL tools feature.
type URLHandler struct {
	cfg     *config.Config
	service *service.URLService
	logger  *slog.Logger
}

// NewURLHandler creates a new URLHandler.
func NewURLHandler(cfg *config.Config, logger *slog.Logger) *URLHandler {
	urlClient := client.NewURLClient(cfg.XGDAPIKey, cfg.VTAPIKey, config.DefaultHTTPClientTimeout)
	return &URLHandler{
		cfg:     cfg,
		service: service.NewURLService(urlClient),
		logger:  logger,
	}
}

func (h *URLHandler) Info() module.Info {
	return module.Info{
		ID:          model.URLModuleID,
		Name:        "URL Tools",
		Description: "URL短縮・安全チェック",
		AlwaysOn:    true,
	}
}

func (h *URLHandler) Commands() []discord.ApplicationCommandCreate {
	return []discord.ApplicationCommandCreate{
		discord.SlashCommandCreate{
			Name:        "url",
			Description: "URLツール（短縮・安全チェック）",
		},
	}
}

func (h *URLHandler) HandleCommand(e *events.ApplicationCommandInteractionCreate) {
	hasXGD := h.cfg.XGDAPIKey != ""
	hasVT := h.cfg.VTAPIKey != ""

	_ = e.CreateMessage(view.BuildURLMainPanel(hasXGD, hasVT))
}

func (h *URLHandler) SettingsPanel(_ snowflake.ID) []discord.LayoutComponent { return nil }
