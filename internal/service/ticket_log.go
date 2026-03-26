// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package service

import (
	"bytes"
	"fmt"
	"html/template"
	"log/slog"
	"slices"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/snowflake/v2"
	"github.com/Sumire-Labs/pedmin/internal/model"
	"github.com/Sumire-Labs/pedmin/internal/view"
)

func (s *TicketService) sendTicketLog(guildID snowflake.ID, ticket *model.Ticket) {
	settings, err := s.LoadSettings(guildID)
	if err != nil || settings.LogChannelID == 0 {
		return
	}

	// Reload to get closed_at/closed_by
	fresh, err := s.store.GetTicketByChannel(ticket.ChannelID)
	if err != nil || fresh == nil {
		fresh = ticket
	}

	msg := view.TicketLog(fresh)
	if _, err := s.client.Rest.CreateMessage(settings.LogChannelID, msg); err != nil {
		s.logger.Error("failed to send ticket log", slog.Any("error", err))
	}
}

func (s *TicketService) sendTranscriptLog(guildID snowflake.ID, ticket *model.Ticket) {
	settings, err := s.LoadSettings(guildID)
	if err != nil || settings.LogChannelID == 0 {
		return
	}

	// Reload to get closed_at/closed_by
	fresh, err := s.store.GetTicketByChannel(ticket.ChannelID)
	if err != nil || fresh == nil {
		fresh = ticket
	}

	file, err := s.generateTranscript(fresh)
	if err != nil {
		s.logger.Error("failed to generate transcript", slog.Any("error", err))
		return
	}

	msg := view.TicketLog(fresh).AddFiles(file)
	if _, err := s.client.Rest.CreateMessage(settings.LogChannelID, msg); err != nil {
		s.logger.Error("failed to send transcript log", slog.Any("error", err))
	}
}

var ticketTranscriptTmpl = template.Must(template.New("transcript").Funcs(template.FuncMap{
	"formatTime": func(t interface{ Format(string) string }) string {
		return t.Format("2006-01-02 15:04:05")
	},
	"isImage": func(contentType *string) bool {
		if contentType == nil {
			return false
		}
		return strings.HasPrefix(*contentType, "image/")
	},
}).Parse(`<!DOCTYPE html>
<html lang="ja">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Ticket #{{printf "%04d" .Ticket.Number}} - {{.Ticket.Subject}}</title>
<style>
* { margin: 0; padding: 0; box-sizing: border-box; }
body { background: #313338; color: #dbdee1; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 15px; line-height: 1.4; }
.header { background: #2b2d31; padding: 20px 24px; border-bottom: 1px solid #1e1f22; }
.header h1 { font-size: 20px; color: #f2f3f5; margin-bottom: 8px; }
.header .meta { font-size: 13px; color: #949ba4; }
.header .meta span { margin-right: 16px; }
.messages { padding: 16px 0; }
.message { padding: 4px 24px; display: flex; gap: 12px; }
.message:hover { background: #2e3035; }
.message.has-header { margin-top: 12px; }
.avatar { width: 40px; height: 40px; border-radius: 50%; flex-shrink: 0; }
.avatar-spacer { width: 40px; flex-shrink: 0; }
.content { min-width: 0; flex: 1; }
.msg-header { display: flex; align-items: baseline; gap: 8px; margin-bottom: 2px; }
.username { font-weight: 600; color: #f2f3f5; font-size: 15px; }
.bot-tag { background: #5865f2; color: #fff; font-size: 10px; padding: 1px 4px; border-radius: 3px; font-weight: 600; vertical-align: middle; }
.timestamp { font-size: 12px; color: #949ba4; }
.text { white-space: pre-wrap; word-wrap: break-word; color: #dbdee1; }
.attachments { margin-top: 4px; }
.attachments img { max-width: 400px; max-height: 300px; border-radius: 8px; margin-top: 4px; display: block; }
.attachments a { color: #00a8fc; text-decoration: none; font-size: 14px; }
.attachments a:hover { text-decoration: underline; }
.footer { background: #2b2d31; padding: 16px 24px; border-top: 1px solid #1e1f22; font-size: 13px; color: #949ba4; text-align: center; }
</style>
</head>
<body>
<div class="header">
	<h1>Ticket #{{printf "%04d" .Ticket.Number}} — {{.Ticket.Subject}}</h1>
	<div class="meta">
		<span>作成者: {{.Ticket.UserID}}</span>
		<span>作成日時: {{formatTime .Ticket.CreatedAt}}</span>
		{{- if .Ticket.ClosedAt}}
		<span>クローズ日時: {{formatTime .Ticket.ClosedAt}}</span>
		{{- end}}
	</div>
</div>
<div class="messages">
{{- $prev := "" -}}
{{- range .Messages -}}
{{- $author := .Author.Username -}}
{{- $showHeader := false -}}
{{- if ne $author $prev}}{{$showHeader = true}}{{end -}}
<div class="message{{if $showHeader}} has-header{{end}}">
	{{- if $showHeader}}
	<img class="avatar" src="{{.Author.EffectiveAvatarURL}}" alt="">
	{{- else}}
	<div class="avatar-spacer"></div>
	{{- end}}
	<div class="content">
		{{- if $showHeader}}
		<div class="msg-header">
			<span class="username">{{.Author.EffectiveName}}</span>
			{{- if .Author.Bot}}<span class="bot-tag">BOT</span>{{end}}
			<span class="timestamp">{{formatTime .CreatedAt}}</span>
		</div>
		{{- end}}
		{{- if .Content}}
		<div class="text">{{.Content}}</div>
		{{- end}}
		{{- if .Attachments}}
		<div class="attachments">
			{{- range .Attachments}}
			{{- if isImage .ContentType}}
			<a href="{{.URL}}" target="_blank"><img src="{{.URL}}" alt="{{.Filename}}"></a>
			{{- else}}
			<a href="{{.URL}}" target="_blank">📎 {{.Filename}}</a>
			{{- end}}
			{{- end}}
		</div>
		{{- end}}
	</div>
</div>
{{$prev = $author}}
{{- end}}
</div>
<div class="footer">
	{{len .Messages}} メッセージ
</div>
</body>
</html>
`))

type ticketTranscriptData struct {
	Ticket   *model.Ticket
	Messages []discord.Message
}

func (s *TicketService) fetchAllMessages(channelID snowflake.ID) ([]discord.Message, error) {
	var all []discord.Message
	page := s.client.Rest.GetMessagesPage(channelID, 0, 100)
	for page.Previous() {
		all = append(all, page.Items...)
	}
	if page.Err != nil && page.Err != rest.ErrNoMorePages {
		return nil, page.Err
	}
	// Previous() fetches newest-first then older pages, so reverse for chronological order
	slices.Reverse(all)
	return all, nil
}

func generateTranscriptHTML(ticket *model.Ticket, messages []discord.Message) ([]byte, error) {
	var buf bytes.Buffer
	if err := ticketTranscriptTmpl.Execute(&buf, ticketTranscriptData{
		Ticket:   ticket,
		Messages: messages,
	}); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *TicketService) generateTranscript(ticket *model.Ticket) (*discord.File, error) {
	messages, err := s.fetchAllMessages(ticket.ChannelID)
	if err != nil {
		return nil, err
	}
	html, err := generateTranscriptHTML(ticket, messages)
	if err != nil {
		return nil, err
	}
	filename := fmt.Sprintf("ticket-%04d.html", ticket.Number)
	return &discord.File{
		Name:   filename,
		Reader: bytes.NewReader(html),
	}, nil
}
