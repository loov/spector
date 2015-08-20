// +build js

package main

import "github.com/egonelbre/spector/trace"

func main() {
	enc := trace.NewEncoder()
	(&trace.StreamStart{}).Encode(enc)

	println(enc.Data)
}
