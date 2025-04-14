package server

import (
	"bufio"
	"fmt"
	"io"
	"log/slog"
	"net"
	"taskcli/manager"
	"taskcli/structs"
)

type Server struct {
	Host    string
	Port    int
	Manager manager.Manager
}

func NewServer(host string, port int) *Server {
	mg := manager.NewManager()
	mg.LoadTasks()
	return &Server{
		Host:    host,
		Port:    port,
		Manager: *mg,
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
				break
			}
			slog.Error("failed to read from connection", slog.String("error", err.Error()))
			break
		}

		var command structs.Command
		if err := command.Unmarshal([]byte(data)); err != nil {
			slog.Error("failed to unmarshal command", slog.String("error", err.Error()))
			conn.Write([]byte("error\n"))
			continue
		}

		slog.Info("received", slog.Any("CMD", command))

		var res structs.Response
		res.Type = command.Type
		res.Success = true

		switch command.Type {
		case structs.ADD:
			if err := s.Manager.AddTask(command); err != nil {
				slog.Error("failed to add task", slog.String("error", err.Error()))
				e := string(err.Error())
				res.Error = &e
				res.Success = false
			}
		case structs.LIST:
			tasks := s.Manager.ListTasks(command.Status)
			res.Tasks = tasks
		case structs.DELETE:
			if err := s.Manager.DeleteTask(command); err != nil {
				slog.Error("failed to delete task", slog.String("error", err.Error()))
				e := string(err.Error())
				res.Error = &e
				res.Success = false
			}
		case structs.UPDATE:
			if command.TaskId == nil {
				slog.Error("task id is nil")
				e := "task id is nil"
				res.Error = &e
				res.Success = false
			} else if err := s.Manager.UpdateTask(command); err != nil {
				slog.Error("failed to update task", slog.String("error", err.Error()))
				e := string(err.Error())
				res.Error = &e
				res.Success = false
			}
		default:
			slog.Error("unknown command type", slog.String("Type", string(command.Type)))
			e := "unknown command type"
			res.Error = &e
			res.Success = false
		}

		b, err := res.Marshal()
		if err != nil {
			slog.Error("failed to marshal response", slog.String("error", err.Error()))
			continue
		}

		b = append(b, byte('\n'))

		conn.Write(b)
	}
}
