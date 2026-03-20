package panel

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

const (
	defaultTimeout     = 10 * time.Second
	powerActionTimeout = 15 * time.Second
)

func (p *Panel) HandleComponent(e *events.ComponentInteractionCreate) {
	if !p.isAllowed(e.User().ID) {
		_ = e.CreateMessage(ephemeralError("このコマンドを使用する権限がありません。"))
		return
	}

	customID := e.Data.CustomID()
	_, rest, _ := strings.Cut(customID, ":")
	action, extra, _ := strings.Cut(rest, ":")

	switch action {
	case "select":
		p.handleSelect(e)
	case "power_start":
		p.handlePower(e, extra, "start")
	case "power_restart":
		p.handlePower(e, extra, "restart")
	case "power_stop":
		p.handlePower(e, extra, "stop")
	case "refresh":
		p.handleRefresh(e, extra)
	case "back":
		p.handleBack(e)
	case "refresh_list":
		p.handleBack(e)
	case "console":
		p.handleConsolePrompt(e, extra)
	}
}

func (p *Panel) handleSelect(e *events.ComponentInteractionCreate) {
	data, ok := e.Data.(discord.StringSelectMenuInteractionData)
	if !ok || len(data.Values) == 0 {
		return
	}
	identifier := data.Values[0]

	_ = e.DeferUpdateMessage()

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	server, res, err := p.GetServerDetail(ctx, identifier)
	if err != nil {
		p.logger.Error("failed to get server detail", slog.Any("error", err))
		_, _ = e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(), BuildErrorPanel(err.Error()))
		return
	}

	_, _ = e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(), BuildServerDetail(*server, res))
}

func (p *Panel) handlePower(e *events.ComponentInteractionCreate, identifier, signal string) {
	_ = e.DeferUpdateMessage()

	ctx, cancel := context.WithTimeout(context.Background(), powerActionTimeout)
	defer cancel()

	server, res, err := p.PowerAction(ctx, identifier, signal)
	if err != nil {
		p.logger.Error("failed to send power action", slog.String("signal", signal), slog.Any("error", err))
		_, _ = e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(), BuildErrorPanel(err.Error()))
		return
	}

	_, _ = e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(), BuildServerDetail(*server, res))
}

func (p *Panel) handleRefresh(e *events.ComponentInteractionCreate, identifier string) {
	_ = e.DeferUpdateMessage()

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	server, res, err := p.GetServerDetail(ctx, identifier)
	if err != nil {
		p.logger.Error("failed to refresh server detail", slog.Any("error", err))
		_, _ = e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(), BuildErrorPanel(err.Error()))
		return
	}

	_, _ = e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(), BuildServerDetail(*server, res))
}

func (p *Panel) handleBack(e *events.ComponentInteractionCreate) {
	_ = e.DeferUpdateMessage()

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

func (p *Panel) handleConsolePrompt(e *events.ComponentInteractionCreate, identifier string) {
	_ = e.Modal(BuildConsoleModal(identifier))
}

func (p *Panel) HandleModal(e *events.ModalSubmitInteractionCreate) {
	if !p.isAllowed(e.User().ID) {
		_ = e.CreateMessage(ephemeralError("このコマンドを使用する権限がありません。"))
		return
	}

	customID := e.Data.CustomID
	_, rest, _ := strings.Cut(customID, ":")
	action, identifier, _ := strings.Cut(rest, ":")

	if action != "console_modal" {
		return
	}

	command := strings.TrimSpace(e.Data.Text(ModuleID + ":cmd"))
	if command == "" {
		_ = e.CreateMessage(ephemeralError("コマンドを入力してください。"))
		return
	}

	_ = e.DeferCreateMessage(true)

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	if err := p.SendConsoleCommand(ctx, identifier, command); err != nil {
		p.logger.Error("failed to send command", slog.Any("error", err))
		_, _ = e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(), BuildConsoleError(err.Error()))
		return
	}

	_, _ = e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(), BuildConsoleResult(command))
}
