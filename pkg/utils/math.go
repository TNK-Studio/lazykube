package utils

import "math"

func MaxFloat64(arr []float64) float64 {
	if len(arr) == 0 {
		return math.MaxFloat64
	}

	max := arr[0]
	for _, num := range arr {
		max = math.Max(max, num)
	}
	return max
}

func MinFloat64(arr []float64) float64 {
	if len(arr) == 0 {
		return float64(math.MinInt64)
	}

	min := arr[0]
	for _, num := range arr {
		min = math.Min(min, num)
	}
	return min
}
