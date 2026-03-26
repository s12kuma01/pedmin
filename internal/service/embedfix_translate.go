// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/Sumire-Labs/pedmin/internal/model"
	"github.com/Sumire-Labs/pedmin/internal/view"
)

// TranslateContent fetches content for the given platform and translates it.
func (s *EmbedFixService) TranslateContent(ctx context.Context, platform, params string) ([]discord.LayoutComponent, error) {
	switch model.Platform(platform) {
	case model.PlatformTwitter:
		return s.translateTwitterContent(ctx, params)
	case model.PlatformReddit:
		return s.translateRedditContent(ctx, params)
	case model.PlatformTikTok:
		return s.translateTikTokContent(ctx, params)
	default:
		// Backwards compatibility: old format where platform is actually screenName
		return s.translateTwitterContent(ctx, platform+":"+params)
	}
}

// IsTranslationAvailable reports whether the translation API is configured.
func (s *EmbedFixService) IsTranslationAvailable() bool {
	return s.translateClient.IsAvailable()
}

func (s *EmbedFixService) translateTwitterContent(ctx context.Context, params string) ([]discord.LayoutComponent, error) {
	screenName, tweetID, _ := strings.Cut(params, ":")

	tweet, err := s.twitterClient.GetTweet(ctx, screenName, tweetID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tweet: %w", err)
	}

	result, err := s.translateClient.Translate(ctx, tweet.Text, "ja")
	if err != nil {
		return nil, fmt.Errorf("failed to translate: %w", err)
	}

	ref := model.EmbedRef{Platform: model.PlatformTwitter, Params: []string{screenName, tweetID}}
	return view.BuildTweetEmbedTranslated(tweet, result, ref), nil
}

func (s *EmbedFixService) translateRedditContent(ctx context.Context, params string) ([]discord.LayoutComponent, error) {
	subreddit, postID, _ := strings.Cut(params, ":")

	post, err := s.redditClient.GetPost(ctx, subreddit, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch reddit post: %w", err)
	}

	text := post.Title
	if post.Selftext != "" {
		text = post.Title + "\n" + post.Selftext
	}

	result, err := s.translateClient.Translate(ctx, text, "ja")
	if err != nil {
		return nil, fmt.Errorf("failed to translate: %w", err)
	}

	ref := model.EmbedRef{Platform: model.PlatformReddit, Params: []string{subreddit, postID}}
	return view.BuildRedditEmbedTranslated(post, result, ref), nil
}

func (s *EmbedFixService) translateTikTokContent(ctx context.Context, params string) ([]discord.LayoutComponent, error) {
	username, videoID, _ := strings.Cut(params, ":")

	video, err := s.tiktokClient.GetVideo(ctx, username, videoID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tiktok video: %w", err)
	}

	result, err := s.translateClient.Translate(ctx, video.Title, "ja")
	if err != nil {
		return nil, fmt.Errorf("failed to translate: %w", err)
	}

	ref := model.EmbedRef{Platform: model.PlatformTikTok, Params: []string{username, videoID}}
	return view.BuildTikTokEmbedTranslated(video, result, ref), nil
}
