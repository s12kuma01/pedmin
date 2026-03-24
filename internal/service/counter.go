// opyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package service

import (
	"fmt"
	"log/slog"
	"regexp"
	"sync"
	"time"

	"github.com/disgoorg/snowflake/v2"

	"github.com/s12kuma01/pedmin/internal/model"
	"github.com/s12kuma01/pedmin/internal/repository"
)

// CounterService handles word counter business logic with an in-memory regex cache.
type CounterService struct {
	store  repository.GuildStore
	logger *slog.Logger
	mu     sync.RWMutex
	cache  map[snowflake.ID][]compiledCounter
}

type compiledCounter struct {
	id    int64
	regex *regexp.Regexp
}

// NewCounterService creates a new CounterService.
func NewCounterService(store repository.GuildStore, logger *slog.Logger) *CounterService {
	return &CounterService{
		store:  store,
		logger: logger,
		cache:  make(map[snowflake.ID][]compiledCounter),
	}
}

// AddCounter creates a new counter and invalidates the guild cache.
func (s *CounterService) AddCounter(guildID snowflake.ID, word string, matchType model.MatchType) (*model.Counter, error) {
	counter := &model.Counter{
		GuildID:   guildID,
		Word:      word,
		MatchType: matchType,
	}
	if err := s.store.CreateCounter(counter); err != nil {
		return nil, err
	}
	s.invalidateCache(guildID)
	return counter, nil
}

// DeleteCounter removes a counter and invalidates the guild cache.
func (s *CounterService) DeleteCounter(id int64, guildID snowflake.ID) error {
	if err := s.store.DeleteCounter(id, guildID); err != nil {
		return err
	}
	s.invalidateCache(guildID)
	return nil
}

// DeleteCounterAndList removes a counter and returns the remaining counters.
func (s *CounterService) DeleteCounterAndList(id int64, guildID snowflake.ID) ([]model.Counter, error) {
	if err := s.DeleteCounter(id, guildID); err != nil {
		return nil, err
	}
	return s.store.GetCounters(guildID)
}

// GetCounters returns all counters for a guild.
func (s *CounterService) GetCounters(guildID snowflake.ID) ([]model.Counter, error) {
	return s.store.GetCounters(guildID)
}

// GetCounter returns a single counter by ID and guild.
func (s *CounterService) GetCounter(id int64, guildID snowflake.ID) (*model.Counter, error) {
	return s.store.GetCounter(id, guildID)
}

// CountCounters returns the number of counters for a guild.
func (s *CounterService) CountCounters(guildID snowflake.ID) (int, error) {
	return s.store.CountCounters(guildID)
}

// ProcessMessage checks a message against all guild counters and records hits.
func (s *CounterService) ProcessMessage(guildID, userID snowflake.ID, content string) {
	if content == "" {
		return
	}

	patterns := s.getOrLoadPatterns(guildID)
	if len(patterns) == 0 {
		return
	}

	var hits []repository.CounterHit
	for _, p := range patterns {
		if p.regex.MatchString(content) {
			hits = append(hits, repository.CounterHit{
				CounterID: p.id,
				GuildID:   guildID,
				UserID:    userID,
			})
		}
	}

	if len(hits) > 0 {
		if err := s.store.RecordHits(hits); err != nil {
			s.logger.Error("failed to record counter hits", slog.Any("error", err))
		}
	}
}

// GetStats returns aggregated hit counts for all counters in a guild.
func (s *CounterService) GetStats(guildID snowflake.ID, period model.StatsPeriod) ([]model.CounterStat, error) {
	since := periodToTime(period)
	return s.store.GetCounterStats(guildID, since)
}

// GetUserRanking returns per-user hit rankings for a specific counter.
func (s *CounterService) GetUserRanking(counterID int64, period model.StatsPeriod) ([]model.CounterUserRank, error) {
	since := periodToTime(period)
	return s.store.GetCounterUserRanking(counterID, since, 10)
}

func (s *CounterService) getOrLoadPatterns(guildID snowflake.ID) []compiledCounter {
	s.mu.RLock()
	if patterns, ok := s.cache[guildID]; ok {
		s.mu.RUnlock()
		return patterns
	}
	s.mu.RUnlock()

	s.mu.Lock()
	defer s.mu.Unlock()

	// Double-check after acquiring write lock.
	if patterns, ok := s.cache[guildID]; ok {
		return patterns
	}

	counters, err := s.store.GetCounters(guildID)
	if err != nil {
		s.logger.Error("failed to load counters for cache", slog.Any("error", err))
		return nil
	}

	patterns := make([]compiledCounter, 0, len(counters))
	for _, c := range counters {
		re, err := buildRegex(c.Word, c.MatchType)
		if err != nil {
			s.logger.Error("failed to compile counter regex",
				slog.String("word", c.Word),
				slog.Any("error", err),
			)
			continue
		}
		patterns = append(patterns, compiledCounter{id: c.ID, regex: re})
	}

	s.cache[guildID] = patterns
	return patterns
}

func (s *CounterService) invalidateCache(guildID snowflake.ID) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.cache, guildID)
}

func buildRegex(word string, matchType model.MatchType) (*regexp.Regexp, error) {
	escaped := regexp.QuoteMeta(word)
	var pattern string
	switch matchType {
	case model.MatchExact:
		pattern = fmt.Sprintf("(?i)^%s$", escaped)
	case model.MatchWord:
		pattern = fmt.Sprintf(`(?i)\b%s\b`, escaped)
	default: // MatchPartial
		pattern = fmt.Sprintf("(?i)%s", escaped)
	}
	return regexp.Compile(pattern)
}

func periodToTime(period model.StatsPeriod) *time.Time {
	now := time.Now()
	var t time.Time
	switch period {
	case model.PeriodToday:
		t = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	case model.PeriodWeek:
		t = now.AddDate(0, 0, -7)
	default: // PeriodAllTime
		return nil
	}
	return &t
}
