package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Oggnjen/udp-holepunching/client/types"
)

var port int

var serverPort string

var conn net.PacketConn

var myIdentifier string

func main() {
	findServerPort()
	go udp()

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Simple Shell")
	fmt.Println("---------------------")
	currentOption := "0"
	for {
		if currentOption == "0" {
			fmt.Println("Choose option:")
			fmt.Println("1. register")
			fmt.Println("2. get client data")
			fmt.Print("-> ")
			text, _ := reader.ReadString('\n')
			currentOption = strings.TrimSpace(text)
			fmt.Println(currentOption, "1")
			continue
		}

		if strings.Compare(currentOption, "1") == 0 {
			fmt.Println("USAO")
			register()
			currentOption = "0"
		}

	}

}

func register() {
	fmt.Println(serverPort)
	udpAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:"+string(serverPort))
	if err != nil {
		panic(err)
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// The message to send as a byte slice
	mess := types.IPAddressPort{
		AddressPort: "127.0.0.1: " + strconv.Itoa(port),
	}
	message, err := json.Marshal(mess)
	if err != nil {
		fmt.Printf("Invalid message")
		return
	}

	// Send the message
	_, err = conn.Write(message)
	if err != nil {
		fmt.Printf("Error sending message: %v\n", err)
		return
	}

	fmt.Printf("Message sent to %s: %s\n", conn.RemoteAddr().String(), string(message))

	// Optionally, read a response from the server
	buffer := make([]byte, 1024)
	// Set a deadline for reading to prevent hanging indefinitely
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	n, _, err := conn.ReadFromUDP(buffer)
	if err != nil {
		fmt.Printf("Error reading response (may be expected if server doesn't reply): %v\n", err)
	} else {
		myIdentifier = string(buffer[:n])
		fmt.Printf("Received response: %s\n", string(buffer[:n]))
	}
}

func findServerPort() {
	response, err := http.Get("http://localhost:8088/get-udp-port")
	if err != nil {
		fmt.Printf("Error making GET request: %s\n", err)
		os.Exit(1)
	}
	// It is crucial to close the response body to prevent resource leaks
	defer response.Body.Close()

	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Status Code: %d\n", response.StatusCode)
	fmt.Printf("Response Body: %s\n", string(body))

	var resPort types.PortResponse
	if err := json.Unmarshal(body, &resPort); err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}
	fmt.Println(strconv.Itoa(resPort.Port))
	serverPort = strconv.Itoa(resPort.Port)
}

func interfaceToStringBytes(data interface{}) ([]byte, bool) {
	if str, ok := data.(string); ok {
		return []byte(str), true
	}
	return nil, false
}

func udp() {

	connTemp, err := net.ListenPacket("udp", ":0")
	conn = connTemp
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

		response := "Message received at " + time.Now().Format(time.RFC3339) + "\n"

		_, err = conn.WriteTo([]byte(response), remoteAddr)

		if err != nil {
			fmt.Printf("Error occured during sending response to client: %v\n", err)
		}
	}
}
