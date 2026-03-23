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

// PingHandler implements module.Module for the ping feature.
type PingHandler struct {
	logger *slog.Logger
}

// NewPingHandler creates a new PingHandler.
func NewPingHandler(logger *slog.Logger) *PingHandler {
	return &PingHandler{logger: logger}
}

func (h *PingHandler) Info() module.Info {
	return module.Info{
		ID:          model.PingModuleID,
		Name:        "Ping",
		Description: "Botの応答確認",
		AlwaysOn:    true,
	}
}

func (h *PingHandler) Commands() []discord.ApplicationCommandCreate {
	return []discord.ApplicationCommandCreate{
		discord.SlashCommandCreate{
			Name:        "ping",
			Description: "Botの応答速度を確認する",
		},
	}
}

func (h *PingHandler) HandleCommand(e *events.ApplicationCommandInteractionCreate) {
	latency := e.Client().Gateway.Latency()
	_ = e.CreateMessage(view.BuildPingResponse(latency))
}

func (h *PingHandler) HandleComponent(_ *events.ComponentInteractionCreate) {}
func (h *PingHandler) HandleModal(_ *events.ModalSubmitInteractionCreate)   {}
func (h *PingHandler) SettingsPanel(_ snowflake.ID) []discord.LayoutComponent {
	return nil
}
