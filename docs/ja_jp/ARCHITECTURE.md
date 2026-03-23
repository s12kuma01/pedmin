# Architecture

## システム概要

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

Pedmin は disgo を使用して Gateway（WebSocket）経由で Discord に接続します。音声再生は disgolink を介して接続された Lavalink が処理します。両サービスは Docker コンテナとして実行されます。

## 標準的な Go レイヤードアーキテクチャ

本プロジェクトは標準的な Go 大規模サービスレイアウト（`cmd/`、`internal/`、`pkg/`）に準拠しています。レイヤーはファイル単位ではなく、**パッケージ**単位で分離されています。

```
┌──────────────────────────────────────────────────────┐
│  cmd/pedmin/main.go (エントリーポイント: DI配線)       │
└─────────────────────┬────────────────────────────────┘
                      │
┌─────────────────────▼────────────────────────────────┐
│  internal/bot/ (フレームワーク: 接続、ルーティング)     │
└─────────────────────┬────────────────────────────────┘
                      │
┌─────────────────────▼────────────────────────────────┐
│  internal/handler/    ← Discord インタラクションハンドラ│
│  internal/service/    ← ビジネスロジック                │
│  internal/view/       ← UIビルダー（純粋関数）          │
│  internal/model/      ← ドメイン型と定数               │
└─────────────────────┬────────────────────────────────┘
                      │
┌─────────────────────▼────────────────────────────────┐
│  internal/repository/ (インフラ: SQLite)               │
│  internal/client/     (インフラ: HTTPクライアント)      │
└──────────────────────────────────────────────────────┘
```

各レイヤーは独立した Go パッケージです。全12の機能モジュールが同じ handler、service、view、model パッケージを共有します。各パッケージ内のファイル名には機能名がプレフィックスとして付けられます（例: `handler/player.go`、`handler/embedfix.go`）。

## パッケージ依存関係グラフ

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
 ├── internal/module      → (末端: disgo 型のみ)
 ├── internal/model       → (末端: disgo, disgolink 型のみ)
 └── pkg/deepl            → (末端: 標準ライブラリのみ)
```

依存関係は下方向に流れます。機能同士が互いに依存することはありません。Bot の機能を必要とするモジュール（settings、embedfix、ticket、logger、rss）は、必要なメソッドのみを持つローカルな `Bot` インターフェースを定義し（ISP）、`bot` パッケージの直接インポートを回避しています。

## レイヤーの責務

| レイヤー | パッケージ | 責務 | 責務外 |
|---------|-----------|------|--------|
| **Handler** | `internal/handler/` | Discord イベントの受信、データ抽出、service 呼び出し、view レスポンスの構築 | ビジネスロジック、API呼び出し、ストアアクセスを含まない |
| **Service** | `internal/service/` | ビジネスロジック、API呼び出し、ストア操作 | Discord UIコンポーネントの構築を含まない |
| **View** | `internal/view/` | Components V2 UIの構築（純粋関数: データ入力 → コンポーネント出力） | 副作用、API呼び出し、ストアアクセスを含まない |
| **Model** | `internal/model/` | ドメイン型、設定構造体、ModuleID 定数 | handler/service/view/repository パッケージのインポートを含まない |
| **Repository** | `internal/repository/` | データ永続化インターフェースと SQLite 実装 | ビジネスロジックを含まない |
| **Client** | `internal/client/` | 外部API HTTPラッパー（Twitter、Reddit、TikTok、Pelican など） | ビジネスロジックを含まない |
| **UI** | `internal/ui/` | 共有 Discord UIヘルパー（EphemeralV2、ErrorMessage、BuildModulePanel） | 機能固有のロジックを含まない |
| **Bot** | `internal/bot/` | Discord クライアントのライフサイクル、モジュールレジストリ、インタラクションルーティング | 機能ロジックを含まない |

## ファイルの責務

### internal/bot/

| ファイル | 責務 |
|---------|------|
| `bot.go` | クライアント初期化、モジュールレジストリ、Start/Close、モジュール状態チェック |
| `commands.go` | スラッシュコマンドのグローバル同期 |
| `router.go` | インタラクションのモジュールへのディスパッチ |
| `voice.go` | VoiceState/VoiceServer イベントの Lavalink への中継 |
| `presence.go` | Bot プレゼンス更新（CPU/RAMモニタリング） |

### internal/handler/（機能ごと）

| ファイルパターン | 責務 |
|----------------|------|
| `{feature}.go` | Handler 構造体、`module.Module` の実装（Info、Commands など） |
| `{feature}_component.go` | ボタン/セレクトメニューの HandleComponent ディスパッチ |
| `{feature}_modal.go` | モーダル送信の HandleModal |
| `{feature}_listener.go` | イベントリスナー（MessageCreate、ReactionAdd など） |

### internal/service/（機能ごと）

| ファイルパターン | 責務 |
|----------------|------|
| `{feature}.go` | Service 構造体、コアビジネスロジック |
| `{feature}_*.go` | 複雑な機能用の追加サービスファイル（poller、lavalink など） |

### internal/view/（機能ごと）

| ファイルパターン | 責務 |
|----------------|------|
| `{feature}.go` | メインUIビルダー関数 |
| `{feature}_*.go` | 追加ビューファイル（settings、queue、helpers など） |

### internal/model/

| ファイル | 責務 |
|---------|------|
| `constants.go` | 全 ModuleID 定数（`PlayerModuleID`、`EmbedFixModuleID` など） |
| `guild.go` | GuildSettings、Ticket、RSSFeed（共有データ型） |
| `player.go` | Queue、QueueManager、LoopMode、PlayerSettings |
| `embedfix.go` | Platform、EmbedRef、EmbedFixSettings、URLマッチャー |
| `twitter.go` / `reddit.go` / `tiktok.go` | 埋め込み表示用 APIレスポンス型 |
| `ticket.go` / `logger.go` | 機能ごとの設定型 |
| `panel.go` / `url.go` / `fuckfetch.go` / `translator.go` | 機能ごとのドメイン型 |

## データフロー

### コマンドインタラクション

```
1. ユーザーが /player と入力
2. Discord Gateway → disgo
3. internal/bot/router.go: onCommandInteraction()
4. コマンド名をモジュール（handler 構造体）にマッチング
5. モジュール有効チェック → handler.HandleCommand(e)
6. handler/player.go: service を呼び出し、view.BuildPlayerUI() で UI を構築
7. Components V2 メッセージで応答
```

### コンポーネントインタラクション

```
1. ユーザーが ⏭（スキップボタン）をクリック
2. Discord Gateway → disgo
3. internal/bot/router.go: onComponentInteraction()
4. CustomID "player:skip" を解析 → moduleID="player"
5. handler.HandleComponent(e) にディスパッチ
6. handler/player_component.go: switch "skip" → service.Skip()
7. service/player.go: queue.Next() + player.Update()
8. handler が view.BuildPlayerUI() で更新された UI を構築
9. 更新されたメッセージで応答
```

### Voice / Lavalink

```
1. ユーザーがモーダル経由でトラックを追加
2. handler/player_modal.go → service.LoadAndPlay()
3. service/player.go: JoinVoiceChannel() → Bot がVCに参加
4. Discord が VoiceState/VoiceServer イベントを送信
5. internal/bot/voice.go が disgolink に中継
6. Lavalink がオーディオをストリーミング
7. service/player_lavalink.go: onTrackEnd → 次の曲を再生
```

### イベントリスナー（Logger）

```
1. ユーザーがメッセージを削除
2. Discord Gateway → disgo
3. handler/logger_listener.go: onMessageDelete()
4. モジュール有効チェック、service 経由で設定を読み込み
5. view/logger_message.go: LoggerMessageDeleteLog()
6. 設定されたチャンネルにログを送信
```

## Components V2 設計

全ての UI は Discord Components V2（`discord.NewMessageCreateV2()`）を使用します。Embed は使用しません。コンテナにアクセントカラーは設定しません。

### レイアウト階層

```
MessageCreate (V2 flag)
 └── ContainerComponent
      ├── TextDisplayComponent (markdown)
      ├── SeparatorComponent (divider)
      ├── MediaGalleryComponent (images)
      ├── SectionComponent (text + thumbnail/button accessory)
      └── ActionRowComponent (buttons, select menus)
```
