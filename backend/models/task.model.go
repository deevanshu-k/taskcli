package models

type Task struct {
	Id         int    `json:"id"`
	Task       string `json:"task"`
	Status     int    `json:"status"`
	Created_at string `json:"created_at"`
}
