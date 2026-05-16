import { defineStore } from 'pinia'
import type { components } from '~/types/api'

type User = components['schemas']['User']

export const useAuthStore = defineStore('auth', () => {
  // Persisted to localStorage so auth survives page reloads.
  // useLocalStorage is SSR-safe: returns null on the server, syncs on the client.
  const user = useLocalStorage<User | null>('auth:user', null)
  const token = useLocalStorage<string | null>('auth:token', null)

  // devEmail is stored separately so the initial fetchMe request (on page reload)
  // can inject X-Dev-User-Email before the full user object is restored.
  const devEmail = useLocalStorage<string | null>('auth:dev_email', null)

  const isLoggedIn = computed(() => user.value !== null)
  const isAdmin = computed(() => user.value?.role === 'admin')
  const currentUserId = computed(() => user.value?.id ?? null)

  function setAuth(newToken: string, newUser: User) {
    token.value = newToken
    user.value = newUser
    devEmail.value = newUser.email ?? null
  }

  function clearAuth() {
    token.value = null
    user.value = null
    devEmail.value = null
  }

  return {
    user,
    token,
    devEmail,
    isLoggedIn,
    isAdmin,
    currentUserId,
    setAuth,
    clearAuth,
  }
})
