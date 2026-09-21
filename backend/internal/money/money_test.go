package money

import (
	"testing"
)

func TestSplit(t *testing.T) {
	tests := []struct {
		name        string
		totalCents  int64
		basisPoints []int32
		want        []int64
	}{
		{
			name:        "leftover cent goes to the largest remainder",
			totalCents:  1000,
			basisPoints: []int32{3333, 3333, 3334},
			want:        []int64{333, 333, 334},
		},
		{
			name:        "single share at 100%",
			totalCents:  12050,
			basisPoints: []int32{10000},
			want:        []int64{12050},
		},
		{
			name:        "zero total",
			totalCents:  0,
			basisPoints: []int32{3333, 3333, 3334},
			want:        []int64{0, 0, 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Split(tt.totalCents, tt.basisPoints)

			if len(got) != len(tt.want) {
				t.Fatalf("Split() returned %d amounts, want %d", len(got), len(tt.want))
			}
			var sum int64
			for i, amount := range got {
				if amount != tt.want[i] {
					t.Errorf("Split()[%d] = %d, want %d", i, amount, tt.want[i])
				}
				sum += amount
			}
			if sum != tt.totalCents {
				t.Errorf("Split() amounts sum to %d, want %d", sum, tt.totalCents)
			}
		})
	}
}
