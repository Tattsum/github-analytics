# セットアップ

## 前提条件

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

## 環境変数

`.env.example` を `.env` にコピーして値を埋めてください（docker-compose とアプリが読み込みます）。

| 変数 | 用途 |
| --- | --- |
| `GITHUB_TOKEN` | バッチのフェッチに使う GitHub Personal Access Token（read系スコープ） |
| `DATABASE_URL` | Postgres 接続文字列。compose 内では host が `postgres`、ホストからは `localhost` |
| `POSTGRES_USER` / `POSTGRES_PASSWORD` / `POSTGRES_DB` / `POSTGRES_PORT` | compose の Postgres サービス設定 |
| `APP_PORT` | アプリの公開ポート（compose） |
| `PORT` | Web サーバのリッスンポート（既定 8090） |
| `ENV` | `development` / `dev` のとき `GET /playground` を公開。本番は未公開 |

```bash
cp .env.example .env
# .env を編集して GITHUB_TOKEN を設定
```

## GitHub Personal Access Token のスコープ

- `public_repo`（公開リポジトリ）
- `repo`（privateリポジトリも対象にする場合）
- `read:org`（組織のメンバーを取得する場合）

次は [使い方（バッチ・Web サーバの実行）](./usage.md) を参照してください。
