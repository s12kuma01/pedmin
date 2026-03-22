package panel

import (
	"context"
	"log/slog"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/s12kuma01/pedmin/ui"
)

func (p *Panel) HandleCommand(e *events.ApplicationCommandInteractionCreate) {
	userID := e.User().ID

	if !p.isAllowed(userID) {
		_ = e.CreateMessage(ui.ErrorMessage("このコマンドを使用する権限がありません。"))
		return
	}

	if p.cfg.PanelURL == "" || p.cfg.PanelAPIKey == "" {
		_ = e.CreateMessage(ui.ErrorMessage("パネルが設定されていません。"))
		return
	}

	_ = e.DeferCreateMessage(false)

	ctx, cancel := context.WithTimeout(context.Background(), p.cfg.HTTPClientTimeout)
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
