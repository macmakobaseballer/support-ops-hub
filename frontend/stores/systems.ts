import { defineStore } from 'pinia'
import type { components } from '~/types/api'

type System = components['schemas']['System']
type UserSummary = components['schemas']['UserSummary']

export const useSystemsStore = defineStore('systems', () => {
  const systems = ref<System[]>([])
  const current = ref<System | null>(null)
  const assignees = ref<UserSummary[]>([])
  const loading = ref(false)

  async function fetchList(params: { customerId?: number, isActive?: boolean } = {}) {
    const api = useNuxtApp().$api
    loading.value = true
    try {
      const q = new URLSearchParams()
      if (params.customerId) q.set('customer_id', String(params.customerId))
      if (params.isActive !== undefined) q.set('is_active', String(params.isActive))
      const query = q.toString()
      const resp = await api.request<{ data: System[] }>(`/systems${query ? `?${query}` : ''}`)
      systems.value = resp.data
    }
    finally {
      loading.value = false
    }
  }

  async function fetchById(id: number) {
    const api = useNuxtApp().$api
    current.value = await api.request<System>(`/systems/${id}`)
  }

  async function fetchAssignees(systemId: number): Promise<UserSummary[]> {
    const api = useNuxtApp().$api
    const resp = await api.request<{ data: UserSummary[] }>(`/systems/${systemId}/assignees`)
    assignees.value = resp.data
    return resp.data
  }

  return { systems, current, assignees, loading, fetchList, fetchById, fetchAssignees }
})
