package embedfix

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
)

func BuildRedditEmbed(post *RedditPost, ref EmbedRef) discord.MessageCreate {
	components := buildRedditComponents(post, ref, "", "")
	return discord.NewMessageCreateV2(discord.NewContainer(components...))
}

func BuildRedditEmbedTranslated(post *RedditPost, result *TranslateResult, ref EmbedRef) []discord.LayoutComponent {
	footer := fmt.Sprintf("🔗 | <t:%d:f> · %sから翻訳", post.CreatedUTC.Unix(), langName(result.DetectedLanguage))
	components := buildRedditComponents(post, ref, result.TranslatedText, footer)
	return []discord.LayoutComponent{discord.NewContainer(components...)}
}

func buildRedditComponents(post *RedditPost, ref EmbedRef, translatedText, footerOverride string) []discord.ContainerSubComponent {
	header := fmt.Sprintf("**r/%s** · u/%s", post.Subreddit, post.Author)

	components := []discord.ContainerSubComponent{
		discord.NewTextDisplay(header),
		discord.NewSmallSeparator(),
	}

	if translatedText != "" {
		components = append(components, discord.NewTextDisplay(translatedText))
	} else {
		components = append(components, discord.NewTextDisplay(fmt.Sprintf("**%s**", post.Title)))
		if post.Selftext != "" {
			components = append(components, discord.NewTextDisplay(post.Selftext))
		}
	}

	// Preview images
	if len(post.Preview) > 0 {
		items := make([]discord.MediaGalleryItem, 0, len(post.Preview))
		for _, imgURL := range post.Preview {
			items = append(items, discord.MediaGalleryItem{
				Media: discord.UnfurledMediaItem{URL: imgURL},
			})
		}
		components = append(components, discord.NewMediaGallery(items...))
	}

	components = append(components, discord.NewSmallSeparator())

	stats := fmt.Sprintf("⬆ %s  💬 %s", formatCount(post.Score), formatCount(post.NumComments))
	components = append(components, discord.NewTextDisplay(stats))

	footer := fmt.Sprintf("🔗 | <t:%d:f>", post.CreatedUTC.Unix())
	if footerOverride != "" {
		footer = footerOverride
	}
	components = append(components, discord.NewTextDisplay(footer))

	// Show translate button for posts with text content
	if post.Selftext != "" && footerOverride == "" {
		customID := fmt.Sprintf("%s:translate:%s:%s:%s", ModuleID, PlatformReddit, ref.Params[0], ref.Params[1])
		components = append(components,
			discord.NewActionRow(
				discord.NewSecondaryButton("🌐 翻訳", customID),
			),
		)
	}

	return components
}
