// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package sqldev

import (
	"database/sql"
)

type Availability struct {
	IDAvailability   int64
	IDProfessional   int64
	InitDatetime     string
	EndDatetime      string
	InitHour         string
	EndHour          string
	TypeAvailability sql.NullInt64
	WeekdayName      string
	Interval         int64
}
