package numbers

import "math"

// Truncat reduces floating point to a cetain accuracy it stops on acc > 10
//
//	acc = 0 - do nothing (3.3333-> 3.3333)
//	acc = 1 - 1 floating point x 10 (3.3333-> 3.3)
//	acc = 2 - 2 floating points x 100 (3.3333-> 3.33)
//	acc = 3 - 3 floating points x 1000 (3.3333-> 3.333)
func Truncat(n float64, acc int) float64 {
	if acc <= 0 {
		return n
	}
	p := 10.0
	for i := 1; i <= acc; i++ {
		p *= 10
		if acc > 10 {
			break
		}
	}
	return float64(int(n*p)) / p
}

// FloorOrCiel returns the nearest fraction for the float truncated
//
//	n: 3.33, fraction: 4, acc:2 - 3.25
//	n: 3.45, fraction: 4, acc:2 - 3.50
//	n: 7.78, fraction: 4, acc:2 - 7.75
//	n: 7.88, fraction: 4, acc:2 - 8.00
func FloorOrCiel(n float64, fraction, acc int) float64 {
	fr := n * float64(fraction)
	rmod := fr - float64(int(fr))
	if rmod > 0.5 {
		n = math.Ceil(fr) / float64(fraction)
	} else {
		n = math.Floor(fr) / float64(fraction)
	}
	return Truncat(n, acc)
}

func Ciel(n float64, fraction, acc int) float64 {
	n = math.Ceil(n*float64(fraction)) / float64(fraction)
	return Truncat(n, acc)
}

func Floor(n float64, fraction, acc int) float64 {
	n = math.Floor(n*float64(fraction)) / float64(fraction)
	return Truncat(n, acc)
}
