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

Pedmin connects to Discord via the Gateway (WebSocket) using disgo. Voice playback is handled by Lavalink, connected through disgolink. Both services run as Docker containers.

## Layered Feature Module Pattern

Each feature is a self-contained Go package with internal layer separation by file:

```
┌─────────────────────────────────────────────────┐
│  main.go (Entrypoint: DI wiring)                │
└────────────────┬────────────────────────────────┘
                 │
┌────────────────▼────────────────────────────────┐
│  bot/ (Framework: Discord connection, routing)   │
└────────────────┬────────────────────────────────┘
                 │
┌────────────────▼────────────────────────────────┐
│  features/*/                                     │
│    module.go           ← Module interface glue   │
│    handler_*.go        ← Request handling        │
│    service.go          ← Business logic          │
│    view_*.go           ← UI building             │
│    (domain files)      ← Data structures         │
└────────────────┬────────────────────────────────┘
                 │
┌────────────────▼────────────────────────────────┐
│  store/ (Infrastructure: SQLite persistence)     │
└─────────────────────────────────────────────────┘
```

All files within a feature share the same Go package. No sub-packages, no circular imports. File names indicate the layer and responsibility.

## Package Dependency Graph

```
main
 ├── config
 ├── bot
 │    ├── module   (interface only)
 │    └── store    (interface only)
 └── features/
      ├── settings   → module
      ├── ping       → module
      ├── avatar     → module
      ├── fuckfetch  → module
      ├── player     → module, disgolink
      ├── ticket     → module, store
      ├── logger     → module, store
      └── rss        → module, store
```

Dependencies flow downward. Features never depend on each other. Modules that need bot functionality (settings, ticket, logger, rss) define a local `Bot` interface with only the methods they need, avoiding direct `bot` package imports.

## File Responsibilities (1 File = 1 Concern)

### bot/
| File | Responsibility |
|------|---------------|
| `bot.go` | Client init, module registry, Start/Close, module state checks |
| `commands.go` | Slash command global sync |
| `router.go` | Interaction dispatch to modules |
| `ui.go` | Shared UI helpers (error messages) |
| `voice.go` | VoiceState/VoiceServer event relay to Lavalink |
| `presence.go` | Bot presence updater (CPU/RAM monitoring) |

### features/player/
| File | Layer | Responsibility |
|------|-------|---------------|
| `module.go` | Module | Info, Commands, empty stubs |
| `handler_command.go` | Handler | `/player` slash command |
| `handler_component.go` | Handler | Button/select switch + delegation |
| `handler_modal.go` | Handler | Add-to-queue modal processing |
| `service.go` | Service | Playback control (pause, skip, volume, track loading) |
| `voice.go` | Service | Voice channel connection |
| `queue.go` | Domain | Queue data structure |
| `queue_manager.go` | Domain | Per-guild queue management |
| `loop_mode.go` | Domain | LoopMode type definition |
| `lavalink.go` | Infra | Lavalink event listeners, node connection |
| `auto_leave.go` | Service | Auto-leave on empty VC |
| `view_player.go` | View | Player UI builder |
| `view_queue.go` | View | Queue list UI builder |
| `view_helpers.go` | View | Progress bar, duration format, thumbnails |

### features/settings/
| File | Layer | Responsibility |
|------|-------|---------------|
| `module.go` | Module | Info, Commands, Bot interface, empty stubs |
| `handler.go` | Handler | Command/component routing |
| `view.go` | View | Main panel, module panel builders |

### features/ticket/
| File | Layer | Responsibility |
|------|-------|---------------|
| `module.go` | Module | Info, Bot/Client/Store deps |
| `handler_component.go` | Handler | Create/close/reopen ticket buttons |
| `handler_modal.go` | Handler | Ticket creation modal |
| `service.go` | Service | Ticket creation/closure logic |
| `settings.go` | Domain | Settings struct & persistence |
| `view_panel.go` | View | Ticket control panel UI |
| `view_ticket.go` | View | Ticket channel message UI |
| `view_log.go` | View | Ticket list/log UI |

### features/logger/
| File | Layer | Responsibility |
|------|-------|---------------|
| `module.go` | Module | Info, Bot/Client/Store deps |
| `listener.go` | Handler | Event listeners (messages, members, bans, roles, channels) |
| `handler.go` | Handler | Component interaction handling |
| `settings.go` | Domain | Logger settings (channel ID, event toggles) |
| `view_settings.go` | View | Settings UI |
| `view_log.go` | View | Log message builders (text, attachments, MediaGallery) |

### features/rss/
| File | Layer | Responsibility |
|------|-------|---------------|
| `module.go` | Module | Info, Bot/Client/Store deps |
| `handler_component.go` | Handler | Add/remove feed UI |
| `handler_modal.go` | Handler | Feed URL input modal |
| `service.go` | Service | Feed fetch & post logic |
| `poller.go` | Infra | Background polling routine |
| `view_settings.go` | View | Settings panel (feed count) |
| `view_manage.go` | View | Add/remove feed UI |
| `view_feed.go` | View | Feed item announcement builder |

## Data Flow

### Command Interaction
```
1. User types /player
2. Discord Gateway → disgo
3. bot/router.go: onCommandInteraction()
4. Match command name → module
5. Check module enabled → module.HandleCommand(e)
6. handler_command.go: build UI via view_player.go
7. Respond with Components V2 message
```

### Component Interaction
```
1. User clicks ⏭ (skip button)
2. Discord Gateway → disgo
3. bot/router.go: onComponentInteraction()
4. Parse CustomID "player:skip" → moduleID="player"
5. Dispatch to player.HandleComponent(e)
6. handler_component.go: switch "skip" → service.go: handleSkip()
7. service.go: queue.Next() + player.Update()
8. Respond with updated UI via view_player.go
```

### Voice / Lavalink
```
1. User adds a track via modal
2. handler_modal.go → service.go: loadAndPlay()
3. voice.go: ensureVoiceConnection() → bot joins VC
4. Discord sends VoiceState/VoiceServer events
5. bot/voice.go relays to disgolink
6. Lavalink streams audio
7. lavalink.go: onTrackEnd → service.go: play next
```

### Event Listeners (Logger)
```
1. User deletes a message
2. Discord Gateway → disgo
3. listener.go: onMessageDelete()
4. Check module enabled, load settings
5. view_log.go: BuildMessageDeleteLog() with text + attachments
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
