package utils

import (
	"fmt"
	"time"
)

const template = "2006-01"

func ParseYearMonth(input string) (start time.Time, end time.Time, err error) {

	start, err = time.Parse(template, input)
	if err != nil {
		err = fmt.Errorf("failed to parse input '%s': %w", input, err)
		return
	}
	end = start.AddDate(0, 1, 0).Add(-time.Nanosecond)

	return
}
