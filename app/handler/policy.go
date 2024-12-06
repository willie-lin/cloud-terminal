package handler

import (
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/labstack/echo/v4"
)

// AddPolicy

func AddPolicy(enforcer *casbin.Enforcer) echo.HandlerFunc {
	return func(c echo.Context) error {
		var policy struct {
			Subject string `json:"subject"`
			Domain  string `json:"domain"`
			Object  string `json:"object"`
			Action  string `json:"action"`
		}

		if err := c.Bind(&policy); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		_, err := enforcer.AddPolicy(policy.Subject, policy.Domain, policy.Object, policy.Action)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusCreated, map[string]string{"message": "Policy added successfully"})
	}

}

// RemovePolicy

func RemovePolicy(enforcer *casbin.Enforcer) echo.HandlerFunc {
	return func(c echo.Context) error {
		var policy struct {
			Subject string `json:"subject"`
			Domain  string `json:"domain"`
			Object  string `json:"object"`
			Action  string `json:"action"`
		}

		if err := c.Bind(&policy); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		_, err := enforcer.RemovePolicy(policy.Subject, policy.Domain, policy.Object, policy.Action)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, map[string]string{"message": "Policy removed successfully"})
	}

}

// GetAllPolicies

func GetAllPolicies(enforcer *casbin.Enforcer) echo.HandlerFunc {
	return func(c echo.Context) error {
		policies := enforcer.GetPolicy()
		return c.JSON(http.StatusOK, policies)
	}

}

// AddRoleForUser

func AddRoleForUser(enforcer *casbin.Enforcer) echo.HandlerFunc {
	return func(c echo.Context) error {
		var roleAssignment struct {
			User   string `json:"user"`
			Role   string `json:"role"`
			Domain string `json:"domain"`
		}

		if err := c.Bind(&roleAssignment); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		_, err := enforcer.AddGroupingPolicy(roleAssignment.User, roleAssignment.Role, roleAssignment.Domain)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusCreated, map[string]string{"message": "Role assigned successfully"})
	}

}

// RemoveRoleForUser

func RemoveRoleForUser(enforcer *casbin.Enforcer) echo.HandlerFunc {
	return func(c echo.Context) error {
		var roleAssignment struct {
			User   string `json:"user"`
			Role   string `json:"role"`
			Domain string `json:"domain"`
		}

		if err := c.Bind(&roleAssignment); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		_, err := enforcer.RemoveGroupingPolicy(roleAssignment.User, roleAssignment.Role, roleAssignment.Domain)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, map[string]string{"message": "Role removed successfully"})
	}

}

// GetForAllRoles

func GetForAllRoles(enforcer *casbin.Enforcer) echo.HandlerFunc {
	return func(c echo.Context) error {
		roles := enforcer.GetAllRoles()
		return c.JSON(http.StatusOK, roles)
	}

}
