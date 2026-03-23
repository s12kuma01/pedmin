package view

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/s12kuma01/pedmin/internal/model"
	"github.com/s12kuma01/pedmin/pkg/deepl"
)

// BuildRedditEmbed builds a Reddit post embed message.
func BuildRedditEmbed(post *model.RedditPost, ref model.EmbedRef) discord.MessageCreate {
	components := BuildRedditComponents(post, ref, "", "")
	return discord.NewMessageCreateV2(discord.NewContainer(components...))
}

// BuildRedditEmbedTranslated builds a translated Reddit post embed as layout components.
func BuildRedditEmbedTranslated(post *model.RedditPost, result *deepl.TranslateResult, ref model.EmbedRef) []discord.LayoutComponent {
	footer := fmt.Sprintf("%s | <t:%d:f> · %sから翻訳", emojiReddit, post.CreatedUTC.Unix(), deepl.LangName(result.DetectedLanguage))
	components := BuildRedditComponents(post, ref, result.TranslatedText, footer)
	return []discord.LayoutComponent{discord.NewContainer(components...)}
}

// BuildRedditComponents builds the Reddit post embed sub-components.
func BuildRedditComponents(post *model.RedditPost, ref model.EmbedRef, translatedText, footerOverride string) []discord.ContainerSubComponent {
	headerText := fmt.Sprintf("**r/%s** · u/%s", post.Subreddit, post.Author)

	var headerComponent discord.ContainerSubComponent
	if post.SubredditIcon != "" {
		headerComponent = discord.NewSection(
			discord.NewTextDisplay(headerText),
		).WithAccessory(discord.NewThumbnail(post.SubredditIcon))
	} else {
		headerComponent = discord.NewTextDisplay(headerText)
	}

	components := []discord.ContainerSubComponent{
		headerComponent,
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

	stats := fmt.Sprintf("%s %s  %s %s",
		emojiUpvote, FormatCount(post.Score),
		emojiMessages, FormatCount(post.NumComments),
	)
	components = append(components, discord.NewTextDisplay(stats))

	footer := fmt.Sprintf("%s | <t:%d:f>", emojiReddit, post.CreatedUTC.Unix())
	if footerOverride != "" {
		footer = footerOverride
	}
	components = append(components, discord.NewTextDisplay(footer))

	// Show translate button for posts with text content
	if post.Selftext != "" && footerOverride == "" {
		customID := fmt.Sprintf("%s:translate:%s:%s:%s", model.EmbedFixModuleID, model.PlatformReddit, ref.Params[0], ref.Params[1])
		components = append(components,
			discord.NewActionRow(
				discord.NewSecondaryButton("\U0001f310 翻訳", customID),
			),
		)
	}

	return components
}
