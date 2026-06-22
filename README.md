# GitHub Analytics

GitHubのチーム活動を定量的に分析し、Web上で可視化するツールです（Findy Teams 風のチーム分析）。
バッチがGitHubからデータを収集して**スナップショット**としてPostgresに蓄積し、GraphQL API + React SPA が
最新スナップショットを「チーム概要 / メンバー横断ランキング・比較 / リポジトリ軸」の3つの切り口で表示します。

## 機能

- **チーム概要（Team Summary）**: チーム全体の合計・集計値
- **メンバー横断のランキング / 比較**: メンバー間で比較可能なスカラー指標（並び替え・順位付け・比較は**フロントエンドで計算**）
- **メンバー個別のドリルダウン**: 年次推移、関与リポジトリ TOP など
- **リポジトリ軸の横断集計**: 全リポジトリ（TOP3 に限定しない）ごとの集計とコントリビュータ内訳
- **時間軸の推移（期間指定）**: 任意の日付範囲（日単位）で絞り込み、日 / 週 / 月の粒度で時系列推移グラフを表示（チーム概要・メンバー詳細）

扱う指標の一覧は [指標（v1）](./docs/architecture.md#指標v1) を参照してください。

## クイックスタート

```bash
# 1. ツールチェーンと依存関係をインストール
mise install
go mod download
make web-install

# 2. .env を用意して GITHUB_TOKEN を設定
cp .env.example .env

# 3. Postgres を起動し、バッチで最初のスナップショットを作成
make db-up
export $(grep -v '^#' .env | xargs)
make batch ARGS="-users user1,user2"

# 4. Web サーバを起動（既定 :8090）
make serve
```

詳細な手順は [セットアップ](./docs/getting-started.md) と [使い方](./docs/usage.md) を参照してください。

## ドキュメント

- [アーキテクチャ](./docs/architecture.md) — ディレクトリ構成 / スナップショット方式 / API・フロントエンド構成 / 指標
- [セットアップ](./docs/getting-started.md) — 前提条件 / 環境変数 / GitHub PAT スコープ
- [使い方](./docs/usage.md) — バッチ実行 / Web サーバ / docker-compose / API 制限
- [開発](./docs/development.md) — 開発フロー / コード生成 / ビルド / テスト / lint / **ドキュメント更新ルール**

> **ドキュメント更新ルール**: コードや仕様を変更したら、同じ変更内で関連ドキュメントを必ず更新します。
> 詳細は [開発ドキュメントのルール](./docs/development.md#ドキュメント更新ルール必須) を参照してください。

## ライセンス

このプロジェクトはMITライセンスの下で公開されています。

## 貢献

プルリクエストやイシューの報告を歓迎します。

1. 機能ブランチを作成（`git checkout -b feature/amazing-feature`）
2. 変更をコミット（**関連ドキュメントの更新を含めること**。[ドキュメント更新ルール](./docs/development.md#ドキュメント更新ルール必須)）
3. ブランチに push
4. プルリクエストを作成
