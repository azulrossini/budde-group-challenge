import { CURRENCY, FULL_PERCENT_BASIS_POINTS, LOCALE, MAX_DECIMAL_PLACES } from './constants'

const PERCENTAGE_PATTERN = /^(\d+)(?:[.,](\d{1,2}))?$/

// parsePercentage parses user input for a percentage into basis points
// (10000 = 100.00%), without using parseFloat/float math. Accepts up to
// MAX_DECIMAL_PLACES digits after '.' or ','; returns null for anything
// else (empty, non-numeric, negative, or too many decimals).
export function parsePercentage(input: string): number | null {
  const match = PERCENTAGE_PATTERN.exec(input.trim())
  if (!match) {
    return null
  }

  const [, wholePart, fractionPart = ''] = match
  const paddedFraction = fractionPart.padEnd(MAX_DECIMAL_PLACES, '0')
  return Number(wholePart) * 100 + Number(paddedFraction)
}

// formatPercentageValue renders basis points as "33.33" (no unit suffix).
export function formatPercentageValue(basisPoints: number): string {
  const whole = Math.trunc(basisPoints / 100)
  const fraction = String(Math.abs(basisPoints % 100)).padStart(MAX_DECIMAL_PLACES, '0')
  return `${whole}.${fraction}`
}

// formatPercentage renders basis points as "33.33 %".
export function formatPercentage(basisPoints: number): string {
  return `${formatPercentageValue(basisPoints)} %`
}

// basisPointsToInput renders basis points as a plain editable value
// ("50", "33.33") — like formatPercentageValue, but without a trailing
// ".00" for whole percentages, since that reads awkwardly in a text input.
export function basisPointsToInput(basisPoints: number): string {
  const whole = Math.trunc(basisPoints / 100)
  const fraction = basisPoints % 100
  if (fraction === 0) {
    return String(whole)
  }
  return `${whole}.${String(Math.abs(fraction)).padStart(MAX_DECIMAL_PLACES, '0')}`
}

// formatCents renders integer cents as a localized euro amount.
export function formatCents(cents: number): string {
  return new Intl.NumberFormat(LOCALE, { style: 'currency', currency: CURRENCY }).format(cents / 100)
}

// split computes each share's amount in cents from its percentage in basis
// points. Mirrors backend/internal/money.Split exactly, including its test
// cases: it assumes basisPoints sums to FULL_PERCENT_BASIS_POINTS (10000)
// — callers must validate that first. Each amount is floor(total * bp /
// 10000); the leftover cents from rounding go one at a time to the shares
// with the largest remainders, ties broken by list order (Array.prototype
// .sort is stable), so the result always sums to totalCents.
export function split(totalCents: number, basisPoints: readonly number[]): number[] {
  const amounts: number[] = []
  const remainders: number[] = []

  let distributed = 0
  for (const bp of basisPoints) {
    const product = totalCents * bp
    const amount = Math.floor(product / FULL_PERCENT_BASIS_POINTS)
    amounts.push(amount)
    remainders.push(product % FULL_PERCENT_BASIS_POINTS)
    distributed += amount
  }

  // order and remainders/amounts always have the same length as
  // basisPoints, so these indexed accesses are safe.
  const order = amounts.map((_, i) => i)
  order.sort((a, b) => remainders[b]! - remainders[a]!)

  const leftover = totalCents - distributed
  for (let i = 0; i < leftover; i++) {
    const index = order[i]!
    amounts[index] = amounts[index]! + 1
  }

  return amounts
}
