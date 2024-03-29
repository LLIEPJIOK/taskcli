package cmds

import (
	"fmt"
	"strconv"
	"taskcli/database"
	"taskcli/task"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/spf13/cobra"
)

var (
	columns = []string{"ID", "Name", "Status", "Created at"}

	RootCommand = &cobra.Command{
		Use:   "info",
		Short: "A CLI task management tool for creating your to do list",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	addCommand = &cobra.Command{
		Use:   "add NAME",
		Short: "Add a new task with specified name",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			task := task.New(args[0])
			return database.Insert(task)
		},
	}

	deleteCommand = &cobra.Command{
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

	updateCommand = &cobra.Command{
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

func init() {
	RootCommand.AddCommand(addCommand)
	RootCommand.AddCommand(deleteCommand)
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
	RootCommand.AddCommand(updateCommand)
	RootCommand.AddCommand(listCommand)
}
