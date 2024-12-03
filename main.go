package main

import (
	"context"
	"fmt"
	"github.com/bykof/gostradamus"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/swaggo/echo-swagger"
	"github.com/willie-lin/cloud-terminal/app/api"
	"github.com/willie-lin/cloud-terminal/app/config"
	"github.com/willie-lin/cloud-terminal/app/database/ent"
	"github.com/willie-lin/cloud-terminal/app/handler"
	"github.com/willie-lin/cloud-terminal/app/logger"
	_ "github.com/willie-lin/cloud-terminal/docs"
	"github.com/willie-lin/cloud-terminal/pkg/utils"
	"go.elastic.co/apm/module/apmechov4"
	"go.uber.org/zap"
	"net/http"
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

	//utils.InitLogger()
	//log, _ := zap.NewDevelopment()
	//log, _ := zap.NewProduction()
	//log := zap.NewProductionEncoderConfig()

	//authEnforcer, err := auth.NewEnforcer("auth_model.conf", "policy.csv")
	//if err != nil {
	//	log.Fatalf("failed to create auth enforcer: %v", err)
	//}
	e := echo.New()

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
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))

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

	//autoMigration := database.AutoMigration
	autoMigration := config.AutoMigration
	autoMigration(client, ctx)

	//debugMode := database.DebugMode
	debugMode := config.DebugMode

	debugMode(err, client, ctx)

	// 设置 Static 中间件
	e.Static("/picture", "picture")

	// CORS middleware
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"https://localhost:3000"},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true,
		AllowMethods:     []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))

	//e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
	//	AllowOrigins: []string{"*"},
	//	AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	//}))

	// 限制IP速率
	rateLimiterConfig := middleware.RateLimiterConfig{
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
	}
	e.Use(middleware.RateLimiterWithConfig(rateLimiterConfig))

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

	e.Use(middleware.Gzip())

	// 定义一个受保护的路由组
	r := e.Group("/admin")
	r.Use(utils.CheckAccessToken)
	// 使用JWT中间件
	r.Use(echojwt.WithConfig(utils.ValidAccessTokenConfig()))

	e.GET("/", handler.Hello(client))
	e.GET("/ip", handler.RealIP())
	e.POST("/api/check-email", api.CheckEmail(client))
	e.POST("/api/check-2FA", handler.Check2FA(client))
	e.POST("/api/login", api.LoginUser(client))
	e.POST("/api/register", api.RegisterUser(client))
	e.POST("/api/reset-password", api.ResetPassword(client))

	// 需要token认证
	r.POST("/enable-2fa", handler.Enable2FA(client))
	r.POST("/confirm-2FA", handler.Confirm2FA(client))
	r.POST("/uploads", handler.UploadFile())
	r.GET("/users", handler.GetAllUsers(client))
	r.POST("/edit-userinfo", handler.UpdateUserInfo(client))
	r.POST("/user/email", handler.GetUserByEmail(client))
	r.POST("/add-user", handler.CreateUser(client))
	r.POST("/update-user", handler.UpdateUser(client))
	r.POST("/delete-user", handler.DeleteUserByUsername(client))

	// role
	r.GET("/roles", handler.GetAllRoles(client))
	r.POST("/add-role", handler.CreateRole(client))
	r.POST("/delete-role", handler.DeleteRoleByName(client))
	r.POST("/check-role-name", handler.CheckRoleName(client))

	// permission
	r.GET("/permissions", handler.GetAllPermissions(client))
	r.POST("/add-permission", handler.CreatePermission(client))
	r.POST("/delete-permission", handler.DeletePermissionByName(client))
	r.POST("/check-permission-name", handler.CheckPermissionName(client))

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	go func() {
		e.Logger.Fatal(e.Start(":80"))
	}()

	e.Logger.Fatal(e.StartTLS(":443", "./cert/cert.pem", "./cert/key.pem"))

	//e.Logger.Fatal(e.StartAutoTLS(":443"))
	//e.Logger.Fatal(e.Start(":2023"))

}
