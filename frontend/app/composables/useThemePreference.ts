export type ThemePreference = 'light' | 'dark' | 'system'

const themePreferenceStorageKey = 'app:theme-preference'

export function useThemePreference() {
  const preference = useState<ThemePreference>('app:theme-preference', () => 'system')
  const systemDark = useState<boolean>('app:theme-system-dark', () => false)
  let mediaQuery: MediaQueryList | null = null

  const resolvedTheme = computed<'light' | 'dark'>(() => {
    if (preference.value === 'system') {
      return systemDark.value ? 'dark' : 'light'
    }

    return preference.value
  })

  const applyTheme = () => {
    if (!import.meta.client) return

    const isDark = resolvedTheme.value === 'dark'
    document.documentElement.classList.toggle('dark', isDark)
    document.documentElement.style.colorScheme = resolvedTheme.value
  }

  const syncSystemTheme = () => {
    systemDark.value = mediaQuery?.matches ?? false
  }

  const setTheme = (nextPreference: ThemePreference) => {
    preference.value = nextPreference
  }

  watch(
    preference,
    (nextPreference) => {
      if (!import.meta.client) return
      window.localStorage.setItem(themePreferenceStorageKey, nextPreference)
      applyTheme()
    },
  )

  watch(resolvedTheme, applyTheme, { immediate: true })

  onMounted(() => {
    const storedPreference = window.localStorage.getItem(themePreferenceStorageKey) as ThemePreference | null
    if (storedPreference === 'light' || storedPreference === 'dark' || storedPreference === 'system') {
      preference.value = storedPreference
    }

    mediaQuery = window.matchMedia('(prefers-color-scheme: dark)')
    syncSystemTheme()
    mediaQuery.addEventListener('change', syncSystemTheme)
    applyTheme()
  })

  onBeforeUnmount(() => {
    mediaQuery?.removeEventListener('change', syncSystemTheme)
  })

  return {
    preference,
    resolvedTheme,
    setTheme,
  }
}
