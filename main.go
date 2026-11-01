package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

func readPacket(listenerConn *net.UDPConn) ([]byte, *net.UDPAddr, error) {
	// allocate memory for the incoming packet
	requestPacketBuffer := make([]byte, 4096)

	// read the request from the client
	n, clientAddr, err := listenerConn.ReadFromUDP(requestPacketBuffer)
	if err != nil {
		return nil, nil, err
	}

	return requestPacketBuffer[:n], clientAddr, nil
}

func forwardPacket(request []byte) ([]byte, error) {
	dnsConn, err := net.Dial("udp", "1.1.1.1:53")
	if err != nil {
		return nil, err
	}
	defer dnsConn.Close()
	dnsConn.SetReadDeadline(time.Now().Add(2 * time.Second))

	fmt.Printf("Sending payload: %v\n", request)

	_, err = dnsConn.Write(request)
	if err != nil {
		return nil, err
	}
	// read the response from the dns server
	responseBuffer := make([]byte, 4096)
	returnSize, err := dnsConn.Read(responseBuffer)
	if err != nil {
		return nil, err
	}
	return responseBuffer[:returnSize], nil
}

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

		requestPacket, clientAddr, err := readPacket(listener)
		if err != nil {
			log.Println("Error receiving packet from client", err)
		}

		responsePacket, err := forwardPacket(requestPacket)
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
