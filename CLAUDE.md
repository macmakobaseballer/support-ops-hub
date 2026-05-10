# Support Ops Hub — CLAUDE.md

AIエージェントがこのプロジェクトで作業する際のアーキテクチャルールと規約。
**これらのルールは提案ではなくハードコンストレイントである。** 逸脱する前に必ずユーザーに確認すること。

---

## プロジェクト概要

問い合わせ管理システム（Support Ops Hub）。
運用・保守チームが顧客からの問い合わせを一元管理するWebアプリケーション。

| 項目 | 内容 |
|------|------|
| フロントエンド | Nuxt.js |
| バックエンド | Go |
| データベース | MySQL 8.0（AWS RDS） |
| インフラ | AWS（ECS Fargate / S3 / CloudFront） |
| API仕様 | [docs/api/openapi.yaml](docs/api/openapi.yaml) |
| アーキテクチャ設計 | [docs/api/architecture.md](docs/api/architecture.md) |

---

## マイクロサービス アーキテクチャルール

### ルール1：サービス境界を越えたDB直接アクセス禁止

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

### ルール2：API Gatewayに業務ロジックを持たせない

API Gateway が担うのは以下のみ：
- JWT検証と `X-User-ID` / `X-User-Role` ヘッダー付与
- パスプレフィックスによるルーティング
- レートリミット / CORS / リクエストログ

ビジネスロジック（ステータス遷移ルール、権限チェックの詳細など）は必ず担当サービスに実装する。

### ルール3：APIバージョニング必須

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
  "code": "TICKET_NOT_FOUND",
  "message": "指定されたチケットは存在しません",
  "details": {}
}
```

- `code`：SCREAMING_SNAKE_CASE。サービス名プレフィックスを付けることを推奨（例: `TICKET_NOT_FOUND`, `SYSTEM_NAME_CONFLICT`）
- `message`：日本語のユーザー向けメッセージ
- `details`：フィールドバリデーションエラーの詳細（省略可）

### ルール5：ステータス遷移はticket-serviceが強制する

有効な遷移：

```
new → in_progress
in_progress → waiting | done
waiting → in_progress | done
done → （なし。終端ステータス）
```

これ以外の遷移は HTTP 422 + `code: INVALID_STATUS_TRANSITION` を返す。
フロントエンド側の表示制御（ボタンの出し分け）とは独立して、**バックエンドで必ず検証する**。

### ルール6：権限チェックの実装方針

- `admin` ロールが必要なエンドポイントは、サービス内で `X-User-Role` ヘッダーを検証する
- JWT はAPI Gatewayで検証済みのため、内部サービスは署名検証不要
- 権限不足は HTTP 403 を返す

管理者専用エンドポイント：
- `POST /customers`, `PUT /customers/{id}`
- `POST /systems`, `PUT /systems/{id}`
- `PUT /systems/{id}/assignees`
- `GET /users`, `GET /users/{id}`

### ルール7：論理削除を原則とする

顧客企業・システム・ユーザーは物理削除せず `is_active = false` で論理削除する。
これはサービス分離後のFK制約廃止に備えた設計である。

### ルール8：`customer_id` はシステム登録後変更不可

`systems.customer_id` は登録後に変更できない（画面仕様 SCR-11 より）。
変更を試みるリクエストは HTTP 422 を返し、DB更新を行わない。

---

## API 設計規約

### エンドポイント命名

- リソース名は複数形の名詞：`/tickets`, `/customers`, `/systems`
- アクション（動詞）は HTTP メソッドで表現する
- サブリソースはネストで表現：`/tickets/{id}/comments`
- 特定のフィールド更新は PATCH + サブパス：`PATCH /tickets/{id}/status`

### フィルタ・ページネーション

一覧エンドポイントの共通クエリパラメータ：

| パラメータ | 型 | 説明 |
|-----------|-----|------|
| page | integer | ページ番号（1始まり） |
| per_page | integer | 件数（デフォルト50、最大200） |

レスポンスには必ず `pagination` オブジェクトを含める：

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

### リレーション取得

`include` クエリパラメータでリレーションを追加取得できる：

```
GET /customers?include=systems
```

N+1クエリを防ぐため、`include` で指定されたリレーションは JOIN または IN句でまとめて取得すること。

### 日時形式

全ての日時フィールドは ISO 8601 形式（`2026-05-10T09:00:00Z`）を使用する。

---

## ドキュメント規約

### OpenAPI仕様の更新

新しいエンドポイントを追加・変更する場合は必ず [docs/api/openapi.yaml](docs/api/openapi.yaml) を同時に更新すること。

チェックリスト：
- [ ] `paths` にエンドポイントを追加
- [ ] 対応する `tags` が存在するか確認（なければ追加）
- [ ] 使用するスキーマを `components/schemas` に定義
- [ ] `x-service` 拡張フィールドでどのサービスが担当するか明記

### 新しいサービスの追加

1. [docs/api/architecture.md](docs/api/architecture.md) のサービスマップを更新
2. [docs/api/openapi.yaml](docs/api/openapi.yaml) に新しいタグを追加
3. [CLAUDE.md](CLAUDE.md) の「所有テーブル」テーブルを更新

---

## コード規約（Go バックエンド）

### ディレクトリ構成（想定）

```
services/
  auth/           ← auth-service
  ticket/         ← ticket-service
  customer/       ← customer-service
  system/         ← system-service
  user/           ← user-service
  analytics/      ← analytics-service
gateway/          ← API Gateway
internal/
  domain/         ← 共有ドメインモデル（ENUMなど）
  middleware/     ← 認証ミドルウェア
```

### サービス間の依存方向

```
gateway → （各サービスのHTTP endpoint）
ticket-service → system-service（担当者プール検証）
analytics-service → （全サービスDBの読み取り / Phase 1）
```

循環依存は禁止。新しい依存が必要になった場合はアーキテクチャ設計書を更新してからコードを書く。

---

## テスト規約

- ステータス遷移ルールの全パターンはユニットテストで網羅する
- 権限チェック（admin/member）は各エンドポイントで統合テストを書く
- 他サービスへの依存がある箇所はインターフェースでモック可能にする

---

## 実装上の注意事項

1. **添付ファイルのS3キー**は `tickets/{ticketId}/{uuid}.{ext}` 形式を使用する
2. **ダウンロードURL**はS3署名付きURLを生成し、有効期限は15分とする
3. **ticket_history の記録**はステータス変更・担当者変更・チケット更新時に自動で行う。APIクライアントからは明示的に呼び出さない
4. **サイドバーデータ**（`GET /analytics/sidebar`）は頻繁にポーリングされる可能性があるため、Redis キャッシュ（TTL: 30秒）を検討すること
