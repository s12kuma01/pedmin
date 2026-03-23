package view

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/s12kuma01/pedmin/internal/model"
	"github.com/s12kuma01/pedmin/pkg/deepl"
)

// BuildTweetEmbed builds a tweet embed message.
func BuildTweetEmbed(tweet *model.Tweet, ref model.EmbedRef) discord.MessageCreate {
	components := BuildTweetComponents(tweet, ref, tweet.Text, "")
	return discord.NewMessageCreateV2(discord.NewContainer(components...))
}

// BuildTweetEmbedTranslated builds a translated tweet embed as layout components.
func BuildTweetEmbedTranslated(tweet *model.Tweet, result *deepl.TranslateResult, ref model.EmbedRef) []discord.LayoutComponent {
	footer := fmt.Sprintf("%s | <t:%d:f> · %sから翻訳", emojiX, tweet.CreatedAt.Unix(), deepl.LangName(result.DetectedLanguage))
	components := BuildTweetComponents(tweet, ref, result.TranslatedText, footer)
	return []discord.LayoutComponent{discord.NewContainer(components...)}
}

// BuildTweetComponents builds the tweet embed sub-components.
func BuildTweetComponents(tweet *model.Tweet, ref model.EmbedRef, text, footerOverride string) []discord.ContainerSubComponent {
	components := []discord.ContainerSubComponent{
		discord.NewSection(
			discord.NewTextDisplay(fmt.Sprintf("**%s** [@%s](https://x.com/%s)", tweet.Author.Name, tweet.Author.ScreenName, tweet.Author.ScreenName)),
		).WithAccessory(discord.NewThumbnail(tweet.Author.AvatarURL)),
		discord.NewSmallSeparator(),
		discord.NewTextDisplay(text),
	}

	if len(tweet.Media) > 0 {
		items := make([]discord.MediaGalleryItem, 0, len(tweet.Media))
		for _, m := range tweet.Media {
			mediaURL := m.URL
			if m.Type == "video" && m.ThumbnailURL != "" {
				mediaURL = m.ThumbnailURL
			}
			items = append(items, discord.MediaGalleryItem{
				Media: discord.UnfurledMediaItem{URL: mediaURL},
			})
		}
		components = append(components, discord.NewMediaGallery(items...))
	}

	components = append(components, discord.NewSmallSeparator())

	stats := fmt.Sprintf("%s %s  %s %s  %s %s  %s %s",
		emojiMessages, FormatCount(tweet.Replies),
		emojiRepost, FormatCount(tweet.Retweets),
		emojiLike, FormatCount(tweet.Likes),
		emojiGraph, FormatCount(tweet.Views),
	)
	components = append(components, discord.NewTextDisplay(stats))

	footer := fmt.Sprintf("%s | <t:%d:f>", emojiX, tweet.CreatedAt.Unix())
	if footerOverride != "" {
		footer = footerOverride
	}
	components = append(components, discord.NewTextDisplay(footer))

	// Show translate button only for non-Japanese tweets and when not already translated
	if tweet.Lang != "ja" && footerOverride == "" {
		customID := fmt.Sprintf("%s:translate:%s:%s:%s", model.EmbedFixModuleID, model.PlatformTwitter, ref.Params[0], ref.Params[1])
		components = append(components,
			discord.NewActionRow(
				discord.NewSecondaryButton("\U0001f310 翻訳", customID),
			),
		)
	}

	return components
}
