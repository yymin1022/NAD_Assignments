/*
 * ChatClient.go
 * 20194094
 * Yongmin Yoo
 */

package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

const SERVER_NAME = "localhost"
const SERVER_PORT = "14094"

var PING_MODE = false
var PING_START_TIME = time.Now()

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run ChatClient.go <nickname>")
		os.Exit(1)
	}

	clientNickname := os.Args[1]

	serverConnection := makeConnection()
	if serverConnection == nil {
		printError("Failed to connect server.")
		os.Exit(1)
	}

	sigintHandler := make(chan os.Signal, 1)
	signal.Notify(sigintHandler, syscall.SIGINT)
	go func() {
		<-sigintHandler
		closeConnection(serverConnection)
		os.Exit(0)
	}()

	defer closeConnection(serverConnection)

	fmt.Fprintln(serverConnection, clientNickname)
	go readMessages(serverConnection)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		message := scanner.Text()
		if strings.HasPrefix(message, "\\") {
			cmd, extra := encodeCommand(message)
			if cmd == "P" {
				PING_MODE = true
				PING_START_TIME = time.Now()
				fmt.Fprintf(serverConnection, "P ping\n")
			} else if cmd != "" {
				fmt.Fprintf(serverConnection, "%s %s\n", cmd, extra)
			} else {
				fmt.Println("Invalid Command.")
			}
			if cmd == "Q" {
				closeConnection(serverConnection)
				os.Exit(0)
			}
		} else {
			fmt.Fprintln(serverConnection, "M"+message)
		}
	}
}

func makeConnection() net.Conn {
	conn, err := net.Dial("tcp4", SERVER_NAME+":"+SERVER_PORT)
	if err != nil {
		printError(err.Error())
		return nil
	}

	localAddr := conn.LocalAddr().(*net.TCPAddr)
	fmt.Printf("Client is running on port %d\n", localAddr.Port)

	return conn
}

func closeConnection(conn net.Conn) {
	fmt.Println("\rClosing Client Program...\nBye bye~")
	if conn != nil {
		_, err := conn.Write([]byte(string(rune(5))))
		if err != nil {
			printError(err.Error())
		}
		_ = conn.Close()
	}
}

func readMessages(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		readData := scanner.Text()
		if PING_MODE == true {
			PING_MODE = false
			PING_END_TIME := time.Now()
			rttValue := PING_END_TIME.Sub(PING_START_TIME).Nanoseconds() / 1000
			fmt.Printf("RTT is %vms\n", rttValue)
		} else {
			fmt.Println(readData[1:])
		}

		if readData[0] == 'K' {
			closeConnection(conn)
			os.Exit(0)
		}
	}
}

func encodeCommand(command string) (string, string) {
	parts := strings.SplitN(command, " ", 2)
	switch parts[0] {
	case "\\ls":
		return "L", ""
	case "\\ping":
		return "P", ""
	case "\\quit":
		return "Q", ""
	case "\\secret", "\\except":
		if len(parts) > 1 {
			return strings.ToUpper(string(parts[0][1])), parts[1]
		}
	}
	return "", ""
}

func printError(msg string) {
	_, err := fmt.Fprintf(os.Stderr, "Error: %s\n", msg)
	if err != nil {
		fmt.Printf("Error: %s\n", msg)
		return
	}
}
