import { computed, ref } from 'vue'
import { getShares, replaceShares } from '../api/billsApi'
import type { components } from '../types/api.gen'
import { ApiError, ErrorCode } from '../types/errors'
import { LoadStatus } from '../types/status'
import { BILL_ID, FULL_PERCENT_BASIS_POINTS, MAX_NAME_LENGTH, MAX_SHARES } from '../utils/constants'
import {
  MSG_AT_LEAST_ONE_SHARE,
  MSG_DUPLICATE_NAME,
  MSG_INTERNAL,
  MSG_INVALID_PERCENTAGE,
  MSG_NAME_REQUIRED,
  MSG_PERCENTAGE_OUT_OF_RANGE,
  nameTooLongMessage,
  sumMismatchMessage,
  tooManySharesMessage,
} from '../utils/messages'
import { basisPointsToInput, formatPercentageValue, parsePercentage, split } from '../utils/money'

type Bill = components['schemas']['Bill']
type BillSharesResponse = components['schemas']['BillSharesResponse']

export interface ShareRowData {
  key: string
  id: number | null
  name: string
  percentageInput: string
}

export interface RowError {
  name?: string
  percentage?: string
}

export function useShares() {
  const status = ref<LoadStatus>(LoadStatus.Idle)
  const isSaving = ref(false)
  const bill = ref<Bill | null>(null)
  const rows = ref<ShareRowData[]>([])
  const loadError = ref<ApiError | null>(null)
  const saveError = ref<ApiError | null>(null)
  const saveSucceeded = ref(false)

  let savedSnapshot = ''
  let nextKey = 0

  const rowErrors = computed<Record<string, RowError>>(() => {
    const errors: Record<string, RowError> = {}
    const seenNames = new Set<string>()

    for (const row of rows.value) {
      const error: RowError = {}
      const name = row.name.trim()

      if (name === '') {
        error.name = MSG_NAME_REQUIRED
      } else if (name.length > MAX_NAME_LENGTH) {
        error.name = nameTooLongMessage(MAX_NAME_LENGTH)
      } else {
        const key = name.toLowerCase()
        if (seenNames.has(key)) {
          error.name = MSG_DUPLICATE_NAME
        }
        seenNames.add(key)
      }

      const basisPoints = parsePercentage(row.percentageInput)
      if (basisPoints === null) {
        error.percentage = MSG_INVALID_PERCENTAGE
      } else if (basisPoints <= 0 || basisPoints > FULL_PERCENT_BASIS_POINTS) {
        error.percentage = MSG_PERCENTAGE_OUT_OF_RANGE
      }

      if (error.name || error.percentage) {
        errors[row.key] = error
      }
    }

    return errors
  })

  const listError = computed<string | null>(() => {
    if (rows.value.length === 0) {
      return MSG_AT_LEAST_ONE_SHARE
    }
    if (rows.value.length > MAX_SHARES) {
      return tooManySharesMessage(MAX_SHARES)
    }
    return null
  })

  const totalBasisPoints = computed(() =>
    rows.value.reduce((total, row) => total + (parsePercentage(row.percentageInput) ?? 0), 0),
  )

  const sumMessage = computed(() => {
    if (listError.value || totalBasisPoints.value === FULL_PERCENT_BASIS_POINTS) {
      return null
    }
    return sumMismatchMessage(formatPercentageValue(totalBasisPoints.value))
  })

  const hasValidationErrors = computed(
    () => listError.value !== null || Object.keys(rowErrors.value).length > 0,
  )

  const hasChanges = computed(() => snapshotOf(rows.value) !== savedSnapshot)

  const canSave = computed(
    () =>
      status.value === LoadStatus.Ready &&
      !isSaving.value &&
      !hasValidationErrors.value &&
      totalBasisPoints.value === FULL_PERCENT_BASIS_POINTS &&
      hasChanges.value,
  )

  const amounts = computed<Record<string, number>>(() => {
    if (!bill.value || hasValidationErrors.value || totalBasisPoints.value !== FULL_PERCENT_BASIS_POINTS) {
      return {}
    }

    const basisPointsList = rows.value.map((row) => parsePercentage(row.percentageInput) ?? 0)
    const splitAmounts = split(bill.value.totalCents, basisPointsList)

    const result: Record<string, number> = {}
    rows.value.forEach((row, i) => {
      result[row.key] = splitAmounts[i] ?? 0
    })
    return result
  })

  function snapshotOf(list: ShareRowData[]): string {
    return JSON.stringify(
      list.map((row) => ({ name: row.name.trim(), percentageBasisPoints: parsePercentage(row.percentageInput) })),
    )
  }

  function applyServerState(response: BillSharesResponse) {
    bill.value = response.bill
    rows.value = response.shares.map((share) => ({
      key: String(share.id),
      id: share.id,
      name: share.name,
      percentageInput: basisPointsToInput(share.percentageBasisPoints),
    }))
    savedSnapshot = snapshotOf(rows.value)
  }

  function toApiError(err: unknown): ApiError {
    return err instanceof ApiError ? err : new ApiError({ code: ErrorCode.InternalError, message: MSG_INTERNAL })
  }

  async function load() {
    status.value = LoadStatus.Loading
    loadError.value = null
    try {
      applyServerState(await getShares(BILL_ID))
      status.value = LoadStatus.Ready
    } catch (err) {
      loadError.value = toApiError(err)
      status.value = LoadStatus.Error
    }
  }

  async function save() {
    if (!canSave.value) {
      return
    }

    isSaving.value = true
    saveError.value = null
    saveSucceeded.value = false

    try {
      const response = await replaceShares(BILL_ID, {
        shares: rows.value.map((row) => ({
          name: row.name.trim(),
          percentageBasisPoints: parsePercentage(row.percentageInput) ?? 0,
        })),
      })
      applyServerState(response)
      saveSucceeded.value = true
    } catch (err) {
      saveError.value = toApiError(err)
    } finally {
      isSaving.value = false
    }
  }

  function addPerson() {
    rows.value.push({ key: `new-${nextKey++}`, id: null, name: '', percentageInput: '' })
    saveSucceeded.value = false
  }

  function removePerson(key: string) {
    rows.value = rows.value.filter((row) => row.key !== key)
    saveSucceeded.value = false
  }

  function updateName(key: string, name: string) {
    const row = rows.value.find((r) => r.key === key)
    if (row) {
      row.name = name
    }
    saveSucceeded.value = false
  }

  function updatePercentage(key: string, percentageInput: string) {
    const row = rows.value.find((r) => r.key === key)
    if (row) {
      row.percentageInput = percentageInput
    }
    saveSucceeded.value = false
  }

  return {
    status,
    isSaving,
    bill,
    rows,
    rowErrors,
    listError,
    sumMessage,
    totalBasisPoints,
    canSave,
    amounts,
    loadError,
    saveError,
    saveSucceeded,
    load,
    save,
    addPerson,
    removePerson,
    updateName,
    updatePercentage,
  }
}
