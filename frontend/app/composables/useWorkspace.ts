import type { Workspace } from '~/types/workspace'

export function useWorkspace() {
  const { get, post, put, delete: del } = useApiClient()
  const auth = useAuthSession()

  const workspaces = useState<Workspace[]>('workspaces', () => [])
  const currentWorkspaceId = useState<string>('workspace:id', () => '')
  const workspacesLoading = useState<boolean>('workspaces:loading', () => false)
  const workspacesError = useState<string | null>('workspaces:error', () => null)

  const currentWorkspace = computed(() =>
    workspaces.value.find(w => w.id === currentWorkspaceId.value) ?? workspaces.value[0]
  )

  const loadWorkspaces = async () => {
    if (!auth.accessToken.value) {
      workspacesLoading.value = false
      return
    }

    workspacesLoading.value = true
    workspacesError.value = null
    try {
      const data = await get<any[]>('/v1/workspaces')
      workspaces.value = data
      if (data.length > 0) {
        const currentStillExists = data.some(workspace => workspace.id === currentWorkspaceId.value)
        const preferredId = import.meta.client ? localStorage.getItem('preferredWorkspaceId') : null
        const preferredWorkspace = preferredId ? data.find(workspace => workspace.id === preferredId) : null

        if (!currentStillExists) {
          currentWorkspaceId.value = preferredWorkspace?.id ?? data[0].id
        }
      } else {
        currentWorkspaceId.value = ''
      }
    } catch (err: any) {
      workspacesError.value = err?.data?.error || err?.message || 'Failed to load workspaces'
      workspaces.value = []
      currentWorkspaceId.value = ''
    } finally {
      workspacesLoading.value = false
    }
  }

  const createWorkspace = async (name: string) => {
    const ws = await post<any>('/v1/workspaces', { name })
    workspaces.value.push(ws)
    currentWorkspaceId.value = ws.id
    return ws
  }

  const renameWorkspace = async (id: string, name: string) => {
    const ws = await put<any>(`/v1/workspaces/${id}`, { name })
    const idx = workspaces.value.findIndex(w => w.id === id)
    if (idx >= 0) {
      workspaces.value[idx] = { ...workspaces.value[idx], ...ws }
    }
  }

  const deleteWorkspace = async (id: string) => {
    await del(`/v1/workspaces/${id}`)
    const nextWorkspaces = workspaces.value.filter(workspace => workspace.id !== id)
    workspaces.value = nextWorkspaces

    if (currentWorkspaceId.value === id) {
      currentWorkspaceId.value = nextWorkspaces[0]?.id ?? ''
    }

    if (import.meta.client && currentWorkspaceId.value) {
      localStorage.setItem('preferredWorkspaceId', currentWorkspaceId.value)
    } else if (import.meta.client) {
      localStorage.removeItem('preferredWorkspaceId')
    }
  }

  watch(currentWorkspaceId, (id) => {
    if (import.meta.client && id) {
      localStorage.setItem('preferredWorkspaceId', id)
    }
  })

  watch(() => auth.accessToken.value, (token) => {
    if (token && workspaces.value.length === 0 && !workspacesLoading.value) {
      loadWorkspaces()
    }
    if (!token) {
      workspaces.value = []
      currentWorkspaceId.value = ''
      workspacesLoading.value = false
      workspacesError.value = null
    }
  }, { immediate: import.meta.client })

  return {
    workspaces,
    workspacesLoading,
    workspacesError,
    currentWorkspaceId,
    currentWorkspace,
    loadWorkspaces,
    createWorkspace,
    renameWorkspace,
    deleteWorkspace,
  }
}
