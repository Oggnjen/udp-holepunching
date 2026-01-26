package logic

import (
	"net"

	"github.com/Oggnjen/udp-holepunching/rendezvous-server/types"
)

type Server struct {
	Clients map[string]types.IPAddressPair
	Conn    *net.UDPConn
	Port    int
}
