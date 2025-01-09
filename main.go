package main

import (
	"fmt"

	"todo/internal/server"
)

const (
	port = 8080
)

func main() {
	fmt.Println("Starting server on port", port)

	server := server.NewServer(port)
	err := server.Start()
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
