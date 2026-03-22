package player

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	disgobot "github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/snowflake/v2"
	"github.com/s12kuma01/pedmin/module"
	"github.com/s12kuma01/pedmin/store"
)

const ModuleID = "player"

type trackedMessage struct {
	channelID snowflake.ID
	messageID snowflake.ID
}

type Player struct {
	lavalink            disgolink.Client
	client              *disgobot.Client
	store               store.GuildStore
	queues              *QueueManager
	messages            sync.Map // map[snowflake.ID]trackedMessage
	defaultVolume       int
	autoLeaveTimeout    time.Duration
	lavalinkTimeout     time.Duration
	lavalinkLoadTimeout time.Duration
	leaveTimers         sync.Map // map[snowflake.ID]*time.Timer
	logger              *slog.Logger
}

func New(link disgolink.Client, client *disgobot.Client, defaultVolume int, autoLeaveTimeout, lavalinkTimeout, lavalinkLoadTimeout time.Duration, guildStore store.GuildStore, logger *slog.Logger) *Player {
	return &Player{
		lavalink:            link,
		client:              client,
		store:               guildStore,
		queues:              NewQueueManager(),
		defaultVolume:       defaultVolume,
		autoLeaveTimeout:    autoLeaveTimeout,
		lavalinkTimeout:     lavalinkTimeout,
		lavalinkLoadTimeout: lavalinkLoadTimeout,
		logger:              logger,
	}
}

func (p *Player) Info() module.Info {
	return module.Info{
		ID:          ModuleID,
		Name:        "ミュージックプレイヤー",
		Description: "様々なソースから音楽を再生するミュージックプレイヤー",
		AlwaysOn:    false,
	}
}

func (p *Player) Commands() []discord.ApplicationCommandCreate {
	return []discord.ApplicationCommandCreate{
		discord.SlashCommandCreate{
			Name:        "player",
			Description: "ミュージックプレイヤーを表示",
		},
	}
}

// getDefaultVolume returns the per-guild default volume, falling back to the global config default.
func (p *Player) getDefaultVolume(guildID snowflake.ID) int {
	settings, err := LoadSettings(p.store, guildID)
	if err != nil || settings.DefaultVolume == nil {
		return p.defaultVolume
	}
	return *settings.DefaultVolume
}

func (p *Player) SettingsSummary(guildID snowflake.ID) string {
	vol := p.getDefaultVolume(guildID)
	return fmt.Sprintf("デフォルト音量: %d%%", vol)
}

func (p *Player) SettingsPanel(guildID snowflake.ID) []discord.LayoutComponent {
	vol := p.getDefaultVolume(guildID)
	return BuildSettingsPanel(vol)
}
