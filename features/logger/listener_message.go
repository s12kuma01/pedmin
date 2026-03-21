package logger

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

func (l *Logger) onMessageUpdate(e *events.GuildMessageUpdate) {
	if e.Message.Author.Bot {
		return
	}
	oldContent := e.OldMessage.Content
	newContent := e.Message.Content
	if oldContent == newContent && AttachmentsEqual(e.OldMessage.Attachments, e.Message.Attachments) {
		return
	}
	l.sendLog(e.GuildID, EventMessageEdit,
		BuildMessageEditLog(e.Message.Author, e.ChannelID, oldContent, newContent, e.OldMessage.Attachments, e.Message.Attachments),
	)
}

func (l *Logger) onMessageDelete(e *events.GuildMessageDelete) {
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
	l.sendLog(e.GuildID, EventMessageDelete,
		BuildMessageDeleteLog(user, e.ChannelID, content, attachments, forwarded),
	)
}
