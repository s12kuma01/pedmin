// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package view

import (
	"fmt"
	"time"

	"github.com/disgoorg/disgo/discord"
)

// BuildPingResponse builds the ping response message.
func BuildPingResponse(latency time.Duration) discord.MessageCreate {
	return discord.NewMessageCreateV2(
		discord.NewContainer(
			discord.NewTextDisplay(
				fmt.Sprintf("🏓 Pong!\n**レイテンシ:** %dms", latency.Milliseconds()),
			),
		),
	).WithEphemeral(true)
}
