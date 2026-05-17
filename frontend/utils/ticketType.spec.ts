import { describe, it, expect } from 'vitest'
import { TICKET_TYPE_LABELS, TICKET_TYPE_DEFAULT_PRIORITY } from './ticketType'

describe('ticketType ユーティリティ', () => {
  describe('TICKET_TYPE_LABELS', () => {
    it('各種別に対応する日本語ラベルを返す', () => {
      expect(TICKET_TYPE_LABELS.question).toBe('操作方法の質問')
      expect(TICKET_TYPE_LABELS.bug).toBe('バグ報告')
      expect(TICKET_TYPE_LABELS.config).toBe('設定変更依頼')
      expect(TICKET_TYPE_LABELS.data).toBe('データ修正依頼')
    })
  })

  describe('TICKET_TYPE_DEFAULT_PRIORITY（種別→優先度自動セット）', () => {
    it('種別を bug に変更すると優先度が high になる', () => {
      expect(TICKET_TYPE_DEFAULT_PRIORITY.bug).toBe('high')
    })

    it('種別を question に変更すると優先度が medium になる', () => {
      expect(TICKET_TYPE_DEFAULT_PRIORITY.question).toBe('medium')
    })

    it('種別を config に変更すると優先度が medium になる', () => {
      expect(TICKET_TYPE_DEFAULT_PRIORITY.config).toBe('medium')
    })

    it('種別を data に変更すると優先度が medium になる', () => {
      expect(TICKET_TYPE_DEFAULT_PRIORITY.data).toBe('medium')
    })
  })
})
