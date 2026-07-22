package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"
	"github.com/willie-lin/cloud-terminal/ent"
	"github.com/willie-lin/cloud-terminal/ent/group"
	"github.com/willie-lin/cloud-terminal/ent/privacy"
	"github.com/willie-lin/cloud-terminal/ent/role"
	"github.com/willie-lin/cloud-terminal/ent/tenant"
	"github.com/willie-lin/cloud-terminal/ent/user"
	"github.com/willie-lin/cloud-terminal/pkg/utils"
	"github.com/willie-lin/cloud-terminal/viewer"
)

// Define the Status enum
type Status int

const (
	Active Status = iota
	Inactive
	Blocked
)

// CreateUser 创建一个新用户（由管理员创建，自动关联到管理员所属租户）
func CreateUser(client *ent.Client) echo.HandlerFunc {
	return func(c *echo.Context) (err error) {
		v := viewer.FromContext(c.Request().Context())
		if v == nil {
			log.Printf("No viewer found in context or not authorized")
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		}

		isSuperAdmin := v.RoleName == "super_admin"
		isTenantAdmin := strings.Contains(strings.ToLower(v.RoleName), "tenant_admin")
		if !isSuperAdmin && !isTenantAdmin {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "Only admins can create users"})
		}

		type UserDTO struct {
			Email    string `json:"email"`
			Password string `json:"password"`
			Username string `json:"username"`
			RoleName string `json:"role_name"`
		}

		dto := new(UserDTO)
		if err := c.Bind(&dto); err != nil {
			log.Printf("Error binding user: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
		}
		if dto.Email == "" || dto.Password == "" || strings.TrimSpace(dto.Username) == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Email, password, and username are required"})
		}
		if len(strings.TrimSpace(dto.Username)) < 6 {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Username must be at least 6 characters long"})
		}

		hashedPassword, err := utils.GenerateFromPassword([]byte(dto.Password), utils.DefaultCost)
		if err != nil {
			log.Printf("Error hashing password: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error hashing password"})
		}

		ctx := privacy.DecisionContext(c.Request().Context(), privacy.Allow)

		userBuilder := client.User.Create().
			SetEmail(strings.TrimSpace(dto.Email)).
			SetUsername(strings.TrimSpace(dto.Username)).
			SetPassword(string(hashedPassword)).
			SetOnline(true).
			SetStatus(true)

		// 自动绑定到创建者（租户管理员）所属的 Group 空间，实现租户空间隔离
		creator, cErr := client.User.Query().Where(user.IDEQ(v.UserID.String())).WithGroup().Only(ctx)
		if cErr == nil && creator != nil && creator.Edges.Group != nil {
			userBuilder.SetGroup(creator.Edges.Group)
		}

		// 绑定指定角色（默认为 'user'，可选 'tenant_admin'）
		targetRoleName := strings.TrimSpace(dto.RoleName)
		if targetRoleName == "" {
			targetRoleName = "user"
		}

		userRole, rErr := client.Role.Query().Where(role.NameEQ(targetRoleName)).Only(ctx)
		if rErr != nil || userRole == nil {
			// 兜底退回
			userRole, _ = client.Role.Query().Where(role.NameEQ("user")).Only(ctx)
		}
		if userRole != nil {
			userBuilder.AddRoleIDs(userRole.ID)
		}

		us, err := userBuilder.Save(ctx)
		if err != nil {
			log.Printf("Error creating user: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to create user: %v", err)})
		}
		return c.JSON(http.StatusCreated, map[string]string{"userID": us.ID, "username": us.Username})
	}
}

// CreateTenantAdmin 超管为指定租户创建管理员账号
func CreateTenantAdmin(client *ent.Client) echo.HandlerFunc {
	return func(c *echo.Context) error {
		v := viewer.FromContext(c.Request().Context())
		if v == nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		}
		if v.RoleName != "super_admin" {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "Only super admins can create tenant admin"})
		}

		tenantID := c.Param("id")
		if _, err := uuid.Parse(tenantID); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid tenant ID"})
		}

		type AdminDTO struct {
			Email    string `json:"email"`
			Password string `json:"password"`
			Username string `json:"username"`
		}
		dto := new(AdminDTO)
		if err := c.Bind(dto); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
		}
		if dto.Email == "" || dto.Password == "" || strings.TrimSpace(dto.Username) == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Email, password, and username are required"})
		}
		if len(strings.TrimSpace(dto.Username)) < 6 {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Username must be at least 6 characters long"})
		}
		dto.Username = strings.TrimSpace(dto.Username)

		// 使用 Privacy 提权上下文，确保建用户与关联角色/组不受拦截
		ctx := privacy.DecisionContext(c.Request().Context(), privacy.Allow)

		// 1. 查询目标租户信息
		tenantObj, tErr := client.Tenant.Query().Where(tenant.IDEQ(tenantID)).Only(ctx)
		if tErr != nil {
			log.Printf("Error querying tenant by ID (%s): %v", tenantID, tErr)
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Tenant not found"})
		}

		// 2. 查询或创建该租户专属 Group
		groupName := strings.TrimSpace(tenantObj.Name) + "_Group"
		g, gErr := client.Group.Query().Where(group.NameEQ(groupName)).Only(ctx)
		if gErr != nil {
			g, gErr = client.Group.Create().SetName(groupName).Save(ctx)
			if gErr != nil {
				log.Printf("Error creating group for tenant (%s): %v", tenantObj.Name, gErr)
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to create tenant group: %v", gErr)})
			}
		}

		// 3. 查找 tenant_admin 角色
		tenantAdminRole, err := client.Role.Query().Where(role.NameEQ("tenant_admin")).Only(ctx)
		if err != nil {
			if ent.IsNotFound(err) {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "tenant_admin role not found, run init first"})
			}
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error querying tenant_admin role"})
		}

		hashedPassword, err := utils.GenerateFromPassword([]byte(dto.Password), utils.DefaultCost)
		if err != nil {
			log.Printf("Error hashing password: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error hashing password"})
		}

		us, err := client.User.Create().
			SetEmail(strings.TrimSpace(dto.Email)).
			SetUsername(dto.Username).
			SetPassword(string(hashedPassword)).
			AddRoleIDs(tenantAdminRole.ID).
			SetGroup(g).
			SetOnline(true).
			SetStatus(true).
			Save(ctx)

		if err != nil {
			log.Printf("Error creating tenant admin: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to create tenant admin: %v", err)})
		}

		return c.JSON(http.StatusCreated, map[string]interface{}{
			"userID":   us.ID,
			"username": us.Username,
			"email":    us.Email,
			"role":     "tenant_admin",
			"group":    g.Name,
		})
	}
}

// GetAllUsers 获取所有用户
func GetAllUsers(client *ent.Client) echo.HandlerFunc {
	return func(c *echo.Context) error {
		users, err := client.User.Query().All(c.Request().Context())
		if err != nil {
			log.Printf("Error querying users: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error querying users from database"})
		}
		return c.JSON(http.StatusOK, users)
	}
}

func GetAllUsersByTenant(client *ent.Client) echo.HandlerFunc {
	return func(c *echo.Context) error {
		v := viewer.FromContext(c.Request().Context())
		if v == nil {
			log.Printf("No viewer found in context")
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "No viewer found in context"})
		}
		tenantID := v.TenantID
		userID := v.UserID
		roleName := v.RoleName
		log.Printf("Viewer info: UserID=%s, TenantID=%s, RoleName=%s", userID, tenantID, roleName)

		ctx := privacy.DecisionContext(c.Request().Context(), privacy.Allow)
		var users []*ent.User
		var err error

		isSuperAdmin := roleName == "super_admin"
		isTenantAdmin := strings.Contains(strings.ToLower(roleName), "tenant_admin")

		if isSuperAdmin {
			// 超级管理员能查看平台全局所有用户
			users, err = client.User.Query().All(ctx)
		} else if isTenantAdmin {
			// 租户管理员：强制做租户组隔离，只能查看本租户 Group 组内的用户！
			currUser, uErr := client.User.Query().Where(user.IDEQ(userID.String())).WithGroup().Only(ctx)
			if uErr == nil && currUser != nil && currUser.Edges.Group != nil {
				groupID := currUser.Edges.Group.ID
				users, err = client.User.Query().
					Where(user.HasGroupWith(group.IDEQ(groupID))).
					All(ctx)
			} else {
				// 兜底退回：仅能查看自己
				users, err = client.User.Query().Where(user.IDEQ(userID.String())).All(ctx)
			}
		} else {
			// 普通用户：仅能查看自己
			users, err = client.User.Query().
				Where(user.IDEQ(userID.String())).
				All(ctx)
		}

		if err != nil {
			log.Printf("Error querying users: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error querying users from database"})
		}

		return c.JSON(http.StatusOK, users)
	}
}

// GetUserByUsername 根据用户名查找
func GetUserByUsername(client *ent.Client) echo.HandlerFunc {
	return func(c *echo.Context) error {
		// DTO
		type UserDTO struct {
			Username string `json:"username"`
		}
		dto := new(UserDTO)
		if err := c.Bind(&dto); err != nil {
			log.Printf("Error binding user: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
		}

		us, err := client.User.Query().Where(user.UsernameEQ(dto.Username)).Only(c.Request().Context())
		if ent.IsNotFound(err) {
			log.Printf("User not found: %v", err)
			return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
		}
		if err != nil {
			log.Printf("Error querying user: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error querying user from database"})
		}
		return c.JSON(http.StatusOK, us.Username)
	}
}

// GetUserByEmail 根据邮箱查找用户
func GetUserByEmail(client *ent.Client) echo.HandlerFunc {
	return func(c *echo.Context) error {

		// EmailDTO
		type EmailDTO struct {
			Email string `json:"email"`
		}
		// user response
		type UserResponse struct {
			Avatar      string `json:"avatar"`
			Nickname    string `json:"nickname"`
			Username    string `json:"username"`
			Email       string `json:"email"`
			PhoneNumber string `json:"phone_number"`
			Bio         string `json:"bio"`
		}

		dto := new(EmailDTO)
		if err := c.Bind(&dto); err != nil {
			log.Printf("Error binding user: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
		}
		// 检查请求的邮箱是否与登录用户的邮箱匹配
		ue, err := client.User.Query().Where(user.EmailEQ(dto.Email)).Only(c.Request().Context())
		if ent.IsNotFound(err) {
			log.Printf("User not found: %v", err)
			return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
		}
		if err != nil {
			log.Printf("Error querying user: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error querying user from database"})
		}
		// Map the user entity to the response struct
		response := &UserResponse{
			Avatar:      ue.Avatar,
			Nickname:    ue.Nickname,
			Username:    ue.Username,
			Email:       ue.Email,
			PhoneNumber: ue.PhoneNumber,
			Bio:         ue.Bio,
		}
		return c.JSON(http.StatusOK, response)
	}
}

// UpdateUser 更新一个用户
func UpdateUser(client *ent.Client) echo.HandlerFunc {
	return func(c *echo.Context) error {
		type UpdateUserDTO struct {
			Username string `json:"username"`
			Nickname string `json:"nickname"`
			Phone    string `json:"phone"`
			Bio      string `json:"bio"`
			Online   bool   `json:"online"`
			Status   bool   `json:"status"`
		}

		dto := new(UpdateUserDTO)
		if err := c.Bind(&dto); err != nil {
			log.Printf("Error binding user: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
		}

		// 从数据库中获取用户
		ua, err := client.User.Query().Where(user.UsernameEQ(dto.Username)).Only(c.Request().Context())
		if ent.IsNotFound(err) {
			log.Printf("User not found: %v", err)
			return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
		}
		if err != nil {
			log.Printf("Error querying user: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error querying user from database"})
		}

		// 更新数据库
		_, err = client.User.UpdateOne(ua).
			SetNickname(dto.Nickname).
			SetPhoneNumber(dto.Phone).
			SetBio(dto.Bio).
			SetOnline(dto.Online).
			SetStatus(dto.Status).
			SetLastLoginTime(time.Now()).
			Save(c.Request().Context())
		if err != nil {
			log.Printf("Error updating user info: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error updating user info in database"})
		}
		//return c.JSON(http.StatusOK, map[string]string{"img": ua.Avatar})
		return c.JSON(http.StatusOK, map[string]string{"message": "User update successful"})
	}
}

// GetUserByUUID 根据ID查找用户

// DeleteUserByUsername 删除一个用户
func DeleteUserByUsername(client *ent.Client) echo.HandlerFunc {
	return func(c *echo.Context) error {
		// UserDTO
		type UserDTO struct {
			Username string `json:"username"`
		}

		dto := new(UserDTO)
		if err := c.Bind(&dto); err != nil {
			log.Printf("Error binding user: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
		}

		ue, err := client.User.Query().Where(user.UsernameEQ(dto.Username)).Only(c.Request().Context())
		if ent.IsNotFound(err) {
			log.Printf("User not found: %v", err)
			return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
		}
		if err != nil {
			log.Printf("Error querying user: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error querying user from database"})
		}

		err = client.User.DeleteOne(ue).Exec(c.Request().Context())
		if ent.IsNotFound(err) {
			log.Printf("User not found: %v", err)
			return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
		}
		if err != nil {
			log.Printf("Error deleting user: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error deleting user"})
		}
		return c.NoContent(http.StatusNoContent)
	}
}

// UploadFile UploadAvatar upload Avatar
func UploadFile() echo.HandlerFunc {
	return func(c *echo.Context) error {
		// 获取上传的文件
		file, err := c.FormFile("file")
		if err != nil {
			return err
		}

		// 检查文件类型
		if !strings.HasPrefix(file.Header.Get("Content-Type"), "image/") {
			return errors.New("只允许上传图片文件")
		}

		// 检查文件大小
		if file.Size > 2*1024*1024 { // 限制为2MB
			return errors.New("文件太大，超过了2MB的限制")
		}

		// 打开文件
		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		// 为文件重新命名
		ext := filepath.Ext(file.Filename)
		newFileName := fmt.Sprintf("%d%s", time.Now().Unix(), ext)

		// 创建目标文件
		dst, err := os.Create(filepath.Join("picture", newFileName))
		if err != nil {
			return err
		}
		defer dst.Close()

		// 将源文件复制到目标文件
		if _, err = io.Copy(dst, src); err != nil {
			return err
		}

		// 服务器地址
		serverAddress := "https://127.0.0.1/"
		// 返回文件的路径
		return c.String(http.StatusOK, serverAddress+filepath.Join("picture", newFileName))
	}
}

// UpdateUserInfo Update user info
func UpdateUserInfo(client *ent.Client) echo.HandlerFunc {
	return func(c *echo.Context) (err error) {
		type UpdateUserDTO struct {
			Email    string `json:"email"`
			Nickname string `json:"nickname"`
			Avatar   string `json:"avatar"`
			Phone    string `json:"phone"`
			Bio      string `json:"bio"`
		}

		dto := new(UpdateUserDTO)
		if err := c.Bind(&dto); err != nil {
			log.Printf("Error binding user: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
		}

		// 从数据库中获取用户
		ua, err := client.User.Query().Where(user.EmailEQ(dto.Email)).Only(c.Request().Context())
		if ent.IsNotFound(err) {
			log.Printf("User not found: %v", err)
			return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
		}
		if err != nil {
			log.Printf("Error querying user: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error querying user from database"})
		}

		// 设置更新字段
		update := client.User.UpdateOne(ua)
		if dto.Nickname != "" {
			update.SetNickname(dto.Nickname)
		}
		if dto.Avatar != "" {
			update.SetAvatar(dto.Avatar)
		}
		if dto.Phone != "" {
			update.SetPhoneNumber(dto.Phone)
		}
		if dto.Bio != "" {
			update.SetBio(dto.Bio)
		}

		// 更新用户信息
		//_, err = client.User.
		//	UpdateOne(ua).
		//	SetNickname(dto.Nickname).
		//	SetAvatar(dto.Avatar).
		//	SetPhone(dto.Phone).
		//	SetBio(dto.Bio).
		//	Save(c.Request().Context())
		_, err = update.Save(c.Request().Context())
		if err != nil {
			log.Printf("Error updating user info: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error updating user info in database"})
		}
		//return c.JSON(http.StatusOK, map[string]string{"img": ua.Avatar})
		return c.JSON(http.StatusOK, map[string]string{"message": "User update successful"})
	}
}

func UpdateUserByUUID(client *ent.Client) echo.HandlerFunc {
	return func(c *echo.Context) error {
		v := viewer.FromContext(c.Request().Context())
		if v == nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		}

		id := c.Param("id")
		if _, err := uuid.Parse(id); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
		}

		type UserPatchDTO struct {
			Nickname    *string                `json:"nickname"`
			Bio         *string                `json:"bio"`
			PhoneNumber *string                `json:"phone_number"`
			Avatar      *string                `json:"avatar"`
			Online      *bool                  `json:"online"`
			Status      *bool                  `json:"status"`
			SSHPublicKey *string               `json:"ssh_public_key"`
			Attributes  *map[string]interface{} `json:"attributes"`
		}

		dto := new(UserPatchDTO)
		if err := c.Bind(dto); err != nil {
			log.Printf("Error binding user update: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
		}

		updater := client.User.UpdateOneID(id)
		if dto.Nickname != nil {
			updater.SetNickname(*dto.Nickname)
		}
		if dto.Bio != nil {
			updater.SetBio(*dto.Bio)
		}
		if dto.PhoneNumber != nil {
			updater.SetPhoneNumber(*dto.PhoneNumber)
		}
		if dto.Avatar != nil {
			updater.SetAvatar(*dto.Avatar)
		}
		if dto.Online != nil {
			updater.SetOnline(*dto.Online)
		}
		if dto.Status != nil {
			updater.SetStatus(*dto.Status)
		}
		if dto.SSHPublicKey != nil {
			updater.SetSSHPublicKey(*dto.SSHPublicKey)
		}
		if dto.Attributes != nil {
			updater.SetAttributes(*dto.Attributes)
		}

		user, err := updater.Save(c.Request().Context())
		if ent.IsNotFound(err) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
		}
		if err != nil {
			log.Printf("Error updating user: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update user"})
		}
		return c.JSON(http.StatusOK, user)
	}
}

// ==================== RESTful ID-based handlers ====================

// GetUserByID gets a single user by ID
func GetUserByID(client *ent.Client) echo.HandlerFunc {
	return func(c *echo.Context) error {
		id := c.Param("id")
		if _, err := uuid.Parse(id); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
		}
		u, err := client.User.Query().
			Where(user.IDEQ(id)).
			Only(c.Request().Context())
		if ent.IsNotFound(err) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
		}
		if err != nil {
			log.Printf("Error querying user: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error querying user"})
		}
		return c.JSON(http.StatusOK, u)
	}
}

// DeleteUserByID deletes a user by ID
func DeleteUserByID(client *ent.Client) echo.HandlerFunc {
	return func(c *echo.Context) error {
		v := viewer.FromContext(c.Request().Context())
		if v == nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		}
		isSuperAdmin := v.RoleName == "super_admin"
		isTenantAdmin := strings.Contains(strings.ToLower(v.RoleName), "tenant_admin")
		if !isSuperAdmin && !isTenantAdmin {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "Only admins can delete users"})
		}

		id := c.Param("id")
		if _, err := uuid.Parse(id); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
		}

		err := client.User.DeleteOneID(id).Exec(c.Request().Context())
		if ent.IsNotFound(err) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
		}
		if err != nil {
			log.Printf("Error deleting user: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete user"})
		}
		return c.NoContent(http.StatusNoContent)
	}
}
