---
name: new-service
description: 新しいバックエンドサービス追加時の必須チェックリスト。横断ルール3に基づく手順。
disable-model-invocation: true
---

新しいバックエンドサービスを追加する場合、以下を **すべて** 更新する。

## チェックリスト

- [ ] `docs/api/architecture.md` のサービスマップを更新
- [ ] `docs/api/openapi.yaml` に新タグを追加 → 変更後は `/openapi-update` を実行
- [ ] `backend/CLAUDE.md` の「所有テーブル」表に新サービスを追記
- [ ] `backend/cmd/<service>/main.go` を作成
- [ ] Phase 1：既存の単一バイナリ・ECS タスクに新サービスのモジュールをリンク（`backend/CLAUDE.md` 構成参照）
- [ ] Phase 2 で独立 ECS タスクに分離する場合は `infra/terraform/envs/*/main.tf` で配線
