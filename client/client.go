package main

import (
	"flag"
)

const serverAddr = "127.0.0.1:5000"

func main() {
	hostAddr := flag.String("host", "127.0.0.1:5000", "server ip addr")
	t := Tcp{
		ServAddr: *hostAddr,
	}
	t.serve()
}
