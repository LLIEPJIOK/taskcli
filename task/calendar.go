package task

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	weeksInLastSixMonths = 26
	daysInLastSixMonths  = 183
)

func printDayCol(day int) {
	var out string
	switch day {
	case 0:
		out = " Mon "
	case 2:
		out = " Wed "
	case 4:
		out = " Fri "
	case 6:
		out = " Sun "
	default:
		out = "     "
	}

	fmt.Print(out)
}

func printMonths() {
	week := time.Now().Add(-(daysInLastSixMonths * time.Hour * 24))
	month := week.Month()
	fmt.Printf("         ")
	for {
		if week.Month() != month {
			fmt.Printf("%s ", week.Month().String()[:3])
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

func printCell(val int) {
	var cell string
	switch {
	case val > 0 && val < 5:
		cell = "\033[1;30;47m"
	case val >= 5 && val < 10:
		cell = "\033[1;30;43m"
	case val >= 10:
		cell = "\033[1;30;42m"
	default:
		cell = "\033[0;37;30m"
	}
	if val == 0 {
		fmt.Printf(cell + "  - " + "\033[0m")
		return
	}

	var number string
	switch {
	case val >= 10 && val < 100:
		number = " %d "
	case val >= 100:
		number = "%d "
	default:
		number = "  %d "
	}
	fmt.Printf(cell+number+"\033[0m", val)
}

func PrintCells() {
	printMonths()
	for i := 0; i < 7; i++ {
		for j := weeksInLastSixMonths + 1; j >= 0; j-- {
			if j == weeksInLastSixMonths+1 {
				printDayCol(i)
			}
			printCell(rand.Int() % 200)
		}
		fmt.Printf("\n")
	}
}
