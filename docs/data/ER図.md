# ER図

**文書バージョン：** 1.0
**作成日：** 2026-05-10
**ステータス：** 確定

## 改訂履歴

| バージョン | 日付 | 変更内容 |
|-----------|------|---------|
| 1.0 | 2026-05-10 | 初版作成 |

---

## ER図（Mermaid）

```mermaid
erDiagram
    customers {
        BIGINT id PK
        VARCHAR name UK
        TEXT notes
        TINYINT is_active
        DATETIME created_at
        DATETIME updated_at
    }

    systems {
        BIGINT id PK
        VARCHAR name
        BIGINT customer_id FK
        TEXT description
        TINYINT is_active
        DATETIME created_at
        DATETIME updated_at
    }

    users {
        BIGINT id PK
        VARCHAR name
        VARCHAR email UK
        VARCHAR password_hash
        ENUM role
        TINYINT is_active
        DATETIME created_at
        DATETIME updated_at
    }

    system_assignees {
        BIGINT system_id FK
        BIGINT user_id FK
    }

    tickets {
        BIGINT id PK
        VARCHAR title
        TEXT description
        ENUM type
        ENUM priority
        ENUM status
        BIGINT customer_id FK
        BIGINT system_id FK
        BIGINT assignee_id FK
        BIGINT created_by FK
        DATETIME received_at
        DATETIME created_at
        DATETIME updated_at
    }

    comments {
        BIGINT id PK
        BIGINT ticket_id FK
        BIGINT author_id FK
        TEXT body
        DATETIME created_at
    }

    ticket_history {
        BIGINT id PK
        BIGINT ticket_id FK
        BIGINT changed_by FK
        VARCHAR field_name
        TEXT old_value
        TEXT new_value
        DATETIME changed_at
    }

    attachments {
        BIGINT id PK
        BIGINT ticket_id FK
        BIGINT uploaded_by FK
        VARCHAR file_name
        VARCHAR file_key
        INT file_size
        VARCHAR content_type
        DATETIME created_at
    }

    customers ||--o{ systems : "has"
    customers ||--o{ tickets : "belongs to"
    systems ||--o{ tickets : "belongs to"
    systems }o--o{ users : "system_assignees"
    users ||--o{ tickets : "assignee (nullable)"
    users ||--o{ tickets : "created_by"
    tickets ||--o{ comments : "has"
    tickets ||--o{ ticket_history : "has"
    tickets ||--o{ attachments : "has"
    users ||--o{ comments : "author"
    users ||--o{ ticket_history : "changed_by"
    users ||--o{ attachments : "uploaded_by"
```

---

## テーブル一覧

| テーブル名 | 概要 |
|---|---|
| `customers` | 顧客企業マスター |
| `systems` | システムマスター（顧客企業に紐づく） |
| `users` | ユーザーマスター（担当者・管理者） |
| `system_assignees` | システム担当者プール（中間テーブル） |
| `tickets` | チケット（問い合わせ本体） |
| `comments` | コメント（対応履歴） |
| `ticket_history` | 変更履歴（監査ログ） |
| `attachments` | 添付ファイル（S3参照） |

---

## 主なリレーション補足

| リレーション | カーディナリティ | 備考 |
|---|---|---|
| customers → systems | 1:N | 1顧客が複数システムを持つ |
| customers → tickets | 1:N | チケットには顧客が必ず紐づく |
| systems → tickets | 1:N | チケットにはシステムが必ず紐づく |
| systems ↔ users | M:N | `system_assignees` 中間テーブルで管理 |
| users → tickets（assignee） | 1:N | 担当者は省略可（NULL許容） |
| users → tickets（created_by） | 1:N | チケット作成者（必須） |
| tickets → comments | 1:N | コメントは時系列順で表示 |
| tickets → ticket_history | 1:N | フィールド変更ごとに1レコード記録 |
| tickets → attachments | 1:N | 添付ファイルのS3パスを管理 |
