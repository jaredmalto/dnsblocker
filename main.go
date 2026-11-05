package main

import (
	"fmt"
	"log"
	"net"

	"github.com/jaredmalto/dnsblocker/internal/server"
)

func main() {
	// open the udp socket
	// create the address to receive
	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:5300")
	if err != nil {
		log.Fatalln(err)
	}
	// establish the connection to listen on
	listener, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Println(err)
	}
	defer listener.Close()

	fmt.Println("listening on ", listener.LocalAddr())

	for {

		requestPacket, clientAddr, err := server.ReadPacket(listener)
		if err != nil {
			log.Println("Error receiving packet from client", err)
		}

		responsePacket, err := server.ForwardPacket(requestPacket)
		if err != nil {
			log.Println("Error sending packet upstream")
		}
		// write to the client address
		_, err = listener.WriteTo(responsePacket, clientAddr)
		if err != nil {
			log.Println("Error sending packet to client: ", err)
			continue
		}

	}

}
