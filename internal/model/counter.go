// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package model

import (
	"time"

	"github.com/disgoorg/snowflake/v2"
)

// MatchType defines how a counter word is matched against messages.
type MatchType string

const (
	MatchPartial MatchType = "partial" // 部分一致 (substring)
	MatchExact   MatchType = "exact"   // 完全一致 (entire message)
	MatchWord    MatchType = "word"    // 単語一致 (word boundary)
)

// AllMatchTypes lists all match types with display labels.
var AllMatchTypes = []struct {
	Key   MatchType
	Label string
}{
	{MatchPartial, "部分一致"},
	{MatchExact, "完全一致"},
	{MatchWord, "単語一致"},
}

// MatchTypeLabel returns the Japanese label for a match type.
func MatchTypeLabel(mt MatchType) string {
	for _, m := range AllMatchTypes {
		if m.Key == mt {
			return m.Label
		}
	}
	return string(mt)
}

// Counter represents a registered word counter for a guild.
type Counter struct {
	ID        int64
	GuildID   snowflake.ID
	Word      string
	MatchType MatchType
	CreatedAt time.Time
}

// CounterStat holds aggregated hit count for a counter.
type CounterStat struct {
	CounterID int64
	Word      string
	MatchType MatchType
	HitCount  int
}

// CounterUserRank holds per-user hit count ranking.
type CounterUserRank struct {
	UserID   snowflake.ID
	HitCount int
}

// StatsPeriod defines the time range for stats queries.
type StatsPeriod string

const (
	PeriodToday   StatsPeriod = "today"
	PeriodWeek    StatsPeriod = "week"
	PeriodAllTime StatsPeriod = "all"
)

// AllStatsPeriods lists all available period filters.
var AllStatsPeriods = []struct {
	Key   StatsPeriod
	Label string
}{
	{PeriodToday, "今日"},
	{PeriodWeek, "今週"},
	{PeriodAllTime, "全期間"},
}
