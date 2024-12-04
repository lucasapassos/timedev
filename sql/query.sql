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
    priority_entry
) VALUES (
    ?, ?, ?, ?, ?, ?
)
RETURNING *;