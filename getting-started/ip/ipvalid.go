package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "provide an ip address to parse!")
		os.Exit(-1)
	}

	ip := net.ParseIP(os.Args[1])
	if ip != nil {
		fmt.Fprintln(os.Stdout, "OK!")
		os.Exit(0)
	}
	fmt.Fprintln(os.Stderr, "bad ip address!")
}
