import type {
  ApiEvent,
  AuthCredentials,
  AuthResponse,
  AuthUser,
  EventTicketResponse,
  RegisterPayload,
} from '~/types/auth'

export function useAuthSession() {
  const accessToken = useState<string | null>('auth:accessToken', () => null)
  const user = useState<AuthUser | null>('auth:user', () => null)
  const loading = useState<boolean>('auth:loading', () => false)
  const lastEvent = useState<ApiEvent | null>('auth:lastEvent', () => null)
  const eventSource = useState<EventSource | null>('auth:eventSource', () => null)

  const isAuthenticated = computed(() => Boolean(accessToken.value && user.value))
  const isManager = computed(() => user.value?.globalRole === 'manager')

  const setSession = (response: AuthResponse) => {
    accessToken.value = response.accessToken
    user.value = {
      userId: response.userId,
      email: response.email,
      username: response.username,
      globalRole: response.globalRole,
    }
  }

  const clearSession = () => {
    accessToken.value = null
    user.value = null
    lastEvent.value = null
    if (eventSource.value) {
      eventSource.value.close()
      eventSource.value = null
    }
  }

  const login = async (payload: AuthCredentials) => {
    const response = await $fetch<AuthResponse>('/api/auth/login', {
      method: 'POST',
      body: payload,
      credentials: 'include',
    })
    setSession(response)
    await connectEvents()
  }

  const register = async (payload: RegisterPayload) => {
    const response = await $fetch<AuthResponse>('/api/auth/register', {
      method: 'POST',
      body: payload,
      credentials: 'include',
    })
    setSession(response)
    await connectEvents()
  }

  const refresh = async () => {
    try {
      const response = await $fetch<AuthResponse>('/api/auth/refresh', {
        method: 'POST',
        credentials: 'include',
      })
      setSession(response)
      return true
    } catch {
      clearSession()
      return false
    }
  }

  const loadMe = async () => {
    if (!accessToken.value) {
      const refreshed = await refresh()
      if (!refreshed) {
        return false
      }
    }

    try {
      user.value = await $fetch<AuthUser>('/api/auth/me', {
        headers: {
          Authorization: `Bearer ${accessToken.value}`,
        },
        credentials: 'include',
      })
      await connectEvents()
      return true
    } catch {
      clearSession()
      return false
    }
  }

  const logout = async () => {
    try {
      await $fetch('/api/auth/logout', {
        method: 'POST',
        headers: accessToken.value
          ? {
              Authorization: `Bearer ${accessToken.value}`,
            }
          : undefined,
        credentials: 'include',
      })
    } finally {
      clearSession()
      await navigateTo('/login')
    }
  }

  const connectEvents = async () => {
    if (!import.meta.client || !accessToken.value) {
      return
    }
    if (eventSource.value) {
      return
    }

    try {
      const response = await $fetch<EventTicketResponse>('/api/events/ticket', {
        method: 'POST',
        headers: {
          Authorization: `Bearer ${accessToken.value}`,
        },
        credentials: 'include',
      })

      const source = new EventSource(`/api/events?ticket=${encodeURIComponent(response.ticket)}`)
      source.addEventListener('connected', () => {
        lastEvent.value = { type: 'connected' }
      })
      source.addEventListener('membership.created', (event) => {
        lastEvent.value = parseServerEvent('membership.created', event)
      })
      source.addEventListener('user.registered', (event) => {
        lastEvent.value = parseServerEvent('user.registered', event)
      })
      source.onerror = () => {
        source.close()
        if (eventSource.value === source) {
          eventSource.value = null
        }
      }
      eventSource.value = source
    } catch {
      eventSource.value = null
    }
  }

  return {
    accessToken,
    user,
    loading,
    lastEvent,
    isAuthenticated,
    isManager,
    login,
    register,
    refresh,
    loadMe,
    logout,
    clearSession,
    connectEvents,
  }
}

function parseServerEvent(type: string, event: Event): ApiEvent {
  const message = event as MessageEvent<string>
  if (!message.data) {
    return { type }
  }

  try {
    return {
      type,
      data: JSON.parse(message.data),
    }
  } catch {
    return {
      type,
      data: message.data,
    }
  }
}
