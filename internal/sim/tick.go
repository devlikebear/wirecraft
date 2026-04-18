package sim

import "time"

type TickID uint64

func (id TickID) Next() TickID {
	return id + 1
}

type Clock struct {
	RateHz int
}

func (c Clock) TargetDuration() time.Duration {
	if c.RateHz <= 0 {
		return 0
	}
	return time.Second / time.Duration(c.RateHz)
}
