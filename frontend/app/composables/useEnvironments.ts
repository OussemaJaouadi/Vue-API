export function useEnvironments() {
  const { get, post, put, delete: del } = useApiClient()

  const environments = useState<any[]>('environments:data', () => [])
  const activeEnvironmentName = useState<string>('environments:active', () => '')
  const envsLoading = useState<boolean>('environments:loading', () => true)

  const loadEnvironments = async () => {
    envsLoading.value = true
    try {
      const data = await get<any[]>('/v1/environments?workspaceId=default')
      environments.value = data
      if (data.length > 0 && !activeEnvironmentName.value) {
        activeEnvironmentName.value = data[0].name
      }
    } catch (err) {
      console.error('Failed to load environments', err)
    } finally {
      envsLoading.value = false
    }
  }

  if (envsLoading.value) {
    loadEnvironments()
  }

  const activeEnvironment = computed(() => {
    if (!activeEnvironmentName.value) return environments.value[0]
    return environments.value.find(env => env.name === activeEnvironmentName.value) ?? environments.value[0]
  })

  const secretVariableCount = computed(() => {
    if (!activeEnvironment.value) return 0
    return activeEnvironment.value.variables.filter((v: any) => v.secret).length
  })

  const handleCreate = async (env: any) => {
    try {
      const created = await post<any>('/v1/environments', {
        workspaceId: 'default',
        name: env.name,
        visibility: env.visibility,
      })
      created.variables = []
      environments.value.push(created)
      activeEnvironmentName.value = created.name
    } catch (err: any) {
      console.error('Failed to create environment', err)
    }
  }

  const deleteEnvironment = async () => {
    if (environments.value.length <= 1) return
    const env = activeEnvironment.value
    if (!env) return

    try {
      await del(`/v1/environments/${env.id}`)
      const index = environments.value.findIndex((e: any) => e.id === env.id)
      environments.value = environments.value.filter((e: any) => e.id !== env.id)
      const nextEnv = environments.value[Math.max(0, index - 1)]
      if (nextEnv) {
        activeEnvironmentName.value = nextEnv.name
      }
    } catch (err: any) {
      console.error('Failed to delete environment', err)
    }
  }

  const addVariable = async () => {
    const env = activeEnvironment.value
    if (!env) return

    try {
      const v = await post<any>(`/v1/environments/${env.id}/variables`, {
        key: 'NEW_KEY',
        value: '',
        secret: false,
      })
      env.variables.push(v)
    } catch (err: any) {
      console.error('Failed to add variable', err)
    }
  }

  const updateVariable = async (variableId: string, updates: { key?: string; value?: string; secret?: boolean }) => {
    const env = activeEnvironment.value
    if (!env) return

    try {
      const updated = await put<any>(`/v1/environments/${env.id}/variables/${variableId}`, updates)
      const index = env.variables.findIndex((v: any) => v.id === variableId)
      if (index !== -1) {
        env.variables[index] = updated
      }
    } catch (err: any) {
      console.error('Failed to update variable', err)
    }
  }

  const deleteVariable = async (variableId: string) => {
    const env = activeEnvironment.value
    if (!env) return

    try {
      await del(`/v1/environments/${env.id}/variables/${variableId}`)
      env.variables = env.variables.filter((v: any) => v.id !== variableId)
    } catch (err: any) {
      console.error('Failed to delete variable', err)
    }
  }

  return {
    environments,
    envsLoading,
    activeEnvironmentName,
    activeEnvironment,
    secretVariableCount,
    handleCreate,
    deleteEnvironment,
    addVariable,
    updateVariable,
    deleteVariable,
    loadEnvironments,
  }
}
