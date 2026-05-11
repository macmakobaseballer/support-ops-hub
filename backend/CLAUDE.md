# Support Ops Hub — backend / CLAUDE.md

Go バックエンドで作業する際のルール。**これらは提案ではなくハードコンストレイントである。** 逸脱する前に必ずユーザーに確認すること。

リポジトリ全体のルールは [../CLAUDE.md](../CLAUDE.md) を参照。

---

## 構成

| 項目 | 内容 |
|------|------|
| 言語・バージョン | Go 1.26.3 |
| モジュール | `github.com/fixer/support-ops-hub/backend` |
| 設計 | モジュラーモノリス。**Phase 1 は単一バイナリにリンクして 1 ECS タスクで起動**し、サービス間は同一プロセス内のモジュール関数で相互呼び出し。Phase 2 で `cmd/<service>` を独立バイナリ・独立 ECS タスクに分離する設計 |
| SQL コード生成 | sqlc（[sqlc.yaml](sqlc.yaml)） |
| API コード生成 | oapi-codegen（[../docs/api/openapi.yaml](../docs/api/openapi.yaml) → `internal/apigen/`） |
| 設定読み込み | envconfig |
| ロギング | 標準 `log/slog` |
| 認証 | JWT |

ディレクトリ構成は [README.md](README.md) を参照。

---

## マイクロサービス アーキテクチャルール

### ルール1：サービス境界を越えた DB 直接アクセス禁止

```
❌ ticket-service が customers テーブルに直接 JOIN する
✅ ticket-service が customer-service API を呼び出す（Phase 2 以降）
✅ MVP(Phase 1) では同一プロセス内のモジュール関数を呼び出す
```

サービスが管理するテーブル：

| サービス | 所有テーブル |
|---------|------------|
| auth-service | users（読み取りのみ） |
| ticket-service | tickets, comments, ticket_history, attachments |
| customer-service | customers |
| system-service | systems, system_assignees |
| user-service | users |
| analytics-service | 全テーブル読み取り専用（Phase 1 限定） |

### ルール2：API Gateway に業務ロジックを持たせない

API Gateway が担うのは以下のみ：
- JWT 検証と `X-User-ID` / `X-User-Role` ヘッダー付与
- パスプレフィックスによるルーティング
- レートリミット / CORS / リクエストログ

ビジネスロジック（ステータス遷移ルール、権限チェックの詳細など）は必ず担当サービスに実装する。

### ルール3：API バージョニング必須

全エンドポイントは `/api/v1/` プレフィックスで始める。
後方互換性を破る変更は `/api/v2/` を新設し、`/api/v1/` を一定期間並行運用する。

```
✅ POST /api/v1/tickets
❌ POST /tickets
❌ POST /api/tickets
```

### ルール4：エラーレスポンス統一形式

全サービスで以下の形式を厳守する：

```json
{
  "code": "VALIDATION_FAILED",
  "message": "入力内容に誤りがあります",
  "details": {
    "customer_id": "必須項目です",
    "name": "100文字以内で入力してください"
  }
}
```

- `code`：SCREAMING_SNAKE_CASE。サービス名プレフィックス推奨（例: `TICKET_NOT_FOUND`, `SYSTEM_NAME_CONFLICT`）
- `message`：日本語のユーザー向けメッセージ
- `details`：フィールドごとのバリデーションエラー詳細（省略可）。**`{ フィールド名: 日本語メッセージ }` のフラット辞書**で返す（[../docs/api/openapi.yaml](../docs/api/openapi.yaml) の `Error` スキーマ参照）

エラー型は `internal/httperr/` に集約し、各ハンドラから利用する。

### ルール5：ステータス遷移は ticket-service が強制する

有効な遷移：

```
new → in_progress
in_progress → waiting | done
waiting → in_progress | done
done → （なし。終端ステータス）
```

これ以外の遷移は HTTP 422 + `code: INVALID_STATUS_TRANSITION` を返す。
フロントエンド側の表示制御とは独立して、**バックエンドで必ず検証する**。

### ルール6：権限チェックの実装方針

- `admin` ロールが必要なエンドポイントは、サービス内で `X-User-Role` ヘッダーを検証する
- JWT は API Gateway で検証済みのため、内部サービスは署名検証不要
- 権限不足は HTTP 403 を返す

管理者専用エンドポイント：
- `POST /customers`, `PUT /customers/{id}`
- `POST /systems`, `PUT /systems/{id}`
- `PUT /systems/{id}/assignees`
- `GET /users`, `GET /users/{id}`

### ルール7：論理削除を原則とする

顧客企業・システム・ユーザーは物理削除せず `is_active = false` で論理削除する。
これはサービス分離後の FK 制約廃止に備えた設計である。

### ルール8：`customer_id` はシステム登録後変更不可

`systems.customer_id` は登録後に変更できない（画面仕様 SCR-11 より）。
変更を試みるリクエストは HTTP 422 を返し、DB 更新を行わない。

---

## API 設計規約

### エンドポイント命名

- リソース名は複数形の名詞：`/tickets`, `/customers`, `/systems`
- アクション（動詞）は HTTP メソッドで表現する
- サブリソースはネスト：`/tickets/{id}/comments`
- 特定のフィールド更新は PATCH + サブパス：`PATCH /tickets/{id}/status`

### フィルタ・ページネーション

ページネーションの適用対象はエンドポイントの種類によって異なる。

#### 検索系一覧 API（必須）

独立したリソースの一覧取得エンドポイント。`page`/`per_page` クエリパラメータとレスポンスの `pagination` オブジェクトが必須。

対象：`GET /tickets`, `GET /customers`, `GET /systems`, `GET /users`

| パラメータ | 型 | 説明 |
|-----------|-----|------|
| page | integer | ページ番号（1始まり） |
| per_page | integer | 件数（デフォルト50、最大200） |

レスポンスに `pagination` オブジェクトを含める：

```json
{
  "data": [...],
  "pagination": {
    "total": 150,
    "page": 1,
    "per_page": 50,
    "total_pages": 3
  }
}
```

#### 1件詳細配下の子一覧 API（適用外）

特定リソース配下の子リソース一覧（例：`GET /tickets/{id}/comments`、`GET /tickets/{id}/attachments`）はページネーション不要。
ただし、並び順は必ず明記すること（例：コメントは `created_at` 昇順）。

### リレーション取得

`include` クエリパラメータでリレーションを追加取得できる：

```
GET /customers?include=systems
```

N+1 クエリを防ぐため、`include` で指定されたリレーションは JOIN または IN 句でまとめて取得すること。

### 日時形式

全ての日時フィールドは ISO 8601 形式（`2026-05-10T09:00:00Z`）を使用する。

---

## Go コード規約

### ディレクトリ構成

実構成は [README.md](README.md) を参照。骨格：

```
backend/
├── cmd/                     ← 各サービスの main エントリポイント
│   ├── gateway/
│   ├── auth/
│   └── ticket/
├── services/                ← サービス別のハンドラ・ドメインロジック
├── gateway/                 ← リバプロ + JWT ミドルウェア
└── internal/
    ├── apigen/              ← oapi-codegen 生成（コミット対象）
    ├── db/                  ← migrations / queries / sqlc 生成
    ├── domain/              ← 共有ドメインモデル（ENUM など）
    ├── httperr/             ← エラー形式統一（ルール4）
    ├── middleware/          ← request ID / slog / recover
    ├── auth/jwt/            ← JWT 発行・検証
    └── config/              ← envconfig
```

### サービス間の依存方向

```
gateway → （各サービスの HTTP endpoint）
ticket-service → system-service（担当者プール検証）
analytics-service → （全サービス DB の読み取り / Phase 1 限定）
```

循環依存は禁止。新しい依存が必要になった場合は [../docs/api/architecture.md](../docs/api/architecture.md) を更新してからコードを書く。

### コード生成（sqlc / oapi-codegen）

- **OpenAPI → Go 型**：`docs/api/openapi.yaml` を `internal/apigen/` に oapi-codegen で生成する。**手書きで `internal/apigen/` を編集しない**
- **SQL → Go 型**：[sqlc.yaml](sqlc.yaml) の設定で sqlc を実行。クエリは `internal/db/queries/` に、生成物は `internal/db/sqlc/` に出力する
- どちらも生成物はコミット対象。CI でも生成し、diff がないことを確認する

### エラーハンドリング

- ハンドラから返すエラーは `internal/httperr/` の型に変換してからレスポンスする
- 業務エラー（バリデーション失敗・遷移不可など）は `httperr.New(code, message, status)` で生成
- 想定外エラー（DB 接続失敗等）は `slog.Error` でログ出力し、`500 INTERNAL_ERROR` を返す
- エラーレスポンスの形式はルール4

### ロギング

- 標準 `log/slog` を使用
- ミドルウェアで request ID を context に注入し、全ログに含める
- ログレベル：DEBUG（ローカル）/ INFO（本番デフォルト）/ WARN（リトライ可能なエラー・想定内だが警戒すべき状態）/ ERROR（5xx 時）

### コンテキスト

- HTTP ハンドラから下層へ `context.Context` を必ず引き回す
- DB クエリ・他サービス呼び出しは context 付き API を使う
- タイムアウトはミドルウェア層で設定する（デフォルト 30 秒）。長時間処理が必要なエンドポイント（分析系等）は個別に上書きする

---

## 実装上の注意事項

1. **添付ファイルの S3 キー**は `tickets/{ticketId}/{uuid}.{ext}` 形式を使用する
2. **ダウンロード URL** は S3 署名付き URL を生成し、有効期限は 15 分とする
3. **ticket_history の記録**はステータス変更・担当者変更・チケット更新時に自動で行う。API クライアントからは明示的に呼び出さない
4. **サイドバーデータ**（`GET /analytics/sidebar`）は頻繁にポーリングされる可能性があるため、**プロセス内 TTL キャッシュ**（30 秒、`sync.Map` + 期限切れエントリ削除、`golang.org/x/sync/singleflight` でサンダリングハード回避）を実装する。ElastiCache / Redis は使わない（コスト最適化方針、[../infra/CLAUDE.md](../infra/CLAUDE.md) 参照）
5. **JWT 失効管理**：ステートレス JWT で運用する（ブラックリストストアを持たない）。アクセストークン TTL は 15 分（`JWT_TTL=15m`）。`POST /auth/logout` はクライアント側でトークン破棄するだけの no-op として実装（200 を返す）。即時失効が必要になった場合の refresh token / DB ベース失効管理は後続マイルストーンで検討
6. **ローカル開発の外部依存**：[../docker-compose.yml](../docker-compose.yml) で MySQL 8.0 / LocalStack（S3 互換）が起動する。AWS SDK は環境変数（`AWS_ENDPOINT_URL_S3` 等）で endpoint を上書きし、LocalStack ↔ 実 AWS を切り替えられる実装にする

---

## ドキュメント更新義務

新しいエンドポイントを追加・変更する場合は **必ず** [../docs/api/openapi.yaml](../docs/api/openapi.yaml) を同時に更新する。

チェックリスト：
- [ ] `paths` にエンドポイントを追加
- [ ] 対応する `tags` が存在するか確認（なければ追加）。タグはサービス単位
- [ ] 使用するスキーマを `components/schemas` に定義
- [ ] エラーレスポンスは `components/schemas/Error` を参照（ルール4）

新しいサービスを追加する場合は [../CLAUDE.md](../CLAUDE.md) の「横断ルール：新サービス追加手順」に従う。

---

## テスト

詳細は [../docs/testing/api-unit-test.md](../docs/testing/api-unit-test.md) を参照。

主要原則：
- ステータス遷移ルール（ルール5）の全パターンをユニットテストで網羅
- 権限チェック（admin/member）を各エンドポイントで検証
- 他サービス依存はインターフェース化してモック可能にする
- カバレッジの数値目標は設けず、必須テストパターン網羅を優先

---

## CI

backend のビルド・テストは [`.github/workflows/backend.yml`](../.github/workflows/backend.yml) で実行される（`go build` / `go vet` / `go test`）。詳細は [../.github/CLAUDE.md](../.github/CLAUDE.md) を参照。

---

## マイルストーン

実装は段階的に進める。現在のマイルストーンと進捗は [README.md](README.md) を参照。
