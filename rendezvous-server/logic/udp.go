package logic

import (
	"encoding/json"
	"fmt"
	"log"
	"net"

	"github.com/Oggnjen/udp-holepunching/rendezvous-server/types"
	"github.com/google/uuid"
)

func (s *Server) StartUdp() {

	conn, err := net.ListenUDP("udp", &net.UDPAddr{Port: 0})
	if err != nil {
		panic("Cannot open UDP port")
	}
	defer conn.Close()

	udpAddr, ok := conn.LocalAddr().(*net.UDPAddr)
	if !ok {
		log.Fatal("Local address is not a UDP address")
	}

	fmt.Printf("Assigned port: %d\n", udpAddr.Port)
	s.Port = udpAddr.Port
	buffer := make([]byte, 1024)

	for {
		n, remoteAddr, err := conn.ReadFrom(buffer)

		if err != nil {
			fmt.Printf("Error occured during reading from UDP %v\n", err)
			continue
		}

		data := string(buffer[:n])

		fmt.Printf("Received %d bytes from %s: %s\n", n, remoteAddr.String(), data)

		identifier := uuid.New().String()
		s.Clients[identifier] = types.IPAddressPair{
			Private: "",
			Public:  remoteAddr.String(),
		}

		message := types.Message{
			Type:    "REGISTRATION",
			Payload: identifier,
		}
		response, err := json.Marshal(message)
		if err != nil {
			fmt.Printf("Error occured during sending response to client: %v\n", err)
			continue
		}
		_, err = conn.WriteTo(response, remoteAddr)

		if err != nil {
			fmt.Printf("Error occured during sending response to client: %v\n", err)
		}
	}
}
