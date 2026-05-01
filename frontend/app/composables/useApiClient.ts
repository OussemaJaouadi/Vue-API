export function useApiClient() {
  const { apiBaseUrl } = useApiConfig()
  const auth = useAuthSession()

  const request = async <T>(path: string, options: any = {}): Promise<T> => {
    const { skipAuth, skipRefresh, ...fetchOptions } = options
    const headers = new Headers(fetchOptions.headers || {})
    if (!skipAuth && auth.accessToken.value) {
      headers.set('Authorization', `Bearer ${auth.accessToken.value}`)
    }

    try {
      const response = await $fetch<unknown>(path, {
        baseURL: apiBaseUrl,
        credentials: 'include',
        ...fetchOptions,
        headers,
      })
      return response as T
    } catch (err: any) {
      const status = err?.statusCode || err?.status || err?.response?.status
      if (!skipRefresh && status === 401 && await auth.refresh()) {
        const retryHeaders = new Headers(fetchOptions.headers || {})
        if (!skipAuth && auth.accessToken.value) {
          retryHeaders.set('Authorization', `Bearer ${auth.accessToken.value}`)
        }

        const response = await $fetch<unknown>(path, {
          baseURL: apiBaseUrl,
          credentials: 'include',
          ...fetchOptions,
          headers: retryHeaders,
        })
        return response as T
      }

      // If the backend returned a JSON error, use that message
      const backendError = err.data?.error
      if (backendError) {
        err.message = backendError
      }
      throw err
    }
  }

  return {
    request,
    get: <T>(path: string, options: any = {}) => request<T>(path, { ...options, method: 'GET' }),
    post: <T>(path: string, body?: any, options: any = {}) => request<T>(path, { ...options, method: 'POST', body }),
    put: <T>(path: string, body?: any, options: any = {}) => request<T>(path, { ...options, method: 'PUT', body }),
    patch: <T>(path: string, body?: any, options: any = {}) => request<T>(path, { ...options, method: 'PATCH', body }),
    delete: <T>(path: string, options: any = {}) => request<T>(path, { ...options, method: 'DELETE' }),
  }
}
