/**
 * SplitFileServer.go
 * ID : 20194094
 * Name : Yongmin Yoo
 **/

package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

var serverPort int
var filenameSuffix string

func main() {
	if len(os.Args) != 2 {
		exitError("Usage: go run server.go <port>")
	}

	serverPortArgument, err := strconv.Atoi(os.Args[1])
	if err != nil {
		exitError("Invalid Argument")
	}

	serverPort = serverPortArgument
	if serverPort/10000 == 4 {
		filenameSuffix = "-part1"
	} else {
		filenameSuffix = "-part2"
	}

	initServer()
}

func initServer() {
	serverListener, err := net.Listen("tcp4", fmt.Sprintf(":%d", serverPort))
	if err != nil {
		exitError(err.Error())
	}
	defer serverListener.Close()
	fmt.Println("Server listening on port", serverPort)

	for {
		conn, err := serverListener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err.Error())
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	readBuffer := make([]byte, 1024)
	for {
		readLength, err := conn.Read(readBuffer)
		if err != nil && err != io.EOF {
			fmt.Println("Error reading:", err.Error())
			return
		}
		if readLength == 0 {
			return
		}

		readData := string(readBuffer[:readLength])
		readDataParts := strings.SplitN(readData, ":", 2)
		if len(readDataParts) != 2 {
			fmt.Println("Invalid message format")
			return
		}

		cmd, filename := readDataParts[0], readDataParts[1]
		if cmd == "PUT" {
			conn.Write([]byte("READY\n"))
			saveHalfFile(conn, filename)
		} else if cmd == "GET" {
			sendHalfFile(conn, filename)
		} else {
			fmt.Println("Unknown command")
			return
		}
	}
}

func saveHalfFile(conn net.Conn, filename string) {
	partFilename := filename + filenameSuffix
	partFile, err := os.Create(partFilename)
	if err != nil {
		fmt.Println("Error creating file:", err.Error())
		return
	}
	defer partFile.Close()

	partFileBuffer := make([]byte, 1024)
	for {
		partFileLength, err := conn.Read(partFileBuffer)
		if err != nil && err != io.EOF {
			fmt.Println("Error reading:", err.Error())
			return
		}
		if partFileLength == 0 {
			return
		}
		if _, err := partFile.Write(partFileBuffer[:partFileLength]); err != nil {
			fmt.Println("Error writing to file:", err.Error())
			return
		}
	}
}

func sendHalfFile(conn net.Conn, filename string) {
	partFilename := filename + filenameSuffix
	partFile, err := os.Open(partFilename)
	if err != nil {
		fmt.Println("Error opening file:", err.Error())
		return
	}
	defer partFile.Close()

	partFileBuffer := make([]byte, 1024)
	for {
		partFileLength, err := partFile.Read(partFileBuffer)
		if err != nil && err != io.EOF {
			fmt.Println("Error reading file:", err.Error())
			return
		}
		if partFileLength == 0 {
			break
		}

		conn.Write(partFileBuffer[:partFileLength])
	}

	conn.Write([]byte("EOF"))
}

func exitError(msg string) {
	_, err := fmt.Fprintf(os.Stderr, "Error: %s\n", msg)
	if err != nil {
		fmt.Printf("Error: %s\n", msg)
	}
	os.Exit(1)
}
