package main

import (
	"log/slog"

	"todo/internal/server"
)

const (
	port = 8080
)

func main() {
	slog.Info("Starting server on port", "port", port)

	server := server.NewServer(port)
	err := server.Start()
	if err != nil {
		slog.Error("Error starting server:", "error", err)
	}
}
