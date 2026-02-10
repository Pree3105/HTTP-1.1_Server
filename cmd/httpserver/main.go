package main

import (
	"httpServer_go/internal/server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const port = 42069

func main() {
	server, err := server.Serve(port) //recieves port and handles requests
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	//sigChan-for shutting down the server properly,server.Server returns immediateley, but if we exit main immedietely ,server will also stop, but we want to wait for a signal(ctrl c) before stopping server
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
