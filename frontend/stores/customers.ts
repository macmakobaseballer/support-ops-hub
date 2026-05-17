import { defineStore } from 'pinia'
import type { components } from '~/types/api'

type Customer = components['schemas']['Customer']

export const useCustomersStore = defineStore('customers', () => {
  const customers = ref<Customer[]>([])
  const current = ref<Customer | null>(null)
  const loading = ref(false)

  async function fetchList(isActive?: boolean) {
    const api = useNuxtApp().$api
    loading.value = true
    try {
      const params = isActive !== undefined ? `?is_active=${isActive}` : ''
      const resp = await api.request<{ data: Customer[] }>(`/customers${params}`)
      customers.value = resp.data
    }
    finally {
      loading.value = false
    }
  }

  async function fetchById(id: number) {
    const api = useNuxtApp().$api
    current.value = await api.request<Customer>(`/customers/${id}`)
  }

  return { customers, current, loading, fetchList, fetchById }
})
