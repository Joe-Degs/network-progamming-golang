package main

import (
	"flag"
	"fmt"
	"net"
	"os"
)

func main() {
	var addr string // the ip address we want the name for
	flag.StringVar(&addr, "addr", "::1", "provide address")
	flag.Parse()

	names, err := net.LookupAddr(addr)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println(names)
}
