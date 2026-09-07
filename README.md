# Multithreaded Go TCP Server

A concurrent TCP echo server and client built with Go's standard `net` package. The server accepts multiple simultaneous clients using a goroutine-per-connection model, echoes back every line it receives, and shuts down gracefully on interrupt.

## Features

- **Concurrent clients** — each connection is handled in its own goroutine, so thousands of clients can be served simultaneously without threads or an event loop.
- **Graceful shutdown** — `SIGINT`/`SIGTERM` stops new connections and waits for active clients to finish via a `sync.WaitGroup`.
- **Read timeouts** — idle connections are dropped after a configurable deadline using `SetReadDeadline`.
- **Uniform I/O** — server and client both operate on the same `net.Conn` interface, reading and writing line-delimited messages.

## Requirements

- Go 1.22 or newer

## Project layout

```text
.
└── tcp-multithreading/
    ├── tcp-server/             # The concurrent TCP echo server
    │   └── main.go
    └── tcp-client/             # A TCP client
        └── main.go
```

## Running the server

```bash
cd tcp-multithreading
cd tcp-server
go run main.go
```

The server listens on `:9000` and logs each connection, message, and disconnect.

## Connecting a client

Using the included Go client:

```bash
cd tcp-multithreading
cd tcp-client
go run main.go
```

Type a message and press Enter — the server replies with `echo: <message>`.
Press `Ctrl+D` to quit.

Or connect with any TCP tool:

```bash
nc localhost 9000
```

## How it works

The server's accept loop is intentionally simple:

```go
for {
    conn, err := listener.Accept()
    if err != nil {
        // handle shutdown / transient errors
    }
    go handle(conn) // each client runs concurrently
}
```

The single `go` keyword is what makes the server concurrent. Goroutines are lightweight (a few KB of stack, multiplexed onto OS threads by the Go runtime), so a goroutine-per-connection scales far beyond a thread-per-connection design.

Graceful shutdown is coordinated with `context` and a `sync.WaitGroup`: on interrupt, the listener is closed to unblock `Accept()`, and the server waits for in-flight connections to drain before exiting.

## Configuration

| Setting        | Location     | Default          |
| -------------- | ------------ | ---------------- |
| Listen address | `tcp-server` | `:9000`          |
| Read deadline  | `tcp-server` | 5 minutes        |
| Server address | `tcp-client` | `localhost:9000` |

## License

MIT
