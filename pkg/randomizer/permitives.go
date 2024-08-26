package randomizer

import (
	"math/rand"

	"app/pkg/numbers"
)

func Bool() bool {
	return rand.Intn(2) == 1 //nolint: gosec // dw
}

func Int(min, max int) int {
	if max < min {
		return -1
	}
	// min, max := 2, 8
	// 2 + (8-2+1) +1 because Intn starts at 0
	// 2 + 1...6
	return min + rand.Intn(max-min+1) //nolint: gosec // dw
}

func Float64(min, max float64, acc int) float64 {
	r := min + rand.Float64()*(max-min) //nolint: gosec // dw
	return numbers.Truncat(r, acc)
}

func SliceElement[S ~[]E, E any](arr S) E {
	n := len(arr)
	if n == 0 {
		var e E
		return e
	}
	return arr[rand.Intn(n)] //nolint: gosec // dw
}

// TODO: add random strings: name, url, long description ...etc
// TODO: add random timestamps in range
