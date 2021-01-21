package main

import (
	"context"
	"fmt"
	"github.com/labstack/echo-contrib/jaegertracing"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/willie-lin/cloud-terminal/pkg/api"
	"github.com/willie-lin/cloud-terminal/pkg/database"
	"golang.org/x/crypto/acme/autocert"
	"net/http"
)

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
	client, err := database.Client()
	if err != nil {
		panic(err)
	}
	fmt.Println("dddd")

	fmt.Println(client)
	fmt.Println("eeee")
	ctx := context.Background()

	autoMigration := database.AutoMigration
	autoMigration(client, ctx)

	debugMode := database.DebugMode

	debugMode(err, client, ctx)

	e.GET("/users", api.GetAllUser())
	e.POST("/user", api.CreateUser())
	e.DELETE("/user", api.DeleteUser())
	e.GET("/user/uname", api.FindUserByUsername())
	e.GET("/user/uid", api.FindUserById())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "hello world!!!")
	})

	defer client.Close()

	e.Logger.Fatal(e.Start(":2021"))

}
