---
name: openapi-update
description: OpenAPI仕様変更後のコード再生成ワークフロー。docs/api/openapi.yaml を変更したら必ず実行する。
disable-model-invocation: true
---

`docs/api/openapi.yaml` を変更した後、以下を順番に実行する。

## 手順

1. **backend 型再生成**
   ```bash
   make oapi-codegen
   ```
   生成先：`backend/internal/apigen/`。生成物はコミット対象。

2. **frontend 型再生成**
   ```bash
   cd frontend && pnpm openapi:types
   ```
   生成先：`frontend/types/api.d.ts`（gitignore 対象。CI で自動生成されるためコミット不要）。

3. **テスト確認**
   ```bash
   cd backend && go test ./...
   cd frontend && pnpm test
   ```
   両方が新スキーマで通ることを確認する。

4. **OpenAPI 変更チェックリスト**
   - [ ] `paths` にエンドポイントを追加済み
   - [ ] 対応する `tags` が存在するか確認（なければ追加）
   - [ ] 使用するスキーマを `components/schemas` に定義済み
   - [ ] エラーレスポンスは `components/schemas/Error` を参照している
