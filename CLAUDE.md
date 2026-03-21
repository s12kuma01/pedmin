# Pedmin - Discord Bot

Pedmin (pepe + administrator) is a modular Discord bot built with Go 1.26.1 and disgo v0.19.2. It serves as a Probot
replacement, featuring Components V2 UI, music playback via Lavalink, and a layered Feature Module architecture. Runs on
Windows Docker Desktop.

## Tech Stack

- **Language**: Go 1.26.1
- **Discord Library**: disgo v0.19.2
- **Lavalink Client**: disgolink v3.1.0
- **Lavalink Server**: Lavalink 4 (Alpine)
- **Data Storage**: SQLite (`modernc.org/sqlite`, pure Go), behind `GuildStore` interface
- **Configuration**: Environment variables (secrets) + TOML file (app settings)

## Commands

```bash
# Build
go build ./...

# Run tests
go test ./...

# Vet
go vet ./...

# Docker
docker compose up        # Start bot + Lavalink
docker compose up -d     # Detached mode
docker compose build     # Rebuild bot image
```

## Architecture: Layered Feature Module Pattern

Each feature is a self-contained module with internal layer separation (handler/service/view), all within the same Go
package.

```
main.go                        # Entrypoint: DI wiring, graceful shutdown
config/config.go               # Env vars + TOML file loading
module/module.go               # Module interface definition
bot/
├── bot.go                     # Client init, module registry, lifecycle
├── commands.go                # Global command sync
├── router.go                  # Interaction → Module dispatch
├── ui.go                      # Shared UI helpers (errorMessage)
├── voice.go                   # VoiceState/VoiceServer → Lavalink relay
└── presence.go                # Bot presence updater (CPU/RAM monitoring)
store/
├── store.go                   # GuildStore interface
├── sqlite_store.go            # SQLite implementation (WAL mode)
└── sqlite_migrations.go       # Schema migrations
features/settings/
├── module.go                  # Info, Commands, Bot interface
├── handler.go                 # HandleCommand / HandleComponent
└── view.go                    # UI builders (mainPanel, modulePanel)
features/ping/
├── module.go                  # Info, Commands
├── handler_command.go         # /ping command
└── view.go                    # Ping response UI builder
features/avatar/
├── module.go                  # Info, Commands
├── handler_command.go         # /avatar command, user resolve
└── view_avatar.go             # Avatar MediaGallery builder
features/embedfix/
├── module.go              # Info, Bot deps, module interface
├── listener.go            # GuildMessageCreate handler, URL detection
├── handler_component.go   # Translate button dispatch
├── client.go              # fxtwitter + Google Translate HTTP clients
├── view.go                # Tweet embed UI builder
└── view_helpers.go        # Relative time, number format, URL regex
features/fuckfetch/
├── module.go                  # Info, Commands
├── handler_command.go         # /fuckfetch command
├── service.go                 # System info gathering
├── view.go                    # Neofetch-style output builder
└── view_helpers.go            # Formatting helpers (bytes, bars, uptime)
features/panel/
├── module.go                  # Info, Commands, permission check
├── handler_command.go         # /panel slash command
├── handler_component.go       # Button/select dispatch + modal handling
├── service.go                 # Server list/detail/power/console operations
├── client.go                  # Pelican API HTTP client
├── view_panel.go              # Server list, detail, error panels
└── view_helpers.go            # Format helpers (bytes, bars, uptime, emoji)
features/player/
├── module.go                  # Info, Commands
├── handler_command.go         # /player slash command
├── handler_component.go       # Button/select switch dispatch
├── handler_modal.go           # Add-to-queue modal
├── service.go                 # Playback logic (Discord API independent)
├── voice.go                   # VC connection helper
├── queue.go                   # Queue data structure
├── queue_manager.go           # Per-guild queue management
├── loop_mode.go               # LoopMode type + constants
├── lavalink.go                # Lavalink event listeners + node connection
├── auto_leave.go              # Auto-leave on empty VC
├── view_player.go             # Player UI builder
├── view_queue.go              # Queue UI builder
└── view_helpers.go            # Progress bar, duration format, thumbnails
features/url/
├── module.go                  # Info, Commands
├── handler_command.go         # /url command
├── handler_component.go       # Button dispatch (shorten/check/back)
├── handler_modal.go           # Modal submission (shorten/check)
├── service.go                 # URL validation, shorten, scan
├── client.go                  # x.gd + VirusTotal HTTP clients
└── view.go                    # Main panel, result, error panels
features/ticket/
├── module.go                  # Info, Commands, Bot/Client/Store deps
├── handler_component.go       # Create/close/reopen ticket buttons
├── handler_settings.go        # Settings UI interactions
├── handler_modal.go           # Ticket creation modal
├── handler_deploy.go          # Panel deployment
├── service.go                 # Ticket creation/closure/settings logic
├── service_log.go             # Log & transcript sending
├── transcript.go              # HTML transcript generation
├── settings.go                # Settings struct & persistence
├── view_settings.go           # Settings panel UI
├── view_panel.go              # Ticket control panel UI
├── view_ticket.go             # Ticket channel message UI
└── view_log.go                # Ticket list/log UI
features/logger/
├── module.go                  # Info, Bot/Client/Store deps
├── listener.go                # Event listeners (messages, members, bans, roles, channels)
├── handler.go                 # Component interaction handling
├── settings.go                # Logger settings (channel ID, event toggles)
├── view_settings.go           # Settings UI
└── view_log.go                # Log message builders (text, attachments, MediaGallery)
features/rss/
├── module.go                  # Info, Bot/Client/Store deps
├── handler_component.go       # Add/remove feed dispatch
├── handler_add_feed.go        # Add feed prompt & validation
├── handler_modal.go           # Feed URL input modal
├── service.go                 # Feed CRUD, validation, post logic
├── service_poll.go            # Single feed poll logic
├── poller.go                  # Background polling routine
├── view_settings.go           # Settings panel (feed count)
├── view_manage.go             # Feed list/detail UI
├── view_feed.go               # Feed item announcement builder
└── view_helpers.go            # Text utilities (stripHTML, truncate)
```

## Key Design Decisions

### 1 File = 1 Responsibility

Every `.go` file has a single, clear responsibility. No file mixes handler logic with UI building or service logic.

### Feature Module Pattern

Each feature (`features/player/`, `features/settings/`) is a self-contained Go package. Internal layers (handler →
service → view) are separated by file, not by package. Same `package player` throughout — no circular import issues.

### Module Interface (`module.Module`)

All features implement: `Info()`, `Commands()`, `HandleCommand()`, `HandleComponent()`, `HandleModal()`,
`SettingsPanel()`, `HandleSettingsComponent()`. Registered in `main.go` via `bot.Register()`.

### CustomID Convention

Component CustomIDs follow `{moduleID}:{action}:{extra}`. Router splits on the first colon to dispatch.

### Components V2

All UI uses `discord.NewMessageCreateV2()`. View files are pure functions: state in → components out. No accent colors
on containers.

### GuildStore Interface

`store.GuildStore` abstracts persistence. SQLite at `data/pedmin.db` with WAL mode. No JSON fallback.

## Documentation

- `docs/ARCHITECTURE.md` - System architecture, layers, data flow
- `docs/MODULE_GUIDE.md` - How to create new modules
- `docs/COMPONENTS_V2.md` - Components V2 reference for disgo
- `docs/LAVALINK.md` - Lavalink integration guide
- `docs/STORE.md` - Data persistence guide
