package api

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"github.com/labstack/gommon/log"
	"github.com/pquerna/otp/totp"
	"github.com/willie-lin/cloud-terminal/ent"
	"github.com/willie-lin/cloud-terminal/ent/accesspolicy"
	"github.com/willie-lin/cloud-terminal/ent/group"
	"github.com/willie-lin/cloud-terminal/ent/privacy"
	"github.com/willie-lin/cloud-terminal/ent/role"
	"github.com/willie-lin/cloud-terminal/ent/schema"
	"github.com/willie-lin/cloud-terminal/ent/tenant"
	"github.com/willie-lin/cloud-terminal/ent/user"
	"github.com/willie-lin/cloud-terminal/pkg/crypto"
	"github.com/willie-lin/cloud-terminal/pkg/utils"
)

// CheckEmail 检查邮箱是否已经存在
func CheckEmail(client *ent.Client) echo.HandlerFunc {
	return func(c *echo.Context) error {
		type UserDTO struct {
			Email string `json:"email"`
		}

		dto := new(UserDTO)
		if err := c.Bind(&dto); err != nil {
			log.Printf("Error binding user: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
		}

		ctx := context.Background()
		exists, err := client.User.Query().Where(user.EmailEQ(dto.Email)).Exist(ctx)
		if err != nil {
			log.Printf("Error checking email: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error checking email from database"})
		}
		return c.JSON(http.StatusOK, map[string]bool{"exists": exists})
	}
}

// RegisterUser 用户注册（已禁用公开注册）
func RegisterUser(client *ent.Client) echo.HandlerFunc {
	return func(c *echo.Context) error {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Registration is disabled. Please contact your administrator."})
	}
}

// _RegisterUserDeprecated 旧的注册逻辑（已废弃）
func _RegisterUserDeprecated(client *ent.Client) echo.HandlerFunc {
	return func(c *echo.Context) error {

		username := utils.GenerateUsername()
		type UserDTO struct {
			Email      string `json:"email"`
			Password   string `json:"password"`
			TenantName string `json:"tenant_name"`
		}

		dto := new(UserDTO)
		if err := c.Bind(&dto); err != nil {
			log.Printf("Error binding user: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
		}

		// 创建密码的哈希值
		hashedPassword, err := utils.GenerateFromPassword([]byte(dto.Password), utils.DefaultCost)
		if err != nil {
			log.Printf("Error hashing password: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error hashing password"})
		}

		// 查询平台
		ctx := context.Background()
		platformName, err := client.Platform.Query().Only(ctx)
		if err != nil && !ent.IsNotFound(err) {
			return err
		}

		// 创建租户
		tenantName := strings.ToLower(dto.TenantName) + "_tenant"
		_, err = client.Tenant.Query().Where(tenant.NameEQ(tenantName)).Only(ctx)
		if err != nil && !ent.IsNotFound(err) {
			return err
		}
		if ent.IsNotFound(err) {
			desc := "Cloud Terminal Tenant"
			if platformName != nil {
				desc = platformName.Name
			}
			_, err = client.Tenant.Create().
				SetName(tenantName).
				SetDescription(desc).
				Save(ctx)
			if err != nil {
				return err
			}
			log.Printf("Created management tenant for platform")
		} else {
			log.Printf("Management tenant already exists.")
		}

		// 查询角色
		tenantRoleName := "tenant_admin"

		r, err := client.Role.Query().Where(role.NameEQ(tenantRoleName)).Only(ctx)
		if err != nil && !ent.IsNotFound(err) {
			return fmt.Errorf("query role %s failed: %w", tenantRoleName, err)
		}
		if ent.IsNotFound(err) {
			r, err = client.Role.Create().
				SetName(tenantRoleName).
				SetIsDefault(true).
				Save(ctx)
			if err != nil {
				return fmt.Errorf("create role %s failed: %w", tenantRoleName, err)
			}
		} else {
			log.Printf("Role: %s (ID: %v) already exists.", r.Name, r.ID)
		}

		// 5. 创建租户管理员策略
		tenantRole, err := client.Role.Query().Where(role.NameEQ(tenantRoleName)).Only(ctx)
		if err != nil {
			return fmt.Errorf("query tenant admin role failed: %w", err)
		}

		tenantPolicyName := tenantRoleName + "_policy"
		tenantPolicy, err := client.AccessPolicy.Query().Where(accesspolicy.NameEQ(tenantPolicyName)).Only(ctx)
		if err != nil && !ent.IsNotFound(err) {
			return fmt.Errorf("query tenant admin policy failed: %w", err)
		}

		if ent.IsNotFound(err) {
			statements := []schema.PolicyStatement{
				{
					Effect: "Allow",
					Actions: []string{
						utils.ActionUserCreate,
						utils.ActionUserRead,
						utils.ActionUserUpdate,
						utils.ActionUserDelete,
						utils.ActionRoleCreate,
						utils.ActionRoleRead,
						utils.ActionRoleUpdate,
						utils.ActionRoleDelete,
						utils.ActionAuditLogRead,
						utils.ActionAuditLogExport,
					},
					Resources: []string{
						utils.ResourceUserAll,
						utils.ResourceRoleAll,
						utils.ResourcePolicyAll,
						utils.ResourceAuditLogAll,
					},
				},
			}

			tenantPolicy, err = client.AccessPolicy.Create().
				SetName(tenantPolicyName).
				SetStatements(statements).
				Save(ctx)
			if err != nil {
				return fmt.Errorf("创建 tenant admin policy 失败: %w", err)
			}
			log.Printf("Created tenant admin policy: %s (ID: %v)", tenantPolicy.Name, tenantPolicy.ID)
		} else if err == nil {
			log.Printf("tenant admin policy: %s (ID: %v) already exists.", tenantPolicy.Name, tenantPolicy.ID)
		} else {
			return fmt.Errorf("unexpected error when querying tenant admin policy: %w", err)
		}

		// 检查策略是否已经关联到角色
		exists, err := client.Role.Query().
			Where(role.IDEQ(tenantRole.ID)).
			Where(role.HasAccessPoliciesWith(accesspolicy.IDEQ(tenantPolicy.ID))).
			Exist(ctx)
		if err != nil {
			return fmt.Errorf("checking role policy existence: %w", err)
		}

		if !exists {
			_, err = tenantRole.Update().Save(ctx)
			if err != nil {
				return fmt.Errorf("关联 tenant_admin 角色和 tenant_admin 策略失败: %w", err)
			}
			log.Printf("关联 tenant_admin 角色和 tenant_admin 策略")
		} else {
			log.Printf("tenant_admin 角色和 tenant_admin 策略已经关联")
		}

		// 创建账户 Group
		groupName := strings.ToLower(dto.TenantName) + "Group"
		act, err := client.Group.Query().Where(group.NameEQ(groupName)).Only(ctx)
		if err != nil && !ent.IsNotFound(err) {
			return err
		}
		if ent.IsNotFound(err) {
			act, err = client.Group.Create().
				SetName(groupName).
				Save(ctx)
			if err != nil {
				return err
			}
			log.Printf("Created Group for platform")
		} else {
			log.Printf("Group already exists.")
		}

		// 6. 创建租户管理员用户并关联到账户
		us, err := client.User.Create().
			SetUsername(username).
			SetEmail(dto.Email).
			SetPassword(string(hashedPassword)).
			SetIsDefault(true).
			SetGroup(act).
			AddRoles(tenantRole).
			Save(ctx)
		return c.JSON(http.StatusCreated, map[string]string{"userID": us.ID})
	}
}

// LoginUser 用户登陆
func LoginUser(client *ent.Client) echo.HandlerFunc {
	return func(c *echo.Context) error {

		type LoginDTO struct {
			Email    string  `json:"email"`
			Password string  `json:"password"`
			OTP      *string `json:"otp,omitempty"`
		}
		dto := new(LoginDTO)
		if err := c.Bind(&dto); err != nil {
			log.Printf("Error binding user: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
		}


		ctx := privacy.DecisionContext(context.Background(), privacy.Allow)
		loginInput := strings.TrimSpace(dto.Email)
		if loginInput == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Email or username is required"})
		}

		// 同时支持通过邮箱(Email)或用户名(Username)登录
		us, err := client.User.Query().
			Where(
				user.Or(
					user.EmailEQ(loginInput),
					user.UsernameEQ(loginInput),
				),
			).
			Only(ctx)

		if ent.IsNotFound(err) {
			log.Printf("User not found for login input '%s': %v", loginInput, err)
			return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
		}
		if err != nil {
			log.Printf("Error querying user: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error querying user from database"})
		}
		if len(dto.Password) == 0 {
			log.Printf("Error: password is empty")
			return c.JSON(http.StatusForbidden, map[string]string{"error": "Password is empty"})
		}

		err = utils.CompareHashAndPassword([]byte(us.Password), []byte(dto.Password))
		if err != nil {
			log.Printf("Error comparing password: %v", err)
			return c.JSON(http.StatusForbidden, map[string]string{"error": "Invalid password"})
		}

		// 检查用户是否已经绑定了二次验证（2FA）
		if us.TotpSecret != "" {
			if dto.OTP == nil {
				log.Printf("Error: OTP is required")
				return c.JSON(http.StatusForbidden, map[string]string{"error": "OTP is required"})
			}
			// 解密存储的 totp_secret（兼容历史明文存量数据）
			plainSecret, err := crypto.DecryptString(us.TotpSecret)
			if err != nil {
				log.Printf("Error decrypting TOTP secret: %v", err)
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error verifying OTP"})
			}
			valid := totp.Validate(*dto.OTP, plainSecret)
			if !valid {
				log.Printf("Error: Invalid OTP")
				return c.JSON(http.StatusForbidden, map[string]string{"error": "Invalid OTP"})
			}
		}
		// update user LastLoginTime
		_, err = client.User.
			UpdateOne(us).
			SetLastLoginTime(time.Now()).
			Save(ctx)
		if err != nil {
			log.Printf("Error updating last login time: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error updating last login time"})
		}
		// 1. 查询 Group 信息（增加容错保护，未绑定 Group 时不阻断登录）
		var groupID uuid.UUID = uuid.Nil
		g, gErr := us.QueryGroup().Only(ctx)
		if gErr == nil && g != nil {
			if parsedGID, err := uuid.Parse(g.ID); err == nil {
				groupID = parsedGID
			}
		}

		// 2. 查询租户信息（获取系统租户或默认首个租户）
		var t *ent.Tenant
		t, err = client.Tenant.Query().Where(tenant.NameEQ(utils.ManagementTenant)).Only(ctx)
		if ent.IsNotFound(err) || t == nil {
			t, _ = client.Tenant.Query().First(ctx)
		}
		if t == nil {
			log.Printf("Error finding tenant for user: %s", us.ID)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Tenant not found"})
		}

		// 3. 获取用户的角色信息（增加容错，默认赋予 'user' 角色）
		roleName := "user"
		r, rErr := us.QueryRoles().First(ctx)
		if rErr == nil && r != nil {
			roleName = r.Name
		} else {
			log.Printf("User %s has no explicit role, falling back to 'user'", us.Username)
		}

		isTenantAdmin := strings.Contains(strings.ToLower(roleName), "tenant_admin")
		isSuperAdmin := strings.ToLower(roleName) == "super_admin"

		// 生成包含租户信息的 accessToken & refreshToken
		accessToken, err := utils.CreateAccessToken(uuid.MustParse(us.ID), uuid.MustParse(t.ID), groupID, us.Email, us.Username, roleName)
		if err != nil {
			log.Printf("Error signing token: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error signing token"})
		}
		refreshToken, err := utils.CreateRefreshToken(uuid.MustParse(us.ID), uuid.MustParse(t.ID), groupID, us.Email, us.Username, roleName)
		if err != nil {
			log.Printf("Error signing refreshToken: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error signing refreshToken"})
		}

		// 4. 将 token 保存在 Cookie 中（兼容 HTTP/HTTPS 开发与生产环境）
		isSecure := c.Scheme() == "https" || c.Request().TLS != nil
		sameSite := http.SameSiteLaxMode
		if isSecure {
			sameSite = http.SameSiteNoneMode
		}

		// 提取请求 Host 并去除端口号，防止 Cookie 因带端口被浏览器拒绝
		cookieDomain := c.Request().Host
		if h, _, err := net.SplitHostPort(cookieDomain); err == nil {
			cookieDomain = h
		}
		if cookieDomain == "localhost" || cookieDomain == "127.0.0.1" {
			cookieDomain = "" // 置空 Domain 适应本地开发
		}

		accessTokenCookie := &http.Cookie{
			Name:     "AccessToken",
			Value:    accessToken,
			Expires:  time.Now().Add(24 * time.Hour),
			SameSite: sameSite,
			Domain:   cookieDomain,
			HttpOnly: true,
			Secure:   isSecure,
			Path:     "/",
		}
		c.SetCookie(accessTokenCookie)

		refreshTokenCookie := &http.Cookie{
			Name:     "RefreshToken",
			Value:    refreshToken,
			Expires:  time.Now().Add(24 * time.Hour),
			SameSite: sameSite,
			Domain:   cookieDomain,
			HttpOnly: true,
			Secure:   isSecure,
			Path:     "/",
		}
		c.SetCookie(refreshTokenCookie)

		gIDStr := ""
		if groupID != uuid.Nil {
			gIDStr = groupID.String()
		}

		return c.JSON(http.StatusOK,
			map[string]interface{}{
				"accessToken":  accessToken,
				"refreshToken": refreshToken,
				"user": map[string]interface{}{
					"id":            us.ID,
					"tenantId":      t.ID,
					"groupId":       gIDStr,
					"email":         us.Email,
					"username":      us.Username,
					"roleName":      roleName,
					"isTenantAdmin": isTenantAdmin,
					"isSuperAdmin":  isSuperAdmin,
				},
			})
	}
}

// ResetPassword reset password
func ResetPassword(client *ent.Client) echo.HandlerFunc {
	return func(c *echo.Context) error {
		type UserDTO struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		dto := new(UserDTO)
		if err := c.Bind(&dto); err != nil {
			log.Printf("Error binding user: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
		}

		ua, err := client.User.Query().Where(user.EmailEQ(dto.Email)).Only(context.Background())
		if ent.IsNotFound(err) {
			log.Printf("User not found: %v", err)
			return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
		}
		if err != nil {
			log.Printf("Error querying user: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error querying user from database"})
		}
		hashedPassword, err := utils.GenerateFromPassword([]byte(dto.Password), utils.DefaultCost)
		if err != nil {
			log.Printf("Error hashing password: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error hashing password"})
		}

		_, err = client.User.
			UpdateOne(ua).SetPassword(string(hashedPassword)).Save(context.Background())
		if err != nil {
			log.Printf("Error updating password: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error updating password in database"})
		}
		return c.JSON(http.StatusOK, map[string]string{"message": "Password reset successful, please log in again."})
	}
}

func LogoutUser() echo.HandlerFunc {
	return func(c *echo.Context) error {
		// 删除访问令牌Cookie
		accessTokenCookie := &http.Cookie{
			Name:     "AccessToken",
			Value:    "",
			Expires:  time.Now().Add(-1 * time.Hour),
			SameSite: http.SameSiteNoneMode,
			HttpOnly: true,
			Secure:   true,
			Path:     "/",
		}
		c.SetCookie(accessTokenCookie)

		// 删除刷新令牌Cookie
		refreshTokenCookie := &http.Cookie{
			Name:     "RefreshToken",
			Value:    "",
			Expires:  time.Now().Add(-1 * time.Hour),
			SameSite: http.SameSiteNoneMode,
			HttpOnly: true,
			Secure:   true,
			Path:     "/",
		}
		c.SetCookie(refreshTokenCookie)

		// 删除CSRF Cookie
		csrfCookie := &http.Cookie{
			Name:     "_csrf",
			Value:    "",
			Expires:  time.Now().Add(-1 * time.Hour),
			SameSite: http.SameSiteNoneMode,
			HttpOnly: true,
			Secure:   true,
			Path:     "/",
		}
		c.SetCookie(csrfCookie)

		return c.JSON(http.StatusOK, map[string]string{"message": "Logout successful"})
	}
}
