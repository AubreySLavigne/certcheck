package certificate

import "time"

type dateRange struct {
	Start time.Time
	End   time.Time
}

func (d *dateRange) contains(t time.Time) bool {
	return d.Start.Before(t) && t.Before(d.End)
}
