/**
 * ChatClient.go
 * ID : 20194094
 * Name : Yongmin Yoo
 **/

package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
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
		msg := readMessage()

		if msg == "\\quit" {
			break
		}

		msg = checkCommand(msg)
		print(msg)
		timeRequest := time.Now().UnixMicro()

		_, err := serverConnection.Write([]byte(msg))
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

		_, err = serverConnection.Read(responseBuffer)

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

func printMenu() {
	fmt.Println()
	fmt.Println("< Select Menu. >")
	fmt.Println("1) Convert Text to UPPER-case Letters")
	fmt.Println("2) Get Server Uptime")
	fmt.Println("3) Get Client IP / Port")
	fmt.Println("4) Get Count of Requests Server Got")
	fmt.Println("5) Exit Client")
}

func readMessage() string {
	var input string

	_, err := fmt.Scanln(&input)
	if err != nil {
		printError(err.Error())
	}

	return input
}

func checkCommand(message string) string {
	if len(message) > 2 && message[0:3] == "\\ls" {
		return "L"
	} else if len(message) > 8 && message[0:9] == "\\secret " {
		return "S" + message[8:len(message)-1]
	} else if len(message) > 8 && message[0:9] == "\\except " {
		return "E" + message[8:len(message)-1]
	} else if len(message) > 4 && message[0:5] == "\\ping" {
		return "P"
	}
	return "N" + message
}

func printError(msg string) {
	_, err := fmt.Fprintf(os.Stderr, "Error: %s\n", msg)
	if err != nil {
		fmt.Printf("Error: %s\n", msg)
		return
	}
}
