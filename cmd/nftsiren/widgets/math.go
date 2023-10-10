package widgets

import "golang.org/x/exp/constraints"

func clamp[T constraints.Integer | constraints.Float](a, min, max T) T {
	if a < min {
		return min
	}
	if a > max {
		return max
	}
	return a
}

func abs[T constraints.Integer | constraints.Float](a T) T {
	if a < 0 {
		return -a
	}
	return a
}
