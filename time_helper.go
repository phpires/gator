package main

import (
	"fmt"
	"time"
)

func parseTime(timeString string) (time.Time, error) {
	formats := []string{
		time.RFC1123,
		time.RFC1123Z,
		time.RFC3339,
		time.RFC3339Nano,
	}

	for _, format := range formats {
		if t, err := time.Parse(format, timeString); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("Couldn't format time: %s. Unknown pattern?", timeString)
}
