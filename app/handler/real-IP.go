package handler

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func RealIP() echo.HandlerFunc {
	return func(c echo.Context) error {
		ip := c.RealIP()
		return c.String(http.StatusOK, "Your IP: "+ip)
	}
}
