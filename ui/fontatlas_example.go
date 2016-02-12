// +build ignore

package main

import (
	"flag"
	"image/png"
	"log"
	"os"

	"github.com/egonelbre/spector/ui/font"
)

func main() {
	flag.Parse()

	atlas, err := font.NewAtlas("~DejaVuSansMono.ttf", 72, 64)
	if err != nil {
		log.Fatal(err)
	}

	atlas.LoadExtendedAscii()
	atlas.DrawDebug()

	result, err := os.Create("~out.png")
	if err != nil {
		log.Fatal(err)
	}
	defer result.Close()

	err = png.Encode(result, atlas.Image)
	if err != nil {
		log.Fatal(err)
	}
}
