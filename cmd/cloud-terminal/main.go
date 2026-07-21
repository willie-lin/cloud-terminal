package main

import (
	"context"
	"fmt"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/labstack/gommon/log"
	"github.com/willie-lin/cloud-terminal/api"
	_ "github.com/willie-lin/cloud-terminal/ent/runtime"
	"github.com/willie-lin/cloud-terminal/handler"
	"github.com/willie-lin/cloud-terminal/middlewarers"
	"github.com/willie-lin/cloud-terminal/pkg/config"
	"github.com/willie-lin/cloud-terminal/pkg/database"
	pkglogger "github.com/willie-lin/cloud-terminal/pkg/logger"
	"github.com/willie-lin/cloud-terminal/pkg/audit"
	"github.com/willie-lin/cloud-terminal/pkg/iam"
	"github.com/willie-lin/cloud-terminal/pkg/sts"
	"github.com/willie-lin/cloud-terminal/pkg/utils"
	"go.uber.org/zap"
	"net/http"
	"os"
	"time"
)

func main() {
	// 1. 加载配置
	cfgPath := "config.yaml"
	if p := os.Getenv("CONFIG_PATH"); p != "" {
		cfgPath = p
	}
	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	fmt.Printf("config loaded: app=%s env=%s\n", cfg.App.Name, cfg.App.Environment)

	// 2. 初始化日志
	if err := pkglogger.Init(cfg.Logger); err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}
	pkglogger.Info("Logger initialized", zap.String("level", cfg.Logger.Level))

	// 3. 初始化 JWT
	utils.AccessTokenSecret = cfg.Server.JWTSecret
	if utils.AccessTokenSecret == "" {
		utils.AccessTokenSecret = "cloud-terminal-default-secret-change-me"
	}

	// 4. 连接数据库
	client, err := database.NewClient(&cfg.Database)
	if err != nil {
		pkglogger.Fatal("failed to connect database", zap.Error(err))
	}
	defer func() {
		_ = client.Close()
	}()
	pkglogger.Info("Database connected", zap.String("driver", cfg.Database.Driver))

	// 审计日志
	auditor := audit.NewInMemoryAuditor(10000)
	pkglogger.Info("Auditor initialized", zap.String("type", "in-memory"))

	// STS 服务
		stsService := sts.New([]byte(cfg.Server.JWTSecret))
	pkglogger.Info("STS service initialized")

	// IAM 引擎
	iamProvider := iam.NewEntPolicyProvider(client)
	iamEvaluator := iam.NewEvaluator(iamProvider)
	pkglogger.Info("IAM evaluator initialized")

	// 5. 自动迁移
	ctx := context.Background()
	if cfg.Database.AutoMigrate {
		if err := database.AutoMigration(ctx, client); err != nil {
			pkglogger.Fatal("auto migration failed", zap.Error(err))
		}
		pkglogger.Info("Database migration completed")
	}

	// 6. 初始化超级管理员
	if err := api.InitSuperAdminAndSuperRoles(client); err != nil {
		pkglogger.Fatal("init super admin failed", zap.Error(err))
	}
	pkglogger.Info("Super admin initialized")

	// 7. 启动 Echo
	e := echo.New()

	e.IPExtractor = echo.ExtractIPDirect()
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(middleware.Gzip())

	var allowedOrigins []string
	if os.Getenv("ENV") == "production" {
		allowedOrigins = []string{"https://app.cloudsec.sbs", "https://cloudsec.sbs"}
	} else {
		allowedOrigins = []string{"https://localhost:3000"}
	}

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     allowedOrigins,
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization, "X-CSRF-Token"},
		AllowCredentials: true,
		AllowMethods:     []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		ContentSecurityPolicy: "default-src 'self'; script-src 'self' https://cloudsec.sbs https://trusted-scripts.com; object-src 'none'; frame-ancestors 'none';",
	}))

	e.Static("/picture", "picture")
	e.Use(utils.SetCSRFToken)

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

	e.Use(middleware.RateLimiterWithConfig(middleware.RateLimiterConfig{
		Skipper: middleware.DefaultSkipper,
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(
			middleware.RateLimiterMemoryStoreConfig{Rate: 10, Burst: 30, ExpiresIn: 3 * time.Minute},
		),
		IdentifierExtractor: func(ctx *echo.Context) (string, error) {
			return ctx.RealIP(), nil
		},
		ErrorHandler: func(context *echo.Context, err error) error {
			return context.JSON(http.StatusForbidden, nil)
		},
		DenyHandler: func(context *echo.Context, identifier string, err error) error {
			return context.JSON(http.StatusTooManyRequests, nil)
		},
	}))

	e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:         "1; mode=block",
		ContentTypeNosniff:    "nosniff",
		XFrameOptions:         "SAMEORIGIN",
		HSTSMaxAge:            3600,
		ContentSecurityPolicy: "default-src 'self'",
	}))

	e.Pre(middleware.MethodOverrideWithConfig(middleware.MethodOverrideConfig{
		Getter: middleware.MethodFromForm("_method"),
	}))

	e.Use(middleware.Decompress())

	// 审计中间件：记录所有 POST/PUT/DELETE 请求
	e.Use(audit.Middleware(auditor))

	// Public routes
	e.GET("/api/csrf-token", func(c *echo.Context) error {
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
	e.POST("/webhook/auth", handler.AuthWebhook(client))
	e.POST("/webhook/config", handler.ConfigWebhook(client))
	e.GET("/api/terminal", handler.TerminalWebSocket(client, stsService))

	// ContainerSSH v2 Webhook（集成 STS + IAM）

	// Task handler
	taskHandler := handler.NewTaskHandler(client, stsService)
	csshHandler := handler.NewContainerSSHHandler(client, stsService, iamEvaluator)
	e.POST("/webhook/v2/auth", csshHandler.AuthWebhookV2())
	e.POST("/webhook/v2/config", csshHandler.ConfigWebhookV2())

	// Protected routes
	r := e.Group("/admin")
	r.Use(middlewarers.AuthenticateAndAuthorize)
	r.Use(utils.CheckAccessToken)

	r.POST("/enable-2fa", handler.Enable2FA(client))
	r.POST("/confirm-2FA", handler.Confirm2FA(client))
	r.POST("/disable-2fa", handler.Disable2FA(client))
	r.POST("/reset-2fa", handler.Reset2FA(client))
	r.POST("/uploads", handler.UploadFile())
	r.POST("/edit-userinfo", handler.UpdateUserInfo(client))
	r.POST("/check-role-name", handler.CheckRoleName(client))

	// === User ===
	r.GET("/users", handler.GetAllUsersByTenant(client))
	r.POST("/users", handler.CreateUser(client))
	r.GET("/users/:id", handler.GetUserByID(client))
	r.PUT("/users/:id", handler.UpdateUserByUUID(client))
	r.DELETE("/users/:id", handler.DeleteUserByID(client))

	// === Role ===
	r.GET("/roles", handler.GetAllRolesByAccountByTenant(client))
	r.POST("/roles", handler.CreateRole(client))
	r.GET("/roles/:id", handler.GetRole(client))
	r.PUT("/roles/:id", handler.UpdateRole(client))
	r.DELETE("/roles/:id", handler.DeleteRoleByID(client))

	// === AccessPolicy ===
	r.GET("/access-policies", handler.GetAllAccessPolicyByAccountByTenant(client))
	r.GET("/access-policies/:id", handler.GetAccessPolicy(client))
	r.POST("/access-policies", handler.CreateAccessPolicy(client))
	r.PUT("/access-policies/:id", handler.UpdateAccessPolicy(client))
	r.DELETE("/access-policies/:id", handler.DeleteAccessPolicy(client))

	// === Tenant ===
	r.GET("/tenants", handler.ListTenants(client))
	r.POST("/tenants", handler.CreateTenant(client))
	r.GET("/tenants/:id", handler.GetTenantByID(client))
	r.PUT("/tenants/:id", handler.UpdateTenant(client))
	r.DELETE("/tenants/:id", handler.DeleteTenantByID(client))
	r.POST("/tenants/:id/admin-user", handler.CreateTenantAdmin(client))

	// === Resource ===
	r.GET("/resources", handler.ListResources(client))
	r.POST("/resources", handler.CreateResource(client))
	r.GET("/resources/:id", handler.GetResource(client))
	r.PUT("/resources/:id", handler.UpdateResource(client))
	r.DELETE("/resources/:id", handler.DeleteResource(client))

	// === Group ===
	r.GET("/groups", handler.ListGroups(client))
	r.POST("/groups", handler.CreateGroup(client))
	r.GET("/groups/:id", handler.GetGroup(client))
	r.PUT("/groups/:id", handler.UpdateGroup(client))
	r.DELETE("/groups/:id", handler.DeleteGroup(client))

	// === Environment ===
	r.GET("/environments", handler.ListEnvironments(client))
	r.POST("/environments", handler.CreateEnvironment(client))
	r.GET("/environments/:id", handler.GetEnvironment(client))
	r.PUT("/environments/:id", handler.UpdateEnvironment(client))
	r.DELETE("/environments/:id", handler.DeleteEnvironment(client))

	// === Session ===
	r.GET("/sessions", handler.ListSessions(client))
	r.GET("/sessions/:id", handler.GetSession(client))
	r.PUT("/sessions/:id", handler.UpdateSession(client))
	r.DELETE("/sessions/:id", handler.DeleteSession(client))

	// === AuditLog ===
	r.GET("/audit-logs", handler.ListAuditLogs(client))
	r.GET("/audit-logs/:id", handler.GetAuditLog(client))

	// === Task ===
	r.POST("/tasks", taskHandler.Create())
	r.GET("/tasks", taskHandler.List())
	r.GET("/tasks/:id", taskHandler.Get())
	r.PUT("/tasks/:id/approve", taskHandler.Approve())
	r.PUT("/tasks/:id/reject", taskHandler.Reject())
	r.DELETE("/tasks/:id", taskHandler.Delete())

	// 8. 启动服务
	go func() {
		pkglogger.Info("Starting HTTP server on :80")
		if err := e.Start(":80"); err != nil {
			pkglogger.Fatal("http server failed", zap.Error(err))
		}
	}()

	port := fmt.Sprintf(":%d", cfg.Server.Port)
	if cfg.Server.Port == 0 {
		port = ":443"
	}
	pkglogger.Info("Starting HTTPS server", zap.String("port", port))
	if err := e.Start(port); err != nil {
		pkglogger.Fatal("https server failed", zap.Error(err))
	}
}
