// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package view

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/s12kuma01/pedmin/internal/model"
	"github.com/s12kuma01/pedmin/internal/ui"
)

// CounterManagePanel builds the counter management panel.
func CounterManagePanel(counters []model.Counter) discord.MessageCreate {
	if len(counters) == 0 {
		return ui.EphemeralV2(
			discord.NewContainer(
				discord.NewTextDisplay("登録されているカウンターはありません。"),
			),
		)
	}

	var options []discord.StringSelectMenuOption
	for _, c := range counters {
		options = append(options, discord.StringSelectMenuOption{
			Label:       c.Word,
			Value:       fmt.Sprintf("%d", c.ID),
			Description: model.MatchTypeLabel(c.MatchType),
		})
	}

	return ui.EphemeralV2(
		discord.NewContainer(
			discord.NewTextDisplay("### カウンター管理"),
			discord.NewSmallSeparator(),
			discord.NewTextDisplay(fmt.Sprintf("登録カウンター: %d/%d", len(counters), model.MaxCountersPerGuild)),
			discord.NewActionRow(
				discord.NewStringSelectMenu(model.CounterModuleID+":manage_select", "カウンターを選択...", options...),
			),
		),
	)
}

// CounterDetail builds the counter detail panel.
func CounterDetail(counter model.Counter) discord.MessageUpdate {
	return discord.NewMessageUpdateV2([]discord.LayoutComponent{
		discord.NewContainer(
			discord.NewTextDisplay(fmt.Sprintf("### %s", counter.Word)),
			discord.NewSmallSeparator(),
			discord.NewTextDisplay(fmt.Sprintf(
				"**マッチタイプ:** %s\n**追加日:** %s",
				model.MatchTypeLabel(counter.MatchType),
				counter.CreatedAt.Format("2006-01-02"),
			)),
			discord.NewLargeSeparator(),
			discord.NewActionRow(
				discord.NewDangerButton("削除", fmt.Sprintf("%s:delete:%d", model.CounterModuleID, counter.ID)),
			),
		),
	})
}

// CounterAddTypePrompt builds the match type selection panel.
func CounterAddTypePrompt() discord.MessageCreate {
	var options []discord.StringSelectMenuOption
	for _, mt := range model.AllMatchTypes {
		options = append(options, discord.StringSelectMenuOption{
			Label: mt.Label,
			Value: string(mt.Key),
		})
	}

	return ui.EphemeralV2(
		discord.NewContainer(
			discord.NewTextDisplay("### カウンター追加"),
			discord.NewSmallSeparator(),
			discord.NewTextDisplay("マッチタイプを選択してください:"),
			discord.NewActionRow(
				discord.NewStringSelectMenu(model.CounterModuleID+":add_type", "マッチタイプを選択...", options...),
			),
		),
	)
}
