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

## Linux / WSL 環境での注意

`pnpm install` 後に以下のエラーが出る場合、oxc のネイティブバインディングが自動インストールされなかったことが原因。

```
Cannot find module './parser.linux-x64-gnu.node'
```

手動で追加する：

```bash
pnpm add -D \
  "@oxc-parser/binding-linux-x64-gnu@0.129.0" \
  "@oxc-transform/binding-linux-x64-gnu@0.129.0" \
  "@oxc-minify/binding-linux-x64-gnu@0.129.0"
```

これらは `devDependencies` に固定済みのため、`pnpm install` を再実行すれば以降は自動で入る。
