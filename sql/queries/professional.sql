-- name: InsertProfessional :one
INSERT INTO professional (
  reference_key,
  nome,
  especialidade
) VALUES (
  @reference_key, @nome, @especialidade
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


-- name: GetProfessionalInfo :one
SELECT
  id_professional,
  reference_key,
  nome,
  especialidade
FROM professional
WHERE reference_key = @reference_key;

-- name: ListAttributesByProfessionalId :many
SELECT
  id_attribute,
  attribute,
  value
FROM attribute
WHERE id_professional = @id_professional;

-- name: CheckProfessionalExists :one
SELECT 1
FROM professional
WHERE id_professional = @id_professional;