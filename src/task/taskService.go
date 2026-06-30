package task

import (
	"fmt"
	"slices"
	"time"

	"github.com/Chibas/TaskTrackerGo/src/storage"
)

type TaskService interface {
	AddTask(description string) error
	UpdateTask(id int, description string) error
	DeleteTask(id int) error
	MarkDone(id int) error
	MarkInProgress(id int) error
	List(status string) ([]storage.Task, error)
}

type taskService struct {
	storage storage.Storage
}

func NewTaskService(storage storage.Storage) TaskService {
	return &taskService{storage: storage}
}

func (s *taskService) AddTask(description string) error {
	taskList, err := s.storage.ReadStorage()
	if err != nil {
		return err
	}

	task := storage.Task{
		Description: description,
		ID:          getMaxId(taskList),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Status:      string(storage.Todo),
	}

	for _, task := range taskList {
		if task.Description == description {
			return fmt.Errorf("Taks already exists")
		}
	}

	taskList = append(taskList, task)
	s.storage.WriteStorage(taskList)
	return nil
}

func (s *taskService) UpdateTask(id int, description string) error {
	taskList, err := s.storage.ReadStorage()
	if err != nil {
		return err
	}
	taskIndex := slices.IndexFunc(taskList, func(t storage.Task) bool {
		return t.ID == id
	})
	if taskIndex < 0 {
		return fmt.Errorf("task with id %d not found", id)
	}

	t := &taskList[taskIndex]
	t.Description = description
	t.UpdatedAt = time.Now()

	s.storage.WriteStorage(taskList)
	return nil
}

func (s *taskService) DeleteTask(id int) error {
	taskList, err := s.storage.ReadStorage()
	if err != nil {
		return err
	}
	taskIndex := slices.IndexFunc(taskList, func(t storage.Task) bool {
		return t.ID == id
	})
	if taskIndex < 0 {
		return fmt.Errorf("task with id %d not found", id)
	}
	taskList = slices.Delete(taskList, taskIndex, taskIndex+1)
	s.storage.WriteStorage(taskList)
	return nil
}

func (s *taskService) MarkDone(id int) error {
	taskList, err := s.storage.ReadStorage()
	if err != nil {
		return err
	}
	taskIndex := slices.IndexFunc(taskList, func(t storage.Task) bool {
		return t.ID == id
	})
	if taskIndex < 0 {
		return fmt.Errorf("task with id %d not found", id)
	}

	t := &taskList[taskIndex]
	t.Status = string(storage.Done)
	t.UpdatedAt = time.Now()
	s.storage.WriteStorage(taskList)
	return nil
}

func (s *taskService) MarkInProgress(id int) error {
	taskList, err := s.storage.ReadStorage()
	if err != nil {
		return err
	}
	taskIndex := slices.IndexFunc(taskList, func(t storage.Task) bool {
		return t.ID == id
	})
	if taskIndex < 0 {
		return fmt.Errorf("task with id %d not found", id)
	}

	t := &taskList[taskIndex]
	t.Status = string(storage.InProgress)
	t.UpdatedAt = time.Now()
	s.storage.WriteStorage(taskList)
	return nil
}

func (s *taskService) List(status string) ([]storage.Task, error) {
	fileData, err := s.storage.ReadStorage()
	if err != nil {
		return nil, err
	}

	if status == "" {
		return fileData, nil
	}

	return slices.DeleteFunc(fileData, func(t storage.Task) bool {
		return t.Status != status
	}), nil

}

func getMaxId(taskList []storage.Task) int {
	maxID := 0
	for _, t := range taskList {
		if t.ID > maxID {
			maxID = t.ID
		}
	}
	return maxID + 1
}
