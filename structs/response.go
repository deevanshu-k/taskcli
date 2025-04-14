package structs

import (
	"encoding/json"
	"fmt"
)

type Response struct {
	Type    CommandType `json:"type"`
	Tasks   []Task      `json:"tasks"`
	Success bool        `json:"success"`
	Error   *string     `json:"error"`
	Message *string     `json:"message"`
}

func NewResponse(commandTye CommandType, tasks []Task, success bool, error *string) *Response {
	return &Response{
		Type:    commandTye,
		Tasks:   tasks,
		Success: success,
		Error:   error,
	}
}

func (r *Response) Marshal() ([]byte, error) {
	data, err := json.Marshal(*r)
	if err != nil {
		return []byte{}, err
	}
	return data, nil
}

func (r *Response) Unmarshal(data []byte) error {
	if err := json.Unmarshal(data, r); err != nil {
		return err
	}
	return nil
}

func (r Response) String() string {
	if r.Message != nil {
		return *r.Message
	}
	return fmt.Sprintf("Type: %s\nTasks: %v\nSuccess: %t\nError: %v\nMessage: %v\n", r.Type, r.Tasks, r.Success, r.Error, r.Message)
}
