/**
 * UDPClient.go
 **/

package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

const SERVER_NAME = "localhost"
const SERVER_PORT = "14094"

func main() {
	serverConnection := makeConnection()

	if serverConnection == nil {
		printError("Failed to connect server.")
		return
	}

	sigintHandler := make(chan os.Signal, 1)
	signal.Notify(sigintHandler, syscall.SIGINT)
	go func() {
		<-sigintHandler
		closeConnection(serverConnection)
		os.Exit(0)
	}()

	exitFlag := false
	for !exitFlag {
		printMenu()
		cmd := readCommand()

		text := ""
		if cmd == 0 {
			continue
		} else if cmd == 5 {
			break
		} else if cmd == 1 {
			fmt.Printf("Input lowercase sentence: ")
			var err error
			text, err = bufio.NewReader(os.Stdin).ReadString('\n')

			if err != nil {
				printError(err.Error())
				continue
			}

			if len(text) >= 1024 {
				printError("Text too long.")
				continue
			}
		}

		timeRequest := time.Now().UnixMicro()

		serverAddr, _ := net.ResolveUDPAddr("udp4", SERVER_NAME+":"+SERVER_PORT)
		_, err := serverConnection.WriteTo([]byte(fmt.Sprintf("%d%s", cmd, text)), serverAddr)
		if err != nil {
			printError(err.Error())
			continue
		}

		responseBuffer := make([]byte, 1024)
		serverTimer := time.NewTimer(time.Second * 5)
		go func() {
			<-serverTimer.C
			exitFlag = true
			printError("Server Timeout.")
			closeConnection(serverConnection)
			os.Exit(0)
		}()

		_, _, err = serverConnection.ReadFrom(responseBuffer)

		if !exitFlag {
			if err != nil {
				printError(err.Error())
				serverTimer.Stop()
				continue
			}
			serverTimer.Stop()
			timeResponse := time.Now().UnixMicro()

			fmt.Printf("\nReply from server: %s\n", string(responseBuffer))
			fmt.Printf("RTT = %.3fms\n", float64(timeResponse-timeRequest)/1000)
		}
	}

	if !exitFlag {
		closeConnection(serverConnection)
	}
}

func makeConnection() net.PacketConn {
	conn, err := net.ListenPacket("udp", ":")

	if err != nil {
		printError(err.Error())
		return nil
	}

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	fmt.Printf("Client is running on port %d\n", localAddr.Port)

	return conn
}

func closeConnection(conn net.PacketConn) {
	fmt.Println("\rClosing Client Program...\nBye bye~")
	if conn != nil {
		_ = conn.Close()
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

	fmt.Print("Input Command: ")
	_, err := fmt.Scanln(&input)
	if err != nil {
		printError(err.Error())
	}

	cmd, err := strconv.ParseInt(input, 10, 0)
	if err != nil {
		printError("Invalid Command.")
		return 0
	}

	if cmd < 1 || cmd > 5 {
		printError("Invalid Command.")
		return 0
	}

	return int(cmd)
}

func printError(msg string) {
	_, err := fmt.Fprintf(os.Stderr, "Error: %s\n", msg)
	if err != nil {
		fmt.Printf("Error: %s\n", msg)
		return
	}
}
