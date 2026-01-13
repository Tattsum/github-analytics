---
name: github-analytics
description: GitHub Analytics プロジェクトの開発ガイドラインとベストプラクティス
metadata:
  short-description: クリーンアーキテクチャとDDDに基づくGitHub活動分析ツールの開発ルール
  version: 1.0.0
  author: GitHub Analytics Team
---

# GitHub Analytics プロジェクトスキル

このスキルは、GitHub Analyticsプロジェクトの開発におけるアーキテクチャ原則、コーディング規約、ベストプラクティスを定義します。

## プロジェクト概要

GitHub Analyticsは、ソフトウェアエンジニアのGitHub活動を定量的に分析し、送別会用の「数字による軌跡」を生成するツールです。

## アーキテクチャ原則

### クリーンアーキテクチャ

このプロジェクトはクリーンアーキテクチャの原則に厳密に従います：

1. **依存関係の方向**: 外側の層は内側の層に依存する。内側の層は外側の層に依存しない。
   - `domain` → 最内層（ビジネスロジック、エンティティ、値オブジェクト）
   - `application` → ユースケース層（ドメイン層に依存）
   - `infrastructure` → インフラ層（外部API、データベースなど）
   - `presentation` → プレゼンテーション層（出力フォーマットなど）
   - `cmd` → エントリーポイント（すべての層を使用）

2. **インターフェースの分離**: 各層はインターフェースを通じて通信する。
   - ドメイン層にインターフェースを定義
   - インフラ層がインターフェースを実装

3. **依存性逆転の原則**: 抽象に依存し、具象に依存しない。

### ドメイン駆動設計（DDD）

1. **エンティティ**: 一意のIDを持つオブジェクト（例: `User`）
2. **値オブジェクト**: 値によって識別されるオブジェクト（例: `Activity`, `YearlyStatistics`）
3. **ドメインサービス**: エンティティや値オブジェクトに属さないビジネスロジック
4. **リポジトリパターン**: データアクセスの抽象化

## コーディング規約

### 命名規則

- **パッケージ名**: 小文字、単数形、簡潔（例: `domain`, `application`）
- **型名**: パスカルケース（例: `UserStatistics`, `ActivityType`）
- **関数名**: パスカルケース（公開）、キャメルケース（非公開）
- **定数**: パスカルケース（例: `ActivityTypeCommit`）
- **エラー変数**: `err` または `Err` で始まる

### ファイル構造

- 1ファイル = 1主要な型（関連する型は同じファイルに配置可能）
- テストファイルは `*_test.go` 形式
- パッケージのドキュメントコメントを必ず記述

### エラーハンドリング

- エラーは必ず処理する（`_` で無視しない）
- エラーメッセージは文脈を含める（`fmt.Errorf("failed to fetch user: %w", err)`）
- カスタムエラー型は `domain/errors.go` に定義

### テスト駆動開発（TDD）

1. **レッド**: 失敗するテストを書く
2. **グリーン**: テストを通す最小限の実装
3. **リファクタ**: コードを改善

### テスト方針

- **ユニットテスト**: 各層のロジックを独立してテスト
- **統合テスト**: 層間の連携をテスト（`integration` タグを使用）
- **モック**: インターフェースを使用してモックを生成
- **カバレッジ**: 最低80%を目標

## レイヤー別ガイドライン

### Domain層 (`domain/`)

- **責務**: ビジネスロジック、エンティティ、値オブジェクト、ドメインサービス
- **依存**: なし（標準ライブラリのみ）
- **テスト**: ビジネスロジックのユニットテスト
- **例**: `User`, `Activity`, `UserStatistics`

### Application層 (`application/`)

- **責務**: ユースケースの実装、ドメイン層の協調
- **依存**: `domain` のみ
- **テスト**: ユースケースのテスト（モックを使用）
- **例**: `StatisticsService`

### Infrastructure層 (`infrastructure/`)

- **責務**: 外部API、データベース、ファイルシステムなど
- **依存**: `domain` のみ（インターフェース経由）
- **テスト**: モックまたは統合テスト
- **例**: `GitHubClient`, `GitHubRepository`

### Presentation層 (`presentation/`)

- **責務**: 出力フォーマット、UI、CLI
- **依存**: `domain`, `application`
- **テスト**: 出力形式のテスト
- **例**: `OutputFormatter`

## コードレビュー基準

1. **アーキテクチャ**: 依存関係の方向が正しいか
2. **テスト**: 十分なテストカバレッジがあるか
3. **エラーハンドリング**: 適切にエラーを処理しているか
4. **命名**: 意図が明確か
5. **ドキュメント**: 公開APIにコメントがあるか

## 禁止事項

1. **循環依存**: パッケージ間の循環依存は禁止
2. **直接依存**: 外側の層から内側の層への直接依存は禁止
3. **グローバル変数**: 可能な限り避ける（設定は依存性注入で）
4. **panic**: エラーハンドリングでpanicは使用しない（リカバリ可能な場合を除く）

## 推奨事項

1. **インターフェース**: テスト容易性のため、インターフェースを積極的に使用
2. **依存性注入**: コンストラクタで依存関係を注入
3. **不変性**: 値オブジェクトは可能な限り不変にする
4. **早期リターン**: ネストを減らすため早期リターンを活用

## ツールとコマンド

- **ビルド**: `make build`
- **テスト**: `make test`
- **カバレッジ**: `make test-coverage`
- **リント**: `make lint` (Goコード)
- **Markdownリント**: `make lint-markdown` (Markdownファイル)
- **JSON/YAMLリント**: `make lint-json-yaml` (JSON/YAMLファイル)
- **GitHub Actionsリント**: `make lint-github-actions` (GitHub Actionsワークフロー)
- **フォーマット**: `make fmt` (Goコード + Markdown + JSON/YAML)
- **全チェック**: `make check` (フォーマット、リント、markdownlint、JSON/YAMLリント、GitHub Actionsリント、vet、テスト)

## Markdown品質管理

- **markdownlint**: すべてのMarkdownファイルは`markdownlint-cli2`で検証される
- **設定ファイル**: `.markdownlint.json`でルールを定義
- **自動修正**: `make lint-markdown-fix`で自動修正可能な問題を修正
- **必須チェック**: `make check`と`make ci`で自動的に実行される
- **Markdownファイル作成時**: 必ず`make lint-markdown`で検証する

## JSON/YAML品質管理

- **prettier**: すべてのJSON/YAMLファイルは`prettier`で検証される
- **設定ファイル**: `.prettierrc.json`でルールを定義
- **自動修正**: `make format-json-yaml`で自動修正可能な問題を修正
- **必須チェック**: `make check`と`make ci`で自動的に実行される
- **JSON/YAMLファイル作成時**: 必ず`make lint-json-yaml`で検証する

## GitHub Actions品質管理

- **actionlint**: すべてのGitHub Actionsワークフローファイル（`.github/workflows/*.yml`）は`actionlint`で検証される
- **自動検証**: `make lint-github-actions`で実行する
- **必須チェック**: `make check`と`make ci`で自動的に実行される
- **ワークフローファイル作成・変更時**: 必ず`make lint-github-actions`で検証する

## 実装時の注意点

### GitHub API使用時

- Rate limitを考慮した実装を必ず行う
- 並列処理は適切に制御する（最大5並列）
- エラーハンドリングは必ず実装する
- タイムアウトを適切に設定する（30分）

### データ取得時

- 組織のメンバー取得には`read:org`スコープが必要
- Privateリポジトリには`repo`スコープが必要
- ページネーションを適切に処理する

### 出力生成時

- 個人名を含めない（公開可能な形式）
- エラーハンドリングを適切に実装
- ファイル作成前にディレクトリの存在確認

### Markdownファイル作成時

- **必須**: すべてのMarkdownファイルは`make lint-markdown`で検証する
- コードブロックには必ず言語指定を付ける（例: `` ```bash ``、`` ```go ``、`` ```text ``）
- リストの前後には空行を入れる
- `.markdownlint.json`のルールに従う
- `make lint-markdown-fix`で自動修正可能な問題を修正する

## 参考資料

- [Clean Architecture by Robert C. Martin](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Domain-Driven Design by Eric Evans](https://www.domainlanguage.com/ddd/)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
