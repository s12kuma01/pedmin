package embedfix

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
)

func BuildTikTokEmbed(video *TikTokVideo, ref EmbedRef) discord.MessageCreate {
	components := buildTikTokComponents(video, ref, "", "")
	return discord.NewMessageCreateV2(discord.NewContainer(components...))
}

func BuildTikTokEmbedTranslated(video *TikTokVideo, result *TranslateResult, ref EmbedRef) []discord.LayoutComponent {
	footer := fmt.Sprintf("🎵 | <t:%d:f> · %sから翻訳", video.CreatedAt.Unix(), langName(result.DetectedLanguage))
	components := buildTikTokComponents(video, ref, result.TranslatedText, footer)
	return []discord.LayoutComponent{discord.NewContainer(components...)}
}

func buildTikTokComponents(video *TikTokVideo, ref EmbedRef, translatedText, footerOverride string) []discord.ContainerSubComponent {
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

	if video.CoverURL != "" {
		components = append(components, discord.NewMediaGallery(
			discord.MediaGalleryItem{
				Media: discord.UnfurledMediaItem{URL: video.CoverURL},
			},
		))
	}

	components = append(components, discord.NewSmallSeparator())

	stats := fmt.Sprintf("▶ %s  ❤ %s  💬 %s  🔗 %s",
		formatCount(video.PlayCount),
		formatCount(video.LikeCount),
		formatCount(video.CommentCount),
		formatCount(video.ShareCount),
	)
	components = append(components, discord.NewTextDisplay(stats))

	footer := fmt.Sprintf("🎵 | <t:%d:f>", video.CreatedAt.Unix())
	if footerOverride != "" {
		footer = footerOverride
	}
	components = append(components, discord.NewTextDisplay(footer))

	// Show translate button for videos with text
	if video.Title != "" && footerOverride == "" {
		customID := fmt.Sprintf("%s:translate:%s:%s:%s", ModuleID, PlatformTikTok, ref.Params[0], ref.Params[1])
		components = append(components,
			discord.NewActionRow(
				discord.NewSecondaryButton("🌐 翻訳", customID),
			),
		)
	}

	return components
}
