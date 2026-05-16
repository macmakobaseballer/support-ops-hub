---
description: Support Ops Hub のローカル開発環境（MySQL・LocalStack・Go gateway・Nuxt フロントエンド）を一括起動する。「ローカルサーバーを起動して」「開発環境を立ち上げて」「/dev-server」「サーバー起動」と言われたら必ずこのスキルを使うこと。
allowed-tools: Bash
---

# dev-server — ローカル開発環境 起動スキル

**作業ディレクトリ**: プロジェクトルート（`/home/makoj/dev/support-ops-hub`）

必ずこの順序で実行する。前のステップが失敗したら次に進まず、エラー内容を報告する。

## Step 1: Docker Compose（MySQL / LocalStack）

```bash
docker compose up -d --wait
```

ヘルスチェックが通るまで待機する（最大2分）。既に起動中でも冪等に動作する。

## Step 2: DB マイグレーション

```bash
export PATH=$PATH:/home/makoj/.local/go/bin:/home/makoj/go/bin
set -a; [ -f .env ] && . .env; set +a
DB_DSN="${DB_DSN:-app:app@tcp(127.0.0.1:3306)/support_ops_hub?parseTime=true&loc=Asia%2FTokyo&charset=utf8mb4}"
goose -dir backend/internal/db/migrations mysql "$DB_DSN" up
```

`goose` が見つからない場合はインストールしてから再実行：
```bash
GOBIN=/home/makoj/go/bin go install github.com/pressly/goose/v3/cmd/goose@v3.24.1
```

## Step 3: バックエンド（gateway）をバックグラウンド起動

```bash
export PATH=$PATH:/home/makoj/.local/go/bin:/home/makoj/go/bin
set -a; [ -f .env ] && . .env; set +a
cd backend && go run ./cmd/gateway/... > /tmp/gateway.log 2>&1 &
echo "Gateway PID: $!"
```

- ポート: `8000`（`.env` の `GATEWAY_PORT=8000`）
- ログ: `tail -f /tmp/gateway.log`

## Step 4: フロントエンドをバックグラウンド起動

```bash
cd frontend && pnpm dev > /tmp/frontend.log 2>&1 &
echo "Frontend PID: $!"
```

- URL: `http://localhost:3000`
- ログ: `tail -f /tmp/frontend.log`

## Step 5: ヘルスチェック（最大 30 秒待機）

```bash
# バックエンド
for i in $(seq 1 15); do
  curl -sf -o /dev/null http://localhost:8000/ && echo "Backend ready" && break
  echo "Waiting for backend... ($i)"; sleep 2
done

# フロントエンド
for i in $(seq 1 15); do
  curl -sf -o /dev/null http://localhost:3000/ && echo "Frontend ready" && break
  echo "Waiting for frontend... ($i)"; sleep 2
done
```

バックエンドは `/` が 404 を返しても「ポート到達可能＝正常起動」と判断する（`/api/v1/...` 以外にルートがないため）。
フロントエンドは `/` が 302（ログインリダイレクト）を返せば正常。

## Step 6: 起動完了を報告

```
起動完了:
  フロントエンド : http://localhost:3000
  バックエンド   : http://localhost:8000
  ログ確認       : tail -f /tmp/gateway.log
                   tail -f /tmp/frontend.log
```

## エラー対処

| 症状 | 確認コマンド |
|------|------------|
| MySQL に接続できない | `docker compose ps` / `docker compose logs mysql` |
| gateway が起動しない | `cat /tmp/gateway.log` |
| フロントが起動しない | `cat /tmp/frontend.log` |
| ポートが既に使用中 | `ss -tlnp \| grep -E '8000\|3000'` |

### 初回セットアップで発生しやすい問題

**`.env` がない場合**  
`.env.example` からコピー後、`DB_DSN` 行をシングルクォートで囲む（`tcp(...)` の括弧が bash 構文エラーになるため）：
```bash
cp .env.example .env
# DB_DSN= の行をシングルクォートで囲んで保存
```

**フロントエンドが `Cannot find module './parser.linux-x64-gnu.node'` で起動しない**  
Linux/WSL 環境で oxc ネイティブバインディングが未インストール。`cd frontend` して実行：
```bash
pnpm add -D \
  "@oxc-parser/binding-linux-x64-gnu@0.129.0" \
  "@oxc-transform/binding-linux-x64-gnu@0.129.0" \
  "@oxc-minify/binding-linux-x64-gnu@0.129.0"
```

## 停止

```bash
pkill -f 'go run ./cmd/gateway'
pkill -f 'pnpm dev'
docker compose down
```
