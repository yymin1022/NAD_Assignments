/**
 * UDPServer.go
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

const UDP_SERVER_PORT = "14094"

var serverResponseCnt int
var serverStartTime time.Time

func main() {
	serverConnection := initServer()
	if serverConnection == nil {
		printError("Failed to Init Server")
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
		if requestAddr != nil {
			fmt.Printf("UDP Connection Request from %s\n", requestAddr.String())

			cmd, _ := strconv.Atoi(string(requestBuffer[0]))
			fmt.Printf("Command %d\n", cmd)

			responseData := getResponse(cmd, string(requestBuffer[1:count]), requestAddr.String())
			_, err := serverConnection.WriteTo([]byte(responseData), requestAddr)
			if err != nil {
				printError("Failed to Send Response")
				continue
			}
			serverResponseCnt++
		}
	}
}

func initServer() net.PacketConn {
	serverConnection, err := net.ListenPacket("udp", ":"+UDP_SERVER_PORT)
	if err != nil {
		printError(err.Error())
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
		curTime := time.Now()
		upTime := int(curTime.Sub(serverStartTime).Seconds())
		upTimeH := upTime / 3600
		upTime %= 3600
		upTimeM := upTime / 60
		upTime %= 60
		return fmt.Sprintf("run time = %02d:%02d:%02d", upTimeH, upTimeM, upTime)
	case 3:
		addrInfo := strings.Split(addr, ":")
		return fmt.Sprintf("client IP = %s, port = %s", addrInfo[0], addrInfo[1])
	case 4:
		return fmt.Sprintf("requests served = %d", serverResponseCnt)
	}
	return ""
}

func printError(msg string) {
	_, err := fmt.Fprintf(os.Stderr, "Error: %s\n", msg)
	if err != nil {
		fmt.Printf("Error: %s\n", msg)
		return
	}
}
