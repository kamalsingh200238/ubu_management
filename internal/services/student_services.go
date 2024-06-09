package services

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/kamalsingh200238/ubu_management/internal/database"
)

func CheckStudentExistByEmail(email string) (bool, database.Student, error) {
	student, err := database.DBQueries.GetStudentByEmail(context.Background(), email)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, student, nil
		}
		return false, student, err
	}
	return true, student, err
}

func CheckIfStudentIsPresidentByID(id int) (bool, database.Society, error) {
	a, err := database.DBQueries.GetSocietyByPresidentId(context.Background(), pgtype.Int4{Int32: int32(id), Valid: true})
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, a, nil
		}
		return false, a, err
	}
	return true, a, nil
}
