package structs

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Status rune

const (
	PENDING    Status = 'P'
	INPROGRESS Status = 'I'
	COMPLETED  Status = 'C'
)

type Task struct {
	Id     string `json:"id"`
	Name   string `json:"task"`
	Status Status `json:"status"`
	Date   string `json:"date"`
}

func NewTask(name string, status Status) *Task {
	id := uuid.NewString()
	data := time.Now().Format("2006-01-02 15:04:05")
	return &Task{
		Id:     id,
		Name:   name,
		Status: status,
		Date:   data,
	}
}

func (t *Task) Marshal() ([]byte, error) {
	data, err := json.Marshal(*t)
	if err != nil {
		return []byte{}, fmt.Errorf("failed to marshal task: %w", err)
	}
	return data, nil
}

func (t *Task) Unmarshal(data []byte) error {
	if err := json.Unmarshal(data, t); err != nil {
		return fmt.Errorf("failed to unmarshal task: %w", err)
	}
	return nil
}
