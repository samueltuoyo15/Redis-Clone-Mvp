package main

import (
	"bufio"
	"fmt"
	"net"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Client disconnected:", conn.RemoteAddr())
			return
		}

		fmt.Printf("Received from %s : %s, conn.RemoteAddr() ", message, conn.RemoteAddr())
		conn.Write([]byte("+PONG\r\n"))
	}
}

func main() {
	fmt.Println("Starting Go Redis server on port 6378")

	listener, err := net.Listen("tcp", ":6378")
	if err != nil {
		panic(err)
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		fmt.Println("New client connected:", conn.RemoteAddr())
		go handleConnection(conn)
	}
}
