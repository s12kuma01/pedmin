// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package handler

import (
	"log/slog"
	"strconv"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgolink/v3/lavalink"
	"github.com/disgoorg/snowflake/v2"
	"github.com/s12kuma01/pedmin/internal/service"
	"github.com/s12kuma01/pedmin/internal/ui"
	"github.com/s12kuma01/pedmin/internal/view"
)

func (h *PlayerHandler) handleComponent(e *events.ComponentInteractionCreate) {
	customID := e.Data.CustomID()
	h.logger.Debug("component interaction received", slog.String("custom_id", customID))
	_, rest, _ := strings.Cut(customID, ":")

	guildID := e.GuildID()
	if guildID == nil {
		h.logger.Warn("component interaction: guildID is nil")
		return
	}

	action, _, _ := strings.Cut(rest, ":")

	switch action {
	case "skip":
		h.handleSkip(e, *guildID)
	case "stop":
		h.handleStop(e, *guildID)
	case "loop":
		h.handleLoop(e, *guildID)
	case "add":
		h.handleAddModal(e)
	case "queue":
		h.handleShowQueue(e, *guildID)
	case "back":
		h.handleBack(e, *guildID)
	case "clear_queue":
		h.handleClearQueue(e, *guildID)
	case "seek_back":
		h.handleSeek(e, *guildID, -service.SeekStep)
	case "seek_forward":
		h.handleSeek(e, *guildID, service.SeekStep)
	case "shuffle":
		h.handleShuffle(e, *guildID)
	case "volume":
		h.handleVolumeSettings(e, *guildID)
	}
}

func (h *PlayerHandler) handleSkip(e *events.ComponentInteractionCreate, guildID snowflake.ID) {
	player := h.service.Lavalink().ExistingPlayer(guildID)
	if player == nil {
		_ = e.DeferUpdateMessage()
		return
	}

	queue := h.service.Queues().Get(guildID)
	next, ok := queue.Next()
	if !ok {
		ctx, cancel := h.service.LavalinkCtx()
		defer cancel()
		_ = player.Update(ctx, lavalink.WithNullTrack())
		h.respondWithPlayerUpdate(e, guildID)
		return
	}

	ctx, cancel := h.service.LavalinkCtx()
	defer cancel()
	if err := player.Update(ctx, lavalink.WithTrack(next)); err != nil {
		h.logger.Error("failed to skip", slog.Any("error", err))
	}
	h.respondWithPlayerUpdate(e, guildID)
}

func (h *PlayerHandler) handleStop(e *events.ComponentInteractionCreate, guildID snowflake.ID) {
	h.service.Stop(guildID)

	playerUI := h.service.BuildPlayerUI(guildID)
	_ = e.UpdateMessage(discord.NewMessageUpdateV2([]discord.LayoutComponent{playerUI}))
}

func (h *PlayerHandler) handleLoop(e *events.ComponentInteractionCreate, guildID snowflake.ID) {
	h.service.CycleLoop(guildID)
	h.respondWithPlayerUpdate(e, guildID)
}

func (h *PlayerHandler) handleSeek(e *events.ComponentInteractionCreate, guildID snowflake.ID, delta lavalink.Duration) {
	player := h.service.Lavalink().ExistingPlayer(guildID)
	if player == nil || player.Track() == nil {
		_ = e.DeferUpdateMessage()
		return
	}
	if player.Track().Info.IsStream {
		_ = e.DeferUpdateMessage()
		return
	}

	h.service.Seek(guildID, delta)
	h.respondWithPlayerUpdate(e, guildID)
}

func (h *PlayerHandler) handleShuffle(e *events.ComponentInteractionCreate, guildID snowflake.ID) {
	h.service.Shuffle(guildID)
	h.respondWithPlayerUpdate(e, guildID)
}

func (h *PlayerHandler) respondWithPlayerUpdate(e *events.ComponentInteractionCreate, guildID snowflake.ID) {
	playerUI := h.service.BuildPlayerUI(guildID)
	_ = e.UpdateMessage(discord.NewMessageUpdateV2([]discord.LayoutComponent{playerUI}))
}

func (h *PlayerHandler) handleVolumeSettings(e *events.ComponentInteractionCreate, guildID snowflake.ID) {
	data, ok := e.Data.(discord.StringSelectMenuInteractionData)
	if !ok || len(data.Values) == 0 {
		return
	}

	vol, err := strconv.Atoi(data.Values[0])
	if err != nil {
		return
	}

	if err := h.service.SaveVolumeSettings(guildID, vol); err != nil {
		h.logger.Error("failed to save player settings", slog.Any("error", err))
	}

	settingsUI := view.BuildPlayerSettingsPanel(vol)
	_ = e.UpdateMessage(ui.BuildModulePanel(h.Info(), true, settingsUI))
}
