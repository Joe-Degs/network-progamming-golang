package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
)

func main() {
	var host string

	flag.StringVar(&host, "host", "::1", "provide a host name")
	flag.Parse()

	var nameTable []struct {
		ip    string
		names []string
	}
	res := &net.Resolver{PreferGo: true}
	ctx := context.Background()

	addrs := resolveHost(ctx, res, host)
	if addrs == nil {
		os.Exit(1)
	}

	for _, addr := range addrs {
		n := resolveAddr(ctx, res, addr)
		if n == nil || n[0] == "" {
			nameTable = append(nameTable, struct {
				ip    string
				names []string
			}{ip: addr, names: nil})
			continue
		}
		nameTable = append(nameTable, struct {
			ip    string
			names []string
		}{ip: addr, names: n})
	}

	for _, name := range nameTable {
		fmt.Println("-------------------------------------------------")
		fmt.Printf("IP Address:       %v\n", name.ip)
		fmt.Printf("Resolved Names:   %v\n", name.names)
	}
	os.Exit(0)
}

func resolveHost(ctx context.Context, res *net.Resolver, host string) []string {
	addrs, err := res.LookupHost(ctx, host)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil
	}
	return addrs
}

func resolveAddr(ctx context.Context, res *net.Resolver, addr string) []string {
	names, err := res.LookupAddr(ctx, addr)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil
	}
	return names
}
