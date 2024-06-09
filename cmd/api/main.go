package main

import (
	"fmt"
	"log/slog"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/kamalsingh200238/ubu_management/internal/database"
	"github.com/kamalsingh200238/ubu_management/internal/routes"
	"github.com/labstack/echo/v4"
	// "github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	// e.Use(middleware.Logger())

	e.Static("/assets", "./public")

	err := database.StartDatabase()
	if err != nil {
		slog.Error("error in starting database", err)
	}
	defer database.CloseDatabase()

	routes.RegisterRoutes(e)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", os.Getenv("PORT"))))
}
