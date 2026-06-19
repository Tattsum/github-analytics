# GitHub Analytics

GitHubのチーム活動を定量的に分析し、Web上で可視化するツールです（Findy Teams 風のチーム分析）。
バッチがGitHubからデータを収集して**スナップショット**としてPostgresに蓄積し、GraphQL API + React SPA が
最新スナップショットを「チーム概要 / メンバー横断ランキング・比較 / リポジトリ軸」の3つの切り口で表示します。

## 機能

- **チーム概要（Team Summary）**: チーム全体の合計・集計値
- **メンバー横断のランキング / 比較**: メンバー間で比較可能なスカラー指標（並び替え・順位付け・比較は**フロントエンドで計算**）
- **メンバー個別のドリルダウン**: 年次推移、関与リポジトリ TOP など
- **リポジトリ軸の横断集計**: 全リポジトリ（TOP3 に限定しない）ごとの集計とコントリビュータ内訳

### 指標（v1）

メンバー軸・リポジトリ軸の双方で次を扱います。

- コミット数
- Pull Request 作成数 / マージ数
- Issue 作成数
- Review 数（PRレビュー）
- 変更行数（additions / deletions、**PR由来のみ**。コミット単位の行数はAPIから取得できません）
- PR / Review 比率

> **v2 の予定（未実装）**: レビュー時間（review time）の指標は PR の timeline 取得が必要なため、
> 意図的に v1 では実装していません。

## アーキテクチャ

Clean Architecture を維持しています。ドメイン層は純粋（インフラ非依存）で、ent / DB コードは
`infrastructure/` 配下にのみ存在します。`domain` / `application` はリポジトリ・サービスの
**インターフェース**にのみ依存し、ent を直接参照しません。

```text
github-analytics/
├── cmd/
│   ├── github-analytics/      # CLI: ファイル出力モード + バッチモード（Postgresへスナップショット保存）
│   └── server/                # Webサーバ: GraphQL API + 埋め込みSPA配信
├── domain/                    # ドメインモデル（純粋。インフラ非依存）
├── application/               # ユースケース・統計計算サービス、Snapshot 型
├── infrastructure/            # GitHub API クライアント / フェッチャー
│   └── ent/                   # ent ORM（生成コード + schema）。DBアクセスはここに限定
│       └── snapshotdb/        # スナップショットの読み書き（SnapshotWriter / SnapshotReader）
├── presentation/              # ファイル出力フォーマッター（CLI file モード用）
├── graph/                     # gqlgen: GraphQLスキーマ（*.graphqls）・生成コード・リゾルバ
├── frontend/                  # React + Vite SPA（urql + graphql-codegen + Recharts）
├── docker-compose.yml         # Postgres + Web アプリ
├── Dockerfile                 # SPAビルド → Goサーバへ frontend/dist を埋め込み
├── gqlgen.yml                 # gqlgen 設定
├── Makefile / mise.toml       # 開発タスク
├── go.mod / go.sum
└── README.md
```

### データの蓄積方法（スナップショット）

バッチ実行 1 回 = 1 スナップショット（`captured_at`）。スナップショットごとに集計済みメトリクスを保存します
（メンバー単位のスカラー、メンバー × 年、メンバー × リポジトリ（全リポジトリ））。Web はデフォルトで
**最新スナップショット**を読み込みます。

### ストレージ / API / フロントエンド

- **ストレージ**: PostgreSQL（Docker）。ORM は ent、ドライバは pgx（stdlib アダプタ）
- **API**: gqlgen による GraphQL（スキーマファースト）。主なクエリ:
  - `members: [MemberStats!]!` — メンバー横断の比較可能スカラー（ランキング・比較用）
  - `member(login: String!): UserStatistics` — ドリルダウン（年次推移・TOPリポジトリ等）
  - `teamSummary: TeamSummary!` — チーム合計・集計
  - `repositories: [RepositoryStats!]!` — リポジトリ軸の横断集計
  - `repository(nameWithOwner: String!): RepositoryStats`
  - 並び替え / 順位付け / 比較は GraphQL ではなく**フロントエンドで計算**します。
- **フロントエンド**: React + Vite の SPA。パッケージマネージャは pnpm。GraphQL クライアントは urql、
  型は graphql-codegen（client preset）、チャートは Recharts。本番は Go バイナリが `frontend/dist` を
  埋め込み**同一オリジン**で配信し、開発時は Vite が `/query` を Go サーバへプロキシします。

## セットアップ

### 前提条件

- [mise](https://mise.jdx.dev/) でツールチェーンを管理（Go 1.26.4 / Node LTS / pnpm）
- Docker / Docker Compose（PostgreSQL 用）
- GitHub Personal Access Token

```bash
# ツールチェーン（Go / Node / pnpm）をインストール
mise install

# Goの依存関係を取得
go mod download

# フロントエンドの依存関係を取得
make web-install   # = cd frontend && pnpm install
```

### 環境変数

`.env.example` を `.env` にコピーして値を埋めてください（docker-compose とアプリが読み込みます）。

| 変数 | 用途 |
| --- | --- |
| `GITHUB_TOKEN` | バッチのフェッチに使う GitHub Personal Access Token（read系スコープ） |
| `DATABASE_URL` | Postgres 接続文字列。compose 内では host が `postgres`、ホストからは `localhost` |
| `POSTGRES_USER` / `POSTGRES_PASSWORD` / `POSTGRES_DB` / `POSTGRES_PORT` | compose の Postgres サービス設定 |
| `APP_PORT` | アプリの公開ポート（compose） |
| `PORT` | Web サーバのリッスンポート（既定 8080） |
| `ENV` | `development` / `dev` のとき `GET /playground` を公開。本番は未公開 |

```bash
cp .env.example .env
# .env を編集して GITHUB_TOKEN を設定
```

### GitHub Personal Access Token のスコープ

- `public_repo`（公開リポジトリ）
- `repo`（privateリポジトリも対象にする場合）
- `read:org`（組織のメンバーを取得する場合）

## バッチの実行（スナップショット保存）

バッチは GitHub からデータを取得・集計し、1 スナップショットを Postgres に書き込みます。冪等で、起動時に
マイグレーションを実行します。`GITHUB_TOKEN` と `DATABASE_URL` が必要です。

```bash
# 1. Postgres を起動
make db-up

# 2. 環境変数を読み込み（.env を利用する場合）
export $(grep -v '^#' .env | xargs)

# 3. バッチを実行（メンバー指定）
make batch ARGS="-users user1,user2"

# 組織の全メンバーを対象にする場合
make batch ARGS="-org myorganization"

# privateリポジトリも含める場合
make batch ARGS="-users user1,user2 -private"
```

`make batch` は内部で `go run ./cmd/github-analytics -mode batch <ARGS>` を実行します。
直接実行する場合は次の通りです。

```bash
GITHUB_TOKEN=... DATABASE_URL=... go run ./cmd/github-analytics -mode batch -users user1,user2
```

> CLI には従来の `file` モード（`output/` にJSON/CSV/テキストを出力）も残っています。
> `-mode file`（既定）で利用でき、Postgres は不要です。

## Web サーバの実行

サーバは最新スナップショットを読み込み、`POST /query` で GraphQL API を、`/` で埋め込み SPA を配信します。
`DATABASE_URL` が必要で、起動時にマイグレーションを実行します。

```bash
# Postgres を起動し、最低 1 回バッチでスナップショットを作成しておくこと
make db-up
make batch ARGS="-users user1,user2"

# サーバを起動（既定 :8080）
make serve

# 開発時に GraphQL playground を使う場合
ENV=development make serve   # GET /playground が公開される
```

### docker-compose で一括起動

Postgres と Web アプリ（SPA ビルド + Go サーバ）をまとめて起動します。

```bash
# .env を用意したうえで
docker compose up -d --build

# アプリは http://localhost:${APP_PORT:-8080} で配信される
# スナップショットの作成は別途バッチ実行が必要（compose 外、またはワンショットで）
make batch ARGS="-org myorganization"
```

## 開発

ツールチェーンは mise（`mise.toml`）、フロントエンドの依存は pnpm で管理します。
開発時は Vite と Go サーバを別々のターミナルで起動し、Vite が `/query` を Go サーバ（:8080）へプロキシします。

```bash
# ターミナル1: Go サーバ（GraphQL API）
make serve

# ターミナル2: Vite 開発サーバ（http://localhost:5173 など）
make dev
```

`mise run <task>` でも同等のタスクを実行できます（`batch` / `serve` / `codegen` / `dev` / `db-up`）。

### コード生成（ent + gqlgen + フロントエンド）

スキーマやリゾルバ、GraphQL クエリを変更したら再生成します。

```bash
# すべて再生成（ent ORM → gqlgen サーバコード → フロントエンドの型）
make codegen

# 個別に実行する場合
make codegen-ent       # infrastructure/ent の ORM コード（go generate）
make codegen-gql       # graph/ の gqlgen サーバコード
make codegen-frontend  # frontend/src/gql の型付きGraphQLクライアント
```

### ビルド

```bash
# CLI（バッチ）をビルド
make build

# Web サーバをビルド（フロントエンドをビルドして frontend/dist を埋め込む）
make build-server
```

### テストの実行

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

### 開発ツールのインストール

```bash
make install-tools
```

golangci-lint v2 / actionlint / markdownlint-cli2 / prettier などをインストールします
（markdownlint-cli2 には Node.js が必要です）。

### コードフォーマット / リント

```bash
# go fmt + golangci-lint fmt + 各種自動修正
make fmt

# Goコードのリント（自動修正は make lint-fix）
make lint

# Markdownのリント（自動修正は make lint-markdown-fix）
make lint-markdown
```

### その他のコマンド

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

## 注意事項

### API制限について

- GitHub GraphQL API の rate limit は 5000 リクエスト/時です
- 本ツールは rate limit を考慮して実装されており、自動的に待機します
- 大量データの場合は処理に時間がかかります（バッチは最大30分のタイムアウト）

### 取得できないデータについて

- 組織の private リポジトリ（適切な権限がない場合）
- 削除されたリポジトリのデータ
- フォーク元リポジトリでの活動（一部）
- コミット単位の変更行数（API制限のため。行数は PR 由来のみ）

## AIエージェント対応

このプロジェクトは Cursor Skills と Claude Skills に対応しています。

- **`.cursorrules`**: Cursor エディタ用のプロジェクトルール
- **`.cursor/skills/`**: Cursor Skills 形式のスキル定義
- **`.claude/skills/`**: Claude Skills 形式のスキル定義

これらにより、AIエージェントがアーキテクチャ原則・コーディング規約・ベストプラクティスを理解し、
適切なコード生成や提案を行います。

## ライセンス

このプロジェクトはMITライセンスの下で公開されています。

## CI/CD

このプロジェクトは GitHub Actions を使用して CI/CD パイプラインを実行します。

- **テスト**: すべてのテストを実行し、カバレッジレポートを生成
- **リント**: Go コードと Markdown ファイルのリントを実行
- **ビルド**: アプリケーションのビルドを確認
- **全チェック**: フォーマット、リント、テストをすべて実行

ワークフローファイルは `.github/workflows/` ディレクトリにあります。

## 貢献

プルリクエストやイシューの報告を歓迎します。

1. 機能ブランチを作成（`git checkout -b feature/amazing-feature`）
2. 変更をコミット
3. ブランチに push
4. プルリクエストを作成
