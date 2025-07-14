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
	UPDATE CommandType = "UPDATE"
	LIST   CommandType = "LIST"
	STATUS CommandType = "STATUS"
	DELETE CommandType = "DELETE"
)

type Command struct {
	Type             CommandType       `json:"type"`
	TaskId           *int              `json:"task_id"`
	Task             *string           `json:"task"`
	Status           *Status           `json:"status"`
	DeleteAll        *bool             `json:"delete_all"`
	NotificationTime *NotificationTime `json:"notification_time"`
	Notify           *Notify           `json:"notify"`
}

func NewCommand(commandType CommandType, taskId *int, task *string, status *Status, deleteAll *bool, notificationTime *NotificationTime, notify *Notify) *Command {
	return &Command{
		Type:             commandType,
		TaskId:           taskId,
		Task:             task,
		Status:           status,
		DeleteAll:        deleteAll,
		NotificationTime: notificationTime,
		Notify:           notify,
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

func (c *Command) SendCommand() (Response, error) {
	var res Response

	// Connect to the server
	conn, err := net.Dial("tcp", fmt.Sprintf("%v:%d", config.Config.Host, config.Config.Port))
	if err != nil {
		return res, fmt.Errorf("failed to connect to server: %w", err)
	}
	defer conn.Close()

	// Marshal the command to JSON
	data, err := c.Marshal()
	if err != nil {
		return res, fmt.Errorf("failed to marshal command: %w", err)
	}

	// Append a newline character to the command
	data = append(data, byte('\n'))

	// Write the command to the connection
	_, err = conn.Write(data)
	if err != nil {
		return res, fmt.Errorf("failed to send command: %w", err)
	}

	// Read the response from the server
	reader := bufio.NewReader(conn)
	response, err := reader.ReadString('\n')
	if err != nil {
		return res, fmt.Errorf("failed to read response: %w", err)
	}

	// Remove the trailing newline character
	response = strings.TrimSpace(response)

	// Unmarshal the response
	if err := res.Unmarshal([]byte(response)); err != nil {
		return res, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Check for errors in the response
	if res.Error != nil {
		return res, fmt.Errorf("error from server: %s", *res.Error)
	}

	return res, nil
}
