# マイクロサービスアーキテクチャ設計

**文書バージョン：** 1.0
**作成日：** 2026-05-10

---

## 1. 設計方針

Support Ops Hub は**ドメイン駆動設計（DDD）の境界づけられたコンテキスト**を基にサービスを分割する。
MVP はモノリスとして実装しながらも、**内部をサービスモジュール単位で分離**することで、将来のマイクロサービス化コストを最小化する。

> **段階的移行戦略**
> Phase 1（MVP）: モジュラーモノリス（単一プロセス、DBは論理分離）
> Phase 2: analytics-service を独立プロセスに切り出し（読み取り専用なので副作用が少ない）
> Phase 3: ticket-service を切り出し（最も高負荷が予想されるため）
> Phase 4: 残りのサービスを順次独立化

---

## 2. サービスマップ

```
                         ┌─────────────────────────────┐
  Nuxt.js (Frontend)     │       API Gateway            │
  ─────────────────────► │  (ルーティング / 認証検証)    │
                         └──────┬────────────────┬──────┘
                                │ 内部HTTP通信    │
          ┌─────────────────────┼────────────────┼─────────────────────┐
          │                     │                │                     │
   ┌──────▼──────┐  ┌───────────▼───────┐  ┌────▼────────┐  ┌────────▼───────┐
   │ auth-service │  │  ticket-service    │  │  analytics- │  │ customer/system│
   │    :8001     │  │      :8002         │  │  service    │  │  /user services│
   │             │  │                   │  │    :8006    │  │  :8003/:8004   │
   │ ・ログイン  │  │ ・チケットCRUD    │  │             │  │   /:8005       │
   │ ・JWT発行  │  │ ・ステータス遷移  │  │ ・集計クエリ│  │ ・マスタ管理   │
   │ ・no-op    │  │ ・コメント        │  │ ・レポート  │  │ ・ユーザー管理 │
   │  logout    │  │ ・変更履歴        │  │             │  │               │
   │            │  │ ・添付ファイル    │  │             │  │               │
   └──────┬──────┘  └───────────┬───────┘  └──────┬──────┘  └───────────────┘
          │                     │                 │（読み取り専用クエリ）
          │         ┌───────────▼────────────────▼───────────────────────────┐
          └────────►│                MySQL (RDS for MySQL 8.0)                 │
                    │  ┌─────────────┐  ┌─────────────┐  ┌─────────────────┐ │
                    │  │  tickets DB  │  │  masters DB  │  │  analytics DB   │ │
                    │  │  tickets     │  │  customers   │  │  (将来: 読み取   │ │
                    │  │  comments    │  │  systems     │  │  り専用レプリカ) │ │
                    │  │  history     │  │  users       │  └─────────────────┘ │
                    │  │  attachments │  │  sys_assign  │                      │
                    │  └─────────────┘  └─────────────┘                       │
                    └────────────────────────────────────────────────────────── ┘

  ※ JWT はステートレス運用（TTL 15 分）。Redis / ElastiCache は使わない。
    サイドバーキャッシュは各サービスプロセス内 TTL キャッシュで実現。
  ※ 上図の "API Gateway" は **アプリ独自の Go バイナリ**（`backend/cmd/gateway`）であり、
    AWS マネージドの API Gateway サービスではない。AWS 上は ALB → ECS on EC2 task
    として配置される（詳細は §4 と §7.1）。
```

---

## 3. サービス定義

### 3.1 auth-service

| 項目 | 内容 |
|------|------|
| 責務 | 認証・JWT 発行 |
| 内部ポート | 8001 |
| DB | users テーブルを参照（user-service DB） |
| 外部依存 | なし |

**JWT 失効戦略：** ステートレス JWT。アクセストークン TTL は 15 分。即時失効ストア（Redis 等）は持たない。`POST /auth/logout` はクライアントがトークンを破棄するだけの no-op として動作する（200 を返す）。最大 15 分のラグを許容する設計とし、refresh token / DB ベース失効管理は後続マイルストーンで検討。

**段階的 Auth 戦略（マイルストーン別）：**

サービス境界は最初から `X-User-ID` / `X-User-Role` ヘッダーで設計し、Gateway の auth 実装だけを段階的に差し替える。サービスコードはヘッダーの出元（dev-auth か JWT か）を意識しない。

| フェーズ | Gateway の auth ミドルウェア | クライアント（FE） |
|---|---|---|
| **M2〜M6** | **dev-auth**：`X-Dev-User-Email` ヘッダーまたは dev cookie から user を引いて DB で role を取得、`X-User-ID` / `X-User-Role` を内部リクエストに注入。`POST /auth/login` はスタブ（email だけでダミートークン返却） | 開発用ユーザー切替セレクタで `X-Dev-User-Email` を localStorage または cookie に保存 |
| **M7** | **本物 JWT**：`POST /auth/login` で bcrypt 検証 + JWT 発行、ミドルウェアは Bearer トークンを検証してヘッダー注入 | **SCR-01 ログイン画面**を追加、JWT を HttpOnly cookie で保存、API Gateway 経由でリクエスト |

この設計により、M2〜M6 で書くサービスコードは M7 で 1 行も変更されない。

**エンドポイント：**

| メソッド | パス | 説明 |
|---------|------|------|
| POST | /auth/login | M2〜M6: dev スタブ / M7: 本物のログイン・JWT 発行（TTL 15 分） |
| POST | /auth/logout | no-op（クライアント側でトークン破棄） |
| GET | /auth/me | 現在ユーザー情報 |

**ステータス遷移ルールはこのサービスが管理しない。** チケットのステータス遷移は ticket-service が担う。

---

### 3.2 ticket-service

| 項目 | 内容 |
|------|------|
| 責務 | チケットのライフサイクル管理（CRUD・ステータス遷移・履歴・コメント） |
| 内部ポート | 8002 |
| DB | tickets, comments, ticket_history, attachments |
| 外部依存 | AWS S3（添付ファイル保存）、customer-service / system-service への参照整合性チェック |

**ステータス遷移ルール（このサービスが強制する）：**

```
new → in_progress
in_progress → waiting | done
waiting → in_progress | done
done → （なし。終端ステータス）
```

**他サービスへの依存注意点：**
チケット登録時に `system_id` の担当者プール検証が必要。MVP では同一DB内参照で解決するが、サービス分離後は system-service API を呼び出す。

---

### 3.3 customer-service

| 項目 | 内容 |
|------|------|
| 責務 | 顧客企業マスタのCRUD |
| 内部ポート | 8003 |
| DB | customers テーブル |
| 外部依存 | なし |

管理者ロールの検証は API Gateway が JWT から行い、サービス内では `X-User-Role` ヘッダーを信頼する。

---

### 3.4 system-service

| 項目 | 内容 |
|------|------|
| 責務 | システムマスタ・担当者プール管理 |
| 内部ポート | 8004 |
| DB | systems, system_assignees テーブル |
| 外部依存 | customer-service（顧客存在確認）、user-service（ユーザー存在確認） |

`customer_id` は登録後変更不可。変更を試みるリクエストは 422 で拒否する。

---

### 3.5 user-service

| 項目 | 内容 |
|------|------|
| 責務 | ユーザーマスタ管理（MVP は閲覧のみ） |
| 内部ポート | 8005 |
| DB | users テーブル |
| 外部依存 | なし |

MVP では管理者がDBに直接ユーザーを登録する運用を想定。将来的にユーザー招待・パスワードリセット機能を追加する。

---

### 3.6 analytics-service

| 項目 | 内容 |
|------|------|
| 責務 | ダッシュボード・レポート集計（読み取り専用） |
| 内部ポート | 8006 |
| DB | 全テーブルの読み取り専用アクセス（Phase 1）→ 将来: 読み取り専用レプリカ |
| 外部依存 | 将来: ticket-service のドメインイベントをサブスクライブ |

**Phase 1 の妥協点：** analytics-service は直接 DB を読み取る。これはサービス分離の原則に反するが、MVP の集計クエリは複雑であり、サービス間通信でデータを取得するコストが大きいため段階的に解決する。

---

## 4. API Gateway の責務

> **用語の整理：** 本書の "API Gateway" は **AWS マネージドサービスの API Gateway ではない**。アプリ独自の Go バイナリ（`backend/cmd/gateway/main.go`）であり、AWS では **ECS on EC2 のタスク** として起動する。AWS 側のルーティング層は **ALB**（詳細は §7.1）。コスト最適化方針で AWS API Gateway は不採用とした。
>
> ```
> Internet → CloudFront (prod) / 直接 → ALB → ECS on EC2 task (cmd/gateway) → 同一プロセス内の各サービスモジュール
> ```

API Gateway（= `cmd/gateway`）が担う横断的関心事：

| 責務 | 実装 |
|------|------|
| JWT検証 | Authorization ヘッダーを検証し、内部サービスへ `X-User-ID`, `X-User-Role` ヘッダーを付与（M2〜M6 は dev-auth、M7 で本物 JWT へ差替、§3.1 参照） |
| ルーティング | パスプレフィックスで内部サービスへルーティング |
| レートリミット | IPアドレスおよびユーザーIDベース |
| リクエストログ | 全リクエストを記録（セキュリティ監査用） |
| CORS | フロントエンドオリジンのみ許可 |
| エラー正規化 | 内部サービスのエラーを統一形式に変換 |

**API Gatewayは業務ロジックを持たない。** ルーティングと横断処理のみ担当する。

---

## 5. サービス間通信

### Phase 1（MVP）：同一プロセス内モジュール呼び出し

```
[ticket-module] → [system-module] : Go パッケージ内関数呼び出し
```

### Phase 2 以降：HTTP / イベント駆動

| 通信種別 | 用途 | 実装 |
|---------|------|------|
| 同期HTTP | データ検証・即時参照 | gRPC or REST |
| 非同期イベント | analytics-service への変更通知 | Kafka or SQS |

**サービス分離後に発生するチケット登録フロー例：**

```
1. POST /tickets (API Gateway)
2. → ticket-service: チケット作成
3.   → system-service API: 担当者がプールに含まれるか検証
4.   → DB INSERT (tickets)
5.   → SQS/Kafka: TICKET_CREATED イベント発行
6. analytics-service: TICKET_CREATED イベントを受信してサマリを更新
```

---

## 6. データ設計方針

### DB分離戦略

| Phase | 分離レベル |
|-------|-----------|
| Phase 1（MVP） | 論理分離（同一MySQL、スキーマを用途別に意識） |
| Phase 2 | analytics-service のみ読み取り専用レプリカを参照 |
| Phase 3 | ticket-service を別MySQL インスタンスに移行 |
| Phase 4 | 各サービス独立DBに完全分離 |

### 他サービスIDの参照方法

サービス分離後、FK制約は**使用しない**。代わりに：
- アプリケーションレイヤーで参照整合性を検証
- 削除は論理削除（`is_active = false`）を原則とする

---

## 7. セキュリティ設計

### 認証フロー

```
Client → POST /auth/login
       ← JWT (payload: {sub: user_id, role, exp})

Client → GET /tickets (Authorization: Bearer <JWT>)
API Gateway: JWTを検証 → X-User-ID, X-User-Role ヘッダーを付与
       → ticket-service (内部。外部から直接アクセス不可)
```

### 権限制御

| エンドポイント | 必要ロール |
|--------------|-----------|
| POST/PUT /customers | admin |
| POST/PUT /systems | admin |
| GET /users | admin |
| その他 | admin または member |

権限チェックはサービス内で `X-User-Role` ヘッダーを参照して行う。内部サービスへの直接アクセスは Security Group で ALB からの受信のみ許可することで防ぐ。

---

## 7.1 AWS デプロイ構成（コスト最適化リビジョン）

| レイヤー | 採用 | コスト・運用上の理由 |
|---|---|---|
| ルーティング層 | **ALB**（AWS マネージド API Gateway は不採用） | アプリ独自の `cmd/gateway`（§4）を ECS タスクとして起動し、AWS マネージド API Gateway のリクエスト課金を回避 |
| コンテナ実行 | **ECS on EC2**（capacity provider + Auto Scaling Group） | Fargate より単価が低く、複数タスク相乗りで効率化 |
| EC2 サイズ | dev/stg: t4g.small × 1、prod: t4g.medium × 2 (Multi-AZ) | Graviton（ARM）で価格性能比を高める |
| RDB | RDS for MySQL 8.0、db.t4g.micro〜small、gp3 20〜50GB | Single-AZ（dev/stg）/ Multi-AZ（prod のみ） |
| KVS | **使わない**（ステートレス JWT + プロセス内 TTL キャッシュ） | ElastiCache を削除して月額 $13〜 を節約 |
| ネットワーク | **NAT なし**。アプリ EC2 は public subnet（SG で ALB からのみ受信）、RDS は private subnet | NAT Gateway を削除して月額 $32×AZ を節約 |
| egress 抑制 | S3 への通信は Gateway 型 VPC Endpoint（無料） | データ転送料の抑制 |
| 静的配信 | CloudFront + S3（prod のみ）。dev/stg は ALB 直結 | dev の固定費削減 |

```
                             ┌────────────────┐
  Client (Nuxt SSR / SPA) ──►│   ALB (https)  │
                             └───────┬────────┘
                                     │ パスベースルーティング
                             ┌───────▼─────────────────────────────────────┐
                             │ ECS Cluster (EC2 launch type)              │
                             │  ┌─────────────────────────────────────┐  │
                             │  │ task: cmd/gateway (本書の API Gateway)│  │
                             │  │   ├─ ・JWT/dev-auth 検証              │  │
                             │  │   ├─ ・ヘッダー注入 (X-User-ID etc.)  │  │
                             │  │   ├─ ・CORS / レートリミット          │  │
                             │  │   └─ 同一プロセス内のモジュール関数で │  │
                             │  │      各 service を呼ぶ                │  │
                             │  └─────────────────────────────────────┘  │
                             └────────────────┬────────────────────────────┘
                                              │ private subnet
                                       ┌──────▼──────┐
                                       │ RDS MySQL 8 │
                                       └─────────────┘
```

詳細は [../../infra/CLAUDE.md](../../infra/CLAUDE.md) および [../../infra/terraform/README.md](../../infra/terraform/README.md) を参照。

---

## 8. 将来拡張ポイント

| 機能 | 担当サービス | 実装方針 |
|------|------------|---------|
| Slack通知 | notification-service（新規） | ticket-service のイベントをサブスクライブ |
| Slack Bot自動起票 | webhook-service（新規） | Slack Event API → ticket-service |
| 障害管理 | incident-service（新規） | ticket-service と同一ドメインイベントを共有 |
| 顧客ポータル | customer-portal-service（新規） | ticket-service のサブセットAPI |
| SLAアラート | sla-service（新規） | ticket-service のイベント + スケジューラ |
