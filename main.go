package main

import (
	"flag"
	"time"

	"github.com/egonelbre/spector/trace"
	"github.com/egonelbre/spector/trace/simulator"
)

func main() {
	flag.Parse()

	stream := simulator.NewStream()
	stream.Start()

	driver := trace.NewDriver(
		stream,
		simulator.NewPrinter(),
		time.Second,
	)
	driver.Run()
}
