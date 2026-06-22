package mapper

import (
	"fmt"
	"time"

	"github.com/Chibas/TaskTrackerGo/src/storage"
)

type CommandMapper interface {
	ParseCommand(s string) (Command, error)
	ExecuteCommand(cmd Command, data any) error
	LogList(l []storage.Task)
}

type commandMapper struct {
	storage storage.Storage
}

func NewCommandMapper() CommandMapper {
	return &commandMapper{
		storage: storage.NewStorage("tasks.json"),
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

func (c *commandMapper) ExecuteCommand(cmd Command, data any) error {
	switch cmd {
	case Add:
		fmt.Printf("ADD %s", data)
	case Update:
		fmt.Printf("Update %s", data)
	case Delete:
		fmt.Printf("Delete %s", data)
	case MarkInProgress:
		fmt.Printf("MarkInProgress %s", data)
	case MarkDone:
		fmt.Printf("MarkDone %s", data)
	case List:
		fmt.Printf("List %s \n", data)
		fileData, err := c.storage.ReadStorage()
		if err != nil {
			return err
		}
		c.LogList(fileData)
	default:
		return fmt.Errorf("Unknown command %s", cmd)
	}
	return nil
}

func (c *commandMapper) LogList(l []storage.Task) {
	for index, task := range l {
		fmt.Printf("#%d Name: %s Status: %s 🔼 Created: %s  🕐 Updated: %s \n", index, task.Description, task.Status, time.Time(task.CreatedAt).Format("02 Jan 2006 15:04:05"), time.Time(task.UpdatedAt).Format("02 Jan 2006 15:04:05"))
	}
}
