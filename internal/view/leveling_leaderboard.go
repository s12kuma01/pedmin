// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package view

import (
	"fmt"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/s12kuma01/pedmin/internal/model"
	"github.com/s12kuma01/pedmin/internal/ui"
)

// LevelingLeaderboard builds the leaderboard as an ephemeral Components V2 message.
func LevelingLeaderboard(entries []model.LeaderboardEntry, page, totalPages int) discord.MessageCreate {
	return ui.EphemeralV2(levelingLeaderboardContainer(entries, page, totalPages))
}

// LevelingLeaderboardUpdate builds the leaderboard as a message update (for pagination).
func LevelingLeaderboardUpdate(entries []model.LeaderboardEntry, page, totalPages int) discord.MessageUpdate {
	return discord.NewMessageUpdateV2([]discord.LayoutComponent{levelingLeaderboardContainer(entries, page, totalPages)})
}

func levelingLeaderboardContainer(entries []model.LeaderboardEntry, page, totalPages int) discord.ContainerComponent {
	components := []discord.ContainerSubComponent{
		discord.NewTextDisplay("### レベルランキング"),
		discord.NewSmallSeparator(),
	}

	if len(entries) == 0 {
		components = append(components, discord.NewTextDisplay("データがありません。"))
	} else {
		medals := []string{"\U0001f947", "\U0001f948", "\U0001f949"}
		var sb strings.Builder
		for _, e := range entries {
			prefix := fmt.Sprintf("**%d.**", e.Rank)
			if e.Rank <= 3 {
				prefix = medals[e.Rank-1]
			}
			sb.WriteString(fmt.Sprintf("%s <@%d> — Lv.%d (%s XP)\n", prefix, e.UserID, e.Level, formatLeaderboardXP(e.TotalXP)))
		}
		components = append(components, discord.NewTextDisplay(sb.String()))
	}

	components = append(components, discord.NewSmallSeparator())
	components = append(components, discord.NewTextDisplay(fmt.Sprintf("ページ %d/%d", page+1, totalPages)))

	// Pagination buttons
	prevBtn := discord.NewSecondaryButton("← 前", fmt.Sprintf("%s:lb_page:%d", model.LevelingModuleID, page-1))
	if page <= 0 {
		prevBtn = prevBtn.AsDisabled()
	}
	nextBtn := discord.NewSecondaryButton("次 →", fmt.Sprintf("%s:lb_page:%d", model.LevelingModuleID, page+1))
	if page >= totalPages-1 {
		nextBtn = nextBtn.AsDisabled()
	}
	components = append(components, discord.NewActionRow(prevBtn, nextBtn))

	return discord.NewContainer(components...)
}

func formatLeaderboardXP(n int) string {
	if n < 1000 {
		return fmt.Sprintf("%d", n)
	}
	if n < 1000000 {
		return fmt.Sprintf("%d,%03d", n/1000, n%1000)
	}
	return fmt.Sprintf("%d,%03d,%03d", n/1000000, (n%1000000)/1000, n%1000)
}
