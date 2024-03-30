/**
 * UDPClient.go
 **/

package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	serverName := "nsl2.cau.ac.kr"
	serverPort := "12000"

	conn, _ := net.ListenPacket("udp", ":")

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	fmt.Printf("Client is running on port %d\n", localAddr.Port)

	fmt.Printf("Input lowercase sentence: ")
	input, _ := bufio.NewReader(os.Stdin).ReadString('\n')

	server_addr, _ := net.ResolveUDPAddr("udp", serverName+":"+serverPort)
	conn.WriteTo([]byte(input), server_addr)

	buffer := make([]byte, 1024)
	conn.ReadFrom(buffer)
	fmt.Printf("Reply from server: %s", string(buffer))

	conn.Close()
}
