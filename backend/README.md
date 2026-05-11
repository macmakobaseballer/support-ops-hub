# Support Ops Hub — backend

Go モジュラーモノリス。サービスは `cmd/{gateway,auth,ticket}` を独立バイナリとして起動するが、Phase 1 では同一プロセス内のモジュール関数として相互呼び出しする（[CLAUDE.md](CLAUDE.md) ルール 1）。

## ディレクトリ

```
backend/
├── cmd/
│   ├── gateway/    API Gateway (M3)
│   ├── auth/       auth-service (M2)
│   └── ticket/     ticket-service (M4)
├── services/       (M2 以降) サービス別ハンドラ・ドメインロジック
├── gateway/        (M3) リバプロ + JWT ミドルウェア
└── internal/
    ├── apigen/     (M1) oapi-codegen 生成（コミット対象）
    ├── db/         (M1) migrations / queries / sqlc 生成
    ├── domain/     (M2) ENUM とドメインルール
    ├── httperr/    (M2) エラー形式統一 (CLAUDE.md ルール 4)
    ├── middleware/ (M2) request ID / slog / recover
    ├── auth/jwt/   (M2) JWT 発行・検証
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
