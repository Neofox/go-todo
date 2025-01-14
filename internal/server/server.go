package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"regotth/web"
)

// Server is the main server struct
// This is just boilerplate code to get the server running
type Server struct {
	port       int
	httpServer *http.Server
}

// NewServer is the constructor for the Server struct
// It sets up the server and returns it
func NewServer(port int) *Server {
	NewServer := &Server{port: port}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      NewServer.RegisterRoutes(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	NewServer.httpServer = server

	if err := web.AssetInit(); err != nil {
		slog.Error("Failed to load asset manifest", "error", err)
	}

	return NewServer
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	slog.Info("Stopping server...")

	return s.httpServer.Shutdown(ctx)
}
