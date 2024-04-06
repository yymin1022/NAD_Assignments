/**
 * UDPServer.go
 **/

package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

const UDP_SERVER_PORT = "14094"

var serverResponseCnt int
var serverStartTime time.Time

func main() {
	serverConnection := initServer()
	if serverConnection == nil {
		fmt.Println("Error: Failed to Init Server")
		return
	}

	sigintHandler := make(chan os.Signal, 1)
	signal.Notify(sigintHandler, syscall.SIGINT)
	go func() {
		<-sigintHandler
		closeServer(serverConnection)
		os.Exit(0)
	}()

	serverResponseCnt = 0
	serverStartTime = time.Now()

	requestBuffer := make([]byte, 1024)

	for {
		count, requestAddr, _ := serverConnection.ReadFrom(requestBuffer)
		fmt.Printf("UDP Connection Request from %s\n", requestAddr.String())

		responseData := getResponse(int(requestBuffer[0]), string(requestBuffer[1:count]), requestAddr.String())
		_, err := serverConnection.WriteTo([]byte(responseData), requestAddr)
		if err != nil {
			fmt.Println("Error: Failed to Send Response")
			continue
		}
		serverResponseCnt++
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

func closeServer(conn net.PacketConn) {
	fmt.Println("\rClosing Server Program...\nBye bye~")
	if conn != nil {
		_ = conn.Close()
	}
}

func getResponse(cmd int, data string, addr string) string {
	switch cmd {
	case 1:
		return strings.ToUpper(data)
	case 2:
		curTime := time.Since(serverStartTime)
		return fmt.Sprintf("run time = %02.0f:%02.0f:%02.0f", curTime.Hours(), curTime.Minutes(), curTime.Seconds())
	case 3:
		addrInfo := strings.Split(addr, ":")
		return fmt.Sprintf("client IP = %s, port = %s", addrInfo[0], addrInfo[1])
	case 4:
		return fmt.Sprintf("requests served = %d", serverResponseCnt)
	}
	return ""
}
