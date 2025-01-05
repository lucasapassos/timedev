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
  AND id_professional = @id_professional
  AND CASE WHEN @deleted = true THEN true ELSE is_deleted = false END;

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
  AND id_blocker = @id_blocker
  AND CASE WHEN @deleted = true THEN true ELSE is_deleted = false END;

-- name: DeleteBlockerById :one
UPDATE blocker
SET is_deleted = true
WHERE id_blocker = @id_blocker
RETURNING *;
