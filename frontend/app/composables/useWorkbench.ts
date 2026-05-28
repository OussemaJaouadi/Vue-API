import { PhGlobe, PhLightning } from '@phosphor-icons/vue'
import { useStorage } from '@vueuse/core'
import type { Component } from 'vue'

export type ApiMethod = 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE' | 'SOCKET'
export type WorkbenchIconKey = 'PhGlobe' | 'PhLightning'

export interface RequestItem {
  id: string
  method: ApiMethod
  name: string
  path: string
  queryParams?: QueryParamItem[]
  headers?: HeaderItem[]
  body?: string
  bodyLanguage?: BodyLanguage
  authConfig?: RequestAuthConfig
  active?: boolean
}

export interface TreeItem {
  id?: string
  name: string
  icon: WorkbenchIconKey
  requests: RequestItem[]
}

export interface QueryParamItem {
  id: string
  key: string
  value: string
  enabled: boolean
}

export interface HeaderItem extends QueryParamItem {}

export interface RequestTarget {
  baseUrl: string
  path: string
}

export type BodyLanguage = 'json' | 'xml' | 'html' | 'yaml' | 'text'
export type AuthMode = 'inherit' | 'bearer' | 'api-key' | 'basic' | 'oauth2' | 'oidc' | 'none'
export type ApiKeyPlacement = 'header' | 'query'
export type OAuthGrant = 'authorization-code-pkce' | 'client-credentials' | 'refresh-token'

export interface RequestAuthConfig {
  mode: AuthMode
  bearerToken: string
  apiKeyName: string
  apiKeyValue: string
  apiKeyPlacement: ApiKeyPlacement
  basicUsername: string
  basicPassword: string
  oauthGrant: OAuthGrant
  oauthAccessToken: string
  oauthClientId: string
  oauthTokenUrl: string
  oauthScopes: string
  oidcIssuerUrl: string
  oidcAudience: string
}

export type WebSocketConnectionState = 'idle' | 'connecting' | 'connected' | 'closing' | 'closed' | 'error'
export type WebSocketEventDirection = 'in' | 'out' | 'system' | 'error'

export interface WebSocketTimelineEvent {
  id: string
  direction: WebSocketEventDirection
  timestamp: string
  title: string
  payload?: string
  sizeBytes?: number
}

export interface WorkbenchResponse {
  status: number
  statusText: string
  duration: number
  size: number
  ttfb: number
  executionTarget: string
  requestId: string
  headers: HeaderItem[]
  body: string
  error?: string
}

export type RequestSaveStatus = 'idle' | 'saving' | 'saved' | 'error'

export const API_METHODS: ApiMethod[] = ['GET', 'POST', 'PUT', 'PATCH', 'DELETE', 'SOCKET']

export const METHOD_LABELS: Record<ApiMethod, string> = {
  GET: 'GET',
  POST: 'POST',
  PUT: 'PUT',
  PATCH: 'PATCH',
  DELETE: 'DELETE',
  SOCKET: 'WS',
}

export const METHOD_COLORS: Record<ApiMethod, string> = {
  GET: 'text-emerald-600 dark:text-emerald-400',
  POST: 'text-blue-600 dark:text-blue-400',
  PUT: 'text-amber-600 dark:text-amber-400',
  PATCH: 'text-purple-600 dark:text-purple-400',
  DELETE: 'text-destructive',
  SOCKET: 'text-primary',
}

export const METHOD_DOT_COLORS: Record<ApiMethod, string> = {
  GET: 'bg-emerald-600 dark:bg-emerald-400',
  POST: 'bg-blue-600 dark:bg-blue-400',
  PUT: 'bg-amber-600 dark:bg-amber-400',
  PATCH: 'bg-purple-600 dark:bg-purple-400',
  DELETE: 'bg-destructive',
  SOCKET: 'bg-primary',
}

export const METHOD_STRIP_COLORS: Record<ApiMethod, string> = {
  GET: 'bg-emerald-600 dark:bg-emerald-400',
  POST: 'bg-blue-600 dark:bg-blue-400',
  PUT: 'bg-amber-600 dark:bg-amber-400',
  PATCH: 'bg-purple-600 dark:bg-purple-400',
  DELETE: 'bg-red-600 dark:bg-red-400',
  SOCKET: 'bg-cyan-600 dark:bg-cyan-400',
}

export const METHOD_BADGE_COLORS: Record<ApiMethod, string> = {
  GET: 'border-emerald-500/40 bg-emerald-500/12 text-emerald-700 dark:text-emerald-300',
  POST: 'border-blue-500/40 bg-blue-500/12 text-blue-700 dark:text-blue-300',
  PUT: 'border-amber-500/40 bg-amber-500/12 text-amber-700 dark:text-amber-300',
  PATCH: 'border-purple-500/40 bg-purple-500/12 text-purple-700 dark:text-purple-300',
  DELETE: 'border-destructive/40 bg-destructive/12 text-destructive',
  SOCKET: 'border-cyan-500/40 bg-cyan-500/12 text-cyan-700 dark:text-cyan-300',
}

export const WORKBENCH_ICONS: Record<WorkbenchIconKey, Component> = {
  PhGlobe,
  PhLightning,
}

export const defaultRequestAuthConfig = (): RequestAuthConfig => ({
  mode: 'bearer',
  bearerToken: '{{accessToken}}',
  apiKeyName: 'x-api-key',
  apiKeyValue: '{{apiKey}}',
  apiKeyPlacement: 'header',
  basicUsername: '{{username}}',
  basicPassword: '{{password}}',
  oauthGrant: 'authorization-code-pkce',
  oauthAccessToken: '{{oauthAccessToken}}',
  oauthClientId: '{{clientId}}',
  oauthTokenUrl: '{{tokenUrl}}',
  oauthScopes: 'openid profile email',
  oidcIssuerUrl: '{{issuerUrl}}',
  oidcAudience: '{{audience}}',
})

const cloneRows = <T extends QueryParamItem>(rows?: T[]) => (rows || []).map(row => ({ ...row }))

const activeRows = <T extends QueryParamItem>(rows: T[]) => rows.filter(row => row.key || row.value)

const requestFromApi = (r: any): RequestItem => ({
  id: r.id,
  method: r.method,
  name: r.name,
  path: r.path,
  queryParams: cloneRows(r.queryParams),
  headers: cloneRows(r.headers),
  body: r.body || '',
  bodyLanguage: r.bodyLanguage || 'json',
  authConfig: {
    ...defaultRequestAuthConfig(),
    ...(r.authConfig || {}),
  },
})

const encodeBasicToken = (username: string, password: string) => {
  const raw = `${username}:${password}`
  if (import.meta.client) {
    return btoa(raw)
  }

  const buffer = (globalThis as any).Buffer
  return buffer ? buffer.from(raw).toString('base64') : raw
}

const authRowsForExecution = (config: RequestAuthConfig) => {
  const headers: HeaderItem[] = []
  const queryParams: QueryParamItem[] = []
  const addHeader = (key: string, value: string) => {
    if (!key || !value) return
    headers.push({ id: crypto.randomUUID(), key, value, enabled: true })
  }
  const addQueryParam = (key: string, value: string) => {
    if (!key || !value) return
    queryParams.push({ id: crypto.randomUUID(), key, value, enabled: true })
  }

  if (config.mode === 'bearer') {
    addHeader('Authorization', `Bearer ${config.bearerToken}`)
  } else if (config.mode === 'api-key') {
    if (config.apiKeyPlacement === 'query') {
      addQueryParam(config.apiKeyName || 'x-api-key', config.apiKeyValue)
    } else {
      addHeader(config.apiKeyName || 'x-api-key', config.apiKeyValue)
    }
  } else if (config.mode === 'basic') {
    addHeader('Authorization', `Basic ${encodeBasicToken(config.basicUsername, config.basicPassword)}`)
  } else if (config.mode === 'oauth2' || config.mode === 'oidc') {
    addHeader('Authorization', `Bearer ${config.oauthAccessToken}`)
  }

  return { headers, queryParams }
}

export function useWorkbench() {
  const treeItems = useState<TreeItem[]>('workbench:tree', () => [])
  const rootRequests = useState<RequestItem[]>('workbench:root-requests', () => [])
  const collectionsLoading = useState<boolean>('workbench:collections-loading', () => true)
  const collectionsError = useState<string | null>('workbench:collections-error', () => null)
  const collectionsWorkspaceId = useState<string>('workbench:collections-workspace-id', () => '')

  const workspaceId = useState<string>('workspace:id', () => '')
  const workspacesLoading = useState<boolean>('workspaces:loading', () => false)

  const settleCollectionsState = () => {
    if (workspacesLoading.value) {
      collectionsLoading.value = true
      collectionsError.value = null
      return
    }

    if (!workspaceId.value) {
      treeItems.value = []
      rootRequests.value = []
      collectionsWorkspaceId.value = ''
      collectionsLoading.value = false
      return
    }

    if (collectionsWorkspaceId.value !== workspaceId.value) {
      loadCollections()
      return
    }

    collectionsLoading.value = false
  }

  onMounted(settleCollectionsState)
  watch([workspaceId, workspacesLoading], settleCollectionsState)

  const loadCollections = async () => {
    const wid = workspaceId.value
    if (!wid) {
      if (!workspacesLoading.value) collectionsLoading.value = false
      return
    }
    collectionsLoading.value = true
    collectionsError.value = null
    try {
      const { get } = useApiClient()
      const data = await get<any>(`/v1/collections?workspaceId=${wid}`)
      treeItems.value = data.collections.map((c: any) => ({
        id: c.id,
        name: c.name,
        icon: c.icon || 'PhGlobe',
        requests: c.requests.map(requestFromApi),
      }))
      rootRequests.value = (data.rootRequests || []).map(requestFromApi)
      collectionsWorkspaceId.value = wid
    } catch (err: any) {
      treeItems.value = []
      rootRequests.value = []
      collectionsWorkspaceId.value = ''
      collectionsError.value = err?.data?.error || err?.message || 'Failed to load collections'
    } finally {
      collectionsLoading.value = false
    }
  }

  const allRequests = computed(() => [
    ...rootRequests.value,
    ...treeItems.value.flatMap(group => group.requests)
  ])

  const activeRequestId = useState<string>('workbench:active-request-id', () => '')
  
  // Persistent Layout State
  const sidebarWidth = useStorage<number>('workbench:sidebar-width', 260)
  const editorPaneHeight = useStorage<number>('workbench:editor-height', 50)
  const responseWidth = useStorage<number>('workbench:response-width', 400)
  const responsePosition = useStorage<'bottom' | 'right'>('workbench:response-position', 'bottom')
  
  const openTabs = useState<RequestItem[]>('workbench:open-tabs', () => [])

  const activeRequest = computed<RequestItem>(() => {
    // Ensure we always return a RequestItem to satisfy consumers.
    const found = allRequests.value.find(request => request.id === activeRequestId.value)
    if (found) return found
    if (openTabs.value[0]) return openTabs.value[0]
    if (allRequests.value[0]) return allRequests.value[0]

    // Fallback placeholder when no requests exist yet.
    return {
      id: crypto.randomUUID(),
      method: 'GET',
      name: 'New Request',
      path: '/',
    }
  })
  
  const loading = useState<boolean>('workbench:loading', () => false)
  const responseData = useState<WorkbenchResponse | null>('workbench:response', () => null)
  const requestTarget = useState<RequestTarget>('workbench:request-target', () => ({ baseUrl: '', path: '' }))
  const queryParams = useState<QueryParamItem[]>('workbench:query-params', () => [])
  const headers = useState<HeaderItem[]>('workbench:headers', () => [])
  const requestBody = useState<string>('workbench:body', () => '')
  const requestBodyLanguage = useState<BodyLanguage>('workbench:body-language', () => 'json')
  const requestAuthConfig = useState<RequestAuthConfig>('workbench:auth-config', defaultRequestAuthConfig)
  const requestSaveStatus = useState<RequestSaveStatus>('workbench:request-save-status', () => 'idle')
  const requestSaveError = useState<string>('workbench:request-save-error', () => '')
  const webSocketState = useState<WebSocketConnectionState>('workbench:ws-state', () => 'idle')
  const webSocketMessage = useState<string>('workbench:ws-message', () => `{
  "type": "ping",
  "requestId": "{{requestId}}"
}`)
  const webSocketMessageLanguage = useState<BodyLanguage>('workbench:ws-message-language', () => 'json')
  const webSocketEvents = useState<WebSocketTimelineEvent[]>('workbench:ws-events', () => [])
  const webSocketExecutionId = useState<string | null>('workbench:ws-execution-id', () => null)

  const auth = useAuthSession()
  watch(() => auth.lastEvent.value, (event) => {
    if (!event || !event.type.startsWith('ws.')) return
    
    const data = event.data as any
    if (data.executionId !== webSocketExecutionId.value) return

    if (event.type === 'ws.connected') {
      webSocketState.value = 'connected'
      loading.value = false
      appendWebSocketEvent({
        direction: 'system',
        title: 'Connected through backend proxy',
        payload: `Target: ${data.target}`,
      })
    } else if (event.type === 'ws.message.in' || event.type === 'ws.message.in.binary') {
      appendWebSocketEvent({
        direction: 'in',
        title: event.type === 'ws.message.in' ? 'Message received' : 'Binary message received',
        payload: data.payload,
        sizeBytes: data.sizeBytes,
      })
    } else if (event.type === 'ws.message.out') {
    } else if (event.type === 'ws.error') {
      appendWebSocketEvent({
        direction: 'error',
        title: 'WebSocket Error',
        payload: data.error,
      })
    } else if (event.type === 'ws.closed') {
      webSocketState.value = 'closed'
      webSocketExecutionId.value = null
      appendWebSocketEvent({
        direction: 'system',
        title: 'Socket closed',
      })
    }
  })

  const isWebSocketRequest = computed(() => activeRequest.value.method === 'SOCKET')

  const hydrateActiveRequestState = () => {
    const request = allRequests.value.find(item => item.id === activeRequestId.value)
    if (!request) return

    queryParams.value = cloneRows(request.queryParams)
    headers.value = cloneRows(request.headers)
    requestBody.value = request.body || ''
    requestBodyLanguage.value = request.bodyLanguage || 'json'
    requestAuthConfig.value = {
      ...defaultRequestAuthConfig(),
      ...(request.authConfig || {}),
    }
  }

  const syncActiveRequestState = () => {
    const request = allRequests.value.find(item => item.id === activeRequestId.value)
    if (!request) return

    request.queryParams = cloneRows(queryParams.value)
    request.headers = cloneRows(headers.value)
    request.body = requestBody.value
    request.bodyLanguage = requestBodyLanguage.value
    request.authConfig = { ...requestAuthConfig.value }
  }

  watch(activeRequestId, hydrateActiveRequestState)
  watch([queryParams, headers, requestBody, requestBodyLanguage, requestAuthConfig], syncActiveRequestState, { deep: true })
  watch([queryParams, headers, requestBody, requestBodyLanguage, requestAuthConfig], () => {
    if (requestSaveStatus.value === 'saving') return

    requestSaveStatus.value = 'idle'
    requestSaveError.value = ''
  }, { deep: true })

  const addFolder = async (name: string = 'New Collection') => {
    const wid = workspaceId.value
    if (!wid) return
    try {
      const { post } = useApiClient()
      const created = await post<any>('/v1/collections', {
        workspaceId: wid,
        name,
        icon: 'PhGlobe',
      })
      treeItems.value.push({
        id: created.id,
        name: created.name,
        icon: created.icon || 'PhGlobe',
        requests: []
      })
    } catch (err) {
      console.error('Failed to create collection', err)
    }
  }

  const addRequest = async (folderName?: string) => {
    const newReq: RequestItem = {
      id: crypto.randomUUID(),
      method: 'GET',
      name: 'New Request',
      path: '/path',
    }
    try {
      const { post } = useApiClient()
      const folderId = folderName ? treeItems.value.find(f => f.name === folderName)?.id : undefined
      if (!workspaceId.value) return
      const created = await post<any>('/v1/collections/requests', {
        workspaceId: workspaceId.value,
        collectionId: folderId || null,
        method: 'GET',
        name: 'New Request',
        path: '/path',
        queryParams: [],
        headers: [],
        body: '',
        bodyLanguage: 'json',
        authConfig: defaultRequestAuthConfig(),
      })
      Object.assign(newReq, requestFromApi(created))
    } catch (err) {
      console.error('Failed to create request', err)
    }

    if (folderName) {
      const folder = treeItems.value.find(f => f.name === folderName)
      if (folder) folder.requests.push(newReq)
    } else {
      rootRequests.value.push(newReq)
    }
    openRequest(newReq)
  }

  const moveRequest = (requestId: string, targetFolderName?: string, targetIndex?: number) => {
    let requestToMove: RequestItem | undefined
    let sourceIndex = -1
    let sourceFolderName: string | undefined
    
    const rootIndex = rootRequests.value.findIndex(r => r.id === requestId)
    if (rootIndex !== -1) {
      requestToMove = rootRequests.value.splice(rootIndex, 1)[0]
      sourceIndex = rootIndex
    } else {
      for (const folder of treeItems.value) {
        const folderIndex = folder.requests.findIndex(r => r.id === requestId)
        if (folderIndex !== -1) {
          requestToMove = folder.requests.splice(folderIndex, 1)[0]
          sourceIndex = folderIndex
          sourceFolderName = folder.name
          break
        }
      }
    }

    if (!requestToMove) return

    const adjustSameContainerIndex = (insertIndex: number) => {
      const sameContainer = sourceFolderName === targetFolderName
      return sameContainer && sourceIndex !== -1 && sourceIndex < insertIndex ? insertIndex - 1 : insertIndex
    }

    if (targetFolderName) {
      const folder = treeItems.value.find(f => f.name === targetFolderName)
      if (folder) {
        const insertIndex = targetIndex !== undefined ? adjustSameContainerIndex(targetIndex) : folder.requests.length
        folder.requests.splice(insertIndex, 0, requestToMove)
      }
    } else {
      const insertIndex = targetIndex !== undefined ? adjustSameContainerIndex(targetIndex) : rootRequests.value.length
      rootRequests.value.splice(insertIndex, 0, requestToMove)
    }
  }

  const deleteItem = async (identifier: string, isFolder: boolean) => {
    const wid = workspaceId.value
    if (!wid) return

    try {
      const { delete: del } = useApiClient()
      if (isFolder) {
        const folder = treeItems.value.find(f => f.id === identifier || f.name === identifier)
        if (folder?.id) {
          await del(`/v1/collections/${folder.id}?workspaceId=${encodeURIComponent(wid)}`)
        }
        if (folder) {
          folder.requests.forEach(r => closeTab(r.id))
        }
        treeItems.value = treeItems.value.filter(f => f.id !== identifier && f.name !== identifier)
      } else {
        await del(`/v1/collections/requests/${identifier}?workspaceId=${encodeURIComponent(wid)}`)
        closeTab(identifier)
        rootRequests.value = rootRequests.value.filter(r => r.id !== identifier)
        treeItems.value.forEach(f => {
          f.requests = f.requests.filter(r => r.id !== identifier)
        })
      }
    } catch (err) {
      console.error('Failed to delete', err)
    }
  }

  const closeTab = (id: string) => {
    const tabIndex = openTabs.value.findIndex(item => item.id === id)
    if (tabIndex === -1) return

    openTabs.value = openTabs.value.filter(item => item.id !== id)
    
    if (activeRequestId.value === id) {
      const candidate = openTabs.value[Math.max(0, tabIndex - 1)] || openTabs.value[0]
      activeRequestId.value = candidate ? candidate.id : ''
    }
  }

  const openRequest = (request: RequestItem) => {
    if (!openTabs.value.some(tab => tab.id === request.id)) {
      openTabs.value.push(request)
    }
    activeRequestId.value = request.id
  }

  const setActiveRequest = (request: RequestItem) => {
    activeRequestId.value = request.id
  }

  const setActiveRequestMethod = (method: ApiMethod) => {
    if (activeRequest && activeRequest.value) {
      activeRequest.value.method = method
    }
  }

  const saveActiveRequestState = async () => {
    const wid = workspaceId.value
    const request = allRequests.value.find(item => item.id === activeRequestId.value)
    if (!wid || !request) return

    syncActiveRequestState()

    const { put } = useApiClient()
    requestSaveStatus.value = 'saving'
    requestSaveError.value = ''

    try {
      const updated = await put<any>(`/v1/collections/requests/${request.id}`, {
        workspaceId: wid,
        method: request.method,
        name: request.name,
        path: request.path,
        queryParams: activeRows(queryParams.value),
        headers: activeRows(headers.value),
        body: requestBody.value,
        bodyLanguage: requestBodyLanguage.value,
        authConfig: requestAuthConfig.value,
      })

      Object.assign(request, requestFromApi(updated))
      hydrateActiveRequestState()
      requestSaveStatus.value = 'saved'
    } catch (err: any) {
      requestSaveStatus.value = 'error'
      requestSaveError.value = err?.data?.error || err?.message || 'Failed to save request'
    }
  }

  const executeActiveRequest = async () => {
    if (isWebSocketRequest.value) {
      if (webSocketState.value === 'connected') {
        await closeWebSocketSession()
        return
      }

      await connectWebSocketSession()
      return
    }

    loading.value = true
    const { post } = useApiClient()
    const authRows = authRowsForExecution(requestAuthConfig.value)

    try {
      const response = await post<WorkbenchResponse>('/execute', {
        method: activeRequest.value.method,
        url: `${requestTarget.value.baseUrl}${activeRequest.value.path}`,
        headers: [...headers.value.filter(h => h.enabled && h.key), ...authRows.headers],
        queryParams: [...queryParams.value.filter(p => p.enabled && p.key), ...authRows.queryParams],
        body: requestBody.value,
      })

      responseData.value = response
    } catch (err: any) {
      responseData.value = {
        status: 0,
        statusText: 'Error',
        duration: 0,
        size: 0,
        ttfb: 0,
        executionTarget: '',
        requestId: '',
        headers: [],
        body: err.message || 'An unknown error occurred',
      }
    } finally {
      loading.value = false
    }
  }

  const appendWebSocketEvent = (event: Omit<WebSocketTimelineEvent, 'id' | 'timestamp'>) => {
    webSocketEvents.value.unshift({
      id: crypto.randomUUID(),
      timestamp: new Date().toLocaleTimeString('en-GB', { hour12: false }),
      ...event,
    })
  }

  const connectWebSocketSession = async () => {
    if (webSocketState.value === 'connecting' || webSocketState.value === 'connected') return

    loading.value = true
    webSocketState.value = 'connecting'
    webSocketEvents.value = []
    const authRows = authRowsForExecution(requestAuthConfig.value)
    
    appendWebSocketEvent({
      direction: 'system',
      title: 'Opening backend websocket session',
      payload: `GET ${requestTarget.value.baseUrl}${activeRequest.value.path}\nUpgrade: websocket\nConnection: Upgrade`,
    })

    const { post } = useApiClient()
    try {
      const resp = await post<{ id: string }>('/execute/ws', {
        method: activeRequest.value.method,
        url: `${requestTarget.value.baseUrl}${activeRequest.value.path}`,
        headers: [...headers.value.filter(h => h.enabled && h.key), ...authRows.headers],
        queryParams: authRows.queryParams,
      })
      webSocketExecutionId.value = resp.id
    } catch (err: any) {
      webSocketState.value = 'error'
      loading.value = false
      appendWebSocketEvent({
        direction: 'error',
        title: 'Connection Failed',
        payload: err.message || 'Could not initiate websocket session',
      })
    }
  }

  const closeWebSocketSession = async () => {
    if (!webSocketExecutionId.value) return

    loading.value = true
    webSocketState.value = 'closing'
    
    const { delete: del } = useApiClient()
    try {
      await del(`/execute/${webSocketExecutionId.value}`)
    } catch (err: any) {
      appendWebSocketEvent({
        direction: 'error',
        title: 'Closure Error',
        payload: err.message,
      })
    } finally {
      loading.value = false
    }
  }

  const sendWebSocketMessage = async () => {
    if (webSocketState.value !== 'connected' || !webSocketExecutionId.value) return

    const payload = webSocketMessage.value.trim()
    const { post } = useApiClient()
    
    try {
      await post(`/execute/${webSocketExecutionId.value}/ws/send`, { payload })
      
      appendWebSocketEvent({
        direction: 'out',
        title: 'Message sent',
        payload,
        sizeBytes: new Blob([payload]).size,
      })
    } catch (err: any) {
      appendWebSocketEvent({
        direction: 'error',
        title: 'Send Error',
        payload: err.message,
      })
    }
  }

  const addQueryParam = () => {
    queryParams.value.push({
      id: crypto.randomUUID(),
      key: '',
      value: '',
      enabled: true,
    })
  }

  const removeQueryParam = (id: string) => {
    queryParams.value = queryParams.value.filter(param => param.id !== id)
  }

  const moveQueryParam = (paramId: string, targetIndex: number) => {
    const sourceIndex = queryParams.value.findIndex(param => param.id === paramId)
    if (sourceIndex === -1) return

    const [param] = queryParams.value.splice(sourceIndex, 1)
    if (!param) return

    const adjustedIndex = sourceIndex < targetIndex ? targetIndex - 1 : targetIndex
    queryParams.value.splice(Math.max(0, adjustedIndex), 0, param)
  }

  const addHeader = () => {
    headers.value.push({
      id: crypto.randomUUID(),
      key: '',
      value: '',
      enabled: true,
    })
  }

  const removeHeader = (id: string) => {
    headers.value = headers.value.filter(header => header.id !== id)
  }

  const moveHeader = (id: string, toIndex: number) => {
    const fromIndex = headers.value.findIndex(h => h.id === id)
    if (fromIndex === -1) return
    const item = headers.value[fromIndex]
    if (!item) return
    headers.value.splice(fromIndex, 1)
    headers.value.splice(toIndex, 0, item)
  }

  return {
    treeItems,
    rootRequests,
    sidebarWidth,
    activeRequestId,
    activeRequest,
    openTabs,
    editorPaneHeight,
    responseWidth,
    responsePosition,
    loading,
    responseData,
    requestTarget,
    queryParams,
    headers,
    requestBody,
    requestBodyLanguage,
    requestAuthConfig,
    requestSaveStatus,
    requestSaveError,
    webSocketState,
    webSocketMessage,
    webSocketMessageLanguage,
    webSocketEvents,
    isWebSocketRequest,
    closeTab,
    openRequest,
    setActiveRequest,
    setActiveRequestMethod,
    saveActiveRequestState,
    executeActiveRequest,
    connectWebSocketSession,
    closeWebSocketSession,
    sendWebSocketMessage,
    addQueryParam,
    removeQueryParam,
    moveQueryParam,
    addHeader,
    removeHeader,
    moveHeader,
    loadCollections,
    collectionsLoading,
    collectionsError,
    addFolder,
    addRequest,
    moveRequest,
    deleteItem,
  }
}
