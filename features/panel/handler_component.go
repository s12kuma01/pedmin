package panel

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
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
	data := e.Data.(discord.StringSelectMenuInteractionData)
	if len(data.Values) == 0 {
		return
	}
	identifier := data.Values[0]

	_ = e.DeferUpdateMessage()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	server, err := p.findServer(ctx, identifier)
	if err != nil {
		p.logger.Error("failed to find server", slog.Any("error", err))
		_, _ = e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(), BuildErrorPanel(err.Error()))
		return
	}

	res, err := p.pelican.GetResources(ctx, identifier)
	if err != nil {
		p.logger.Error("failed to get resources", slog.Any("error", err))
		_, _ = e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(), BuildErrorPanel(err.Error()))
		return
	}

	_, _ = e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(), BuildServerDetail(*server, res))
}

func (p *Panel) handlePower(e *events.ComponentInteractionCreate, identifier, signal string) {
	_ = e.DeferUpdateMessage()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := p.pelican.SendPowerAction(ctx, identifier, signal); err != nil {
		p.logger.Error("failed to send power action", slog.String("signal", signal), slog.Any("error", err))
		_, _ = e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(), BuildErrorPanel(err.Error()))
		return
	}

	// Wait for state transition
	time.Sleep(2 * time.Second)

	p.refreshDetail(e, identifier)
}

func (p *Panel) handleRefresh(e *events.ComponentInteractionCreate, identifier string) {
	_ = e.DeferUpdateMessage()
	p.refreshDetail(e, identifier)
}

func (p *Panel) handleBack(e *events.ComponentInteractionCreate) {
	_ = e.DeferUpdateMessage()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	servers, err := p.pelican.ListServers(ctx)
	if err != nil {
		p.logger.Error("failed to list servers", slog.Any("error", err))
		_, _ = e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(), BuildErrorPanel(err.Error()))
		return
	}

	msg := BuildServerList(servers)
	_, _ = e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(), discord.NewMessageUpdateV2(msg.Components))
}

func (p *Panel) handleConsolePrompt(e *events.ComponentInteractionCreate, identifier string) {
	_ = e.Modal(discord.ModalCreate{
		CustomID: ModuleID + ":console_modal:" + identifier,
		Title:    "コンソールコマンド",
		Components: []discord.LayoutComponent{
			discord.NewLabel("コマンド",
				discord.NewShortTextInput(ModuleID+":cmd").
					WithRequired(true).
					WithPlaceholder("say hello"),
			),
		},
	})
}

func (p *Panel) refreshDetail(e *events.ComponentInteractionCreate, identifier string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	server, err := p.findServer(ctx, identifier)
	if err != nil {
		p.logger.Error("failed to find server", slog.Any("error", err))
		_, _ = e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(), BuildErrorPanel(err.Error()))
		return
	}

	res, err := p.pelican.GetResources(ctx, identifier)
	if err != nil {
		p.logger.Error("failed to get resources", slog.Any("error", err))
		_, _ = e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(), BuildErrorPanel(err.Error()))
		return
	}

	_, _ = e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(), BuildServerDetail(*server, res))
}

func (p *Panel) findServer(ctx context.Context, identifier string) (*Server, error) {
	servers, err := p.pelican.ListServers(ctx)
	if err != nil {
		return nil, err
	}
	for _, s := range servers {
		if s.Identifier == identifier {
			return &s, nil
		}
	}
	return nil, ErrNotFound
}
