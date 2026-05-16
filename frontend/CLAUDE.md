# Support Ops Hub — frontend / CLAUDE.md

Nuxt 3 フロントエンドで作業する際のルール。**これらは提案ではなくハードコンストレイントである。** 逸脱する前に必ずユーザーに確認すること。

リポジトリ全体のルールは [../CLAUDE.md](../CLAUDE.md) を参照。

---

## 構成

| 項目 | 内容 |
|------|------|
| フレームワーク | Nuxt 3（`compatibilityDate: 2025-05-01`） |
| 言語 | TypeScript（`strict: true`） |
| パッケージマネージャ | pnpm 10.13.1（npm / yarn は使わない） |
| 状態管理 | Pinia（`@pinia/nuxt`） |
| スタイル | Tailwind CSS（`@nuxtjs/tailwindcss`） |
| ユーティリティ | VueUse（`@vueuse/nuxt`） |
| Lint | ESLint（`@nuxt/eslint`） |
| API 型生成 | openapi-typescript（`pnpm openapi:types`） |

ディレクトリ構成と開発コマンドは [README.md](README.md) を参照。

---

## 画面と仕様の対応

各 `pages/` のページは [../docs/screens/画面仕様書.md](../docs/screens/画面仕様書.md) の SCR-* と 1:1 対応させる。下表は推奨ファイルパス（命名上の制約がなければそのまま採用）。

| Screen | 想定ファイル |
|--------|------------|
| SCR-01 ログイン | `pages/login.vue` |
| SCR-02 ダッシュボード | `pages/index.vue` |
| SCR-03 チケット一覧 | `pages/tickets/index.vue` |
| SCR-04 チケット詳細 | `pages/tickets/[id].vue` |
| SCR-05 チケット作成/編集モーダル | `components/tickets/TicketFormModal.vue` |
| SCR-06 顧客別レポート | `pages/customers/[id]/report.vue` |
| SCR-07 システム別ダッシュボード | `pages/systems/[id]/dashboard.vue` |
| SCR-08 顧客マスタ | `pages/customers/index.vue` |
| SCR-09 顧客作成/編集モーダル | `components/customers/CustomerFormModal.vue` |
| SCR-10 システムマスタ | `pages/systems/index.vue` |
| SCR-11 システム作成/編集モーダル | `components/systems/SystemFormModal.vue` |
| SCR-12 ユーザー管理 | `pages/admin/users/index.vue` |

実装時は対応する SCR を必ず確認すること。仕様と乖離する場合は画面仕様書を先に更新する。

---

## アーキテクチャルール

### ルール F1：API は集約レイヤー経由で呼ぶ

`useFetch` / `$fetch` をページ・コンポーネントから直接呼ばない。必ず以下の層を経由する：

```
component → store / composable → API client (lib/api/)
```

理由：
- テスト容易性（モック差し替えが容易、[../docs/testing/frontend-unit-test.md](../docs/testing/frontend-unit-test.md) 第7章）
- 型安全性（`types/api.d.ts` を一箇所で使う）
- 認証ヘッダ・エラー変換のロジックを一箇所に集約

### ルール F2：型は OpenAPI から生成する

リクエスト・レスポンスの型は手書きしない。`pnpm openapi:types` で [../docs/api/openapi.yaml](../docs/api/openapi.yaml) から `types/api.d.ts` を生成し、それを使う。

- `types/api.d.ts` は **gitignore 対象**。CI で常に生成する
- backend 側の OpenAPI 更新があったら必ず再生成する
- カスタム型が必要な場合は `types/api.d.ts` を import して派生型を定義する

### ルール F3：エラー表示はバックエンドの `code` を解釈する

API クライアントは backend のエラーレスポンス（`{ code, message, details }`、[../backend/CLAUDE.md](../backend/CLAUDE.md) ルール4）をパースし、画面側は以下の方針で扱う：

- `message` を基本的にそのまま表示（日本語ユーザー向けメッセージ）
- フォームバリデーションエラー（`details` 直下の `{ フィールド名: メッセージ }` 構造、[../backend/CLAUDE.md](../backend/CLAUDE.md) ルール4）は対応フィールド横に表示
- 401 はログイン画面へリダイレクト
- 403 はトーストまたは画面メッセージで「権限がありません」を表示
- 5xx は共通のエラートーストを表示

エラー処理ロジックは API クライアント層に集約し、各画面で重複実装しない。

### ルール F4：ステータス表示制御はバックエンドルールと一致させる

チケットのステータス遷移ボタンは、[../backend/CLAUDE.md](../backend/CLAUDE.md) ルール5の遷移マトリクスと一致させる：

| 現在のステータス | 表示するアクション |
|------------|----------------|
| new | 「対応中にする」 |
| in_progress | 「保留にする」「完了にする」 |
| waiting | 「対応中に戻す」「完了にする」 |
| done | （ステータス変更ボタンなし） |

**バックエンド側で検証されるが、フロントは表示制御として独立して実装する。** 不正遷移を試みる UI を表示しない。

### ルール F5：admin 専用 UI の表示制御

[../backend/CLAUDE.md](../backend/CLAUDE.md) ルール6 に対応する admin 専用エンドポイントを呼ぶ UI は、ロールベースの表示制御を行う：

- ユーザー管理画面（SCR-12）自体のアクセス制御
- 顧客作成・編集ボタン（SCR-08）
- システム作成・編集ボタン（SCR-10）
- システム担当者編集（SCR-11）

ロール情報は Pinia の auth store に保持し、画面側は `v-if="auth.isAdmin"` 等で制御する。

### ルール F6：`customer_id` フォームのロック

システム編集モーダル（SCR-11）では、編集モード時に「顧客企業」セレクトを `disabled` にする（[../backend/CLAUDE.md](../backend/CLAUDE.md) ルール8）。

新規作成モード時のみ編集可能。

---

## 状態管理（Pinia）

### ストア配置

```
stores/
├── auth.ts          ← ログインユーザー情報、ロール（シングルトン状態のため単数形）
├── tickets.ts       ← チケット一覧・詳細キャッシュ
├── customers.ts     ← 顧客マスタ
├── systems.ts       ← システムマスタ
└── users.ts         ← ユーザー管理（admin 専用）
```

### ストア設計方針

- `setup` ストア（Composition API スタイル）を使用する
- API 呼び出しは action 内で行う（コンポーネントから直接 API client を呼ばない）
- 楽観的更新を入れる場合は失敗時のロールバックを必ず実装する

---

## コンポーネント設計

- `components/` 配下は機能別にディレクトリ分け（`components/tickets/`、`components/customers/` 等）
- 表示専用コンポーネント（Atom 相当）は `components/ui/` 配下
- props は必ず TypeScript の型を付ける（`defineProps<{ ... }>()`）
- emit も型を付ける（`defineEmits<{ update: [value: string] }>()`）
- `data-testid` 属性をユーザー操作要素に付与（[../docs/testing/frontend-unit-test.md](../docs/testing/frontend-unit-test.md) 第5.3節）

---

## スタイル

- Tailwind CSS のユーティリティクラスを基本とする
- カスタム CSS は最小限に。必要な場合は `assets/css/` に置く
- カラー・余白等のデザイントークンは Tailwind の設定で統一する

---

## ランタイム設定

API ベース URL は `runtimeConfig.public.apiBase` で管理する：

```ts
// nuxt.config.ts
runtimeConfig: {
  public: {
    apiBase: process.env.NUXT_PUBLIC_API_BASE ?? 'http://127.0.0.1:8000/api/v1',
  },
},
```

API クライアント以外の場所で `useRuntimeConfig()` を呼んで URL を組み立てない。

---

## テスト

詳細は [../docs/testing/frontend-unit-test.md](../docs/testing/frontend-unit-test.md) を参照。

主要原則：
- ステータス遷移ボタン・admin 表示制御・`customer_id` ロック等の必須パターンを網羅
- API クライアントを `vi.mock` で差し替え（実ネットワーク呼び出し禁止）
- スナップショットテストは原則禁止
- カバレッジ数値目標は設けず、必須パターン優先

---

## ルール F7：UI デザインシステム

**プロトタイプ（[`prototype/index.html`](../prototype/index.html)）のビジュアルを忠実に再現すること。** 色・余白・コンポーネント形状はハードコンストレイント。デザインの独自判断は禁止。

### カラーパレット

| トークン名 | 値 | 用途 |
|---|---|---|
| `primary` | `#005461` | サイドバー背景・ボタン hover |
| `secondary` | `#018790` | ボタン・リンク・sidebar hover |
| `accent` | `#00B7B5` | sidebar active・progress bar・Chart 境界線 |
| `page-bg` | `#F4F4F4` | ページ背景 |

デザイントークンは `tailwind.config.ts` の `theme.extend.colors` に定義済み。コンポーネント内でハードコードしない。

### ステータスバッジ

| ステータス | 背景色 | 文字色 | クラス |
|---|---|---|---|
| 新規受付（new） | `#dbeafe` | `#1d4ed8` | `.badge-status-new` |
| 対応中（in_progress） | `#f5ede4` | `#7c4d33` | `.badge-status-in_progress` |
| 確認待ち（waiting） | `#ede9fe` | `#7c3aed` | `.badge-status-waiting` |
| 完了（done） | `#dcfce7` | `#166534` | `.badge-status-done` |

### 優先度バッジ

| 優先度 | 背景色 | 文字色 | クラス |
|---|---|---|---|
| 高（high） | `#fee2e2` | `#dc2626` | `.badge-priority-high` |
| 中（medium） | `#fef3c7` | `#b45309` | `.badge-priority-medium` |
| 低（low） | `#f0fdf4` | `#166534` | `.badge-priority-low` |

バッジはすべて `font-weight: 500`、`border-radius` は Tailwind の `rounded` を使用。

### レイアウト構造

```
┌─────────────────────────────────────────────┐
│  .sidebar (固定, 240px)  │  .main-content   │
│  ┌─────────────────────┐ │  ┌─────────────┐ │
│  │ .sidebar-logo       │ │  │ .page-header│ │ ← sticky
│  ├─────────────────────┤ │  ├─────────────┤ │
│  │ .sidebar-section    │ │  │ .content-   │ │
│  │ .sidebar-item       │ │  │  area       │ │
│  │ .customer-header    │ │  │             │ │
│  │   .system-nav-item  │ │  └─────────────┘ │
│  ├─────────────────────┤ │                  │
│  │ .sidebar-footer     │ │                  │
│  └─────────────────────┘ │                  │
└─────────────────────────────────────────────┘
```

- `margin-left` は `240px`（`spacing.sidebar` トークン）固定。レスポンシブ対応は MVP スコープ外。
- ページヘッダーは `position: sticky; top: 0; z-index: 100`。

### コンポーネントクラス一覧

以下のクラスはすべて `assets/css/main.css` の `@layer components` に定義済み。新規コンポーネントでは必ず使用し、個別に再定義しない。

| クラス | 用途 | 主なスタイル |
|---|---|---|
| `.stat-card` | KPI数値カード | bg-white、border-left 4px、radius 10px、shadow |
| `.stat-label` | stat-card の小見出し | 0.78rem、uppercase、#64748b |
| `.stat-value` | stat-card の数値 | 2.2rem、font-bold |
| `.data-table` | テーブルラッパー | bg-white、radius 10px、overflow-hidden |
| `.filter-bar` | フィルター行 | bg-white、radius 10px、flex wrap |
| `.form-card` | フォーム・詳細カード | bg-white、radius 10px、padding 24px |
| `.form-section-title` | form-card 内セクション見出し | 0.75rem、uppercase、border-bottom |
| `.chart-card` | グラフカード | bg-white、radius 10px、padding 20px |
| `.chart-card-title` | chart-card 内タイトル | 0.75rem、uppercase、#64748b |
| `.comment-item` | 対応履歴の1件 | border-bottom #f1f5f9 |
| `.comment-meta` | コメント投稿者・日時 | 0.78rem、#94a3b8 |
| `.login-container` | ログイン背景 | グラデーション(135deg, #005461→#018790) |
| `.login-card` | ログインカード | bg-white、radius 12px、400px、shadow |
| `.btn-primary` | 主要アクションボタン | bg #018790、hover #005461 |
| `.btn-outline-primary` | サブアクションボタン | border/text #018790 |

### stat-card の border-left カラー

stat-card の左ボーダー色はステータスごとに inline style で指定する（Tailwind の JIT では動的クラスが purge されるため）：

```vue
<div class="stat-card" style="border-color: #3b82f6;">  <!-- 新規受付 -->
<div class="stat-card" style="border-color: #7c4d33;">  <!-- 対応中 -->
<div class="stat-card" style="border-color: #8b5cf6;">  <!-- 確認待ち -->
<div class="stat-card" style="border-color: #22c55e;">  <!-- 完了 -->
```

### アイコン

Bootstrap Icons（`bi-*`）を使用する。パッケージ: `bootstrap-icons`。

```html
<!-- Vue コンポーネントでの使用例 -->
<i class="bi bi-speedometer2"></i>   <!-- ダッシュボード -->
<i class="bi bi-ticket-perforated"></i>  <!-- 問い合わせ -->
<i class="bi bi-building"></i>       <!-- 顧客企業 -->
<i class="bi bi-server"></i>         <!-- システム -->
<i class="bi bi-people"></i>         <!-- ユーザー管理 -->
<i class="bi bi-headset" style="color:#3b82f6;"></i>  <!-- ロゴ -->
```

CSS は `assets/css/main.css` で `@import "bootstrap-icons/font/bootstrap-icons.css"` を追加する（パッケージインストール後）。

### Chart.js カラー

| グラフ | 設定 |
|---|---|
| 棒グラフ（月次） | `backgroundColor: '#00B7B530'`、`borderColor: '#00B7B5'`、`borderWidth: 2`、`borderRadius: 4` |
| ドーナツ（種別内訳） | `['#00B7B5', '#ef4444', '#f59e0b', '#018790']`（操作方法・バグ・設定変更・データ修正） |

### サイドバーフッターのロールバッジ

```vue
<!-- admin -->
<span style="font-size:0.65rem; background:#00B7B5; color:white; border-radius:4px; padding:1px 6px; font-weight:600;">管理者</span>
<!-- member -->
<span style="font-size:0.65rem; background:rgba(255,255,255,0.2); color:rgba(255,255,255,0.8); border-radius:4px; padding:1px 6px; font-weight:600;">担当者</span>
```

---

## ドキュメント更新義務

- 画面の追加・変更時は [../docs/screens/画面仕様書.md](../docs/screens/画面仕様書.md) を必ず更新する
- API 型を変更した場合は backend 側で [../docs/api/openapi.yaml](../docs/api/openapi.yaml) を更新後、`pnpm openapi:types` で型を再生成する

---

## CI

frontend のビルド・型チェックは [`.github/workflows/frontend.yml`](../.github/workflows/frontend.yml) で実行される。詳細は [../.github/CLAUDE.md](../.github/CLAUDE.md) を参照。
