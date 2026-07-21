package handler

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"github.com/labstack/gommon/log"

	"github.com/willie-lin/cloud-terminal/ent"
	"github.com/willie-lin/cloud-terminal/ent/task"
	"github.com/willie-lin/cloud-terminal/pkg/sts"
	entresource "github.com/willie-lin/cloud-terminal/ent/resource"
	entuser "github.com/willie-lin/cloud-terminal/ent/user"
	"github.com/willie-lin/cloud-terminal/viewer"
)

// ─── DTO ──────────────────────────────────────────────────────

type CreateTaskDTO struct {
	ResourceID     string `json:"resource_id"`
	Reason         string `json:"reason"`
	DurationMinutes int    `json:"duration_minutes"`
}

type ApproveTaskDTO struct {
	ReviewerComment string `json:"reviewer_comment,omitempty"`
}

type RejectTaskDTO struct {
	ReviewerComment string `json:"reviewer_comment,omitempty"`
}

type TaskHandler struct {
	client     *ent.Client
	stsService *sts.Service
}

func NewTaskHandler(client *ent.Client, stsService *sts.Service) *TaskHandler {
	return &TaskHandler{client: client, stsService: stsService}
}

// ─── Create ──────────────────────────────────────────────────

func (h *TaskHandler) Create() echo.HandlerFunc {
	return func(c *echo.Context) error {
		v := viewer.FromContext(c.Request().Context())
		if v == nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		}

		dto := new(CreateTaskDTO)
		if err := c.Bind(dto); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
		}
		if dto.Reason == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "reason is required"})
		}
		if dto.DurationMinutes <= 0 {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "duration_minutes must be positive"})
		}
		if dto.ResourceID == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "resource_id is required"})
		}

		// 检查资源是否存在
		exists, err := h.client.Resource.Query().
			Where(entresource.IDEQ(dto.ResourceID)).
			Exist(c.Request().Context())
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error checking resource"})
		}
		if !exists {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "resource not found"})
		}

		// 检查是否有同一资源 pending 的重复申请
		dup, _ := h.client.Task.Query().
			Where(
				task.HasRequesterWith(entuser.IDEQ(v.UserID.String())),
				task.HasResourceWith(entresource.IDEQ(dto.ResourceID)),
				task.StatusEQ(task.StatusPending),
			).
			Exist(c.Request().Context())
		_ = dup  // ignore errors on duplicate check

		if dup {
			return c.JSON(http.StatusConflict, map[string]string{"error": "You already have a pending request for this resource"})
		}

		t, err := h.client.Task.Create().
			SetReason(dto.Reason).
			SetDurationMinutes(dto.DurationMinutes).
			SetRequesterID(v.UserID.String()).
			SetResourceID(dto.ResourceID).
			Save(c.Request().Context())
		if err != nil {
			log.Printf("Error creating task: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create task"})
		}

		return c.JSON(http.StatusCreated, t)
	}
}

// ─── List ────────────────────────────────────────────────────

func (h *TaskHandler) List() echo.HandlerFunc {
	return func(c *echo.Context) error {
		v := viewer.FromContext(c.Request().Context())
		if v == nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		}

		query := h.client.Task.Query()

		// 普通用户只看自己的申请
		isSuperAdmin := v.RoleName == "super_admin"
		isTenantAdmin := strings.Contains(strings.ToLower(v.RoleName), "tenant_admin")
		_ = isTenantAdmin
		if !isSuperAdmin {
			query = query.Where(task.HasRequesterWith(entuser.IDEQ(v.UserID.String())))
		}

		if status := c.QueryParam("status"); status != "" {
			query = query.Where(task.StatusEQ(task.Status(status)))
		}
		if resourceID := c.QueryParam("resource_id"); resourceID != "" {
			query = query.Where(task.HasResourceWith(entresource.IDEQ(resourceID)))
		}
		if requesterID := c.QueryParam("requester_id"); requesterID != "" && (isSuperAdmin) {
			query = query.Where(task.HasRequesterWith(entuser.IDEQ(requesterID)))
		}

		// 关联加载
		query = query.WithRequester().WithResource()

		tasks, err := query.Order(ent.Desc(task.FieldCreatedAt)).All(c.Request().Context())
		if err != nil {
			log.Printf("Error querying tasks: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error querying tasks"})
		}
		return c.JSON(http.StatusOK, tasks)
	}
}

// ─── Get ─────────────────────────────────────────────────────

func (h *TaskHandler) Get() echo.HandlerFunc {
	return func(c *echo.Context) error {
		id := c.Param("id")
		if _, err := uuid.Parse(id); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid task ID"})
		}

		t, err := h.client.Task.Query().
			Where(task.IDEQ(id)).
			WithRequester().
			WithResource().
			WithReviewer().
			Only(c.Request().Context())
		if ent.IsNotFound(err) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Task not found"})
		}
		if err != nil {
			log.Printf("Error querying task: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error querying task"})
		}
		return c.JSON(http.StatusOK, t)
	}
}

// ─── Approve ─────────────────────────────────────────────────

func (h *TaskHandler) Approve() echo.HandlerFunc {
	return func(c *echo.Context) error {
		v := viewer.FromContext(c.Request().Context())
		if v == nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		}

		isSuperAdmin := v.RoleName == "super_admin"
		isTenantAdmin := strings.Contains(strings.ToLower(v.RoleName), "tenant_admin")
		if !isSuperAdmin && !isTenantAdmin {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "Only admins can approve tasks"})
		}

		id := c.Param("id")
		if _, err := uuid.Parse(id); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid task ID"})
		}

		dto := new(ApproveTaskDTO)
		_ = c.Bind(dto)

		// 获取 task + resource
		t, err := h.client.Task.Query().
			Where(task.IDEQ(id)).
			WithResource().
			Only(c.Request().Context())
		if err != nil {
			if ent.IsNotFound(err) {
				return c.JSON(http.StatusNotFound, map[string]string{"error": "Task not found"})
			}
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error querying task"})
		}
		if t.Status != task.StatusPending {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("Task is already %s", t.Status)})
		}

		// 签发 STS token
		resourceURN := t.Edges.Resource.Urn
		ttl := time.Duration(t.DurationMinutes) * time.Minute

		resp, err := h.stsService.IssueToken(c.Request().Context(), &sts.IssueRequest{
			UserID:      v.UserID.String(),
			ResourceURN: resourceURN,
			TTL:         ttl,
		}, nil, nil)
		if err != nil {
			log.Printf("Error issuing STS token: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to issue token"})
		}

		// 更新 task
		now := time.Now()
		t, err = h.client.Task.UpdateOneID(id).
			SetStatus(task.StatusApproved).
			SetReviewedAt(now).
			SetReviewerID(v.UserID.String()).
			SetNillableReviewerComment(strPtr(dto.ReviewerComment)).
			SetIssuedToken(resp.Token).
			SetExpiresAt(resp.ExpiresAt).
			Save(c.Request().Context())
		if err != nil {
			log.Printf("Error updating task: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update task"})
		}

		return c.JSON(http.StatusOK, t)
	}
}

// ─── Reject ──────────────────────────────────────────────────

func (h *TaskHandler) Reject() echo.HandlerFunc {
	return func(c *echo.Context) error {
		v := viewer.FromContext(c.Request().Context())
		if v == nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		}

		isSuperAdmin := v.RoleName == "super_admin"
		isTenantAdmin := strings.Contains(strings.ToLower(v.RoleName), "tenant_admin")
		if !isSuperAdmin && !isTenantAdmin {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "Only admins can reject tasks"})
		}

		id := c.Param("id")
		if _, err := uuid.Parse(id); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid task ID"})
		}

		dto := new(RejectTaskDTO)
		_ = c.Bind(dto)

		t, err := h.client.Task.Query().Where(task.IDEQ(id)).Only(c.Request().Context())
		if err != nil {
			if ent.IsNotFound(err) {
				return c.JSON(http.StatusNotFound, map[string]string{"error": "Task not found"})
			}
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error querying task"})
		}
		if t.Status != task.StatusPending {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("Task is already %s", t.Status)})
		}

		now := time.Now()
		t, err = h.client.Task.UpdateOneID(id).
			SetStatus(task.StatusRejected).
			SetReviewedAt(now).
			SetReviewerID(v.UserID.String()).
			SetNillableReviewerComment(strPtr(dto.ReviewerComment)).
			Save(c.Request().Context())
		if err != nil {
			log.Printf("Error rejecting task: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to reject task"})
		}

		return c.JSON(http.StatusOK, t)
	}
}

// ─── Delete ──────────────────────────────────────────────────

func (h *TaskHandler) Delete() echo.HandlerFunc {
	return func(c *echo.Context) error {
		v := viewer.FromContext(c.Request().Context())
		if v == nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		}
		isSuperAdmin := v.RoleName == "super_admin"
		if !isSuperAdmin {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "Only super admins can delete tasks"})
		}

		id := c.Param("id")
		if _, err := uuid.Parse(id); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid task ID"})
		}

		err := h.client.Task.DeleteOneID(id).Exec(c.Request().Context())
		if ent.IsNotFound(err) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Task not found"})
		}
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete task"})
		}
		return c.NoContent(http.StatusNoContent)
	}
}

// ─── helpers ─────────────────────────────────────────────────

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
