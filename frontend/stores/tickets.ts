import { defineStore } from 'pinia'
import type { components } from '~/types/api'

type Ticket = components['schemas']['Ticket']

export const useTicketsStore = defineStore('tickets', () => {
  const tickets = ref<Ticket[]>([])

  return { tickets }
})
