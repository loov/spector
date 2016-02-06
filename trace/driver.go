package trace

import (
	"sync/atomic"
	"time"
)

type Driver struct {
	interval time.Duration
	reader   Reader
	handler  Handler

	killed int32
}

func NewDriver(reader Reader, handler Handler, interval time.Duration) *Driver {
	return &Driver{
		interval: interval,
		reader:   reader,
		handler:  handler,
	}
}

func (driver *Driver) Start()      { go driver.Run() }
func (driver *Driver) Stop()       { atomic.StoreInt32(&driver.killed, 1) }
func (driver *Driver) Alive() bool { return atomic.LoadInt32(&driver.killed) == 0 }

func (driver *Driver) Run() {
	tick := time.NewTicker(driver.interval)
	defer tick.Stop()
	for range tick.C {
		if !driver.Alive() {
			return
		}

		for _, event := range driver.reader.Next() {
			driver.handler.Handle(event)
		}
	}
}
