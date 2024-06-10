package controllers

import (
	"net/http"

	"github.com/kamalsingh200238/ubu_management/internal/templates"
	"github.com/kamalsingh200238/ubu_management/internal/utils"
	"github.com/labstack/echo/v4"
)

func ShowStudentDashboard(c echo.Context) error {
	return utils.Render(c, http.StatusOK, templates.StudentDashboard())
}
