import type {
  HeaderItem,
  RequestTarget,
  WebSocketTimelineEvent,
} from '~/composables/useWorkbench'
import type {
  QueryParamItem,
  TreeItem,
  WorkbenchResponse,
} from '~/composables/useWorkbench'

export const mockWorkbenchTree: TreeItem[] = [
  {
    name: 'Authentication',
    icon: 'PhGlobe',
    requests: [
      { id: 'auth-login', method: 'POST', name: 'Login', path: '/auth/login', active: true },
      { id: 'auth-refresh', method: 'POST', name: 'Refresh token', path: '/auth/refresh' },
      { id: 'auth-me', method: 'GET', name: 'Current user', path: '/auth/me' },
    ],
  },
  {
    name: 'Realtime',
    icon: 'PhLightning',
    requests: [
      { id: 'events-ticket', method: 'POST', name: 'Create event ticket', path: '/events/ticket' },
      { id: 'events-stream', method: 'SOCKET', name: 'Event stream', path: '/events' },
    ],
  },
]

export const mockActiveRequestId = 'auth-login'

export const mockOpenTabIds = [
  'auth-login',
  'auth-refresh',
  'auth-me',
  'events-stream',
]

export const mockRequestTarget: RequestTarget = {
  baseUrl: '{{apiBaseUrl}}',
  path: '/auth/login',
}

export const mockQueryParams: QueryParamItem[] = [
  { id: 'param-base-url', key: 'base_url', value: '{{apiBaseUrl}}', enabled: true },
  { id: 'param-workspace-id', key: 'workspace_id', value: '{{currentWorkspace}}', enabled: false },
  { id: 'param-include-meta', key: 'include_meta', value: 'true', enabled: true },
]

export const mockHeaders: HeaderItem[] = [
  { id: 'header-authorization', key: 'Authorization', value: 'Bearer {{accessToken}}', enabled: true },
  { id: 'header-content-type', key: 'Content-Type', value: 'application/json', enabled: true },
]

export const mockBody = `{
  "login": "owner@example.com",
  "password": "••••••••••••"
}`

export const mockWorkbenchResponse: WorkbenchResponse = {
  status: 200,
  statusText: 'OK',
  duration: 126,
  size: 1433,
  ttfb: 82,
  executionTarget: 'backend',
  requestId: 'req_018f6c93-7b2a',
  headers: [
    { id: 'res-content-type', key: 'content-type', value: 'application/json', enabled: true },
    { id: 'res-cache-control', key: 'cache-control', value: 'no-store', enabled: true },
  ],
  body: `{
  "accessToken": "••••••••••••••••••••••••••••••••",
  "userId": "018f6c93-7b2a-7f5c-a7b1",
  "globalRole": "manager"
}`,
}

export const mockWorkbenchResponses: Record<string, WorkbenchResponse> = {
  'auth-login': mockWorkbenchResponse,
  'auth-refresh': {
    status: 200,
    statusText: 'OK',
    duration: 94,
    size: 921,
    ttfb: 61,
    executionTarget: 'backend',
    requestId: 'req_018f6c94-refresh',
    headers: [
      { id: 'res-refresh-content-type', key: 'content-type', value: 'application/json', enabled: true },
      { id: 'res-refresh-cache-control', key: 'cache-control', value: 'no-store', enabled: true },
    ],
    body: `{
  "accessToken": "••••••••••••••••••••••••••••••••",
  "expiresIn": 900,
  "tokenType": "Bearer"
}`,
  },
  'auth-me': {
    status: 200,
    statusText: 'OK',
    duration: 48,
    size: 1126,
    ttfb: 33,
    executionTarget: 'backend',
    requestId: 'req_018f6c95-me',
    headers: [
      { id: 'res-me-content-type', key: 'content-type', value: 'application/json', enabled: true },
      { id: 'res-me-vary', key: 'vary', value: 'Authorization', enabled: true },
    ],
    body: `{
  "id": "018f6c93-7b2a-7f5c-a7b1",
  "email": "oussema@test.io",
  "username": "oussema_admin",
  "globalRole": "manager",
  "workspaceMemberships": [
    {
      "workspaceId": "018f6c98-core",
      "role": "manager",
      "status": "active"
    }
  ]
}`,
  },
  'events-ticket': {
    status: 201,
    statusText: 'Created',
    duration: 72,
    size: 614,
    ttfb: 55,
    executionTarget: 'backend',
    requestId: 'req_018f6c96-ticket',
    headers: [
      { id: 'res-ticket-content-type', key: 'content-type', value: 'application/json', enabled: true },
      { id: 'res-ticket-cache-control', key: 'cache-control', value: 'no-store', enabled: true },
    ],
    body: `{
  "ticket": "evt_tk_72f8f7a188b84c9c",
  "expiresIn": 30,
  "streamUrl": "/events?ticket=evt_tk_72f8f7a188b84c9c"
}`,
  },
  'events-stream': {
    status: 200,
    statusText: 'Stream Open',
    duration: 12,
    size: 0,
    ttfb: 12,
    executionTarget: 'backend',
    requestId: 'req_018f6c97-stream',
    headers: [
      { id: 'res-stream-content-type', key: 'content-type', value: 'text/event-stream', enabled: true },
      { id: 'res-stream-cache-control', key: 'cache-control', value: 'no-cache', enabled: true },
    ],
    body: `event: workspace.member.added
data: {"userId":"018f6c93-7b2a-7f5c-a7b1","workspaceId":"018f6c98-core","role":"developer"}

event: request.executed
data: {"requestId":"req_018f6c97-stream","status":200}`,
  },
}

export const mockWebSocketEvents: WebSocketTimelineEvent[] = [
  {
    id: 'ws-event-system-ready',
    direction: 'system',
    timestamp: '21:42:10',
    title: 'Session scaffold ready',
    payload: 'Backend will own the target websocket and stream events to this timeline.',
  },
  {
    id: 'ws-event-initial-inbound',
    direction: 'in',
    timestamp: '21:42:11',
    title: 'Mock inbound event',
    payload: `{
  "event": "workspace.member.added",
  "userId": "018f6c93-7b2a-7f5c-a7b1",
  "workspaceId": "018f6c98-core"
}`,
    sizeBytes: 128,
  },
]
