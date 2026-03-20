package panel

import (
	"context"
	"log/slog"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

func (p *Panel) HandleCommand(e *events.ApplicationCommandInteractionCreate) {
	userID := e.User().ID

	if !p.isAllowed(userID) {
		_ = e.CreateMessage(ephemeralError("このコマンドを使用する権限がありません。"))
		return
	}

	if p.cfg.PanelURL == "" || p.cfg.PanelAPIKey == "" {
		_ = e.CreateMessage(ephemeralError("パネルが設定されていません。"))
		return
	}

	_ = e.DeferCreateMessage(false)

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	servers, err := p.ListServersWithStatus(ctx)
	if err != nil {
		p.logger.Error("failed to list servers", slog.Any("error", err))
		_, _ = e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(), BuildErrorPanel(err.Error()))
		return
	}

	msg := BuildServerList(servers)
	_, _ = e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(), discord.NewMessageUpdateV2(msg.Components))
}

func ephemeralError(text string) discord.MessageCreate {
	return discord.NewMessageCreateV2(
		discord.NewContainer(
			discord.NewTextDisplay(text),
		),
	).WithEphemeral(true)
}
