import { defineStore } from 'pinia'
import type { components } from '~/types/api'

type DashboardSummary = components['schemas']['DashboardSummary']
type CustomerReport = components['schemas']['CustomerReport']

export const useAnalyticsStore = defineStore('analytics', () => {
  const dashboard = ref<DashboardSummary | null>(null)
  const customerReport = ref<CustomerReport | null>(null)
  const loading = ref(false)

  async function fetchDashboard() {
    const api = useNuxtApp().$api
    loading.value = true
    try {
      dashboard.value = await api.request<DashboardSummary>('/analytics/dashboard')
    }
    finally {
      loading.value = false
    }
  }

  async function fetchCustomerReport(customerId: number) {
    const api = useNuxtApp().$api
    loading.value = true
    try {
      customerReport.value = await api.request<CustomerReport>(`/analytics/customers/${customerId}/report`)
    }
    finally {
      loading.value = false
    }
  }

  return { dashboard, customerReport, loading, fetchDashboard, fetchCustomerReport }
})
