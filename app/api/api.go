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
	"github.com/willie-lin/cloud-terminal/app/database/ent/privacy"
	"github.com/willie-lin/cloud-terminal/app/database/ent/role"
	"github.com/willie-lin/cloud-terminal/app/database/ent/user"
	"github.com/willie-lin/cloud-terminal/pkg/utils"
	"net/http"
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

		tx, err := client.Tx(context.Background())
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "create transaction error"})
		}
		defer func() {
			if r := recover(); r != nil || err != nil {
				err := tx.Rollback()
				if err != nil {
					return
				}
			}
		}()

		// 创建新租户
		tenant, err := tx.Tenant.Create().SetName(dto.TenantName).Save(context.Background())
		if err != nil {
			err := tx.Rollback()
			if err != nil {
				return err
			}
			log.Printf("Error creating tenant: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error creating tenant in database"})
		}

		// 获取所有默认角色
		//defaultRoles, err := tx.Role.Query().All(context.Background())
		defaultRoles, err := tx.Role.Query().Where(role.IsDefaultEQ(true)).All(context.Background())
		if err != nil {
			err := tx.Rollback()
			if err != nil {
				return err
			}
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error getting default roles"})
		}
		// 通过事务关联租户和角色 (关键修改在这里)
		for _, r := range defaultRoles {
			err = tx.Tenant.UpdateOne(tenant).AddRoles(r).Exec(context.Background()) // 直接在事务中使用 AddRoles
			if err != nil {
				err := tx.Rollback()
				if err != nil {
					return err
				}
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error assigning role to tenant"})
			}
		}

		// 创建用户并设置为租户的超级管理员
		us, err := tx.User.Create().
			SetEmail(dto.Email).
			SetUsername(username).
			SetPassword(string(hashedPassword)).
			SetTenantID(tenant.ID).
			Save(context.Background())
		if err != nil {
			err := tx.Rollback()
			if err != nil {
				return err
			}
			log.Printf("Error creating user: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "create user failed"})
		}
		//us, err := client.User.Create().
		//	SetEmail(dto.Email).
		//	SetUsername(username).
		//	SetPassword(string(hashedPassword)).
		//	SetTenantID(tenant.ID).
		//	Save(context.Background())
		//if err != nil {
		//	log.Printf("Error creating user: %v", err)
		//	return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error creating user in database"})
		//}
		//检查管理员角色是否存在，如果不存在则创建
		//adminRole, err := client.Role.Query().Where(role.NameEQ("Admin")).Only(context.Background())

		adminRole, err := tx.Role.Query().Where(role.NameEQ("Admin")).Only(context.Background())

		if ent.IsNotFound(err) {
			err := tx.Rollback()
			if err != nil {
				return err
			}
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "admin role is not found"})
		}
		if err = tx.User.UpdateOne(us).AddRoles(adminRole).Exec(context.Background()); err != nil {
			err := tx.Rollback()
			if err != nil {
				return err
			}
			log.Printf("Error assigning super admin role to user: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "add role to user failed"})
		}

		if err := tx.Commit(); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "commit transaction error"})
		}

		//if err = client.User.UpdateOne(us).AddRoles(adminRole).Exec(context.Background()); err != nil {
		//	log.Printf("Error assigning super admin role to user: %v", err)
		//	return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error assigning super admin role"})
		//}
		//return c.JSON(http.StatusCreated, map[string]string{"userID": us.ID.String()})
		return c.JSON(http.StatusCreated, map[string]string{"message": "Tenant and admin created successfully"})
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
		us, err := client.User.Query().Where(user.EmailEQ(dto.Email)).WithTenant().Only(ctx)
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

		// 查询租户信息，通过边查询获取用户关联的租户
		tenant, err := us.QueryTenant().Only(ctx)
		if err != nil {
			log.Printf("Error finding tenant: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error finding tenant"})
		}

		// 获取用户的第一个角色ID
		//role, err := us.QueryRoles().First(context.Background())
		role, err := us.QueryRoles().Only(ctx)
		//role, err := us.QueryRoles().All(ctx)
		if err != nil {
			log.Printf("Error querying roles: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error querying roles"})
		}
		fmt.Println(role.ID)
		fmt.Println(role.Name)
		fmt.Println(role.Description)

		// 生成包含租户信息的accessToken
		accessToken, err := utils.CreateAccessToken(us.ID, tenant.ID, us.Email, us.Username, role.Name)
		if err != nil {
			log.Printf("Error signing token: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error signing token"})
		}
		// 生成包含租户信息的RefreshToken
		refreshToken, err := utils.CreateRefreshToken(us.ID, tenant.ID, us.Email, us.Username, role.Name)
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
					"id":       us.ID,
					"tenantId": tenant.ID,
					"email":    us.Email,
					"username": us.Username,
					"roleName": role.Name,
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
