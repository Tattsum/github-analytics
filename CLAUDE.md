# CLAUDE.md

このファイルは Claude（および Claude Code）がこのリポジトリで作業する際のエントリーポイントです。

## プロジェクトルールの参照（必須）

**このプロジェクトのアーキテクチャ原則・コーディング規約・レイヤー別ガイドライン・禁止事項は
[`.cursorrules`](.cursorrules) に集約されています。** Claude は作業前に必ず `.cursorrules` を参照し、
その内容に従うこと。`CLAUDE.md` と `.cursorrules` は同じルールセットを指し、**お互いを参照します**。

- ルールの正本（アーキテクチャ・命名・エラーハンドリング・テスト方針・禁止事項など）→ [`.cursorrules`](.cursorrules)
- このファイル（`CLAUDE.md`）は Claude 向けの入口であり、`.cursorrules` を補完します
- どちらか一方を更新したら、もう一方の記述と矛盾しないように両方を確認すること

## ドキュメント

実行・開発手順は `docs/` 配下にまとまっています。

- [アーキテクチャ](docs/architecture.md) — ディレクトリ構成 / スナップショット方式 / API・フロントエンド構成 / 指標
- [セットアップ](docs/getting-started.md) — 前提条件 / 環境変数 / GitHub PAT スコープ
- [使い方](docs/usage.md) — バッチ実行 / Web サーバ / docker-compose / API 制限
- [開発](docs/development.md) — 開発フロー / コード生成 / ビルド / テスト / lint

## ドキュメント更新ルール（必須）

コードや仕様を変更したら、**同じ変更（PR / コミット）内で関連ドキュメントを必ず更新する**。
変更種別と更新先ドキュメントの対応表・チェックリストは
[`docs/development.md`](docs/development.md#ドキュメント更新ルール必須) を参照。
