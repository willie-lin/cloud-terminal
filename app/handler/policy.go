package handler

import (
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/labstack/echo/v4"
)

type PolicyHandler struct {
	enforcer *casbin.Enforcer
}

func NewPolicyHandler(enforcer *casbin.Enforcer) *PolicyHandler {
	return &PolicyHandler{enforcer: enforcer}
}

func (h *PolicyHandler) AddPolicy(c echo.Context) error {
	var policy struct {
		Subject string `json:"subject"`
		Domain  string `json:"domain"`
		Object  string `json:"object"`
		Action  string `json:"action"`
	}

	if err := c.Bind(&policy); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	_, err := h.enforcer.AddPolicy(policy.Subject, policy.Domain, policy.Object, policy.Action)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "Policy added successfully"})
}

func (h *PolicyHandler) RemovePolicy(c echo.Context) error {
	var policy struct {
		Subject string `json:"subject"`
		Domain  string `json:"domain"`
		Object  string `json:"object"`
		Action  string `json:"action"`
	}

	if err := c.Bind(&policy); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	_, err := h.enforcer.RemovePolicy(policy.Subject, policy.Domain, policy.Object, policy.Action)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Policy removed successfully"})
}

func (h *PolicyHandler) GetAllPolicies(c echo.Context) error {
	policies := h.enforcer.GetPolicy()
	return c.JSON(http.StatusOK, policies)
}

func (h *PolicyHandler) AddRoleForUser(c echo.Context) error {
	var roleAssignment struct {
		User   string `json:"user"`
		Role   string `json:"role"`
		Domain string `json:"domain"`
	}

	if err := c.Bind(&roleAssignment); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	_, err := h.enforcer.AddGroupingPolicy(roleAssignment.User, roleAssignment.Role, roleAssignment.Domain)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "Role assigned successfully"})
}

func (h *PolicyHandler) RemoveRoleForUser(c echo.Context) error {
	var roleAssignment struct {
		User   string `json:"user"`
		Role   string `json:"role"`
		Domain string `json:"domain"`
	}

	if err := c.Bind(&roleAssignment); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	_, err := h.enforcer.RemoveGroupingPolicy(roleAssignment.User, roleAssignment.Role, roleAssignment.Domain)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Role removed successfully"})
}

func (h *PolicyHandler) GetAllRoles(c echo.Context) error {
	roles := h.enforcer.GetAllRoles()
	return c.JSON(http.StatusOK, roles)
}
