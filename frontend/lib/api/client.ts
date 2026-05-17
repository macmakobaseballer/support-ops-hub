export type ApiError = {
  code: string
  message: string
  details?: Record<string, string>
}

export type ApiClientOptions = {
  baseURL: string
  getToken?: () => string | null
  /** Returns the dev user email to inject as X-Dev-User-Email header (dev mode only). */
  getDevUserEmail?: () => string | null
  onUnauthorized?: () => void
  onForbidden?: (err: ApiError) => void
  /** 422 Unprocessable Entity: validation errors or business rule violations. */
  onValidationError?: (err: ApiError) => void
  onServerError?: (err: ApiError) => void
}

export class ApiClient {
  constructor(private readonly options: ApiClientOptions) {}

  async request<T>(path: string, init: RequestInit = {}): Promise<T> {
    const headers = new Headers(init.headers)
    if (!headers.has('Content-Type')) {
      headers.set('Content-Type', 'application/json')
    }

    const token = this.options.getToken?.()
    if (token) {
      headers.set('Authorization', `Bearer ${token}`)
    }

    const devEmail = this.options.getDevUserEmail?.()
    if (devEmail) {
      headers.set('X-Dev-User-Email', devEmail)
    }

    const resp = await fetch(`${this.options.baseURL}${path}`, {
      ...init,
      headers,
    })

    if (!resp.ok) {
      const err = await this.parseError(resp)
      this.dispatchError(resp.status, err)
      throw err
    }

    if (resp.status === 204 || resp.headers.get('content-length') === '0') {
      return undefined as T
    }

    return resp.json() as Promise<T>
  }

  private async parseError(resp: Response): Promise<ApiError> {
    try {
      return await resp.json() as ApiError
    }
    catch {
      return {
        code: 'UNKNOWN_ERROR',
        message: 'エラーが発生しました',
      }
    }
  }

  private dispatchError(status: number, err: ApiError): void {
    if (status === 401) {
      this.options.onUnauthorized?.()
      return
    }
    if (status === 403) {
      this.options.onForbidden?.(err)
      return
    }
    if (status === 422) {
      this.options.onValidationError?.(err)
      return
    }
    if (status >= 500) {
      this.options.onServerError?.(err)
    }
  }
}
