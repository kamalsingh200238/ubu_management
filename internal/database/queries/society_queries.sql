-- name: AddSociety :one
INSERT INTO
  societies (name, active, president_id)
VALUES
  ($1, $2, $3)
RETURNING
  *;

-- name: UpdateSociety :one
UPDATE societies
SET
  name = $1,
  active = $2,
  president_id = $3
WHERE
  id = $4
RETURNING
  *;

-- name: SetSocietyActiveStatus :one
UPDATE societies
SET
  active = $1
WHERE
  id = $2
RETURNING
  *;

-- name: GetAllSocietiesWithPresidentWithStudentCount :many
SELECT
  a.id AS society_id,
  a.name AS society_name,
  a.president_id,
  a.active AS society_active,
  s.name AS president_name,
  s.email AS president_email,
  s.password_hash AS president_password,
  COALESCE(sa.enrolled_count, 0) AS enrolled_student_count
FROM
  societies AS a
  LEFT JOIN students AS s ON s.id = a.president_id
  LEFT JOIN (
    SELECT
      society_id,
      COUNT(*) AS enrolled_count
    FROM
      student_societies
    GROUP BY
      society_id
  ) AS sa ON a.id = sa.society_id
ORDER BY
  a.id;

-- name: GetSocietyWithPresidentBySocietyId :one
SELECT
  s.id AS president_id,
  s.name AS president_name,
  s.email AS president_email,
  s.password_hash AS president_password_hash,
  so.id AS society_id,
  so.name AS society_name,
  so.president_id AS society_president_id,
  so.active AS society_active
FROM
  students AS s
  JOIN societies AS so ON s.id = so.president_id
WHERE
  so.id = $1
LIMIT
  1;

-- name: GetAllSocietiesStudentIsEnrolledIn :many
SELECT
  s.*
FROM
  societies AS s
  JOIN student_societies AS ss ON ss.society_id = s.id
WHERE
  ss.student_id = $1;

-- name: GetAllSocietiesStudentIsNotEnrolledIn :many
SELECT
  a.*
FROM
  societies AS a
  LEFT JOIN student_societies AS ss ON a.id = ss.society_id
  AND ss.student_id = $1
WHERE
  ss.student_id IS NULL
ORDER BY
  a.id;

-- name: GetAllStudentsEnrolledInSocietyOrderByStudentName :many
SELECT
  s.id AS student_id,
  s.name AS student_name,
  ss.society_id AS society_id
FROM
  students AS s
  JOIN student_societies AS ss ON s.id = ss.student_id
  AND ss.society_id = $1
ORDER BY
  s.name;

-- name: EnrollStudentInSociety :one
INSERT INTO
  student_societies (student_id, society_id)
VALUES
  ($1, $2)
RETURNING
  *;

-- name: LeaveSociety :one
DELETE FROM student_societies
WHERE
  student_id = $1
  AND society_id = $2
RETURNING
  *;
