export function useEnvironments() {
  const { get, post, put, delete: del } = useApiClient()

  const environments = useState<any[]>('environments:data', () => [])
  const activeEnvironmentName = useState<string>('environments:active', () => '')
  const envsLoading = useState<boolean>('environments:loading', () => true)
  const envsError = useState<string | null>('environments:error', () => null)
  const envsWorkspaceId = useState<string>('environments:workspace-id', () => '')

  const loadEnvironments = async () => {
    const wid = workspaceId.value
    if (!wid) {
      if (!workspacesLoading.value) envsLoading.value = false
      return
    }
    envsLoading.value = true
    envsError.value = null
    try {
      const data = await get<any[]>(`/v1/environments?workspaceId=${wid}`)
      environments.value = data
      envsWorkspaceId.value = wid
      if (data.length > 0 && !activeEnvironmentName.value) {
        activeEnvironmentName.value = data[0].name
      }
    } catch (err: any) {
      environments.value = []
      activeEnvironmentName.value = ''
      envsWorkspaceId.value = ''
      envsError.value = err?.data?.error || err?.message || 'Failed to load environments'
    } finally {
      envsLoading.value = false
    }
  }

  const workspaceId = useState<string>('workspace:id', () => '')
  const workspacesLoading = useState<boolean>('workspaces:loading', () => false)

  const settleEnvironmentState = () => {
    if (workspacesLoading.value) {
      envsLoading.value = true
      envsError.value = null
      return
    }

    if (!workspaceId.value) {
      environments.value = []
      activeEnvironmentName.value = ''
      envsWorkspaceId.value = ''
      envsLoading.value = false
      return
    }

    if (envsWorkspaceId.value !== workspaceId.value) {
      loadEnvironments()
      return
    }

    envsLoading.value = false
  }

  watch([workspaceId, workspacesLoading], settleEnvironmentState, { immediate: true })

  const activeEnvironment = computed(() => {
    if (!activeEnvironmentName.value) return environments.value[0]
    return environments.value.find(env => env.name === activeEnvironmentName.value) ?? environments.value[0]
  })

  const secretVariableCount = computed(() => {
    if (!activeEnvironment.value) return 0
    return activeEnvironment.value.variables.filter((v: any) => v.secret).length
  })

  const handleCreate = async (env: any) => {
    const wid = workspaceId.value
    if (!wid) return
    try {
      const created = await post<any>('/v1/environments', {
        workspaceId: wid,
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
    const wid = workspaceId.value
    const env = activeEnvironment.value
    if (!env || !wid) return

    try {
      await del(`/v1/environments/${env.id}?workspaceId=${encodeURIComponent(wid)}`)
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
    const wid = workspaceId.value
    const env = activeEnvironment.value
    if (!env || !wid) return

    try {
      const v = await post<any>(`/v1/environments/${env.id}/variables`, {
        workspaceId: wid,
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
    const wid = workspaceId.value
    const env = activeEnvironment.value
    if (!env || !wid) return

    try {
      const updated = await put<any>(`/v1/environments/${env.id}/variables/${variableId}`, {
        workspaceId: wid,
        ...updates,
      })
      const index = env.variables.findIndex((v: any) => v.id === variableId)
      if (index !== -1) {
        env.variables[index] = updated
      }
    } catch (err: any) {
      console.error('Failed to update variable', err)
    }
  }

  const deleteVariable = async (variableId: string) => {
    const wid = workspaceId.value
    const env = activeEnvironment.value
    if (!env || !wid) return

    try {
      await del(`/v1/environments/${env.id}/variables/${variableId}?workspaceId=${encodeURIComponent(wid)}`)
      env.variables = env.variables.filter((v: any) => v.id !== variableId)
    } catch (err: any) {
      console.error('Failed to delete variable', err)
    }
  }

  return {
    environments,
    envsLoading,
    envsError,
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
