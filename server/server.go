package server

import (
	"bufio"
	"fmt"
	"io"
	"log/slog"
	"net"
	"strings"
)

type Server struct {
	Host string
	Port int
}

func NewServer(host string, port int) *Server {
	return &Server{
		Host: host,
		Port: port,
	}
}

func (s *Server) Start() {
	if err := s.BindAndListen(); err != nil {
		slog.Error("failed to start server", slog.String("error", err.Error()))
	}
}

func (s *Server) BindAndListen() error {
	socket, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Host, s.Port))
	if err != nil {
		return fmt.Errorf("failed to bind to socket: %w", err)
	}
	defer socket.Close()

	for {
		conn, err := socket.Accept()
		if err != nil {
			return fmt.Errorf("failed to accept connection: %w", err)
		}

		slog.Info("client connected", slog.String("remote_addr", conn.RemoteAddr().String()))

		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		data, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				slog.Info("client disconnected")
				break
			}
			slog.Error("failed to read from connection", slog.String("error", err.Error()))
			break
		}

		data = strings.Trim(data, "\n")
		data = strings.TrimSpace(data)

		slog.Info("received", slog.String("data", data))
		conn.Write([]byte("ok\n"))
	}
}
