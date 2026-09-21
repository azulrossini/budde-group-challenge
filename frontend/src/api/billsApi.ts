import { get, putJson } from './client'
import type { components } from '../types/api.gen'

type BillSharesResponse = components['schemas']['BillSharesResponse']
type ReplaceSharesRequest = components['schemas']['ReplaceSharesRequest']

const BASE_PATH = '/bills'
const SHARES_SEGMENT = 'shares'

export function getShares(billId: number): Promise<BillSharesResponse> {
  return get<BillSharesResponse>(`${BASE_PATH}/${billId}/${SHARES_SEGMENT}`)
}

export function replaceShares(billId: number, payload: ReplaceSharesRequest): Promise<BillSharesResponse> {
  return putJson<BillSharesResponse>(`${BASE_PATH}/${billId}/${SHARES_SEGMENT}`, payload)
}
