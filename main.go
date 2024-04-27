package main

import (
	"log"

	"github.com/LLIEPJIOK/taskcli/cmds"
	"github.com/LLIEPJIOK/taskcli/database"
)

func main() {
	defer database.Close()
	if err := cmds.RootCommand.Execute(); err != nil {
		log.Fatal("cannot execute root command:", err)
	}
}
