package utils

import "strconv"

// FormatDuration converts a duration in seconds to a human-readable string.
// e.g. 1441 → "24 min", 3600 → "1 hr", 5400 → "1 hr 30 min"
func FormatDuration(seconds int) string {
	hours := seconds / 3600
	minutes := (seconds % 3600) / 60

	if hours > 0 && minutes > 0 {
		return strconv.Itoa(hours) + " hr " + strconv.Itoa(minutes) + " min"
	}
	if hours > 0 {
		return strconv.Itoa(hours) + " hr"
	}
	return strconv.Itoa(minutes) + " min"
}
