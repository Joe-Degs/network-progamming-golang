package main

import (
	"log"
	"net"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", ":1123")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := net.ListenUDP("udp", addr)
		if err != nil {
			log.Fatal(err)
		}
		go handleConn(conn)
	}

}

func handleConn(conn *net.UDPConn) {
	defer conn.Close()
	_, raddr, err := conn.ReadFromUDP(make([]byte, 48))
	if err != nil {
		log.Fatal(err)
	}
	if _, err := conn.WriteToUDP([]byte("A string message"), raddr); err != nil {
		log.Fatal(err)
	}
}
