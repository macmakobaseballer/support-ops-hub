SHELL := /bin/bash

# ===========================================================================
# Support Ops Hub — top-level Makefile
# ===========================================================================

.PHONY: help up down logs ps backend-build backend-vet front-install front-dev \
        tf-fmt tf-validate-dev clean \
        tools-install check-backend check-frontend check-infra check-ci \
        check-secrets check-deps check-all \
        sqlc oapi-codegen codegen db-migrate db-migrate-status

help:
	@echo "Support Ops Hub — make targets"
	@echo ""
	@echo "Docker:"
	@echo "  up               docker compose で MySQL/LocalStack を起動"
	@echo "  down             docker compose を停止 (ボリュームは残す)"
	@echo "  logs             コンテナログを追従"
	@echo "  ps               コンテナ状態を表示"
	@echo ""
	@echo "Backend:"
	@echo "  backend-build    go build ./..."
	@echo "  backend-vet      go vet ./..."
	@echo "  db-migrate       goose で MySQL マイグレーション適用 (.env の DB_DSN を使用)"
	@echo "  db-migrate-status  goose の適用状況を表示"
	@echo ""
	@echo "Frontend:"
	@echo "  front-install    pnpm install"
	@echo "  front-dev        pnpm dev (http://localhost:3000)"
	@echo ""
	@echo "Terraform:"
	@echo "  tf-fmt           terraform fmt -recursive (チェックのみ)"
	@echo "  tf-validate-dev  dev 環境の terraform validate"
	@echo ""
	@echo "Quality (品質チェック基盤):"
	@echo "  tools-install    品質チェックツール (golangci-lint v2/govulncheck/trivy/tflint/actionlint/zizmor) を ~/.local/bin / ~/go/bin に導入"
	@echo "  check-backend    backend: golangci-lint + go test -cover (race は CI で実行)"
	@echo "  check-frontend   frontend: lint + typecheck + vitest"
	@echo "  check-infra      terraform fmt/validate + tflint + trivy config"
	@echo "  check-ci         actionlint + zizmor (.github/workflows/)"
	@echo "  check-secrets    trivy fs --scanners secret (リポ全体)"
	@echo "  check-deps       govulncheck + pnpm audit (依存変更後に手動で)"
	@echo "  check-all        上記 5 つを順次 (check-deps は含まない)"
	@echo ""
	@echo "Cleanup:"
	@echo "  clean            生成物・ボリュームを削除"

# --- docker compose ---

up:
	docker compose up -d --wait

down:
	docker compose down

logs:
	docker compose logs -f --tail=100

ps:
	docker compose ps

# --- backend ---

backend-build:
	cd backend && go build ./...

backend-vet:
	cd backend && go vet ./...

# --- codegen ---

SQLC_VERSION         ?= v1.27.0
OAPI_CODEGEN_VERSION ?= v2.7.0

sqlc:
	cd backend && sqlc generate

oapi-codegen:
	oapi-codegen --config backend/internal/apigen/oapi-codegen.yaml \
	             docs/api/openapi.yaml

codegen: sqlc oapi-codegen
	@echo "codegen complete"

# --- database ---

GOOSE_VERSION ?= v3.24.1
DB_DSN        ?= app:app@tcp(127.0.0.1:3306)/support_ops_hub?parseTime=true&loc=Asia%2FTokyo&charset=utf8mb4

db-migrate:
	@set -a; [ -f .env ] && . ./.env; set +a; \
	goose -dir backend/internal/db/migrations mysql "$${DB_DSN:-$(DB_DSN)}" up

db-migrate-status:
	@set -a; [ -f .env ] && . ./.env; set +a; \
	goose -dir backend/internal/db/migrations mysql "$${DB_DSN:-$(DB_DSN)}" status

# --- frontend ---

front-install:
	cd frontend && pnpm install

front-dev:
	cd frontend && pnpm dev

# --- terraform ---

tf-fmt:
	cd infra/terraform && terraform fmt -recursive -check

tf-validate-dev:
	cd infra/terraform/envs/dev && terraform init -backend=false -input=false && terraform validate

# --- quality tools ---

TRIVY_VERSION ?= 0.70.0
ACTIONLINT_VERSION ?= 1.7.12
ZIZMOR_VERSION ?= 1.24.1
GOLANGCI_LINT_VERSION ?= v2.12.2

tools-install:
	@echo "==> Installing Go-based tools..."
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)
	go install golang.org/x/vuln/cmd/govulncheck@latest
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@$(SQLC_VERSION)
	go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@$(OAPI_CODEGEN_VERSION)
	go install github.com/pressly/goose/v3/cmd/goose@$(GOOSE_VERSION)
	@echo "==> Installing trivy v$(TRIVY_VERSION)..."
	mkdir -p $$HOME/.local/bin
	cd /tmp && curl -fsSL -o trivy.tar.gz https://github.com/aquasecurity/trivy/releases/download/v$(TRIVY_VERSION)/trivy_$(TRIVY_VERSION)_Linux-64bit.tar.gz \
		&& tar -xzf trivy.tar.gz trivy && mv trivy $$HOME/.local/bin/trivy && chmod +x $$HOME/.local/bin/trivy && rm trivy.tar.gz
	@echo "==> Installing tflint (latest)..."
	cd /tmp && curl -fsSL -o tflint.zip https://github.com/terraform-linters/tflint/releases/latest/download/tflint_linux_amd64.zip \
		&& python3 -c "import zipfile; zipfile.ZipFile('tflint.zip').extractall('.')" \
		&& mv tflint $$HOME/.local/bin/tflint && chmod +x $$HOME/.local/bin/tflint && rm tflint.zip
	@echo "==> Installing actionlint v$(ACTIONLINT_VERSION)..."
	cd /tmp && curl -fsSL -o actionlint.tar.gz https://github.com/rhysd/actionlint/releases/download/v$(ACTIONLINT_VERSION)/actionlint_$(ACTIONLINT_VERSION)_linux_amd64.tar.gz \
		&& tar -xzf actionlint.tar.gz actionlint && mv actionlint $$HOME/.local/bin/actionlint && chmod +x $$HOME/.local/bin/actionlint && rm actionlint.tar.gz
	@echo "==> Installing zizmor v$(ZIZMOR_VERSION)..."
	cd /tmp && curl -fsSL -o zizmor.tar.gz https://github.com/woodruffw/zizmor/releases/download/v$(ZIZMOR_VERSION)/zizmor-x86_64-unknown-linux-gnu.tar.gz \
		&& tar -xzf zizmor.tar.gz && mv zizmor $$HOME/.local/bin/zizmor && chmod +x $$HOME/.local/bin/zizmor && rm zizmor.tar.gz
	@echo "==> Done. Ensure ~/.local/bin and ~/go/bin are in PATH."

check-backend:
	cd backend && golangci-lint run ./...
	# -race は CGO 必須のためローカルでは省略。CI (.github/workflows/backend.yml)
	# が `go test -race -cover ./...` を必ず走らせる。
	cd backend && go test -cover ./...

check-frontend:
	cd frontend && pnpm lint && pnpm typecheck && pnpm test

check-infra:
	cd infra/terraform && terraform fmt -recursive -check
	cd infra/terraform/envs/dev && terraform init -backend=false -input=false && terraform validate
	cd infra/terraform && tflint --init && TFLINT_CONFIG_FILE=$$(pwd)/.tflint.hcl tflint --recursive
	trivy config infra/terraform/ --severity HIGH,CRITICAL

check-ci:
	actionlint .github/workflows/*.yml
	zizmor --min-severity=high .github/workflows/

check-secrets:
	trivy fs --scanners secret --severity HIGH,CRITICAL --skip-dirs node_modules,.nuxt,.terraform .

check-deps:
	cd backend && govulncheck ./...
	cd frontend && pnpm audit --prod --audit-level=high

check-all: check-backend check-frontend check-infra check-ci check-secrets

# --- cleanup ---

clean:
	docker compose down -v
	rm -rf backend/tmp frontend/.nuxt frontend/.output frontend/node_modules
