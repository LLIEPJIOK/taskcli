package cmds

import (
	"fmt"
	"log"
	"time"

	"github.com/LLIEPJIOK/taskcli/database"
	"github.com/LLIEPJIOK/taskcli/task"

	"github.com/spf13/cobra"
)

var calendarCommand = &cobra.Command{
	Use:   "calendar",
	Short: "Show calendar with statics on your task for last 6 month",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		printCells()
		return nil
	},
}

func printDayCol(day int) {
	var out string
	switch day {
	case 2:
		out = " Tue "
	case 4:
		out = " Thu "
	case 6:
		out = " Fri "
	default:
		out = "     "
	}

	fmt.Print(out)
}

func printCell(val, quantity int) {
	if val == -1 {
		fmt.Printf("    ")
		return
	}
	var cell string
	switch {
	case quantity > 0 && quantity < 5:
		cell = "\033[1;30;47m"
	case quantity >= 5 && quantity < 10:
		cell = "\033[1;30;46m"
	case quantity >= 10:
		cell = "\033[1;30;42m"
	default:
		cell = "\033[0;37;30m"
	}

	number := "  %d "
	if val >= 10 {
		number = " %d "
	}
	fmt.Printf(cell+number+"\033[0m", val)
}

func printCells() {
	dateTo := time.Now()
	for dateTo.Weekday() != time.Sunday {
		dateTo = dateTo.AddDate(0, 0, 1)
	}
	dateFrom := time.Now().AddDate(0, -6, 0)
	for dateFrom.Weekday() != time.Monday {
		dateFrom = dateFrom.AddDate(0, 0, -1)
	}
	sixMonthAgoDate := time.Now().AddDate(0, -6, 0)

	daysQuantity, err := database.GetDateWithQuantity()
	if err != nil {
		log.Println(err)
		return
	}

	// print months
	monthTypedDate := dateFrom.AddDate(0, 0, 6)
	month := sixMonthAgoDate.AddDate(0, -1, 0).Month()
	fmt.Print("     ")
	for monthTypedDate.Before(dateTo) {
		if monthTypedDate.Day() <= 7 && monthTypedDate.Month() != sixMonthAgoDate.Month() {
			fmt.Print("        ")
		}
		if month != monthTypedDate.Month() {
			month = monthTypedDate.Month()
			fmt.Printf(" %s ", month.String()[:3])
		} else {
			fmt.Print("    ")
		}
		monthTypedDate = monthTypedDate.AddDate(0, 0, 7)
	}
	fmt.Println()

	// print days
	for i := 0; i < 7; i++ {
		typedDate := dateFrom
		printDayCol(int(typedDate.Weekday()))
		for typedDate.Before(dateTo) {
			if typedDate.Day() <= 7 && typedDate.Month() != sixMonthAgoDate.Month() {
				printCell(-1, -1)
				printCell(-1, -1)
			}
			if !typedDate.After(sixMonthAgoDate) || typedDate.After(time.Now()) {
				printCell(-1, -1)
			} else {
				printCell(typedDate.Day(), daysQuantity[task.Day{Year: typedDate.Year(), Month: int(typedDate.Month()), Day: typedDate.Day()}])
			}
			typedDate = typedDate.AddDate(0, 0, 7)
		}
		dateFrom = dateFrom.AddDate(0, 0, 1)
		fmt.Println()
	}
	fmt.Println()
}
