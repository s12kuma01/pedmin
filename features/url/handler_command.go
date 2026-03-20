package url

import "github.com/disgoorg/disgo/events"

func (u *URL) HandleCommand(e *events.ApplicationCommandInteractionCreate) {
	hasXGD := u.cfg.XGDAPIKey != ""
	hasVT := u.cfg.VTAPIKey != ""

	_ = e.CreateMessage(BuildMainPanel(hasXGD, hasVT))
}
