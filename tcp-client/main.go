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

	input := bufio.NewScanner(os.Stdin)
	reply := bufio.NewScanner(conn)
	for input.Err() == nil && input.Scan() {
		text := input.Text()
		fmt.Fprintf(conn, "%s\n", text)
		if text == "exit" {
			break
		}
		for reply.Err() == nil && reply.Scan() {
			replyText := reply.Text()
			fmt.Printf("received from server: %s\n", replyText)
		}
	}

	fmt.Printf("disconnected from server: %s\n", conn.RemoteAddr())
}
