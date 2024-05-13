/*
 * ChatServer.go
 * 20194094
 * Yongmin Yoo
 */

package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

const SERVER_PORT = "14094"
const MAX_CLIENT = 8

var clients = make(map[string]net.Conn)

func main() {
	listener, err := net.Listen("tcp4", ":"+SERVER_PORT)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	fmt.Println("Server is running on port", SERVER_PORT)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		if len(clients) == MAX_CLIENT {
			fmt.Fprintln(conn, "KChatting Room is Full. Cannot connect")
			conn.Close()
		}

		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	reader := bufio.NewReader(conn)
	nickname, _ := reader.ReadString('\n')
	nickname = strings.TrimSpace(nickname)

	if _, exists := clients[nickname]; exists {
		fmt.Fprintln(conn, "KNickname is already used by another user. Cannot connect")
		conn.Close()
		return
	}

	fmt.Fprintf(conn, "M[Welcome %s to CAU Net-Class Chat Room at %s.]\n", nickname, conn.LocalAddr().String())
	fmt.Fprintf(conn, "M[There are %d users in the room]\n", len(clients))

	clients[nickname] = conn
	broadcast(fmt.Sprintf("MWelcome %s to the chat.", nickname), nickname)
	fmt.Printf("%s Joined from %s. There are %d users in the room.\n", nickname, conn.RemoteAddr().String(), len(clients))

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		text := scanner.Text()
		if strings.HasPrefix(text, "M") {
			broadcast(fmt.Sprintf("M%s: %s", nickname, text[1:]), nickname)
		} else {
			if command, extra := decodeCommand(text); command != "" {
				runCommand(command, extra, conn, nickname)
			} else if text != "Q" {
				fmt.Fprintf(conn, "KInvalid command.\n")
				fmt.Printf("Invalid command: %s\n", text)
			}
		}
	}

	delete(clients, nickname)
	broadcast(fmt.Sprintf("M%s has left the chat.", nickname), nickname)
	fmt.Printf("%s left the room. There are %d users in the room.\n", nickname, len(clients))
	conn.Close()
}

func decodeCommand(text string) (string, string) {
	parts := strings.SplitN(text, " ", 2)
	if len(parts) > 1 {
		return parts[0], parts[1]
	}
	return "", ""
}

func runCommand(cmd, extra string, conn net.Conn, nickname string) {
	switch cmd {
	case "L":
		listUsers(conn)
	case "P":
		fmt.Fprintf(conn, "P\n")
	case "Q":
		conn.Close()
	case "S":
		handleSecret(extra, conn, nickname)
	case "E":
		handleExcept(extra, conn, nickname)
	default:
		fmt.Fprintf(conn, "MInvalid command received: %s\n", cmd)
	}
}

func listUsers(conn net.Conn) {
	for nick, clientConn := range clients {
		fmt.Fprintf(conn, "I%s - %s\n", nick, clientConn.RemoteAddr().String())
	}
}

func handleSecret(details string, conn net.Conn, nickname string) {
	parts := strings.SplitN(details, " ", 2)
	if len(parts) < 2 {
		return
	}
	target, message := parts[0], parts[1]
	if targetConn, ok := clients[target]; ok {
		fmt.Fprintln(targetConn, fmt.Sprintf("Mfrom: %s> %s", nickname, message))
	}
}

func handleExcept(details string, conn net.Conn, nickname string) {
	parts := strings.SplitN(details, " ", 2)
	if len(parts) < 2 {
		return
	}
	except, message := parts[0], parts[1]
	for nick, clientConn := range clients {
		if nick != except && conn != clientConn {
			fmt.Fprintln(clientConn, fmt.Sprintf("Mfrom %s> %s", nickname, message))
		}
	}
}

func broadcast(message, skip string) {
	for nick, conn := range clients {
		if nick != skip {
			fmt.Fprintln(conn, message)
		}
	}
}
