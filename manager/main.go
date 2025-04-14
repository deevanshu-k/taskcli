package manager

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"taskcli/config"
	"taskcli/structs"
)

type Manager struct {
	mu         *sync.Mutex
	tasks      map[int]structs.Task
	storageDir string
}

func NewManager() *Manager {
	return &Manager{
		mu:         &sync.Mutex{},
		tasks:      make(map[int]structs.Task),
		storageDir: config.StorageDir + "/taskcli.db",
	}
}

func (m *Manager) LoadTasks() error {
	file, err := os.OpenFile(m.storageDir, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return fmt.Errorf("failed to open storage file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Bytes()

		var task structs.Task
		if err := task.Unmarshal(line); err != nil {
			return fmt.Errorf("failed to unmarshal task: %v", err)
		}
		m.tasks[task.Id] = task
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner error: %v", err)
	}

	return nil
}

func (m *Manager) AddTask(cmd structs.Command) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if cmd.Task == nil {
		return fmt.Errorf("task is nil")
	}

	id := 0
	for _, task := range m.tasks {
		if task.Id > id {
			id = task.Id
		}
	}
	id++
	task := structs.NewTask(id, *cmd.Task, structs.PENDING)

	file, err := os.OpenFile(m.storageDir, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open storage file: %v", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	b, err := task.Marshal()
	if err != nil {
		return fmt.Errorf("failed to marshal task: %v", err)
	}
	b = append(b, byte('\n'))
	if _, err := writer.Write(b); err != nil {
		return fmt.Errorf("failed to write task to file: %v", err)
	}
	writer.Flush()

	m.tasks[task.Id] = *task

	return nil
}

func (m *Manager) ListTasks(status *structs.Status) (tasks []structs.Task) {
	for _, task := range m.tasks {
		if status == nil || task.Status == *status {
			tasks = append(tasks, task)
		}
	}
	return
}

func (m *Manager) DeleteTask(cmd structs.Command) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if cmd.DeleteAll != nil && *cmd.DeleteAll {
		for id := range m.tasks {
			delete(m.tasks, id)
		}
	} else if cmd.TaskId != nil {
		delete(m.tasks, *cmd.TaskId)
	}

	file, err := os.OpenFile(m.storageDir, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return fmt.Errorf("failed to open storage file: %v", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, task := range m.tasks {
		b, err := task.Marshal()
		if err != nil {
			return fmt.Errorf("failed to marshal task: %v", err)
		}
		b = append(b, byte('\n'))
		if _, err := writer.Write(b); err != nil {
			return fmt.Errorf("failed to write task to file: %v", err)
		}
	}
	writer.Flush()

	return nil
}
