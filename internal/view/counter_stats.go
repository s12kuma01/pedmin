// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package view

import (
	"fmt"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/Sumire-Labs/pedmin/internal/model"
	"github.com/Sumire-Labs/pedmin/internal/ui"
)

// CounterStatsPanel builds the stats panel with word hit counts.
func CounterStatsPanel(stats []model.CounterStat, period model.StatsPeriod) discord.MessageCreate {
	periodLabel := periodLabel(period)

	var statsText strings.Builder
	if len(stats) == 0 {
		statsText.WriteString("データがありません。")
	} else {
		for i, st := range stats {
			fmt.Fprintf(&statsText, "%d. **%s** — %d回\n", i+1, st.Word, st.HitCount)
		}
	}

	// Period select menu
	var periodOptions []discord.StringSelectMenuOption
	for _, p := range model.AllStatsPeriods {
		opt := discord.StringSelectMenuOption{
			Label: p.Label,
			Value: string(p.Key),
		}
		if p.Key == period {
			opt = opt.WithDefault(true)
		}
		periodOptions = append(periodOptions, opt)
	}

	components := []discord.ContainerSubComponent{
		discord.NewTextDisplay("### ワードカウンター統計"),
		discord.NewSmallSeparator(),
		discord.NewTextDisplay(fmt.Sprintf("**期間:** %s", periodLabel)),
		discord.NewSmallSeparator(),
		discord.NewTextDisplay(statsText.String()),
		discord.NewLargeSeparator(),
		discord.NewActionRow(
			discord.NewStringSelectMenu(model.CounterModuleID+":stats_period", "期間を選択...", periodOptions...),
		),
	}

	// Counter select menu for user ranking (only if there are stats with hits)
	var counterOptions []discord.StringSelectMenuOption
	for _, st := range stats {
		if st.HitCount > 0 {
			counterOptions = append(counterOptions, discord.StringSelectMenuOption{
				Label:       st.Word,
				Value:       fmt.Sprintf("%d", st.CounterID),
				Description: fmt.Sprintf("%d回", st.HitCount),
			})
		}
	}
	if len(counterOptions) > 0 {
		components = append(components,
			discord.NewActionRow(
				discord.NewStringSelectMenu(model.CounterModuleID+":stats_select", "ランキングを見る...", counterOptions...),
			),
		)
	}

	return ui.EphemeralV2(discord.NewContainer(components...))
}

// CounterUserRanking builds the user ranking panel for a specific counter.
func CounterUserRanking(word string, ranks []model.CounterUserRank, period model.StatsPeriod) discord.MessageUpdate {
	periodLabel := periodLabel(period)

	var rankText strings.Builder
	if len(ranks) == 0 {
		rankText.WriteString("データがありません。")
	} else {
		medals := []string{"\U0001f947", "\U0001f948", "\U0001f949"}
		for i, r := range ranks {
			prefix := fmt.Sprintf("%d.", i+1)
			if i < len(medals) {
				prefix = medals[i]
			}
			fmt.Fprintf(&rankText, "%s <@%d> — %d回\n", prefix, r.UserID, r.HitCount)
		}
	}

	return discord.NewMessageUpdateV2([]discord.LayoutComponent{
		discord.NewContainer(
			discord.NewTextDisplay(fmt.Sprintf("### \"%s\" のランキング", word)),
			discord.NewSmallSeparator(),
			discord.NewTextDisplay(fmt.Sprintf("**期間:** %s", periodLabel)),
			discord.NewSmallSeparator(),
			discord.NewTextDisplay(rankText.String()),
		),
	})
}

func periodLabel(period model.StatsPeriod) string {
	for _, p := range model.AllStatsPeriods {
		if p.Key == period {
			return p.Label
		}
	}
	return "全期間"
}
