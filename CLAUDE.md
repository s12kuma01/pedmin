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
deepl/
├── client.go                  # DeepL translation API client (shared)
└── lang.go                    # Language code → Japanese name mapping
ui/
├── ui.go                      # Shared UI helpers (EphemeralV2, ErrorMessage)
├── format.go                  # FormatBytes, BuildBar, FormatUptime
└── settings_panel.go          # BuildModulePanel, BuildMainPanel (shared settings UI)
bot/
├── bot.go                     # Client init, module registry, lifecycle
├── commands.go                # Global command sync
├── router.go                  # Interaction → Module dispatch
├── voice.go                   # VoiceState/VoiceServer → Lavalink relay
└── presence.go                # Bot presence updater (CPU/RAM monitoring)
store/
├── store.go                   # SettingsStore/TicketStore/RSSStore/GuildStore interfaces
├── module_settings.go         # Generic LoadModuleSettings/SaveModuleSettings helpers
├── sqlite_store.go            # SQLite implementation (WAL mode)
├── sqlite_migrations.go       # Schema migrations
├── sqlite_modules.go          # Module settings persistence
├── sqlite_ticket.go           # Ticket persistence
└── sqlite_rss.go              # RSS feed persistence
features/settings/
├── module.go                  # Info, Commands, Bot interface
├── handler_command.go         # /settings slash command
├── handler_component.go       # Select/toggle/back dispatch
└── view.go                    # (delegates to ui/settings_panel.go)
features/ping/
├── module.go                  # Info, Commands
├── handler_command.go         # /ping command
└── view.go                    # Ping response UI builder
features/avatar/
├── module.go                  # Info, Commands
├── handler_command.go         # /avatar command, user resolve
└── view_avatar.go             # Avatar MediaGallery builder
features/embedfix/
├── module.go                  # Info, Bot deps, module interface
├── listener.go                # GuildMessageCreate listener, URL detection
├── handler_component.go       # Translate button + platform settings dispatch
├── service_embed.go           # URL processing, platform-specific embed sending
├── service_translate.go       # Translation workflow per platform
├── settings.go                # Settings struct & persistence
├── domain.go                  # Platform type, EmbedRef, URL regex matching
├── client_twitter.go          # FxTwitter API client
├── client_reddit.go           # Reddit JSON API client
├── client_tiktok.go           # TikTok proxy API client
├── view_twitter.go            # Tweet embed UI builder
├── view_reddit.go             # Reddit post embed UI builder
├── view_tiktok.go             # TikTok video embed UI builder
├── view_settings.go           # Platform toggle settings panel
└── view_helpers.go            # Emoji constants, formatCount
features/translator/
├── module.go                  # Info, Bot deps
├── listener.go                # MessageReactionAdd listener
├── service.go                 # Fetch message → translate → post
├── view.go                    # Translation embed UI builder
└── view_helpers.go            # Flag emoji → language code mapping
features/fuckfetch/
├── module.go                  # Info, Commands
├── handler_command.go         # /fuckfetch command
├── service.go                 # System info gathering
└── view.go                    # Neofetch-style output builder
features/panel/
├── module.go                  # Info, Commands, permission check
├── handler_command.go         # /panel slash command
├── handler_component.go       # Button/select dispatch
├── handler_modal.go           # Console command modal
├── service.go                 # Server list/detail/power/console operations
├── client.go                  # Pelican API HTTP client + domain types
├── client_actions.go          # Power/console action methods
├── view_panel.go              # Server list, detail panels
├── view_console.go            # Console result/error panels
└── view_helpers.go            # Format helpers (bytes, bars, uptime, emoji)
features/player/
├── module.go                  # Info, Commands
├── handler_command.go         # /player slash command
├── handler_component.go       # Button/select switch dispatch
├── handler_modal.go           # Add-to-queue modal
├── handler_queue.go           # Queue page navigation
├── service.go                 # Playback logic (Discord API independent)
├── settings.go                # Per-guild volume settings
├── voice.go                   # VC connection helper
├── queue.go                   # Queue data structure
├── queue_manager.go           # Per-guild queue management
├── loop_mode.go               # LoopMode type + constants
├── lavalink.go                # Lavalink event listeners + node connection
├── auto_leave.go              # Auto-leave on empty VC
├── message_tracker.go         # Player message tracking for updates
├── view_player.go             # Player UI builder
├── view_queue.go              # Queue UI builder
├── view_settings.go           # Volume settings panel
└── view_helpers.go            # Progress bar, duration format, thumbnails
features/url/
├── module.go                  # Info, Commands
├── handler_command.go         # /url command
├── handler_component.go       # Button dispatch (shorten/check/back)
├── handler_modal.go           # Modal submission (shorten/check)
├── service.go                 # URL validation, shorten, scan
├── client.go                  # URLClient struct + x.gd shorten
├── client_virustotal.go       # VirusTotal scan client
└── view.go                    # Main panel, result, error panels
features/ticket/
├── module.go                  # Info, Commands, Bot/Client/Store deps
├── handler_component.go       # Create/close/reopen ticket buttons
├── handler_settings.go        # Settings UI interactions
├── handler_modal.go           # Ticket creation modal
├── handler_deploy.go          # Panel deployment
├── service.go                 # Ticket creation/closure logic
├── service_log.go             # Log & transcript sending
├── service_settings.go        # Category/log channel/role updates
├── transcript.go              # HTML transcript generation
├── settings.go                # Settings struct & persistence
├── view_settings.go           # Settings panel UI
├── view_panel.go              # Ticket control panel UI
├── view_ticket.go             # Ticket channel message UI
└── view_log.go                # Ticket list/log UI
features/logger/
├── module.go                  # Info, Bot/Client/Store deps
├── listener.go                # Listener setup + sendLog helper
├── listener_message.go        # Message edit/delete listeners
├── listener_guild.go          # Member/ban/role/channel listeners
├── handler_component.go       # Settings component handling
├── settings.go                # Logger settings (channel ID, event toggles)
├── view_settings.go           # Settings panel UI
├── view_message_log.go        # Message edit/delete log builders
├── view_guild_log.go          # Member/ban/role/channel log builders
├── view_structure_log.go      # Channel structure change log builders
└── view_attachment.go         # Attachment diff & display
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
`SettingsPanel()`. Optional: `SettingsSummarizer`, `VoiceStateListener`. Registered in `main.go` via `bot.Register()`.

### CustomID Convention

Component CustomIDs follow `{moduleID}:{action}:{extra}`. Router splits on the first colon to dispatch.

### Components V2

All UI uses `discord.NewMessageCreateV2()`. View files are pure functions: state in → components out. No accent colors
on containers.

### GuildStore Interface

`store.GuildStore` abstracts persistence, composed from `SettingsStore`, `TicketStore`, and `RSSStore` sub-interfaces
(ISP). Generic `LoadModuleSettings[T]`/`SaveModuleSettings[T]` helpers reduce per-module boilerplate. SQLite at
`data/pedmin.db` with WAL mode. No JSON fallback.

## Documentation

- `docs/ARCHITECTURE.md` - System architecture, layers, data flow
- `docs/MODULE_GUIDE.md` - How to create new modules
- `docs/COMPONENTS_V2.md` - Components V2 reference for disgo
- `docs/LAVALINK.md` - Lavalink integration guide
- `docs/STORE.md` - Data persistence guide
