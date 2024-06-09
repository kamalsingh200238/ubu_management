-- name: GetStudentByEmail :one
SELECT
  *
FROM
  students
WHERE
  email = $1
LIMIT
  1;

-- name: GetStudentById :one
SELECT
  *
FROM
  students
WHERE
  id = $1
LIMIT
  1;

-- name: AddStudent :one
INSERT INTO
  students (name, email, password_hash)
VALUES
  ($1, $2, $3)
RETURNING
  *;

-- name: GetSocietyByPresidentId :one
SELECT
  *
FROM
  societies
WHERE
  president_id = $1;
