# Pedmin - Discord Bot

Pedmin (pepe + administrator) is a modular Discord bot built with Go 1.26.1 and disgo v0.19.2. It serves as a Probot
replacement, featuring Components V2 UI, music playback via Lavalink, and a standard Go layered architecture. Runs on
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
go build ./cmd/pedmin/...

# Run tests
go test ./...

# Vet
go vet ./...

# Docker
docker compose up        # Start bot + Lavalink
docker compose up -d     # Detached mode
docker compose build     # Rebuild bot image
```

## Architecture: Standard Go Layered Pattern

The project follows standard Go large-service conventions with `cmd/`, `internal/`, and `pkg/` directories. Layers are
separated by package: handler (controller) → service (business logic) → repository (persistence), with shared model and
view packages.

```
cmd/pedmin/main.go                 # Entrypoint: DI wiring, graceful shutdown
config/
├── config.go                      # Env vars + TOML file loading
└── defaults.go                    # Default configuration values
migrations/
├── embed.go                       # embed.FS export for SQL files
├── 001_guild_modules.sql          # Guild modules + settings tables
├── 002_tickets.sql                # Tickets table
├── 003_rss_feeds.sql              # RSS feeds + seen items tables
├── 004_counters.sql               # Word counters + hit logs tables
├── 005_leveling.sql               # User XP + role rewards tables
└── 006_component_panels.sql       # Component builder panels table
pkg/deepl/
├── client.go                      # DeepL translation API client
└── lang.go                        # Language code → Japanese name mapping
pkg/rankcard/
├── rankcard.go                    # Canvas rank card image generation (fogleman/gg)
└── font.go                        # Embedded Noto Sans JP font
internal/
├── module/
│   └── module.go                  # Module interface definition
├── bot/
│   ├── bot.go                     # Client init, module registry, lifecycle
│   ├── commands.go                # Global command sync
│   ├── router.go                  # Interaction → Module dispatch
│   ├── voice.go                   # VoiceState/VoiceServer → Lavalink relay
│   └── presence.go                # Bot presence updater (CPU/RAM monitoring)
├── model/                         # Domain types, settings, constants
│   ├── constants.go               # All ModuleID constants
│   ├── guild.go                   # GuildSettings, Ticket, RSSFeed
│   ├── player.go                  # Queue, QueueManager, LoopMode, PlayerSettings
│   ├── embedfix.go                # Platform, EmbedRef, EmbedFixSettings, URL matchers
│   ├── twitter.go                 # Tweet, TweetAuthor, TweetMedia
│   ├── reddit.go                  # RedditPost
│   ├── tiktok.go                  # TikTokVideo, TikTokAuthor
│   ├── ticket.go                  # TicketSettings
│   ├── logger.go                  # LoggerSettings, event constants
│   ├── panel.go                   # Server, ServerLimits, Resources
│   ├── url.go                     # VTResult
│   ├── fuckfetch.go               # SystemInfo
│   ├── translator.go              # FlagToLang mapping
│   ├── counter.go                 # Counter, MatchType, CounterStat
│   ├── leveling.go                # UserXP, LevelingSettings, XP formula
│   ├── autorole.go                # AutoroleSettings
│   └── builder.go                 # ComponentPanel, PanelComponent types
├── repository/                    # Data persistence (GuildStore interface)
│   ├── repository.go              # GuildStore = SettingsStore + TicketStore + RSSStore + CounterStore + LevelingStore + PanelBuilderStore
│   ├── module_settings.go         # Generic LoadModuleSettings/SaveModuleSettings helpers
│   ├── sqlite.go                  # SQLite implementation (WAL mode) + migrations
│   ├── sqlite_modules.go          # Module settings persistence
│   ├── sqlite_ticket.go           # Ticket persistence
│   ├── sqlite_rss.go              # RSS feed persistence
│   ├── sqlite_counter.go          # Word counter persistence
│   ├── sqlite_leveling.go         # Leveling XP + role rewards persistence
│   └── sqlite_builder.go          # Component panel persistence
├── client/                        # External API clients
│   ├── twitter.go                 # FxTwitter API client
│   ├── reddit.go                  # Reddit JSON API client
│   ├── tiktok.go                  # TikTok proxy API client
│   ├── panel.go                   # Pelican API client + actions
│   ├── url.go                     # x.gd URL shortener client
│   └── virustotal.go              # VirusTotal scan client
├── service/                       # Business logic (Discord event-independent)
│   ├── player.go                  # Playback control, queue operations
│   ├── player_lavalink.go         # Lavalink event listeners + node connection
│   ├── player_autoleave.go        # Auto-leave timer logic
│   ├── player_progress.go         # Progress ticker + message tracker
│   ├── embedfix.go                # URL detection, platform embed sending
│   ├── embedfix_translate.go      # Translation workflow per platform
│   ├── ticket.go                  # Ticket creation/closure logic
│   ├── ticket_log.go              # Log, transcript sending, HTML generation
│   ├── ticket_settings.go         # Category/log channel/role updates
│   ├── logger.go                  # Logger settings load/save
│   ├── rss.go                     # Feed CRUD, validation, poll logic
│   ├── rss_poller.go              # Background polling routine
│   ├── translator.go              # Translation service
│   ├── panel.go                   # Server list/detail/power/console
│   ├── url.go                     # URL validation, shorten, scan
│   ├── fuckfetch.go               # System info gathering
│   ├── counter.go                 # Word counter logic + regex cache
│   ├── leveling.go                # XP processing, cooldowns, multipliers
│   ├── leveling_voice.go          # Voice XP session tracking + ticker
│   ├── leveling_rankcard.go       # Rank card image generation
│   ├── autorole.go                # Auto role assignment on member join
│   └── builder.go                 # Component panel CRUD + deploy
├── handler/                       # Discord interaction handlers (module.Module impl)
│   ├── ping.go                    # /ping command
│   ├── avatar.go                  # /avatar command
│   ├── fuckfetch.go               # /fuckfetch command
│   ├── settings.go                # /settings command + component dispatch
│   ├── player.go                  # /player command + module interface
│   ├── player_component.go        # Player button/select dispatch
│   ├── player_modal.go            # Add-to-queue modal
│   ├── player_queue.go            # Queue page navigation
│   ├── embedfix.go                # Module interface + listener setup
│   ├── embedfix_component.go      # Translate button + settings dispatch
│   ├── embedfix_listener.go       # MessageCreate listener
│   ├── translator.go              # Module interface + listener setup
│   ├── translator_listener.go     # ReactionAdd listener
│   ├── panel.go                   # /panel command
│   ├── panel_component.go         # Panel button/select dispatch
│   ├── panel_modal.go             # Console command modal
│   ├── ticket.go                  # Module interface
│   ├── ticket_component.go        # Create/close/reopen buttons
│   ├── ticket_modal.go            # Ticket creation modal
│   ├── ticket_deploy.go           # Panel deployment
│   ├── ticket_settings.go         # Settings UI interactions
│   ├── logger.go                  # Module interface + listener setup
│   ├── logger_listener.go         # Message/guild event listeners
│   ├── logger_component.go        # Settings component handling
│   ├── rss.go                     # Module interface
│   ├── rss_component.go           # Add/remove feed dispatch
│   ├── rss_add_feed.go            # Add feed prompt/validation
│   ├── rss_modal.go               # Feed URL input modal
│   ├── url.go                     # /url command
│   ├── url_component.go           # Shorten/check/back buttons
│   ├── url_modal.go               # Modal submission
│   ├── counter.go                 # Counter module interface + listener
│   ├── counter_component.go       # Counter settings/stats dispatch
│   ├── counter_modal.go           # Counter word input modal
│   ├── leveling.go                # Leveling module + /rank, /leaderboard
│   ├── leveling_listener.go       # Message XP listener
│   ├── leveling_component_settings.go  # Leveling settings tab dispatch
│   ├── leveling_component_rewards.go   # Rewards/multipliers dispatch
│   ├── leveling_modal_settings.go      # XP/cooldown/voice modals
│   ├── leveling_modal_rewards.go       # Reward/multiplier modals
│   ├── autorole.go                # Autorole module + member join listener
│   ├── builder.go                 # Builder module interface
│   ├── builder_component_add.go   # Add component flows
│   ├── builder_component_manage.go # Component management
│   ├── builder_component_deploy.go # Deploy/preview/delete flows
│   ├── builder_modal.go           # Panel create/rename modals
│   └── builder_modal_components.go # Component input modals
├── view/                          # UI builders (pure functions)
│   ├── ping.go                    # Ping response
│   ├── avatar.go                  # Avatar MediaGallery
│   ├── fuckfetch.go               # Neofetch-style output
│   ├── player.go                  # Player UI
│   ├── player_queue.go            # Queue UI
│   ├── player_settings.go         # Volume settings panel
│   ├── player_helpers.go          # Progress bar, duration format
│   ├── embedfix_twitter.go        # Tweet embed
│   ├── embedfix_reddit.go         # Reddit embed
│   ├── embedfix_tiktok.go         # TikTok embed
│   ├── embedfix_settings.go       # Platform toggle settings
│   ├── embedfix_helpers.go        # Emoji constants, formatCount
│   ├── translator.go              # Translation embed
│   ├── panel.go                   # Server list/detail panels
│   ├── panel_console.go           # Console result/error panels
│   ├── panel_helpers.go           # Format helpers
│   ├── ticket_panel.go            # Ticket control panel
│   ├── ticket_ticket.go           # Ticket channel message
│   ├── ticket_log.go              # Ticket list/log
│   ├── ticket_settings.go         # Ticket settings panel
│   ├── logger_settings.go         # Logger settings panel
│   ├── logger_message.go          # Message edit/delete logs
│   ├── logger_guild.go            # Member/ban/role/channel logs
│   ├── logger_structure.go        # Channel structure logs
│   ├── logger_attachment.go       # Attachment diff/display
│   ├── rss_settings.go            # RSS settings panel
│   ├── rss_manage.go              # Feed list/detail
│   ├── rss_feed.go                # Feed item announcement
│   ├── rss_helpers.go             # Text utilities
│   ├── url.go                     # URL panels
│   ├── counter_settings.go        # Counter settings panel
│   ├── counter_manage.go          # Counter list/detail + add type prompt
│   ├── counter_stats.go           # Counter stats + user ranking
│   ├── leveling_settings.go       # Leveling general settings tab
│   ├── leveling_settings_rewards.go    # Leveling rewards tab
│   ├── leveling_settings_multipliers.go # Leveling multipliers tab
│   ├── leveling_leaderboard.go    # Leaderboard with pagination
│   ├── leveling_levelup.go        # Level-up notification (placeholder)
│   ├── autorole_settings.go       # Autorole settings panel
│   ├── builder_list.go            # Panel list UI
│   ├── builder_edit.go            # Panel edit mode UI
│   ├── builder_manage.go          # Component management UI
│   ├── builder_deploy.go          # Deploy flow UI
│   ├── builder_helpers.go         # Builder view helpers
│   └── builder_render.go          # Panel render (JSON → Components V2)
└── ui/                            # Shared Discord UI helpers
    ├── ui.go                      # EphemeralV2, ErrorMessage
    ├── format.go                  # FormatBytes, BuildBar, FormatUptime
    └── settings_panel.go          # BuildModulePanel, BuildMainPanel
```

## Key Design Decisions

### Standard Go Project Layout

The project follows `cmd/`, `internal/`, `pkg/` conventions. Layers are separated by package, not by file within a
package. This enables clear dependency boundaries and testability.

### Layer Separation

- **handler** (controller): Receives Discord events, extracts data, calls service, builds view, sends response
- **service** (business logic): Discord event-independent, operates on primitive types and model types
- **view** (UI builder): Pure functions — state in → Discord components out
- **model** (domain): Shared types, constants, settings structs
- **repository** (persistence): Database interfaces and implementations
- **client** (external APIs): HTTP clients for third-party services

### Module Interface (`module.Module`)

All features implement: `Info()`, `Commands()`, `HandleCommand()`, `HandleComponent()`, `HandleModal()`,
`SettingsPanel()`. Optional: `SettingsSummarizer`, `VoiceStateListener`. Handler structs implement this interface.
Registered in `main.go` via `bot.Register()`.

### CustomID Convention

Component CustomIDs follow `{moduleID}:{action}:{extra}`. Router splits on the first colon to dispatch.
ModuleID constants are centralized in `internal/model/constants.go`.

### Components V2

All UI uses `discord.NewMessageCreateV2()`. View files are pure functions: state in → components out. No accent colors
on containers.

### GuildStore Interface

`repository.GuildStore` abstracts persistence, composed from `SettingsStore`, `TicketStore`, `RSSStore`,
`CounterStore`, `LevelingStore`, and `PanelBuilderStore` sub-interfaces (ISP). Generic `LoadModuleSettings[T]`/`SaveModuleSettings[T]` helpers reduce per-module boilerplate.
SQLite at `data/pedmin.db` with WAL mode. Migrations loaded via `embed.FS` from `migrations/`.

### Import Dependency Graph (no cycles)

```
cmd/pedmin/main.go → config, internal/bot, internal/handler, internal/service, internal/repository, internal/client
internal/handler   → internal/service, internal/view, internal/model, internal/module, internal/ui
internal/service   → internal/repository, internal/model, internal/view, internal/client, pkg/deepl, pkg/rankcard
internal/view      → internal/model, internal/ui
internal/bot       → internal/module, internal/repository, internal/ui, config
internal/repository → internal/model
internal/client    → internal/model
internal/ui        → internal/module
internal/module    → (leaf: disgo types only)
internal/model     → (leaf: disgo, disgolink types only)
pkg/deepl          → (leaf: stdlib only)
pkg/rankcard       → (leaf: fogleman/gg, golang.org/x/image)
```

## Documentation

- `docs/ARCHITECTURE.md` - System architecture, layers, data flow
- `docs/MODULE_GUIDE.md` - How to create new modules
- `docs/COMPONENTS_V2.md` - Components V2 reference for disgo
- `docs/LAVALINK.md` - Lavalink integration guide
- `docs/STORE.md` - Data persistence guide
