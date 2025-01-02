-- name: InsertProfessional :one
INSERT INTO professional (
  reference_key,
  nome,
  especialidade
) VALUES (
  ?, ?, ?
)
RETURNING *;

-- name: InsertAttribute :one
INSERT INTO attribute (
  id_professional,
  attribute,
  value
) VALUES (
  @id_professional,
  @attribute,
  @value
) RETURNING *;

-- name: InsertAvailability :one
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
RETURNING *;

-- name: ListAvailability :one
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
WHERE id_availability = ? LIMIT 1;

-- name: GetExistingSlot :one
SELECT id_slot
FROM slot s
WHERE 1=1
  AND is_deleted = 0
	AND id_professional = ?
	AND datetime(?) between datetime(slot) and datetime(slot, concat(s."interval" - 1, ' minute'))
    AND priority_entry = ?;

-- name: UpdateSlot :one
UPDATE slot
SET status_entry = @status_entry,
    priority_entry = @priority_entry,
    owner = @owner,
    external_id = @external_id,
    updated_at = CURRENT_TIMESTAMP
WHERE id_slot == @id_slot
RETURNING *;

-- name: InsertSlot :one
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
  @id_professional,
  @id_availability,
  @slot,
  @weekday_name,
  @interval,
  @priority_entry,
  @status_entry,
  @is_deleted,
  @id_blocker
)
RETURNING slot;

-- name: ListSlots :many
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
  AND CASE WHEN @deleted == true THEN 1 ELSE is_deleted == 0 END
  AND CASE WHEN @is_hour == true THEN time(datetime(slot, '-3 hour')) between time(@init_hour) and time(@end_hour) ELSE 1 END
  AND datetime(slot) between datetime(@slot_init) and datetime(@slot_end)
  AND CASE WHEN @is_professional == true THEN p.reference_key in (sqlc.slice('reference_key')) ELSE 1 END
  AND CASE WHEN @is_open == true THEN s.status_entry == 'open' ELSE 1 END
  AND CASE WHEN @is_especialidade == true THEN p.especialidade in (sqlc.slice('especialidade')) ELSE 1 END
  AND CASE WHEN @is_idclinica == true THEN s.id_professional in (
    SELECT a.id_professional FROM attribute a WHERE attribute == 'idclinica' and value in (sqlc.slice('idclinica'))
  ) ELSE 1 END
ORDER BY s.slot;

-- name: ListSlotsByIdAvailability :many
SELECT
  id_slot
FROM slot
WHERE 1=1
  AND is_deleted = 0
  AND id_availability == @id_availability;

-- name: DeleteSlotById :exec
UPDATE slot
SET is_deleted = 1,
  updated_at = CURRENT_TIMESTAMP,
  deleted_at = CURRENT_TIMESTAMP
WHERE id_slot == @id_slot;

-- name: DeleteAvailabilityById :one
UPDATE availability
SET is_deleted = 1
WHERE id_availability == @id_availability
RETURNING *;

-- name: GetProfessionalInfo :one
SELECT
  id_professional,
  reference_key,
  nome,
  especialidade
FROM professional
WHERE reference_key == @reference_key;

-- name: ListAttributesByProfessionalId :many
SELECT
  id_attribute,
  attribute,
  value
FROM attribute
WHERE id_professional == @id_professional;

-- name: ListAvailabilityByProfessionalId :many
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
  AND id_professional == @id_professional
  AND CASE WHEN @deleted == true THEN 1 ELSE is_deleted == 0 END;

-- name: GetSlotById :one
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
  AND id_slot == @id_slot
  AND CASE WHEN @deleted == true THEN 1 ELSE is_deleted == 0 END;

-- name: ListBlockerByProfessional :many
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
  AND id_professional == @id_professional
  AND CASE WHEN @deleted == true THEN 1 ELSE is_deleted ==0 END;

-- name: InsertBlocker :one
INSERT INTO blocker (
  id_professional,
  title,
  description,
  init_datetime,
  end_datetime
) VALUES (
  @id_professional, 
  @title, 
  @description,
  @init_datetime,
  @end_datetime
)
RETURNING *;

-- name: GetBlockerById :one
SELECT
  id_blocker,
  id_professional,
  title,
  init_datetime,
  end_datetime,
  is_deleted
FROM blocker
WHERE 1=1
  AND id_blocker == @id_blocker
  AND CASE WHEN @deleted == true THEN 1 ELSE is_deleted == 0 END;

-- name: DeleteBlockerById :one
UPDATE blocker
SET is_deleted = 1
WHERE id_blocker == @id_blocker
RETURNING *;

-- name: UpdateSlotSetBlocker :many
UPDATE slot
SET status_entry = @status_entry,
	  id_blocker = @id_blocker,
    updated_at = CURRENT_TIMESTAMP
WHERE 1=1
	AND id_professional = @id_professional
  AND slot >= @init_blocker AND slot <= @end_blocker
RETURNING *;

-- name: CreateSlot :one
INSERT INTO slot(
  id_professional,
  id_availability,
  slot,
  weekday_name,
  interval,
  priority_entry,
  status_entry
) VALUES (
  @id_professional,
  @id_availability,
  @slot,
  @weekday_name,
  @interval,
  @priority_entry,
  @status_entry
)
RETURNING *;

-- name: CheckProfessionalExists :one
SELECT 1
FROM professional
WHERE id_professional = @id_professional;