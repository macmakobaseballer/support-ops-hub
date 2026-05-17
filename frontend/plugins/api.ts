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
    getDevUserEmail: () => {
      if (!import.meta.dev) return null
      // On page reload user may be null while devEmail is already restored from localStorage
      return (authStore.user?.email ?? authStore.devEmail) as string | null
    },
    onUnauthorized: () => {
      authStore.clearAuth()
      navigateTo('/login')
    },
    onForbidden: () => {
      setError('権限がありません')
    },
    onValidationError: (err) => {
      // 呼び出し元が catch していない場合のフォールバック通知
      setError(err.message || '入力内容に誤りがあります')
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
