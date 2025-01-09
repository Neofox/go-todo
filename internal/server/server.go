package server

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Server struct {
	port       int
	httpServer *http.Server
}

func NewServer(port int) *Server {
	NewServer := &Server{port: port}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      NewServer.RegisterRoutes(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	NewServer.httpServer = server

	return NewServer
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	fmt.Println("Stopping server...")

	return s.httpServer.Shutdown(ctx)
}
