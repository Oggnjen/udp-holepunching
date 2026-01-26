package logic

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Oggnjen/udp-holepunching/rendezvous-server/types"
)

func (s *Server) StartHttpServer() {
	http.HandleFunc("/get-udp-port", s.getUdpPortHandler)
	http.HandleFunc("/start-communication", s.getClientHandler)

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

func (s *Server) getClientHandler(w http.ResponseWriter, r *http.Request) {
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

		// s.Conn.WriteToUDPAddrPort()

	}
}
