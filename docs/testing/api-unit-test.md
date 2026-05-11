# API ユニットテスト規約

Go バックエンドのユニットテストに関する規約を定義する。
[backend/CLAUDE.md](../../backend/CLAUDE.md) のアーキテクチャルールと同様、**これらは提案ではなくハードコンストレイントである。** 逸脱する前に必ずチームで合意すること。

---

## 1. 目的とスコープ

本ドキュメントが扱う「ユニットテスト」とは、**外部システム（DB・他サービス・S3・Redis）に依存せず、Go プロセス内で完結するテスト**を指す。

| 種別 | 本ドキュメントの対象 | 依存 |
|------|------------------|------|
| ユニットテスト | ✅ 対象 | なし（モック使用） |
| インテグレーションテスト | ❌ 対象外 | 実DB・テスト用コンテナ等 |
| E2E テスト | ❌ 対象外 | 起動中の全サービス |

インテグレーション・E2E の規約は別途定義する。

---

## 2. テストフレームワーク

以下のライブラリを標準採用する。**新しいライブラリの導入は事前にチームで合意すること。**

| 用途 | ライブラリ |
|------|-----------|
| テストランナー | 標準 `testing` |
| アサーション | `github.com/stretchr/testify/assert` / `require` |
| モック生成 | `go.uber.org/mock`（gomock の後継） |
| HTTP テスト | 標準 `net/http/httptest` |
| 時刻・ID 制御 | 独自インターフェース注入（後述） |

`go-sqlmock` は **使わない**。SQL を発行する箇所は repository 層に閉じ込め、上位層のテストではインターフェース経由でモックする。SQL の正しさは integration テストで検証する。

---

## 3. テストの分類とレイヤー

各サービスは以下のレイヤー構造を持つ。テストはレイヤーごとに書く。

```
handler → service (usecase) → repository
              ↓
           domain（純粋ロジック）
```

| レイヤー | テスト対象 | 依存の扱い |
|---------|----------|-----------|
| domain | ステータス遷移ルール、バリデーション、ID 生成等の純粋ロジック | 依存なし |
| service / usecase | 業務ロジック・他サービス呼び出し | repository・他 service をモック |
| handler | リクエスト/レスポンス変換、HTTP ステータス、エラー形式 | service をモック、`httptest` 使用 |
| repository | （ユニットテスト対象外） | integration で検証 |

**handler は service をモックして単体で検証する。** handler テストの中で実際の DB やビジネスロジックを動かしてはならない。

---

## 4. ファイル配置・命名規約

### 4.1 ファイル配置

- テストファイルは **テスト対象と同じディレクトリ**に置く（`xxx.go` の隣に `xxx_test.go`）
- パッケージは原則 **同パッケージ** (`package ticket`)
- 公開 API だけをテストする場合のみ **外部パッケージ** (`package ticket_test`) を使う

### 4.2 テスト関数命名

```
Test{対象}_{条件}_{期待結果}
```

例：
```go
func TestUpdateStatus_FromNewToInProgress_Succeeds(t *testing.T) { ... }
func TestUpdateStatus_FromDone_Returns422(t *testing.T) { ... }
func TestCreateTicket_AsMember_Returns403(t *testing.T) { ... }
```

- 単に `TestUpdateStatus` のような命名は禁止（テーブル駆動テスト全体を包む親関数に限り可）
- 日本語のテスト関数名は使わない（grep しづらい）。日本語が必要ならテーブル駆動の `name` フィールドに書く

---

## 5. テーブル駆動テストの規約

ステータス遷移・権限チェック等の網羅テストは **必ずテーブル駆動**で書く。

```go
func TestUpdateStatus(t *testing.T) {
    tests := []struct {
        name     string
        from     domain.Status
        to       domain.Status
        wantErr  bool
        wantCode string
    }{
        {"new から in_progress は成功", domain.StatusNew, domain.StatusInProgress, false, ""},
        {"done からの遷移は不可", domain.StatusDone, domain.StatusInProgress, true, "INVALID_STATUS_TRANSITION"},
        // ...
    }

    for _, tt := range tests {
        tt := tt
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel()
            err := domain.ValidateStatusTransition(tt.from, tt.to)
            if tt.wantErr {
                require.Error(t, err)
                assert.Equal(t, tt.wantCode, err.Code())
            } else {
                require.NoError(t, err)
            }
        })
    }
}
```

ルール：
- `name` フィールド必須。**日本語で「何を検証するか」を書く**
- `t.Run(tt.name, ...)` で必ずサブテスト化（失敗時に該当ケースが特定できる）
- ループ変数キャプチャ対策の `tt := tt` を入れる（Go 1.22 以降は不要だが互換のため記述）

---

## 6. 並列実行

- ユニットテストはすべて `t.Parallel()` を付ける
- 共有状態（パッケージグローバル変数・`os.Setenv` 等）に依存するテストは並列禁止 → そもそもプロダクトコード側で共有状態を持たない設計にする
- 時刻・乱数・UUID はインターフェース注入する（次節）

---

## 7. モック戦略

### 7.1 外部依存はインターフェース化

- 他サービス呼び出し・S3 SDK・Redis クライアント等は必ずインターフェースを定義し、本物の実装と並列にモックを用意する
- モックは `gomock` の `mockgen` で自動生成する（`//go:generate` で記述）

### 7.2 時刻と UUID は注入する

```go
type Clock interface { Now() time.Time }
type IDGenerator interface { NewID() string }
```

- `time.Now()` / `uuid.New()` をプロダクトコードから直接呼ばない
- service の struct にフィールドとして持たせ、テストでは固定値を返すモックを注入する

### 7.3 DB はモックしない

SQL モック（`go-sqlmock`）は使わず、**repository インターフェースをモックする**。

```go
// ❌ ダメな例：handler や service の中で sqlmock を使う
// ✅ 良い例：repository インターフェースをモック化
type TicketRepository interface {
    FindByID(ctx context.Context, id string) (*Ticket, error)
    Save(ctx context.Context, t *Ticket) error
}
```

実際の SQL の正しさは integration テストで検証する。

---

## 8. 必須テストパターン

[backend/CLAUDE.md](../../backend/CLAUDE.md) のアーキテクチャルールに対応する以下のテストは **必ず書く**。

### 8.1 ステータス遷移（backend/CLAUDE.md ルール5）

domain 層で以下の遷移を全パターン網羅する：

| from | to | 期待結果 |
|------|----|---------|
| new | in_progress | 成功 |
| new | waiting | 422 INVALID_STATUS_TRANSITION |
| new | done | 422 INVALID_STATUS_TRANSITION |
| in_progress | waiting | 成功 |
| in_progress | done | 成功 |
| in_progress | new | 422 |
| waiting | in_progress | 成功 |
| waiting | done | 成功 |
| waiting | new | 422 |
| done | * | 422（終端ステータス） |

### 8.2 権限チェック（backend/CLAUDE.md ルール6）

admin 専用エンドポイントは handler テストで以下を網羅：

- `X-User-Role: admin` → 200/201
- `X-User-Role: member` → 403
- `X-User-Role` ヘッダー欠落 → 403

対象エンドポイント（backend/CLAUDE.md ルール6 参照）：
- `POST /customers`, `PUT /customers/{id}`
- `POST /systems`, `PUT /systems/{id}`, `PUT /systems/{id}/assignees`
- `GET /users`, `GET /users/{id}`

### 8.3 エラーレスポンス形式（backend/CLAUDE.md ルール4）

エラーを返すハンドラのテストでは以下を検証：

- レスポンスボディに `code`・`message` フィールドが存在
- `code` が SCREAMING_SNAKE_CASE
- HTTP ステータスコードが意図通り

ヘルパー関数 `assertErrorResponse(t, resp, 422, "INVALID_STATUS_TRANSITION")` を共通化する。

### 8.4 `customer_id` 変更不可（backend/CLAUDE.md ルール8）

`PUT /systems/{id}` で `customer_id` を変更しようとするリクエストが 422 を返し、repository の更新メソッドが呼び出されない（モックの呼び出し回数 0）ことを検証する。

### 8.5 論理削除（backend/CLAUDE.md ルール7）

削除系エンドポイントが repository の「`is_active = false` 更新」メソッドを呼び出し、物理削除メソッドを呼ばないことをモックで検証する。

### 8.6 ページネーション（backend/CLAUDE.md「フィルタ・ページネーション」）

検索系一覧 API (`GET /tickets`・`/customers`・`/systems`・`/users`) のハンドラテストでは：

- `per_page` 未指定時にデフォルト値（50）が使われる
- `per_page` 上限（200）超過時に上限に丸められる、または 422 を返す（実装方針に従う）
- レスポンスに `pagination` オブジェクトが含まれる

---

## 9. アサーション規約

- **`require`**：失敗したら以降のアサーションが意味をなさないとき（`require.NoError`・`require.NotNil` 等）
- **`assert`**：失敗しても他のアサーションを継続したいとき
- 1 つのテスト関数（または 1 サブテスト）で **assertion を 10 個以上書かない**。検証観点が多すぎる場合はテーブル駆動に分解する
- `assert.Equal` の引数順は `(t, expected, actual)`。逆に書かない
- `assert.True(t, len(xs) == 3)` は禁止。`assert.Len(t, xs, 3)` を使う

---

## 10. テストヘルパー

### 10.1 builder パターン

ドメインオブジェクトのテスト用インスタンス生成は **builder（functional options）パターン**で共通化する：

```go
// testutil/ticket_builder.go
func NewTicket(opts ...func(*domain.Ticket)) *domain.Ticket {
    t := &domain.Ticket{
        ID:     "tkt_default",
        Status: domain.StatusNew,
        // ...
    }
    for _, opt := range opts {
        opt(t)
    }
    return t
}

func WithStatus(s domain.Status) func(*domain.Ticket) {
    return func(t *domain.Ticket) { t.Status = s }
}
```

### 10.2 `testdata/` の使い分け

- JSON リクエスト/レスポンスのサンプル等、ファイルで管理した方が読みやすいものは `testdata/` に置く
- 1 行で書ける程度の値は inline で書く（過剰な fixture 化を避ける）

---

## 11. カバレッジ

**数値目標は設けない。** カバレッジ率を満たすために無意味なテストを書くことを避けるため、以下の方針とする：

- 必須テストパターン（第8章）を網羅していれば OK とする
- `go test -cover` の数値は参考情報として CI で出力するに留める
- レビューでは「何を検証しているか」を見る。アサーションのないテスト・getter/setter だけのテスト等は指摘対象

---

## 12. アンチパターン集

以下はレビューで指摘し、修正を求める。

| アンチパターン | なぜダメか |
|--------------|----------|
| `time.Sleep` でタイミング待ち | 不安定。`Clock` インターフェースで時刻制御する |
| `time.Now()` / `uuid.New()` をプロダクトコードで直接呼ぶ | テストで固定できない |
| テスト間で状態を共有（パッケージ変数・グローバル `map`） | 並列実行で壊れる |
| 1 つのテスト関数で複数のシナリオを混ぜる | 失敗原因が分かりにくい。テーブル駆動に分ける |
| handler テストで実 DB を起動 | ユニットではない（integration へ移す） |
| `http.DefaultClient` で外部 API を叩く | テストが不安定かつ遅い。モックする |
| エラー検証が `err != nil` だけ | エラー種別（`code`）まで検証する |
| `assert.True(t, len(xs) == 3)` | `assert.Len(t, xs, 3)` を使う |
| 日本語のテスト関数名 | grep しづらい。`name` フィールドに日本語を書く |
| アサーションのないテスト | 何を検証しているか不明 |
| モックの戻り値設定だけして呼び出し検証をしない | 「呼ばれない」ことを検証すべき場面（ルール8等）で漏れる |

---

## 改訂履歴

| 日付 | 内容 |
|------|------|
| 2026-05-11 | 初版 |
