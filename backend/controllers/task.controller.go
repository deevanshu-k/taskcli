package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"taskclibackend/models"
)

func Ping(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "application/json")
		type ping struct {
			Pong bool `json:"pong"`
		}
		if err := json.NewEncoder(w).Encode(ping{Pong: true}); err != nil {
			http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
			fmt.Println("Error encoding JSON:", err)
		}
	}
}

func GetTasks(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Fetch tasks from db
		rows, err := db.Query("SELECT id,task,status,created_at FROM tasks")
		if err != nil {
			http.Error(w, "Database query error", http.StatusInternalServerError)
			fmt.Println("Database query error:", err)
			return
		}
		// Extract tasks from query response
		var tasks []models.Task
		for rows.Next() {
			var task models.Task
			err := rows.Scan(&task.Id, &task.Task, &task.Status, &task.Created_at)
			if err != nil {
				http.Error(w, "Error scanning row", http.StatusInternalServerError)
				fmt.Println("Error scanning row:", err)
				return
			}
			tasks = append(tasks, task)
		}

		w.Header().Add("content-type", "application/json")
		if err := json.NewEncoder(w).Encode(tasks); err != nil {
			http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
			fmt.Println("Error encoding JSON:", err)
		}
	}
}
