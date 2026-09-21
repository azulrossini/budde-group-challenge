import { describe, expect, it } from 'vitest'
import { basisPointsToInput, formatCents, formatPercentage, parsePercentage, split } from './money'

describe('parsePercentage', () => {
  it.each([
    ['33.33', 3333],
    ['33,5', 3350],
    ['100', 10000],
    ['0.01', 1],
    ['  50  ', 5000],
  ])('parses %s as %i basis points', (input, want) => {
    expect(parsePercentage(input)).toBe(want)
  })

  it.each([['33.333'], ['abc'], ['-5'], ['']])('rejects %s', (input) => {
    expect(parsePercentage(input)).toBeNull()
  })
})

describe('formatCents', () => {
  it('formats cents as euros', () => {
    expect(formatCents(12050)).toBe('€120.50')
    expect(formatCents(0)).toBe('€0.00')
  })
})

describe('formatPercentage', () => {
  it('renders basis points with a % suffix', () => {
    expect(formatPercentage(3333)).toBe('33.33 %')
    expect(formatPercentage(10000)).toBe('100.00 %')
  })
})

describe('basisPointsToInput', () => {
  it('drops trailing zeros for whole percentages', () => {
    expect(basisPointsToInput(5000)).toBe('50')
  })

  it('keeps the fraction for non-whole percentages', () => {
    expect(basisPointsToInput(3333)).toBe('33.33')
  })
})

describe('split', () => {
  it('gives the leftover cent to the largest remainder', () => {
    expect(split(1000, [3333, 3333, 3334])).toEqual([333, 333, 334])
  })

  it('gives everything to a single share at 100%', () => {
    expect(split(12050, [10000])).toEqual([12050])
  })

  it('splits a zero total into zeros', () => {
    expect(split(0, [3333, 3333, 3334])).toEqual([0, 0, 0])
  })

  it.each([
    [1000, [3333, 3333, 3334]],
    [12050, [10000]],
    [0, [3333, 3333, 3334]],
    [999, [2500, 2500, 2500, 2500]],
  ])('sums to the total for total=%i basisPoints=%j', (totalCents, basisPoints) => {
    const amounts = split(totalCents, basisPoints)
    expect(amounts.reduce((sum, amount) => sum + amount, 0)).toBe(totalCents)
  })
})
