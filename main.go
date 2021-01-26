package main

import (
	"context"
	"fmt"
	"github.com/labstack/echo-contrib/jaegertracing"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/willie-lin/cloud-terminal/pkg/api"
	"github.com/willie-lin/cloud-terminal/pkg/config"
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
	//client, err := database.Client()
	client, err := config.NewClient()
	if err != nil {
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

	e.GET("/users", api.GetAllUser())
	e.POST("/user", api.CreateUser())
	e.PUT("/user", api.UpdateUser())
	e.PUT("/user/uid", api.UpdateUserById())
	e.PUT("/test", api.TestBindJson())

	e.DELETE("/user", api.DeleteUser())
	e.DELETE("/user/uid", api.DeleteUserById())
	e.GET("/user/uname", api.FindUserByUsername())
	e.GET("/user/uid", api.FindUserById())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "hello world!!!")
	})

	defer client.Close()

	e.Logger.Fatal(e.Start(":2021"))

}
