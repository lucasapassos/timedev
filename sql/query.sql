-- -- name: GetAuthor :one
-- SELECT * FROM authors
-- WHERE id = ? LIMIT 1;

-- -- name: ListAuthors :many
-- SELECT * FROM authors
-- ORDER BY name;

-- -- name: CreateAuthor :one
-- INSERT INTO authors (
--   name, bio
-- ) VALUES (
--   ?, ?
-- )
-- RETURNING *;

-- -- name: UpdateAuthor :exec
-- UPDATE authors
-- set name = ?,
-- bio = ?
-- WHERE id = ?;

-- -- name: DeleteAuthor :exec
-- DELETE FROM authors
-- WHERE id = ?;

-- name: InsertProfessional :one
INSERT INTO professional (
  nome,
  especialidade
) VALUES (
  ?, ?
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
    priority_entry
) VALUES (
    ?, ?, ?, ?, ?, ?, ?, ?, ?
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
    priority_entry
FROM availability
WHERE id_availability = ? LIMIT 1;

-- name: GetExistingSlot :one
SELECT 1
FROM slot s
WHERE 1=1
	AND id_professional = ?
	AND datetime(?) between datetime(slot) and datetime(slot, concat(s."interval" - 1, ' minute'))
    AND priority_entry = ?;

-- name: InsertSlot :exec
INSERT INTO slot (
    id_professional,
    id_availability,
    slot,
    weekday_name,
    interval,
    priority_entry,
    status_entry
) VALUES (
    ?, ?, ?, ?, ?, ?, 'Aberto'
)
RETURNING *;

-- name: ListSlots :many
SELECT
  id_slot,
  id_professional,
  id_availability,
  slot,
  weekday_name,
  interval,
  priority_entry,
  status_entry
FROM slot
WHERE 1=1
  AND datetime(slot) between datetime(@slot_init) and datetime(@slot_end)
  AND CASE WHEN @is_professional == true THEN id_professional == @id_professional ELSE 1 END
