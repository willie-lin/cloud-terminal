package main

import (
	"context"
	"fmt"
	"github.com/bykof/gostradamus"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/jaegertracing"
	"github.com/labstack/echo-contrib/session"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/swaggo/echo-swagger"
	"github.com/willie-lin/cloud-terminal/app/api"
	"github.com/willie-lin/cloud-terminal/app/config"
	"github.com/willie-lin/cloud-terminal/app/database/ent"
	_ "github.com/willie-lin/cloud-terminal/app/database/ent/runtime"
	"github.com/willie-lin/cloud-terminal/app/handler"
	"github.com/willie-lin/cloud-terminal/app/logger"
	"github.com/willie-lin/cloud-terminal/app/middlewarers"
	_ "github.com/willie-lin/cloud-terminal/docs"
	"github.com/willie-lin/cloud-terminal/pkg/utils"
	"go.elastic.co/apm/module/apmechov4"
	"go.uber.org/zap"
	"net/http"
	"os"
	"time"
)

const versionFile = "/app/VERSION"

// @title Swagger Example API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host petstore.swagger.io
// @BasePath /v2
func main() {

	e := echo.New()

	// Enable tracing middleware
	c := jaegertracing.New(e, nil)
	defer c.Close()

	// 使用重定向中间件将http连接重定向到https
	e.Pre(middleware.HTTPSRedirect())

	// 设置主机策略
	// e.AutoTLSManager.HostPolicy = autocert.HostWhitelist("<DOMAIN>")

	// 缓存证书以避免达到速率限制
	//e.AutoTLSManager.Cache = autocert.DirCache("/var/www/.cache")

	// 设置日志
	zapLogger, _ := zap.NewProduction()
	e.Use(logger.ZapLogger(zapLogger))

	e.IPExtractor = echo.ExtractIPDirect()
	e.IPExtractor = echo.ExtractIPFromXFFHeader()

	e.Use(apmechov4.Middleware())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(middleware.Gzip())

	var allowedOrigins []string
	if os.Getenv("ENV") == "production" {
		allowedOrigins = []string{"https://app.cloudsec.sbs", "https://cloudsec.sbs"}
	} else {
		allowedOrigins = []string{"https://localhost:3000"}
	}

	// CORS middleware
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     allowedOrigins,
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization, "X-CSRF-Token"},
		AllowCredentials: true,
		AllowMethods:     []string{echo.GET, echo.PUT, echo.POST, echo.DELETE, http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
		//MaxAge:           300,
	}))

	// 设置CSP头
	e.Use(middleware.SecureWithConfig(middleware.SecureConfig{ContentSecurityPolicy: "default-src 'self'; script-src 'self' https://cloudsec.sbs https://trusted-scripts.com; object-src 'none'; frame-ancestors 'none';"}))
	// 设置 Static 中间件
	e.Static("/picture", "picture")

	// 设置session中间件
	e.Use(session.Middleware(sessions.NewCookieStore(securecookie.GenerateRandomKey(64))))

	// 设置 CSRF 令牌
	e.Use(utils.SetCSRFToken)
	//fmt.Println("4444444444444444444444")
	//
	//e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
	//	return func(c echo.Context) error {
	//		fmt.Printf("请求方法: %s, 路径: %s\n", c.Request().Method, c.Request().URL.Path)
	//		fmt.Printf("所有请求头: %+v\n", c.Request().Header)
	//		fmt.Printf("所有cookie: %+v\n", c.Request().Cookies())
	//		csrfToken := c.Request().Header.Get("X-CSRF-Token")
	//		fmt.Printf("Header中的CSRF令牌: %s\n", csrfToken)
	//		cookieToken, err := c.Cookie("_csrf")
	//		if err == nil {
	//			fmt.Printf("Cookie中的CSRF令牌: %s\n", cookieToken.Value)
	//		} else {
	//			fmt.Println("Cookie中没有CSRF令牌")
	//		}
	//		return next(c)
	//	}
	//})
	//fmt.Println("5555555555555555555555555555")

	// 使用 CSRF 中间件
	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup:    "header:X-CSRF-Token,cookie:_csrf,form:_csrf,query:_csrf",
		CookieName:     "_csrf",
		CookiePath:     "/",
		CookieDomain:   "localhost",
		CookieSecure:   true,
		CookieHTTPOnly: true,
		CookieSameSite: http.SameSiteNoneMode,
		CookieMaxAge:   3600,
	}))

	// 限制IP速率
	e.Use(middleware.RateLimiterWithConfig(middleware.RateLimiterConfig{
		Skipper: middleware.DefaultSkipper,
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(
			middleware.RateLimiterMemoryStoreConfig{Rate: 10, Burst: 30, ExpiresIn: 3 * time.Minute},
		),
		IdentifierExtractor: func(ctx echo.Context) (string, error) {
			id := ctx.RealIP()
			return id, nil
		},
		ErrorHandler: func(context echo.Context, err error) error {
			return context.JSON(http.StatusForbidden, nil)
		},
		DenyHandler: func(context echo.Context, identifier string, err error) error {
			return context.JSON(http.StatusTooManyRequests, nil)
		},
	}))

	//Secure 安全
	e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:         "1; mode=block",
		ContentTypeNosniff:    "nosniff",
		XFrameOptions:         "SAMEORIGIN",
		HSTSMaxAge:            3600,
		ContentSecurityPolicy: "default-src 'self'",
		//Secure: "max-age=31536000; includeSubDomains",
	}))

	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Skipper:      middleware.DefaultSkipper,
		ErrorMessage: "custom timeout error message returns to client",
		OnTimeoutRouteErrorHandler: func(err error, c echo.Context) {
			log.Print(c.Path())
		},
		Timeout: 30 * time.Second,
	}))

	e.Pre(middleware.MethodOverrideWithConfig(middleware.MethodOverrideConfig{
		Getter: middleware.MethodFromForm("_method"),
	}))

	e.Use(middleware.Decompress())

	client, err := config.NewClient()
	fmt.Println(client)
	if err != nil {
		log.Fatal("opening ent client", zap.Error(err))
		return
	}
	dateTime := gostradamus.Now()
	fmt.Println(dateTime)
	fmt.Println("dddd")

	fmt.Println(client)
	fmt.Println("eeee")

	defer func(client *ent.Client) {
		err := client.Close()
		if err != nil {

		}
	}(client)
	ctx := context.Background()

	// 运行自动迁移工具
	config.AutoMigration(client, ctx)
	fmt.Println("config auto migration")
	config.DebugMode(err, client, ctx)

	//
	//// 初始化全局权限
	//if err := api.InitializeGlobalPermissions(client); err != nil {
	//	log.Fatalf("failed to initialize global permissions: %v", err)
	//}
	//
	//// 初始化默认角色
	//if err = api.InitializeTenantRolesAndPermissions(client); err != nil {
	//	log.Fatalf("Error initializing tenant roles and permissions: %v", err)
	//}
	if err := api.InitSuperAdminAndSuperRoles(client); err != nil {
		log.Fatal("init superadmin and superadmin roles failed", zap.Error(err))
	}

	// 迁移租户
	//database.AssignDefaultTenant(client)

	// 打开暂时会报错
	//e.Use(middleware.AddTrailingSlashWithConfig(middleware.TrailingSlashConfig{
	//	RedirectCode: http.StatusMovedPermanently,
	//}))

	//初始化Casbin enforcer
	//enforcer, err := casbin.NewEnforcer("./app/casbin/auth_model.conf", "./app/casbin/policy.csv")
	//if err != nil {
	//	log.Fatalf("创建casbin enforcer失败: %v", err)
	//}

	// 初始化处理器
	// Viewer 中间件应在所有请求处理之前应用

	// 在所有路由之前应用中间件
	e.GET("/api/csrf-token", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	e.GET("/", handler.Hello(client))
	e.GET("/ip", handler.RealIP())
	e.POST("/api/check-email", api.CheckEmail(client))
	e.POST("/api/check-2FA", handler.Check2FA(client))
	e.POST("/api/check-tenant-name", handler.CheckTenantName(client))

	e.POST("/api/login", api.LoginUser(client))
	e.POST("/api/logout", api.LogoutUser())
	e.POST("/api/register", api.RegisterUser(client))
	e.POST("/api/reset-password", api.ResetPassword(client))

	// Viewer 中间件应在所有请求处理之前应用
	//e.Use(utils.WithViewer)

	// 定义一个受保护的路由组
	r := e.Group("/admin")
	r.Use(middlewarers.AuthenticateAndAuthorize)
	//r.Use(utils.WithViewer)
	//r.Use(utils.CheckAccessToken)

	//e.Use(utils.WithViewer)
	r.Use(echojwt.WithConfig(utils.ValidAccessTokenConfig()))

	// 定义会话检查端点
	//e.GET("/api/check-session", handler.CheckSession)

	// 需要token认证
	r.POST("/enable-2fa", handler.Enable2FA(client))
	r.POST("/confirm-2FA", handler.Confirm2FA(client))
	r.POST("/uploads", handler.UploadFile())
	//r.GET("/users", handler.GetAllUsers(client))
	r.GET("/users", handler.GetAllUsersByTenant(client))
	r.POST("/user/email", handler.GetUserByEmail(client))
	r.POST("/edit-userinfo", handler.UpdateUserInfo(client))
	//r.POST("/user/email", handler.GetUserByEmail(client), middlewarers.Authorize(enforcer))

	r.POST("/add-user", handler.CreateUser(client))
	r.POST("/update-user", handler.UpdateUser(client))
	r.POST("/delete-user", handler.DeleteUserByUsername(client))

	// role
	//r.GET("/roles", handler.GetAllRoles(client))
	r.GET("/roles", handler.GetAllRolesByAccountByTenant(client))
	//r.POST("/add-role", handler.CreateRole(client), middlewarers.Authorize(enforcer))
	r.POST("/add-role", handler.CreateRole(client))
	//r.POST("/delete-role", handler.DeleteRoleByName(client), middlewarers.Authorize(enforcer))
	r.POST("/delete-role", handler.DeleteRoleByName(client))
	r.POST("/check-role-name", handler.CheckRoleName(client))

	// AccessPolicy
	r.GET("/access-policies", handler.GetAllAccessPolicyByAccountByTenant(client))
	//r.GET("/permissions", handler.GetAllPermissionsByTenant(client))
	////r.POST("/add-permission", handler.CreatePermission(client), middlewarers.Authorize(enforcer))
	//r.POST("/add-permission", handler.CreatePermission(client))
	////r.POST("/delete-permission", handler.DeletePermissionByName(client), middlewarers.Authorize(enforcer))
	//r.POST("/delete-permission", handler.DeletePermissionByName(client))
	//r.POST("/check-permission-name", handler.CheckPermissionName(client))

	//
	//e.POST("/tenants", handler.CreateTenant(client), middlewarers.Authorize(enforcer))
	e.POST("/tenants", handler.CreateTenant(client))
	//e.GET("/tenants/:id", handler.GetTenantByName(client), middlewarers.Authorize(enforcer))
	e.GET("/tenants/:id", handler.GetTenantByName(client))

	//e.POST("/policies", handler.AddPolicy(enforcer), middlewarers.Authorize(enforcer))
	//e.DELETE("/policies", handler.RemovePolicy(enforcer), middlewarers.Authorize(enforcer))
	//e.GET("/policies", handler.GetAllPolicies(enforcer), middlewarers.Authorize(enforcer))

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	//测试使用本地的证书
	go func() {
		e.Logger.Fatal(e.Start(":80"))
	}()

	e.Logger.Fatal(e.StartTLS(":443", "./cert/cert.pem", "./cert/key.pem"))

	//生产编译 docker image
	//e.Logger.Fatal(e.Start(":8080"))

}
