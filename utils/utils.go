package utils

import "fmt"

func SecondsToReadable(seconds int) string {
	days := seconds / 86400
	seconds %= 86400
	hours := seconds / 3600
	seconds %= 3600
	minutes := seconds / 60
	seconds %= 60

	timeString := ""

	if days > 0 {
		timeString += fmt.Sprintf("%dd ", days)
	}
	if hours > 0 || days > 0 { // Include hours if there are any days
		timeString += fmt.Sprintf("%dh ", hours)
	}
	if minutes > 0 || hours > 0 || days > 0 { // Include minutes if there are any hours or days
		timeString += fmt.Sprintf("%dm ", minutes)
	}
	if seconds > 0 || minutes > 0 || hours > 0 || days > 0 { // Include seconds if there are any minutes, hours, or days
		timeString += fmt.Sprintf("%ds", seconds)
	}

	return timeString
}
