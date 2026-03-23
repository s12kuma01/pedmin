package bot

import (
	"log/slog"

	"github.com/disgoorg/disgo/discord"
)

func (b *Bot) SyncCommands() error {
	var commands []discord.ApplicationCommandCreate

	for _, m := range b.modules {
		commands = append(commands, m.Commands()...)
	}

	_, err := b.Client.Rest.SetGlobalCommands(b.Client.ApplicationID, commands)
	if err != nil {
		return err
	}

	b.Logger.Info("synced slash commands", slog.Int("count", len(commands)))
	return nil
}
