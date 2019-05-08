package main

import (
	"flag"
)

func main() {
	port := flag.String("port", ":5000", "serve port")
	flag.Parse()
	t := &Tcp{
		ServAddr: *port,
	}
	t.serve()
}
