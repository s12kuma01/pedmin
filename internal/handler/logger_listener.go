// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package handler

import (
	disgobot "github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/Sumire-Labs/pedmin/internal/model"
	"github.com/Sumire-Labs/pedmin/internal/view"
)

// SetupLoggerListeners registers all logger event listeners on the Discord client.
func SetupLoggerListeners(client *disgobot.Client, h *LoggerHandler) {
	client.AddEventListeners(
		disgobot.NewListenerFunc(h.onMessageUpdate),
		disgobot.NewListenerFunc(h.onMessageDelete),
		disgobot.NewListenerFunc(h.onMemberJoin),
		disgobot.NewListenerFunc(h.onMemberLeave),
		disgobot.NewListenerFunc(h.onBan),
		disgobot.NewListenerFunc(h.onUnban),
		disgobot.NewListenerFunc(h.onRoleCreate),
		disgobot.NewListenerFunc(h.onRoleUpdate),
		disgobot.NewListenerFunc(h.onRoleDelete),
		disgobot.NewListenerFunc(h.onChannelCreate),
		disgobot.NewListenerFunc(h.onChannelUpdate),
		disgobot.NewListenerFunc(h.onChannelDelete),
	)
}

func (h *LoggerHandler) onMessageUpdate(e *events.GuildMessageUpdate) {
	if e.Message.Author.Bot {
		return
	}
	oldContent := e.OldMessage.Content
	newContent := e.Message.Content
	if oldContent == newContent && view.LoggerAttachmentsEqual(e.OldMessage.Attachments, e.Message.Attachments) {
		return
	}
	h.sendLog(e.GuildID, model.EventMessageEdit,
		view.LoggerMessageEditLog(e.Message.Author, e.ChannelID, oldContent, newContent, e.OldMessage.Attachments, e.Message.Attachments),
	)
}

func (h *LoggerHandler) onMessageDelete(e *events.GuildMessageDelete) {
	var user *discord.User
	content := e.Message.Content
	attachments := e.Message.Attachments
	forwarded := len(e.Message.MessageSnapshots) > 0

	// For forwarded messages, extract content from the snapshot.
	if forwarded {
		snap := e.Message.MessageSnapshots[0].Message
		if content == "" {
			content = snap.Content
		}
		if len(attachments) == 0 {
			attachments = snap.Attachments
		}
	}

	if e.Message.Author.ID != 0 {
		user = &e.Message.Author
	}
	if user != nil && user.Bot {
		return
	}
	h.sendLog(e.GuildID, model.EventMessageDelete,
		view.LoggerMessageDeleteLog(user, e.ChannelID, content, attachments, forwarded),
	)
}

func (h *LoggerHandler) onMemberJoin(e *events.GuildMemberJoin) {
	h.sendLog(e.GuildID, model.EventMemberJoin,
		view.LoggerMemberJoinLog(e.Member),
	)
}

func (h *LoggerHandler) onMemberLeave(e *events.GuildMemberLeave) {
	h.sendLog(e.GuildID, model.EventMemberLeave,
		view.LoggerMemberLeaveLog(e.User),
	)
}

func (h *LoggerHandler) onBan(e *events.GuildBan) {
	h.sendLog(e.GuildID, model.EventBanAdd,
		view.LoggerBanLog(e.User),
	)
}

func (h *LoggerHandler) onUnban(e *events.GuildUnban) {
	h.sendLog(e.GuildID, model.EventBanRemove,
		view.LoggerUnbanLog(e.User),
	)
}

func (h *LoggerHandler) onRoleCreate(e *events.RoleCreate) {
	h.sendLog(e.GuildID, model.EventRoleChange,
		view.LoggerRoleCreateLog(e.Role),
	)
}

func (h *LoggerHandler) onRoleUpdate(e *events.RoleUpdate) {
	h.sendLog(e.GuildID, model.EventRoleChange,
		view.LoggerRoleUpdateLog(e.Role),
	)
}

func (h *LoggerHandler) onRoleDelete(e *events.RoleDelete) {
	h.sendLog(e.GuildID, model.EventRoleChange,
		view.LoggerRoleDeleteLog(e.Role),
	)
}

func (h *LoggerHandler) onChannelCreate(e *events.GuildChannelCreate) {
	h.sendLog(e.GuildID, model.EventChannelChange,
		view.LoggerChannelCreateLog(e.Channel),
	)
}

func (h *LoggerHandler) onChannelUpdate(e *events.GuildChannelUpdate) {
	h.sendLog(e.GuildID, model.EventChannelChange,
		view.LoggerChannelUpdateLog(e.Channel),
	)
}

func (h *LoggerHandler) onChannelDelete(e *events.GuildChannelDelete) {
	h.sendLog(e.GuildID, model.EventChannelChange,
		view.LoggerChannelDeleteLog(e.Channel),
	)
}
