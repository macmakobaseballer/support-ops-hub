---
name: api-design
description: API設計詳細規約（エンドポイント命名・ページネーション・フィルタ・リレーション・日時形式）。新規エンドポイント実装時に参照。
---

# API設計詳細規約

## エンドポイント命名

- リソース名は複数形の名詞：`/tickets`, `/customers`, `/systems`
- アクション（動詞）は HTTP メソッドで表現する
- サブリソースはネスト：`/tickets/{id}/comments`
- 特定のフィールド更新は PATCH + サブパス：`PATCH /tickets/{id}/status`

## フィルタ・ページネーション

### 検索系一覧 API（必須）

独立したリソースの一覧取得エンドポイントには `page`/`per_page` クエリパラメータとレスポンスの `pagination` オブジェクトが必須。

対象：`GET /tickets`, `GET /customers`, `GET /systems`, `GET /users`

| パラメータ | 型 | 説明 |
|-----------|-----|------|
| page | integer | ページ番号（1始まり） |
| per_page | integer | 件数（デフォルト50、最大200） |

レスポンス形式：

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

### 1件詳細配下の子一覧 API（適用外）

`GET /tickets/{id}/comments`、`GET /tickets/{id}/attachments` 等はページネーション不要。ただし並び順は必ず明記（例：コメントは `created_at` 昇順）。

## リレーション取得

`include` クエリパラメータでリレーションを追加取得できる：

```
GET /customers?include=systems
```

N+1 クエリを防ぐため、`include` で指定されたリレーションは JOIN または IN 句でまとめて取得すること。

## 日時形式

全ての日時フィールドは ISO 8601 形式（`2026-05-10T09:00:00Z`）を使用する。
