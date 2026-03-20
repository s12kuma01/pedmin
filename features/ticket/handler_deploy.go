package ticket

import (
	"log/slog"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
)

func (t *Ticket) handleDeployPrompt(e *events.ComponentInteractionCreate) {
	_ = e.CreateMessage(discord.NewMessageCreateV2(
		discord.NewContainer(
			discord.NewTextDisplay("パネルを設置するチャンネルを選択してください:"),
			discord.NewActionRow(
				discord.NewChannelSelectMenu(ModuleID+":deploy_channel", "チャンネルを選択...").
					WithChannelTypes(discord.ChannelTypeGuildText),
			),
		),
	).WithEphemeral(true))
}

func (t *Ticket) handleDeployChannelSelect(e *events.ComponentInteractionCreate) {
	data, ok := e.Data.(discord.ChannelSelectMenuInteractionData)
	if !ok || len(data.Values) == 0 {
		return
	}
	channelID := data.Values[0]
	_ = e.UpdateMessage(discord.NewMessageUpdateV2([]discord.LayoutComponent{
		discord.NewContainer(
			discord.NewTextDisplay("パネルを設置するチャンネル: <#"+channelID.String()+">"),
			discord.NewActionRow(
				discord.NewSuccessButton("設置する", ModuleID+":deploy_confirm:"+channelID.String()),
				discord.NewSecondaryButton("キャンセル", ModuleID+":deploy_cancel"),
			),
		),
	}))
}

func (t *Ticket) handleDeployConfirm(e *events.ComponentInteractionCreate, channelIDStr string) {
	_ = e.DeferUpdateMessage()

	channelID, err := snowflake.Parse(channelIDStr)
	if err != nil {
		t.logger.Error("failed to parse channel ID", slog.Any("error", err))
		return
	}

	panel := BuildTicketPanel()
	if _, err := t.client.Rest.CreateMessage(channelID, panel); err != nil {
		t.logger.Error("failed to deploy ticket panel", slog.Any("error", err))
	}
}
