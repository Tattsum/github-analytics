# GitHub Analytics

ソフトウェアエンジニアのGitHub活動を定量的に分析し、送別会用の「数字による軌跡」を生成するツールです。

## 機能

- **基本活動量の集計**
  - コミット数（年別・累計）
  - Pull Request作成数・マージ数
  - Issue作成数
  - Reviewコメント数（PRレビュー）

- **技術的インパクトの分析**
  - 変更行数（additions / deletions）
  - 最もコミット数の多いリポジトリ TOP3
  - 長期間関与しているリポジトリ（最初と最後の活動日時）

- **継続性・キャリアの変遷**
  - GitHub活動開始年
  - 年ごとの活動ピーク
  - ロール変化を想起させる指標（PR作成 → Review比率の変化）

## 出力形式

各ユーザーについて、以下の形式で出力されます：

1. **JSON形式** (`{username}_statistics.json`)
   - 機械可読な形式で全データを出力

2. **CSV形式** (`{username}_statistics.csv`)
   - スプレッドシートで分析可能な形式

3. **テキスト要約** (`{username}_summary.txt`)
   - 「この人を数字で表すと」
   - 「エンジニアとしての特徴」
   - 「役割の変化が読み取れるポイント」

4. **プレゼン用短文** (`{username}_presentation.txt`)
   - スライド1枚に使える短文（箇条書き3〜4行）

## セットアップ

### 前提条件

- Go 1.25以上
- GitHub Personal Access Token

### インストール

```bash
# 依存関係のインストール
go mod download

# ビルド
go build -o github-analytics ./cmd/github-analytics
```

### GitHub Personal Access Tokenの取得

1. GitHubにログイン
2. Settings → Developer settings → Personal access tokens → Tokens (classic)
3. "Generate new token (classic)"をクリック
4. 以下のスコープを選択：
   - `public_repo` (公開リポジトリ用)
   - `repo` (privateリポジトリも対象にする場合)
   - `read:org` (組織のメンバーを取得する場合)
5. トークンを生成し、安全な場所に保存

## 使用方法

### 基本的な使用方法

```bash
# 環境変数にトークンを設定
export GITHUB_TOKEN="your_token_here"

# 特定のユーザーを分析
./github-analytics -users user1,user2

# または
go run ./cmd/github-analytics -users user1,user2
```

### 組織のメンバーを分析する場合

```bash
export GITHUB_TOKEN="your_token_here"
# 組織の全メンバーを分析
./github-analytics -org myorganization
```

**注意**: 組織のメンバーを取得するには、GitHub Personal Access Tokenに`read:org`スコープが必要です。

### Privateリポジトリも対象にする場合

```bash
export GITHUB_TOKEN="your_token_here"
./github-analytics -users user1,user2 -private
```

### 出力先を指定する場合

```bash
./github-analytics -users user1,user2 -output ./results
```

### ヘルプの表示

```bash
./github-analytics -help
```

### 出力先の確認

実行後、指定した出力ディレクトリ（デフォルト: `output/`）に以下のファイルが生成されます：

```text
output/
├── user1_statistics.json
├── user1_statistics.csv
├── user1_summary.txt
├── user1_presentation.txt
├── user2_statistics.json
├── user2_statistics.csv
├── user2_summary.txt
└── user2_presentation.txt
```

## 注意事項

### API制限について

- GitHub GraphQL APIのrate limitは5000リクエスト/時です
- 本ツールはrate limitを考慮して実装されており、自動的に待機します
- 大量のデータがある場合は、処理に時間がかかる場合があります（最大30分のタイムアウト）

### 取得できないデータについて

以下のデータは、GitHub APIの制限により取得できない場合があります：

- 組織のprivateリポジトリ（適切な権限がない場合）
- 削除されたリポジトリのデータ
- フォーク元のリポジトリでの活動（一部）

### 並列処理について

- 複数ユーザーのデータは並列で取得されます
- 各ユーザー内でのリポジトリ取得も並列で実行されます（最大5並列）
- rate limitを考慮して実装されているため、安全に実行できます

### コマンドライン引数

- `-users`: 分析対象のGitHubユーザー名（カンマ区切り、例: `user1,user2`）
- `-org`: 分析対象のGitHub組織名（指定した場合、組織の全メンバーを分析）
- `-output`: 出力ディレクトリ（デフォルト: `output`）
- `-private`: privateリポジトリも対象にする（フラグを指定）
- `-help`: ヘルプを表示

**注意**: `-users`と`-org`のどちらか一方を指定する必要があります。

## プロジェクト構造

```text
github-analytics/
├── cmd/
│   └── github-analytics/
│       └── main.go              # メインアプリケーション
├── domain/                      # ドメインモデル
│   ├── user.go
│   ├── activity.go
│   └── statistics.go
├── infrastructure/              # インフラストラクチャ層
│   ├── github_client.go         # GitHub APIクライアント
│   ├── github_repository.go     # データ取得リポジトリ
│   └── github_data_fetcher.go  # データフェッチャー
├── application/                 # アプリケーション層
│   └── statistics_service.go   # 統計計算サービス
├── presentation/                # プレゼンテーション層
│   └── output_formatter.go     # 出力フォーマッター
├── .claude/
│   └── skills/                  # Claude Skills（AIエージェント用ルール）
│       └── github-analytics.md
├── .cursor/
│   └── skills/                  # Cursor Skills（AIエージェント用ルール）
│       └── github-analytics.md
├── .cursorrules                 # Cursor用プロジェクトルール
├── go.mod
├── go.sum
└── README.md
```

## AIエージェント対応

このプロジェクトは、Cursor SkillsとClaude Skillsに対応しています：

- **`.cursorrules`**: Cursorエディタ用のプロジェクトルール
- **`.cursor/skills/`**: Cursor Skills形式のスキル定義
- **`.claude/skills/`**: Claude Skills形式のスキル定義

これらのスキルファイルにより、AIエージェントがプロジェクトのアーキテクチャ原則、コーディング規約、ベストプラクティスを理解し、適切なコード生成や提案を行います。

### Cursorでの使用方法

CursorのNightlyリリースでは、「Import Agent Skills」設定を有効にすることで、`.claude/skills/`や`.cursor/skills/`ディレクトリ内のスキルを自動的に読み込みます。

### Claudeでの使用方法

Claude CodeやClaude Desktopでは、`.claude/skills/`ディレクトリ内のスキルが自動的に読み込まれ、プロジェクト固有のガイドラインに従った支援が行われます。

## 開発

### 開発ツールのインストール

```bash
make install-tools
```

これにより、以下のツールがインストールされます：

- golangci-lint v2 (Goコードのリント)
- markdownlint-cli2 (Markdownファイルのリント)

**注意**: markdownlint-cli2のインストールにはNode.js (v18以上)が必要です。

### テストの実行

```bash
# すべてのテストを実行
make test

# テストとカバレッジレポートを生成
make test-coverage

# 短いテストのみ実行（統合テストをスキップ）
make test-short
```

### コードフォーマット

```bash
# go fmt + golangci-lint fmt + 自動修正
make fmt
```

### リント

```bash
# Goコードのリントを実行
make lint

# Goコードのリントを実行し、自動修正可能な問題を修正
make lint-fix

# Markdownファイルのリントを実行
make lint-markdown

# Markdownファイルのリントを実行し、自動修正可能な問題を修正
make lint-markdown-fix
```

### その他のコマンド

```bash
# ビルド
make build

# すべてのチェック（フォーマット、リント、markdownlint、vet、テスト）
make check

# CI用（ツールインストール、チェック、テスト、カバレッジ）
make ci

# 依存関係の更新
make deps

# クリーンアップ
make clean
```

## ライセンス

このプロジェクトはMITライセンスの下で公開されています。

## CI/CD

このプロジェクトはGitHub Actionsを使用してCI/CDパイプラインを実行します：

- **テスト**: すべてのテストを実行し、カバレッジレポートを生成
- **リント**: GoコードとMarkdownファイルのリントを実行
- **ビルド**: アプリケーションのビルドを確認
- **全チェック**: フォーマット、リント、テストをすべて実行

ワークフローファイルは`.github/workflows/`ディレクトリにあります。

## 貢献

プルリクエストやイシューの報告を歓迎します。

### 開発フロー

1. 機能ブランチを作成（`git checkout -b feature/amazing-feature`）
2. 変更をコミット（`git commit -m 'Add amazing feature'`）
3. ブランチにpush（`git push origin feature/amazing-feature`）
4. プルリクエストを作成

### コミットメッセージ

- 明確で簡潔なメッセージを心がける
- 変更の理由を説明する
- 関連するIssue番号があれば記載する

### GitHubへの初回push

```bash
# リモートリポジトリを設定（GitHubでリポジトリを作成後）
git remote add origin https://github.com/Tattsum/github-analytics.git

# またはSSHを使用する場合
git remote add origin git@github.com:Tattsum/github-analytics.git

# メインブランチにpush
git push -u origin main
```
