package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"taskclibackend/models"
)

type ReturnError struct {
	Message string `json:"message"`
}

func writeError(w http.ResponseWriter, msg string, status int) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ReturnError{Message: msg})
}

func Ping(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "application/json")
		type ping struct {
			Pong bool `json:"pong"`
		}
		if err := json.NewEncoder(w).Encode(ping{Pong: true}); err != nil {
			fmt.Println("Error while encoding json in response")
			writeError(w, "Something went wrong!", http.StatusInternalServerError)
			fmt.Println("Error encoding JSON:", err)
		}
	}
}

func GetTasks(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "application/json")
		// Ensure the request method is GET
		if r.Method != http.MethodGet {
			writeError(w, "404 not found!", http.StatusNotFound)
			return
		}
		// Fetch tasks from db
		rows, err := db.Query("SELECT id,task,status,created_at FROM tasks")
		if err != nil {
			writeError(w, "Something went wrong!", http.StatusInternalServerError)
			fmt.Println("Database query error:", err)
			return
		}
		// Extract tasks from query response
		var tasks []models.Task
		for rows.Next() {
			var task models.Task
			err := rows.Scan(&task.Id, &task.Task, &task.Status, &task.Created_at)
			if err != nil {
				writeError(w, "Something went wrong!", http.StatusInternalServerError)
				fmt.Println("Error scanning row:", err)
				return
			}
			tasks = append(tasks, task)
		}

		if len(tasks) == 0 {
			w.Write([]byte("[]"))
			return
		}
		if err := json.NewEncoder(w).Encode(tasks); err != nil {
			writeError(w, "Something went wrong!", http.StatusInternalServerError)
			fmt.Println("Error encoding JSON:", err)
		}
	}
}

func CreateTasks(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "application/json")
		// Ensure the request method is POST
		if r.Method != http.MethodPost {
			writeError(w, "404 not found!", http.StatusNotFound)
			return
		}

		// Read body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Println("Error while reading body: ", err)
			writeError(w, "Something went wrong!", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		var tasks struct {
			Tasks []string `json:"tasks"`
		}
		err = json.Unmarshal(body, &tasks)
		if err != nil {
			fmt.Println("Error parsing JSON: ", err)
			writeError(w, "Something went wrong!", http.StatusInternalServerError)
			return
		}

		// Validate data
		if len(tasks.Tasks) == 0 {
			fmt.Println("Body is empty!")
			writeError(w, "Need at least one task!", http.StatusBadRequest)
			return
		}
		// Prepare Query
		query := "INSERT INTO tasks(task,status) VALUES "
		for i := 0; i < len(tasks.Tasks); i++ {
			query += "('" + tasks.Tasks[i] + "',0)"
			if i != len(tasks.Tasks)-1 {
				query += ","
			}
		}
		stmt, err := db.Prepare(query)
		if err != nil {
			fmt.Println("Error while prepareing query: ", err)
			writeError(w, "Something went wrong!", http.StatusInternalServerError)
			return
		}

		// Execute query
		result, err := stmt.Exec()
		if err != nil {
			fmt.Println("Error while executing query: ", err)
			writeError(w, "Something went wrong!", http.StatusInternalServerError)
			return
		}

		count, err := result.RowsAffected()
		if err != nil {
			fmt.Println("Error while getting id: ", err)
			writeError(w, "Something went wrong!", http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(struct {
			Count int64 `json:"count"`
		}{Count: count}); err != nil {
			fmt.Println("Error encoding JSON: ", err)
			writeError(w, "Something went wrong!", http.StatusInternalServerError)
			return
		}
	}
}

func DeleteTasks(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "application/json")
		// Ensure the request method is DELETE
		if r.Method != http.MethodDelete {
			writeError(w, "404 not found!", http.StatusNotFound)
			return
		}

		// Read body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Println("Error while reading body: ", err)
			writeError(w, "Something went wrong!", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		var tasks struct {
			Ids       []int `json:"ids"`
			DeleteAll bool  `json:"delete_all"`
		}
		err = json.Unmarshal(body, &tasks)
		if err != nil {
			fmt.Println("Error parsing JSON: ", err)
			writeError(w, "Something went wrong!", http.StatusInternalServerError)
			return
		}
		// Validate data
		if !tasks.DeleteAll && len(tasks.Ids) == 0 {
			fmt.Println("Body is empty!")
			writeError(w, "Need at least one task id!", http.StatusBadRequest)
			return
		}

		// Prepare Query
		var query string
		if !tasks.DeleteAll {
			query = "DELETE FROM tasks WHERE id in ("
			for i := 0; i < len(tasks.Ids); i++ {
				query += strconv.Itoa(tasks.Ids[i])
				if i != len(tasks.Ids)-1 {
					query += ","
				}
			}
			query += ")"
		} else {
			query = "TRUNCATE TABLE tasks"
		}

		fmt.Println(query)
		stmt, err := db.Prepare(query)
		if err != nil {
			fmt.Println("Error while prepareing query: ", err)
			writeError(w, "Something went wrong!", http.StatusInternalServerError)
			return
		}

		// Execute query
		result, err := stmt.Exec()
		if err != nil {
			fmt.Println("Error while executing query: ", err)
			writeError(w, "Something went wrong!", http.StatusInternalServerError)
			return
		}

		count, err := result.RowsAffected()
		if err != nil {
			fmt.Println("Error while getting id: ", err)
			writeError(w, "Something went wrong!", http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(struct {
			Count int64 `json:"count"`
		}{Count: count}); err != nil {
			fmt.Println("Error encoding JSON: ", err)
			writeError(w, "Something went wrong!", http.StatusInternalServerError)
			return
		}
	}
}

func UpdateTask(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "application/json")
		// Ensure the request method is DELETE
		if r.Method != http.MethodPatch {
			writeError(w, "404 not found!", http.StatusNotFound)
			return
		}

		// Read body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Println("Error while reading body: ", err)
			writeError(w, "Something went wrong!", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		var tasks struct {
			Id     *int    `json:"id"`
			Task   *string `json:"task,omitempty"`
			Status *int    `json:"status,omitempty"`
		}
		err = json.Unmarshal(body, &tasks)
		if err != nil {
			fmt.Println("Error parsing JSON: ", err)
			writeError(w, "Something went wrong!", http.StatusInternalServerError)
			return
		}

		// Validate data
		if tasks.Id == nil {
			fmt.Println("Id is empty!")
			writeError(w, "Id is required!", http.StatusBadRequest)
			return
		}
		if tasks.Status == nil && tasks.Task == nil {
			fmt.Println("both status and task is missing!")
			writeError(w, "Status or task is required!", http.StatusBadRequest)
			return
		}
		if tasks.Status != nil && (*tasks.Status < 0 || *tasks.Status > 2) {
			fmt.Println("Status value is incorrect!")
			writeError(w, "Status can only have 0,1,2 as value!", http.StatusBadRequest)
			return
		}

		// Prepare Query
		query := "UPDATE tasks SET"
		if tasks.Status != nil {
			query += " status = " + strconv.Itoa(*tasks.Status) + " "
		}
		if tasks.Task != nil {
			if tasks.Status != nil {
				query += ","
			}
			query += " task = '" + *tasks.Task + "' "
		}
		query += "WHERE id = " + strconv.Itoa(*tasks.Id)
		fmt.Println(query)
		stmt, err := db.Prepare(query)
		if err != nil {
			fmt.Println("Error while prepareing query: ", err)
			writeError(w, "Something went wrong!", http.StatusInternalServerError)
			return
		}

		// Execute query
		result, err := stmt.Exec()
		if err != nil {
			fmt.Println("Error while executing query: ", err)
			writeError(w, "Something went wrong!", http.StatusInternalServerError)
			return
		}

		count, err := result.RowsAffected()
		if err != nil {
			fmt.Println("Error while getting id: ", err)
			writeError(w, "Something went wrong!", http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(struct {
			Count int64 `json:"count"`
		}{Count: count}); err != nil {
			fmt.Println("Error encoding JSON: ", err)
			writeError(w, "Something went wrong!", http.StatusInternalServerError)
			return
		}

	}
}
