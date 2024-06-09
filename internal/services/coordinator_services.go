package services

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/kamalsingh200238/ubu_management/internal/database"
)

func CheckCoordinatorExistsByEmail(email string) (bool, database.SocietyCoordinator, error) {
	coordinator, err := database.DBQueries.GetCoordinatorByEmail(context.Background(), email)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, coordinator, nil
		}
		return false, coordinator, err
	}
	return true, coordinator, err
}
