package valourgo

import (
	"fmt"
	"strings"
	"time"
)

type Time struct {
	time.Time
}

func (t *Time) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	if s == "" {
		return nil
	}

	layouts := []string{
		time.RFC3339,
		"2006-01-02T15:04:05",
	}

	for _, l := range layouts {
		if tt, err := time.Parse(l, s); err == nil {
			t.Time = tt
			return nil
		}
	}

	return fmt.Errorf("invalid time format: %s", s)
}
