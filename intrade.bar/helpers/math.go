package helpers

import "math"

func RoundUp(input float64) int{
	tmp := math.Round(input)
	if input - tmp < 0 {
		return int(tmp) + 1
	}

	return int(tmp)

}