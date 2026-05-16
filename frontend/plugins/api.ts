import { ApiClient } from '~/lib/api/client'
import { useAuthStore } from '~/stores/auth'
import { useApiError } from '~/composables/useApiError'

export default defineNuxtPlugin(() => {
  const config = useRuntimeConfig()
  const authStore = useAuthStore()
  const { setError } = useApiError()

  const client = new ApiClient({
    baseURL: config.public.apiBase as string,
    getToken: () => authStore.token,
    onUnauthorized: () => {
      authStore.clearAuth()
      navigateTo('/login')
    },
    onForbidden: () => {
      setError('権限がありません')
    },
    onServerError: (err) => {
      setError(err.message || 'サーバーエラーが発生しました')
    },
  })

  return {
    provide: {
      api: client,
    },
  }
})
