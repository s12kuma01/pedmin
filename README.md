# Pedmin

Pedmin (pepe + administrator) is a modular Discord bot built with Go and [disgo](https://github.com/disgoorg/disgo). Features Components V2 UI, music playback via Lavalink, and a layered Feature Module architecture.

## Features

- **Feature Module System** - Each feature is a self-contained module (handler/service/view layers)
- **Settings Panel** - `/settings` command with interactive admin UI, per-guild module toggle
- **Music Player** - `/player` with Lavalink-powered playback, queue management, loop modes, and rich V2 UI
- **Support Tickets** - Channel-based ticket system with creation panel, close/reopen, and transcript logging
- **Server Logger** - Configurable event logging (message edit/delete with attachment support, member join/leave, ban/unban, role/channel changes)
- **RSS Feeds** - Background RSS feed monitoring with automatic announcements
- **Avatar Viewer** - `/avatar` command with server/global avatar display via MediaGallery
- **System Info** - `/fuckfetch` neofetch-style system information display
- **Components V2** - Modern Discord UI with containers, sections, and interactive controls
- **SQLite Storage** - Per-guild settings and data with WAL mode for concurrent access

## Quick Start

1. Copy `.env.example` to `.env` and fill in your bot token and app ID
2. Run with Docker:
   ```bash
   docker compose up
   ```

## Development

```bash
go build ./...    # Build
go test ./...     # Test
go vet ./...      # Lint
```

## Documentation

See the [`docs/`](docs/) directory for detailed guides:
- [Architecture](docs/ARCHITECTURE.md) - Layered Feature Module design
- [Module Development](docs/MODULE_GUIDE.md) - Step-by-step guide
- [Components V2](docs/COMPONENTS_V2.md) - disgo V2 component reference
- [Lavalink Integration](docs/LAVALINK.md) - Music playback setup
- [Data Store](docs/STORE.md) - SQLite persistence

## License

See [LICENSE.md](LICENSE.md).
