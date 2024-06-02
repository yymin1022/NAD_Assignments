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
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
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

	serverListener := initServer()

	sigintHandler := make(chan os.Signal, 1)
	signal.Notify(sigintHandler, syscall.SIGINT)
	go func() {
		<-sigintHandler
		closeServer(serverListener)
		os.Exit(0)
	}()

	for {
		if serverListener != nil {
			conn, _ := serverListener.Accept()
			go handleConnection(conn)
		}
	}
}

func closeServer(serverConn net.Listener) {
	fmt.Println("\rClosing Server Program...\nBye bye~")
	serverConn.Close()
}

func initServer() net.Listener {
	serverListener, err := net.Listen("tcp4", fmt.Sprintf(":%d", serverPort))
	if err != nil {
		exitError(err.Error())
	}
	fmt.Println("Server listening on port", serverPort)

	return serverListener
}

func handleConnection(conn net.Conn) {
	defer func() {
		if conn != nil {
			conn.Close()
		}
	}()

	readBuffer := make([]byte, 1024)
	for {
		if conn == nil {
			return
		}
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
	partFilename := getPartedFilename(filename)
	partFile, err := os.Create(partFilename)
	if err != nil {
		fmt.Println("Error creating file:", err.Error())
		conn.Write([]byte("ERROR\n"))
		return
	}
	defer partFile.Close()

	partFileBuffer := make([]byte, 1024)
	for {
		partFileLength, err := conn.Read(partFileBuffer)
		if err != nil && err != io.EOF {
			fmt.Println("Error reading:", err.Error())
			conn.Write([]byte("ERROR\n"))
			return
		}
		if partFileLength == 0 {
			return
		}
		if _, err := partFile.Write(partFileBuffer[:partFileLength]); err != nil {
			fmt.Println("Error writing to file:", err.Error())
			conn.Write([]byte("ERROR\n"))
			return
		}
	}
}

func sendHalfFile(conn net.Conn, filename string) {
	partFilename := getPartedFilename(filename)
	partFile, err := os.Open(partFilename)
	if err != nil {
		fmt.Println("Error opening file:", err.Error())
		conn.Write([]byte("ERROR\n"))
		return
	}
	defer partFile.Close()

	partFileBuffer := make([]byte, 1024)
	for {
		partFileLength, err := partFile.Read(partFileBuffer)
		if err != nil {
			if err == io.EOF {
				break
			}

			fmt.Println("Error reading file:", err.Error())
			conn.Write([]byte("ERROR\n"))
			return
		}
		if partFileLength == 0 {
			break
		}

		conn.Write(append([]byte("N"), partFileBuffer[:partFileLength]...))
	}

	conn.Write([]byte("EOF"))
}

func getPartedFilename(filename string) string {
	ext := filepath.Ext(filename)
	base := strings.TrimSuffix(filename, ext)
	return base + filenameSuffix + ext
}

func exitError(msg string) {
	_, err := fmt.Fprintf(os.Stderr, "Error: %s\n", msg)
	if err != nil {
		fmt.Printf("Error: %s\n", msg)
	}
	os.Exit(1)
}
