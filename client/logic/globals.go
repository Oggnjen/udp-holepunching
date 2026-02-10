package logic

import (
	"net"

	"github.com/Oggnjen/udp-holepunching/client/types"
)

// var Port int

// var ServerPort string

// // var Conn net.PacketConn

// var MyIdentifier string

// var PeerIpAddress types.IPAddressPair

// var Messages chan types.Message

type Client struct {
	Port               int
	Conn               net.PacketConn
	Message            chan types.Message
	Identifier         string
	PeerAddress        types.IPAddressPair
	ServerUdpPort      string
	ContactSuccess     bool
	PeerContactSuccess bool
	ServerUrl          string
	ChatInitiated      bool
}
