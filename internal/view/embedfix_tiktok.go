// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package view

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/Sumire-Labs/pedmin/internal/model"
	"github.com/Sumire-Labs/pedmin/pkg/deepl"
)

// BuildTikTokEmbed builds a TikTok video embed message.
func BuildTikTokEmbed(video *model.TikTokVideo, ref model.EmbedRef) discord.MessageCreate {
	components := BuildTikTokComponents(video, ref, "", "")
	return discord.NewMessageCreateV2(discord.NewContainer(components...))
}

// BuildTikTokEmbedTranslated builds a translated TikTok video embed as layout components.
func BuildTikTokEmbedTranslated(video *model.TikTokVideo, result *deepl.TranslateResult, ref model.EmbedRef) []discord.LayoutComponent {
	footer := fmt.Sprintf("%s | <t:%d:R> · %sから翻訳", emojiTikTok, video.CreatedAt.Unix(), deepl.LangName(result.DetectedLanguage))
	components := BuildTikTokComponents(video, ref, result.TranslatedText, footer)
	return []discord.LayoutComponent{discord.NewContainer(components...)}
}

// BuildTikTokComponents builds the TikTok video embed sub-components.
func BuildTikTokComponents(video *model.TikTokVideo, ref model.EmbedRef, translatedText, footerOverride string) []discord.ContainerSubComponent {
	authorLine := fmt.Sprintf("**%s** @%s", video.Author.Nickname, video.Author.UniqueID)

	components := []discord.ContainerSubComponent{
		discord.NewSection(
			discord.NewTextDisplay(authorLine),
		).WithAccessory(discord.NewThumbnail(video.Author.Avatar)),
		discord.NewSmallSeparator(),
	}

	if translatedText != "" {
		components = append(components, discord.NewTextDisplay(translatedText))
	} else if video.Title != "" {
		components = append(components, discord.NewTextDisplay(video.Title))
	}

	// Try video URL first (Discord may render it inline), fallback to cover image
	mediaURL := video.CoverURL
	if video.VideoURL != "" {
		mediaURL = video.VideoURL
	}
	if mediaURL != "" {
		components = append(components, discord.NewMediaGallery(
			discord.MediaGalleryItem{
				Media: discord.UnfurledMediaItem{URL: mediaURL},
			},
		))
	}

	components = append(components, discord.NewSmallSeparator())

	stats := fmt.Sprintf("%s %s  %s %s  %s %s  %s %s",
		emojiPlay, FormatCount(video.PlayCount),
		emojiLike, FormatCount(video.LikeCount),
		emojiMessages, FormatCount(video.CommentCount),
		emojiShare, FormatCount(video.ShareCount),
	)
	components = append(components, discord.NewTextDisplay(stats))

	footer := fmt.Sprintf("%s | <t:%d:R>", emojiTikTok, video.CreatedAt.Unix())
	if footerOverride != "" {
		footer = footerOverride
	}
	components = append(components, discord.NewTextDisplay(footer))

	// Show translate button for videos with text
	if video.Title != "" && footerOverride == "" {
		customID := fmt.Sprintf("%s:translate:%s:%s:%s", model.EmbedFixModuleID, model.PlatformTikTok, ref.Params[0], ref.Params[1])
		components = append(components,
			discord.NewActionRow(
				discord.NewSecondaryButton("\U0001f310 翻訳", customID),
			),
		)
	}

	return components
}
