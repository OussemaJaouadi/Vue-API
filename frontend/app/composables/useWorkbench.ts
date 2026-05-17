import { PhGlobe, PhLightning } from '@phosphor-icons/vue'
import { useStorage } from '@vueuse/core'
import type { Component } from 'vue'
import {
  mockWebSocketEvents,
  mockActiveRequestId,
  mockBody,
  mockHeaders,
  mockOpenTabIds,
  mockQueryParams,
  mockRequestTarget,
  mockWorkbenchResponse,
  mockWorkbenchResponses,
  mockWorkbenchTree,
} from '~/mock-data/workbench'

export type ApiMethod = 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE' | 'SOCKET'
export type WorkbenchIconKey = 'PhGlobe' | 'PhLightning'

export interface RequestItem {
  id: string
  method: ApiMethod
  name: string
  path: string
  active?: boolean
}

export interface TreeItem {
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
  size: string
  ttfb: number
  decoded: string
  executionTarget: string
  requestId: string
  headers: HeaderItem[]
  body: string
}

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

export function useWorkbench() {
  const treeItems = useState<TreeItem[]>('workbench:tree', () => structuredClone(mockWorkbenchTree))
  const rootRequests = useState<RequestItem[]>('workbench:root-requests', () => [])
  
  onMounted(() => {
    const hasSocket = treeItems.value.some(g => g.requests.some(r => r.method === 'SOCKET'))
    if (!hasSocket) {
      treeItems.value = structuredClone(mockWorkbenchTree)
    }
  })

  const allRequests = computed(() => [
    ...rootRequests.value,
    ...treeItems.value.flatMap(group => group.requests)
  ])

  const activeRequestId = useState<string>('workbench:active-request-id', () => mockActiveRequestId)
  
  // Persistent Layout State
  const sidebarWidth = useStorage<number>('workbench:sidebar-width', 260)
  const editorPaneHeight = useStorage<number>('workbench:editor-height', 50)
  const responseWidth = useStorage<number>('workbench:response-width', 400)
  const responsePosition = useStorage<'bottom' | 'right'>('workbench:response-position', 'bottom')
  
  const openTabs = useState<RequestItem[]>('workbench:open-tabs', () => {
    return mockOpenTabIds
      .map(id => allRequests.value.find(request => request.id === id))
      .filter((request): request is RequestItem => Boolean(request))
  })

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
  const responseData = useState<WorkbenchResponse | null>('workbench:response', () => structuredClone(mockWorkbenchResponse))
  const requestTarget = useState<RequestTarget>('workbench:request-target', () => structuredClone(mockRequestTarget))
  const queryParams = useState<QueryParamItem[]>('workbench:query-params', () => structuredClone(mockQueryParams))
  const headers = useState<HeaderItem[]>('workbench:headers', () => structuredClone(mockHeaders))
  const requestBody = useState<string>('workbench:body', () => mockBody)
  const requestBodyLanguage = useState<BodyLanguage>('workbench:body-language', () => 'json')
  const webSocketState = useState<WebSocketConnectionState>('workbench:ws-state', () => 'idle')
  const webSocketMessage = useState<string>('workbench:ws-message', () => `{
  "type": "ping",
  "requestId": "{{requestId}}"
}`)
  const webSocketMessageLanguage = useState<BodyLanguage>('workbench:ws-message-language', () => 'json')
  const webSocketEvents = useState<WebSocketTimelineEvent[]>('workbench:ws-events', () => structuredClone(mockWebSocketEvents))

  const isWebSocketRequest = computed(() => activeRequest.value.method === 'SOCKET')

  const addFolder = (name: string = 'New Collection') => {
    treeItems.value.push({
      name,
      icon: 'PhGlobe',
      requests: []
    })
  }

  const addRequest = (folderName?: string) => {
    const newReq: RequestItem = {
      id: crypto.randomUUID(),
      method: 'GET',
      name: 'New Request',
      path: '/path',
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

  const deleteItem = (id: string, isFolder: boolean) => {
    if (isFolder) {
      treeItems.value = treeItems.value.filter(f => f.name !== id)
    } else {
      rootRequests.value = rootRequests.value.filter(r => r.id !== id)
      treeItems.value.forEach(f => {
        f.requests = f.requests.filter(r => r.id !== id)
      })
      if (activeRequestId.value === id) {
        activeRequestId.value = allRequests.value[0]?.id || ''
      }
    }
  }

  const closeTab = (tab: RequestItem) => {
    if (openTabs.value.length === 1) return
    const tabIndex = openTabs.value.findIndex(item => item.id === tab.id)
    openTabs.value = openTabs.value.filter(item => item.id !== tab.id)
    if (activeRequestId.value === tab.id) {
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

    await new Promise(resolve => window.setTimeout(resolve, 260))

    const response = mockWorkbenchResponses[activeRequest.value.id] || mockWorkbenchResponse
    responseData.value = {
      ...structuredClone(response),
      requestId: `req_${crypto.randomUUID().slice(0, 13)}`,
    }
    loading.value = false
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
    appendWebSocketEvent({
      direction: 'system',
      title: 'Opening backend websocket session',
      payload: `GET ${requestTarget.value.baseUrl}${activeRequest.value.path}\nUpgrade: websocket\nConnection: Upgrade`,
    })

    await new Promise(resolve => window.setTimeout(resolve, 320))

    webSocketState.value = 'connected'
    loading.value = false
    appendWebSocketEvent({
      direction: 'system',
      title: 'Connected through backend proxy',
      payload: '101 Switching Protocols\nexecutionId: ws_mock_018f6c97',
      sizeBytes: 54,
    })
  }

  const closeWebSocketSession = async () => {
    if (webSocketState.value !== 'connected') return

    loading.value = true
    webSocketState.value = 'closing'
    appendWebSocketEvent({
      direction: 'system',
      title: 'Closing websocket session',
      payload: 'client requested close',
    })

    await new Promise(resolve => window.setTimeout(resolve, 180))

    webSocketState.value = 'closed'
    loading.value = false
    appendWebSocketEvent({
      direction: 'system',
      title: 'Socket closed',
      payload: 'code=1000 reason=normal closure',
    })
  }

  const sendWebSocketMessage = async () => {
    if (webSocketState.value !== 'connected') return

    const payload = webSocketMessage.value.trim()
    appendWebSocketEvent({
      direction: 'out',
      title: 'Message sent',
      payload,
      sizeBytes: new Blob([payload]).size,
    })

    await new Promise(resolve => window.setTimeout(resolve, 240))

    appendWebSocketEvent({
      direction: 'in',
      title: 'Message received',
      payload: `{
  "type": "pong",
  "receivedAt": "${new Date().toISOString()}",
  "echo": ${JSON.stringify(payload)}
}`,
      sizeBytes: 126,
    })
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
    webSocketState,
    webSocketMessage,
    webSocketMessageLanguage,
    webSocketEvents,
    isWebSocketRequest,
    closeTab,
    openRequest,
    setActiveRequest,
    setActiveRequestMethod,
    executeActiveRequest,
    connectWebSocketSession,
    closeWebSocketSession,
    sendWebSocketMessage,
    addQueryParam,
    removeQueryParam,
    moveQueryParam,
    addHeader,
    removeHeader,
    addFolder,
    addRequest,
    moveRequest,
    deleteItem,
  }
}
