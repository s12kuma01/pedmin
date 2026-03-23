package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"

	disgobot "github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/disgolink/v3/lavalink"
	"github.com/disgoorg/snowflake/v2"
	"github.com/s12kuma01/pedmin/internal/model"
	"github.com/s12kuma01/pedmin/internal/repository"
	"github.com/s12kuma01/pedmin/internal/view"
)

// SeekStep is the seek increment in milliseconds (10 seconds).
const SeekStep = lavalink.Duration(10_000)

var errNotInVoiceChannel = errors.New("user not in voice channel")

// PlayerService manages music playback, queues, tracked messages, leave timers, and progress tickers.
type PlayerService struct {
	lavalink            disgolink.Client
	client              *disgobot.Client
	store               repository.GuildStore
	queues              *model.QueueManager
	messages            sync.Map // map[snowflake.ID]model.TrackedMessage
	defaultVolume       int
	autoLeaveTimeout    time.Duration
	lavalinkTimeout     time.Duration
	lavalinkLoadTimeout time.Duration
	leaveTimers         sync.Map // map[snowflake.ID]*time.Timer
	progressTickers     sync.Map // map[snowflake.ID]context.CancelFunc
	logger              *slog.Logger
}

// NewPlayerService creates a new PlayerService.
func NewPlayerService(
	link disgolink.Client,
	client *disgobot.Client,
	defaultVolume int,
	autoLeaveTimeout, lavalinkTimeout, lavalinkLoadTimeout time.Duration,
	guildStore repository.GuildStore,
	logger *slog.Logger,
) *PlayerService {
	return &PlayerService{
		lavalink:            link,
		client:              client,
		store:               guildStore,
		queues:              model.NewQueueManager(),
		defaultVolume:       defaultVolume,
		autoLeaveTimeout:    autoLeaveTimeout,
		lavalinkTimeout:     lavalinkTimeout,
		lavalinkLoadTimeout: lavalinkLoadTimeout,
		logger:              logger,
	}
}

// LavalinkCtx creates a context with the lavalink timeout.
func (s *PlayerService) LavalinkCtx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), s.lavalinkTimeout)
}

// Lavalink returns the lavalink client.
func (s *PlayerService) Lavalink() disgolink.Client {
	return s.lavalink
}

// Queues returns the queue manager.
func (s *PlayerService) Queues() *model.QueueManager {
	return s.queues
}

// Client returns the Discord client.
func (s *PlayerService) Client() *disgobot.Client {
	return s.client
}

// JoinVoiceChannel joins the user's current voice channel.
func (s *PlayerService) JoinVoiceChannel(guildID, userID snowflake.ID) error {
	voiceState, ok := s.client.Caches.VoiceState(guildID, userID)
	if !ok || voiceState.ChannelID == nil {
		return errNotInVoiceChannel
	}
	return s.client.UpdateVoiceState(context.Background(), guildID, voiceState.ChannelID, false, true)
}

// GetDefaultVolume returns the per-guild default volume, falling back to the global config default.
func (s *PlayerService) GetDefaultVolume(guildID snowflake.ID) int {
	settings, err := repository.LoadModuleSettings(s.store, guildID, model.PlayerModuleID, func() *model.PlayerSettings {
		return &model.PlayerSettings{}
	})
	if err != nil || settings.DefaultVolume == nil {
		return s.defaultVolume
	}
	return *settings.DefaultVolume
}

// SaveVolumeSettings saves the per-guild volume setting.
func (s *PlayerService) SaveVolumeSettings(guildID snowflake.ID, volume int) error {
	settings := &model.PlayerSettings{DefaultVolume: &volume}
	return repository.SaveModuleSettings(s.store, guildID, model.PlayerModuleID, settings)
}

// Skip skips to the next track in the queue.
func (s *PlayerService) Skip(guildID snowflake.ID) {
	player := s.lavalink.ExistingPlayer(guildID)
	if player == nil {
		return
	}

	queue := s.queues.Get(guildID)
	next, ok := queue.Next()
	if !ok {
		ctx, cancel := s.LavalinkCtx()
		defer cancel()
		_ = player.Update(ctx, lavalink.WithNullTrack())
		return
	}

	ctx, cancel := s.LavalinkCtx()
	defer cancel()
	if err := player.Update(ctx, lavalink.WithTrack(next)); err != nil {
		s.logger.Error("failed to skip", slog.Any("error", err))
	}
}

// Stop stops playback, destroys the player, and leaves voice.
func (s *PlayerService) Stop(guildID snowflake.ID) {
	s.StopProgressTicker(guildID)
	player := s.lavalink.ExistingPlayer(guildID)
	if player != nil {
		ctx, cancel := s.LavalinkCtx()
		_ = player.Destroy(ctx)
		cancel()
		s.lavalink.RemovePlayer(guildID)
	}
	s.queues.Delete(guildID)
	_ = s.client.UpdateVoiceState(context.Background(), guildID, nil, false, false)
}

// CycleLoop cycles the loop mode for the guild queue.
func (s *PlayerService) CycleLoop(guildID snowflake.ID) {
	queue := s.queues.Get(guildID)
	queue.CycleLoop()
}

// Seek adjusts the player position by delta milliseconds.
func (s *PlayerService) Seek(guildID snowflake.ID, delta lavalink.Duration) {
	player := s.lavalink.ExistingPlayer(guildID)
	if player == nil || player.Track() == nil {
		return
	}
	if player.Track().Info.IsStream {
		return
	}

	newPos := player.Position() + delta
	if newPos < 0 {
		newPos = 0
	}
	if length := player.Track().Info.Length; newPos > length {
		newPos = length
	}

	ctx, cancel := s.LavalinkCtx()
	defer cancel()
	if err := player.Update(ctx, lavalink.WithPosition(newPos)); err != nil {
		s.logger.Error("failed to seek", slog.Any("error", err))
	}
}

// Shuffle shuffles the queue for the guild.
func (s *PlayerService) Shuffle(guildID snowflake.ID) {
	queue := s.queues.Get(guildID)
	queue.Shuffle()
}

// ClearQueue clears the queue for the guild.
func (s *PlayerService) ClearQueue(guildID snowflake.ID) {
	queue := s.queues.Get(guildID)
	queue.Clear()
}

// LoadAndPlay loads a track query and starts playback.
// It uses a callback to send the followup message result.
func (s *PlayerService) LoadAndPlay(guildID, userID snowflake.ID, query string, sendFollowup func(text string)) {
	if !strings.HasPrefix(query, "http") {
		query = "ytsearch:" + query
	}

	node := s.lavalink.BestNode()
	if node == nil {
		s.logger.Error("no lavalink node available")
		sendFollowup("❌ 音楽サーバーに接続できません。")
		return
	}

	loadCtx, loadCancel := context.WithTimeout(context.Background(), s.lavalinkLoadTimeout)
	defer loadCancel()
	node.LoadTracksHandler(loadCtx, query, disgolink.NewResultHandler(
		func(track lavalink.Track) {
			s.playTrack(guildID, userID, track)
			sendFollowup(fmt.Sprintf("🎵 キューに追加: **%s**", track.Info.Title))
		},
		func(playlist lavalink.Playlist) {
			if len(playlist.Tracks) == 0 {
				sendFollowup("❌ プレイリストにトラックがありません。")
				return
			}
			queue := s.queues.Get(guildID)
			queue.Add(playlist.Tracks...)
			queue.SetCurrent(0)

			_ = s.JoinVoiceChannel(guildID, userID)
			player := s.lavalink.Player(guildID)
			ctx, cancel := s.LavalinkCtx()
			_ = player.Update(ctx, lavalink.WithTrack(playlist.Tracks[0]))
			cancel()
			sendFollowup(fmt.Sprintf("🎵 プレイリストから %d 曲を追加しました", len(playlist.Tracks)))
		},
		func(tracks []lavalink.Track) {
			if len(tracks) == 0 {
				sendFollowup("❌ 検索結果が見つかりません。")
				return
			}
			s.playTrack(guildID, userID, tracks[0])
			sendFollowup(fmt.Sprintf("🎵 キューに追加: **%s**", tracks[0].Info.Title))
		},
		func() {
			s.logger.Info("no matches found", slog.String("query", query))
			sendFollowup("❌ 検索結果が見つかりません。")
		},
		func(err error) {
			s.logger.Error("load failed", slog.Any("error", err))
			sendFollowup("❌ トラックの読み込みに失敗しました。")
		},
	))
}

func (s *PlayerService) playTrack(guildID, userID snowflake.ID, track lavalink.Track) {
	queue := s.queues.Get(guildID)
	queue.Add(track)

	// Only set current when this is the first track; otherwise the queue
	// already has a current track and this one should wait its turn.
	if queue.Len() == 1 {
		queue.SetCurrent(0)
	}

	_ = s.JoinVoiceChannel(guildID, userID)
	player := s.lavalink.Player(guildID)

	if player.Track() == nil {
		current, ok := queue.Current()
		if !ok {
			current = track
		}
		ctx, cancel := s.LavalinkCtx()
		_ = player.Update(ctx, lavalink.WithTrack(current))
		cancel()
	}
}

// BuildPlayerUI is a convenience wrapper that builds the player UI from current state.
func (s *PlayerService) BuildPlayerUI(guildID snowflake.ID) discord.ContainerComponent {
	player := s.lavalink.Player(guildID)
	queue := s.queues.Get(guildID)
	return view.BuildPlayerUI(player, queue)
}

// BuildQueueUI is a convenience wrapper that builds the queue UI from current state.
func (s *PlayerService) BuildQueueUI(guildID snowflake.ID) discord.ContainerComponent {
	player := s.lavalink.Player(guildID)
	queue := s.queues.Get(guildID)
	return view.BuildQueueUI(queue, player)
}

// Shutdown stops all background goroutines. Call during graceful shutdown.
func (s *PlayerService) Shutdown() {
	s.StopAllProgressTickers()
}
