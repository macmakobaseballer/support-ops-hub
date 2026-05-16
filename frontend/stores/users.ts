import { defineStore } from 'pinia'
import type { components } from '~/types/api'

type User = components['schemas']['User']

export const useUsersStore = defineStore('users', () => {
  const users = ref<User[]>([])

  return { users }
})
