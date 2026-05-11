SHELL := /bin/bash

# ===========================================================================
# Support Ops Hub — top-level Makefile (M0)
# More targets are added in later milestones (M1: migrate/sqlc/oapi, M2+: run-*)
# ===========================================================================

.PHONY: help up down logs ps backend-build backend-vet front-install front-dev \
        tf-fmt tf-validate-dev clean

help:
	@echo "Support Ops Hub — make targets"
	@echo ""
	@echo "  up               docker compose で MySQL/LocalStack を起動"
	@echo "  down             docker compose を停止 (ボリュームは残す)"
	@echo "  logs             コンテナログを追従"
	@echo "  ps               コンテナ状態を表示"
	@echo ""
	@echo "  backend-build    go build ./..."
	@echo "  backend-vet      go vet ./..."
	@echo ""
	@echo "  front-install    pnpm install"
	@echo "  front-dev        pnpm dev (http://localhost:3000)"
	@echo ""
	@echo "  tf-fmt           terraform fmt -recursive (チェックのみ)"
	@echo "  tf-validate-dev  dev 環境の terraform validate"
	@echo ""
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

# --- cleanup ---

clean:
	docker compose down -v
	rm -rf backend/tmp frontend/.nuxt frontend/.output frontend/node_modules
