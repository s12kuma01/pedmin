package ping

import (
	"github.com/disgoorg/disgo/events"
)

func (p *Ping) HandleCommand(e *events.ApplicationCommandInteractionCreate) {
	latency := e.Client().Gateway.Latency()
	_ = e.CreateMessage(BuildPingUI(latency))
}
