import { describe, it, expect } from 'vitest'
import {
  TICKET_STATUS_LABELS,
  TICKET_STATUS_BADGE_CLASS,
  TICKET_STATUS_BORDER_COLOR,
  VALID_STATUS_TRANSITIONS,
  STATUS_TRANSITION_LABELS,
  STATUS_TRANSITION_LABELS_FROM_WAITING,
} from './ticketStatus'

describe('ticketStatus ユーティリティ', () => {
  describe('TICKET_STATUS_LABELS', () => {
    it('各ステータスに対応する日本語ラベルを返す', () => {
      expect(TICKET_STATUS_LABELS.new).toBe('新規受付')
      expect(TICKET_STATUS_LABELS.in_progress).toBe('対応中')
      expect(TICKET_STATUS_LABELS.waiting).toBe('確認待ち')
      expect(TICKET_STATUS_LABELS.done).toBe('完了')
    })
  })

  describe('TICKET_STATUS_BADGE_CLASS', () => {
    it('各ステータスに対応するバッジクラスを返す', () => {
      expect(TICKET_STATUS_BADGE_CLASS.new).toBe('badge-status-new')
      expect(TICKET_STATUS_BADGE_CLASS.in_progress).toBe('badge-status-in_progress')
      expect(TICKET_STATUS_BADGE_CLASS.waiting).toBe('badge-status-waiting')
      expect(TICKET_STATUS_BADGE_CLASS.done).toBe('badge-status-done')
    })
  })

  describe('VALID_STATUS_TRANSITIONS（ステータス遷移ボタン表示制御 F4）', () => {
    it('new からは in_progress のみ遷移可能', () => {
      expect(VALID_STATUS_TRANSITIONS.new).toEqual(['in_progress'])
    })

    it('in_progress からは waiting と done に遷移可能', () => {
      expect(VALID_STATUS_TRANSITIONS.in_progress).toContain('waiting')
      expect(VALID_STATUS_TRANSITIONS.in_progress).toContain('done')
      expect(VALID_STATUS_TRANSITIONS.in_progress).toHaveLength(2)
    })

    it('waiting からは in_progress と done に遷移可能', () => {
      expect(VALID_STATUS_TRANSITIONS.waiting).toContain('in_progress')
      expect(VALID_STATUS_TRANSITIONS.waiting).toContain('done')
      expect(VALID_STATUS_TRANSITIONS.waiting).toHaveLength(2)
    })

    it('done からは遷移不可（終端ステータス）', () => {
      expect(VALID_STATUS_TRANSITIONS.done).toEqual([])
    })

    it('new → waiting は不可（ボタンが表示されない）', () => {
      expect(VALID_STATUS_TRANSITIONS.new).not.toContain('waiting')
    })

    it('new → done は不可（ボタンが表示されない）', () => {
      expect(VALID_STATUS_TRANSITIONS.new).not.toContain('done')
    })
  })

  describe('STATUS_TRANSITION_LABELS', () => {
    it('in_progress への遷移は「対応開始する」', () => {
      expect(STATUS_TRANSITION_LABELS.in_progress).toBe('対応開始する')
    })

    it('waiting への遷移は「確認待ちへ」', () => {
      expect(STATUS_TRANSITION_LABELS.waiting).toBe('確認待ちへ')
    })

    it('done への遷移は「完了にする」', () => {
      expect(STATUS_TRANSITION_LABELS.done).toBe('完了にする')
    })
  })

  describe('STATUS_TRANSITION_LABELS_FROM_WAITING', () => {
    it('waiting から in_progress への遷移は「対応再開する」', () => {
      expect(STATUS_TRANSITION_LABELS_FROM_WAITING.in_progress).toBe('対応再開する')
    })

    it('waiting から done への遷移は「完了にする」', () => {
      expect(STATUS_TRANSITION_LABELS_FROM_WAITING.done).toBe('完了にする')
    })
  })
})
