package libs

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
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

func lastId() (int, error) {
	records, err := AllData()
	if err != nil {
		return 0, errors.New("fail to ready db")
	}
	if len(records) == 0 {
		return 0, nil
	}

	var tasks []Task
	for _, record := range records {
		id, _ := strconv.Atoi(record[0])
		description := record[1]
		var status int
		if record[2] == Pending.String() {
			status = int(Pending)
		} else if record[2] == Inprogress.String() {
			status = int(Inprogress)
		} else {
			status = int(Complete)
		}

		task := Task{
			Id:          id,
			Description: description,
			Status:      StatusEnum(status),
			CreatedAt:   record[3],
		}
		tasks = append(tasks, task)
	}

	return tasks[len(tasks)-1].Id, nil
}

func CreateTask(tasks []string, status StatusEnum) error {
	file, err := os.OpenFile("data.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error while opening the db %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	// Data to append
	t := time.Now()
	id, err := lastId()
	if err != nil {
		return fmt.Errorf("error while reading the db  %w", err)
	}
	newRecords := [][]string{}
	for _, task := range tasks {
		id = id + 1
		newRecord := []string{fmt.Sprint(id), task, fmt.Sprint(int(status)), t.Format("02-01-2006")}
		newRecords = append(newRecords, newRecord)
	}

	// Write the new record to the CSV file
	err = writer.WriteAll(newRecords)
	writer.Flush()
	if err != nil {
		return fmt.Errorf("failed to write record to file %w", err)
	}

	return nil
}

func UpdateStatus(id string, status string) error {
	records, err := AllData()
	if err != nil {
		return err
	}

	for i := 0; i < len(records); i++ {
		if records[i][0] == id {
			records[i][2] = status
			err := reFreshData(records)
			if err != nil {
				return err
			}
			return nil
		}
	}

	return errors.New("task with this id not exist")
}

func UpdateTask(id string, task string) error {
	records, err := AllData()
	if err != nil {
		return err
	}

	for i := 0; i < len(records); i++ {
		if records[i][0] == id {
			records[i][1] = task
			err := reFreshData(records)
			if err != nil {
				return err
			}
			return nil
		}
	}

	return errors.New("task with this id not exist")
}

func DeleteAll() error {
	return reFreshData([][]string{})
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
