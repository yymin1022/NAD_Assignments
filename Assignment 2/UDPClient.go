/**
 * UDPClient.go
 **/

package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

const SERVER_NAME = "localhost"
const SERVER_PORT = "14094"

func main() {
	serverConnection := makeConnection()

	sigintHandler := make(chan os.Signal, 1)
	signal.Notify(sigintHandler, syscall.SIGINT)
	go func() {
		<-sigintHandler
		closeConnection(serverConnection)
		os.Exit(0)
	}()

	if serverConnection == nil {
		fmt.Println("Error: Failed to Connect")
		return
	}

	for {
		printMenu()
		cmd := readCommand()

		text := ""
		if cmd == 0 {
			continue
		} else if cmd == 5 {
			break
		} else if cmd == 1 {
			fmt.Printf("Input lowercase sentence: ")
			_, err := fmt.Scanf("%s", &text)

			if err != nil {
				fmt.Printf("Error: %s\n", err.Error())
				continue
			}
		}

		serverAddr, _ := net.ResolveUDPAddr("udp", SERVER_NAME+":"+SERVER_PORT)
		_, err := serverConnection.WriteTo([]byte(string(rune(cmd))+text), serverAddr)
		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
			continue
		}

		responseBuffer := make([]byte, 1024)
		_, _, err = serverConnection.ReadFrom(responseBuffer)
		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
			continue
		}
		fmt.Printf("Reply from server: %s\n", string(responseBuffer))
	}

	closeConnection(serverConnection)
}

func makeConnection() net.PacketConn {
	conn, err := net.ListenPacket("udp", ":")

	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return nil
	}

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	fmt.Printf("Client is running on port %d\n", localAddr.Port)

	return conn
}

func closeConnection(conn net.PacketConn) {
	fmt.Println("\rClosing Client Program...\nBye bye~")
	if conn != nil {
		err := conn.Close()
		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
		}
	}
}

func printMenu() {
	fmt.Println()
	fmt.Println("< Select Menu. >")
	fmt.Println("1) Convert Text to UPPER-case Letters")
	fmt.Println("2) Get Server Uptime")
	fmt.Println("3) Get Client IP / Port")
	fmt.Println("4) Get Count of Requests Server Got")
	fmt.Println("5) Exit Client")
}

func readCommand() int {
	var input string

	_, err := fmt.Scanln(&input)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}

	cmd, err := strconv.ParseInt(input, 10, 0)
	if err != nil {
		fmt.Println("Error: Invalid Command")
		return 0
	}

	if cmd < 1 || cmd > 5 {
		fmt.Println("Error: Invalid Command")
		return 0
	}

	return int(cmd)
}
