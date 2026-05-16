// Runs on the client only (.client.ts). Restores auth state from localStorage on page reload.
export default defineNuxtPlugin(async () => {
  const authStore = useAuthStore()

  // devEmail is persisted in localStorage. If present, the user was logged in before the reload.
  // We call fetchMe() which injects X-Dev-User-Email via getDevUserEmail in plugins/api.ts.
  if (authStore.devEmail && !authStore.user) {
    const { fetchMe } = useAuth()
    try {
      await fetchMe()
    }
    catch {
      // Token / session is no longer valid — reset to avoid a stuck state
      authStore.clearAuth()
    }
  }
})
