package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"
)

type serverState string

const (
	StateInit  serverState = "init"
	StateReady serverState = "ready"
	StateDone  serverState = "done"
)

// Serve, listen & handle run in multiple goroutines, atomic.Bool can get boolean val from multiple goroutines at the same time without causing race conditions.
type Server struct {
	State    serverState
	Listener net.Listener
	closed   atomic.Bool
}

func (s *Server) listen() {
	for {
		if s.closed.Load() {
			return
		}
		conn, err := s.Listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			log.Println("error: ", err)
			continue
		}
		go s.handle(conn) //separate goroutine for each conn
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	response := "HTTP/1.1 200 OK\r\n" +
		"Content-Type: text/plain\r\n" +
		"Content-Length: 14\r\n\r\n" +
		"Hello World!\r\n"

	_, err := conn.Write([]byte(response))
	if err != nil {
		log.Println("error:", err)
	}
}

func Serve(port int) (*Server, error) {
	address := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}
	server := &Server{State: StateInit, Listener: listener}
	go server.listen()
	server.State = StateReady
	return server, nil
}

func (s *Server) Close() error {
	s.closed.Store(true)
	//close listener
	err := s.Listener.Close()
	if err != nil {
		return err
	}
	s.State = StateDone
	return nil
}
