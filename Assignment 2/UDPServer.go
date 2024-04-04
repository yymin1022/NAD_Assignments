/**
 * UDPServer.go
 **/

package main

import (
	"fmt"
	"net"
	"strings"
)

const UDP_SERVER_PORT = "14094"

func main() {
	serverConnection := initServer()
	if serverConnection == nil {
		fmt.Println("Error: Failed to Init Server")
		return
	}

	requestBuffer := make([]byte, 1024)

	for {
		count, requestAddr, _ := serverConnection.ReadFrom(requestBuffer)
		fmt.Printf("UDP Connection Request from %s\n", requestAddr.String())

		responseData := getResponse(int(requestBuffer[0]), string(requestBuffer[1:count]))
		_, err := serverConnection.WriteTo([]byte(responseData), requestAddr)
		if err != nil {
			fmt.Println("Error: Failed to Send Response")
			continue
		}
	}
}

func initServer() net.PacketConn {
	serverConnection, err := net.ListenPacket("udp", ":"+UDP_SERVER_PORT)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return nil
	}

	fmt.Printf("Server is ready to receive on port %s\n", UDP_SERVER_PORT)

	return serverConnection
}

func getResponse(cmd int, data string) string {
	fmt.Printf("cmd is %d, data is %s\n", cmd, data)
	return strings.ToUpper(data)
}
