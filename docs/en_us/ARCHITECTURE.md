# Architecture

## System Overview

```
┌─────────────┐     Gateway      ┌──────────────┐
│   Discord    │ ◄──────────────► │   Pedmin     │
│   API        │   Interactions   │   Bot        │
└─────────────┘                   └──────┬───────┘
                                         │
                                         │ disgolink
                                         │
                                  ┌──────▼───────┐
                                  │   Lavalink    │
                                  │   Server      │
                                  └──────────────┘
```

Pedmin connects to Discord via the Gateway (WebSocket) using disgo. Voice playback is handled by Lavalink, connected
through disgolink. Both services run as Docker containers.

## Standard Go Layered Architecture

The project follows the standard Go large-service layout (`cmd/`, `internal/`, `pkg/`). Layers are separated by
**package**, not by file.

```
┌──────────────────────────────────────────────────────┐
│  cmd/pedmin/main.go (Entrypoint: DI wiring)          │
└─────────────────────┬────────────────────────────────┘
                      │
┌─────────────────────▼────────────────────────────────┐
│  internal/bot/ (Framework: connection, routing)       │
└─────────────────────┬────────────────────────────────┘
                      │
┌─────────────────────▼────────────────────────────────┐
│  internal/handler/    ← Discord interaction handlers  │
│  internal/service/    ← Business logic                │
│  internal/view/       ← UI builders (pure functions)  │
│  internal/model/      ← Domain types & constants      │
└─────────────────────┬────────────────────────────────┘
                      │
┌─────────────────────▼────────────────────────────────┐
│  internal/repository/ (Infrastructure: SQLite)        │
│  internal/client/     (Infrastructure: HTTP clients)  │
└──────────────────────────────────────────────────────┘
```

Each layer is its own Go package. All 12 feature modules share the same handler, service, view, and model packages.
File names within each package are prefixed by feature name (e.g., `handler/player.go`, `handler/embedfix.go`).

## Package Dependency Graph

```
cmd/pedmin/main.go
 ├── config
 ├── internal/bot         → internal/module, internal/repository, internal/ui
 ├── internal/handler     → internal/service, internal/view, internal/model, internal/module, internal/ui
 ├── internal/service     → internal/repository, internal/model, internal/client, pkg/deepl
 ├── internal/view        → internal/model, internal/ui
 ├── internal/repository  → internal/model
 ├── internal/client      → internal/model
 ├── internal/ui          → internal/module
 ├── internal/module      → (leaf: disgo types only)
 ├── internal/model       → (leaf: disgo, disgolink types only)
 └── pkg/deepl            → (leaf: stdlib only)
```

Dependencies flow downward. Features never depend on each other. Modules that need bot functionality (settings,
embedfix, ticket, logger, rss) define a local `Bot` interface with only the methods they need (ISP), avoiding direct
`bot` package imports.

## Layer Responsibilities

| Layer | Package | Does | Does NOT |
|-------|---------|------|----------|
| **Handler** | `internal/handler/` | Receive Discord events, extract data, call service, build view response | Contain business logic, API calls, or store access |
| **Service** | `internal/service/` | Business logic, API calls, store operations | Build Discord UI components |
| **View** | `internal/view/` | Build Components V2 UI (pure functions: data in → components out) | Have side effects, call APIs, or access store |
| **Model** | `internal/model/` | Domain types, settings structs, ModuleID constants | Import handler/service/view/repository packages |
| **Repository** | `internal/repository/` | Data persistence interfaces and SQLite implementation | Contain business logic |
| **Client** | `internal/client/` | External API HTTP wrappers (Twitter, Reddit, TikTok, Pelican, etc.) | Contain business logic |
| **UI** | `internal/ui/` | Shared Discord UI helpers (EphemeralV2, ErrorMessage, BuildModulePanel) | Contain feature-specific logic |
| **Bot** | `internal/bot/` | Discord client lifecycle, module registry, interaction routing | Contain feature logic |

## File Responsibilities

### internal/bot/

| File | Responsibility |
|------|----------------|
| `bot.go` | Client init, module registry, Start/Close, module state checks |
| `commands.go` | Slash command global sync |
| `router.go` | Interaction dispatch to modules |
| `voice.go` | VoiceState/VoiceServer event relay to Lavalink |
| `presence.go` | Bot presence updater (CPU/RAM monitoring) |

### internal/handler/ (per feature)

| File pattern | Responsibility |
|---|---|
| `{feature}.go` | Handler struct, `module.Module` implementation (Info, Commands, etc.) |
| `{feature}_component.go` | HandleComponent dispatch for buttons/selects |
| `{feature}_modal.go` | HandleModal for modal submissions |
| `{feature}_listener.go` | Event listeners (MessageCreate, ReactionAdd, etc.) |

### internal/service/ (per feature)

| File pattern | Responsibility |
|---|---|
| `{feature}.go` | Service struct, core business logic |
| `{feature}_*.go` | Additional service files for complex features (poller, lavalink, etc.) |

### internal/view/ (per feature)

| File pattern | Responsibility |
|---|---|
| `{feature}.go` | Primary UI builder functions |
| `{feature}_*.go` | Additional view files (settings, queue, helpers, etc.) |

### internal/model/

| File | Responsibility |
|------|----------------|
| `constants.go` | All ModuleID constants (`PlayerModuleID`, `EmbedFixModuleID`, etc.) |
| `guild.go` | GuildSettings, Ticket, RSSFeed (shared data types) |
| `player.go` | Queue, QueueManager, LoopMode, PlayerSettings |
| `embedfix.go` | Platform, EmbedRef, EmbedFixSettings, URL matchers |
| `twitter.go` / `reddit.go` / `tiktok.go` | API response types for embed display |
| `ticket.go` / `logger.go` | Per-feature settings types |
| `panel.go` / `url.go` / `fuckfetch.go` / `translator.go` | Per-feature domain types |

## Data Flow

### Command Interaction

```
1. User types /player
2. Discord Gateway → disgo
3. internal/bot/router.go: onCommandInteraction()
4. Match command name → module (handler struct)
5. Check module enabled → handler.HandleCommand(e)
6. handler/player.go: call service, build UI via view.BuildPlayerUI()
7. Respond with Components V2 message
```

### Component Interaction

```
1. User clicks ⏭ (skip button)
2. Discord Gateway → disgo
3. internal/bot/router.go: onComponentInteraction()
4. Parse CustomID "player:skip" → moduleID="player"
5. Dispatch to handler.HandleComponent(e)
6. handler/player_component.go: switch "skip" → service.Skip()
7. service/player.go: queue.Next() + player.Update()
8. handler builds updated UI via view.BuildPlayerUI()
9. Respond with updated message
```

### Voice / Lavalink

```
1. User adds a track via modal
2. handler/player_modal.go → service.LoadAndPlay()
3. service/player.go: JoinVoiceChannel() → bot joins VC
4. Discord sends VoiceState/VoiceServer events
5. internal/bot/voice.go relays to disgolink
6. Lavalink streams audio
7. service/player_lavalink.go: onTrackEnd → play next
```

### Event Listeners (Logger)

```
1. User deletes a message
2. Discord Gateway → disgo
3. handler/logger_listener.go: onMessageDelete()
4. Check module enabled, load settings via service
5. view/logger_message.go: LoggerMessageDeleteLog()
6. Send log to configured channel
```

## Components V2 Design

All UI uses Discord Components V2 (`discord.NewMessageCreateV2()`). No embeds. No accent colors on containers.

### Layout Hierarchy

```
MessageCreate (V2 flag)
 └── ContainerComponent
      ├── TextDisplayComponent (markdown)
      ├── SeparatorComponent (divider)
      ├── MediaGalleryComponent (images)
      ├── SectionComponent (text + thumbnail/button accessory)
      └── ActionRowComponent (buttons, select menus)
```
