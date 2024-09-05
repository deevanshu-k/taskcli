package libs

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

type Task struct {
	Id          int        `json:"id"`
	Description string     `json:"task"`
	Status      StatusEnum `json:"status"`
	CreatedAt   string     `json:"created_at"`
}

type StatusEnum int

const (
	Pending    StatusEnum = 0
	Inprogress StatusEnum = 1
	Complete   StatusEnum = 2
)

// Implement the Stringer interface
func (d StatusEnum) String() string {
	return [...]string{"Pending", "In-progress", "Complete"}[d]
}

func AllData() ([][]string, error) {
	base_url := os.Getenv("BASE_URL")
	resp, err := http.Get(base_url + "/tasks")
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var tasks []Task
	err = json.Unmarshal(body, &tasks)
	if err != nil {
		return nil, err
	}

	var data [][]string
	for i := 0; i < len(tasks); i++ {
		st := "0"
		if tasks[i].Status == Inprogress {
			st = "1"
		}
		if tasks[i].Status == Complete {
			st = "2"
		}
		data = append(data, []string{strconv.Itoa(tasks[i].Id), tasks[i].Description, st, tasks[i].CreatedAt})
	}

	return data, nil
}

func CreateTask(tasks []string, status StatusEnum) error {
	base_url := os.Getenv("BASE_URL")
	var reqBody struct {
		Tasks []string `json:"tasks"`
	}
	reqBody.Tasks = tasks
	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}
	_, err = http.Post(base_url+"/createTasks", "application/json", bytes.NewReader(reqBodyBytes))
	if err != nil {
		return err
	}
	return nil
}

func UpdateStatus(id int, status int) error {
	base_url := os.Getenv("BASE_URL")
	var reqBody struct {
		Id     int `json:"id"`
		Status int `json:"status"`
	}
	reqBody.Id = id
	reqBody.Status = status
	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	// Create a new PATCH request
	req, err := http.NewRequest(http.MethodPatch, base_url+"/updateTask", bytes.NewReader(reqBodyBytes))
	if err != nil {
		return err
	}
	// Set request headers
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Find count
	var resBodyJson struct {
		Count int `json:"count"`
	}
	if err := json.Unmarshal(resBody, &resBodyJson); err != nil {
		return err
	}
	if resBodyJson.Count == 0 {
		return fmt.Errorf("task with this id not exist")
	}
	return nil
}

func UpdateTask(id int, task string) error {
	base_url := os.Getenv("BASE_URL")
	var reqBody struct {
		Id   int    `json:"id"`
		Task string `json:"task"`
	}
	reqBody.Id = id
	reqBody.Task = task
	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	// Create a new PATCH request
	req, err := http.NewRequest(http.MethodPatch, base_url+"/updateTask", bytes.NewReader(reqBodyBytes))
	if err != nil {
		return err
	}
	// Set request headers
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Find count
	var resBodyJson struct {
		Count int `json:"count"`
	}
	if err := json.Unmarshal(resBody, &resBodyJson); err != nil {
		return err
	}
	if resBodyJson.Count == 0 {
		return fmt.Errorf("task with this id not exist")
	}
	return nil
}

func DeleteAll() error {
	base_url := os.Getenv("BASE_URL")
	var reqBody struct {
		Ids       []int `json:"ids"`
		DeleteAll bool  `json:"delete_all"`
	}
	reqBody.DeleteAll = true
	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	// Create a new DELETE request
	req, err := http.NewRequest(http.MethodDelete, base_url+"/deleteTasks", bytes.NewReader(reqBodyBytes))
	if err != nil {
		return err
	}
	// Set request headers
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Find count
	var resBodyJson struct {
		Count int `json:"count"`
	}
	if err := json.Unmarshal(resBody, &resBodyJson); err != nil {
		return err
	}
	return nil
}

func DeleteByIds(ids []string) error {
	// mapping of ids for fast lookup
	idsMap := make(map[string]bool)
	for _, id := range ids {
		idsMap[id] = true
	}

	// Get all data
	records, err := AllData()
	if err != nil {
		return err
	}

	// Filter data
	var filteredData = [][]string{}
	var anyChange = false
	for _, record := range records {
		if idsMap[record[0]] {
			anyChange = true
		} else {
			filteredData = append(filteredData, record)
		}
	}

	if !anyChange {
		return errors.New("no tasks found")
	}

	// Save the filtered data
	if err := reFreshData(filteredData); err != nil {
		return err
	}

	return nil
}

func reFreshData(data [][]string) error {
	file, err := os.OpenFile("data.csv", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("error while opening the db %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.WriteAll(data)
	return nil
}
