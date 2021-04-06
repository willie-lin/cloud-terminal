package main

import (
	"context"
	"fmt"
	"github.com/bykof/gostradamus"
	"go.elastic.co/apm/module/apmechov4"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/swaggo/echo-swagger"
	_ "github.com/willie-lin/cloud-terminal/docs"
	"github.com/willie-lin/cloud-terminal/pkg/api"
	"github.com/willie-lin/cloud-terminal/pkg/config"
	"github.com/willie-lin/cloud-terminal/pkg/handler"
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

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host petstore.swagger.io
// @BasePath /v2
func main() {

	//utils.InitLogger()
	log, _ := zap.NewDevelopment()
	//log, _ := zap.NewProduction()
	//log := zap.NewProductionEncoderConfig()
	e := echo.New()
	//e.AutoTLSManager.Cache = autocert.DirCache("/var/www/.cache")
	//c := jaegertracing.New(e, nil)
	//defer c.Close()
	e.Use(apmechov4.Middleware())
	//e.Use(logger.ZapLogger(log))

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())

	//e.Use(middleware.Gzip())

	// 连接 数据库
	//client, err := database.Client()
	//client, err := database.Client()
	client, err := config.NewClient()
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
	e.GET("/users", handler.GetAllUser(client))
	e.GET("/user/uname", handler.FindUserByUsername(client))
	e.GET("/user/uid", handler.FindUserById(client))
	e.POST("/user", handler.CreateUser(client))
	e.POST("/api/login", api.Login(client))
	e.PUT("/user", handler.UpdateUser(client))
	e.PUT("/user/uid", handler.UpdateUserById(client))
	e.PUT("/test", handler.TestBindJson(client))

	e.DELETE("/user", handler.DeleteUser(client))
	e.DELETE("/user/uid", handler.DeleteUserById(client))
	// UserGroup
	e.GET("/groups", handler.GetAllGroups(client))
	e.GET("/group/name", handler.FindGroupByName(client))
	e.POST("/group", handler.CreateGroup(client))
	e.DELETE("/group/uid", handler.DeleteGroupById(client))
	e.DELETE("/group/name", handler.DeleteGroup(client))
	e.DELETE("/user/uid", handler.DeleteUserById(client))
	e.DELETE("/user/uid", handler.DeleteUserById(client))

	e.POST("/user2group", handler.AddUserToGroup(client))
	e.PUT("/user4group", handler.DeleteUserFromGroup(client))

	e.POST("/group2user", handler.AddGroupToUser(client))
	e.PUT("/group4user", handler.DeleteUserFromGroup(client))

	e.GET("/group/group8user", handler.FindGroupWithUser(client))
	e.GET("/group/user_group_name", handler.FindUserByGroupName(client))
	e.GET("/group/group_user_username", handler.FindGroupByUsername(client))

	e.GET("/group/user_with_group", handler.GetAllUsersWithGroups(client))
	e.GET("/group/group_with_user", handler.GetAllGroupsWithUsers(client))

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.Logger.Fatal(e.Start(":2021"))

}
