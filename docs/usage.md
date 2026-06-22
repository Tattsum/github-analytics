# 使い方

セットアップが未了の場合は [セットアップ](./getting-started.md) を先に実施してください。

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

## docker-compose で一括起動

Postgres と Web アプリ（SPA ビルド + Go サーバ）をまとめて起動します。

```bash
# .env を用意したうえで
docker compose up -d --build

# アプリは http://localhost:${APP_PORT:-8080} で配信される
# スナップショットの作成は別途バッチ実行が必要（compose 外、またはワンショットで）
make batch ARGS="-org myorganization"
```

## API 制限について

- GitHub GraphQL API の rate limit は 5000 リクエスト/時です
- 本ツールは rate limit を考慮して実装されており、自動的に待機します
- 大量データの場合は処理に時間がかかります（バッチは最大30分のタイムアウト）

取得できないデータの詳細は [アーキテクチャ](./architecture.md#取得できないデータについて) を参照してください。
