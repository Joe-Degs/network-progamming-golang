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

func main() {
	var host string
	flag.StringVar(&host, "e", ":1123", "server address")
	flag.Parse()

	// create a udp host address
	addr, err := net.ResolveUDPAddr("udp", host)
	if err != nil {
		throwErr(err)
	}

	// setup a udp listening connection
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		throwErr(err)
	}
	defer conn.Close()

	fmt.Printf("listening for time requests on (udp) %s\n", conn.LocalAddr())

	_, raddr, err := conn.ReadFromUDP(make([]byte, 48))
	if err != nil {
		throwErr(err)
	}
	if raddr == nil {
		throwErr(errors.New("missing the requester's address"))
	}

	// conjure ntp seconds from ntp epoch till now
	secs, fracs := getNTPSeconds(time.Now())
	resp := make([]byte, 48)
	// write seconds a uint32(4bytes) into resp [40:43]
	binary.BigEndian.PutUint32(resp[40:], uint32(secs))
	binary.BigEndian.PutUint32(resp[44:], uint32(fracs))

	if _, err := conn.WriteToUDP(resp, raddr); err != nil {
		throwErr(err)
	}
}

func getNTPOffset() float64 {
	//now_sec, now_frac := t.Second(), t.NanoSecond()
	ntpEpoch := time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
	unixEpoch := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	return unixEpoch.Sub(ntpEpoch).Seconds()
}

func getNTPSeconds(t time.Time) (int64, int64) {
	return t.Unix() + int64(getNTPOffset()), int64(t.Nanosecond())
}

func throwErr(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(0b00000001)
}
