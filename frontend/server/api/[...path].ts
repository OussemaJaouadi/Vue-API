import {
  createError,
  defineEventHandler,
  getHeaders,
  getMethod,
  getQuery,
  getRouterParam,
  readRawBody,
} from 'h3'

const hopByHopHeaders = new Set([
  'connection',
  'content-length',
  'host',
  'keep-alive',
  'proxy-authenticate',
  'proxy-authorization',
  'te',
  'trailer',
  'transfer-encoding',
  'upgrade',
])

export default defineEventHandler(async (event) => {
  const config = useRuntimeConfig(event)
  const baseUrl = String(config.apiBaseUrl || '')
  const path = getRouterParam(event, 'path')

  if (!baseUrl) {
    throw createError({
      statusCode: 500,
      statusMessage: 'Backend API base URL is not configured',
    })
  }

  if (!path) {
    throw createError({
      statusCode: 404,
      statusMessage: 'API route not found',
    })
  }

  const targetUrl = new URL(path, baseUrl.endsWith('/') ? baseUrl : `${baseUrl}/`)
  const query = getQuery(event)
  for (const [key, value] of Object.entries(query)) {
    if (Array.isArray(value)) {
      for (const item of value) {
        targetUrl.searchParams.append(key, String(item))
      }
      continue
    }

    if (value !== undefined) {
      targetUrl.searchParams.set(key, String(value))
    }
  }

  const headers = new Headers()
  for (const [key, value] of Object.entries(getHeaders(event))) {
    if (!hopByHopHeaders.has(key.toLowerCase()) && value !== undefined) {
      headers.set(key, value)
    }
  }

  const method = getMethod(event)
  const init: RequestInit = {
    method,
    headers,
  }

  if (method !== 'GET' && method !== 'HEAD') {
    const body = await readRawBody(event, false)
    if (body) {
      init.body = body
    }
  }

  const response = await fetch(targetUrl, init)
  const responseHeaders = new Headers(response.headers)
  responseHeaders.delete('content-encoding')
  responseHeaders.delete('content-length')

  return new Response(response.body, {
    status: response.status,
    statusText: response.statusText,
    headers: responseHeaders,
  })
})
