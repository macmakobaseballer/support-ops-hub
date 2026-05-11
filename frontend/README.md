# Support Ops Hub — frontend

Nuxt 3 + TypeScript + Pinia + Tailwind。

## 開発

```bash
pnpm install
pnpm dev          # http://localhost:3000
pnpm typecheck    # vue-tsc
pnpm lint         # ESLint
pnpm openapi:types  # docs/api/openapi.yaml から types/api.d.ts を生成
```

## ディレクトリ

```
frontend/
├── pages/        Nuxt ファイルベースルーティング（SCR-* を実装）
├── app.vue       ルート
├── nuxt.config.ts
├── package.json
└── types/        (M1) openapi-typescript 生成（gitignore 対象）
```

M0 時点では `pages/index.vue` の仮トップ画面のみ。実装は M2 以降。
