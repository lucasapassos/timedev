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
    priority_entry
) VALUES (
    @id_professional,
    @init_datetime,
    @end_datetime,
    @init_hour,
    @end_hour,
    @type_availability,
    @weekday_name,
    @interval,
    @resting,
    @priority_entry
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
WHERE id_availability = @id_availability;

-- name: DeleteAvailabilityById :one
UPDATE availability
SET is_deleted = TRUE
WHERE id_availability = @id_availability
RETURNING *;

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
  AND id_professional = @id_professional
  AND CASE WHEN @deleted = true THEN true ELSE is_deleted = false END;