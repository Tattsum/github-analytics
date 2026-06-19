.PHONY: help build build-server test test-coverage lint lint-fix lint-markdown lint-markdown-fix fmt clean install-tools run
.PHONY: db-up db-down batch serve dev web-install web-build codegen codegen-ent codegen-gql codegen-frontend

# 変数定義
BINARY_NAME=github-analytics
SERVER_BINARY=server
CMD_PATH=./cmd/github-analytics
SERVER_PATH=./cmd/server
COVERAGE_FILE=coverage.out
COVERAGE_HTML=coverage.html
# バッチ・サーバが参照するPostgres接続文字列。.env やシェルで上書き可能。
DATABASE_URL?=postgres://github_analytics:github_analytics@localhost:5432/github_analytics?sslmode=disable

# デフォルトターゲット
help: ## このヘルプメッセージを表示
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## CLI（バッチ）をビルド
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BINARY_NAME) $(CMD_PATH)
	@echo "Build complete: $(BINARY_NAME)"

build-server: web-build ## Webサーバをビルド（frontend/dist を埋め込み）
	@echo "Building $(SERVER_BINARY)..."
	@go build -o $(SERVER_BINARY) $(SERVER_PATH)
	@echo "Build complete: $(SERVER_BINARY)"

test: ## テストを実行
	@echo "Running tests..."
	@go test -v -race -coverprofile=$(COVERAGE_FILE) ./...

test-coverage: test ## テストを実行し、カバレッジレポートを生成
	@echo "Generating coverage report..."
	@go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@go tool cover -func=$(COVERAGE_FILE)
	@echo "Coverage report generated: $(COVERAGE_HTML)"

test-short: ## 短いテストを実行（統合テストをスキップ）
	@echo "Running short tests..."
	@go test -v -short ./...

lint: ## golangci-lintを実行
	@echo "Running golangci-lint..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found. Run 'make install-tools' first."; \
		exit 1; \
	fi

lint-fix: ## golangci-lintを実行し、自動修正可能な問題を修正
	@echo "Running golangci-lint with auto-fix..."
	@golangci-lint run --fix

lint-markdown: ## markdownlintを実行
	@echo "Running markdownlint..."
	@if command -v markdownlint-cli2 >/dev/null 2>&1; then \
		markdownlint-cli2 '**/*.md' '!node_modules/**' '!vendor/**' '!output/**' '!.git/**'; \
	elif [ -f node_modules/.bin/markdownlint-cli2 ]; then \
		npx markdownlint-cli2 '**/*.md' '!node_modules/**' '!vendor/**' '!output/**' '!.git/**'; \
	else \
		echo "markdownlint-cli2 not found. Run 'make install-tools' first."; \
		exit 1; \
	fi

lint-json: ## JSONファイルのリントを実行
	@echo "Running JSON linter..."
	@if [ -f node_modules/.bin/prettier ]; then \
		npm run lint:json; \
	else \
		echo "prettier not found. Run 'make install-tools' first."; \
		exit 1; \
	fi

lint-yaml: ## YAMLファイルのリントを実行
	@echo "Running YAML linter..."
	@if [ -f node_modules/.bin/prettier ]; then \
		npm run lint:yaml; \
	else \
		echo "prettier not found. Run 'make install-tools' first."; \
		exit 1; \
	fi

lint-json-yaml: lint-json lint-yaml ## JSONとYAMLファイルのリントを実行

lint-github-actions: ## GitHub Actionsワークフローのリントを実行
	@echo "Running actionlint..."
	@if command -v actionlint >/dev/null 2>&1; then \
		actionlint; \
	elif [ -f $(go env GOPATH)/bin/actionlint ]; then \
		$(go env GOPATH)/bin/actionlint; \
	else \
		echo "actionlint not found. Run 'make install-tools' first."; \
		exit 1; \
	fi

lint-markdown-fix: ## markdownlintを実行し、自動修正可能な問題を修正
	@echo "Running markdownlint with auto-fix..."
	@if command -v markdownlint-cli2-fix >/dev/null 2>&1; then \
		markdownlint-cli2-fix '**/*.md' '!node_modules/**' '!vendor/**' '!output/**' '!.git/**'; \
	elif [ -f node_modules/.bin/markdownlint-cli2-fix ]; then \
		npx markdownlint-cli2-fix '**/*.md' '!node_modules/**' '!vendor/**' '!output/**' '!.git/**'; \
	else \
		echo "markdownlint-cli2-fix not found. Run 'make install-tools' first."; \
		exit 1; \
	fi

format-json: ## JSONファイルをフォーマット
	@echo "Formatting JSON files..."
	@if [ -f node_modules/.bin/prettier ]; then \
		npm run format:json; \
	else \
		echo "prettier not found. Run 'make install-tools' first."; \
		exit 1; \
	fi

format-yaml: ## YAMLファイルをフォーマット
	@echo "Formatting YAML files..."
	@if [ -f node_modules/.bin/prettier ]; then \
		npm run format:yaml; \
	else \
		echo "prettier not found. Run 'make install-tools' first."; \
		exit 1; \
	fi

format-json-yaml: format-json format-yaml ## JSONとYAMLファイルをフォーマット

fmt: ## コードをフォーマット
	@echo "Formatting code..."
	@go fmt ./...
	@golangci-lint fmt
	@golangci-lint run --fix || true
	@echo "Formatting Markdown files..."
	@$(MAKE) lint-markdown-fix || echo "markdownlint-fix failed or not available"
	@echo "Formatting JSON and YAML files..."
	@$(MAKE) format-json-yaml || echo "JSON/YAML formatting failed or not available"

clean: ## ビルド成果物とカバレッジファイルを削除
	@echo "Cleaning..."
	@rm -f $(BINARY_NAME) $(SERVER_BINARY)
	@rm -f $(COVERAGE_FILE) $(COVERAGE_HTML)
	@rm -rf frontend/dist
	@go clean -cache
	@echo "Clean complete"

install-tools: ## 開発ツールをインストール
	@echo "Installing development tools..."
	@go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
	@go install github.com/rhysd/actionlint/cmd/actionlint@latest
	@if command -v npm >/dev/null 2>&1; then \
		echo "Installing Node.js tools (markdownlint-cli2, prettier)..."; \
		npm install || echo "npm install failed"; \
	else \
		echo "npm not found. Please install Node.js to use markdownlint and prettier."; \
	fi
	@echo "Tools installed"

run: build ## アプリケーションをビルドして実行
	@echo "Running $(BINARY_NAME)..."
	@./$(BINARY_NAME)

deps: ## 依存関係を更新
	@echo "Updating dependencies..."
	@go get -u ./...
	@go mod tidy

deps-check: ## 依存関係の更新をチェック
	@echo "Checking for dependency updates..."
	@go list -u -m all

vet: ## go vetを実行
	@echo "Running go vet..."
	@go vet ./...

check: fmt lint lint-markdown lint-json-yaml lint-github-actions vet test ## フォーマット、リント、markdownlint、JSON/YAMLリント、GitHub Actionsリント、vet、テストをすべて実行

ci: install-tools check test-coverage ## CI用: ツールインストール、チェック、テスト、カバレッジ

# --- Web アプリ（Postgres / バッチ / サーバ / フロントエンド）----------------

db-up: ## docker-compose で Postgres を起動（ヘルスチェック待ち）
	@echo "Starting PostgreSQL via docker-compose..."
	@docker compose up -d --wait postgres
	@echo "PostgreSQL is ready."

db-down: ## Postgres を停止（データボリュームは保持）
	@echo "Stopping PostgreSQL..."
	@docker compose down

batch: ## GitHubから収集し、1スナップショットをPostgresへ保存（要 GITHUB_TOKEN / DATABASE_URL）
	@echo "Running batch (snapshot to Postgres)..."
	@DATABASE_URL="$(DATABASE_URL)" go run $(CMD_PATH) -mode batch $(ARGS)

serve: ## Webサーバを起動（最新スナップショットを配信、要 DATABASE_URL）
	@echo "Starting web server on :$${PORT:-8080}..."
	@DATABASE_URL="$(DATABASE_URL)" go run $(SERVER_PATH)

dev: ## フロントエンド開発サーバ（Vite）を起動。別ターミナルで `make serve` を実行すること
	@echo "Starting Vite dev server (proxies /query -> Go server on :8080)..."
	@echo "Run 'make serve' in another terminal to start the Go backend."
	@cd frontend && pnpm dev

web-install: ## フロントエンドの依存関係をインストール
	@echo "Installing frontend dependencies..."
	@cd frontend && pnpm install

web-build: ## フロントエンドをビルドして frontend/dist を生成
	@echo "Building frontend (frontend/dist)..."
	@cd frontend && pnpm install && pnpm codegen && pnpm build

codegen: codegen-ent codegen-gql codegen-frontend ## すべてのコード生成（ent + gqlgen + フロントエンド）を実行

codegen-ent: ## ent のORMコードを生成（infrastructure/ent）
	@echo "Generating ent code..."
	@go generate ./infrastructure/ent/...

codegen-gql: ## gqlgen のGraphQLサーバコードを生成（graph/）
	@echo "Generating gqlgen code..."
	@go run github.com/99designs/gqlgen generate

codegen-frontend: ## graphql-codegen で型付きGraphQLクライアントを生成（frontend/src/gql）
	@echo "Generating frontend GraphQL types..."
	@cd frontend && pnpm codegen

.DEFAULT_GOAL := help
