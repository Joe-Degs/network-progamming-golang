package main

import (
	"flag"
	"fmt"
	"math"
	"net"
	"os"
)

var (
	cidr string
)

func init() {
	flag.StringVar(&cidr, "c", "", "the CIDR address")
}

// this program  implements a CIDR subnet calculator
// it takes a CIDR addres and calculates the range, total hosts,
// wildcard mask, etc.
func main() {
	flag.Parse()

	if cidr == "" {
		fmt.Fprintln(os.Stderr, "CIDR address missing")
		os.Exit(1)
	}

	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Given IPv4 block 192.168.100.14/24
	// The followings use the IPNet to get:
	// - the routing address for the subnet
	// - one-bits of the network mask (24 out of 32 total)
	// - the subnet mask
	// - total hsots in the network (2 ^ (host identifier bits))
	// wildcard, inverse of subnet mask
	// max address of subnet mask
	ones, totalBits := ipnet.Mask.Size()
	size := totalBits - ones
	totalHosts := math.Pow(2, float64(size))
	wildcardIP := wildcard(net.IP(ipnet.Mask))
	last := lastIP(ip, net.IPMask(wildcardIP))

	fmt.Println()
	fmt.Printf("CIDR: %s\n", cidr)
	fmt.Println("------------------------------------------------")
	fmt.Printf("CIDR Block:         %s\n", cidr)
	fmt.Printf("Network:            %s\n", ipnet.IP)
	fmt.Printf("IP Range:           %s - %s\n", ip, last)
	fmt.Printf("Total Hosts:        %0.0f\n", totalHosts)
	fmt.Printf("Netmask:            %s\n", net.IP(ipnet.Mask))
	fmt.Printf("Wildcard Mask:      %s\n", wildcardIP)
	fmt.Println()
}

// wildcard returns the opposite of the netmask
// the netmask for the network.
func wildcard(mask net.IP) net.IP {
	var ipVal net.IP

	// loops through all the octets and reverses them
	for _, octet := range mask {
		ipVal = append(ipVal, ^octet)
	}
	return ipVal
}

// lastIP calculates the highest address range
// starting at the given IP
func lastIP(ip net.IP, mask net.IPMask) net.IP {
	ipIn := ip.To4()
	if ipIn == nil {
		ipIn = ip.To16()
		if ipIn == nil {
			return nil
		}
	}
	var ipVal net.IP

	// apply network mask to each octet
	for i, octet := range ipIn {
		ipVal = append(ipVal, octet|mask[i])
	}
	return ipVal
}
