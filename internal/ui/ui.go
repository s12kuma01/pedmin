// opyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package ui

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
)

// EphemeralV2 wraps layout components into an ephemeral Components V2 message.
func EphemeralV2(components ...discord.LayoutComponent) discord.MessageCreate {
	return discord.NewMessageCreateV2(components...).WithEphemeral(true)
}

// EphemeralError builds an ephemeral error message with ❌ prefix.
func EphemeralError(text string) discord.MessageCreate {
	return EphemeralV2(
		discord.NewContainer(
			discord.NewTextDisplay(fmt.Sprintf("❌ %s", text)),
		),
	)
}

// ErrorMessage builds an ephemeral message with plain text in a container.
func ErrorMessage(text string) discord.MessageCreate {
	return EphemeralV2(
		discord.NewContainer(
			discord.NewTextDisplay(text),
		),
	)
}
