package utils

import (
	"log/slog"

	"github.com/angelofallars/htmx-go"
	"github.com/labstack/echo/v4"
)

func Redirect(c echo.Context, newPath string) error {
	if htmx.IsHTMX(c.Request()) {
		if err := htmx.NewResponse().Redirect(newPath).Write(c.Response().Writer); err != nil {
			slog.Error("error in redirecting", err)
			return err
		}
		return c.String(302, "redirect")
	}
	return c.Redirect(302, newPath)
}
