package main

import (
	"context"
	"fmt"
	"github.com/labstack/echo-contrib/jaegertracing"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/swaggo/echo-swagger"
	"github.com/willie-lin/cloud-terminal/pkg/api"
	"github.com/willie-lin/cloud-terminal/pkg/config"
	"github.com/willie-lin/cloud-terminal/pkg/database"
	"github.com/willie-lin/cloud-terminal/pkg/handler"
	"go.uber.org/zap"
	"golang.org/x/crypto/acme/autocert"
	"net/http"
)

const versionFile = "/app/VERSION"

func createLogger(encoding string) (*zap.Logger, error) {
	if encoding == "json" {
		return zap.NewProduction()
	}
	return zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "console",
		EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}.Build()
}

func main() {
	log, _ := zap.NewDevelopment()
	e := echo.New()
	e.AutoTLSManager.Cache = autocert.DirCache("/var/www/.cache")
	c := jaegertracing.New(e, nil)
	defer c.Close()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(middleware.Gzip())

	// 连接 数据库
	//client, err := database.Client()
	client, err := database.Client()
	//client, err := config.NewClient()
	if err != nil {
		log.Fatal("opening ent client", zap.Error(err))
		panic(err)
	}
	fmt.Println("dddd")

	fmt.Println(client)
	fmt.Println("eeee")
	ctx := context.Background()

	//autoMigration := database.AutoMigration
	autoMigration := config.AutoMigration
	autoMigration(client, ctx)

	//debugMode := database.DebugMode
	debugMode := config.DebugMode

	debugMode(err, client, ctx)

	e.GET("/users", handler.GetAllUser(client))
	e.POST("/user", handler.CreateUser(client))
	e.POST("/api/login", api.Login(client))
	e.PUT("/user", handler.UpdateUser(client))
	e.PUT("/user/uid", handler.UpdateUserById(client))
	e.PUT("/test", handler.TestBindJson(client))

	e.DELETE("/user", handler.DeleteUser(client))
	e.DELETE("/user/uid", handler.DeleteUserById(client))
	e.GET("/user/uname", handler.FindUserByUsername(client))
	e.GET("/user/uid", handler.FindUserById(client))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "hello world!!!")
	})

	defer client.Close()

	e.Logger.Fatal(e.Start(":2021"))

}
