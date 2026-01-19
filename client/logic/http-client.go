package logic

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/Oggnjen/udp-holepunching/client/types"
)

func (c *Client) FindServerPort() {
	response, err := http.Get("http://157.90.241.190:8088/get-udp-port")
	if err != nil {
		fmt.Printf("Error making GET request: %s\n", err)
		os.Exit(1)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %s\n", err)
		os.Exit(1)
	}

	var resPort types.PortResponse
	if err := json.Unmarshal(body, &resPort); err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}
	c.ServerUdpPort = strconv.Itoa(resPort.Port)
	// ServerPort = strconv.Itoa(resPort.Port)
}

func (c *Client) GetClientData(clientId string) {
	response, err := http.Get("http://157.90.241.190:8088/get-client?identifier=" + clientId)
	if err != nil {
		fmt.Printf("Error making GET request: %s\n", err)
		os.Exit(1)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %s\n", err)
		os.Exit(1)
	}

	var peerResponse types.IPAddressPair
	if err := json.Unmarshal(body, &peerResponse); err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}
	c.PeerAddress = peerResponse
}
