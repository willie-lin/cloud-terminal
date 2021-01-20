package database

import (
	"context"
	"fmt"
	"github.com/facebook/ent/dialect/sql"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/willie-lin/cloud-terminal/pkg/database/ent"
	"github.com/willie-lin/cloud-terminal/pkg/database/ent/migrate"
	"log"
	"time"
)

//type Client struct {
//	Ent *ent.Client
//	CTX context.Context
//	Log *log.Logger
//}

//type Option

//func NewClient(opts ...Option) *Client {
//	cfg := config{log: log.Println, hooks: &hooks{}}
//	cfg.options(opts...)
//	client := &Client{config: cfg}
//	client.init()
//	return client
//}

func Client() (*ent.Client, error) {
	drv, err := sql.Open("mysql", "root:root1234@tcp(127.0.0.1:3306)/terminal?charset=utf8&parseTime=true")
	//drv, err := sql.Open("mysql", "root:root1234@tcp(127.0.0.1:3306)/ent")
	if err != nil {
		return nil, err
	}
	db := drv.DB()
	//db.SetMaxIdleConns(maxIdleConns)
	//db.SetConnMaxLifetime(connMaxLifetime)
	//db.SetMaxOpenConns(maxOpenConns)

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxLifetime(time.Hour)
	fmt.Println("ccc")
	return ent.NewClient(ent.Driver(drv)), nil

}

//func Clients() (*ent.Client, error) {
//	db, err := sql.Open(driverName, dataSourceName)
//	if err != nil {
//		return nil, err
//	}
//	//db := drv.DB()
//	db.SetMaxIdleConns(maxIdleConns)
//	db.SetConnMaxLifetime(connMaxLifetime)
//	db.SetMaxOpenConns(maxOpenConns)
//	drv := entsql.OpenDB(driverName, db)
//	return ent.NewClient(ent.Driver(drv)), nil
//}

func AutoMigration(client *ent.Client, ctx context.Context) {
	if err := client.Schema.Create(ctx); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
}

func DebugMode(err error, client *ent.Client, ctx context.Context) {
	err = client.Debug().Schema.Create(
		ctx,
		migrate.WithDropIndex(true),
		migrate.WithDropColumn(true),
	)
	if err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
}
