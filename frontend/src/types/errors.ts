import type { components } from './api.gen'

export const ErrorCode = {
  ValidationError: 'VALIDATION_ERROR',
  BadRequest: 'BAD_REQUEST',
  NotFound: 'NOT_FOUND',
  Unauthorized: 'UNAUTHORIZED',
  InternalError: 'INTERNAL_ERROR',
  // Frontend-only codes: the backend never sends these, client.ts assigns
  // them for failures that never reach the server.
  Network: 'NETWORK',
  Timeout: 'TIMEOUT',
} as const

export type ErrorCode = (typeof ErrorCode)[keyof typeof ErrorCode]

export type FieldError = components['schemas']['FieldError']

export interface ApiErrorParams {
  code: ErrorCode
  message: string
  details?: FieldError[]
  requestId?: string
  status?: number
}

export class ApiError extends Error {
  readonly code: ErrorCode
  readonly details?: FieldError[]
  readonly requestId?: string
  readonly status?: number

  constructor(params: ApiErrorParams) {
    super(params.message)
    this.name = 'ApiError'
    this.code = params.code
    this.details = params.details
    this.requestId = params.requestId
    this.status = params.status
  }
}
