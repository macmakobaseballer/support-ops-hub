export type ApiError = {
  code: string
  message: string
  details?: Record<string, string>
}

export class ApiClient {
  constructor(private readonly baseURL: string) {}

  async request<T>(path: string, options: RequestInit = {}): Promise<T> {
    const resp = await fetch(`${this.baseURL}${path}`, {
      ...options,
      headers: {
        'Content-Type': 'application/json',
        ...options.headers,
      },
    })

    if (!resp.ok) {
      const err = await resp.json().catch((): ApiError => ({
        code: 'UNKNOWN_ERROR',
        message: 'エラーが発生しました',
      })) as ApiError
      // Rule F3: caller handles 401/403/5xx by inspecting err.code
      throw err
    }

    return resp.json() as Promise<T>
  }
}
