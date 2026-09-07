package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:9000")
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	log.Printf("connected to server: %s", conn.RemoteAddr())

	reply := bufio.NewScanner(conn)
	go func() {
		for reply.Err() == nil && reply.Scan() {
			line := reply.Text()
			fmt.Printf("received: %s\n", line)
		}
		log.Printf("server disconnected: %s", conn.RemoteAddr())
	}()

	input := bufio.NewScanner(os.Stdin)
	for input.Err() == nil && input.Scan() {
		text := input.Text()
		if text == "exit" {
			break
		}
		fmt.Fprintf(conn, "%s\n", text)
	}

	fmt.Printf("disconnected from server: %s\n", conn.RemoteAddr())
}
