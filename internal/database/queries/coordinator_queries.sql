-- name: GetCoordinatorByEmail :one
SELECT
  *
FROM
  society_coordinators
WHERE
  email = $1
LIMIT
  1;

-- name: AddCoordinator :one
INSERT INTO
  society_coordinators (name, email, password_hash)
VALUES
  ($1, $2, $3)
RETURNING
  *;
