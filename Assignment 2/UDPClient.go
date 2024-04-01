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

const SERVER_NAME = "nsl2.cau.ac.kr"
const SERVER_PORT = "12000"

func main() {
	conn := makeConnection()

	if conn == nil {
		return
	}

	printMenu()

	fmt.Printf("Input lowercase sentence: ")
	input, _ := bufio.NewReader(os.Stdin).ReadString('\n')

	server_addr, _ := net.ResolveUDPAddr("udp", SERVER_NAME+":"+SERVER_PORT)
	conn.WriteTo([]byte(input), server_addr)

	buffer := make([]byte, 1024)
	conn.ReadFrom(buffer)
	fmt.Printf("Reply from server: %s", string(buffer))

	closeConnection(conn)
}

func makeConnection() net.PacketConn {
	conn, err := net.ListenPacket("udp", ":")

	if err != nil {
		fmt.Printf("Error: %s", err.Error())
		return nil
	}

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	fmt.Printf("Client is running on port %d\n", localAddr.Port)

	return conn
}

func closeConnection(conn net.PacketConn) {
	if conn != nil {
		err := conn.Close()
		if err != nil {
			fmt.Printf("Error: %s", err.Error())
		}
	}
}

func printMenu() {
	fmt.Println("1) Convert Text to UPPER-case Letters")
	fmt.Println("2) Get Server Uptime")
	fmt.Println("3) Get Client IP / Port")
	fmt.Println("4) Get Count of Requests Server Got")
	fmt.Println("5) Exit Client")
}
