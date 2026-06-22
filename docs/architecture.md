# アーキテクチャ

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
├── docs/                      # プロジェクトドキュメント（本ディレクトリ）
├── docker-compose.yml         # Postgres + Web アプリ
├── Dockerfile                 # SPAビルド → Goサーバへ frontend/dist を埋め込み
├── gqlgen.yml                 # gqlgen 設定
├── Makefile / mise.toml       # 開発タスク
├── go.mod / go.sum
└── README.md
```

## データの蓄積方法（スナップショット）

バッチ実行 1 回 = 1 スナップショット（`captured_at`）。スナップショットごとに集計済みメトリクスを保存します
（メンバー単位のスカラー、メンバー × 年、メンバー × 日、メンバー × リポジトリ（全リポジトリ））。Web はデフォルトで
**最新スナップショット**を読み込みます。

メンバー × 日（`MemberDayStat`）は活動を `YYYY-MM-DD`（UTC 基準で丸めた日）単位に集計したもので、
任意の日付範囲での絞り込みと時系列推移グラフのデータ源になります。日付範囲フィルタと週 / 月へのバケット集約は
ランキング・比較と同様に**フロントエンドで計算**します。

## ストレージ / API / フロントエンド

- **ストレージ**: PostgreSQL（Docker）。ORM は ent、ドライバは pgx（stdlib アダプタ）
- **API**: gqlgen による GraphQL（スキーマファースト）。主なクエリ:
  - `members: [MemberStats!]!` — メンバー横断の比較可能スカラー（ランキング・比較用）
  - `member(login: String!): UserStatistics` — ドリルダウン（年次推移・日次推移・TOPリポジトリ等）
  - `teamSummary: TeamSummary!` — チーム合計・集計
  - `teamDailyStats: [DailyStatistics!]!` — チーム全体の日次合計（日付昇順の時系列）
  - `repositories: [RepositoryStats!]!` — リポジトリ軸の横断集計
  - `repository(nameWithOwner: String!): RepositoryStats`
  - 並び替え / 順位付け / 比較・日付範囲の絞り込みは GraphQL ではなく**フロントエンドで計算**します。
- **フロントエンド**: React + Vite の SPA。パッケージマネージャは pnpm。GraphQL クライアントは urql、
  型は graphql-codegen（client preset）、チャートは Recharts。本番は Go バイナリが `frontend/dist` を
  埋め込み**同一オリジン**で配信し、開発時は Vite が `/query` を Go サーバへプロキシします。

## 指標（v1）

メンバー軸・リポジトリ軸の双方で次を扱います。

- コミット数
- Pull Request 作成数 / マージ数
- Issue 作成数
- Review 数（PRレビュー）
- 変更行数（additions / deletions、**PR由来のみ**。コミット単位の行数はAPIから取得できません）
- PR / Review 比率

これらの指標はメンバー軸・リポジトリ軸に加え、**時間軸**でも扱えます。チーム概要・メンバー詳細では、
任意の日付範囲（日単位）で絞り込み、日 / 週 / 月のいずれかの粒度で時系列推移グラフを表示できます。

> **v2 の予定（未実装）**: レビュー時間（review time）の指標は PR の timeline 取得が必要なため、
> 意図的に v1 では実装していません。

## 取得できないデータについて

- 組織の private リポジトリ（適切な権限がない場合）
- 削除されたリポジトリのデータ
- フォーク元リポジトリでの活動（一部）
- コミット単位の変更行数（API制限のため。行数は PR 由来のみ）
