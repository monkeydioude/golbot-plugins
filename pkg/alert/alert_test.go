package alert

import (
	"fmt"
	"testing"
	"time"
)

func TestICanGetDurationForACloseTime(t *testing.T) {
	goal := 900 * time.Second
	tn := time.Now()
	tn = tn.Add(goal)
	trial := fmt.Sprintf("%d:%d", tn.Hour(), tn.Minute())

	d := getDuration(trial)

	if d == 0 || d != goal {
		t.Errorf("Duration should have been == %f", goal.Seconds())
	}
}

func TestICanGetDurationForAFarTime(t *testing.T) {
	goal := (15 * time.Hour) + (55 * time.Minute)
	tn := time.Now()
	tn = tn.Add(goal)
	trial := fmt.Sprintf("%d:%d", tn.Hour(), tn.Minute())

	d := getDuration(trial)

	if d == 0 || d != goal {
		t.Errorf("Duration should have been == %f", goal.Seconds())
	}
}
