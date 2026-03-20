package ping

import (
	"fmt"
	"time"

	"github.com/disgoorg/disgo/discord"
)

// BuildPingUI builds the ping response message.
func BuildPingUI(latency time.Duration) discord.MessageCreate {
	return discord.NewMessageCreateV2(
		discord.NewContainer(
			discord.NewTextDisplay(
				fmt.Sprintf("🏓 Pong!\n**レイテンシ:** %dms", latency.Milliseconds()),
			),
		),
	).WithEphemeral(true)
}
