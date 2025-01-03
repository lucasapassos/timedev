// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: query.sql

package models

import (
	"context"
	"database/sql"
	"strings"
	"time"
)

const checkProfessionalExists = `-- name: CheckProfessionalExists :one
SELECT 1
FROM professional
WHERE id_professional = ?1
`

func (q *Queries) CheckProfessionalExists(ctx context.Context, idProfessional int64) (int64, error) {
	row := q.db.QueryRowContext(ctx, checkProfessionalExists, idProfessional)
	var column_1 int64
	err := row.Scan(&column_1)
	return column_1, err
}

const createSlot = `-- name: CreateSlot :one
INSERT INTO slot(
  id_professional,
  id_availability,
  slot,
  weekday_name,
  interval,
  priority_entry,
  status_entry
) VALUES (
  ?1,
  ?2,
  ?3,
  ?4,
  ?5,
  ?6,
  ?7
)
RETURNING id_slot, inserted_at, updated_at, id_availability, id_professional, slot, weekday_name, interval, priority_entry, status_entry, external_id, owner, is_deleted, deleted_at, id_blocker
`

type CreateSlotParams struct {
	IDProfessional int64         `json:"id_professional"`
	IDAvailability sql.NullInt64 `json:"id_availability"`
	Slot           time.Time     `json:"slot"`
	WeekdayName    string        `json:"weekday_name"`
	Interval       int64         `json:"interval"`
	PriorityEntry  int64         `json:"priority_entry"`
	StatusEntry    string        `json:"status_entry"`
}

func (q *Queries) CreateSlot(ctx context.Context, arg CreateSlotParams) (Slot, error) {
	row := q.db.QueryRowContext(ctx, createSlot,
		arg.IDProfessional,
		arg.IDAvailability,
		arg.Slot,
		arg.WeekdayName,
		arg.Interval,
		arg.PriorityEntry,
		arg.StatusEntry,
	)
	var i Slot
	err := row.Scan(
		&i.IDSlot,
		&i.InsertedAt,
		&i.UpdatedAt,
		&i.IDAvailability,
		&i.IDProfessional,
		&i.Slot,
		&i.WeekdayName,
		&i.Interval,
		&i.PriorityEntry,
		&i.StatusEntry,
		&i.ExternalID,
		&i.Owner,
		&i.IsDeleted,
		&i.DeletedAt,
		&i.IDBlocker,
	)
	return i, err
}

const deleteAvailabilityById = `-- name: DeleteAvailabilityById :one
UPDATE availability
SET is_deleted = 1
WHERE id_availability == ?1
RETURNING id_availability, id_professional, init_datetime, end_datetime, init_hour, end_hour, type_availability, weekday_name, interval, resting, priority_entry, is_deleted
`

func (q *Queries) DeleteAvailabilityById(ctx context.Context, idAvailability int64) (Availability, error) {
	row := q.db.QueryRowContext(ctx, deleteAvailabilityById, idAvailability)
	var i Availability
	err := row.Scan(
		&i.IDAvailability,
		&i.IDProfessional,
		&i.InitDatetime,
		&i.EndDatetime,
		&i.InitHour,
		&i.EndHour,
		&i.TypeAvailability,
		&i.WeekdayName,
		&i.Interval,
		&i.Resting,
		&i.PriorityEntry,
		&i.IsDeleted,
	)
	return i, err
}

const deleteBlockerById = `-- name: DeleteBlockerById :one
UPDATE blocker
SET is_deleted = 1
WHERE id_blocker == ?1
RETURNING id_blocker, title, description, id_professional, init_datetime, end_datetime, is_deleted
`

func (q *Queries) DeleteBlockerById(ctx context.Context, idBlocker int64) (Blocker, error) {
	row := q.db.QueryRowContext(ctx, deleteBlockerById, idBlocker)
	var i Blocker
	err := row.Scan(
		&i.IDBlocker,
		&i.Title,
		&i.Description,
		&i.IDProfessional,
		&i.InitDatetime,
		&i.EndDatetime,
		&i.IsDeleted,
	)
	return i, err
}

const deleteSlotById = `-- name: DeleteSlotById :exec
UPDATE slot
SET is_deleted = 1,
  updated_at = CURRENT_TIMESTAMP,
  deleted_at = CURRENT_TIMESTAMP
WHERE id_slot == ?1
`

func (q *Queries) DeleteSlotById(ctx context.Context, idSlot int64) error {
	_, err := q.db.ExecContext(ctx, deleteSlotById, idSlot)
	return err
}

const getBlockerById = `-- name: GetBlockerById :one
SELECT
  id_blocker,
  id_professional,
  title,
  init_datetime,
  end_datetime,
  is_deleted
FROM blocker
WHERE 1=1
  AND id_blocker == ?1
  AND CASE WHEN ?2 == true THEN 1 ELSE is_deleted == 0 END
`

type GetBlockerByIdParams struct {
	IDBlocker int64       `json:"id_blocker"`
	Deleted   interface{} `json:"deleted"`
}

type GetBlockerByIdRow struct {
	IDBlocker      int64     `json:"id_blocker"`
	IDProfessional int64     `json:"id_professional"`
	Title          string    `json:"title"`
	InitDatetime   time.Time `json:"init_datetime"`
	EndDatetime    time.Time `json:"end_datetime"`
	IsDeleted      int64     `json:"is_deleted"`
}

func (q *Queries) GetBlockerById(ctx context.Context, arg GetBlockerByIdParams) (GetBlockerByIdRow, error) {
	row := q.db.QueryRowContext(ctx, getBlockerById, arg.IDBlocker, arg.Deleted)
	var i GetBlockerByIdRow
	err := row.Scan(
		&i.IDBlocker,
		&i.IDProfessional,
		&i.Title,
		&i.InitDatetime,
		&i.EndDatetime,
		&i.IsDeleted,
	)
	return i, err
}

const getExistingSlot = `-- name: GetExistingSlot :one
SELECT id_slot
FROM slot s
WHERE 1=1
  AND is_deleted = 0
	AND id_professional = ?
	AND datetime(?) between datetime(slot) and datetime(slot, concat(s."interval" - 1, ' minute'))
    AND priority_entry = ?
`

type GetExistingSlotParams struct {
	IDProfessional int64       `json:"id_professional"`
	Datetime       interface{} `json:"datetime"`
	PriorityEntry  int64       `json:"priority_entry"`
}

func (q *Queries) GetExistingSlot(ctx context.Context, arg GetExistingSlotParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, getExistingSlot, arg.IDProfessional, arg.Datetime, arg.PriorityEntry)
	var id_slot int64
	err := row.Scan(&id_slot)
	return id_slot, err
}

const getProfessionalInfo = `-- name: GetProfessionalInfo :one
SELECT
  id_professional,
  reference_key,
  nome,
  especialidade
FROM professional
WHERE reference_key == ?1
`

type GetProfessionalInfoRow struct {
	IDProfessional int64  `json:"id_professional"`
	ReferenceKey   string `json:"reference_key"`
	Nome           string `json:"nome"`
	Especialidade  string `json:"especialidade"`
}

func (q *Queries) GetProfessionalInfo(ctx context.Context, referenceKey string) (GetProfessionalInfoRow, error) {
	row := q.db.QueryRowContext(ctx, getProfessionalInfo, referenceKey)
	var i GetProfessionalInfoRow
	err := row.Scan(
		&i.IDProfessional,
		&i.ReferenceKey,
		&i.Nome,
		&i.Especialidade,
	)
	return i, err
}

const getSlotById = `-- name: GetSlotById :one
SELECT
  id_slot,
  inserted_at,
  updated_at,
  id_professional,
  id_availability,
  slot,
  weekday_name,
  interval,
  priority_entry,
  status_entry,
  owner,
  external_id,
  is_deleted
FROM slot
WHERE 1=1
  AND id_slot == ?1
  AND CASE WHEN ?2 == true THEN 1 ELSE is_deleted == 0 END
`

type GetSlotByIdParams struct {
	IDSlot  int64       `json:"id_slot"`
	Deleted interface{} `json:"deleted"`
}

type GetSlotByIdRow struct {
	IDSlot         int64          `json:"id_slot"`
	InsertedAt     time.Time      `json:"inserted_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	IDProfessional int64          `json:"id_professional"`
	IDAvailability sql.NullInt64  `json:"id_availability"`
	Slot           time.Time      `json:"slot"`
	WeekdayName    string         `json:"weekday_name"`
	Interval       int64          `json:"interval"`
	PriorityEntry  int64          `json:"priority_entry"`
	StatusEntry    string         `json:"status_entry"`
	Owner          sql.NullString `json:"owner"`
	ExternalID     sql.NullString `json:"external_id"`
	IsDeleted      int64          `json:"is_deleted"`
}

func (q *Queries) GetSlotById(ctx context.Context, arg GetSlotByIdParams) (GetSlotByIdRow, error) {
	row := q.db.QueryRowContext(ctx, getSlotById, arg.IDSlot, arg.Deleted)
	var i GetSlotByIdRow
	err := row.Scan(
		&i.IDSlot,
		&i.InsertedAt,
		&i.UpdatedAt,
		&i.IDProfessional,
		&i.IDAvailability,
		&i.Slot,
		&i.WeekdayName,
		&i.Interval,
		&i.PriorityEntry,
		&i.StatusEntry,
		&i.Owner,
		&i.ExternalID,
		&i.IsDeleted,
	)
	return i, err
}

const insertAttribute = `-- name: InsertAttribute :one
INSERT INTO attribute (
  id_professional,
  attribute,
  value
) VALUES (
  ?1,
  ?2,
  ?3
) RETURNING id_attribute, id_professional, attribute, value
`

type InsertAttributeParams struct {
	IDProfessional int64  `json:"id_professional"`
	Attribute      string `json:"attribute"`
	Value          string `json:"value"`
}

func (q *Queries) InsertAttribute(ctx context.Context, arg InsertAttributeParams) (Attribute, error) {
	row := q.db.QueryRowContext(ctx, insertAttribute, arg.IDProfessional, arg.Attribute, arg.Value)
	var i Attribute
	err := row.Scan(
		&i.IDAttribute,
		&i.IDProfessional,
		&i.Attribute,
		&i.Value,
	)
	return i, err
}

const insertAvailability = `-- name: InsertAvailability :one
INSERT INTO availability (
    id_professional,
    init_datetime,
    end_datetime,
    init_hour,
    end_hour,
    type_availability,
    weekday_name,
    interval,
    resting,
    priority_entry,
    is_deleted
) VALUES (
    ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 0
)
RETURNING id_availability, id_professional, init_datetime, end_datetime, init_hour, end_hour, type_availability, weekday_name, interval, resting, priority_entry, is_deleted
`

type InsertAvailabilityParams struct {
	IDProfessional   int64     `json:"id_professional"`
	InitDatetime     time.Time `json:"init_datetime"`
	EndDatetime      time.Time `json:"end_datetime"`
	InitHour         string    `json:"init_hour"`
	EndHour          string    `json:"end_hour"`
	TypeAvailability int64     `json:"type_availability"`
	WeekdayName      string    `json:"weekday_name"`
	Interval         int64     `json:"interval"`
	Resting          int64     `json:"resting"`
	PriorityEntry    int64     `json:"priority_entry"`
}

func (q *Queries) InsertAvailability(ctx context.Context, arg InsertAvailabilityParams) (Availability, error) {
	row := q.db.QueryRowContext(ctx, insertAvailability,
		arg.IDProfessional,
		arg.InitDatetime,
		arg.EndDatetime,
		arg.InitHour,
		arg.EndHour,
		arg.TypeAvailability,
		arg.WeekdayName,
		arg.Interval,
		arg.Resting,
		arg.PriorityEntry,
	)
	var i Availability
	err := row.Scan(
		&i.IDAvailability,
		&i.IDProfessional,
		&i.InitDatetime,
		&i.EndDatetime,
		&i.InitHour,
		&i.EndHour,
		&i.TypeAvailability,
		&i.WeekdayName,
		&i.Interval,
		&i.Resting,
		&i.PriorityEntry,
		&i.IsDeleted,
	)
	return i, err
}

const insertBlocker = `-- name: InsertBlocker :one
INSERT INTO blocker (
  id_professional,
  title,
  description,
  init_datetime,
  end_datetime
) VALUES (
  ?1, 
  ?2, 
  ?3,
  ?4,
  ?5
)
RETURNING id_blocker, title, description, id_professional, init_datetime, end_datetime, is_deleted
`

type InsertBlockerParams struct {
	IDProfessional int64          `json:"id_professional"`
	Title          string         `json:"title"`
	Description    sql.NullString `json:"description"`
	InitDatetime   time.Time      `json:"init_datetime"`
	EndDatetime    time.Time      `json:"end_datetime"`
}

func (q *Queries) InsertBlocker(ctx context.Context, arg InsertBlockerParams) (Blocker, error) {
	row := q.db.QueryRowContext(ctx, insertBlocker,
		arg.IDProfessional,
		arg.Title,
		arg.Description,
		arg.InitDatetime,
		arg.EndDatetime,
	)
	var i Blocker
	err := row.Scan(
		&i.IDBlocker,
		&i.Title,
		&i.Description,
		&i.IDProfessional,
		&i.InitDatetime,
		&i.EndDatetime,
		&i.IsDeleted,
	)
	return i, err
}

const insertProfessional = `-- name: InsertProfessional :one
INSERT INTO professional (
  reference_key,
  nome,
  especialidade
) VALUES (
  ?, ?, ?
)
RETURNING id_professional, reference_key, especialidade, nome
`

type InsertProfessionalParams struct {
	ReferenceKey  string `json:"reference_key"`
	Nome          string `json:"nome"`
	Especialidade string `json:"especialidade"`
}

func (q *Queries) InsertProfessional(ctx context.Context, arg InsertProfessionalParams) (Professional, error) {
	row := q.db.QueryRowContext(ctx, insertProfessional, arg.ReferenceKey, arg.Nome, arg.Especialidade)
	var i Professional
	err := row.Scan(
		&i.IDProfessional,
		&i.ReferenceKey,
		&i.Especialidade,
		&i.Nome,
	)
	return i, err
}

const insertSlot = `-- name: InsertSlot :one
INSERT INTO slot (
    id_professional,
    id_availability,
    slot,
    weekday_name,
    interval,
    priority_entry,
    status_entry,
    is_deleted,
    id_blocker
) VALUES (
  ?1,
  ?2,
  ?3,
  ?4,
  ?5,
  ?6,
  ?7,
  ?8,
  ?9
)
RETURNING slot
`

type InsertSlotParams struct {
	IDProfessional int64         `json:"id_professional"`
	IDAvailability sql.NullInt64 `json:"id_availability"`
	Slot           time.Time     `json:"slot"`
	WeekdayName    string        `json:"weekday_name"`
	Interval       int64         `json:"interval"`
	PriorityEntry  int64         `json:"priority_entry"`
	StatusEntry    string        `json:"status_entry"`
	IsDeleted      int64         `json:"is_deleted"`
	IDBlocker      sql.NullInt64 `json:"id_blocker"`
}

func (q *Queries) InsertSlot(ctx context.Context, arg InsertSlotParams) (time.Time, error) {
	row := q.db.QueryRowContext(ctx, insertSlot,
		arg.IDProfessional,
		arg.IDAvailability,
		arg.Slot,
		arg.WeekdayName,
		arg.Interval,
		arg.PriorityEntry,
		arg.StatusEntry,
		arg.IsDeleted,
		arg.IDBlocker,
	)
	var slot time.Time
	err := row.Scan(&slot)
	return slot, err
}

const listAttributesByProfessionalId = `-- name: ListAttributesByProfessionalId :many
SELECT
  id_attribute,
  attribute,
  value
FROM attribute
WHERE id_professional == ?1
`

type ListAttributesByProfessionalIdRow struct {
	IDAttribute int64  `json:"id_attribute"`
	Attribute   string `json:"attribute"`
	Value       string `json:"value"`
}

func (q *Queries) ListAttributesByProfessionalId(ctx context.Context, idProfessional int64) ([]ListAttributesByProfessionalIdRow, error) {
	rows, err := q.db.QueryContext(ctx, listAttributesByProfessionalId, idProfessional)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListAttributesByProfessionalIdRow
	for rows.Next() {
		var i ListAttributesByProfessionalIdRow
		if err := rows.Scan(&i.IDAttribute, &i.Attribute, &i.Value); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listAvailability = `-- name: ListAvailability :one
SELECT 
    id_availability,
    id_professional,
    init_datetime,
    end_datetime,
    init_hour,
    end_hour,
    type_availability,
    weekday_name,
    interval,
    priority_entry,
    is_deleted
FROM availability
WHERE id_availability = ? LIMIT 1
`

type ListAvailabilityRow struct {
	IDAvailability   int64     `json:"id_availability"`
	IDProfessional   int64     `json:"id_professional"`
	InitDatetime     time.Time `json:"init_datetime"`
	EndDatetime      time.Time `json:"end_datetime"`
	InitHour         string    `json:"init_hour"`
	EndHour          string    `json:"end_hour"`
	TypeAvailability int64     `json:"type_availability"`
	WeekdayName      string    `json:"weekday_name"`
	Interval         int64     `json:"interval"`
	PriorityEntry    int64     `json:"priority_entry"`
	IsDeleted        int64     `json:"is_deleted"`
}

func (q *Queries) ListAvailability(ctx context.Context, idAvailability int64) (ListAvailabilityRow, error) {
	row := q.db.QueryRowContext(ctx, listAvailability, idAvailability)
	var i ListAvailabilityRow
	err := row.Scan(
		&i.IDAvailability,
		&i.IDProfessional,
		&i.InitDatetime,
		&i.EndDatetime,
		&i.InitHour,
		&i.EndHour,
		&i.TypeAvailability,
		&i.WeekdayName,
		&i.Interval,
		&i.PriorityEntry,
		&i.IsDeleted,
	)
	return i, err
}

const listAvailabilityByProfessionalId = `-- name: ListAvailabilityByProfessionalId :many
SELECT
  id_availability,
  init_datetime,
  end_datetime,
  init_hour,
  end_hour,
  type_availability,
  weekday_name,
  interval,
  priority_entry,
  is_deleted
FROM availability
WHERE 1=1
  AND id_professional == ?1
  AND CASE WHEN ?2 == true THEN 1 ELSE is_deleted == 0 END
`

type ListAvailabilityByProfessionalIdParams struct {
	IDProfessional int64       `json:"id_professional"`
	Deleted        interface{} `json:"deleted"`
}

type ListAvailabilityByProfessionalIdRow struct {
	IDAvailability   int64     `json:"id_availability"`
	InitDatetime     time.Time `json:"init_datetime"`
	EndDatetime      time.Time `json:"end_datetime"`
	InitHour         string    `json:"init_hour"`
	EndHour          string    `json:"end_hour"`
	TypeAvailability int64     `json:"type_availability"`
	WeekdayName      string    `json:"weekday_name"`
	Interval         int64     `json:"interval"`
	PriorityEntry    int64     `json:"priority_entry"`
	IsDeleted        int64     `json:"is_deleted"`
}

func (q *Queries) ListAvailabilityByProfessionalId(ctx context.Context, arg ListAvailabilityByProfessionalIdParams) ([]ListAvailabilityByProfessionalIdRow, error) {
	rows, err := q.db.QueryContext(ctx, listAvailabilityByProfessionalId, arg.IDProfessional, arg.Deleted)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListAvailabilityByProfessionalIdRow
	for rows.Next() {
		var i ListAvailabilityByProfessionalIdRow
		if err := rows.Scan(
			&i.IDAvailability,
			&i.InitDatetime,
			&i.EndDatetime,
			&i.InitHour,
			&i.EndHour,
			&i.TypeAvailability,
			&i.WeekdayName,
			&i.Interval,
			&i.PriorityEntry,
			&i.IsDeleted,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listBlockerByProfessional = `-- name: ListBlockerByProfessional :many
SELECT
  id_blocker,
  id_professional,
  title,
  description,
  init_datetime,
  end_datetime,
  is_deleted
FROM blocker
WHERE 1=1
  AND id_professional == ?1
  AND CASE WHEN ?2 == true THEN 1 ELSE is_deleted ==0 END
`

type ListBlockerByProfessionalParams struct {
	IDProfessional int64       `json:"id_professional"`
	Deleted        interface{} `json:"deleted"`
}

type ListBlockerByProfessionalRow struct {
	IDBlocker      int64          `json:"id_blocker"`
	IDProfessional int64          `json:"id_professional"`
	Title          string         `json:"title"`
	Description    sql.NullString `json:"description"`
	InitDatetime   time.Time      `json:"init_datetime"`
	EndDatetime    time.Time      `json:"end_datetime"`
	IsDeleted      int64          `json:"is_deleted"`
}

func (q *Queries) ListBlockerByProfessional(ctx context.Context, arg ListBlockerByProfessionalParams) ([]ListBlockerByProfessionalRow, error) {
	rows, err := q.db.QueryContext(ctx, listBlockerByProfessional, arg.IDProfessional, arg.Deleted)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListBlockerByProfessionalRow
	for rows.Next() {
		var i ListBlockerByProfessionalRow
		if err := rows.Scan(
			&i.IDBlocker,
			&i.IDProfessional,
			&i.Title,
			&i.Description,
			&i.InitDatetime,
			&i.EndDatetime,
			&i.IsDeleted,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listSlots = `-- name: ListSlots :many
SELECT
  s.id_slot,
  s.status_entry,
  s.inserted_at,
  s.updated_at,
  p.reference_key,
  s.id_availability,
  s.slot,
  p.especialidade,
  s.weekday_name,
  s.interval,
  s.priority_entry,
  s.owner,
  s.external_id,
  s.is_deleted,
  s.deleted_at,
  s.id_blocker
FROM slot s
LEFT JOIN professional p on s.id_professional = p.id_professional
WHERE 1=1
  AND CASE WHEN ?1 == true THEN 1 ELSE is_deleted == 0 END
  AND CASE WHEN ?2 == true THEN time(datetime(slot, '-3 hour')) between time(?3) and time(?4) ELSE 1 END
  AND datetime(slot) between datetime(?5) and datetime(?6)
  AND CASE WHEN ?7 == true THEN p.reference_key in (/*SLICE:reference_key*/?) ELSE 1 END
  AND CASE WHEN ?9 == true THEN s.status_entry == 'open' ELSE 1 END
  AND CASE WHEN ?10 == true THEN p.especialidade in (/*SLICE:especialidade*/?) ELSE 1 END
  AND CASE WHEN ?12 == true THEN s.id_professional in (
    SELECT a.id_professional FROM attribute a WHERE attribute == 'idclinica' and value in (/*SLICE:idclinica*/?)
  ) ELSE 1 END
ORDER BY s.slot
`

type ListSlotsParams struct {
	Deleted         interface{} `json:"deleted"`
	IsHour          interface{} `json:"is_hour"`
	InitHour        interface{} `json:"init_hour"`
	EndHour         interface{} `json:"end_hour"`
	SlotInit        interface{} `json:"slot_init"`
	SlotEnd         interface{} `json:"slot_end"`
	IsProfessional  interface{} `json:"is_professional"`
	ReferenceKey    []string    `json:"reference_key"`
	IsOpen          interface{} `json:"is_open"`
	IsEspecialidade interface{} `json:"is_especialidade"`
	Especialidade   []string    `json:"especialidade"`
	IsIdclinica     interface{} `json:"is_idclinica"`
	Idclinica       []string    `json:"idclinica"`
}

type ListSlotsRow struct {
	IDSlot         int64          `json:"id_slot"`
	StatusEntry    string         `json:"status_entry"`
	InsertedAt     time.Time      `json:"inserted_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	ReferenceKey   sql.NullString `json:"reference_key"`
	IDAvailability sql.NullInt64  `json:"id_availability"`
	Slot           time.Time      `json:"slot"`
	Especialidade  sql.NullString `json:"especialidade"`
	WeekdayName    string         `json:"weekday_name"`
	Interval       int64          `json:"interval"`
	PriorityEntry  int64          `json:"priority_entry"`
	Owner          sql.NullString `json:"owner"`
	ExternalID     sql.NullString `json:"external_id"`
	IsDeleted      int64          `json:"is_deleted"`
	DeletedAt      sql.NullTime   `json:"deleted_at"`
	IDBlocker      sql.NullInt64  `json:"id_blocker"`
}

func (q *Queries) ListSlots(ctx context.Context, arg ListSlotsParams) ([]ListSlotsRow, error) {
	query := listSlots
	var queryParams []interface{}
	queryParams = append(queryParams, arg.Deleted)
	queryParams = append(queryParams, arg.IsHour)
	queryParams = append(queryParams, arg.InitHour)
	queryParams = append(queryParams, arg.EndHour)
	queryParams = append(queryParams, arg.SlotInit)
	queryParams = append(queryParams, arg.SlotEnd)
	queryParams = append(queryParams, arg.IsProfessional)
	if len(arg.ReferenceKey) > 0 {
		for _, v := range arg.ReferenceKey {
			queryParams = append(queryParams, v)
		}
		query = strings.Replace(query, "/*SLICE:reference_key*/?", strings.Repeat(",?", len(arg.ReferenceKey))[1:], 1)
	} else {
		query = strings.Replace(query, "/*SLICE:reference_key*/?", "NULL", 1)
	}
	queryParams = append(queryParams, arg.IsOpen)
	queryParams = append(queryParams, arg.IsEspecialidade)
	if len(arg.Especialidade) > 0 {
		for _, v := range arg.Especialidade {
			queryParams = append(queryParams, v)
		}
		query = strings.Replace(query, "/*SLICE:especialidade*/?", strings.Repeat(",?", len(arg.Especialidade))[1:], 1)
	} else {
		query = strings.Replace(query, "/*SLICE:especialidade*/?", "NULL", 1)
	}
	queryParams = append(queryParams, arg.IsIdclinica)
	if len(arg.Idclinica) > 0 {
		for _, v := range arg.Idclinica {
			queryParams = append(queryParams, v)
		}
		query = strings.Replace(query, "/*SLICE:idclinica*/?", strings.Repeat(",?", len(arg.Idclinica))[1:], 1)
	} else {
		query = strings.Replace(query, "/*SLICE:idclinica*/?", "NULL", 1)
	}
	rows, err := q.db.QueryContext(ctx, query, queryParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListSlotsRow
	for rows.Next() {
		var i ListSlotsRow
		if err := rows.Scan(
			&i.IDSlot,
			&i.StatusEntry,
			&i.InsertedAt,
			&i.UpdatedAt,
			&i.ReferenceKey,
			&i.IDAvailability,
			&i.Slot,
			&i.Especialidade,
			&i.WeekdayName,
			&i.Interval,
			&i.PriorityEntry,
			&i.Owner,
			&i.ExternalID,
			&i.IsDeleted,
			&i.DeletedAt,
			&i.IDBlocker,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listSlotsByIdAvailability = `-- name: ListSlotsByIdAvailability :many
SELECT
  id_slot
FROM slot
WHERE 1=1
  AND is_deleted = 0
  AND id_availability == ?1
`

func (q *Queries) ListSlotsByIdAvailability(ctx context.Context, idAvailability sql.NullInt64) ([]int64, error) {
	rows, err := q.db.QueryContext(ctx, listSlotsByIdAvailability, idAvailability)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []int64
	for rows.Next() {
		var id_slot int64
		if err := rows.Scan(&id_slot); err != nil {
			return nil, err
		}
		items = append(items, id_slot)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateSlot = `-- name: UpdateSlot :one
UPDATE slot
SET status_entry = ?1,
    priority_entry = ?2,
    owner = ?3,
    external_id = ?4,
    updated_at = CURRENT_TIMESTAMP
WHERE id_slot == ?5
RETURNING id_slot, inserted_at, updated_at, id_availability, id_professional, slot, weekday_name, interval, priority_entry, status_entry, external_id, owner, is_deleted, deleted_at, id_blocker
`

type UpdateSlotParams struct {
	StatusEntry   string         `json:"status_entry"`
	PriorityEntry int64          `json:"priority_entry"`
	Owner         sql.NullString `json:"owner"`
	ExternalID    sql.NullString `json:"external_id"`
	IDSlot        int64          `json:"id_slot"`
}

func (q *Queries) UpdateSlot(ctx context.Context, arg UpdateSlotParams) (Slot, error) {
	row := q.db.QueryRowContext(ctx, updateSlot,
		arg.StatusEntry,
		arg.PriorityEntry,
		arg.Owner,
		arg.ExternalID,
		arg.IDSlot,
	)
	var i Slot
	err := row.Scan(
		&i.IDSlot,
		&i.InsertedAt,
		&i.UpdatedAt,
		&i.IDAvailability,
		&i.IDProfessional,
		&i.Slot,
		&i.WeekdayName,
		&i.Interval,
		&i.PriorityEntry,
		&i.StatusEntry,
		&i.ExternalID,
		&i.Owner,
		&i.IsDeleted,
		&i.DeletedAt,
		&i.IDBlocker,
	)
	return i, err
}

const updateSlotSetBlocker = `-- name: UpdateSlotSetBlocker :many
UPDATE slot
SET status_entry = ?1,
	  id_blocker = ?2,
    updated_at = CURRENT_TIMESTAMP
WHERE 1=1
	AND id_professional = ?3
  AND slot >= ?4 AND slot <= ?5
RETURNING id_slot, inserted_at, updated_at, id_availability, id_professional, slot, weekday_name, interval, priority_entry, status_entry, external_id, owner, is_deleted, deleted_at, id_blocker
`

type UpdateSlotSetBlockerParams struct {
	StatusEntry    string        `json:"status_entry"`
	IDBlocker      sql.NullInt64 `json:"id_blocker"`
	IDProfessional int64         `json:"id_professional"`
	InitBlocker    time.Time     `json:"init_blocker"`
	EndBlocker     time.Time     `json:"end_blocker"`
}

func (q *Queries) UpdateSlotSetBlocker(ctx context.Context, arg UpdateSlotSetBlockerParams) ([]Slot, error) {
	rows, err := q.db.QueryContext(ctx, updateSlotSetBlocker,
		arg.StatusEntry,
		arg.IDBlocker,
		arg.IDProfessional,
		arg.InitBlocker,
		arg.EndBlocker,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Slot
	for rows.Next() {
		var i Slot
		if err := rows.Scan(
			&i.IDSlot,
			&i.InsertedAt,
			&i.UpdatedAt,
			&i.IDAvailability,
			&i.IDProfessional,
			&i.Slot,
			&i.WeekdayName,
			&i.Interval,
			&i.PriorityEntry,
			&i.StatusEntry,
			&i.ExternalID,
			&i.Owner,
			&i.IsDeleted,
			&i.DeletedAt,
			&i.IDBlocker,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
