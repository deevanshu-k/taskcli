package structs

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"taskcli/config"
)

type CommandType string

const (
	ADD    CommandType = "ADD"
	REMOVE CommandType = "REMOVE"
	UPDATE CommandType = "UPDATE"
)

type Command struct {
	Type      CommandType `json:"type"`
	TaskId    *string     `json:"task_id"`
	Task      *string     `json:"task"`
	Status    *Status     `json:"status"`
	DeleteAll *bool       `json:"delete_all"`
}

func NewCommand(commandType CommandType, taskId *string, task *string, status *Status, deleteAll *bool) *Command {
	return &Command{
		Type:      commandType,
		TaskId:    taskId,
		Task:      task,
		Status:    status,
		DeleteAll: deleteAll,
	}
}

func (c *Command) Marshal() ([]byte, error) {
	data, err := json.Marshal(*c)
	if err != nil {
		return []byte{}, err
	}
	return data, nil
}

func (c *Command) Unmarshal(data []byte) error {
	if err := json.Unmarshal(data, c); err != nil {
		return err
	}
	return nil
}

func (c *Command) SendCommand() error {
	conn, err := net.Dial("tcp", fmt.Sprintf("%v:%d", config.Config.Host, config.Config.Port))
	if err != nil {
		return fmt.Errorf("failed to connect to server: %w", err)
	}
	defer conn.Close()
	data, err := c.Marshal()
	if err != nil {
		return fmt.Errorf("failed to marshal command: %w", err)
	}
	_, err = conn.Write(data)
	if err != nil {
		return fmt.Errorf("failed to send command: %w", err)
	}

	_, err = conn.Write([]byte("\n"))
	if err != nil {
		return fmt.Errorf("failed to send command: %w", err)
	}

	// Read the response from the server
	reader := bufio.NewReader(conn)
	res, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	res = strings.TrimSpace(res)

	if res == "ok" {
		return nil
	}
	return fmt.Errorf("%s", res)
}
