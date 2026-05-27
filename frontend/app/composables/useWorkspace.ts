import type { Workspace } from '~/types/workspace'

export function useWorkspace() {
  const { get, post, put } = useApiClient()

  const workspaces = useState<Workspace[]>('workspaces', () => [])
  const currentWorkspaceId = useState<string>('workspace:id', () => '')
  const workspacesLoading = useState<boolean>('workspaces:loading', () => false)

  const currentWorkspace = computed(() =>
    workspaces.value.find(w => w.id === currentWorkspaceId.value) ?? workspaces.value[0]
  )

  const loadWorkspaces = async () => {
    workspacesLoading.value = true
    try {
      const data = await get<any[]>('/v1/workspaces')
      workspaces.value = data
      if (data.length > 0 && !currentWorkspaceId.value) {
        currentWorkspaceId.value = data[0].id
      }
    } catch (err) {
      console.error('Failed to load workspaces', err)
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

  watch(currentWorkspaceId, (id) => {
    if (id) {
      localStorage.setItem('preferredWorkspaceId', id)
    }
  })

  if (workspaces.value.length === 0 && !workspacesLoading.value) {
    loadWorkspaces()
  }

  return {
    workspaces,
    workspacesLoading,
    currentWorkspaceId,
    currentWorkspace,
    loadWorkspaces,
    createWorkspace,
    renameWorkspace,
  }
}
