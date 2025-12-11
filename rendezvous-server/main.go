package main

import (
	"net"

	"github.com/Oggnjen/udp-holepunching/rendezvous-server/types"
)

var clients map[string]types.IPAddressPair

func main() {
	clients = make(map[string]types.IPAddressPair)

	udpAddress, err := net.ResolveUDPAddr("udp", ":8088")

}
