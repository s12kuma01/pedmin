# Lavalink Integration Guide

## Setup

### Docker Compose
Lavalink runs as a Docker container alongside the bot:
```yaml
services:
  lavalink:
    image: ghcr.io/lavalink-devs/lavalink:4-alpine
    volumes:
      - ./lavalink/application.yml:/opt/Lavalink/application.yml
```

### Configuration (`lavalink/application.yml`)
Key settings:
- `server.port`: WebSocket/REST port (default: 2333)
- `lavalink.server.password`: Authentication password
- `lavalink.server.sources`: Enable/disable audio sources
- `lavalink.plugins`: Plugin dependencies and repositories

### Environment Variables
| Variable | Default | Description |
|----------|---------|-------------|
| `LAVALINK_HOST` | `lavalink:2333` | Lavalink address |
| `LAVALINK_PASSWORD` | `youshallnotpass` | Lavalink auth password |

## disgolink API

### Connecting to a Node
```go
link := disgolink.New(botAppID)
_, err := link.AddNode(ctx, disgolink.NodeConfig{
    Name:     "main",
    Address:  "lavalink:2333",
    Password: "youshallnotpass",
})
```

### Loading Tracks
```go
node := link.BestNode()

// With handler callbacks
node.LoadTracksHandler(ctx, "ytsearch:query", disgolink.NewResultHandler(
    func(track lavalink.Track) { /* single track */ },
    func(playlist lavalink.Playlist) { /* playlist loaded */ },
    func(tracks []lavalink.Track) { /* search results */ },
    func() { /* no matches */ },
    func(err error) { /* load failed */ },
))

// Or direct result
result, err := node.LoadTracks(ctx, "ytsearch:query")
```

### Player Operations
```go
player := link.Player(guildID)

// Play a track
player.Update(ctx, lavalink.WithTrack(track))

// Pause/Resume
player.Update(ctx, lavalink.WithPaused(true))

// Volume (0-200)
player.Update(ctx, lavalink.WithVolume(50))

// Seek
player.Update(ctx, lavalink.WithPosition(30000)) // 30 seconds

// Stop (clear track)
player.Update(ctx, lavalink.WithNullTrack())

// Destroy player
player.Destroy(ctx)
link.RemovePlayer(guildID)
```

### Reading Player State
```go
player.Track()    // *lavalink.Track (nil if nothing playing)
player.Paused()   // bool
player.Position() // lavalink.Duration (milliseconds)
player.Volume()   // int (0-200)
```

## Event Listeners

Register listeners on the disgolink client:
```go
link.AddListeners(
    disgolink.NewListenerFunc(func(player disgolink.Player, event lavalink.TrackStartEvent) {
        // Track started playing
    }),
    disgolink.NewListenerFunc(func(player disgolink.Player, event lavalink.TrackEndEvent) {
        // Track ended - check event.Reason
        if event.Reason == lavalink.TrackEndReasonFinished {
            // Play next track
        }
    }),
)
```

### Event Types
| Event | When |
|-------|------|
| `TrackStartEvent` | Track begins playing |
| `TrackEndEvent` | Track finishes/fails/is replaced |
| `TrackExceptionEvent` | Playback error |
| `TrackStuckEvent` | Track stuck (no audio frames) |
| `WebSocketClosedEvent` | Voice WebSocket closed |

### TrackEndReason Values
| Reason | Meaning | Action |
|--------|---------|--------|
| `TrackEndReasonFinished` | Normal completion | Play next |
| `TrackEndReasonLoadFailed` | Failed to load | Play next / notify |
| `TrackEndReasonStopped` | Manually stopped | Do nothing |
| `TrackEndReasonReplaced` | Another track started | Do nothing |
| `TrackEndReasonCleanup` | Player destroyed | Do nothing |

## Voice Connection

For Lavalink to work, voice state/server updates must be forwarded:
```go
// In bot initialization
bot.WithEventListenerFunc(func(e *events.GuildVoiceStateUpdate) {
    if e.VoiceState.UserID != client.ApplicationID {
        return
    }
    link.OnVoiceStateUpdate(ctx, e.VoiceState.GuildID, e.VoiceState.ChannelID, e.VoiceState.SessionID)
})

bot.WithEventListenerFunc(func(e *events.VoiceServerUpdate) {
    link.OnVoiceServerUpdate(ctx, e.GuildID, e.Token, *e.Endpoint)
})
```

To join a voice channel:
```go
client.UpdateVoiceState(ctx, guildID, &channelID, false, true) // selfMute=false, selfDeaf=true
```

## Plugins

### lavasearch-plugin
Provides enhanced search across multiple sources:
```yaml
plugins:
  lavasearch:
    sources:
      - youtube
      - soundcloud
```

### lavalyrics-plugin
Provides lyrics for playing tracks:
```yaml
plugins:
  lavalyrics:
    sources:
      - youtube
```

## Troubleshooting

| Problem | Solution |
|---------|----------|
| `no lavalink node available` | Check Lavalink is running and accessible |
| Connection refused | Verify host/port in env vars matches Lavalink config |
| No audio | Ensure bot is deafened, voice events are forwarded |
| Track load fails | Check Lavalink logs for source-specific errors |
| WebSocket closed (4014) | Bot needs `GuildVoiceStates` intent |
