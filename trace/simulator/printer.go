package simulator

import (
	"fmt"

	"github.com/egonelbre/spector/trace"
)

type Printer struct{}

func NewPrinter() *Printer {
	return &Printer{}
}

func (printer *Printer) Handle(event trace.Event) {
	fmt.Println(event)
}
