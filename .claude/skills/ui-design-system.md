---
name: ui-design-system
description: UIデザインシステム仕様（プロトタイプ準拠）。コンポーネント実装・スタイル変更時に参照。カラーパレット・バッジ・レイアウト・CSSクラス・アイコン・Chart設定を含む。
---

# UIデザインシステム

**`prototype/index.html` のビジュアルを忠実に再現すること。** 色・余白・コンポーネント形状はハードコンストレイント。デザインの独自判断は禁止。

## カラーパレット

| トークン名 | 値 | 用途 |
|---|---|---|
| `primary` | `#005461` | サイドバー背景・ボタン hover |
| `secondary` | `#018790` | ボタン・リンク・sidebar hover |
| `accent` | `#00B7B5` | sidebar active・progress bar・Chart 境界線 |
| `page-bg` | `#F4F4F4` | ページ背景 |

デザイントークンは `tailwind.config.ts` の `theme.extend.colors` に定義済み。コンポーネント内でハードコードしない。

## ステータスバッジ

| ステータス | 背景色 | 文字色 | クラス |
|---|---|---|---|
| 新規受付（new） | `#dbeafe` | `#1d4ed8` | `.badge-status-new` |
| 対応中（in_progress） | `#f5ede4` | `#7c4d33` | `.badge-status-in_progress` |
| 確認待ち（waiting） | `#ede9fe` | `#7c3aed` | `.badge-status-waiting` |
| 完了（done） | `#dcfce7` | `#166534` | `.badge-status-done` |

## 優先度バッジ

| 優先度 | 背景色 | 文字色 | クラス |
|---|---|---|---|
| 高（high） | `#fee2e2` | `#dc2626` | `.badge-priority-high` |
| 中（medium） | `#fef3c7` | `#b45309` | `.badge-priority-medium` |
| 低（low） | `#f0fdf4` | `#166534` | `.badge-priority-low` |

バッジはすべて `font-weight: 500`、`border-radius` は Tailwind の `rounded` を使用。

## レイアウト構造

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

## コンポーネントクラス一覧

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

## stat-card の border-left カラー

stat-card の左ボーダー色はステータスごとに inline style で指定する（Tailwind の JIT では動的クラスが purge されるため）：

```vue
<div class="stat-card" style="border-color: #3b82f6;">  <!-- 新規受付 -->
<div class="stat-card" style="border-color: #7c4d33;">  <!-- 対応中 -->
<div class="stat-card" style="border-color: #8b5cf6;">  <!-- 確認待ち -->
<div class="stat-card" style="border-color: #22c55e;">  <!-- 完了 -->
```

## アイコン

Bootstrap Icons（`bi-*`）を使用する。パッケージ: `bootstrap-icons`。

```html
<i class="bi bi-speedometer2"></i>      <!-- ダッシュボード -->
<i class="bi bi-ticket-perforated"></i> <!-- 問い合わせ -->
<i class="bi bi-building"></i>          <!-- 顧客企業 -->
<i class="bi bi-server"></i>            <!-- システム -->
<i class="bi bi-people"></i>            <!-- ユーザー管理 -->
<i class="bi bi-headset" style="color:#3b82f6;"></i>  <!-- ロゴ -->
```

CSS は `assets/css/main.css` で `@import "bootstrap-icons/font/bootstrap-icons.css"` を追加する（パッケージインストール後）。

## Chart.js カラー

| グラフ | 設定 |
|---|---|
| 棒グラフ（月次） | `backgroundColor: '#00B7B530'`、`borderColor: '#00B7B5'`、`borderWidth: 2`、`borderRadius: 4` |
| ドーナツ（種別内訳） | `['#00B7B5', '#ef4444', '#f59e0b', '#018790']`（操作方法・バグ・設定変更・データ修正） |

## サイドバーフッターのロールバッジ

```vue
<!-- admin -->
<span style="font-size:0.65rem; background:#00B7B5; color:white; border-radius:4px; padding:1px 6px; font-weight:600;">管理者</span>
<!-- member -->
<span style="font-size:0.65rem; background:rgba(255,255,255,0.2); color:rgba(255,255,255,0.8); border-radius:4px; padding:1px 6px; font-weight:600;">担当者</span>
```
