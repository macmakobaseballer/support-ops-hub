import { defineStore } from 'pinia'
import type { components } from '~/types/api'

type User = components['schemas']['User']

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const token = ref<string | null>(null)

  const isLoggedIn = computed(() => user.value !== null)
  const isAdmin = computed(() => user.value?.role === 'admin')
  const currentUserId = computed(() => user.value?.id ?? null)

  function setAuth(newToken: string, newUser: User) {
    token.value = newToken
    user.value = newUser
  }

  function clearAuth() {
    token.value = null
    user.value = null
  }

  return {
    user,
    token,
    isLoggedIn,
    isAdmin,
    currentUserId,
    setAuth,
    clearAuth,
  }
})
