package embedfix

import (
	"context"
	"log/slog"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

func (ef *EmbedFix) processMessageURLs(ctx context.Context, e *events.GuildMessageCreate) {
	refs := extractEmbedURLs(e.Message.Content)
	if len(refs) == 0 {
		return
	}

	settings, err := LoadSettings(ef.store, e.GuildID)
	if err != nil {
		ef.logger.Error("failed to load embedfix settings", slog.Any("error", err))
		settings = defaultSettings()
	}

	// Filter refs by enabled platforms
	var enabledRefs []EmbedRef
	for _, ref := range refs {
		if settings.IsPlatformEnabled(ref.Platform) {
			enabledRefs = append(enabledRefs, ref)
		}
	}
	if len(enabledRefs) == 0 {
		return
	}

	// Suppress embeds on the original message (best-effort)
	if _, err := ef.client.Rest.UpdateMessage(e.ChannelID, e.MessageID,
		discord.NewMessageUpdate().WithSuppressEmbeds(true)); err != nil {
		ef.logger.Debug("failed to suppress embeds",
			slog.Any("error", err),
			slog.String("message_id", e.MessageID.String()),
		)
	}

	for _, ref := range enabledRefs {
		switch ref.Platform {
		case PlatformTwitter:
			ef.processTwitterEmbed(ctx, e, ref)
		case PlatformReddit:
			ef.processRedditEmbed(ctx, e, ref)
		case PlatformTikTok:
			ef.processTikTokEmbed(ctx, e, ref)
		}
	}
}

func (ef *EmbedFix) processTwitterEmbed(ctx context.Context, e *events.GuildMessageCreate, ref EmbedRef) {
	screenName, tweetID := ref.Params[0], ref.Params[1]

	tweet, err := ef.twitterClient.GetTweet(ctx, screenName, tweetID)
	if err != nil {
		ef.logger.Warn("failed to fetch tweet",
			slog.String("screen_name", screenName),
			slog.String("tweet_id", tweetID),
			slog.Any("error", err),
		)
		return
	}

	msg := BuildTweetEmbed(tweet, ref)
	if _, err = ef.client.Rest.CreateMessage(e.ChannelID, msg.WithMessageReferenceByID(e.MessageID)); err != nil {
		ef.logger.Warn("failed to send tweet embed",
			slog.String("tweet_id", tweetID),
			slog.Any("error", err),
		)
	}
}

func (ef *EmbedFix) processRedditEmbed(ctx context.Context, e *events.GuildMessageCreate, ref EmbedRef) {
	subreddit, postID := ref.Params[0], ref.Params[1]

	post, err := ef.redditClient.GetPost(ctx, subreddit, postID)
	if err != nil {
		ef.logger.Warn("failed to fetch reddit post",
			slog.String("subreddit", subreddit),
			slog.String("post_id", postID),
			slog.Any("error", err),
		)
		return
	}

	msg := BuildRedditEmbed(post, ref)
	if _, err = ef.client.Rest.CreateMessage(e.ChannelID, msg.WithMessageReferenceByID(e.MessageID)); err != nil {
		ef.logger.Warn("failed to send reddit embed",
			slog.String("post_id", postID),
			slog.Any("error", err),
		)
	}
}

func (ef *EmbedFix) processTikTokEmbed(ctx context.Context, e *events.GuildMessageCreate, ref EmbedRef) {
	username, videoID := ref.Params[0], ref.Params[1]

	video, err := ef.tiktokClient.GetVideo(ctx, username, videoID)
	if err != nil {
		ef.logger.Warn("failed to fetch tiktok video",
			slog.String("username", username),
			slog.String("video_id", videoID),
			slog.Any("error", err),
		)
		return
	}

	msg := BuildTikTokEmbed(video, ref)
	if _, err = ef.client.Rest.CreateMessage(e.ChannelID, msg.WithMessageReferenceByID(e.MessageID)); err != nil {
		ef.logger.Warn("failed to send tiktok embed",
			slog.String("video_id", videoID),
			slog.Any("error", err),
		)
	}
}
