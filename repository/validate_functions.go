package repository

import (
	"regexp"
	"time"
)

func IsValidTypeAvailability(typeAv int64) bool {
	switch typeAv {
	case 0, 2, 3:
		return true
	}
	return false
}

func IsValidDatetime(dt string) bool {
	layout := "2006-01-02 15:04:05"
	_, err := time.Parse(layout, dt)
	return err == nil
}

func IsValidHour(hr string) bool {
	matched, _ := regexp.MatchString(`^[012][0-9]:[012][0-9]$`, hr)
	if matched {
		return true
	} else {
		return false
	}
}
