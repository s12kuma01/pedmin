package rss

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"log/slog"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
	"github.com/mmcdole/gofeed"
	"github.com/s12kuma01/pedmin/store"
)

const MaxFeedsPerGuild = 10

func (r *RSS) AddFeed(ctx context.Context, guildID snowflake.ID, channelID snowflake.ID, url string) (*store.RSSFeed, error) {
	count, err := r.store.CountRSSFeeds(guildID)
	if err != nil {
		return nil, fmt.Errorf("failed to count feeds: %w", err)
	}
	if count >= MaxFeedsPerGuild {
		return nil, fmt.Errorf("フィード数が上限（%d件）に達しています", MaxFeedsPerGuild)
	}

	parser := gofeed.NewParser()
	parsed, err := parser.ParseURLWithContext(url, ctx)
	if err != nil {
		return nil, fmt.Errorf("フィードの取得に失敗しました: %w", err)
	}

	feed := &store.RSSFeed{
		GuildID:   guildID,
		URL:       url,
		ChannelID: channelID,
		Title:     parsed.Title,
	}

	if err := r.store.CreateRSSFeed(feed); err != nil {
		if errors.Is(err, store.ErrDuplicateFeed) {
			return nil, fmt.Errorf("このフィードは既に登録されています")
		}
		return nil, fmt.Errorf("failed to create feed: %w", err)
	}

	// Mark all existing items as seen to prevent flood
	var hashes []string
	for _, item := range parsed.Items {
		hashes = append(hashes, itemHash(item))
	}
	if len(hashes) > 0 {
		if err := r.store.MarkItemsSeen(feed.ID, hashes); err != nil {
			r.logger.Warn("failed to mark existing items as seen", slog.Any("error", err))
		}
	}

	if len(parsed.Items) > 0 {
		msg := BuildFeedAnnouncement(feed.Title, parsed.Items[0])
		if _, err := (*r.client).Rest.CreateMessage(channelID, msg); err != nil {
			r.logger.Warn("failed to post preview", slog.Any("error", err))
		}
	}

	return feed, nil
}

func (r *RSS) RemoveFeed(feedID int64, guildID snowflake.ID) error {
	return r.store.DeleteRSSFeed(feedID, guildID)
}

// GetFeeds returns all RSS feeds for a guild.
func (r *RSS) GetFeeds(guildID snowflake.ID) ([]store.RSSFeed, error) {
	return r.store.GetRSSFeeds(guildID)
}

// GetFeed returns a single feed by ID within a guild.
func (r *RSS) GetFeed(guildID snowflake.ID, feedID int64) (*store.RSSFeed, error) {
	feeds, err := r.store.GetRSSFeeds(guildID)
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
func (r *RSS) DeleteFeedAndList(feedID int64, guildID snowflake.ID) ([]store.RSSFeed, error) {
	if err := r.RemoveFeed(feedID, guildID); err != nil {
		return nil, err
	}
	return r.GetFeeds(guildID)
}

func itemHash(item *gofeed.Item) string {
	key := item.GUID
	if key == "" {
		key = item.Link
	}
	h := sha256.Sum256([]byte(key))
	return fmt.Sprintf("%x", h)
}

func ephemeralV2(components ...discord.LayoutComponent) discord.MessageCreate {
	return discord.NewMessageCreateV2(components...).WithEphemeral(true)
}
