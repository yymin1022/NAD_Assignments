/**
 * UDPServer.go
 **/

package main

import (
	"bytes"
	"fmt"
	"net"
)

const SERVER_PORT = "14094"

func main() {
	conn, _ := net.ListenPacket("udp", ":"+SERVER_PORT)
	fmt.Printf("Server is ready to receive on port %s\n", SERVER_PORT)

	buffer := make([]byte, 1024)

	for {
		count, r_addr, _ := conn.ReadFrom(buffer)
		fmt.Printf("UDP message from %s\n", r_addr.String())
		conn.WriteTo(bytes.ToUpper(buffer[:count]), r_addr)
	}
}
