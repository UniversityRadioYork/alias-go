package utils

import (
	"testing"
	"time"
)

func TestUtils_IsHistoricalOfficerValid(t *testing.T) {

	to, _ := time.Parse("2006/01/02", "2016/01/01")
	now1, _ := time.Parse("2006/01/02", "2016/05/01")
	now2, _ := time.Parse("2006/01/02", "2016/01/01")

	c := Config{
		configData: configData{
			StandDownPeriod: 28,
		},
	}

	b, err := c.IsHistoricalOfficerValid(
		now1,
		to,
	)

	if err != nil {
		t.Error(err)
	}

	if b {
		t.Error("Failed #1")
	}

	b, err = c.IsHistoricalOfficerValid(
		now2,
		to,
	)

	if err != nil {
		t.Error(err)
	}

	if !b {
		t.Error("Failed #2")
	}

}
