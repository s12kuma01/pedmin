package ping

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

func (p *Ping) HandleCommand(e *events.ApplicationCommandInteractionCreate) {
	latency := e.Client().Gateway.Latency()

	_ = e.CreateMessage(discord.NewMessageCreateV2(
		discord.NewContainer(
			discord.NewTextDisplay(
				fmt.Sprintf("🏓 Pong!\n**レイテンシ:** %dms", latency.Milliseconds()),
			),
		).WithAccentColor(0x00B894),
	).WithEphemeral(true))
}
