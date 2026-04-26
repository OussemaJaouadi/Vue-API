export function useApiClient() {
  const { apiBaseUrl } = useApiConfig()

  const request = async <T>(path: string, options: any = {}): Promise<T> => {
    try {
      return await $fetch<T>(path, {
        baseURL: apiBaseUrl,
        ...options,
      })
    } catch (err: any) {
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
