# Components V2 リファレンス (disgo v0.19.2)

## コンポーネントの種類

### レイアウトコンポーネント (トップレベル、メッセージ/モーダルで使用)

| 型                      | コンストラクタ                                                | 用途                          |
|-------------------------|---------------------------------------------------------------|-------------------------------|
| `ContainerComponent`    | `discord.NewContainer(subs...)`                               | コンポーネントをグループ化    |
| `ActionRowComponent`    | `discord.NewActionRow(comps...)`                              | ボタン/セレクトを配置 (最大5) |
| `TextDisplayComponent`  | `discord.NewTextDisplay(content)`                             | Markdown テキストブロック     |
| `SectionComponent`      | `discord.NewSection(subs...)`                                 | テキストとアクセサリをグループ化 |
| `SeparatorComponent`    | `discord.NewLargeSeparator()` / `discord.NewSmallSeparator()` | 視覚的な区切り線              |
| `MediaGalleryComponent` | `discord.NewMediaGallery(items...)`                           | 画像ギャラリー表示            |
| `LabelComponent`        | `discord.NewLabel(label, comp)`                               | モーダル用ラベル (V2)         |

### インタラクティブコンポーネント (ActionRow 内)

| 型                          | コンストラクタ                                                   | 用途               |
|-----------------------------|------------------------------------------------------------------|--------------------|
| `ButtonComponent`           | `discord.NewPrimaryButton(label, customID)`                      | クリック可能なボタン |
| `StringSelectMenuComponent` | `discord.NewStringSelectMenu(customID, placeholder, options...)` | ドロップダウン選択   |

### ボタンスタイル

```go
discord.NewPrimaryButton(label, customID)   // 青
discord.NewSecondaryButton(label, customID) // グレー
discord.NewSuccessButton(label, customID)   // 緑
discord.NewDangerButton(label, customID)    // 赤
discord.NewLinkButton(label, url)           // グレー (リンクアイコン付き)
```

### アクセサリコンポーネント (Section 内)

| 型                   | コンストラクタ                | 用途                   |
|----------------------|-------------------------------|------------------------|
| `ThumbnailComponent` | `discord.NewThumbnail(url)`   | 小さい画像             |
| `ButtonComponent`    | (上記と同じ)                  | ボタンアクセサリ       |

### メディアコンポーネント

| 型                  | コンストラクタ                  | 用途                             |
|---------------------|---------------------------------|----------------------------------|
| `MediaGalleryItem`  | `Media` フィールドを持つ構造体  | ギャラリー内の単一メディアアイテム |
| `UnfurledMediaItem` | `URL` フィールドを持つ構造体    | メディア URL 参照                |

## V2 メッセージの作成

### 新規メッセージ

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

### エフェメラルメッセージ

```go
msg := discord.NewMessageCreateV2(components...).WithEphemeral(true)
```

### メッセージの更新 (コンポーネントインタラクション用)

```go
update := discord.NewMessageUpdateV2([]discord.LayoutComponent{
    discord.NewContainer(
        discord.NewTextDisplay("Updated content"),
    ),
})
```

## MediaGallery

`MediaGallery` と `MediaGalleryItem` を使って画像を表示する。

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

avatar モジュールのアバター表示や、logger モジュールの添付ファイルログに使用されている。

## インタラクションへの応答

### コマンド応答

```go
func (m *MyMod) HandleCommand(e *events.ApplicationCommandInteractionCreate) {
    // 即時応答
    _ = e.CreateMessage(discord.NewMessageCreateV2(components...))

    // または遅延応答 (「考え中...」が表示される)
    _ = e.DeferCreateMessage(true) // true = エフェメラル
    // ... 処理を行い、REST 経由でフォローアップ
}
```

### コンポーネント応答

```go
func (m *MyMod) HandleComponent(e *events.ComponentInteractionCreate) {
    // コンポーネントが配置されているメッセージを更新
    _ = e.UpdateMessage(discord.NewMessageUpdateV2(components))

    // または新しいエフェメラルメッセージを作成
    _ = e.CreateMessage(discord.NewMessageCreateV2(components...).WithEphemeral(true))

    // または確認応答のみ (見た目の変化なし)
    _ = e.DeferUpdateMessage()
}
```

### モーダル応答

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

## 本プロジェクトの UI パターン

### 管理パネル (settings)

Container にモジュール一覧のセレクトメニュー、詳細ビューにはトグルボタンと戻るボタンを配置。

### メディアプレイヤー (player)

Container 内に Section でトラック情報 + サムネイル、プログレスバーテキスト、2つの ActionRow でコントロールを配置。

### リストビュー (queue)

Container 内にナンバリングされたトラックリストを TextDisplay で表示し、ナビゲーションボタンを配置。

### ログメッセージ (logger)

Container にタイトル、セパレータ、本文テキストを配置。画像添付は MediaGallery で表示し、画像以外のファイルはテキストとして一覧表示。

### アバター表示 (avatar)

Container 内の MediaGallery でサーバーアバターおよびグローバルアバターを表示。

## Section とアクセサリ

```go
discord.NewSection(
    discord.NewTextDisplay("### Title"),
    discord.NewTextDisplay("Subtitle text"),
).WithAccessory(discord.NewThumbnail("https://example.com/image.png"))
```

## エフェメラル vs チャンネルメッセージ

| 種類            | 使用場面                                                         |
|-----------------|------------------------------------------------------------------|
| **エフェメラル** | 設定、エラーメッセージ、確認 - そのユーザーだけに表示すべき場合   |
| **チャンネル**   | プレイヤー UI、アナウンス - 全員に表示すべき場合                  |
