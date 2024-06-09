package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

var (
	database   = os.Getenv("DB_DATABASE")
	password   = os.Getenv("DB_PASSWORD")
	username   = os.Getenv("DB_USERNAME")
	port       = os.Getenv("DB_PORT")
	host       = os.Getenv("DB_HOST")
	DBInstance *pgx.Conn
	DBQueries  *Queries
)

func StartDatabase() error {
	var err error
	dbConnString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", username, password, host, port, database)
	DBInstance, err = pgx.Connect(context.Background(), dbConnString)
	if err != nil {
		return fmt.Errorf("error connecting to the database: %v", err)
	}

	DBQueries = New(DBInstance)

	// Check if the connection is successful
	err = DBInstance.Ping(context.Background())
	if err != nil {
		// Close the connection if there's an error
		DBInstance.Close(context.Background())
		return fmt.Errorf("error pinging database: %v", err)
	}

	// Connection successful
	fmt.Println("Successfully connected to DB!")
	return nil
}

func CloseDatabase() {
	if DBInstance != nil {
		DBInstance.Close(context.Background())
		fmt.Println("Database connection closed.")
	}
}
