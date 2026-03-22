package embedfix

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
)

func BuildInstagramEmbed(post *InstagramPost, ref EmbedRef) discord.MessageCreate {
	components := buildInstagramComponents(post, ref, "", "")
	return discord.NewMessageCreateV2(discord.NewContainer(components...))
}

func BuildInstagramEmbedTranslated(post *InstagramPost, result *TranslateResult, ref EmbedRef) []discord.LayoutComponent {
	footer := fmt.Sprintf("📷 | <t:%d:f> · %sから翻訳", post.CreatedAt.Unix(), langName(result.DetectedLanguage))
	components := buildInstagramComponents(post, ref, result.TranslatedText, footer)
	return []discord.LayoutComponent{discord.NewContainer(components...)}
}

func buildInstagramComponents(post *InstagramPost, ref EmbedRef, translatedText, footerOverride string) []discord.ContainerSubComponent {
	components := []discord.ContainerSubComponent{
		discord.NewTextDisplay(fmt.Sprintf("**%s**", post.AuthorName)),
		discord.NewSmallSeparator(),
	}

	if translatedText != "" {
		components = append(components, discord.NewTextDisplay(translatedText))
	} else if post.Title != "" {
		components = append(components, discord.NewTextDisplay(post.Title))
	}

	if post.ThumbnailURL != "" {
		components = append(components, discord.NewMediaGallery(
			discord.MediaGalleryItem{
				Media: discord.UnfurledMediaItem{URL: post.ThumbnailURL},
			},
		))
	}

	components = append(components, discord.NewSmallSeparator())

	footer := fmt.Sprintf("📷 | <t:%d:f>", post.CreatedAt.Unix())
	if footerOverride != "" {
		footer = footerOverride
	}
	components = append(components, discord.NewTextDisplay(footer))

	// Show translate button for posts with text
	if post.Title != "" && footerOverride == "" {
		customID := fmt.Sprintf("%s:translate:%s:%s", ModuleID, PlatformInstagram, ref.Params[0])
		components = append(components,
			discord.NewActionRow(
				discord.NewSecondaryButton("🌐 翻訳", customID),
			),
		)
	}

	return components
}
