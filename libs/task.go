package libs

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Task struct {
	Id          int
	Description string
	Status      StatusEnum
	CreatedAt   string
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
	file, err := os.Open("data.csv")
	if err != nil {
		return nil, errors.New("fail to open db")
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, errors.New("fail to ready db")
	}
	return records, nil
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

func CreateTask(description string, status StatusEnum) error {
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
	newRecord := []string{fmt.Sprint(id + 1), description, fmt.Sprint(int(status)), t.Format("02-01-2006")}

	// Write the new record to the CSV file
	err = writer.Write(newRecord)
	writer.Flush()
	if err != nil {
		return fmt.Errorf("failed to write record to file %w", err)
	}

	return nil
}
