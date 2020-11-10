package server

import (
	"bytes"
	"io"
	"log"
	"net"
	"strings"
	"sync"
)

type HandleFunc func(conn net.Conn)

type Server struct {
	addr     string
	mu       sync.RWMutex
	handlers map[string]HandleFunc
}

func NewServer(addr string) *Server {
	return &Server{addr: addr, handlers: make(map[string]HandleFunc)}
}

func (s *Server) Register(path string, handler HandleFunc) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.handlers[path] = handler
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Print(err)
		return err
	}

	defer func() {
		err = listener.Close()
		if err != nil {
			log.Print(err)
		}
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			return err
		}

		go s.handle(conn)
	}

	return nil
}

func (s *Server) handle(conn net.Conn) {
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Print(err)
		}
	}()

	buffer := make([]byte, 4096)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				log.Printf("%v", buffer[:n])
				return
			}
			log.Print(err)
			return
		}

		data := buffer[:n]
		requestLineDelim := []byte{'\r', '\n'}
		requestLineEnd := bytes.Index(data, requestLineDelim)
		if requestLineEnd == -1 {
			return
		}

		requestLine := string(data[:requestLineEnd])
		parts := strings.Split(requestLine, " ")
		if len(parts) != 3 {
			return
		}

		path, version := parts[1], parts[2]

		if version != "HTTP/1.1" {
			return
		}

		handler := func(conn net.Conn) {
			err := conn.Close()
			if err != nil {
				log.Print(err)
			}
		}
		s.mu.RLock()
		for i := 0; i < len(s.handlers); i++ {
			if handl, ok := s.handlers[path]; ok {
				handler = handl
				break
			}
		}
		s.mu.RUnlock()
		handler(conn)
	}
}
