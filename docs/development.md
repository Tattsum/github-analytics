# 開発

ツールチェーンは mise（`mise.toml`）、フロントエンドの依存は pnpm で管理します。
開発時は Vite と Go サーバを別々のターミナルで起動し、Vite が `/query` を Go サーバ（:8080）へプロキシします。

```bash
# ターミナル1: Go サーバ（GraphQL API）
make serve

# ターミナル2: Vite 開発サーバ（http://localhost:5173 など）
make dev
```

`mise run <task>` でも同等のタスクを実行できます（`batch` / `serve` / `codegen` / `dev` / `db-up`）。

## ドキュメント更新ルール（必須）

**コードや仕様を変更したら、同じ変更（PR / コミット）内で関連ドキュメントを必ず更新します。**
ドキュメントとコードの乖離はバグと同等に扱い、レビューでも指摘対象とします。

更新が必要になる代表的なケースと、対応するドキュメント:

| 変更内容 | 更新するドキュメント |
| --- | --- |
| ディレクトリ構成・レイヤー責務・スナップショット方式・指標の変更 | [`docs/architecture.md`](./architecture.md) |
| 前提ツール・環境変数・PAT スコープの変更 | [`docs/getting-started.md`](./getting-started.md) |
| CLI フラグ・バッチ / サーバ起動手順・compose の変更 | [`docs/usage.md`](./usage.md) |
| Make ターゲット・コード生成・テスト / lint 手順の変更 | 本ファイル（`docs/development.md`） |
| 上記いずれかで概要や入口が変わる場合 | [`README.md`](../README.md) |

チェックリスト:

- [ ] 追加・変更した挙動が、いずれかのドキュメントに反映されているか
- [ ] 削除した機能・フラグ・環境変数の記述がドキュメントから除去されているか
- [ ] サンプルコマンド・コードブロックが実際に動作する内容になっているか
- [ ] `make lint-markdown` が通るか

## コード生成（ent + gqlgen + フロントエンド）

スキーマやリゾルバ、GraphQL クエリを変更したら再生成します。

```bash
# すべて再生成（ent ORM → gqlgen サーバコード → フロントエンドの型）
make codegen

# 個別に実行する場合
make codegen-ent       # infrastructure/ent の ORM コード（go generate）
make codegen-gql       # graph/ の gqlgen サーバコード
make codegen-frontend  # frontend/src/gql の型付きGraphQLクライアント
```

## ビルド

```bash
# CLI（バッチ）をビルド
make build

# Web サーバをビルド（フロントエンドをビルドして frontend/dist を埋め込む）
make build-server
```

## テストの実行

```bash
# すべてのテストを実行
make test

# テストとカバレッジレポートを生成
make test-coverage

# 短いテストのみ実行（統合テストをスキップ）
make test-short

# フロントエンドのテスト
cd frontend && pnpm test
```

### 視覚リグレッションテスト（Playwright・ローカルのみ）

レスポンシブ対応の非劣化を担保するため、Playwright でページのスクリーンショットを比較します。
GraphQL レスポンスは `tests/visual` 配下の固定フィクスチャを `page.route` でモックするため、
Go バックエンドや Postgres の起動は不要です（CI には組み込みません）。

```bash
cd frontend

# 全 4 ページ × 3 幅（1280 / 768 / 375）＋ ハンバーガー展開を比較
pnpm test:visual

# 意図した見た目の変更を反映してベースラインを更新
pnpm test:visual:update
```

ベースライン画像は `tests/visual/__screenshots__/` にコミットします。
PC 幅（1280）のベースラインは emotion 移行前に取得した非劣化基準で、移行後に差分ゼロを確認済みです。

## 開発ツールのインストール

```bash
make install-tools
```

golangci-lint v2 / actionlint / markdownlint-cli2 / prettier などをインストールします
（markdownlint-cli2 には Node.js が必要です）。

## コードフォーマット / リント

```bash
# go fmt + go fix（2回連続）+ golangci-lint fmt + 各種自動修正
make fmt

# Goコードのリント（自動修正は make lint-fix）
make lint

# Markdownのリント（自動修正は make lint-markdown-fix）
make lint-markdown
```

`make fmt` は `go fix ./...` を **2回連続**で実行し、Goコードを最新イディオムへモダナイズする。
2回実行するのは、1回目の適用で生じた変化からさらに修正候補が出ることがあり、1回では収束しない場合があるため
（2回目で差分が出なくなる＝fixed point を確認する）。

## その他のコマンド

```bash
# すべてのチェック（フォーマット、リント、markdownlint、JSON/YAMLリント、GitHub Actionsリント、vet、テスト）
make check

# CI用（ツールインストール、チェック、テスト、カバレッジ）
make ci

# 依存関係の更新
make deps

# クリーンアップ（バイナリ・カバレッジ・frontend/dist を削除）
make clean

# Postgres を停止
make db-down
```

`make help` で全ターゲットの一覧を表示できます。

## CI/CD

このプロジェクトは GitHub Actions を使用して CI/CD パイプラインを実行します。

- **テスト**: push / PR ごとにすべてのテストを実行
- **リント**: Go コードと Markdown ファイルのリントを実行
- **ビルド**: アプリケーションのビルドを確認
- **全チェック**: フォーマット、リント、テストをすべて実行
- **カバレッジ**: 毎日 JST 10:00 に定期実行（`workflow_dispatch` で手動実行も可）。push / PR では実行しない

ワークフローファイルは `.github/workflows/` ディレクトリにあります。

## AIエージェント対応

このプロジェクトは Cursor Skills と Claude Skills に対応しています。

- **`.cursorrules`**: Cursor エディタ用のプロジェクトルール
- **`.cursor/skills/`**: Cursor Skills 形式のスキル定義
- **`.claude/skills/`**: Claude Skills 形式のスキル定義

これらにより、AIエージェントがアーキテクチャ原則・コーディング規約・ベストプラクティスを理解し、
適切なコード生成や提案を行います。
