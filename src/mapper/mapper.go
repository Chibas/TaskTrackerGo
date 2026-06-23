package mapper

import (
	"fmt"
	"strconv"

	"github.com/Chibas/TaskTrackerGo/src/storage"
	"github.com/Chibas/TaskTrackerGo/src/task"
)

type CommandMapper interface {
	ParseCommand(s string) (Command, error)
	ExecuteCommand(cmd Command, arg1 any, arg2 any) error
	LogList(l []storage.Task)
}

type commandMapper struct {
	taskService task.TaskService
}

func NewCommandMapper(taskService task.TaskService) CommandMapper {
	return &commandMapper{
		taskService: taskService,
	}
}

func (c *commandMapper) ParseCommand(s string) (Command, error) {
	switch Command(s) {
	case Add, Update, Delete, MarkInProgress, MarkDone, List:
		return Command(s), nil
	default:
		return "", fmt.Errorf("Unknown command %s", s)
	}
}

func (c *commandMapper) ExecuteCommand(cmd Command, arg1 any, arg2 any) error {
	switch cmd {
	case Add:
		description, ok := arg1.(string)
		if !ok {
			return fmt.Errorf("add requires a description string")
		}
		return c.taskService.AddTask(description)

	case Update:
		idStr, ok := arg1.(string)
		if !ok {
			return fmt.Errorf("update requires an id")
		}
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return fmt.Errorf("invalid id: %w", err)
		}
		description, ok := arg2.(string)
		if !ok {
			return fmt.Errorf("update requires a description string")
		}
		return c.taskService.UpdateTask(id, description)

	case Delete:
		idStr, ok := arg1.(string)
		if !ok {
			return fmt.Errorf("delete requires an id")
		}
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return fmt.Errorf("invalid id: %w", err)
		}
		return c.taskService.DeleteTask(id)

	case MarkInProgress:
		idStr, ok := arg1.(string)
		if !ok {
			return fmt.Errorf("mark-in-progress requires an id")
		}
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return fmt.Errorf("invalid id: %w", err)
		}
		return c.taskService.MarkInProgress(id)

	case MarkDone:
		idStr, ok := arg1.(string)
		if !ok {
			return fmt.Errorf("mark-done requires an id")
		}
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return fmt.Errorf("invalid id: %w", err)
		}
		return c.taskService.MarkDone(id)

	case List:
		fileData, err := c.taskService.List()
		if err != nil {
			return err
		}
		c.LogList(fileData)
		return nil

	default:
		return fmt.Errorf("Unknown command %s", cmd)
	}

}

func (c *commandMapper) LogList(l []storage.Task) {
	for _, task := range l {
		fmt.Printf("#%d | Name: %s | Status: %s | 🔼 Created: %s  🕐 Updated: %s \n", task.ID, task.Description, task.Status, task.CreatedAt.Format("02 Jan 2006 15:04:05"), task.UpdatedAt.Format("02 Jan 2006 15:04:05"))
	}
}
