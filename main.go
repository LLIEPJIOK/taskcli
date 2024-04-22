package main

import (
	"log"
	"taskcli/cmds"
	"taskcli/database"
	"taskcli/task"
)

func main() {
	defer database.CloseDB()
	task.PrintCells()
	if err := cmds.RootCommand.Execute(); err != nil {
		log.Fatal("cannot execute root command:", err)
	}
}
