import type { components } from '~/types/api'

type TicketPriority = components['schemas']['TicketPriority']

export const TICKET_PRIORITY_LABELS: Record<TicketPriority, string> = {
  high: '高',
  medium: '中',
  low: '低',
}

export const TICKET_PRIORITY_BADGE_CLASS: Record<TicketPriority, string> = {
  high: 'badge-priority-high',
  medium: 'badge-priority-medium',
  low: 'badge-priority-low',
}
