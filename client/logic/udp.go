package logic

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/Oggnjen/udp-holepunching/client/types"
)

func (c *Client) StartUdp() {

	defer c.Conn.Close()

	udpAddr, ok := c.Conn.LocalAddr().(*net.UDPAddr)
	if !ok {
		log.Fatal("Local address is not a UDP address")
	}

	c.Port = udpAddr.Port
	buffer := make([]byte, 1024)

	for {
		n, remoteAddr, err := c.Conn.ReadFrom(buffer)

		if err != nil {
			fmt.Printf("Error occured during reading from UDP %v\n", err)
			continue
		}

		data := buffer[:n]
		var message types.Message
		err = json.Unmarshal(data, &message)

		if err != nil {
			fmt.Printf("Error occured during reading from UDP %v\n", err)
			continue
		}

		if message.Type == "REGISTRATION" {
			c.Identifier = message.Payload
			fmt.Println("Your identifier is: " + c.Identifier)
		}

		if !c.ContactSuccess {
			var message types.Message
			err := json.Unmarshal(data, &message)
			if err != nil {
				fmt.Println("Error occured during unmarshaling request data", err)
				continue
			}
			if message.Type == "HELLO" {
				c.ContactSuccess = true
				response := types.Message{
					Type:    "SUCCESS",
					Payload: "...",
				}

				res, err := json.Marshal(response)
				if err != nil {
					fmt.Printf("Error occured during sending response to client: %v\n", err)
					continue
				}

				_, err = c.Conn.WriteTo(res, remoteAddr)

				if err != nil {
					fmt.Printf("Error occured during sending response to client: %v\n", err)
				}
			}
		}
		if message.Type != "TEXT" {
			if message.Type != "HEART-BEAT" {
				c.Message <- message
			}
		} else {
			fmt.Println(message.Payload)
		}
	}
}

func (c *Client) Register() {
	mess := types.IPAddressPort{
		AddressPort: "127.0.0.1: " + strconv.Itoa(c.Port),
	}
	message, err := json.Marshal(mess)
	if err != nil {
		fmt.Printf("Invalid message")
		return
	}
	remoteAddr, _ := net.ResolveUDPAddr("udp", "127.0.0.1:"+string(c.ServerUdpPort))

	_, err = c.Conn.WriteTo(message, remoteAddr)

	if err != nil {
		fmt.Printf("Error sending message: %v\n", err)
		return
	}
}

func (c *Client) StartChatting() {
	interval := 3 * time.Second
	ticker := time.NewTicker(interval)

	defer ticker.Stop()

	// punch a hole
	for {
		select {
		case mess := <-c.Message:
			if mess.Type == "SUCCESS" {
				c.PeerContactSuccess = true
				fmt.Println("Successfull contact!")
				go c.startHeartBeat()
				return
			}
		case t := <-ticker.C:
			mess := types.Message{
				Type:    "HELLO",
				Payload: "...",
			}
			message, err := json.Marshal(mess)
			if err != nil {
				fmt.Printf("Invalid message")
				return
			}
			remoteAddr, _ := net.ResolveUDPAddr("udp", c.PeerAddress.Public)

			_, err = c.Conn.WriteTo(message, remoteAddr)

			if err != nil {
				fmt.Printf("Error sending message: %v\n", err)
				return
			}
			fmt.Println("Next trying next hole punching at:", t.Add(interval))
		}
	}
}

func (c *Client) startHeartBeat() {
	interval := 3 * time.Second
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Resolve once
	remoteAddr, err := net.ResolveUDPAddr("udp", c.PeerAddress.Public)
	if err != nil {
		fmt.Printf("Failed to resolve peer address %q: %v\n", c.PeerAddress.Public, err)
		return
	}
	for {
		select {

		case <-ticker.C:
			mess := types.Message{
				Type:    "HEART-BEAT",
				Payload: "...",
			}
			message, err := json.Marshal(mess)
			if err != nil {
				fmt.Printf("Failed to marshal heartbeat: %v\n", err)
				return
			}

			_, err = c.Conn.WriteTo(message, remoteAddr)
			if err != nil {
				fmt.Printf("Error sending heartbeat: %v\n", err)
				return
			}
		}
	}
}

func (c *Client) SendMessage(payload string) {
	mess := types.Message{
		Type:    "TEXT",
		Payload: payload,
	}
	message, err := json.Marshal(mess)
	if err != nil {
		fmt.Printf("Invalid message")
		return
	}
	remoteAddr, _ := net.ResolveUDPAddr("udp", c.PeerAddress.Public)

	_, err = c.Conn.WriteTo(message, remoteAddr)

	if err != nil {
		fmt.Printf("Error sending message: %v\n", err)
		return
	}
}
