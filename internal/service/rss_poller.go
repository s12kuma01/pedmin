package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/disgoorg/snowflake/v2"
	"github.com/s12kuma01/pedmin/internal/model"
	"github.com/s12kuma01/pedmin/internal/repository"
)

// RSSPollerBot is the interface the poller uses to check if a module is enabled.
type RSSPollerBot interface {
	IsModuleEnabled(guildID snowflake.ID, moduleID string) bool
}

// RSSPoller runs background polling for RSS feeds.
type RSSPoller struct {
	bot          RSSPollerBot
	service      *RSSService
	store        repository.GuildStore
	pollInterval time.Duration
	cancel       context.CancelFunc
	logger       *slog.Logger
}

// NewRSSPoller creates a new RSSPoller.
func NewRSSPoller(bot RSSPollerBot, svc *RSSService, store repository.GuildStore, pollInterval time.Duration, logger *slog.Logger) *RSSPoller {
	return &RSSPoller{
		bot:          bot,
		service:      svc,
		store:        store,
		pollInterval: pollInterval,
		logger:       logger,
	}
}

// StartPoller starts the background polling goroutine.
func (p *RSSPoller) StartPoller(ctx context.Context) {
	ctx, p.cancel = context.WithCancel(ctx)
	go p.pollLoop(ctx)
}

// StopPoller stops the background polling goroutine.
func (p *RSSPoller) StopPoller() {
	if p.cancel != nil {
		p.cancel()
	}
}

func (p *RSSPoller) pollLoop(ctx context.Context) {
	p.logger.Info("rss poller started", slog.Duration("interval", p.pollInterval))

	ticker := time.NewTicker(p.pollInterval)
	defer ticker.Stop()

	// Prune old seen items on startup
	p.pruneOldItems()

	for {
		select {
		case <-ctx.Done():
			p.logger.Info("rss poller stopped")
			return
		case <-ticker.C:
			p.pollAllFeeds(ctx)
		}
	}
}

func (p *RSSPoller) pollAllFeeds(ctx context.Context) {
	feeds, err := p.store.GetAllRSSFeeds()
	if err != nil {
		p.logger.Error("failed to get all rss feeds", slog.Any("error", err))
		return
	}

	for _, feed := range feeds {
		select {
		case <-ctx.Done():
			return
		default:
		}

		if !p.bot.IsModuleEnabled(feed.GuildID, model.RSSModuleID) {
			continue
		}

		if err := p.service.PollSingleFeed(ctx, feed); err != nil {
			p.logger.Warn("failed to poll feed",
				slog.Int64("feed_id", feed.ID),
				slog.String("url", feed.URL),
				slog.Any("error", err),
			)
		}
	}

	// Prune old seen items periodically
	p.pruneOldItems()
}

func (p *RSSPoller) pruneOldItems() {
	cutoff := time.Now().Add(-30 * 24 * time.Hour)
	if err := p.store.PruneSeenItems(cutoff); err != nil {
		p.logger.Warn("failed to prune seen items", slog.Any("error", err))
	}
}
