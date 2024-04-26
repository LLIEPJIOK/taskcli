package cmds

import (
	"fmt"
	"strconv"

	"github.com/LLIEPJIOK/taskcli/database"

	"github.com/spf13/cobra"
)

var deleteCommand = &cobra.Command{
	Use:   "delete ID",
	Short: "Delete a task with specified id",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("cannot parse id: %w", err)
		}
		return database.Delete(uint(id))
	},
}
