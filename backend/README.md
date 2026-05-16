# Support Ops Hub — backend

Go モジュラーモノリス。サービスは `cmd/{gateway,auth,ticket}` を独立バイナリとして起動するが、Phase 1 では同一プロセス内のモジュール関数として相互呼び出しする（[CLAUDE.md](CLAUDE.md) ルール 1）。

## ディレクトリ

```
backend/
├── cmd/
│   ├── gateway/    API Gateway (M2: dev-auth / M7: JWT 検証)
│   ├── auth/       auth-service (M7: bcrypt + JWT 発行)
│   └── ticket/     ticket-service (M3)
├── services/       (M3 以降) サービス別ハンドラ・ドメインロジック
├── gateway/        (M2) リバプロ + dev-auth → M7 で JWT ミドルウェアに差替
└── internal/
    ├── apigen/     (M1) oapi-codegen 生成（コミット対象）
    ├── db/         (M1) migrations / queries / sqlc 生成
    ├── domain/     (M2) ENUM とドメインルール
    ├── httperr/    (M2) エラー形式統一 (CLAUDE.md ルール 4)
    ├── middleware/ (M2) request ID / slog / recover
    ├── auth/jwt/   (M7) JWT 発行・検証（M2 では未使用）
    └── config/     (M2) envconfig
```

## 開発

`backend/` 単体ではなく、リポジトリルートの `make` ターゲットを使用する。

```bash
make up              # MySQL 起動（リポジトリルート）
make db-migrate      # goose で 8 テーブルを適用（.env の DB_DSN）
make backend-build   # go build ./...
make backend-vet     # go vet ./...
make codegen         # sqlc + oapi-codegen 再生成
```

マイグレーションは [pressly/goose](https://github.com/pressly/goose) を使用する。`internal/db/migrations/` の SQL は `-- +goose Up` / `-- +goose Down` 形式。適用状況は `make db-migrate-status` で確認できる。

## マイルストーン進捗

**M1 完了**。主な成果物：

| 成果物 | パス |
|--------|------|
| DB マイグレーション（8テーブル） | `internal/db/migrations/` |
| sqlc クエリ定義 | `internal/db/queries/` |
| sqlc 生成（Go 型・クエリ関数） | `internal/db/sqlc/` |
| oapi-codegen 生成（OpenAPI モデル） | `internal/apigen/models.gen.go` |

`go build ./...` / `go vet ./...` / `go test ./...` は全パス。

M2 以降の実装先：
- `internal/domain/` — ENUM 型エイリアス・ドメインルール
- `internal/httperr/` — エラー形式統一
- `internal/middleware/` — リクエスト ID / ロギング / recover
- `internal/config/` — envconfig
- `gateway/` — dev-auth → M7 で JWT に差替
