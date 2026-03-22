package embedfix

import (
	"context"
	"log/slog"

	disgobot "github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

func SetupListeners(client *disgobot.Client, ef *EmbedFix) {
	client.AddEventListeners(
		disgobot.NewListenerFunc(ef.onMessageCreate),
	)
}

func (ef *EmbedFix) onMessageCreate(e *events.GuildMessageCreate) {
	if e.Message.Author.Bot {
		return
	}

	guildID := e.GuildID
	if !ef.bot.IsModuleEnabled(guildID, ModuleID) {
		return
	}

	refs := extractEmbedURLs(e.Message.Content)
	if len(refs) == 0 {
		return
	}

	// Suppress embeds on the original message (best-effort)
	_, err := ef.client.Rest.UpdateMessage(e.ChannelID, e.MessageID,
		discord.NewMessageUpdate().WithSuppressEmbeds(true))
	if err != nil {
		ef.logger.Debug("failed to suppress embeds",
			slog.Any("error", err),
			slog.String("message_id", e.MessageID.String()),
		)
	}

	ctx := context.Background()
	for _, ref := range refs {
		switch ref.Platform {
		case PlatformTwitter:
			ef.handleTwitterEmbed(ctx, e, ref)
		case PlatformReddit:
			ef.handleRedditEmbed(ctx, e, ref)
		case PlatformTikTok:
			ef.handleTikTokEmbed(ctx, e, ref)
		case PlatformInstagram:
			ef.handleInstagramEmbed(ctx, e, ref)
		}
	}
}

func (ef *EmbedFix) handleTwitterEmbed(ctx context.Context, e *events.GuildMessageCreate, ref EmbedRef) {
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

func (ef *EmbedFix) handleRedditEmbed(ctx context.Context, e *events.GuildMessageCreate, ref EmbedRef) {
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

func (ef *EmbedFix) handleTikTokEmbed(ctx context.Context, e *events.GuildMessageCreate, ref EmbedRef) {
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

func (ef *EmbedFix) handleInstagramEmbed(ctx context.Context, e *events.GuildMessageCreate, ref EmbedRef) {
	if !ef.instagramClient.IsAvailable() {
		return
	}

	shortcode := ref.Params[0]

	post, err := ef.instagramClient.GetPost(ctx, shortcode)
	if err != nil {
		ef.logger.Warn("failed to fetch instagram post",
			slog.String("shortcode", shortcode),
			slog.Any("error", err),
		)
		return
	}

	msg := BuildInstagramEmbed(post, ref)
	if _, err = ef.client.Rest.CreateMessage(e.ChannelID, msg.WithMessageReferenceByID(e.MessageID)); err != nil {
		ef.logger.Warn("failed to send instagram embed",
			slog.String("shortcode", shortcode),
			slog.Any("error", err),
		)
	}
}
