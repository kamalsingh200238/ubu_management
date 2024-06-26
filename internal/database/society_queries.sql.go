// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: society_queries.sql

package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const addSociety = `-- name: AddSociety :one
INSERT INTO
  societies (name, active, president_id)
VALUES
  ($1, $2, $3)
RETURNING
  id, name, president_id, active
`

type AddSocietyParams struct {
	Name        string
	Active      pgtype.Bool
	PresidentID pgtype.Int4
}

func (q *Queries) AddSociety(ctx context.Context, arg AddSocietyParams) (Society, error) {
	row := q.db.QueryRow(ctx, addSociety, arg.Name, arg.Active, arg.PresidentID)
	var i Society
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.PresidentID,
		&i.Active,
	)
	return i, err
}

const enrollStudentInSociety = `-- name: EnrollStudentInSociety :one
INSERT INTO
  student_societies (student_id, society_id)
VALUES
  ($1, $2)
RETURNING
  student_id, society_id
`

type EnrollStudentInSocietyParams struct {
	StudentID int32
	SocietyID int32
}

func (q *Queries) EnrollStudentInSociety(ctx context.Context, arg EnrollStudentInSocietyParams) (StudentSociety, error) {
	row := q.db.QueryRow(ctx, enrollStudentInSociety, arg.StudentID, arg.SocietyID)
	var i StudentSociety
	err := row.Scan(&i.StudentID, &i.SocietyID)
	return i, err
}

const getAllSocietiesStudentIsEnrolledIn = `-- name: GetAllSocietiesStudentIsEnrolledIn :many
SELECT
  s.id, s.name, s.president_id, s.active
FROM
  societies AS s
  JOIN student_societies AS ss ON ss.society_id = s.id
WHERE
  ss.student_id = $1
`

func (q *Queries) GetAllSocietiesStudentIsEnrolledIn(ctx context.Context, studentID int32) ([]Society, error) {
	rows, err := q.db.Query(ctx, getAllSocietiesStudentIsEnrolledIn, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Society
	for rows.Next() {
		var i Society
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.PresidentID,
			&i.Active,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAllSocietiesStudentIsNotEnrolledIn = `-- name: GetAllSocietiesStudentIsNotEnrolledIn :many
SELECT
  a.id, a.name, a.president_id, a.active
FROM
  societies AS a
  LEFT JOIN student_societies AS ss ON a.id = ss.society_id
  AND ss.student_id = $1
WHERE
  ss.student_id IS NULL
ORDER BY
  a.id
`

func (q *Queries) GetAllSocietiesStudentIsNotEnrolledIn(ctx context.Context, studentID int32) ([]Society, error) {
	rows, err := q.db.Query(ctx, getAllSocietiesStudentIsNotEnrolledIn, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Society
	for rows.Next() {
		var i Society
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.PresidentID,
			&i.Active,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAllSocietiesWithPresidentWithStudentCount = `-- name: GetAllSocietiesWithPresidentWithStudentCount :many
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
  a.id
`

type GetAllSocietiesWithPresidentWithStudentCountRow struct {
	SocietyID            int32
	SocietyName          string
	PresidentID          pgtype.Int4
	SocietyActive        pgtype.Bool
	PresidentName        pgtype.Text
	PresidentEmail       pgtype.Text
	PresidentPassword    pgtype.Text
	EnrolledStudentCount int64
}

func (q *Queries) GetAllSocietiesWithPresidentWithStudentCount(ctx context.Context) ([]GetAllSocietiesWithPresidentWithStudentCountRow, error) {
	rows, err := q.db.Query(ctx, getAllSocietiesWithPresidentWithStudentCount)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAllSocietiesWithPresidentWithStudentCountRow
	for rows.Next() {
		var i GetAllSocietiesWithPresidentWithStudentCountRow
		if err := rows.Scan(
			&i.SocietyID,
			&i.SocietyName,
			&i.PresidentID,
			&i.SocietyActive,
			&i.PresidentName,
			&i.PresidentEmail,
			&i.PresidentPassword,
			&i.EnrolledStudentCount,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAllStudentsEnrolledInSocietyOrderByStudentName = `-- name: GetAllStudentsEnrolledInSocietyOrderByStudentName :many
SELECT
  s.id AS student_id,
  s.name AS student_name,
  ss.society_id AS society_id
FROM
  students AS s
  JOIN student_societies AS ss ON s.id = ss.student_id
  AND ss.society_id = $1
ORDER BY
  s.name
`

type GetAllStudentsEnrolledInSocietyOrderByStudentNameRow struct {
	StudentID   int32
	StudentName string
	SocietyID   int32
}

func (q *Queries) GetAllStudentsEnrolledInSocietyOrderByStudentName(ctx context.Context, societyID int32) ([]GetAllStudentsEnrolledInSocietyOrderByStudentNameRow, error) {
	rows, err := q.db.Query(ctx, getAllStudentsEnrolledInSocietyOrderByStudentName, societyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAllStudentsEnrolledInSocietyOrderByStudentNameRow
	for rows.Next() {
		var i GetAllStudentsEnrolledInSocietyOrderByStudentNameRow
		if err := rows.Scan(&i.StudentID, &i.StudentName, &i.SocietyID); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getSocietyWithPresidentBySocietyId = `-- name: GetSocietyWithPresidentBySocietyId :one
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
  1
`

type GetSocietyWithPresidentBySocietyIdRow struct {
	PresidentID           int32
	PresidentName         string
	PresidentEmail        string
	PresidentPasswordHash string
	SocietyID             int32
	SocietyName           string
	SocietyPresidentID    pgtype.Int4
	SocietyActive         pgtype.Bool
}

func (q *Queries) GetSocietyWithPresidentBySocietyId(ctx context.Context, id int32) (GetSocietyWithPresidentBySocietyIdRow, error) {
	row := q.db.QueryRow(ctx, getSocietyWithPresidentBySocietyId, id)
	var i GetSocietyWithPresidentBySocietyIdRow
	err := row.Scan(
		&i.PresidentID,
		&i.PresidentName,
		&i.PresidentEmail,
		&i.PresidentPasswordHash,
		&i.SocietyID,
		&i.SocietyName,
		&i.SocietyPresidentID,
		&i.SocietyActive,
	)
	return i, err
}

const leaveSociety = `-- name: LeaveSociety :one
DELETE FROM student_societies
WHERE
  student_id = $1
  AND society_id = $2
RETURNING
  student_id, society_id
`

type LeaveSocietyParams struct {
	StudentID int32
	SocietyID int32
}

func (q *Queries) LeaveSociety(ctx context.Context, arg LeaveSocietyParams) (StudentSociety, error) {
	row := q.db.QueryRow(ctx, leaveSociety, arg.StudentID, arg.SocietyID)
	var i StudentSociety
	err := row.Scan(&i.StudentID, &i.SocietyID)
	return i, err
}

const setSocietyActiveStatus = `-- name: SetSocietyActiveStatus :one
UPDATE societies
SET
  active = $1
WHERE
  id = $2
RETURNING
  id, name, president_id, active
`

type SetSocietyActiveStatusParams struct {
	Active pgtype.Bool
	ID     int32
}

func (q *Queries) SetSocietyActiveStatus(ctx context.Context, arg SetSocietyActiveStatusParams) (Society, error) {
	row := q.db.QueryRow(ctx, setSocietyActiveStatus, arg.Active, arg.ID)
	var i Society
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.PresidentID,
		&i.Active,
	)
	return i, err
}

const updateSociety = `-- name: UpdateSociety :one
UPDATE societies
SET
  name = $1,
  active = $2,
  president_id = $3
WHERE
  id = $4
RETURNING
  id, name, president_id, active
`

type UpdateSocietyParams struct {
	Name        string
	Active      pgtype.Bool
	PresidentID pgtype.Int4
	ID          int32
}

func (q *Queries) UpdateSociety(ctx context.Context, arg UpdateSocietyParams) (Society, error) {
	row := q.db.QueryRow(ctx, updateSociety,
		arg.Name,
		arg.Active,
		arg.PresidentID,
		arg.ID,
	)
	var i Society
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.PresidentID,
		&i.Active,
	)
	return i, err
}
