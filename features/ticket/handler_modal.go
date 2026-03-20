package ticket

import (
	"fmt"
	"log/slog"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

func (t *Ticket) handleModal(e *events.ModalSubmitInteractionCreate) {
	if e.Data.CustomID != ModuleID+":create_modal" {
		return
	}

	guildID := e.GuildID()
	if guildID == nil {
		return
	}

	subject := ""
	if ti, ok := e.Data.TextInput(ModuleID + ":subject"); ok {
		subject = ti.Value
	}
	description := ""
	if ti, ok := e.Data.TextInput(ModuleID + ":description"); ok {
		description = ti.Value
	}

	if subject == "" {
		_ = e.CreateMessage(discord.NewMessageCreateV2(
			discord.NewContainer(
				discord.NewTextDisplay("件名を入力してください。"),
			),
		).WithEphemeral(true))
		return
	}

	_ = e.DeferCreateMessage(true)

	channelID, _, err := t.createTicket(*guildID, e.User().ID, subject, description)
	if err != nil {
		t.logger.Error("failed to create ticket", slog.Any("error", err))
		_, _ = e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(), discord.NewMessageUpdateV2([]discord.LayoutComponent{
			discord.NewContainer(
				discord.NewTextDisplay("チケットの作成に失敗しました。"),
			),
		}))
		return
	}

	_, _ = e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(), discord.NewMessageUpdateV2([]discord.LayoutComponent{
		discord.NewContainer(
			discord.NewTextDisplay(fmt.Sprintf("チケットを作成しました: <#%d>", channelID)),
		),
	}))
}
