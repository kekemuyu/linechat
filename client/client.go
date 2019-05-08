package main

import (
	"flag"
	"fmt"
)

const serverAddr = "127.0.0.1:5000"

func main() {
	hostAddr := flag.String("host", "kekemuyu.com:3000", "server ip addr")
	flag.Parse()
	t := &Tcp{
		ServAddr: *hostAddr,
	}
	fmt.Println(t.ServAddr)
	t.serve()
}
