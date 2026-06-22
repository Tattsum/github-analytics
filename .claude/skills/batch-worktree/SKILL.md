---
name: batch-worktree
description: github-analytics のバッチ（-mode batch / GitHub から収集して Postgres へ1スナップショット保存）を、使い捨ての git worktree 内で実行する。コミット済みコードのスナップショットから動くため、メインチェックアウトや別 worktree で進行中の作業に一切影響されず、並行してバッチを回せる。「バッチを worktree で回して」「作業と別にバッチ実行」「隔離してバッチ」等で発動。
---

# batch-worktree

進行中の開発作業（メインチェックアウトの未コミット編集や、別の git worktree）から**独立**して、github-analytics のバッチを実行するためのスキル。

## なぜ worktree か

`make batch` を普段の作業ツリーで直接実行すると、編集中のコードを巻き込んでビルド・実行してしまう。このスキルは指定した**コミット済み ref**（既定: `HEAD`）から使い捨ての worktree を切り、そこでバッチを実行する。これにより:

- バッチは安定したコードスナップショットから動く（編集中ファイルの影響を受けない）
- ユーザーはメインツリー／別 worktree で並行して作業を続けられる

隔離されるのは**コードだけ**。Postgres は共有で、バッチは従来どおり DB に1スナップショットを書き込む。

## 実行手順

スキルディレクトリの `run-batch.sh` を使う。バッチ対象（`-users` か `-org`、任意で `-team` / `-private`）は必ず指定する（未指定だとバッチ側が fatal で終了する）。

```bash
# 特定ユーザー
.claude/skills/batch-worktree/run-batch.sh -- -users user1,user2

# 組織のチーム + private リポジトリ含む
.claude/skills/batch-worktree/run-batch.sh -- -org myorg -team my-team -private

# 特定コミット/ブランチのコードでバッチを回す
.claude/skills/batch-worktree/run-batch.sh --ref main -- -org myorg
```

スクリプトがやること:

1. `--ref`（既定 `HEAD`）を short commit に解決し、`$TMPDIR` 配下にユニークな worktree を `--detach` で作成。
2. worktree は gitignore 済みファイルを持たないため、**origin リポジトリ root** の `.env` / `.envrc` から
   `GITHUB_TOKEN` / `DATABASE_URL` を読み込む（トークンは出力しない）。`GITHUB_TOKEN` が無ければエラー終了。
3. `docker compose up -d --wait postgres` で Postgres を冪等に起動（共有）。
4. worktree 内で `go run ./cmd/github-analytics -mode batch <args>` を実行。
5. `trap` で**必ず** worktree を `worktree remove --force` + `prune` で撤去（`--keep` でデバッグ用に残せる）。

## 長時間バッチ（並行作業向け）

バッチは最大30分かかりうる。ユーザーが作業を続けながら回したい場合は、Bash ツールの `run_in_background` でこのスクリプトを起動し、完了通知を待つ。worktree 隔離と合わせて、メイン作業を止めずにバッチを完走できる。

## 注意

- 対象指定（`-users` / `-org`）が無いとバッチは fatal 終了する。呼ぶ前に対象を確認すること。
- `--keep` で残した worktree は手動で `git worktree remove --force <path>` する。
- DB は共有のため、同時に複数バッチを走らせるとスナップショットが複数書かれる。隔離はコード単位であって DB 単位ではない点に留意。
