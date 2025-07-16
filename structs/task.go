package structs

import (
	"encoding/json"
	"fmt"
	"time"
)

type Status rune

const (
	PENDING    Status = 'P'
	INPROGRESS Status = 'I'
	COMPLETED  Status = 'C'
)

func (s Status) String() string {
	switch s {
	case PENDING:
		return "Pending"
	case INPROGRESS:
		return "In-Progress"
	case COMPLETED:
		return "Completed"
	default:
		return "Unknown"
	}
}

type Notify string

const (
	ON  Notify = "on"
	OFF Notify = "off"
)

func (n Notify) String() string {
	switch n {
	case ON:
		return "On"
	case OFF:
		return "Off"
	default:
		return "-"
	}
}

type NotificationTime struct {
	Hour   int8 `json:"hour"`
	Minute int8 `json:"minute"`
}

func (nt NotificationTime) IsBefore(hour int8, minute int8) bool {
	return nt.Hour < hour || (nt.Hour == hour && nt.Minute < minute)
}

func (nt NotificationTime) String() string {
	return fmt.Sprintf("%02d:%02d", nt.Hour, nt.Minute)
}

type Task struct {
	Id               int              `json:"id"`
	Name             string           `json:"task"`
	Status           Status           `json:"status"`
	Date             string           `json:"date"`
	NotificationTime NotificationTime `json:"notification_time"`
	Notify           Notify           `json:"notify"`
}

func NewTask(id int, name string, status Status, notificationTime NotificationTime) *Task {
	data := time.Now().Format("2006-01-02 15:04:05")
	return &Task{
		Id:               id,
		Name:             name,
		Status:           status,
		Date:             data,
		NotificationTime: notificationTime,
		Notify:           OFF,
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
