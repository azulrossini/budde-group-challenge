package money

import "sort"

const FullPercentBasisPoints = 10000

// Split computes each share's amount in cents from its percentage in basis
// points. It assumes basisPoints sums to FullPercentBasisPoints (10000) —
// callers must validate that before calling. Each amount is floor(total *
// bp / 10000); the leftover cents from rounding are distributed one at a
// time to the shares with the largest remainders, ties broken by list
// order, so the result always sums to totalCents.
func Split(totalCents int64, basisPoints []int32) []int64 {
	amounts := make([]int64, len(basisPoints))
	remainders := make([]int64, len(basisPoints))

	var distributed int64
	for i, bp := range basisPoints {
		product := totalCents * int64(bp)
		amounts[i] = product / FullPercentBasisPoints
		remainders[i] = product % FullPercentBasisPoints
		distributed += amounts[i]
	}

	order := make([]int, len(basisPoints))
	for i := range order {
		order[i] = i
	}
	sort.SliceStable(order, func(a, b int) bool {
		return remainders[order[a]] > remainders[order[b]]
	})

	leftover := totalCents - distributed
	for i := int64(0); i < leftover; i++ {
		amounts[order[i]]++
	}

	return amounts
}
