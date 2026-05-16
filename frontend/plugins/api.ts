import { ApiClient } from '~/lib/api/client'

export default defineNuxtPlugin(() => {
  const config = useRuntimeConfig()
  const client = new ApiClient(config.public.apiBase as string)
  return {
    provide: {
      api: client,
    },
  }
})
