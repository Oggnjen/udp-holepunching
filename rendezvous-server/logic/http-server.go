package logic

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/netip"

	"github.com/Oggnjen/udp-holepunching/rendezvous-server/types"
)

func (s *Server) StartHttpServer() {
	http.HandleFunc("/get-udp-port", s.getUdpPortHandler)
	http.HandleFunc("/start-communication", s.startCommunicationHandler)

	if err := http.ListenAndServe("0.0.0.0:8088", nil); err != nil {
		log.Fatal(err)
	}
}

func (s *Server) getUdpPortHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {

		w.Header().Set("Content-Type", "application/json")
		response := types.PortResponse{
			Port: s.Port,
		}
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(response)
	}
}

func (s *Server) startCommunicationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		defer r.Body.Close()
		w.Header().Set("Content-Type", "application/json")

		var startCommunication types.StartCommunication

		err := json.NewDecoder(r.Body).Decode(&startCommunication)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		clientAddr, ok := s.Clients[startCommunication.Peer]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(clientAddr)

		clientInitiatorAddr, ok := s.Clients[startCommunication.PeerInitiator]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		adrr, err := netip.ParseAddrPort(clientAddr.Public)

		mess := types.Message{
			Type:    "INITIATE_CHAT",
			Payload: clientInitiatorAddr.Public,
		}

		messBytes, err := json.Marshal(mess)
		if err != nil {
			fmt.Println("Error occured during sending udp message on: " + clientAddr.Public)
			return
		}

		fmt.Println("SENDING UDP ON" + adrr.String())
		_, err = s.Conn.WriteToUDPAddrPort(messBytes, adrr)

		if err != nil {
			fmt.Println("Error occured during sending udp message on: " + clientAddr.Public)
			return
		}

	}
}
