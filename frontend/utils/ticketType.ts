import type { components } from '~/types/api'

type TicketType = components['schemas']['TicketType']

export const TICKET_TYPE_LABELS: Record<TicketType, string> = {
  question: '操作方法の質問',
  bug: 'バグ報告',
  config: '設定変更依頼',
  data: 'データ修正依頼',
}

// Default priority for each type (frontend-only, backend stores what it receives)
export const TICKET_TYPE_DEFAULT_PRIORITY: Record<TicketType, components['schemas']['TicketPriority']> = {
  question: 'medium',
  bug: 'high',
  config: 'medium',
  data: 'medium',
}
