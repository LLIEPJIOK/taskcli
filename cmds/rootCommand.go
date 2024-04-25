package cmds

import "github.com/spf13/cobra"

var RootCommand = &cobra.Command{
	Use:   "info",
	Short: "A CLI task management tool for creating your to do list",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	RootCommand.AddCommand(addCommand)
	RootCommand.AddCommand(deleteCommand)
	RootCommand.AddCommand(updateCommand)
	RootCommand.AddCommand(listCommand)
	RootCommand.AddCommand(calendarCommand)
}
