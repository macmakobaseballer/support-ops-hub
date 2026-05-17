# Frontend Unit Test Results — M3

実行日時: 2026-05-17

## 結果サマリー

| テストファイル | 結果 | テスト数 |
|---|---|---|
| tests/smoke.test.ts | ✅ PASS | 1 |
| utils/ticketStatus.spec.ts | ✅ PASS | 13 |
| utils/ticketType.spec.ts | ✅ PASS | 5 |
| components/tickets/TicketFormModal.spec.ts | ✅ PASS | 4 |
| **合計** | **✅ 全 PASS** | **23** |

## M3 新規テスト

### ticketStatus ユーティリティ（ステータス遷移ボタン表示制御 F4）
- new → in_progress のみ遷移可能 ✅
- in_progress → waiting / done ✅
- waiting → in_progress / done ✅
- done → 遷移なし（終端） ✅
- new → waiting は不可 ✅
- new → done は不可 ✅
- ラベル・バッジクラス・ボタンラベル各マッピング全 PASS ✅

### ticketType ユーティリティ（種別→優先度自動セット）
- bug → high ✅
- question → medium ✅
- config → medium ✅
- data → medium ✅

### TicketFormModal コンポーネント（F4・F6 対応）
- タイトルが空のまま保存 → バリデーションエラー表示 ✅
- 編集時は顧客企業セレクトが disabled ✅
- 編集時はシステムセレクトが disabled ✅
- 新規登録時は顧客企業セレクトが enabled ✅

## 詳細ログ


> support-ops-hub-frontend@ test /home/makoj/dev/support-ops-hub/frontend
> vitest run


 RUN  v2.1.9 /home/makoj/dev/support-ops-hub/frontend

 ✓ utils/ticketType.spec.ts (5 tests) 5ms
 ✓ tests/smoke.test.ts (1 test) 3ms
 ✓ utils/ticketStatus.spec.ts (13 tests) 11ms
 ✓ components/tickets/TicketFormModal.spec.ts (4 tests) 62ms

 Test Files  4 passed (4)
      Tests  23 passed (23)
   Start at  15:10:02
   Duration  1.65s (transform 516ms, setup 0ms, collect 640ms, tests 81ms, environment 1.82s, prepare 775ms)

