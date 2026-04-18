package sim

import (
	"testing"
	"time"
)

func TestClockTargetDuration(t *testing.T) {
	clock := Clock{RateHz: 20}

	if got := clock.TargetDuration(); got != 50*time.Millisecond {
		t.Fatalf("TargetDuration() = %s, want %s", got, 50*time.Millisecond)
	}
}

func TestClockRejectsInvalidRate(t *testing.T) {
	clock := Clock{RateHz: 0}

	if got := clock.TargetDuration(); got != 0 {
		t.Fatalf("TargetDuration() for invalid rate = %s, want 0", got)
	}
}

func TestTickIDNextIsMonotonic(t *testing.T) {
	var id TickID = 41

	next := id.Next()

	if next != 42 {
		t.Fatalf("Next() = %d, want 42", next)
	}
	if next <= id {
		t.Fatalf("Next() = %d, want greater than %d", next, id)
	}
}
