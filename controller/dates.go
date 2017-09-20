package main

import "fmt"
import "time"

func main() {
	// Start date as tomorrow
	startdate := time.Now().Local().AddDate(0, 0, 1)
	fmt.Println("Start date: ", startdate.Format("2006-01-02"))
	// End date as 3 months from tomorrow
	enddate := startdate.AddDate(0, 3, 0)
	fmt.Println("End date: ", enddate.Format("2006-01-02"))

	// Iterate over dates to print all allowed dates
	for d := startdate; d != enddate; d = d.AddDate(0, 0, 1) {
		// Exclude weekends (0 - Sunday, 6 - Saturday)
		if d.Weekday() != 0 {
			fmt.Println("Dates: ", d.Format("2006-01-02"), d.Weekday())
		}
	}
}
