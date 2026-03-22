package embedfix

import (
	"context"
	"fmt"
	"strings"

	"github.com/disgoorg/disgo/discord"
)

// translateContent fetches content for the given platform and translates it.
func (ef *EmbedFix) translateContent(ctx context.Context, platform, params string) ([]discord.LayoutComponent, error) {
	switch Platform(platform) {
	case PlatformTwitter:
		return ef.translateTwitterContent(ctx, params)
	case PlatformReddit:
		return ef.translateRedditContent(ctx, params)
	case PlatformTikTok:
		return ef.translateTikTokContent(ctx, params)
	default:
		// Backwards compatibility: old format where platform is actually screenName
		return ef.translateTwitterContent(ctx, platform+":"+params)
	}
}

func (ef *EmbedFix) translateTwitterContent(ctx context.Context, params string) ([]discord.LayoutComponent, error) {
	screenName, tweetID, _ := strings.Cut(params, ":")

	tweet, err := ef.twitterClient.GetTweet(ctx, screenName, tweetID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tweet: %w", err)
	}

	result, err := ef.translateClient.Translate(ctx, tweet.Text, "ja")
	if err != nil {
		return nil, fmt.Errorf("failed to translate: %w", err)
	}

	ref := EmbedRef{Platform: PlatformTwitter, Params: []string{screenName, tweetID}}
	return BuildTweetEmbedTranslated(tweet, result, ref), nil
}

func (ef *EmbedFix) translateRedditContent(ctx context.Context, params string) ([]discord.LayoutComponent, error) {
	subreddit, postID, _ := strings.Cut(params, ":")

	post, err := ef.redditClient.GetPost(ctx, subreddit, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch reddit post: %w", err)
	}

	text := post.Title
	if post.Selftext != "" {
		text = post.Title + "\n" + post.Selftext
	}

	result, err := ef.translateClient.Translate(ctx, text, "ja")
	if err != nil {
		return nil, fmt.Errorf("failed to translate: %w", err)
	}

	ref := EmbedRef{Platform: PlatformReddit, Params: []string{subreddit, postID}}
	return BuildRedditEmbedTranslated(post, result, ref), nil
}

func (ef *EmbedFix) translateTikTokContent(ctx context.Context, params string) ([]discord.LayoutComponent, error) {
	username, videoID, _ := strings.Cut(params, ":")

	video, err := ef.tiktokClient.GetVideo(ctx, username, videoID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tiktok video: %w", err)
	}

	result, err := ef.translateClient.Translate(ctx, video.Title, "ja")
	if err != nil {
		return nil, fmt.Errorf("failed to translate: %w", err)
	}

	ref := EmbedRef{Platform: PlatformTikTok, Params: []string{username, videoID}}
	return BuildTikTokEmbedTranslated(video, result, ref), nil
}
