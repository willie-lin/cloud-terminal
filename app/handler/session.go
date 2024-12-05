package handler

import (
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"net/http"
)

func CheckSession(c echo.Context) error {
	sess, _ := session.Get("session", c)
	if sess.IsNew {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Session expired or not found"})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Session is active"})
}
