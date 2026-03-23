// opyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package handler

import (
	"context"

	disgobot "github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/s12kuma01/pedmin/internal/model"
)

// SetupTranslatorListeners registers the translator's event listeners on the Discord client.
func SetupTranslatorListeners(client *disgobot.Client, h *TranslatorHandler) {
	client.AddEventListeners(
		disgobot.NewListenerFunc(h.onMessageReactionAdd),
	)
}

func (h *TranslatorHandler) onMessageReactionAdd(e *events.GuildMessageReactionAdd) {
	// Ignore bot reactions
	if e.Member.User.Bot {
		return
	}

	if !h.bot.IsModuleEnabled(e.GuildID, model.TranslatorModuleID) {
		return
	}

	if !h.service.IsAvailable() {
		return
	}

	// Check if emoji is a flag
	if e.Emoji.Name == nil {
		return
	}
	targetLang, ok := model.FlagToLang[*e.Emoji.Name]
	if !ok {
		return
	}

	h.service.ProcessTranslation(context.Background(), e.ChannelID, e.MessageID, targetLang)
}
