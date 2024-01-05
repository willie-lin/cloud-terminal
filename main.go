package main

import (
	"context"
	"fmt"
	"github.com/bykof/gostradamus"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/gommon/log"
	"github.com/willie-lin/cloud-terminal/app/api"
	"github.com/willie-lin/cloud-terminal/app/logger"
	"go.elastic.co/apm/module/apmechov4"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/swaggo/echo-swagger"
	"github.com/willie-lin/cloud-terminal/app/config"
	"github.com/willie-lin/cloud-terminal/app/handler"
	_ "github.com/willie-lin/cloud-terminal/docs"
	"go.uber.org/zap"
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
	e := echo.New()
	// 使用重定向中间件将http连接重定向到https
	//e.Pre(middleware.HTTPSRedirect())

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
	// CORS middleware
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))

	e.Use(middleware.Gzip())

	// 连接 数据库
	//client, err := database.Client()
	//client, err := database.Client()
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

	defer client.Close()
	ctx := context.Background()

	//autoMigration := database.AutoMigration
	autoMigration := config.AutoMigration
	autoMigration(client, ctx)

	//debugMode := database.DebugMode
	debugMode := config.DebugMode

	debugMode(err, client, ctx)

	//v1 := e.Group("/api/v1")
	//v1.Use()
	e.GET("/", handler.Hello(client))
	e.GET("/ip", handler.RealIP())
	e.POST("api/enable-2fa", handler.Enable2FA(client))
	e.POST("/api/confirm-2FA", handler.Confirm2FA(client))
	e.POST("/api/check-2FA", handler.Check2FA(client))

	e.POST("/api/check-email", api.CheckEmail(client))
	e.POST("/api/login", api.LoginUser(client))
	e.POST("/api/register", api.RegisterUser(client))
	e.POST("/api/reset-password", api.ForgotPassword(client))

	e.POST("/api/uploads", handler.UploadFile())
	e.GET("/api/users", handler.GetAllUsers(client))
	e.POST("/api/edit-userinfo", handler.UpdateUserInfo(client))

	e.POST("/api/user/email", handler.GetUserByEmail(client))

	//e.POST("/api/login", api.Login(client))
	//e.GET("/user/uid", handler.FindUserById(client))
	//e.PUT("/user", handler.UpdateUser(client))
	//e.PUT("/user/uid", handler.UpdateUserById(client))
	//e.PUT("/test", handler.TestBindJson(client))
	//
	//e.DELETE("/user", handler.DeleteUser(client))
	//e.DELETE("/user/uid", handler.DeleteUserById(client))
	//// UserGroup
	//e.GET("/groups", handler.GetAllGroups(client))
	//e.GET("/group/name", handler.FindGroupByName(client))
	//e.POST("/group", handler.CreateGroup(client))
	//e.DELETE("/group/uid", handler.DeleteGroupById(client))
	//e.DELETE("/group/name", handler.DeleteGroup(client))
	//e.DELETE("/user/uid", handler.DeleteUserById(client))
	//e.DELETE("/user/uid", handler.DeleteUserById(client))
	//
	//e.POST("/user2group", handler.AddUserToGroup(client))
	//e.PUT("/user4group", handler.DeleteUserFromGroup(client))
	//
	//e.POST("/group2user", handler.AddGroupToUser(client))
	//e.PUT("/group4user", handler.DeleteUserFromGroup(client))
	//
	//e.GET("/group/group8user", handler.FindGroupWithUser(client))
	//e.GET("/group/user_group_name", handler.FindUserByGroupName(client))
	//e.GET("/group/group_user_username", handler.FindGroupByUsername(client))
	//
	//e.GET("/group/user_with_group", handler.GetAllUsersWithGroups(client))
	//e.GET("/group/group_with_user", handler.GetAllGroupsWithUsers(client))

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	//go func() {
	//	e.Logger.Fatal(e.Start(":80"))
	//}()
	//
	//e.Logger.Fatal(e.StartTLS(":443", "./cert/cert.pem", "./cert/key.pem"))

	e.Logger.Fatal(e.Start(":2023"))

}
