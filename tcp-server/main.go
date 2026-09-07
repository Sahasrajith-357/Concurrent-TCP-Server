package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer listener.Close()
	log.Printf("Server is listening on: %v", listener.Addr())

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		<-ctx.Done()
		log.Println("shutting down server...")
		listener.Close()
	}()

	var wg sync.WaitGroup

	for {
		conn, err := listener.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				log.Println("server stopped accepting new connections")
				wg.Wait()
				log.Println("all connections closed, exiting")
				return
			default:
				log.Printf("failed to accept connection: %v", err)
				continue
			}
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			handle(conn)
		}()
	}
}

func handle(connection net.Conn) {
	defer connection.Close()
	addr := connection.RemoteAddr()
	log.Printf("client connected: %s", addr)

	scanner := bufio.NewScanner(connection)
	for {
		connection.SetReadDeadline(time.Now().Add(5 * time.Minute))

		if scanner.Err() != nil || !scanner.Scan() {
			break
		}

		line := scanner.Text()
		log.Printf("[%s] received: %q", addr, line)

		if _, err := fmt.Fprintf(connection, "echo: %s\n", line); err != nil {
			log.Printf("[%s] write error: %v", addr, err)
			break
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("[%s] read error: %v", addr, err)
	}

	log.Printf("client disconnected: %s", connection.RemoteAddr())
}
