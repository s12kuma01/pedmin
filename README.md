# Pedmin

Pedmin は、既存の多機能BOTの完全な代替として開発された、ユーザー・開発者フレンドリーなOSS多機能BOTです。

## 説明

ProbotやMee6などのBOTの代替となるべくオープンソースで開発された完全無料なDiscord多機能BOT。
Pedminは、レイヤードアーキテクチャに重きを置いた、新しく最適化されたコードで書かれています。
すべての実装でEmbedではなくComponents V2を使用しており、既存のBOTよりも優れたUI/UXを実現しています。

## 招待

> [!WARNING]
> 現在開発中なため、不定期の頻繁な再起動があります。
https://discord.com/oauth2/authorize?client_id=1484236709611704571

## 機能

- **設定パネル** — `/settings` ギルドごとのモジュール有効化/無効化を管理する管理用UI
- **音楽プレイヤー** — `/player` Jockie Musicの代替として作成されたシンプルで高音質、使いやすい音楽プレイヤー
- **サポートチケット** — Ticket Toolの代替として作成されたサポートチケット機能、設定パネルから有効化できます
- **サーバーログ** — 既存の多機能BOTのloggerの代替として作成されたイベントログ機能
- **RSS フィード** — MonitoRSSの代替として作成されたRSSフィード監視・自動アナウンス機能
- **Embed Fix** — 既存のEmbed修正BOTの代替、現在Xのみ対応
- **サーバーパネル** — `/panel` Pelicanパネルの操作用 (限定されたユーザーのみ)
- **URL ツール** — `/url` URL短縮（x.gd）& 安全スキャン（VirusTotal）
- **アバター表示** — `/avatar` サーバー/グローバルアバターを MediaGallery で表示
- **システム情報** — `/fuckfetch` neofetchからインスパイアされたシステム情報表示

## 技術スタック

| 技術 | バージョン / ライブラリ |
|------|------------------------|
| 言語 | Go 1.26.1 |
| Discord | [disgo](https://github.com/disgoorg/disgo) v0.19.2 |
| Lavalink クライアント | [disgolink](https://github.com/disgoorg/disgolink) v3.1.0 |
| Lavalink サーバー | Lavalink 4 (Alpine) |
| データベース | SQLite ([modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite)) |
| 設定 | 環境変数 + TOML ([BurntSushi/toml](https://github.com/BurntSushi/toml)) |

## アーキテクチャ

各機能は独立した Feature Module として実装され、内部で handler/service/view レイヤーに分離されています。

```
main.go          # エントリポイント: DI 配線、グレースフルシャットダウン
config/          # 環境変数 + TOML 設定読み込み
bot/             # Discord 接続、インタラクションルーティング
module/          # Module インターフェース定義
features/        # 機能モジュール (11個)
store/           # GuildStore インターフェース + SQLite 実装
```

詳細は [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) を参照。

## ドキュメント

[`docs/`](docs/) ディレクトリに詳細ガイドがあります:

- [Architecture](docs/ARCHITECTURE.md) — レイヤード Feature Module 設計
- [Module Development](docs/MODULE_GUIDE.md) — モジュール作成ガイド
- [Components V2](docs/COMPONENTS_V2.md) — disgo V2 コンポーネントリファレンス
- [Lavalink Integration](docs/LAVALINK.md) — 音楽再生セットアップ
- [Data Store](docs/STORE.md) — SQLite データ永続化

## ライセンス

[MPL-2.0](LICENSE.md)
