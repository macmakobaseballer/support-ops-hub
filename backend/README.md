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
make backend-build   # go build ./...
make backend-vet     # go vet ./...
```

## マイルストーン進捗

現在は M0（骨格のみ）。`cmd/*/main.go` は空エントリで `go build ./...` を通すことだけが目的。
