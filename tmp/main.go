package main

import (
	"fmt"
	"log"
	"time"
)

func CalculateWeekdayBetween(start, end time.Time, targetWeekday time.Weekday, typeInterval int) []time.Time {
	start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
	end = time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, end.Location())

	if start.After(end) {
		start, end = end, start
	}

	var days_weekday []time.Time

	flag_biweekly := true

	current := start

	for !current.After(end) {

		if typeInterval == 1 || typeInterval == 2 {
			if current.Weekday() == targetWeekday && flag_biweekly {
				matchingDay := time.Date(
					current.Year(),
					current.Month(),
					current.Day(),
					0, 0, 0, 0, current.Location(),
				)
				days_weekday = append(days_weekday, matchingDay)

				if typeInterval == 2 {
					flag_biweekly = !flag_biweekly
				}
			} else if current.Weekday() == targetWeekday && !flag_biweekly {
				flag_biweekly = !flag_biweekly
			}
		}

		if typeInterval == 3 {
			if current.Weekday() == targetWeekday && len(days_weekday) == 0 {
				matchingDay := time.Date(
					current.Year(),
					current.Month(),
					current.Day(),
					0, 0, 0, 0, current.Location(),
				)
				days_weekday = append(days_weekday, matchingDay)
			} else if current.Weekday() == targetWeekday &&
				len(days_weekday) > 0 &&
				days_weekday[len(days_weekday)-1].Month() != current.Month() {
				matchingDay := time.Date(
					current.Year(),
					current.Month(),
					current.Day(),
					0, 0, 0, 0, current.Location(),
				)
				days_weekday = append(days_weekday, matchingDay)
			}
		}

		current = current.AddDate(0, 0, 1)
	}

	return days_weekday
}

// SplitTimeRange splits a time range into intervals of the specified duration
func SplitTimeRange(hourStart, hourEnd time.Time, interval time.Duration) []time.Time {
	if hourStart.After(hourEnd) {
		return nil
	}

	var slots []time.Time
	// slots = append(slots, current)

	current := hourStart
	// Keep adding intervals until we reach or exceed the end time
	for current.Add(interval).Before(hourEnd) {
		current = current.Add(interval)
		slots = append(slots, current)
	}

	return slots
}

func ComputeAgenda(initialHour, endHour string, daysToCompute []time.Time, duration time.Duration) ([]time.Time, error) {
	layout := "2006-01-02 15:04:05"

	startTime, err := time.Parse(layout, fmt.Sprintf("1970-01-01 %s", initialHour))
	if err != nil {
		log.Fatal(err)
	}

	endTime, err := time.Parse(layout, fmt.Sprintf("1970-01-01 %s", endHour))
	if err != nil {
		log.Fatal(err)
	}

	var list_slots []time.Time

	for _, dayCompute := range daysToCompute {
		tempStartDayTime := time.Date(
			dayCompute.Year(),
			dayCompute.Month(),
			dayCompute.Day(),
			startTime.Hour(),
			startTime.Minute(),
			startTime.Second(),
			startTime.Nanosecond(),
			startTime.Location(),
		)

		tempEndDayTime := time.Date(
			dayCompute.Year(),
			dayCompute.Month(),
			dayCompute.Day(),
			endTime.Hour(),
			endTime.Minute(),
			endTime.Second(),
			endTime.Nanosecond(),
			endTime.Location(),
		)

		slots := SplitTimeRange(tempStartDayTime, tempEndDayTime, duration)

		list_slots = append(list_slots, slots...)

	}

	return list_slots, nil
}

type BasicTime struct {
	StartHour     time.Time
	EndHour       time.Time
	StartDate     time.Time
	EndDate       time.Time
	WeekDayValids time.Weekday
	Interval      time.Duration
	TypeInterval  int
}

func main() {
	// Example usage
	layout := "2006-01-02 15:04:05"
	start, _ := time.Parse(layout, "2024-10-01 00:00:00")
	end, _ := time.Parse(layout, "2024-11-07 11:30:00")

	// // Set interval to 30 minutes
	// interval := 20 * time.Minute

	// slots := SplitTimeRange(start, end, interval)

	// fmt.Println("Time slots:")
	// for i, slot := range slots {
	// 	fmt.Printf("Slot %d: %s\n", i+1, slot.Format(layout))
	// }

	weekDayValids := CalculateWeekdayBetween(start, end, time.Tuesday, 3)

	slots, err := ComputeAgenda("16:00:00", "19:00:00", weekDayValids, 60*time.Minute)
	if err != nil {
		log.Fatal(err)
	}

	for i, slot := range slots {
		fmt.Printf("Slot %d: %s\n", i+1, slot.Format(layout))
	}

	// monthDaysValid := CalculateMonthlyDays(2, start, end)
	// fmt.Println(monthDaysValid)

}
