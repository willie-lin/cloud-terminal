package handler

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v5"
	"net/http"
)

var sessionStore *sessions.CookieStore

func InitSessionStore() {
	sessionStore = sessions.NewCookieStore([]byte("cloud-terminal-secret-key-change-in-production"))
}

func GetSessionStore() *sessions.CookieStore {
	return sessionStore
}

func CheckSession(c *echo.Context) error {
	sess, err := sessionStore.Get(c.Request(), "session")
	if err != nil || sess.IsNew {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Session expired or not found"})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Session is active"})
}
