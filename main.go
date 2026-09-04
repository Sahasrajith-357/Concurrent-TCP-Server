package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer listener.Close()
	log.Printf("Server is listening on: %v", listener.Addr())

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("failed to accept connection: %v", err)
			continue
		}
		log.Printf("client connected: %s", conn.RemoteAddr())

		handle(conn)
	}
}

func handle(connection net.Conn) {
	defer connection.Close()

	scanner := bufio.NewScanner(connection)
	for scanner.Err() == nil && scanner.Scan() {
		line := scanner.Text()
		log.Printf("received: %q", line)

		fmt.Fprintf(connection, "echo: %s\n", line)
	}

	log.Printf("client disconnected: %s", connection.RemoteAddr())
}
