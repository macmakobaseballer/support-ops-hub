import type { ApiClient } from '~/lib/api/client'
import type { components } from '~/types/api'
import { useAuthStore } from '~/stores/auth'

type User = components['schemas']['User']

export function useAuth() {
  const nuxtApp = useNuxtApp()
  const api = nuxtApp.$api as unknown as ApiClient
  const authStore = useAuthStore()

  async function login(email: string): Promise<void> {
    const resp = await api.request<{ token: string; user: User }>(
      '/auth/login',
      {
        method: 'POST',
        body: JSON.stringify({ email, password: 'dev-password' }),
      },
    )
    authStore.setAuth(resp.token, resp.user)
  }

  async function logout(): Promise<void> {
    try {
      await api.request<undefined>('/auth/logout', { method: 'POST' })
    }
    finally {
      authStore.clearAuth()
    }
  }

  async function fetchMe(): Promise<void> {
    const user = await api.request<User>('/auth/me')
    if (authStore.token) {
      authStore.setAuth(authStore.token, user)
    }
  }

  return { login, logout, fetchMe }
}
