package helper

import "time"

// MonthsBetween calculates the number of months between two dates
func MonthsBetween(start, end time.Time) int {
	if end.Before(start) {
		start, end = end, start
	}
	years := end.Year() - start.Year()
	months := int(end.Month()) - int(start.Month())
	if end.Day() < start.Day() {
		months -= 1
	}
	return years*12 + months
}
