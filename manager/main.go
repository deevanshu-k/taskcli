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
	tasks      map[string]structs.Task
	storageDir string
}

func NewManager() *Manager {
	return &Manager{
		mu:         &sync.Mutex{},
		tasks:      make(map[string]structs.Task),
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

	task := structs.NewTask(*cmd.Task, structs.PENDING)

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
