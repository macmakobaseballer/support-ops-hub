export function useApiError() {
  const error = useState<string | null>('global-error', () => null)

  function setError(msg: string) {
    error.value = msg
  }

  function clearError() {
    error.value = null
  }

  return { error: readonly(error), setError, clearError }
}
