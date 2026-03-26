// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package service

import (
	"context"
	"log/slog"

	disgobot "github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
	"github.com/Sumire-Labs/pedmin/internal/client"
	"github.com/Sumire-Labs/pedmin/internal/model"
	"github.com/Sumire-Labs/pedmin/internal/repository"
	"github.com/Sumire-Labs/pedmin/internal/view"
	"github.com/Sumire-Labs/pedmin/pkg/deepl"
)

// EmbedFixService handles URL processing, platform-specific embed sending, and embed suppression.
type EmbedFixService struct {
	store           repository.GuildStore
	twitterClient   *client.FxTwitterClient
	redditClient    *client.RedditClient
	tiktokClient    *client.TikTokClient
	translateClient *deepl.TranslateClient
	discordClient   *disgobot.Client
	logger          *slog.Logger
}

// NewEmbedFixService creates a new EmbedFixService.
func NewEmbedFixService(
	store repository.GuildStore,
	twitterClient *client.FxTwitterClient,
	redditClient *client.RedditClient,
	tiktokClient *client.TikTokClient,
	translateClient *deepl.TranslateClient,
	discordClient *disgobot.Client,
	logger *slog.Logger,
) *EmbedFixService {
	return &EmbedFixService{
		store:           store,
		twitterClient:   twitterClient,
		redditClient:    redditClient,
		tiktokClient:    tiktokClient,
		translateClient: translateClient,
		discordClient:   discordClient,
		logger:          logger,
	}
}

// LoadSettings loads embedfix settings for a guild.
func (s *EmbedFixService) LoadSettings(guildID snowflake.ID) (*model.EmbedFixSettings, error) {
	settings, err := repository.LoadModuleSettings(s.store, guildID, model.EmbedFixModuleID, model.DefaultEmbedFixSettings)
	if err != nil {
		return nil, err
	}
	if settings.Platforms == nil {
		return model.DefaultEmbedFixSettings(), nil
	}
	return settings, nil
}

// SaveSettings saves embedfix settings for a guild.
func (s *EmbedFixService) SaveSettings(guildID snowflake.ID, settings *model.EmbedFixSettings) error {
	return repository.SaveModuleSettings(s.store, guildID, model.EmbedFixModuleID, settings)
}

// ProcessMessageURLs extracts and processes embed URLs from a message.
func (s *EmbedFixService) ProcessMessageURLs(ctx context.Context, guildID, channelID, messageID snowflake.ID, content string) {
	refs := model.ExtractEmbedURLs(content)
	if len(refs) == 0 {
		return
	}

	settings, err := s.LoadSettings(guildID)
	if err != nil {
		s.logger.Error("failed to load embedfix settings", slog.Any("error", err))
		settings = model.DefaultEmbedFixSettings()
	}

	// Filter refs by enabled platforms
	var enabledRefs []model.EmbedRef
	for _, ref := range refs {
		if settings.IsPlatformEnabled(ref.Platform) {
			enabledRefs = append(enabledRefs, ref)
		}
	}
	if len(enabledRefs) == 0 {
		return
	}

	// Suppress embeds on the original message (best-effort)
	if _, err := s.discordClient.Rest.UpdateMessage(channelID, messageID,
		discord.NewMessageUpdate().WithSuppressEmbeds(true)); err != nil {
		s.logger.Debug("failed to suppress embeds",
			slog.Any("error", err),
			slog.String("message_id", messageID.String()),
		)
	}

	for _, ref := range enabledRefs {
		switch ref.Platform {
		case model.PlatformTwitter:
			s.processTwitterEmbed(ctx, channelID, messageID, ref)
		case model.PlatformReddit:
			s.processRedditEmbed(ctx, channelID, messageID, ref)
		case model.PlatformTikTok:
			s.processTikTokEmbed(ctx, channelID, messageID, ref)
		}
	}
}

func (s *EmbedFixService) processTwitterEmbed(ctx context.Context, channelID, messageID snowflake.ID, ref model.EmbedRef) {
	screenName, tweetID := ref.Params[0], ref.Params[1]

	tweet, err := s.twitterClient.GetTweet(ctx, screenName, tweetID)
	if err != nil {
		s.logger.Warn("failed to fetch tweet",
			slog.String("screen_name", screenName),
			slog.String("tweet_id", tweetID),
			slog.Any("error", err),
		)
		return
	}

	msg := view.BuildTweetEmbed(tweet, ref)
	if _, err = s.discordClient.Rest.CreateMessage(channelID, msg.WithMessageReferenceByID(messageID).WithAllowedMentions(&discord.AllowedMentions{})); err != nil {
		s.logger.Warn("failed to send tweet embed",
			slog.String("tweet_id", tweetID),
			slog.Any("error", err),
		)
	}
}

func (s *EmbedFixService) processRedditEmbed(ctx context.Context, channelID, messageID snowflake.ID, ref model.EmbedRef) {
	subreddit, postID := ref.Params[0], ref.Params[1]

	post, err := s.redditClient.GetPost(ctx, subreddit, postID)
	if err != nil {
		s.logger.Warn("failed to fetch reddit post",
			slog.String("subreddit", subreddit),
			slog.String("post_id", postID),
			slog.Any("error", err),
		)
		return
	}

	msg := view.BuildRedditEmbed(post, ref)
	if _, err = s.discordClient.Rest.CreateMessage(channelID, msg.WithMessageReferenceByID(messageID).WithAllowedMentions(&discord.AllowedMentions{})); err != nil {
		s.logger.Warn("failed to send reddit embed",
			slog.String("post_id", postID),
			slog.Any("error", err),
		)
	}
}

func (s *EmbedFixService) processTikTokEmbed(ctx context.Context, channelID, messageID snowflake.ID, ref model.EmbedRef) {
	username, videoID := ref.Params[0], ref.Params[1]

	video, err := s.tiktokClient.GetVideo(ctx, username, videoID)
	if err != nil {
		s.logger.Warn("failed to fetch tiktok video",
			slog.String("username", username),
			slog.String("video_id", videoID),
			slog.Any("error", err),
		)
		return
	}

	msg := view.BuildTikTokEmbed(video, ref)
	if _, err = s.discordClient.Rest.CreateMessage(channelID, msg.WithMessageReferenceByID(messageID).WithAllowedMentions(&discord.AllowedMentions{})); err != nil {
		s.logger.Warn("failed to send tiktok embed",
			slog.String("video_id", videoID),
			slog.Any("error", err),
		)
	}
}
