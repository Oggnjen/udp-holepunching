package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/Oggnjen/udp-holepunching/client/logic"
	"github.com/Oggnjen/udp-holepunching/client/types"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	client := NewClient()
	client.FindServerPort()
	go client.StartUdp()

	reader := bufio.NewReader(os.Stdin)

	currentOption := "0"
	for {
		if client.Identifier != "" && currentOption != "C" && client.PeerAddress.Public == "" {
			fmt.Printf("Your identifier is: %s\n", client.Identifier)
			fmt.Println("Send your identifier to someone you want to chat with!")
		}

		if currentOption == "0" {
			fmt.Println("Choose option:")
			if client.Identifier == "" {
				fmt.Println("1. register")
			}
			fmt.Println("2. start chatting with someone")
			fmt.Print("-> ")
			text, _ := reader.ReadString('\n')
			currentOption = strings.TrimSpace(text)
			continue
		}

		if strings.Compare(currentOption, "1") == 0 {
			client.Register()
			currentOption = "X"
		}

		if strings.Compare(currentOption, "2") == 0 {
			fmt.Print("insert client identifier -> ")
			text, _ := reader.ReadString('\n')
			clientIdentifier := strings.TrimSpace(text)
			client.StartCommunication(clientIdentifier)
			fmt.Println("Starting chatting session...")
			client.StartChatting()
			currentOption = "C"
		}

		if strings.Compare(currentOption, "X") == 0 {
			fmt.Println("Press any key to continue...")
			reader.ReadString('\n')
			currentOption = "0"
		}

		if strings.Compare(currentOption, "C") == 0 || client.PeerContactSuccess {
			text, _ := reader.ReadString('\n')
			client.SendMessage(text)
		}

	}

}

func NewClient() *logic.Client {
	serverUrl := os.Getenv("SERVER_URL")
	if serverUrl == "" {
		fmt.Println("SERVER_URL not set, using default: localhost")
		serverUrl = "localhost"
	} else {
		fmt.Println("SERVER_URL:", serverUrl)
	}
	connTemp, err := net.ListenUDP("udp", &net.UDPAddr{Port: 0})
	if err != nil {
		panic("Cannot open UDP port")
	}
	return &logic.Client{
		Conn:           connTemp,
		Message:        make(chan types.Message, 10),
		ContactSuccess: false,
		ServerUrl:      serverUrl,
	}
}
