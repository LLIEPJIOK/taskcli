package cmds

import (
	"fmt"
	"taskcli/database"
	"taskcli/task"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/spf13/cobra"
)

var (
	columns     = []string{"ID", "Name", "Status", "Created at"}
	listCommand = &cobra.Command{
		Use:   "list",
		Short: "List all your tasks",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			task, err := database.GetAllTasks()
			if err != nil {
				return err
			}
			fmt.Print(setupTable(task))
			return nil
		},
	}
)

func setupTable(tasks []task.Task) *table.Table {
	var rows [][]string
	for _, task := range tasks {
		rows = append(rows, []string{
			fmt.Sprint(task.ID),
			task.Name,
			task.Status,
			task.CreationTime.Format("2006-01-02"),
		})
	}
	t := table.New().
		Border(lipgloss.HiddenBorder()).
		Headers(columns...).
		Rows(rows...).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == 0 {
				return lipgloss.NewStyle().
					Foreground(lipgloss.Color("212")).
					Border(lipgloss.NormalBorder()).
					BorderTop(false).
					BorderLeft(false).
					BorderRight(false).
					BorderBottom(true).
					Bold(true)
			}
			if row%2 == 0 {
				return lipgloss.NewStyle().Foreground(lipgloss.Color("246"))
			}
			return lipgloss.NewStyle()
		}).Width(80)
	return t
}
