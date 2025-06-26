package randomizer

import (
	"math/rand"

	"app/pkg/numbers"
)

func Bool() bool {
	return rand.Intn(2) == 1 //nolint: gosec // dw
}

func Int(minval, maxval int) int {
	if maxval < minval {
		return -1
	}
	// min, max := 2, 8
	// 2 + (8-2+1) +1 because Intn starts at 0
	// 2 + 1...6
	return minval + rand.Intn(maxval-minval+1) //nolint: gosec // dw
}

func Float64(minval, maxval float64, acc int) float64 {
	r := minval + rand.Float64()*(maxval-minval) //nolint: gosec // dw
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
