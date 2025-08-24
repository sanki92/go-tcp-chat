package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

var clients = make(map[net.Conn]string)
var mutex = sync.Mutex{}
var broadcast = make(chan string)

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	fmt.Println("Chat server is running on :8080")

	go broadcaster()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error acceptiong:", err)
			continue
		}

		mutex.Lock()
		clients[conn] = ""
		mutex.Unlock()

		fmt.Println("New client connected:", conn.RemoteAddr())

		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer func() {
		mutex.Lock()
		delete(clients, conn)
		mutex.Unlock()
		fmt.Println("Client disconnected:", conn.RemoteAddr())
		conn.Close()
	}()

	reader := bufio.NewReader(conn)
	var nickname string
	for {
		conn.Write([]byte("Enter your nickname: "))
		nickname, _ = reader.ReadString('\n')
		nickname = strings.TrimSpace(nickname)

		taken := false
		mutex.Lock()
		for _, value := range clients {
			if nickname == value {
				taken = true
				break
			}
		}
		mutex.Unlock()

		if taken {
			conn.Write([]byte("Nickname already taken\n"))
		} else {
			break
		}
	}
	mutex.Lock()
	clients[conn] = nickname
	mutex.Unlock()

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		msg = strings.TrimSpace(msg)
		
		if strings.HasPrefix(msg, "@") {
			parts := strings.SplitN(msg, " ", 2)
			if len(parts) < 2 {
				conn.Write([]byte("Usage: @<nickname> <msg>\n"))
				continue
			}
			targetNick,_ := strings.CutPrefix(parts[0], "@")
			privateMsg := parts[1]

			userFound := false
			mutex.Lock()
			for key, value := range clients {
				if value == targetNick {
					fmt.Fprintf(key, "[private from %s] %s\n", clients[conn], privateMsg)
					userFound = true
					break
				}
			}
			mutex.Unlock()

			if !userFound {
				conn.Write([]byte("Username not found!\n"))
			}
		} else {
			if msg != "" {
				msg := fmt.Sprintf("[%s] %s", clients[conn], msg)
				broadcast <- msg
			}
		}
	}
}

func broadcaster() {
	for {
		msg := <-broadcast

		mutex.Lock()
		for conn := range clients {
			_, err := fmt.Fprintln(conn, msg)
			if err != nil {
				conn.Close()
				delete(clients, conn)
			}
		}
		mutex.Unlock()
		fmt.Println("Broadcast:", msg)
	}
}
