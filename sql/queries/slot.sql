
-- name: GetExistingSlot :one
SELECT id_slot
FROM slot s
WHERE 1=1
  AND is_deleted = false
	AND id_professional = @id_professional
	AND @slot between slot and slot + (INTERVAL '1 min' * (interval -1))
    AND priority_entry = @priority_entry;

-- name: UpdateSlot :one
UPDATE slot
SET status_entry = @status_entry,
    priority_entry = @priority_entry,
    owner = @owner,
    id_external = @id_external,
    updated_at = CURRENT_TIMESTAMP
WHERE id_slot = @id_slot
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
    id_blocker
) VALUES (
  @id_professional,
  @id_availability,
  @slot,
  @weekday_name,
  @interval,
  @priority_entry,
  @status_entry,
  @id_blocker
)
RETURNING id_slot, slot;

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
  s.id_external,
  s.is_deleted,
  s.deleted_at,
  s.id_blocker
FROM slot s
LEFT JOIN professional p on s.id_professional = p.id_professional
WHERE 1=1
  AND CASE WHEN @deleted = true THEN true ELSE is_deleted = false END
  AND CASE WHEN @is_hour = true THEN cast(concat(extract(hour from slot), ':', extract(minute from slot)) as time) between cast(@init_hour::varchar as time) and cast(@end_hour::varchar as time) ELSE true END
  AND slot between @slot_init and @slot_end
  AND CASE WHEN @is_priority_entry = true THEN s.priority_entry = ANY(@priority_entry::integer[]) ELSE true END
  AND CASE WHEN @is_professional = true THEN p.reference_key = ANY(@reference_key::varchar[]) ELSE true END
  AND CASE WHEN @is_open = true THEN s.status_entry = 'open' ELSE true END
  AND CASE WHEN @is_especialidade = true THEN p.especialidade = ANY(@especialidade::varchar[]) ELSE true END
  AND CASE WHEN @is_idclinica = true THEN s.id_professional in (
    SELECT a.id_professional FROM attribute a WHERE attribute = 'idclinica' and value = ANY(@idclinica::varchar[])
  ) ELSE true END
ORDER BY s.slot;

-- name: ListSlotsSummary :many
SELECT
  date_trunc(@timing::varchar, s.slot)::timestamp as timeref,
	s.status_entry,
	count(*) as counting
FROM slot s
LEFT JOIN professional p on s.id_professional = p.id_professional
WHERE 1=1
  AND CASE WHEN @deleted = true THEN true ELSE is_deleted = false END
  AND CASE WHEN @is_hour = true THEN cast(concat(extract(hour from slot), ':', extract(minute from slot)) as time) between cast(@init_hour::varchar as time) and cast(@end_hour::varchar as time) ELSE true END
  AND slot between @slot_init and @slot_end
  AND CASE WHEN @is_professional = true THEN p.reference_key = ANY(@reference_key::varchar[]) ELSE true END
  AND CASE WHEN @is_priority_entry = true THEN s.priority_entry = ANY(@priority_entry::integer[]) ELSE true END
  AND CASE WHEN @is_open = true THEN s.status_entry = 'open' ELSE true END
  AND CASE WHEN @is_especialidade = true THEN p.especialidade = ANY(@especialidade::varchar[]) ELSE true END
  AND CASE WHEN @is_idclinica = true THEN s.id_professional in (
    SELECT a.id_professional FROM attribute a WHERE attribute = 'idclinica' and value = ANY(@idclinica::varchar[])
  ) ELSE true END
GROUP BY 1,2
ORDER BY 1,2;

-- name: ListSlotsByIdAvailability :many
SELECT
  id_slot
FROM slot
WHERE 1=1
  AND is_deleted = FALSE
  AND id_availability = @id_availability;

-- name: DeleteSlotById :exec
UPDATE slot
SET is_deleted = TRUE,
  updated_at = CURRENT_TIMESTAMP,
  deleted_at = CURRENT_TIMESTAMP
WHERE id_slot = @id_slot;

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
  id_external,
  is_deleted
FROM slot
WHERE 1=1
  AND id_slot = @id_slot
  AND CASE WHEN @deleted = true THEN true ELSE is_deleted = FALSE END;

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