# フロントエンド ユニットテスト規約

Nuxt.js フロントエンドのユニットテストに関する規約を定義する。
[backend/CLAUDE.md](../../backend/CLAUDE.md) のアーキテクチャルールと同様、**これらは提案ではなくハードコンストレイントである。** 逸脱する前に必ずチームで合意すること。

---

## 1. 目的とスコープ

本ドキュメントが扱う「ユニットテスト」とは、**ブラウザ・実 API・実ルーティングを介さず、Node.js プロセス内で完結するテスト**を指す。

| 種別 | 本ドキュメントの対象 | 依存 |
|------|------------------|------|
| ユニットテスト | ✅ 対象 | なし（モック使用）。DOM は `happy-dom` |
| コンポーネント結合テスト | △ 一部対象 | 複数コンポーネントを組み合わせた `mount` テストはここに含める |
| E2E テスト | ❌ 対象外 | ブラウザ＋起動中のバックエンド |

E2E の規約は別途定義する。

---

## 2. テストフレームワーク

以下のライブラリを標準採用する。**新しいライブラリの導入は事前にチームで合意すること。**

| 用途 | ライブラリ |
|------|-----------|
| テストランナー | `vitest` |
| コンポーネントマウント | `@vue/test-utils` |
| Nuxt 機能モック | `@nuxt/test-utils`（`mockNuxtImport` 等） |
| DOM 実装 | `happy-dom`（`jsdom` より速い） |
| アサーション | `vitest` 標準の `expect` |
| Pinia セットアップ | `@pinia/testing` |

`@testing-library/vue` は採用しない（`@vue/test-utils` の `wrapper.find` / `wrapper.trigger` で十分。ライブラリ二重化を避ける）。

---

## 3. テスト対象と優先順位

優先度の高い順にテストを書く。下に行くほどテストの ROI（労力対効果）が下がる。

| 優先度 | 対象 | テストする内容 |
|-------|------|------------|
| 1 | `utils/` `lib/` の純関数 | 入出力・エッジケース |
| 2 | `composables/` | 状態遷移・副作用（モック注入） |
| 3 | Pinia store | actions の副作用・getter のロジック |
| 4 | 表示専用コンポーネント（Atom/Molecule 相当） | props → DOM・emit イベント |
| 5 | フォーム・モーダル等の振る舞いを持つコンポーネント | ユーザー操作・バリデーション |
| 6 | ページコンポーネント (`pages/`) | E2E でカバー。ユニットでは最小限 |

ページコンポーネントを過剰にユニットでテストしない。複数コンポーネント・ストア・ルーティングが絡む検証は E2E に委ねる。

---

## 4. ファイル配置・命名規約

### 4.1 ファイル配置

- テストファイルは **テスト対象と同じディレクトリ**に置く
  - `components/TicketCard.vue` → `components/TicketCard.spec.ts`
  - `composables/useTicket.ts` → `composables/useTicket.spec.ts`
- 拡張子は `.spec.ts` で統一（`.test.ts` は使わない）

### 4.2 テスト構成

`describe` / `it` の階層は **対象 → シナリオ → 期待結果** の順で記述する。

```ts
describe('TicketCard', () => {
  describe('ステータス表示', () => {
    it('status が new のとき「未対応」バッジを表示する', () => { ... })
    it('status が done のとき「完了」バッジを表示する', () => { ... })
  })

  describe('クリック時の挙動', () => {
    it('カードクリックで click イベントを emit する', () => { ... })
  })
})
```

- `it` の説明は **日本語**で書く（テスト失敗時のメッセージが読みやすい）
- `it('should ...')` のような英語テンプレは使わない

---

## 5. コンポーネントテストの方針

### 5.1 `mount` を基本とする

- デフォルトは `mount`。子コンポーネントも含めてレンダリングする
- `shallowMount` は **子コンポーネントが重く、当該テストで関係がない** 場合のみ使う（理由をコメントで明記）

### 5.2 検証するもの・しないもの

| ✅ 検証する | ❌ 検証しない |
|----------|-----------|
| レンダリングされた DOM の内容 | コンポーネント内部の `ref` / `data` の値 |
| emit されたイベントとペイロード | 内部メソッドの呼び出し |
| props の変更による再レンダリング | 実装で使っているライブラリの API |
| ユーザー操作（クリック・入力）の結果 | CSS クラス名そのもの（意味のあるテストデータ属性で検証する） |

**実装詳細ではなく外から見える振る舞いを検証する。** リファクタリングでテストが壊れたら、それは大体テストが間違っている。

### 5.3 要素特定の優先順位

```
1. role / aria-label による特定（アクセシビリティ）
2. data-testid 属性
3. テキスト内容
4. （最終手段）CSS セレクタ
```

`wrapper.find('.btn-primary')` のような CSS クラス依存の特定は禁止。`data-testid="submit-button"` を付与して特定する。

### 5.4 ユーザー操作

```ts
await wrapper.find('[data-testid="status-select"]').setValue('in_progress')
await wrapper.find('[data-testid="submit-button"]').trigger('click')
expect(wrapper.emitted('update')).toEqual([[{ status: 'in_progress' }]])
```

非同期 DOM 更新を待つために `await` を必ず付ける。

---

## 6. composables / Pinia store のテスト

### 6.1 composables

純粋に setup 関数として呼び出す。Nuxt の auto-import に依存しないようテスト内で明示 import する。

```ts
import { useTicketStatus } from './useTicketStatus'

describe('useTicketStatus', () => {
  it('done への遷移後は canTransitionTo が false を返す', () => {
    const { status, canTransitionTo, transitionTo } = useTicketStatus('new')
    transitionTo('in_progress')
    transitionTo('done')
    expect(canTransitionTo('in_progress')).toBe(false)
  })
})
```

副作用（API 呼び出し）を含む composable は依存をオプション引数で注入できるようにし、テストではモックを渡す。

### 6.2 Pinia store

`@pinia/testing` の `createTestingPinia` を使う。

```ts
import { createTestingPinia } from '@pinia/testing'
import { setActivePinia } from 'pinia'

beforeEach(() => {
  setActivePinia(createTestingPinia({ stubActions: false }))
})
```

- actions のロジックを検証したいときは `stubActions: false`
- コンポーネントから store を使うテストでは `stubActions: true`（デフォルト）にして actions の呼び出しだけ検証する

---

## 7. API クライアントのモック戦略

### 7.1 基本方針：クライアントモジュールを `vi.mock`

API クライアントを `lib/api/` 等にまとめ、テストでは `vi.mock` でモジュールごと差し替える。

```ts
import { vi } from 'vitest'
import * as ticketApi from '~/lib/api/ticket'

vi.mock('~/lib/api/ticket')

it('fetch エラー時にエラーバナーを表示する', async () => {
  vi.mocked(ticketApi.listTickets).mockRejectedValue({
    code: 'INTERNAL_ERROR',
    message: 'サーバーエラーが発生しました',
  })
  // ...
})
```

### 7.2 MSW は使わない（ユニットでは）

[MSW](https://mswjs.io/) によるネットワーク層モックはユニットテストでは導入しない。理由：

- ユニットでは「コンポーネントが API クライアントをどう使うか」を検証すれば十分
- MSW は E2E / コンポーネント結合テスト側で別途検討する

---

## 8. Nuxt 機能のモック

`@nuxt/test-utils` の `mockNuxtImport` で Nuxt の auto-import を差し替える。

| Nuxt 機能 | モック方法 |
|----------|-----------|
| `useRoute` | `mockNuxtImport('useRoute', () => () => ({ params: { id: '1' } }))` |
| `useRouter` / `navigateTo` | 同上。push/replace の呼び出しは `vi.fn()` で記録 |
| `useFetch` / `$fetch` | API クライアント経由でアクセスする方針（第7章）なので、原則直接モックしない |
| `useNuxtApp` | プラグイン依存のテストでのみモック |
| `useRuntimeConfig` | 環境変数依存のテストでのみモック |

**ページコンポーネントが `useFetch` を直接呼ぶ実装は避け、composable / store / API クライアント層に集約する。** これによりテスト容易性が大きく向上する。

---

## 9. スナップショットテスト

**基本禁止。** 以下の理由で運用が破綻しやすい：

- DOM の些細な変更で大量に壊れ、内容を確認せず `--update` してしまいがち
- 「何を検証しているか」がスナップショットファイルからは読み取れない

**例外として許可される用途：**

- 表示専用かつ DOM が固定された小さなコンポーネント（アイコン、バッジ等）
- 生成された SVG 文字列の固定検証
- ライブラリの出力など、人間が定義した期待値より生成物が信頼できるケース

例外を使う場合は、テスト内に **なぜスナップショットを使うか** のコメントを書くこと。

---

## 10. 必須テストパターン

[backend/CLAUDE.md](../../backend/CLAUDE.md) のアーキテクチャルールと画面仕様に対応する以下のテストは **必ず書く**。

### 10.1 ステータス遷移ボタンの表示制御（バックエンドルール5に対応）

チケット詳細画面（[SCR-04](../screens/screenshots/SCR-04-ticket-detail.png) / [SCR-05](../screens/screenshots/SCR-05-ticket-modal.png)）のステータス変更 UI が、現在のステータスに応じて遷移可能なボタンだけを表示することを検証する。

| 現在のステータス | 表示すべきボタン |
|------------|-------------|
| new | 「対応中にする」 |
| in_progress | 「保留にする」「完了にする」 |
| waiting | 「対応中に戻す」「完了にする」 |
| done | （ステータス変更ボタンなし） |

**注意：** バックエンドでも検証されるが（ルール5）、フロントは表示制御として独立してテストする。

### 10.2 権限による要素表示（backend/CLAUDE.md ルール6）

ユーザーロール（admin / member）によって表示が変わる UI を網羅する：

- 顧客マスタ画面（[SCR-08](../screens/screenshots/SCR-08-customer-master.png)）の「新規作成」ボタンが admin のみに表示される
- システムマスタ画面（[SCR-10](../screens/screenshots/SCR-10-system-master.png)）の編集ボタンが admin のみに表示される
- ユーザー管理画面（[SCR-12](../screens/screenshots/SCR-12-user-management.png)）自体が admin のみアクセス可能

ロールは Pinia store のモックで切り替えてテストする。

### 10.3 API エラーレスポンスのユーザー表示（backend/CLAUDE.md ルール4）

API クライアントが `{ code, message, details }` 形式のエラーを返したとき：

- `message` フィールドの内容が画面に表示される
- HTTP ステータスや `code` に応じた振る舞い（401 でログイン画面遷移、422 でフォームエラー表示、5xx でトースト表示等）

エラーハンドリングのテストはエラー種別ごとに網羅する。

### 10.4 `customer_id` フォーム項目の編集ロック（backend/CLAUDE.md ルール8）

システム編集モーダル（[SCR-11](../screens/screenshots/SCR-11-system-modal.png)）で：

- 新規作成時：顧客企業の選択肢が編集可能
- 編集時：顧客企業の選択肢が `disabled` で表示される

### 10.5 検索系一覧の `page`/`per_page` 連動

チケット一覧（[SCR-03](../screens/screenshots/SCR-03-ticket-list.png)）等のページネーション UI が：

- ページ送り時に store / API クライアントへ `page` / `per_page` を渡す
- レスポンスの `pagination.total_pages` が表示に反映される

---

## 11. アサーション規約

| 用途 | 推奨 |
|------|-----|
| プリミティブ値の比較 | `expect(x).toBe(y)` |
| オブジェクト・配列の比較 | `expect(x).toEqual(y)` |
| 部分一致 | `expect(x).toMatchObject(y)` |
| 要素の表示判定 | `expect(wrapper.find('[data-testid="..."]').exists()).toBe(true)` |
| テキスト内容 | `expect(wrapper.text()).toContain('完了')` |
| emit 検証 | `expect(wrapper.emitted('click')).toHaveLength(1)` |
| 非同期エラー | `await expect(fn()).rejects.toThrow()` |

- 1 つの `it` で **assertion を 10 個以上書かない**。検証観点が多すぎる場合は `it` を分割する
- `expect(x).toBe(true)` だけのテスト（何を比較しているか分からない）は禁止

---

## 12. テストヘルパー

### 12.1 マウントファクトリ

複数のテストで同じ props・provide・global 設定を使うコンポーネントは、マウントファクトリを用意する：

```ts
// components/__tests__/factory.ts
export function mountTicketCard(overrides: Partial<Props> = {}) {
  return mount(TicketCard, {
    props: { ticket: defaultTicket, ...overrides },
    global: {
      plugins: [createTestingPinia()],
    },
  })
}
```

### 12.2 fixture

API レスポンス・ドメインオブジェクトの fixture は `tests/fixtures/` または対象隣接の `__fixtures__/` に置く。builder（functional options）パターン推奨：

```ts
export function makeTicket(overrides: Partial<Ticket> = {}): Ticket {
  return { id: 'tkt_default', status: 'new', /* ... */, ...overrides }
}
```

---

## 13. カバレッジ

**数値目標は設けない。** API 側（[docs/testing/api-unit-test.md](api-unit-test.md)）と同じ方針：

- 必須テストパターン（第10章）を網羅していれば OK
- `vitest --coverage` の数値は CI で参考情報として出力する
- レビューでは「何を検証しているか」を見る

---

## 14. アンチパターン集

以下はレビューで指摘し、修正を求める。

| アンチパターン | なぜダメか |
|--------------|----------|
| `await new Promise(r => setTimeout(r, 100))` で DOM 更新待ち | 不安定。`await wrapper.vm.$nextTick()` か `await flushPromises()` を使う |
| `Date.now()` を直接呼ぶプロダクトコード | テストで固定できない。日時取得は composable / util 経由に集約してモック可能にする |
| `vi.useFakeTimers()` を `afterEach` で戻し忘れる | 他テストに影響。`afterEach(() => vi.useRealTimers())` を徹底 |
| CSS クラス名で要素を特定 (`wrapper.find('.text-red-500')`) | スタイル変更で壊れる。`data-testid` を使う |
| コンポーネントの内部 `ref` を直接読む | 実装詳細。DOM 出力か emit で検証する |
| スナップショットの安易な利用 | 何を検証しているか分からない。第9章参照 |
| `it.only` / `describe.only` のコミット | 他のテストが実行されなくなる。lint で検出する |
| 巨大な `mount` 結果に対する膨大なアサーション | テストの意図が散らかる。`it` を分割する |
| `await wrapper.trigger('click')` の `await` 忘れ | DOM が更新前に検証され、偽陽性/陰性が出る |
| `useFetch` をページ内で直接呼ぶ実装 | テストできない。composable / store / API クライアント層に集約する |
| API クライアントを `vi.mock` せず実 fetch を期待 | ユニットでネットワークに依存する。テスト不安定 |

---

## 改訂履歴

| 日付 | 内容 |
|------|------|
| 2026-05-11 | 初版 |
