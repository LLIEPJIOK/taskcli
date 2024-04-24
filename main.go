package main

import (
	"log"
	"taskcli/cmds"
	"taskcli/database"
)

func main() {
	defer database.CloseDB()
	cmds.PrintCells()
	if err := cmds.RootCommand.Execute(); err != nil {
		log.Fatal("cannot execute root command:", err)
	}
}
