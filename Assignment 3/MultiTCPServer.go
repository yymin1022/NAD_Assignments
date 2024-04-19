/**
 * TCPServer.go
 * ID : 20194094
 * Name : Yongmin Yoo
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

const TCP_SERVER_PORT = "24094"

var serverResponseCnt int
var serverStartTime time.Time

func main() {
	serverListener := initServer()
	if serverListener == nil {
		printError("Failed to Init Server")
		return
	}

	sigintHandler := make(chan os.Signal, 1)
	signal.Notify(sigintHandler, syscall.SIGINT)
	go func() {
		<-sigintHandler
		closeServer(serverListener)
		os.Exit(0)
	}()

	clientCnt := 0
	clientID := 0
	go func() {
		for {
			serverTimer := time.NewTimer(time.Second * 10)
			<-serverTimer.C
			printLog(fmt.Sprintf("Number of clients connected = %d", clientCnt))
		}
	}()

	serverResponseCnt = 0
	serverStartTime = time.Now()

	for {
		serverConnection, _ := serverListener.Accept()
		go func() {
			clientCnt++
			clientID++
			curClientID := clientID
			printLog(fmt.Sprintf("Client %d connected. Number of clients connected = %d", curClientID, clientCnt))
			for serverConnection != nil {
				requestAddr := serverConnection.RemoteAddr()
				if requestAddr != nil {
					requestBuffer := make([]byte, 1024)
					count, _ := serverConnection.Read(requestBuffer)
					cmd, _ := strconv.Atoi(string(requestBuffer[0]))

					if count == 0 {
						_ = serverConnection.Close()
						break
					}

					responseData := getResponse(cmd, string(requestBuffer[1:count]), requestAddr.String())
					if responseData == "" {
						_ = serverConnection.Close()
						break
					}

					printLog(fmt.Sprintf("TCP Connection Request from %s", requestAddr.String()))
					printLog(fmt.Sprintf("Command %d", cmd))
					_, err := serverConnection.Write([]byte(responseData))
					if err != nil {
						printError("Failed to Send Response")
						continue
					}
					serverResponseCnt++
				}
			}
			clientCnt--
			printLog(fmt.Sprintf("Client %d connected. Number of clients connected = %d", curClientID, clientCnt))
		}()
	}
}

func initServer() net.Listener {
	serverListener, err := net.Listen("tcp4", ":"+TCP_SERVER_PORT)
	if err != nil {
		printError(err.Error())
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
		return strings.TrimSpace(strings.ToUpper(data))
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

func printLog(msg string) {
	_, err := fmt.Printf("[Time: HH:MM:SS] %s\n", msg)
	if err != nil {
		printError("STDOUT Print Error")
		return
	}
}

func printError(msg string) {
	_, err := fmt.Fprintf(os.Stderr, "Error: %s\n", msg)
	if err != nil {
		fmt.Printf("Error: %s\n", msg)
		return
	}
}
