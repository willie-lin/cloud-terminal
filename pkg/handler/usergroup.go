package handler

import (
	"context"
	"github.com/goccy/go-json"
	"github.com/labstack/echo/v4"
	"github.com/willie-lin/cloud-terminal/pkg/database/ent"
	"github.com/willie-lin/cloud-terminal/pkg/database/ent/usergroup"
	"github.com/willie-lin/cloud-terminal/pkg/utils"
	"go.uber.org/zap"
	"log"
	"net/http"
	"time"
)

// 创建用户组
func CreateUserGroup(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {

		ug := new(ent.UserGroup)

		log, _ := zap.NewDevelopment()
		if err := json.NewDecoder(c.Request().Body).Decode(&ug); err != nil {
			log.Fatal("json decode error", zap.Error(err))
			return err
		}

		ug.ID = utils.UUID()

		ugs, err := client.UserGroup.Create().
			SetID(ug.ID).
			SetName(ug.Name).
			SetMembers(ug.Members).
			SetCreatedAt(time.Now()).
			SetUpdatedAt(time.Now()).Save(context.Background())
		if err != nil {
			log.Fatal("Create UserGroup Error", zap.Error(err))
			return err
		}
		return c.JSON(http.StatusOK, &ugs)
	}
}

// 删除分组
func DeleteUserGroup(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {

		ug := new(ent.UserGroup)
		if err := json.NewDecoder(c.Request().Body).Decode(&ug); err != nil {
			log.Fatal("Delete UserGroup Error", zap.Error(err))
			return err
		}

		ugs, err := client.UserGroup.Query().Where(usergroup.NameEQ(ug.Name)).Only(context.Background())
		if err != nil {
			log.Fatal("Query UserGroup Error:", zap.Error(err))
			return err
		}
		err = client.UserGroup.DeleteOne(ugs).Exec(context.Background())
		if err != nil {
			log.Fatal("Delete UserGroup Error ", zap.Error(err))
			return err
		}
		return c.NoContent(http.StatusOK)
	}
}

// 根据ID删除用户组
func DeleteUserGroupById(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {

		ug := new(ent.UserGroup)
		if err := json.NewDecoder(c.Request().Body).Decode(&ug); err != nil {
			log.Fatal("Delete UserGroup Error", zap.Error(err))
			return err
		}

		ugs, err := client.UserGroup.Query().Where(usergroup.NameEQ(ug.Name)).Only(context.Background())
		if err != nil {
			log.Fatal("Query UserGroup Error:", zap.Error(err))
			return err
		}
		err = client.UserGroup.DeleteOneID(ugs.ID).Exec(context.Background())
		if err != nil {
			log.Fatal("Delete UserGroup Error ", zap.Error(err))
			return err
		}
		return c.NoContent(http.StatusOK)
	}
}
