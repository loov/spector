package main

import (
	"flag"
	"net"
	"net/http"
	"os"
)

var (
	addr = flag.String("listen", ":8000", "address to listen on")
)

func main() {
	flag.Parse()
	if os.Getenv("HOST") != "" || os.Getenv("PORT") != "" {
		*addr = net.JoinHostPort(os.Getenv("HOST"), os.Getenv("PORT"))
	}

	http.ListenAndServe(*addr, http.FileServer(http.Dir(".")))
}
