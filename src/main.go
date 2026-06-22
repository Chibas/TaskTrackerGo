package main

import (
	"fmt"
	"os"

	"github.com/Chibas/TaskTrackerGo/src/mapper"
)

func main() {
	var data string
	argsCmd := os.Args[1:]

	if len(argsCmd) == 0 {
		fmt.Println("no command provided")
		return
	}

	if len(argsCmd) > 1 {
		data = argsCmd[1]
	}

	m := mapper.NewCommandMapper()
	cmd, err := m.ParseCommand(argsCmd[0])

	if err != nil {
		fmt.Println("parse error:", err)
		return
	}

	fmt.Printf("Data %s \n", data)

	m.ExecuteCommand(cmd, data)
}
