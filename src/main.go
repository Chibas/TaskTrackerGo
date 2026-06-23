package main

import (
	"fmt"
	"os"

	"github.com/Chibas/TaskTrackerGo/src/mapper"
	"github.com/Chibas/TaskTrackerGo/src/storage"
	"github.com/Chibas/TaskTrackerGo/src/task"
)

func main() {
	argsCmd := os.Args[1:]

	if len(argsCmd) == 0 {
		fmt.Println("no command provided")
		return
	}

	var arg1 any
	var arg2 any
	if len(argsCmd) > 1 {
		arg1 = argsCmd[1]
	}
	if len(argsCmd) > 2 {
		arg2 = argsCmd[2]
	}

	storageSvc := storage.NewStorage("tasks.json")
	taskSvc := task.NewTaskService(storageSvc)
	m := mapper.NewCommandMapper(taskSvc)
	cmd, err := m.ParseCommand(argsCmd[0])

	if err != nil {
		fmt.Println("parse error:", err)
		return
	}

	fmt.Printf("Data %v \n", arg1)

	if err := m.ExecuteCommand(cmd, arg1, arg2); err != nil {
		fmt.Println("execute error:", err)
	}

	if cmd != mapper.List && err == nil {
		m.ExecuteCommand(mapper.Command(mapper.List), nil, nil)
	}
}
