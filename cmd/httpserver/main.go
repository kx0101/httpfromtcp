package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/kx0101/httpfromtcp/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	defer server.Close()
	fmt.Printf("server started on port %d\n", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	fmt.Println("shutting down server")
}
