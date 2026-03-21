# Components V2 Reference (disgo v0.19.2)

## Component Types

### Layout Components (top-level, used in messages/modals)

| Type                    | Constructor                                                   | Use                           |
|-------------------------|---------------------------------------------------------------|-------------------------------|
| `ContainerComponent`    | `discord.NewContainer(subs...)`                               | Groups components             |
| `ActionRowComponent`    | `discord.NewActionRow(comps...)`                              | Holds buttons/selects (max 5) |
| `TextDisplayComponent`  | `discord.NewTextDisplay(content)`                             | Markdown text block           |
| `SectionComponent`      | `discord.NewSection(subs...)`                                 | Groups text with an accessory |
| `SeparatorComponent`    | `discord.NewLargeSeparator()` / `discord.NewSmallSeparator()` | Visual divider                |
| `MediaGalleryComponent` | `discord.NewMediaGallery(items...)`                           | Image gallery display         |
| `LabelComponent`        | `discord.NewLabel(label, comp)`                               | Labels for modals (V2)        |

### Interactive Components (inside ActionRow)

| Type                        | Constructor                                                      | Use              |
|-----------------------------|------------------------------------------------------------------|------------------|
| `ButtonComponent`           | `discord.NewPrimaryButton(label, customID)`                      | Clickable button |
| `StringSelectMenuComponent` | `discord.NewStringSelectMenu(customID, placeholder, options...)` | Dropdown select  |

### Button Styles

```go
discord.NewPrimaryButton(label, customID)   // Blue
discord.NewSecondaryButton(label, customID) // Gray
discord.NewSuccessButton(label, customID)   // Green
discord.NewDangerButton(label, customID)    // Red
discord.NewLinkButton(label, url)           // Gray with link icon
```

### Accessory Components (inside Section)

| Type                 | Constructor                 | Use              |
|----------------------|-----------------------------|------------------|
| `ThumbnailComponent` | `discord.NewThumbnail(url)` | Small image      |
| `ButtonComponent`    | (same as above)             | Button accessory |

### Media Components

| Type                | Constructor               | Use                          |
|---------------------|---------------------------|------------------------------|
| `MediaGalleryItem`  | struct with `Media` field | Single media item in gallery |
| `UnfurledMediaItem` | struct with `URL` field   | Media URL reference          |

## Creating V2 Messages

### New Message

```go
msg := discord.NewMessageCreateV2(
    discord.NewContainer(
        discord.NewTextDisplay("## Title"),
        discord.NewLargeSeparator(),
        discord.NewTextDisplay("Content here"),
        discord.NewActionRow(
            discord.NewPrimaryButton("Click me", "mymod:action"),
        ),
    ),
)
```

### Ephemeral Message

```go
msg := discord.NewMessageCreateV2(components...).WithEphemeral(true)
```

### Update Message (for component interactions)

```go
update := discord.NewMessageUpdateV2([]discord.LayoutComponent{
    discord.NewContainer(
        discord.NewTextDisplay("Updated content"),
    ),
})
```

## MediaGallery

Display images using `MediaGallery` with `MediaGalleryItem`:

```go
discord.NewMediaGallery(
    discord.MediaGalleryItem{
        Media: discord.UnfurledMediaItem{URL: "https://cdn.example.com/image.png"},
    },
    discord.MediaGalleryItem{
        Media: discord.UnfurledMediaItem{URL: "https://cdn.example.com/image2.png"},
    },
)
```

Used in the avatar module for avatar display and in the logger module for attachment logging.

## Responding to Interactions

### Command Response

```go
func (m *MyMod) HandleCommand(e *events.ApplicationCommandInteractionCreate) {
    // Immediate response
    _ = e.CreateMessage(discord.NewMessageCreateV2(components...))

    // Or deferred (shows "thinking...")
    _ = e.DeferCreateMessage(true) // true = ephemeral
    // ... do work, then follow up via REST
}
```

### Component Response

```go
func (m *MyMod) HandleComponent(e *events.ComponentInteractionCreate) {
    // Update the message the component is on
    _ = e.UpdateMessage(discord.NewMessageUpdateV2(components))

    // Or create a new ephemeral message
    _ = e.CreateMessage(discord.NewMessageCreateV2(components...).WithEphemeral(true))

    // Or just acknowledge (no visible change)
    _ = e.DeferUpdateMessage()
}
```

### Modal Response

```go
_ = e.Modal(discord.ModalCreate{
    CustomID: "mymod:my_modal",
    Title:    "My Modal",
    Components: []discord.LayoutComponent{
        discord.NewLabel("Field Name",
            discord.NewShortTextInput("mymod:field").
                WithPlaceholder("Enter value").
                WithRequired(true),
        ),
    },
})
```

## UI Patterns in This Project

### Admin Panel (settings)

Container with select menu for module list, detail view with toggle button and back button.

### Media Player (player)

Container with Section for track info + thumbnail, progress bar text, two ActionRows for controls.

### List View (queue)

Container with numbered track list as TextDisplay, navigation buttons.

### Log Messages (logger)

Container with title, separator, body text. Image attachments displayed via MediaGallery, non-image files listed as
text.

### Avatar Display (avatar)

Container with MediaGallery showing server and/or global avatar.

## Section with Accessory

```go
discord.NewSection(
    discord.NewTextDisplay("### Title"),
    discord.NewTextDisplay("Subtitle text"),
).WithAccessory(discord.NewThumbnail("https://example.com/image.png"))
```

## Ephemeral vs Channel Messages

| Type          | Use When                                                           |
|---------------|--------------------------------------------------------------------|
| **Ephemeral** | Settings, error messages, confirmations - only the user should see |
| **Channel**   | Player UI, announcements - everyone should see                     |
