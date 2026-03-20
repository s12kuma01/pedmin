package ticket

import "github.com/disgoorg/disgo/discord"

func BuildTicketPanel() discord.MessageCreate {
	return discord.NewMessageCreateV2(
		discord.NewContainer(
			discord.NewTextDisplay("## 🎫 チケットサポート"),
			discord.NewSmallSeparator(),
			discord.NewTextDisplay("サポートが必要な場合は、下のボタンからチケットを作成してください。"),
			discord.NewActionRow(
				discord.NewPrimaryButton("チケットを作成", ModuleID+":create"),
			),
		),
	).WithAllowedMentions(&discord.AllowedMentions{})
}
