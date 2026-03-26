// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package view

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
	"github.com/Sumire-Labs/pedmin/pkg/deepl"
)

// BuildTranslationEmbed builds the translation result message.
func BuildTranslationEmbed(translatedText string, sourceLang, targetLang string, authorID snowflake.ID, messageID snowflake.ID) discord.MessageCreate {
	components := []discord.ContainerSubComponent{
		discord.NewTextDisplay(translatedText),
		discord.NewSmallSeparator(),
		discord.NewTextDisplay(fmt.Sprintf(
			"-# 翻訳元: <@%d>\n-# %s → %s",
			authorID, deepl.LangName(sourceLang), deepl.LangName(targetLang),
		)),
	}

	return discord.NewMessageCreateV2(
		discord.NewContainer(components...),
	).WithMessageReferenceByID(messageID).WithAllowedMentions(&discord.AllowedMentions{})
}
