export function useApiConfig() {
  const config = useRuntimeConfig()

  return {
    appName: config.public.appName,
    apiBaseUrl: '/api',
  }
}
