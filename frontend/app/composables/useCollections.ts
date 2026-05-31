import { useStorage } from '@vueuse/core'

export function useCollections() {
  const workbench = useWorkbench()
  const workspaceId = useState<string>('workspace:id', () => '')

  const requestCount = computed(() => workbench.treeItems.value.reduce((count, group) => count + group.requests.length, workbench.rootRequests.value.length))
  const activeCollectionName = useState<string>('collections:active-collection', () => 'all')
  const expandedCollections = useStorage<Record<string, Record<string, boolean>>>('collections:expanded-by-workspace', {})
  const workspaceExpandedCollections = computed({
    get: () => expandedCollections.value[workspaceId.value || 'global'] || { all: true },
    set: (value) => {
      expandedCollections.value = {
        ...expandedCollections.value,
        [workspaceId.value || 'global']: value,
      }
    },
  })

  const environmentPolicyFor = (name: string) => {
    const collection = workbench.treeItems.value.find(group => group.name === name)
    return collection?.environmentPolicy ?? {
      defaultEnvironment: '',
      allowedEnvironments: [],
      visibility: 'project' as const,
      roles: [],
    }
  }

  const activeCollection = computed(() => {
    if (activeCollectionName.value === 'all') return null
    return workbench.treeItems.value.find(group => group.name === activeCollectionName.value) ?? null
  })

  watch([workspaceId, () => workbench.treeItems.value.map(group => group.name).join('|')], () => {
    if (activeCollectionName.value === 'all') return
    if (!workbench.treeItems.value.some(group => group.name === activeCollectionName.value)) {
      activeCollectionName.value = 'all'
    }
  })

  const displayedCollections = computed(() => {
    if (activeCollection.value) return [activeCollection.value]
    return workbench.treeItems.value
  })

  const activeRequestCount = computed(() => {
    if (activeCollection.value) return activeCollection.value.requests.length
    return requestCount.value
  })

  const selectCollection = (name: string) => {
    activeCollectionName.value = name
    workspaceExpandedCollections.value = {
      ...workspaceExpandedCollections.value,
      [name]: true,
    }
  }

  const toggleCollection = (name: string) => {
    workspaceExpandedCollections.value = {
      ...workspaceExpandedCollections.value,
      [name]: !workspaceExpandedCollections.value[name],
    }
  }

  const exportCollections = async () => {
    if (!import.meta.client) return
    if (!workspaceId.value) {
      throw new Error('No workspace selected')
    }

    const { get } = useApiClient()
    const exportParams = new URLSearchParams({ workspaceId: workspaceId.value })
    if (activeCollection.value?.id) {
      exportParams.set('collectionId', activeCollection.value.id)
    }

    const exportPayload = await get<unknown>(`/v1/collections/export?${exportParams.toString()}`)
    const payload = JSON.stringify(exportPayload, null, 2)
    const blob = new Blob([payload], { type: 'application/json' })
    const url = URL.createObjectURL(blob)
    const link = document.createElement('a')
    const exportScope = activeCollection.value?.name ?? 'collections'
    const exportSlug = exportScope.toLowerCase().replace(/[^a-z0-9]+/g, '-').replace(/^-|-$/g, '') || 'collections'
    link.href = url
    link.download = `vue-api-workbench-${exportSlug}-${new Date().toISOString().slice(0, 10)}.json`
    link.click()
    URL.revokeObjectURL(url)
  }

  const importCollections = async (fileName: string, content: string) => {
    if (!workspaceId.value) {
      throw new Error('No workspace selected')
    }

    const { post } = useApiClient()
    return await post<{
      format: string
      collectionsCreated: number
      requestsCreated: number
      warnings: string[]
    }>('/v1/collections/import', {
      workspaceId: workspaceId.value,
      fileName,
      content,
    })
  }

  const previewImportContent = async (fileName: string, content: string) => {
    if (!workspaceId.value) {
      throw new Error('No workspace selected')
    }

    const { post } = useApiClient()
    return await post<{
      fileName: string
      format: string
      status: 'ready' | 'unsupported' | 'error'
      summary: string
      details: string[]
    }>('/v1/collections/import/preview', {
      workspaceId: workspaceId.value,
      fileName,
      content,
    })
  }

  return {
    requestCount,
    activeCollectionName,
    expandedCollections: workspaceExpandedCollections,
    activeCollection,
    displayedCollections,
    activeRequestCount,
    environmentPolicyFor,
    selectCollection,
    toggleCollection,
    exportCollections,
    importCollections,
    previewImportContent,
  }
}
