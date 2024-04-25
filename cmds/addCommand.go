package cmds

import (
	"taskcli/database"
	"taskcli/task"

	"github.com/spf13/cobra"
)

var addCommand = &cobra.Command{
	Use:   "add NAME",
	Short: "Add a new task with specified name",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		task := task.New(args[0])
		return database.Insert(task)
	},
}
