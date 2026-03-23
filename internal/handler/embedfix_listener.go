package handler

import (
	"context"

	"github.com/disgoorg/disgo/events"
	"github.com/s12kuma01/pedmin/internal/model"
)

func (h *EmbedFixHandler) onMessageCreate(e *events.GuildMessageCreate) {
	if e.Message.Author.Bot {
		return
	}
	if !h.bot.IsModuleEnabled(e.GuildID, model.EmbedFixModuleID) {
		return
	}

	h.service.ProcessMessageURLs(context.Background(), e.GuildID, e.ChannelID, e.MessageID, e.Message.Content)
}
