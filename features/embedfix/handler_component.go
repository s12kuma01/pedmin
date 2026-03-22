package embedfix

import (
	"context"
	"log/slog"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

func (ef *EmbedFix) handleComponent(e *events.ComponentInteractionCreate) {
	customID := e.Data.CustomID()
	_, rest, _ := strings.Cut(customID, ":")
	action, rest, _ := strings.Cut(rest, ":")
	if action != "translate" {
		return
	}

	_ = e.DeferUpdateMessage()

	if ef.translateClient.apiKey == "" {
		_, _ = e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(),
			discord.NewMessageUpdateV2([]discord.LayoutComponent{
				discord.NewContainer(
					discord.NewTextDisplay("翻訳APIキーが設定されていないため、翻訳できません。"),
				),
			}))
		return
	}

	// Parse platform from customID: embedfix:translate:{platform}:{params...}
	platform, rest, _ := strings.Cut(rest, ":")

	ctx := context.Background()

	switch Platform(platform) {
	case PlatformTwitter:
		ef.translateTwitter(ctx, e, rest)
	case PlatformReddit:
		ef.translateReddit(ctx, e, rest)
	case PlatformTikTok:
		ef.translateTikTok(ctx, e, rest)
	case PlatformInstagram:
		ef.translateInstagram(ctx, e, rest)
	default:
		// Backwards compatibility: old format where platform is actually screenName
		ef.translateTwitter(ctx, e, platform+":"+rest)
	}
}

func (ef *EmbedFix) translateTwitter(ctx context.Context, e *events.ComponentInteractionCreate, params string) {
	screenName, tweetID, _ := strings.Cut(params, ":")

	tweet, err := ef.twitterClient.GetTweet(ctx, screenName, tweetID)
	if err != nil {
		ef.logger.Warn("failed to fetch tweet for translation",
			slog.String("tweet_id", tweetID),
			slog.Any("error", err),
		)
		ef.respondTranslateError(e, "ツイートの取得に失敗しました。")
		return
	}

	result, err := ef.translateClient.Translate(ctx, tweet.Text, "ja")
	if err != nil {
		ef.logger.Warn("failed to translate tweet",
			slog.String("tweet_id", tweetID),
			slog.Any("error", err),
		)
		ef.respondTranslateError(e, "翻訳に失敗しました。")
		return
	}

	ref := EmbedRef{Platform: PlatformTwitter, Params: []string{screenName, tweetID}}
	ui := BuildTweetEmbedTranslated(tweet, result, ref)
	_, _ = e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(),
		discord.NewMessageUpdateV2(ui))
}

func (ef *EmbedFix) translateReddit(ctx context.Context, e *events.ComponentInteractionCreate, params string) {
	subreddit, postID, _ := strings.Cut(params, ":")

	post, err := ef.redditClient.GetPost(ctx, subreddit, postID)
	if err != nil {
		ef.respondTranslateError(e, "投稿の取得に失敗しました。")
		return
	}

	text := post.Title
	if post.Selftext != "" {
		text = post.Title + "\n" + post.Selftext
	}

	result, err := ef.translateClient.Translate(ctx, text, "ja")
	if err != nil {
		ef.respondTranslateError(e, "翻訳に失敗しました。")
		return
	}

	ref := EmbedRef{Platform: PlatformReddit, Params: []string{subreddit, postID}}
	ui := BuildRedditEmbedTranslated(post, result, ref)
	_, _ = e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(),
		discord.NewMessageUpdateV2(ui))
}

func (ef *EmbedFix) translateTikTok(ctx context.Context, e *events.ComponentInteractionCreate, params string) {
	username, videoID, _ := strings.Cut(params, ":")

	video, err := ef.tiktokClient.GetVideo(ctx, username, videoID)
	if err != nil {
		ef.respondTranslateError(e, "動画の取得に失敗しました。")
		return
	}

	result, err := ef.translateClient.Translate(ctx, video.Title, "ja")
	if err != nil {
		ef.respondTranslateError(e, "翻訳に失敗しました。")
		return
	}

	ref := EmbedRef{Platform: PlatformTikTok, Params: []string{username, videoID}}
	ui := BuildTikTokEmbedTranslated(video, result, ref)
	_, _ = e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(),
		discord.NewMessageUpdateV2(ui))
}

func (ef *EmbedFix) translateInstagram(ctx context.Context, e *events.ComponentInteractionCreate, params string) {
	shortcode := params

	post, err := ef.instagramClient.GetPost(ctx, shortcode)
	if err != nil {
		ef.respondTranslateError(e, "投稿の取得に失敗しました。")
		return
	}

	if post.Title == "" {
		ef.respondTranslateError(e, "翻訳するテキストがありません。")
		return
	}

	result, err := ef.translateClient.Translate(ctx, post.Title, "ja")
	if err != nil {
		ef.respondTranslateError(e, "翻訳に失敗しました。")
		return
	}

	ref := EmbedRef{Platform: PlatformInstagram, Params: []string{shortcode}}
	ui := BuildInstagramEmbedTranslated(post, result, ref)
	_, _ = e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(),
		discord.NewMessageUpdateV2(ui))
}

func (ef *EmbedFix) respondTranslateError(e *events.ComponentInteractionCreate, msg string) {
	_, _ = e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(),
		discord.NewMessageUpdateV2([]discord.LayoutComponent{
			discord.NewContainer(
				discord.NewTextDisplay(msg),
			),
		}))
}
