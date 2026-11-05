package server

import (
	"net"
	"time"
)

func ReadPacket(listenerConn *net.UDPConn) ([]byte, *net.UDPAddr, error) {
	// allocate memory for the incoming packet
	requestPacketBuffer := make([]byte, 4096)

	// read the request from the client
	n, clientAddr, err := listenerConn.ReadFromUDP(requestPacketBuffer)
	if err != nil {
		return nil, nil, err
	}

	return requestPacketBuffer[:n], clientAddr, nil
}

func ForwardPacket(request []byte) ([]byte, error) {
	dnsConn, err := net.Dial("udp", "1.1.1.1:53")
	if err != nil {
		return nil, err
	}
	defer dnsConn.Close()
	dnsConn.SetReadDeadline(time.Now().Add(2 * time.Second))

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
