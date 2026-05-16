import { defineStore } from 'pinia'
import type { components } from '~/types/api'

type Customer = components['schemas']['Customer']

export const useCustomersStore = defineStore('customers', () => {
  const customers = ref<Customer[]>([])

  return { customers }
})
