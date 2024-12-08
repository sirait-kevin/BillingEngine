package helper

import (
	"time"
)

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

// WeeksBetween calculates the number of weeks between two dates
func WeeksBetween(start, end time.Time) int {
	if end.Before(start) {
		start, end = end, start
	}
	duration := end.Sub(start)
	weeks := int(duration.Hours() / (24 * 7))
	return weeks
}

// YearsBetween calculates the number of years between two dates
func YearsBetween(start, end time.Time) int {
	if end.Before(start) {
		start, end = end, start
	}

	years := end.Year() - start.Year()

	// Adjust for the cases where the end year is incomplete
	if end.Month() < start.Month() || (end.Month() == start.Month() && end.Day() < start.Day()) {
		years--
	}

	return years
}
