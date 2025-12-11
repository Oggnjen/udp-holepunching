package main

import (
	"net"
)

func main() {

	pc, err := net.ListenPacket("udp", ":0")
	if err != nil {
		panic(err)
	}
	defer pc.Close()
	// pc.
}
