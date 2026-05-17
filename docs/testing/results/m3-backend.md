# Backend Unit Test Results — M3

実行日時: 2026-05-17

## 結果サマリー

| パッケージ | 結果 | カバレッジ | テスト数 |
|---|---|---|---|
| gateway | ✅ PASS | 33.3% | 4 |
| internal/domain | ✅ PASS | 87.5% | 17 |
| internal/httperr | ✅ PASS | 84.6% | 6 |
| internal/middleware | ✅ PASS | 40.9% | 3 |
| services/auth | ✅ PASS | 80.0% | 10 |
| services/ticket | ✅ PASS | 35.3% | 14+7サブテスト |

## M3 新規テスト（services/ticket）

### ステータス遷移（ルール5 全パターン）
- `TestUpdateStatus_ValidTransitions` — 有効遷移5パターン全 PASS
  - new→in_progress, in_progress→waiting, in_progress→done, waiting→in_progress, waiting→done
- `TestUpdateStatus_InvalidTransitions_Returns422` — 無効遷移7パターン全 PASS（422 + INVALID_STATUS_TRANSITION）
  - new→waiting, new→done, done→new, done→in_progress, done→waiting, in_progress→new, waiting→new

### 権限チェック（ルール6）
- `TestCreateTicket_NoUserID_Returns401` — X-User-ID なし → 401 ✅
- `TestUpdateStatus_NoUserID_Returns401` — X-User-ID なし → 401 ✅

### エラーレスポンス形式（ルール4）
- `TestGetTicket_NotFound_HasErrorCode` — code="TICKET_NOT_FOUND" 確認 ✅

### バリデーション
- `TestCreateTicket_MissingTitle_Returns400` — タイトル必須バリデーション ✅
- `TestCreateTicket_InvalidType_Returns400` — 無効な種別 ✅
- `TestUpdateStatus_InvalidStatusValue_Returns400` — 無効なステータス文字列 ✅

### 正常系
- `TestGetTicket_Found_Returns200` — 200 + 正しい DTO ✅
- `TestCreateTicket_Valid_Returns201` — 201 + id=1 ✅
- `TestUpdateStatus_NotFound_Returns404` — 404 ✅
- `TestUpdateAssignee_NullAssignee_Returns200` — null assignee で 200 ✅

## 詳細ログ

	github.com/macmakobaseballer/support-ops-hub/backend/cmd/auth			github.com/macmakobaseballer/support-ops-hub/backend/cmd/gateway		coverage: 0.0% of statements
	github.com/macmakobaseballer/support-ops-hub/backend/cmd/ticket		=== RUN   TestDevAuth_MissingHeader_Returns401
--- PASS: TestDevAuth_MissingHeader_Returns401 (0.00s)
=== RUN   TestDevAuth_UserNotFound_Returns401
2026/05/17 14:54:21 WARN dev-auth: user not found request_id=""
--- PASS: TestDevAuth_UserNotFound_Returns401 (0.00s)
=== RUN   TestDevAuth_InactiveUser_Returns403
--- PASS: TestDevAuth_InactiveUser_Returns403 (0.00s)
=== RUN   TestDevAuth_ValidUser_InjectsHeaders
--- PASS: TestDevAuth_ValidUser_InjectsHeaders (0.00s)
PASS
coverage: 33.3% of statements
ok  	github.com/macmakobaseballer/support-ops-hub/backend/gateway	0.003s	coverage: 33.3% of statements
	github.com/macmakobaseballer/support-ops-hub/backend/internal/apigen		coverage: 0.0% of statements
	github.com/macmakobaseballer/support-ops-hub/backend/internal/auth/jwt		coverage: 0.0% of statements
	github.com/macmakobaseballer/support-ops-hub/backend/internal/config		coverage: 0.0% of statements
	github.com/macmakobaseballer/support-ops-hub/backend/internal/db/sqlc		coverage: 0.0% of statements
=== RUN   TestTicketStatus_IsValidTransition
=== RUN   TestTicketStatus_IsValidTransition/new->in_progress
=== RUN   TestTicketStatus_IsValidTransition/new->waiting
=== RUN   TestTicketStatus_IsValidTransition/new->done
=== RUN   TestTicketStatus_IsValidTransition/new->new
=== RUN   TestTicketStatus_IsValidTransition/in_progress->waiting
=== RUN   TestTicketStatus_IsValidTransition/in_progress->done
=== RUN   TestTicketStatus_IsValidTransition/in_progress->new
=== RUN   TestTicketStatus_IsValidTransition/in_progress->in_progress
=== RUN   TestTicketStatus_IsValidTransition/waiting->in_progress
=== RUN   TestTicketStatus_IsValidTransition/waiting->done
=== RUN   TestTicketStatus_IsValidTransition/waiting->new
=== RUN   TestTicketStatus_IsValidTransition/waiting->waiting
=== RUN   TestTicketStatus_IsValidTransition/done->new
=== RUN   TestTicketStatus_IsValidTransition/done->in_progress
=== RUN   TestTicketStatus_IsValidTransition/done->waiting
=== RUN   TestTicketStatus_IsValidTransition/done->done
=== RUN   TestTicketStatus_IsValidTransition/unknown->in_progress
--- PASS: TestTicketStatus_IsValidTransition (0.00s)
    --- PASS: TestTicketStatus_IsValidTransition/new->in_progress (0.00s)
    --- PASS: TestTicketStatus_IsValidTransition/new->waiting (0.00s)
    --- PASS: TestTicketStatus_IsValidTransition/new->done (0.00s)
    --- PASS: TestTicketStatus_IsValidTransition/new->new (0.00s)
    --- PASS: TestTicketStatus_IsValidTransition/in_progress->waiting (0.00s)
    --- PASS: TestTicketStatus_IsValidTransition/in_progress->done (0.00s)
    --- PASS: TestTicketStatus_IsValidTransition/in_progress->new (0.00s)
    --- PASS: TestTicketStatus_IsValidTransition/in_progress->in_progress (0.00s)
    --- PASS: TestTicketStatus_IsValidTransition/waiting->in_progress (0.00s)
    --- PASS: TestTicketStatus_IsValidTransition/waiting->done (0.00s)
    --- PASS: TestTicketStatus_IsValidTransition/waiting->new (0.00s)
    --- PASS: TestTicketStatus_IsValidTransition/waiting->waiting (0.00s)
    --- PASS: TestTicketStatus_IsValidTransition/done->new (0.00s)
    --- PASS: TestTicketStatus_IsValidTransition/done->in_progress (0.00s)
    --- PASS: TestTicketStatus_IsValidTransition/done->waiting (0.00s)
    --- PASS: TestTicketStatus_IsValidTransition/done->done (0.00s)
    --- PASS: TestTicketStatus_IsValidTransition/unknown->in_progress (0.00s)
PASS
coverage: 87.5% of statements
ok  	github.com/macmakobaseballer/support-ops-hub/backend/internal/domain	0.018s	coverage: 87.5% of statements
=== RUN   TestNew_SetsFields
--- PASS: TestNew_SetsFields (0.00s)
=== RUN   TestWrite_StatusAndContentType
--- PASS: TestWrite_StatusAndContentType (0.00s)
=== RUN   TestWrite_JSONBody
--- PASS: TestWrite_JSONBody (0.00s)
=== RUN   TestUnauthorized_Returns401
--- PASS: TestUnauthorized_Returns401 (0.00s)
=== RUN   TestForbidden_Returns403
--- PASS: TestForbidden_Returns403 (0.00s)
=== RUN   TestInternalError_Returns500
--- PASS: TestInternalError_Returns500 (0.00s)
PASS
coverage: 84.6% of statements
ok  	github.com/macmakobaseballer/support-ops-hub/backend/internal/httperr	0.005s	coverage: 84.6% of statements
=== RUN   TestRequestID_GeneratesWhenAbsent
--- PASS: TestRequestID_GeneratesWhenAbsent (0.00s)
=== RUN   TestRequestID_PreservesExistingHeader
--- PASS: TestRequestID_PreservesExistingHeader (0.00s)
=== RUN   TestRequestIDFromContext_EmptyWhenNotSet
--- PASS: TestRequestIDFromContext_EmptyWhenNotSet (0.00s)
PASS
coverage: 40.9% of statements
ok  	github.com/macmakobaseballer/support-ops-hub/backend/internal/middleware	0.020s	coverage: 40.9% of statements
	github.com/macmakobaseballer/support-ops-hub/backend/services/analytics		coverage: 0.0% of statements
=== RUN   TestHandleLogin_Success
--- PASS: TestHandleLogin_Success (0.00s)
=== RUN   TestHandleLogin_UnknownEmail_Returns401
--- PASS: TestHandleLogin_UnknownEmail_Returns401 (0.00s)
=== RUN   TestHandleLogin_InactiveUser_Returns403
--- PASS: TestHandleLogin_InactiveUser_Returns403 (0.00s)
=== RUN   TestHandleLogin_EmptyEmail_Returns400
--- PASS: TestHandleLogin_EmptyEmail_Returns400 (0.00s)
=== RUN   TestHandleLogin_InvalidJSON_Returns400
--- PASS: TestHandleLogin_InvalidJSON_Returns400 (0.00s)
=== RUN   TestHandleMe_Success
--- PASS: TestHandleMe_Success (0.00s)
=== RUN   TestHandleMe_NoHeader_Returns401
--- PASS: TestHandleMe_NoHeader_Returns401 (0.00s)
=== RUN   TestHandleMe_InvalidID_Returns401
--- PASS: TestHandleMe_InvalidID_Returns401 (0.00s)
=== RUN   TestHandleMe_UnknownID_Returns401
--- PASS: TestHandleMe_UnknownID_Returns401 (0.00s)
=== RUN   TestHandleMe_InactiveUser_Returns403
--- PASS: TestHandleMe_InactiveUser_Returns403 (0.00s)
PASS
coverage: 80.0% of statements
ok  	github.com/macmakobaseballer/support-ops-hub/backend/services/auth	0.018s	coverage: 80.0% of statements
	github.com/macmakobaseballer/support-ops-hub/backend/services/customer		coverage: 0.0% of statements
	github.com/macmakobaseballer/support-ops-hub/backend/services/system		coverage: 0.0% of statements
=== RUN   TestGetTicket_Found_Returns200
=== PAUSE TestGetTicket_Found_Returns200
=== RUN   TestGetTicket_NotFound_HasErrorCode
=== PAUSE TestGetTicket_NotFound_HasErrorCode
=== RUN   TestCreateTicket_Valid_Returns201
=== PAUSE TestCreateTicket_Valid_Returns201
=== RUN   TestCreateTicket_MissingTitle_Returns400
=== PAUSE TestCreateTicket_MissingTitle_Returns400
=== RUN   TestCreateTicket_InvalidType_Returns400
=== PAUSE TestCreateTicket_InvalidType_Returns400
=== RUN   TestCreateTicket_NoUserID_Returns401
=== PAUSE TestCreateTicket_NoUserID_Returns401
=== RUN   TestUpdateStatus_ValidTransitions
=== PAUSE TestUpdateStatus_ValidTransitions
=== RUN   TestUpdateStatus_InvalidTransitions_Returns422
=== PAUSE TestUpdateStatus_InvalidTransitions_Returns422
=== RUN   TestUpdateStatus_InvalidStatusValue_Returns400
=== PAUSE TestUpdateStatus_InvalidStatusValue_Returns400
=== RUN   TestUpdateStatus_NoUserID_Returns401
=== PAUSE TestUpdateStatus_NoUserID_Returns401
=== RUN   TestUpdateStatus_NotFound_Returns404
=== PAUSE TestUpdateStatus_NotFound_Returns404
=== RUN   TestUpdateAssignee_NullAssignee_Returns200
=== PAUSE TestUpdateAssignee_NullAssignee_Returns200
=== CONT  TestGetTicket_Found_Returns200
=== CONT  TestUpdateStatus_ValidTransitions
--- PASS: TestGetTicket_Found_Returns200 (0.00s)
=== CONT  TestCreateTicket_NoUserID_Returns401
=== RUN   TestUpdateStatus_ValidTransitions/new->in_progress
--- PASS: TestCreateTicket_NoUserID_Returns401 (0.00s)
=== CONT  TestCreateTicket_InvalidType_Returns400
--- PASS: TestCreateTicket_InvalidType_Returns400 (0.00s)
=== CONT  TestCreateTicket_MissingTitle_Returns400
--- PASS: TestCreateTicket_MissingTitle_Returns400 (0.00s)
=== CONT  TestCreateTicket_Valid_Returns201
=== PAUSE TestUpdateStatus_ValidTransitions/new->in_progress
--- PASS: TestCreateTicket_Valid_Returns201 (0.00s)
=== CONT  TestGetTicket_NotFound_HasErrorCode
=== RUN   TestUpdateStatus_ValidTransitions/in_progress->waiting
=== PAUSE TestUpdateStatus_ValidTransitions/in_progress->waiting
--- PASS: TestGetTicket_NotFound_HasErrorCode (0.00s)
=== RUN   TestUpdateStatus_ValidTransitions/in_progress->done
=== CONT  TestUpdateStatus_NoUserID_Returns401
=== PAUSE TestUpdateStatus_ValidTransitions/in_progress->done
=== RUN   TestUpdateStatus_ValidTransitions/waiting->in_progress
=== PAUSE TestUpdateStatus_ValidTransitions/waiting->in_progress
--- PASS: TestUpdateStatus_NoUserID_Returns401 (0.00s)
=== RUN   TestUpdateStatus_ValidTransitions/waiting->done
=== CONT  TestUpdateAssignee_NullAssignee_Returns200
--- PASS: TestUpdateAssignee_NullAssignee_Returns200 (0.00s)
=== CONT  TestUpdateStatus_NotFound_Returns404
=== PAUSE TestUpdateStatus_ValidTransitions/waiting->done
=== CONT  TestUpdateStatus_ValidTransitions/new->in_progress
--- PASS: TestUpdateStatus_NotFound_Returns404 (0.00s)
=== CONT  TestUpdateStatus_ValidTransitions/waiting->done
=== CONT  TestUpdateStatus_InvalidStatusValue_Returns400
--- PASS: TestUpdateStatus_InvalidStatusValue_Returns400 (0.00s)
=== CONT  TestUpdateStatus_InvalidTransitions_Returns422
=== RUN   TestUpdateStatus_InvalidTransitions_Returns422/new->waiting
=== PAUSE TestUpdateStatus_InvalidTransitions_Returns422/new->waiting
=== RUN   TestUpdateStatus_InvalidTransitions_Returns422/new->done
=== PAUSE TestUpdateStatus_InvalidTransitions_Returns422/new->done
=== RUN   TestUpdateStatus_InvalidTransitions_Returns422/done->new
=== PAUSE TestUpdateStatus_InvalidTransitions_Returns422/done->new
=== RUN   TestUpdateStatus_InvalidTransitions_Returns422/done->in_progress
=== PAUSE TestUpdateStatus_InvalidTransitions_Returns422/done->in_progress
=== RUN   TestUpdateStatus_InvalidTransitions_Returns422/done->waiting
=== PAUSE TestUpdateStatus_InvalidTransitions_Returns422/done->waiting
=== RUN   TestUpdateStatus_InvalidTransitions_Returns422/in_progress->new
=== PAUSE TestUpdateStatus_InvalidTransitions_Returns422/in_progress->new
=== RUN   TestUpdateStatus_InvalidTransitions_Returns422/waiting->new
=== PAUSE TestUpdateStatus_InvalidTransitions_Returns422/waiting->new
=== CONT  TestUpdateStatus_InvalidTransitions_Returns422/new->waiting
=== CONT  TestUpdateStatus_InvalidTransitions_Returns422/waiting->new
=== CONT  TestUpdateStatus_ValidTransitions/waiting->in_progress
=== CONT  TestUpdateStatus_InvalidTransitions_Returns422/in_progress->new
=== CONT  TestUpdateStatus_InvalidTransitions_Returns422/done->waiting
=== CONT  TestUpdateStatus_InvalidTransitions_Returns422/done->in_progress
=== CONT  TestUpdateStatus_ValidTransitions/in_progress->done
=== CONT  TestUpdateStatus_InvalidTransitions_Returns422/done->new
=== CONT  TestUpdateStatus_InvalidTransitions_Returns422/new->done
--- PASS: TestUpdateStatus_InvalidTransitions_Returns422 (0.00s)
    --- PASS: TestUpdateStatus_InvalidTransitions_Returns422/new->waiting (0.00s)
    --- PASS: TestUpdateStatus_InvalidTransitions_Returns422/waiting->new (0.00s)
    --- PASS: TestUpdateStatus_InvalidTransitions_Returns422/in_progress->new (0.00s)
    --- PASS: TestUpdateStatus_InvalidTransitions_Returns422/done->waiting (0.00s)
    --- PASS: TestUpdateStatus_InvalidTransitions_Returns422/done->in_progress (0.00s)
    --- PASS: TestUpdateStatus_InvalidTransitions_Returns422/done->new (0.00s)
    --- PASS: TestUpdateStatus_InvalidTransitions_Returns422/new->done (0.00s)
=== CONT  TestUpdateStatus_ValidTransitions/in_progress->waiting
--- PASS: TestUpdateStatus_ValidTransitions (0.00s)
    --- PASS: TestUpdateStatus_ValidTransitions/waiting->done (0.00s)
    --- PASS: TestUpdateStatus_ValidTransitions/new->in_progress (0.00s)
    --- PASS: TestUpdateStatus_ValidTransitions/waiting->in_progress (0.00s)
    --- PASS: TestUpdateStatus_ValidTransitions/in_progress->waiting (0.00s)
    --- PASS: TestUpdateStatus_ValidTransitions/in_progress->done (0.00s)
PASS
coverage: 35.3% of statements
ok  	github.com/macmakobaseballer/support-ops-hub/backend/services/ticket	0.022s	coverage: 35.3% of statements
