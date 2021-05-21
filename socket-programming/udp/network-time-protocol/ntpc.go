package main

import (
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"time"
)

// this program implements  a trivial network time protocol client over udp
// Its uses ntp version 3 data packet format which is 48 bytes long datagram
// for both request and response.
// Usage:
// ntpc -e <host endpoint>
func main() {
	var host string
	flag.StringVar(&host, "e", "us.pool.ntp.org:123", "NTP Host")
	flag.Parse()

	// req datagram is a 48-byte long slice
	// that is used in sending time request to the server
	req := make([]byte, 48)

	// req is initiealized with 0x1b which is a request setting
	// for the time server
	req[0] = 0x1b

	// response 48 byte slice for recieving ntp responses from the server
	resp := make([]byte, 48)

	raddr, err := net.ResolveUDPAddr("udp", host)
	if err != nil {
		throwErr(err)
	}

	// create a net.UDPConn connection.
	conn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		throwErr(err)
	}

	defer func() {
		if err := conn.Close(); err != nil {
			throwErr(err)
		}
	}()

	// send a time request
	if _, err = conn.Write(req); err != nil {
		throwErr(err)
	}

	//  recieve a fricken response
	println("reading now ...")
	// this never worked, don't know if its my terrible connection
	// or its the  ntp server thats slow to respond but its not cool.
	read, err := conn.Read(resp)
	println("done reading ...")
	if err != nil {
		throwErr(err)
	}
	if read != 48 {
		throwErr(errors.New("did not get all packets succesfully!"))
	}

	// ntp data comes in the big-endian format
	// with 64-bit value containing the server time in seconds
	// first 32 is the seconds and the rest 32 reps fractions
	// so 4 bytes and it starts from [40:]
	// this is the number of seconds since 1900 (ntp epoch)
	secs := binary.BigEndian.Uint32(resp[40:])
	fracs := binary.BigEndian.Uint32(resp[44:])

	// ntp epoch is 1900
	// unix epoch is 1970
	ntpEpoch := time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
	unixEpoch := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	offset := unixEpoch.Sub(ntpEpoch).Seconds()
	now := float64(secs) - offset
	fmt.Printf("%v\n", time.Unix(int64(now), int64(fracs)))
}

func throwErr(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(0b1)
}
