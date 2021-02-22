package handler

import (
	"context"
	"github.com/goccy/go-json"
	"github.com/labstack/echo/v4"
	"github.com/willie-lin/cloud-terminal/pkg/database/ent"
	"github.com/willie-lin/cloud-terminal/pkg/database/ent/group"

	"github.com/willie-lin/cloud-terminal/pkg/utils"
	"go.uber.org/zap"

	"net/http"
	"time"
)

// 创建用户组
func CreateGroup(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {

		g := new(ent.Group)

		log, _ := zap.NewDevelopment()
		if err := json.NewDecoder(c.Request().Body).Decode(&g); err != nil {
			log.Fatal("json decode error", zap.Error(err))
			return err
		}

		g.ID = utils.UUID()

		ugs, err := client.Group.Create().
			SetID(g.ID).
			SetName(g.Name).
			SetCreatedAt(time.Now()).
			SetUpdatedAt(time.Now()).Save(context.Background())
		if err != nil {
			log.Fatal("Create Group Error", zap.Error(err))
			return err
		}
		return c.JSON(http.StatusOK, &ugs)
	}
}

// 删除分组
func DeleteGroup(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {

		g := new(ent.Group)

		log, _ := zap.NewDevelopment()

		if err := json.NewDecoder(c.Request().Body).Decode(&g); err != nil {
			log.Fatal("Delete Group Error", zap.Error(err))
			return err
		}

		ugs, err := client.Group.Query().Where(group.NameEQ(g.Name)).Only(context.Background())
		if err != nil {
			log.Fatal("Query Group Error:", zap.Error(err))
			return err
		}
		err = client.Group.DeleteOne(ugs).Exec(context.Background())
		if err != nil {
			log.Fatal("Delete Group Error ", zap.Error(err))
			return err
		}
		return c.NoContent(http.StatusOK)
	}
}

// 根据ID删除用户组
func DeleteGroupById(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {

		g := new(ent.Group)

		log, _ := zap.NewDevelopment()

		if err := json.NewDecoder(c.Request().Body).Decode(&g); err != nil {
			log.Fatal("Delete Group Error", zap.Error(err))
			return err
		}

		gs, err := client.Group.Query().Where(group.NameEQ(g.Name)).Only(context.Background())
		if err != nil {
			log.Fatal("Query Group Error:", zap.Error(err))
			return err
		}
		err = client.Group.DeleteOneID(gs.ID).Exec(context.Background())
		if err != nil {
			log.Fatal("Delete Group Error ", zap.Error(err))
			return err
		}
		return c.NoContent(http.StatusOK)
	}
}

// 查寻所有用户组
func GetAllGroups(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {

		//client, err := database.Client()
		//client, err := config.NewClient()
		//if err != nil {
		//	panic(err)
		//}
		//user := new(ent.User)
		log, _ := zap.NewDevelopment()
		groups, err := client.Group.Query().All(context.Background())
		if err != nil {
			log.Fatal("GetAll Groups Error: ", zap.Error(err))
			return err
		}

		return c.JSON(http.StatusOK, groups)
	}
}
