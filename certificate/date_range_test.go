package certificate

import (
	"testing"
	"time"
)

func TestDateRangeContains(t *testing.T) {

	r := dateRange{
		Start: time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC),
		End:   time.Date(2019, time.February, 14, 1, 2, 3, 4, time.UTC),
	}

	tests := []struct {
		input   time.Time
		expects bool
	}{
		{input: time.Date(2018, time.January, 2, 0, 0, 0, 0, time.UTC), expects: false},
		{input: r.Start, expects: false},
		{input: time.Date(2019, time.January, 2, 0, 0, 0, 0, time.UTC), expects: true},
		{input: r.End, expects: false},
		{input: time.Date(2020, time.January, 2, 0, 0, 0, 0, time.UTC), expects: false},
	}

	for _, test := range tests {

		if res := r.contains(test.input); res != test.expects {
			t.Errorf("Is target date within time range? Got %t, Expected %t", res, test.expects)
		}
	}
}
