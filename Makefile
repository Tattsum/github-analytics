.PHONY: help build test test-coverage lint lint-fix lint-markdown lint-markdown-fix fmt clean install-tools run

# 変数定義
BINARY_NAME=github-analytics
CMD_PATH=./cmd/github-analytics
COVERAGE_FILE=coverage.out
COVERAGE_HTML=coverage.html

# デフォルトターゲット
help: ## このヘルプメッセージを表示
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## アプリケーションをビルド
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BINARY_NAME) $(CMD_PATH)
	@echo "Build complete: $(BINARY_NAME)"

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
	@rm -f $(BINARY_NAME)
	@rm -f $(COVERAGE_FILE) $(COVERAGE_HTML)
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

.DEFAULT_GOAL := help
