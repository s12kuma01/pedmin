package service

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"log/slog"

	disgobot "github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/snowflake/v2"
	"github.com/mmcdole/gofeed"
	"github.com/s12kuma01/pedmin/internal/model"
	"github.com/s12kuma01/pedmin/internal/repository"
	"github.com/s12kuma01/pedmin/internal/view"
)

// RSSService handles RSS feed CRUD and posting logic.
type RSSService struct {
	store  repository.GuildStore
	client *disgobot.Client
	logger *slog.Logger
}

// NewRSSService creates a new RSSService.
func NewRSSService(store repository.GuildStore, client *disgobot.Client, logger *slog.Logger) *RSSService {
	return &RSSService{
		store:  store,
		client: client,
		logger: logger,
	}
}

// AddFeed validates and adds a new RSS feed, posting a preview of the latest item.
func (s *RSSService) AddFeed(ctx context.Context, guildID, channelID snowflake.ID, url string) (*model.RSSFeed, error) {
	count, err := s.store.CountRSSFeeds(guildID)
	if err != nil {
		return nil, fmt.Errorf("failed to count feeds: %w", err)
	}
	if count >= model.MaxRSSFeedsPerGuild {
		return nil, fmt.Errorf("フィード数が上限（%d件）に達しています", model.MaxRSSFeedsPerGuild)
	}

	parser := gofeed.NewParser()
	parsed, err := parser.ParseURLWithContext(url, ctx)
	if err != nil {
		return nil, fmt.Errorf("フィードの取得に失敗しました: %w", err)
	}

	feed := &model.RSSFeed{
		GuildID:   guildID,
		URL:       url,
		ChannelID: channelID,
		Title:     parsed.Title,
	}

	if err := s.store.CreateRSSFeed(feed); err != nil {
		if errors.Is(err, model.ErrDuplicateFeed) {
			return nil, fmt.Errorf("このフィードは既に登録されています")
		}
		return nil, fmt.Errorf("failed to create feed: %w", err)
	}

	// Mark all existing items as seen to prevent flood
	var hashes []string
	for _, item := range parsed.Items {
		hashes = append(hashes, rssItemHash(item))
	}
	if len(hashes) > 0 {
		if err := s.store.MarkItemsSeen(feed.ID, hashes); err != nil {
			s.logger.Warn("failed to mark existing items as seen", slog.Any("error", err))
		}
	}

	if len(parsed.Items) > 0 {
		msg := view.RSSFeedAnnouncement(feed.Title, parsed.Items[0])
		if _, err := (*s.client).Rest.CreateMessage(channelID, msg); err != nil {
			s.logger.Warn("failed to post preview", slog.Any("error", err))
		}
	}

	return feed, nil
}

// RemoveFeed deletes an RSS feed.
func (s *RSSService) RemoveFeed(feedID int64, guildID snowflake.ID) error {
	return s.store.DeleteRSSFeed(feedID, guildID)
}

// GetFeeds returns all RSS feeds for a guild.
func (s *RSSService) GetFeeds(guildID snowflake.ID) ([]model.RSSFeed, error) {
	return s.store.GetRSSFeeds(guildID)
}

// GetFeed returns a single feed by ID within a guild.
func (s *RSSService) GetFeed(guildID snowflake.ID, feedID int64) (*model.RSSFeed, error) {
	feeds, err := s.store.GetRSSFeeds(guildID)
	if err != nil {
		return nil, err
	}
	for _, f := range feeds {
		if f.ID == feedID {
			return &f, nil
		}
	}
	return nil, fmt.Errorf("feed not found")
}

// DeleteFeedAndList deletes a feed and returns the remaining feed list.
func (s *RSSService) DeleteFeedAndList(feedID int64, guildID snowflake.ID) ([]model.RSSFeed, error) {
	if err := s.RemoveFeed(feedID, guildID); err != nil {
		return nil, err
	}
	return s.GetFeeds(guildID)
}

// CountFeeds returns the number of RSS feeds for a guild.
func (s *RSSService) CountFeeds(guildID snowflake.ID) (int, error) {
	return s.store.CountRSSFeeds(guildID)
}

// PollSingleFeed checks a single feed for new items and posts them.
func (s *RSSService) PollSingleFeed(ctx context.Context, feed model.RSSFeed) error {
	parser := gofeed.NewParser()
	parsed, err := parser.ParseURLWithContext(feed.URL, ctx)
	if err != nil {
		return fmt.Errorf("failed to parse feed: %w", err)
	}

	// Find new items (iterate in reverse to post oldest first)
	var newItems []*gofeed.Item
	var newHashes []string
	for i := len(parsed.Items) - 1; i >= 0; i-- {
		item := parsed.Items[i]
		hash := rssItemHash(item)
		seen, err := s.store.IsItemSeen(feed.ID, hash)
		if err != nil {
			return fmt.Errorf("failed to check seen: %w", err)
		}
		if !seen {
			newItems = append(newItems, item)
			newHashes = append(newHashes, hash)
		}
	}

	if len(newItems) == 0 {
		return nil
	}

	for _, item := range newItems {
		msg := view.RSSFeedAnnouncement(feed.Title, item)
		if _, err := (*s.client).Rest.CreateMessage(feed.ChannelID, msg); err != nil {
			s.logger.Warn("failed to post feed item",
				slog.Int64("feed_id", feed.ID),
				slog.Any("error", err),
			)
		}
	}

	if err := s.store.MarkItemsSeen(feed.ID, newHashes); err != nil {
		return fmt.Errorf("failed to mark items seen: %w", err)
	}

	s.logger.Info("posted new feed items",
		slog.Int64("feed_id", feed.ID),
		slog.Int("count", len(newItems)),
	)
	return nil
}

func rssItemHash(item *gofeed.Item) string {
	key := item.GUID
	if key == "" {
		key = item.Link
	}
	h := sha256.Sum256([]byte(key))
	return fmt.Sprintf("%x", h)
}
