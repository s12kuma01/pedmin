package ticket

import (
	"strings"

	"github.com/disgoorg/disgo/events"
)

func (t *Ticket) handleComponent(e *events.ComponentInteractionCreate) {
	customID := e.Data.CustomID()
	_, rest, _ := strings.Cut(customID, ":")
	action, extra, _ := strings.Cut(rest, ":")

	guildID := e.GuildID()
	if guildID == nil {
		return
	}

	switch action {
	case "category":
		t.handleCategorySelect(e, *guildID)
	case "log_prompt":
		t.handleLogPrompt(e)
	case "log_channel":
		t.handleLogChannelSelect(e, *guildID)
	case "role_prompt":
		t.handleRolePrompt(e)
	case "role":
		t.handleRoleSelect(e, *guildID)
	case "deploy_prompt":
		t.handleDeployPrompt(e)
	case "deploy_channel":
		t.handleDeployChannelSelect(e)
	case "deploy_confirm":
		if extra == "" {
			return
		}
		t.handleDeployConfirm(e, extra)
	case "deploy_cancel":
		_ = e.DeferUpdateMessage()
	case "create":
		if !t.bot.IsModuleEnabled(*guildID, ModuleID) {
			return
		}
		_ = e.Modal(BuildCreateTicketModal())
	case "close":
		t.archiveTicket(e, *guildID)
	case "delete":
		t.deleteTicket(e, *guildID)
	}
}
