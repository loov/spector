// +build js

package main

import "github.com/egonelbre/spector/trace"

func main() {
	enc := trace.NewEncoder()
	(&trace.StreamStart{
		ProcessID: 1562,
		MachineID: 412,

		Time: 1242,
		Freq: 10,
	}).Encode(enc)

	println(enc.Data)
}
