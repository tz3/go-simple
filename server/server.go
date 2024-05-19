package server

import (
	"log"
	"net/http"
)

// Server defines the server struct.
type Server struct {
	Port string
}

// NewServer init a new server with defined port.
func NewServer(port string) *Server {
	return &Server{
		Port: port,
	}
}

// RegisterHandler register handler with handler function.
func (s *Server) RegisterHandler(route string, handler http.HandlerFunc) {
	http.HandleFunc(route, handler)
}

// Start start the server with a given the server port.
func (s *Server) Start() {
	log.Printf("Starting server on port %s...\n", s.Port)
	err := http.ListenAndServe(":"+s.Port, nil)
	if err != nil {
		log.Fatal("Failed to start server:", err)
		return
	}
}
