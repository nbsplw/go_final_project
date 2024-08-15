package pkg

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

const timeFormat = "20060102"

func NextDate(now time.Time, date string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("repeat is required")
	}

	t, err := time.Parse(timeFormat, date)
	if err != nil {
		return "", err
	}

	switch repeat[0] {
	case 'y':
		return addYear(t, now)
	case 'd':
		return addDays(t, now, repeat)
	case 'w':
		return AddWeeks(t, now, repeat)
	case 'm':
		return AddMonths(t, now, repeat)
	}
	return "", nil
}

func addYear(t time.Time, now time.Time) (string, error) {
	for {
		t = t.AddDate(1, 0, 0)
		if t.After(now) {
			break
		}
	}
	return t.Format(timeFormat), nil
}

func addDays(t time.Time, now time.Time, repeat string) (string, error) {
	rep := strings.Split(repeat, " ")
	if len(rep) != 2 {
		return "", errors.New("no interval in query")
	}
	daysNumber, err := strconv.Atoi(rep[1])
	if err != nil {
		return "", err
	}
	if daysNumber > 400 || daysNumber < 1 {
		return "", errors.New("incorrect interval")
	} else {
		for {
			t = t.AddDate(0, 0, daysNumber)
			if t.After(now) {
				break
			}
		}
		return t.Format(timeFormat), nil
	}
}

func AddWeeks(t time.Time, now time.Time, repeat string) (string, error) {
	var weekDays []int
	if len(repeat) < 3 {
		return "", errors.New("days interval is empty")
	}

	weekDaysStr := strings.Split(repeat[2:], ",")
	if len(weekDaysStr) == 0 {
		return "", errors.New("weeks interval is empty")
	}

	for _, day := range weekDaysStr {
		dayNumber, err := strconv.Atoi(day)
		if err != nil {
			return "", err
		}
		if dayNumber < 1 || dayNumber > 7 {
			return "", fmt.Errorf("incorrect day: %d", dayNumber)
		}
		weekDays = append(weekDays, dayNumber)
	}

	for i, day := range weekDays {
		if day == 7 {
			weekDays[i] = 0
		}
	}

	sort.Ints(weekDays)
	var nextWeekDay int
	for _, wd := range weekDays {
		if wd >= int(t.Weekday()) {
			nextWeekDay = wd
			break
		}
	}
	if nextWeekDay == 0 {
		nextWeekDay = weekDays[0]
	}
	for {
		t = t.AddDate(0, 0, 1)
		if t.After(now) && int(t.Weekday()) == nextWeekDay {
			return t.Format(timeFormat), nil
		}
	}
}

func AddMonths(t time.Time, now time.Time, repeat string) (string, error) {
	var listOfDays time.Time
	var err error
	if len(repeat) < 3 {
		return "", errors.New("incorrect month interval")
	}

	repSlice := strings.Split(repeat, " ")

	if len(repSlice) == 2 {
		listOfDays, err = getNextDate(now, t.Format(timeFormat), repeat)
		if err != nil {
			return "", err
		}
	} else if len(repSlice) == 3 {
		listOfDays = getNextMonthDate(now, t.Format(timeFormat), repeat)
	} else {
		return "", errors.New("incorrect repeat interval")
	}

	return listOfDays.Format(timeFormat), nil
}

func getNextMonthDate(now time.Time, target string, rule string) time.Time {
	ruleParts := strings.Split(rule, " ")
	daysPart := strings.Split(ruleParts[1], ",")
	monthsPart := strings.Split(ruleParts[2], ",")

	days := make([]int, len(daysPart))
	for i, day := range daysPart {
		days[i], _ = strconv.Atoi(day)
	}

	months := make([]int, len(monthsPart))
	for i, month := range monthsPart {
		months[i], _ = strconv.Atoi(month)
	}

	targetTime, _ := time.Parse(timeFormat, target)
	var nearestDate time.Time

	for {
		for _, day := range days {
			for _, month := range months {
				var date time.Time
				if day < 0 {

					endOfMonth := time.Date(targetTime.Year(), time.Month(month+1), 0, 0, 0, 0, 0, time.UTC)
					date = endOfMonth.AddDate(0, 0, day+1)
				} else {
					date = time.Date(targetTime.Year(), time.Month(month), day, 0, 0, 0, 0, time.UTC)
				}

				if date.Before(now) {
					date = date.AddDate(1, 0, 0)
				}

				if (nearestDate.IsZero() || date.Before(nearestDate)) && date.After(now) {
					nearestDate = date
				}
			}
		}
		if !nearestDate.IsZero() {
			break
		}

		targetTime = targetTime.AddDate(1, 0, 0)
	}

	return nearestDate
}

func getNextDate(now time.Time, target string, rule string) (time.Time, error) {

	targetTime, _ := time.Parse(timeFormat, target)

	ruleParts := strings.SplitN(rule, " ", 2)
	if len(ruleParts) != 2 {
		return time.Time{}, fmt.Errorf("invalid rule format")
	}

	var monthDays []int
	if strings.Contains(ruleParts[1], ",") {
		daysParts := strings.Split(ruleParts[1], ",")
		monthDays = make([]int, len(daysParts))
		for i, part := range daysParts {
			day, err := strconv.Atoi(strings.TrimSpace(part))
			if err != nil {
				return time.Time{}, fmt.Errorf("invalid day in rule: %v", err)
			}
			if day < -2 || day > 31 {
				return time.Time{}, fmt.Errorf("incorrect days interval")
			}
			monthDays[i] = day
		}
	} else {
		day, err := strconv.Atoi(strings.TrimSpace(ruleParts[1]))
		if err != nil {
			return time.Time{}, fmt.Errorf("invalid day in rule: %v", err)
		}
		if day < -2 || day > 31 {
			return time.Time{}, fmt.Errorf("incorrect repeat interval")
		}
		monthDays = []int{day}
	}

	var nearestDate time.Time
	for {
		for _, day := range monthDays {
			var date time.Time
			endOfMonth := time.Date(targetTime.Year(), targetTime.Month()+1, 0, 0, 0, 0, 0, time.UTC)
			if day < 0 {

				date = endOfMonth.AddDate(0, 0, day+1)
			} else {
				if day > endOfMonth.Day() {
					targetTime = targetTime.AddDate(0, 1, 0)
					endOfMonth = time.Date(targetTime.Year(), targetTime.Month()+1, 0, 0, 0, 0, 0, time.UTC)
				}
				date = time.Date(targetTime.Year(), targetTime.Month(), day, 0, 0, 0, 0, time.UTC)
			}

			if date.Before(now) || date.After(endOfMonth) {
				date = date.AddDate(0, 1, 0)
			}
			if (nearestDate.IsZero() || date.Before(nearestDate)) && date.After(now) {
				nearestDate = date
				break
			}
		}
		if !nearestDate.IsZero() {
			break
		}
		targetTime = targetTime.AddDate(0, 1, 0)
	}

	return nearestDate, nil
}
