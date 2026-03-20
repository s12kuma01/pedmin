package bot

import "github.com/disgoorg/disgo/discord"

func errorMessage(text string) discord.MessageCreate {
	return discord.NewMessageCreateV2(
		discord.NewContainer(
			discord.NewTextDisplay(text),
		),
	).WithEphemeral(true)
}
