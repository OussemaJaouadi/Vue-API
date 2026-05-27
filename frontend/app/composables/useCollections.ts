export function useCollections() {
  const workbench = useWorkbench()
  const { environments } = useEnvironments()

  const requestCount = computed(() => workbench.treeItems.value.reduce((count, group) => count + group.requests.length, workbench.rootRequests.value.length))
  const activeCollectionName = useState<string>('collections:active-collection', () => 'all')
  const expandedCollections = useState<Record<string, boolean>>('collections:expanded', () => ({
    all: true,
  }))

  const environmentPolicyFor = (name: string) => {
    const envs = environments.value
    const defaultEnv = envs.length > 0 ? envs[0].name : '—'
    const allowed = envs.length > 0 ? envs.map((e: any) => e.name) : []
    return {
      defaultEnvironment: defaultEnv,
      allowedEnvironments: allowed,
      visibility: 'project' as const,
      roles: ['manager', 'developer'],
    }
  }

  const createWorkbenchExport = () => ({
    schema: 'vue-api-workbench.collection.v1',
    exportedAt: new Date().toISOString(),
    collections: workbench.treeItems.value.map(group => ({
      name: group.name,
      icon: group.icon,
      requests: group.requests.map(request => ({
        id: request.id,
        method: request.method,
        name: request.name,
        path: request.path,
      })),
    })),
    rootRequests: workbench.rootRequests.value.map(request => ({
      id: request.id,
      method: request.method,
      name: request.name,
      path: request.path,
    })),
  })

  const activeCollection = computed(() => {
    if (activeCollectionName.value === 'all') return null
    return workbench.treeItems.value.find(group => group.name === activeCollectionName.value) ?? null
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
    expandedCollections.value[name] = true
  }

  const toggleCollection = (name: string) => {
    expandedCollections.value[name] = !expandedCollections.value[name]
  }

  const exportCollections = () => {
    if (!import.meta.client) return

    const payload = JSON.stringify(createWorkbenchExport(), null, 2)
    const blob = new Blob([payload], { type: 'application/json' })
    const url = URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `vue-api-workbench-collections-${new Date().toISOString().slice(0, 10)}.json`
    link.click()
    URL.revokeObjectURL(url)
  }

  const detectJsonImport = (payload: any) => {
    if (payload?.schema === 'vue-api-workbench.collection.v1') {
      const collections = Array.isArray(payload.collections) ? payload.collections.length : 0
      const requests = Array.isArray(payload.collections)
        ? payload.collections.reduce((count: number, collection: any) => count + (Array.isArray(collection.requests) ? collection.requests.length : 0), 0)
        : 0

      return {
        format: 'Workbench export',
        status: 'ready' as const,
        summary: `${collections} collections / ${requests} requests`,
        details: ['Can be imported into local workbench state now.', 'Backend persistence will validate and store this later.'],
      }
    }

    if (payload?.openapi) {
      const paths = payload.paths && typeof payload.paths === 'object' ? Object.keys(payload.paths).length : 0

      return {
        format: `OpenAPI ${payload.openapi}`,
        status: 'unsupported' as const,
        summary: `${paths} paths detected`,
        details: ['Frontend detected the spec shape.', 'Backend parser will map operations into folders and requests.'],
      }
    }

    if (payload?.swagger === '2.0') {
      const paths = payload.paths && typeof payload.paths === 'object' ? Object.keys(payload.paths).length : 0

      return {
        format: 'Swagger 2.0',
        status: 'unsupported' as const,
        summary: `${paths} paths detected`,
        details: ['Frontend detected the legacy spec shape.', 'Backend parser will normalize it before import.'],
      }
    }

    if (payload?.info?._postman_id || payload?.item) {
      const items = Array.isArray(payload.item) ? payload.item.length : 0

      return {
        format: 'Postman collection',
        status: 'unsupported' as const,
        summary: `${items} top-level items detected`,
        details: ['Frontend detected a Postman collection.', 'Postman normalization is documented as a later parser adapter.'],
      }
    }

    return {
      format: 'Unknown JSON',
      status: 'error' as const,
      summary: 'No supported collection shape detected',
      details: ['Expected Workbench export, OpenAPI, Swagger, or Postman collection JSON.'],
    }
  }

  return {
    requestCount,
    activeCollectionName,
    expandedCollections,
    activeCollection,
    displayedCollections,
    activeRequestCount,
    environmentPolicyFor,
    selectCollection,
    toggleCollection,
    exportCollections,
    detectJsonImport,
  }
}
