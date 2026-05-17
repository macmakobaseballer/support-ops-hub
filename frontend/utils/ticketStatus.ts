import type { components } from '~/types/api'

type TicketStatus = components['schemas']['TicketStatus']

export const TICKET_STATUS_LABELS: Record<TicketStatus, string> = {
  new: '新規受付',
  in_progress: '対応中',
  waiting: '確認待ち',
  done: '完了',
}

export const TICKET_STATUS_BADGE_CLASS: Record<TicketStatus, string> = {
  new: 'badge-status-new',
  in_progress: 'badge-status-in_progress',
  waiting: 'badge-status-waiting',
  done: 'badge-status-done',
}

// Border color for stat-card (must be inline style because Tailwind JIT purges dynamic classes)
export const TICKET_STATUS_BORDER_COLOR: Record<TicketStatus, string> = {
  new: '#3b82f6',
  in_progress: '#7c4d33',
  waiting: '#8b5cf6',
  done: '#22c55e',
}

// Valid status transitions for UI button display (mirrors backend ルール5)
export const VALID_STATUS_TRANSITIONS: Record<TicketStatus, TicketStatus[]> = {
  new: ['in_progress'],
  in_progress: ['waiting', 'done'],
  waiting: ['in_progress', 'done'],
  done: [],
}

export const STATUS_TRANSITION_LABELS: Partial<Record<TicketStatus, string>> = {
  in_progress: '対応開始する',
  waiting: '確認待ちへ',
  done: '完了にする',
}

export const STATUS_TRANSITION_LABELS_FROM_WAITING: Partial<Record<TicketStatus, string>> = {
  in_progress: '対応再開する',
  done: '完了にする',
}
