package components

import "time"

type Clock struct {
	BaseClock *time.Ticker
}

func (c *Clock) Start() {
	c.BaseClock = time.NewTicker(100 * time.Millisecond)
}
