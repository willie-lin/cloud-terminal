package api

import (
	"context"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/pquerna/otp/totp"
	"github.com/willie-lin/cloud-terminal/app/database/ent"
	"github.com/willie-lin/cloud-terminal/app/database/ent/accesspolicy"
	"github.com/willie-lin/cloud-terminal/app/database/ent/account"
	"github.com/willie-lin/cloud-terminal/app/database/ent/privacy"
	"github.com/willie-lin/cloud-terminal/app/database/ent/role"
	"github.com/willie-lin/cloud-terminal/app/database/ent/schema"
	"github.com/willie-lin/cloud-terminal/app/database/ent/tenant"
	"github.com/willie-lin/cloud-terminal/app/database/ent/user"
	"github.com/willie-lin/cloud-terminal/pkg/utils"
	"net/http"
	"strings"
	"time"
)

// CheckEmail 检查邮箱是否已经存在
func CheckEmail(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		type UserDTO struct {
			Email string `json:"email"`
		}

		dto := new(UserDTO)
		if err := c.Bind(&dto); err != nil {
			log.Printf("Error binding user: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
		}
		ctx := privacy.DecisionContext(context.Background(), privacy.Allow)
		exists, err := client.User.Query().Where(user.EmailEQ(dto.Email)).Exist(ctx)

		if err != nil {
			log.Printf("Error checking email: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error checking email from database"})
		}
		return c.JSON(http.StatusOK, map[string]bool{"exists": exists})
	}
}

// RegisterUser 用户注册
func RegisterUser(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		// 使用 privacy.DecisionContext 跳过隐私检查
		ctx := privacy.DecisionContext(context.Background(), privacy.Allow)

		username := utils.GenerateUsername()
		type UserDTO struct {
			Email      string `json:"email"`
			Password   string `json:"password"`
			TenantName string `json:"tenant_name"` // 新增租户名称字段
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
		platformName, err := client.Platform.Query().Only(ctx)
		if err != nil && !ent.IsNotFound(err) {
			return err
		}

		// 创建租户
		tenantName := strings.ToLower(dto.TenantName) + "_tenant"
		tt, err := client.Tenant.Query().Where(tenant.NameEQ(tenantName)).Only(ctx)
		if err != nil && !ent.IsNotFound(err) {
			return err
		}
		if ent.IsNotFound(err) {
			tt, err = client.Tenant.Create().
				SetName(tenantName).
				SetPlatform(platformName).
				Save(ctx)
			if err != nil {
				return err
			}
			log.Printf("Created management tenant for %s platform", platformName)
		} else {
			log.Printf("Management tenant already exists for %s platform.", platformName)
		}

		// 查询角色
		tenantRoleName := tenantName + "_admin" // Replace with your desired role name

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
						// ActionTenantCreate Tenant Actions
						// ActionUserCreate User Actions
						utils.ActionUserCreate,
						utils.ActionUserRead,
						utils.ActionUserUpdate,
						utils.ActionUserDelete,

						// ActionRoleCreate ActionUserCreate User Actions
						utils.ActionRoleCreate,
						utils.ActionRoleRead,
						utils.ActionRoleUpdate,
						utils.ActionRoleDelete,

						// ActionProjectCreate Project Actions
						utils.ActionProjectCreate,
						utils.ActionProjectRead,
						utils.ActionProjectUpdate,
						utils.ActionProjectDelete,

						// ActionAuditLogRead Audit Log Actions
						utils.ActionAuditLogRead,
						utils.ActionAuditLogExport,
					}, // 超级管理员拥有所有操作权限
					Resources: []string{
						utils.ResourceUserAll, // 匹配所有租户和账户的用户
						utils.ResourceAccountAll,
						utils.ResourceRoleAll,
						utils.ResourcePolicyAll,
						utils.ResourceProjectAll,
						utils.ResourceAuditLogAll,
					}, // 超级管理员拥有所有资源权限
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
			_, err = tenantRole.Update().AddAccessPolicies(tenantPolicy).Save(ctx)
			if err != nil {
				return fmt.Errorf("关联 tenant_admin 角色和 tenant_admin 策略失败: %w", err)
			}
			log.Printf("关联 tenant_admin 角色和 tenant_admin 策略")
		} else {
			log.Printf("tenant_admin 角色和 tenant_admin 策略已经关联")
		}

		// 创建账户
		accountName := strings.ToLower(dto.TenantName) + "Account"
		act, err := client.Account.Query().Where(account.NameEQ(accountName)).Only(ctx)
		if err != nil && !ent.IsNotFound(err) {
			return err
		}
		if ent.IsNotFound(err) {
			act, err = client.Account.Create().
				SetName(accountName).
				SetTenant(tt).
				AddRoles(tenantRole).
				AddAccessPolicies(tenantPolicy).
				Save(ctx)
			if err != nil {
				return err
			}
			log.Printf("Created Account Service for %s platform", platformName)
		} else {
			log.Printf("MAccount Service already exists for %s platform.", platformName)
		}

		// 6. 创建租户管理员用户并关联到账户
		us, err := client.User.Create().
			SetUsername(username).
			SetEmail(dto.Email).
			SetPassword(string(hashedPassword)).
			SetIsDefault(true).
			SetAccount(act).
			SetRole(tenantRole).
			Save(ctx)
		return c.JSON(http.StatusCreated, map[string]string{"userID": us.ID.String()})
		//return c.JSON(http.StatusCreated, map[string]string{"message": "Tenant and admin created successfully"})
	}
}

//type jwtCustomClaims struct {
//	Email string `json:"email"`
//	jwt.RegisteredClaims
//}

// LoginUser 用户登陆
func LoginUser(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {

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

		// 使用决策上下文进行查询，跳过隐私规则
		ctx := privacy.DecisionContext(context.Background(), privacy.Allow)

		//fmt.Println(dto.OTP)
		us, err := client.User.Query().Where(user.EmailEQ(dto.Email)).Only(ctx)
		if ent.IsNotFound(err) {
			log.Printf("User not found: %v", err)
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

		// 假设 us.Password 是数据库中存储的哈希值
		err = utils.CompareHashAndPassword([]byte(us.Password), []byte(dto.Password))
		if err != nil {
			log.Printf("Error comparing password: %v", err)
			return c.JSON(http.StatusForbidden, map[string]string{"error": "Invalid password"})
		}

		// 检查用户是否已经绑定了二次验证（2FA）
		if us.TotpSecret != "" {
			// 用户已经启用了OTP，所以必须提供OTP
			if dto.OTP == nil {
				log.Printf("Error: OTP is required")
				return c.JSON(http.StatusForbidden, map[string]string{"error": "OTP is required"})
			}
			// 验证用户提供的OTP
			valid := totp.Validate(*dto.OTP, us.TotpSecret)
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
		// 查询账户信息，通过边查询获取用户关联的账户
		sa, err := us.QueryAccount().Only(ctx)
		if ent.IsNotFound(err) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Account not found"})
		} else if err != nil {
			log.Printf("Error finding Account: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
		}
		fmt.Println(sa.ID)

		// 查询租户信息，通过边查询获取用户关联的租户
		tenant, err := sa.QueryTenant().Only(ctx)
		if ent.IsNotFound(err) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Tenant not found"})
		} else if err != nil {
			log.Printf("Error finding tenant: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
		}
		fmt.Println(tenant.ID)

		// 获取用户的第一个角色ID
		//role, err := us.QueryRoles().First(context.Background())
		r, err := us.QueryRole().Only(ctx)
		//r, err := sa.QueryRoles().Only(ctx)
		//role, err := us.QueryRoles().All(ctx)
		if err != nil {
			log.Printf("Error querying roles: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error querying roles"})
		}

		// 关键部分：判断是否为租户管理员
		isTenantAdmin := strings.Contains(strings.ToLower(r.Name), "tenant_admin") // 或更精确的匹配逻辑

		fmt.Println(r.ID)
		fmt.Println(r.Name)
		fmt.Println(r.Description)

		// 生成包含租户信息的accessToken
		accessToken, err := utils.CreateAccessToken(us.ID, tenant.ID, sa.ID, us.Email, us.Username, r.Name)
		if err != nil {
			log.Printf("Error signing token: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error signing token"})
		}
		// 生成包含租户信息的RefreshToken
		refreshToken, err := utils.CreateRefreshToken(us.ID, tenant.ID, sa.ID, us.Email, us.Username, r.Name)
		if err != nil {
			log.Printf("Error signing refreshToken: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error signing refreshToken"})
		}
		// 将token保存在HTTP-only的cookie中，并设置相关的属性
		accessTokenCookie := new(http.Cookie)
		accessTokenCookie.Name = "AccessToken"
		accessTokenCookie.Value = accessToken
		accessTokenCookie.Expires = time.Now().Add(24 * time.Hour)
		accessTokenCookie.SameSite = http.SameSiteNoneMode
		accessTokenCookie.Domain = c.Request().Host
		accessTokenCookie.HttpOnly = true
		accessTokenCookie.Secure = true
		accessTokenCookie.Path = "/"
		c.SetCookie(accessTokenCookie)

		// 创建另一个cookie来保存RefreshToken
		refreshTokenCookie := new(http.Cookie)
		refreshTokenCookie.Name = "RefreshToken"
		refreshTokenCookie.Value = refreshToken
		refreshTokenCookie.Expires = time.Now().Add(24 * time.Hour) // RefreshToken通常有更长的过期时间
		refreshTokenCookie.SameSite = http.SameSiteNoneMode
		refreshTokenCookie.Domain = c.Request().Host
		refreshTokenCookie.HttpOnly = true
		refreshTokenCookie.Secure = true
		refreshTokenCookie.Path = "/"
		c.SetCookie(refreshTokenCookie)

		// 登录成功后，保存用户的登录信息到session
		sess, _ := session.Get("session", c)
		if err != nil {
			log.Printf("Error getting session: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error getting session"})
		}
		sess.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   3600, // 设置session的过期时间
			HttpOnly: true,
		}
		sess.Values["username"] = us.Username // 保存用户名到session
		sess.Values["email"] = us.Email       // 保存用户名到session
		err = sess.Save(c.Request(), c.Response())
		if err != nil {
			log.Printf("Error saving session: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error saving session"})
		}

		// 如果认证成功，设置用户和租户信息到上下文
		//c.Set("user", us)
		//c.Set("tenant", tenant)

		//return c.JSON(http.StatusOK, map[string]string{"message": "Login successful"})
		//return c.JSON(http.StatusOK, map[string]string{"message": "Login successful", "refreshToken": refreshToken})
		// 返回包含用户信息的响应
		return c.JSON(http.StatusOK,
			map[string]interface{}{
				"accessToken":  accessToken,
				"refreshToken": refreshToken,
				"user": map[string]interface{}{
					"id":            us.ID,
					"tenantId":      tenant.ID,
					"accountId":     sa.ID,
					"email":         us.Email,
					"username":      us.Username,
					"roleName":      r.Name,
					"isTenantAdmin": isTenantAdmin, // 添加 isTenantAdmin 字段
				},
			})
	}
}

// ResetPassword  reset password
func ResetPassword(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
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
		// 使用你的方法来创建密码的哈希值
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
	return func(c echo.Context) error {
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

		// 删除Session Cookie
		sessionCookie := &http.Cookie{
			Name:     "session",
			Value:    "",
			Expires:  time.Now().Add(-1 * time.Hour),
			SameSite: http.SameSiteNoneMode,
			HttpOnly: true,
			Secure:   true,
			Path:     "/",
		}
		c.SetCookie(sessionCookie)

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

		// 返回登出成功的响应
		return c.JSON(http.StatusOK, map[string]string{"message": "Logout successful"})
	}
}
