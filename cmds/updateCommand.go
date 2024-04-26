package cmds

import (
	"strconv"
	"time"

	"github.com/LLIEPJIOK/taskcli/database"
	"github.com/LLIEPJIOK/taskcli/task"

	"github.com/spf13/cobra"
)

var updateCommand = &cobra.Command{
	Use:   "update ID",
	Short: "Delete a task with specified id",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			return err
		}
		prog, err := cmd.Flags().GetInt("status")
		if err != nil {
			return err
		}
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		var status string
		switch prog {
		case int(task.InProgress):
			status = task.InProgress.String()
		case int(task.Done):
			status = task.Done.String()
		default:
			status = task.ToDo.String()
		}
		newTask := task.Task{ID: uint(id), Name: name, Status: status, CreationTime: time.Time{}}
		return database.Update(&newTask)
	},
}

func init() {
	updateCommand.Flags().StringP(
		"name",
		"n",
		"",
		"specify a name for your task",
	)
	updateCommand.Flags().IntP(
		"status",
		"s",
		int(task.ToDo),
		"specify a status for your task",
	)
}
