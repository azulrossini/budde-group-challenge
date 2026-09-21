import { beforeEach, describe, expect, it, vi } from 'vitest'
import * as billsApi from '../api/billsApi'
import { useShares } from './useShares'

vi.mock('../api/billsApi')

describe('useShares', () => {
  beforeEach(() => {
    vi.mocked(billsApi.getShares).mockResolvedValue({
      bill: { id: 1, description: 'Team dinner', totalCents: 12050 },
      shares: [{ id: 1, name: 'Alice', percentageBasisPoints: 5000, amountCents: 6025 }],
      totalBasisPoints: 5000,
    })
  })

  it('canSave is false at a total of 9999 and true at 10000, with changes', async () => {
    const shares = useShares()
    await shares.load()

    const key = shares.rows.value[0]!.key

    shares.updatePercentage(key, '99.99')
    expect(shares.totalBasisPoints.value).toBe(9999)
    expect(shares.canSave.value).toBe(false)

    shares.updatePercentage(key, '100')
    expect(shares.totalBasisPoints.value).toBe(10000)
    expect(shares.canSave.value).toBe(true)
  })
})
