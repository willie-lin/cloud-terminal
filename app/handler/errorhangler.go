package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/willie-lin/cloud-terminal/app/database/ent"
	"net/http"
)

func handleError(c echo.Context, err error) error {
	if ent.IsNotFound(err) {
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}
	return nil
}
