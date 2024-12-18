// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package models

import (
	"time"
)

type Attribute struct {
	IDAttribute    int64  `json:"id_attribute"`
	IDProfessional int64  `json:"id_professional"`
	Attribute      string `json:"attribute"`
	Value          string `json:"value"`
}

type Availability struct {
	IDAvailability   int64  `json:"id_availability"`
	IDProfessional   int64  `json:"id_professional"`
	InitDatetime     string `json:"init_datetime"`
	EndDatetime      string `json:"end_datetime"`
	InitHour         string `json:"init_hour"`
	EndHour          string `json:"end_hour"`
	TypeAvailability int64  `json:"type_availability"`
	WeekdayName      string `json:"weekday_name"`
	Interval         int64  `json:"interval"`
	PriorityEntry    int64  `json:"priority_entry"`
}

type Professional struct {
	IDProfessional int64  `json:"id_professional"`
	Especialidade  string `json:"especialidade"`
	Nome           string `json:"nome"`
}

type Slot struct {
	IDSlot         int64     `json:"id_slot"`
	IDAvailability int64     `json:"id_availability"`
	IDProfessional int64     `json:"id_professional"`
	Slot           time.Time `json:"slot"`
	WeekdayName    string    `json:"weekday_name"`
	Interval       int64     `json:"interval"`
	PriorityEntry  int64     `json:"priority_entry"`
	StatusEntry    string    `json:"status_entry"`
}
