package repository

import (
	"errors"
	"fmt"
	"log"
	"time"
)

func computeWeekdayName(weekday string) (time.Weekday, error) {
	var weekdayName time.Weekday
	switch caseValue := weekday; caseValue {
	case "Monday":
		weekdayName = time.Monday
	case "Tuesday":
		weekdayName = time.Tuesday
	case "Friday":
		weekdayName = time.Friday
	case "Thursday":
		weekdayName = time.Thursday
	case "Wednesday":
		weekdayName = time.Wednesday
	case "Sunday":
		weekdayName = time.Sunday
	case "Saturday":
		weekdayName = time.Saturday
	default:
		return time.Monday, errors.New("weekday name not found")
	}

	return weekdayName, nil
}

func ComputeSlots(startDatetime, endDatetime time.Time, weekday string, intervalDuration, typeavailability int, hour_init, hour_end string) ([]time.Time, error) {
	weekdayName, err := computeWeekdayName(weekday)
	if err != nil {
		return nil, err
	}

	weekDayValids := CalculateWeekdayBetween(startDatetime, endDatetime, weekdayName, typeavailability)

	slots, err := ComputeAgenda(hour_init, hour_end, weekDayValids, time.Duration(intervalDuration)*time.Minute)
	if err != nil {
		return nil, err
	}

	return slots, nil
}

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

		if typeInterval == 0 || typeInterval == 2 {
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

	// TODO Add to function the power to use inconsecutives slots (like, 45 (atd) + 15 (break) = 1h) or for every 2 sessions, skip the third.

	var slots []time.Time
	current := hourStart
	slots = append(slots, current)

	// Keep adding intervals until we reach or exceed the end time
	for current.Add(interval).Before(hourEnd) {
		current = current.Add(interval)
		slots = append(slots, current)
	}

	return slots
}

func ComputeAgenda(initialHour, endHour string, daysToCompute []time.Time, duration time.Duration) ([]time.Time, error) {
	layout := "2006-01-02 15:04:05"

	startTime, err := time.Parse(layout, fmt.Sprintf("1970-01-01 %s:00", initialHour))
	if err != nil {
		log.Fatal(err)
	}

	endTime, err := time.Parse(layout, fmt.Sprintf("1970-01-01 %s:00", endHour))
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
			dayCompute.Location(),
		)

		tempEndDayTime := time.Date(
			dayCompute.Year(),
			dayCompute.Month(),
			dayCompute.Day(),
			endTime.Hour(),
			endTime.Minute(),
			endTime.Second(),
			endTime.Nanosecond(),
			dayCompute.Location(),
		)

		slots := SplitTimeRange(tempStartDayTime, tempEndDayTime, duration)

		list_slots = append(list_slots, slots...)

	}

	return list_slots, nil
}
