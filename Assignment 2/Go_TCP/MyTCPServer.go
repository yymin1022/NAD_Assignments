/**
 * TCPServer.go
 **/

package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const TCP_SERVER_PORT = "14094"

var serverResponseCnt int
var serverStartTime time.Time

func main() {
	serverListener := initServer()
	if serverListener == nil {
		fmt.Println("Error: Failed to Init Server")
		return
	}

	sigintHandler := make(chan os.Signal, 1)
	signal.Notify(sigintHandler, syscall.SIGINT)
	go func() {
		<-sigintHandler
		closeServer(serverListener)
		os.Exit(0)
	}()

	serverResponseCnt = 0
	serverStartTime = time.Now()

	requestBuffer := make([]byte, 1024)

	for {
		serverConnection, _ := serverListener.Accept()
		for serverConnection != nil {
			requestAddr := serverConnection.RemoteAddr()
			if requestAddr != nil {
				fmt.Printf("TCP Connection Request from %s\n", requestAddr.String())

				count, _ := serverConnection.Read(requestBuffer)
				cmd, _ := strconv.Atoi(string(requestBuffer[0]))
				fmt.Printf("Command %d\n", cmd)

				responseData := getResponse(cmd, string(requestBuffer[1:count]), requestAddr.String())

				if responseData == "" {
					_ = serverConnection.Close()
					break
				}
				_, err := serverConnection.Write([]byte(responseData))
				if err != nil {
					fmt.Println("Error: Failed to Send Response")
					continue
				}
				serverResponseCnt++
			}
		}
	}
}

func initServer() net.Listener {
	serverListener, err := net.Listen("tcp", ":"+TCP_SERVER_PORT)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return nil
	}

	fmt.Printf("Server is ready to receive on port %s\n", TCP_SERVER_PORT)

	return serverListener
}

func closeServer(listener net.Listener) {
	fmt.Println("\rClosing Server Program...\nBye bye~")
	if listener != nil {
		_ = listener.Close()
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
