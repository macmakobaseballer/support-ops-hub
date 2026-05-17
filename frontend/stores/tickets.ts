import { defineStore } from 'pinia'
import type { components } from '~/types/api'

type Ticket = components['schemas']['Ticket']
type TicketCreate = components['schemas']['TicketCreate']
type TicketUpdate = components['schemas']['TicketUpdate']
type TicketStatus = components['schemas']['TicketStatus']
type Pagination = components['schemas']['Pagination']

export interface TicketFilters {
  status?: TicketStatus | ''
  priority?: string
  type?: string
  customer_id?: number | ''
  system_id?: number | ''
  keyword?: string
  page?: number
  per_page?: number
}

export const useTicketsStore = defineStore('tickets', () => {
  const list = ref<Ticket[]>([])
  const current = ref<Ticket | null>(null)
  const pagination = ref<Pagination | null>(null)
  const loading = ref(false)

  async function fetchList(filters: TicketFilters = {}) {
    const api = useNuxtApp().$api
    loading.value = true
    try {
      const params = new URLSearchParams()
      if (filters.status) params.set('status', filters.status)
      if (filters.priority) params.set('priority', filters.priority)
      if (filters.type) params.set('type', filters.type)
      if (filters.customer_id) params.set('customer_id', String(filters.customer_id))
      if (filters.system_id) params.set('system_id', String(filters.system_id))
      if (filters.keyword) params.set('keyword', filters.keyword)
      if (filters.page) params.set('page', String(filters.page))
      if (filters.per_page) params.set('per_page', String(filters.per_page))

      const query = params.toString()
      const resp = await api.request<{ data: Ticket[], pagination: Pagination }>(
        `/tickets${query ? `?${query}` : ''}`,
      )
      list.value = resp.data
      pagination.value = resp.pagination
    }
    finally {
      loading.value = false
    }
  }

  async function fetchById(id: number) {
    const api = useNuxtApp().$api
    current.value = await api.request<Ticket>(`/tickets/${id}`)
  }

  async function create(payload: TicketCreate): Promise<Ticket> {
    const api = useNuxtApp().$api
    const ticket = await api.request<Ticket>('/tickets', {
      method: 'POST',
      body: JSON.stringify(payload),
    })
    return ticket
  }

  async function update(id: number, payload: TicketUpdate): Promise<Ticket> {
    const api = useNuxtApp().$api
    const ticket = await api.request<Ticket>(`/tickets/${id}`, {
      method: 'PUT',
      body: JSON.stringify(payload),
    })
    if (current.value?.id === id) {
      current.value = ticket
    }
    return ticket
  }

  async function updateStatus(id: number, status: TicketStatus): Promise<Ticket> {
    const api = useNuxtApp().$api
    const ticket = await api.request<Ticket>(`/tickets/${id}/status`, {
      method: 'PATCH',
      body: JSON.stringify({ status }),
    })
    if (current.value?.id === id) {
      current.value = ticket
    }
    return ticket
  }

  async function updateAssignee(id: number, assigneeId: number | null): Promise<Ticket> {
    const api = useNuxtApp().$api
    const ticket = await api.request<Ticket>(`/tickets/${id}/assignee`, {
      method: 'PATCH',
      body: JSON.stringify({ assignee_id: assigneeId }),
    })
    if (current.value?.id === id) {
      current.value = ticket
    }
    return ticket
  }

  return { list, current, pagination, loading, fetchList, fetchById, create, update, updateStatus, updateAssignee }
})
