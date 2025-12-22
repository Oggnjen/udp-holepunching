package main

import (
	"github.com/Oggnjen/udp-holepunching/rendezvous-server/logic"
	"github.com/Oggnjen/udp-holepunching/rendezvous-server/types"
)

func main() {
	server := createServer()
	go server.StartUdp()
	server.StartHttpServer()
}

func createServer() *logic.Server {
	Clients := make(map[string]types.IPAddressPair)
	return &logic.Server{
		Clients: Clients,
	}
}
