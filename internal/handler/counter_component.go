// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package handler

import (
	"log/slog"
	"strconv"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/Sumire-Labs/pedmin/internal/model"
	"github.com/Sumire-Labs/pedmin/internal/ui"
	"github.com/Sumire-Labs/pedmin/internal/view"
)

func (h *CounterHandler) HandleComponent(e *events.ComponentInteractionCreate) {
	customID := e.Data.CustomID()
	_, rest, _ := strings.Cut(customID, ":")
	action, extra, _ := strings.Cut(rest, ":")

	guildID := e.GuildID()
	if guildID == nil {
		return
	}

	switch action {
	case "add_prompt":
		h.counterHandleAddPrompt(e)

	case "add_type":
		h.counterHandleAddType(e)

	case "manage":
		h.counterHandleManage(e)

	case "manage_select":
		h.counterHandleManageSelect(e)

	case "delete":
		h.counterHandleDelete(e, extra)

	case "stats":
		h.counterHandleStats(e, model.PeriodAllTime)

	case "stats_period":
		h.counterHandleStatsPeriod(e)

	case "stats_select":
		h.counterHandleStatsSelect(e)
	}
}

func (h *CounterHandler) counterHandleAddPrompt(e *events.ComponentInteractionCreate) {
	_ = e.CreateMessage(view.CounterAddTypePrompt())
}

func (h *CounterHandler) counterHandleAddType(e *events.ComponentInteractionCreate) {
	data, ok := e.Data.(discord.StringSelectMenuInteractionData)
	if !ok || len(data.Values) == 0 {
		return
	}

	matchType := data.Values[0]

	_ = e.Modal(discord.ModalCreate{
		CustomID: model.CounterModuleID + ":add_modal:" + matchType,
		Title:    "カウンター追加",
		Components: []discord.LayoutComponent{
			discord.NewLabel("ワード",
				discord.NewShortTextInput(model.CounterModuleID+":word").
					WithRequired(true).
					WithPlaceholder("カウントするワードを入力"),
			),
		},
	})
}

func (h *CounterHandler) counterHandleManage(e *events.ComponentInteractionCreate) {
	counters, err := h.service.GetCounters(*e.GuildID())
	if err != nil {
		h.logger.Error("failed to get counters", slog.Any("error", err))
		_ = e.CreateMessage(ui.EphemeralError("カウンター一覧の取得に失敗しました。"))
		return
	}

	_ = e.CreateMessage(view.CounterManagePanel(counters))
}

func (h *CounterHandler) counterHandleManageSelect(e *events.ComponentInteractionCreate) {
	data, ok := e.Data.(discord.StringSelectMenuInteractionData)
	if !ok || len(data.Values) == 0 {
		return
	}

	counterID, err := strconv.ParseInt(data.Values[0], 10, 64)
	if err != nil {
		return
	}

	counter, err := h.service.GetCounter(counterID, *e.GuildID())
	if err != nil {
		h.logger.Error("failed to get counter", slog.Any("error", err))
		_ = e.UpdateMessage(discord.NewMessageUpdateV2([]discord.LayoutComponent{
			discord.NewContainer(
				discord.NewTextDisplay("カウンターが見つかりませんでした。"),
			),
		}))
		return
	}

	_ = e.UpdateMessage(view.CounterDetail(*counter))
}

func (h *CounterHandler) counterHandleDelete(e *events.ComponentInteractionCreate, counterIDStr string) {
	counterID, err := strconv.ParseInt(counterIDStr, 10, 64)
	if err != nil {
		return
	}

	counters, err := h.service.DeleteCounterAndList(counterID, *e.GuildID())
	if err != nil {
		h.logger.Error("failed to delete counter", slog.Any("error", err))
		_ = e.CreateMessage(ui.EphemeralError("カウンターの削除に失敗しました。"))
		return
	}

	if len(counters) == 0 {
		_ = e.UpdateMessage(discord.NewMessageUpdateV2([]discord.LayoutComponent{
			discord.NewContainer(
				discord.NewTextDisplay("登録されているカウンターはありません。"),
			),
		}))
		return
	}

	msg := view.CounterManagePanel(counters)
	_ = e.UpdateMessage(discord.NewMessageUpdateV2(msg.Components))
}

func (h *CounterHandler) counterHandleStats(e *events.ComponentInteractionCreate, period model.StatsPeriod) {
	stats, err := h.service.GetStats(*e.GuildID(), period)
	if err != nil {
		h.logger.Error("failed to get counter stats", slog.Any("error", err))
		_ = e.CreateMessage(ui.EphemeralError("統計の取得に失敗しました。"))
		return
	}

	_ = e.CreateMessage(view.CounterStatsPanel(stats, period))
}

func (h *CounterHandler) counterHandleStatsPeriod(e *events.ComponentInteractionCreate) {
	data, ok := e.Data.(discord.StringSelectMenuInteractionData)
	if !ok || len(data.Values) == 0 {
		return
	}

	period := model.StatsPeriod(data.Values[0])

	stats, err := h.service.GetStats(*e.GuildID(), period)
	if err != nil {
		h.logger.Error("failed to get counter stats", slog.Any("error", err))
		_ = e.UpdateMessage(discord.NewMessageUpdateV2([]discord.LayoutComponent{
			discord.NewContainer(
				discord.NewTextDisplay("統計の取得に失敗しました。"),
			),
		}))
		return
	}

	msg := view.CounterStatsPanel(stats, period)
	_ = e.UpdateMessage(discord.NewMessageUpdateV2(msg.Components))
}

func (h *CounterHandler) counterHandleStatsSelect(e *events.ComponentInteractionCreate) {
	data, ok := e.Data.(discord.StringSelectMenuInteractionData)
	if !ok || len(data.Values) == 0 {
		return
	}

	counterID, err := strconv.ParseInt(data.Values[0], 10, 64)
	if err != nil {
		return
	}

	// Extract current period from the select menu default value
	period := model.PeriodAllTime

	counter, err := h.service.GetCounter(counterID, *e.GuildID())
	if err != nil {
		h.logger.Error("failed to get counter", slog.Any("error", err))
		return
	}

	ranks, err := h.service.GetUserRanking(counterID, period)
	if err != nil {
		h.logger.Error("failed to get user ranking", slog.Any("error", err))
		return
	}

	_ = e.UpdateMessage(view.CounterUserRanking(counter.Word, ranks, period))
}
