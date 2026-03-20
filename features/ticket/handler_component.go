package ticket

import (
	"log/slog"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
)

func (t *Ticket) handleComponent(e *events.ComponentInteractionCreate) {
	customID := e.Data.CustomID()
	parts := strings.SplitN(customID, ":", 3)
	action := parts[1]

	guildID := e.GuildID()
	if guildID == nil {
		return
	}

	switch action {
	case "category":
		data := e.Data.(discord.ChannelSelectMenuInteractionData)
		if len(data.Values) == 0 {
			return
		}
		settings, err := LoadSettings(t.store, *guildID)
		if err != nil {
			t.logger.Error("failed to load settings", slog.Any("error", err))
			return
		}
		settings.CategoryID = data.Values[0]
		if err := SaveSettings(t.store, *guildID, settings); err != nil {
			t.logger.Error("failed to save settings", slog.Any("error", err))
		}
		_ = e.DeferUpdateMessage()

	case "log_prompt":
		_ = e.CreateMessage(discord.NewMessageCreateV2(
			discord.NewContainer(
				discord.NewTextDisplay("ログチャンネルを選択してください:"),
				discord.NewActionRow(
					discord.NewChannelSelectMenu(ModuleID+":log_channel", "ログチャンネルを選択...").
						WithChannelTypes(discord.ChannelTypeGuildText),
				),
			),
		).WithEphemeral(true))

	case "log_channel":
		data := e.Data.(discord.ChannelSelectMenuInteractionData)
		if len(data.Values) == 0 {
			return
		}
		settings, err := LoadSettings(t.store, *guildID)
		if err != nil {
			t.logger.Error("failed to load settings", slog.Any("error", err))
			return
		}
		settings.LogChannelID = data.Values[0]
		if err := SaveSettings(t.store, *guildID, settings); err != nil {
			t.logger.Error("failed to save settings", slog.Any("error", err))
		}
		_ = e.DeferUpdateMessage()

	case "role_prompt":
		_ = e.CreateMessage(discord.NewMessageCreateV2(
			discord.NewContainer(
				discord.NewTextDisplay("サポートロールを選択してください:"),
				discord.NewActionRow(
					discord.NewRoleSelectMenu(ModuleID+":role", "サポートロールを選択..."),
				),
			),
		).WithEphemeral(true))

	case "role":
		data := e.Data.(discord.RoleSelectMenuInteractionData)
		if len(data.Values) == 0 {
			return
		}
		settings, err := LoadSettings(t.store, *guildID)
		if err != nil {
			t.logger.Error("failed to load settings", slog.Any("error", err))
			return
		}
		settings.SupportRoleID = data.Values[0]
		if err := SaveSettings(t.store, *guildID, settings); err != nil {
			t.logger.Error("failed to save settings", slog.Any("error", err))
		}
		_ = e.DeferUpdateMessage()

	case "deploy_prompt":
		_ = e.CreateMessage(discord.NewMessageCreateV2(
			discord.NewContainer(
				discord.NewTextDisplay("パネルを設置するチャンネルを選択してください:"),
				discord.NewActionRow(
					discord.NewChannelSelectMenu(ModuleID+":deploy_channel", "チャンネルを選択...").
						WithChannelTypes(discord.ChannelTypeGuildText),
				),
			),
		).WithEphemeral(true))

	case "deploy_channel":
		data := e.Data.(discord.ChannelSelectMenuInteractionData)
		if len(data.Values) == 0 {
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

	case "deploy_confirm":
		if len(parts) < 3 {
			return
		}
		t.handleDeployConfirm(e, parts[2])

	case "deploy_cancel":
		_ = e.DeferUpdateMessage()

	case "create":
		if !t.bot.IsModuleEnabled(*guildID, ModuleID) {
			return
		}
		_ = e.Modal(discord.ModalCreate{
			CustomID: ModuleID + ":create_modal",
			Title:    "チケットを作成",
			Components: []discord.LayoutComponent{
				discord.NewLabel("件名",
					discord.NewShortTextInput(ModuleID+":subject").
						WithPlaceholder("チケットの件名を入力").
						WithRequired(true).
						WithMaxLength(100),
				),
				discord.NewLabel("説明",
					discord.NewParagraphTextInput(ModuleID+":description").
						WithPlaceholder("詳しい内容を入力してください").
						WithRequired(false).
						WithMaxLength(1000),
				),
			},
		})

	case "close":
		t.archiveTicket(e, *guildID)

	case "delete":
		t.deleteTicket(e, *guildID)
	}
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
