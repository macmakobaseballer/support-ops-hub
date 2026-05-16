import { defineStore } from 'pinia'
import type { components } from '~/types/api'

type System = components['schemas']['System']

export const useSystemsStore = defineStore('systems', () => {
  const systems = ref<System[]>([])

  return { systems }
})
