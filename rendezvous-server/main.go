package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/Oggnjen/udp-holepunching/rendezvous-server/types"
	"github.com/google/uuid"
)

var clients map[string]types.IPAddressPair

var port int

func main() {
	clients = make(map[string]types.IPAddressPair)
	go udp()
	http.HandleFunc("/get-udp-port", getUdpPortHandler)
	http.HandleFunc("/get-client", getClientHandler)

	if err := http.ListenAndServe("0.0.0.0:8088", nil); err != nil {
		log.Fatal(err)
	}

	udp()
}

func getUdpPortHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {

		w.Header().Set("Content-Type", "application/json")
		response := types.PortResponse{
			Port: port,
		}
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(response)
	}
}

func getClientHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Header().Set("Content-Type", "application/json")

		clientIdentifier := r.URL.Query().Get("identifier")

		clientAddr, ok := clients[clientIdentifier]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(clientAddr)

	}
}

func udp() {

	conn, err := net.ListenPacket("udp", ":0")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	udpAddr, ok := conn.LocalAddr().(*net.UDPAddr)
	if !ok {
		log.Fatal("Local address is not a UDP address")
	}

	fmt.Printf("Assigned port: %d\n", udpAddr.Port)
	port = udpAddr.Port
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
		clients[identifier] = types.IPAddressPair{
			Private: "",
			Public:  remoteAddr.String(),
		}
		_, err = conn.WriteTo([]byte(identifier), remoteAddr)

		if err != nil {
			fmt.Printf("Error occured during sending response to client: %v\n", err)
		}
	}
}
