import type { components } from '../types/api.gen'
import { ApiError, ErrorCode, type FieldError } from '../types/errors'
import { REQUEST_TIMEOUT_MS } from '../utils/constants'
import { MSG_INTERNAL, MSG_NETWORK, MSG_TIMEOUT } from '../utils/messages'

type ErrorResponseBody = components['schemas']['ErrorResponse']

const HEADER_CONTENT_TYPE = 'Content-Type'
const CONTENT_TYPE_JSON = 'application/json'

interface RequestOptions {
  method?: string
  body?: unknown
}

// request is the only place in the app that calls fetch. It applies a
// timeout and converts every failure into an ApiError.
async function request<T>(path: string, options: RequestOptions = {}): Promise<T> {
  const controller = new AbortController()
  const timeoutId = setTimeout(() => controller.abort(), REQUEST_TIMEOUT_MS)

  let response: Response
  try {
    response = await fetch(path, {
      method: options.method ?? 'GET',
      headers: options.body !== undefined ? { [HEADER_CONTENT_TYPE]: CONTENT_TYPE_JSON } : undefined,
      body: options.body !== undefined ? JSON.stringify(options.body) : undefined,
      signal: controller.signal,
    })
  } catch (err) {
    if (err instanceof DOMException && err.name === 'AbortError') {
      throw new ApiError({ code: ErrorCode.Timeout, message: MSG_TIMEOUT })
    }
    throw new ApiError({ code: ErrorCode.Network, message: MSG_NETWORK })
  } finally {
    clearTimeout(timeoutId)
  }

  if (!response.ok) {
    throw await toApiError(response)
  }

  return (await response.json()) as T
}

async function toApiError(response: Response): Promise<ApiError> {
  let body: ErrorResponseBody | undefined
  try {
    body = (await response.json()) as ErrorResponseBody
  } catch {
    body = undefined
  }

  if (body && typeof body.code === 'string' && typeof body.message === 'string') {
    return new ApiError({
      code: body.code as ErrorCode,
      message: body.message,
      details: body.details as FieldError[] | undefined,
      requestId: body.requestId,
      status: response.status,
    })
  }

  return new ApiError({ code: ErrorCode.InternalError, message: MSG_INTERNAL, status: response.status })
}

export function get<T>(path: string): Promise<T> {
  return request<T>(path)
}

export function putJson<T>(path: string, body: unknown): Promise<T> {
  return request<T>(path, { method: 'PUT', body })
}
