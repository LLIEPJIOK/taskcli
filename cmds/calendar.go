package cmds

import (
	"fmt"
	"log"
	"taskcli/database"
	"taskcli/task"
	"time"
)

const (
	weeksInLastSixMonths = 26
	daysInLastSixMonths  = 183
)

func printDayCol(day int) {
	var out string
	switch day {
	case 1:
		out = " Mon "
	case 3:
		out = " Wed "
	case 5:
		out = " Fri "
	case 0:
		out = " Sun "
	default:
		out = "     "
	}

	fmt.Print(out)
}

func printMonths() {
	week := time.Now().Add(-(daysInLastSixMonths * time.Hour * 24))
	month := week.Month()
	fmt.Printf("            ")
	for {
		if week.Month() != month {
			fmt.Printf("        %s ", week.Month().String()[:3])
			month = week.Month()
		} else {
			fmt.Printf("    ")
		}

		week = week.Add(7 * time.Hour * 24)
		if week.After(time.Now()) {
			break
		}
	}
	fmt.Printf("\n")
}

func printCell(val, quantity int) {
	if val == -1 {
		fmt.Printf("\033[0;37;30m" + "  - " + "\033[0m")
		return
	}
	var cell string
	switch {
	case quantity == 0:
		cell = "\033[0;37;30m"
	case quantity > 0 && quantity < 5:
		cell = "\033[1;30;45m"
	case quantity >= 5 && quantity < 10:
		cell = "\033[1;30;43m"
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

func PrintCells() {
	dateTo := time.Now()
	for dateTo.Weekday() != time.Sunday {
		dateTo = dateTo.AddDate(0, 0, 1)
	}
	dateFrom := time.Now().AddDate(0, -6, 0)
	for dateFrom.Weekday() != time.Monday {
		dateFrom = dateFrom.AddDate(0, 0, -1)
	}

	printMonths()

	daysQuantity, err := database.GetDateWithQuantity()
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(daysQuantity)
	for i := 0; i < 7; i++ {
		typedDate := dateFrom
		printDayCol(int(typedDate.Weekday()))
		for typedDate.Before(dateTo) {
			if typedDate.Day() <= 7 {
				printCell(-1, -1)
				printCell(-1, -1)
			}
			printCell(typedDate.Day(), daysQuantity[task.Day{Year: typedDate.Year(), Month: int(typedDate.Month()), Day: typedDate.Day()}])
			typedDate = typedDate.AddDate(0, 0, 7)
		}
		dateFrom = dateFrom.AddDate(0, 0, 1)

		fmt.Printf("\n")
	}
}
