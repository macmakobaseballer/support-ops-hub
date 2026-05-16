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

`prototype/index.html` のビジュアルを忠実に再現すること。色・余白・コンポーネント形状はハードコンストレイント。デザインの独自判断は禁止。

カラーパレット・バッジ・レイアウト・CSSクラス一覧・アイコン・Chart.js 設定の詳細はスキル `/ui-design-system` を参照。

---

## ドキュメント更新義務

- 画面の追加・変更時は [../docs/screens/画面仕様書.md](../docs/screens/画面仕様書.md) を必ず更新する
- API 型を変更した場合は backend 側で [../docs/api/openapi.yaml](../docs/api/openapi.yaml) を更新後、`pnpm openapi:types` で型を再生成する

---

## CI

frontend のビルド・型チェックは [`.github/workflows/frontend.yml`](../.github/workflows/frontend.yml) で実行される。詳細は [../.github/CLAUDE.md](../.github/CLAUDE.md) を参照。
