package util

import (
	"fmt"
	"math"
)

func FloatToTimestamp(seconds float64) string {
	seconds = math.Round(seconds*1000) / 1000

	wholeSeconds := int64(seconds)
	milliseconds := int((seconds - float64(wholeSeconds)) * 1000)

	hours := wholeSeconds / 3600
	remaining := wholeSeconds % 3600
	minutes := remaining / 60
	secs := remaining % 60

	result := fmt.Sprintf("%02d:%02d:%02d.%03d", hours, minutes, secs, milliseconds)
	return result
}
