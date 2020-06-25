package internal

import (
	"fmt"
	"time"
)

const dateLayout1 = "2006-01-02 15:04:05"

func parseDate(layout string, date string) (time.Time, error) {
	t, err := time.Parse(layout, date)
	if err != nil {
		return t, fmt.Errorf("unexpected date value '%s'. the date must be in %s format", date, layout)
	}
	return t, nil
}
